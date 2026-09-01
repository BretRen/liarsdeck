package room

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"pdnode.com/play/liarsbar-web/internal/game"
	"pdnode.com/play/liarsbar-web/internal/model"
)

type Room struct {
	Code      string
	Clients   map[*Client]bool
	Game      *game.Game
	Hub       *Hub
	CreatedAt int64
	mu        sync.Mutex
}

func NewRoom(code string, hub *Hub, options ...any) *Room {
	return &Room{
		Code:      code,
		Clients:   make(map[*Client]bool),
		Game:      game.NewGame(code, options...),
		Hub:       hub,
		CreatedAt: time.Now().Unix(),
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[client] = true
}

func (r *Room) Broadcast() {
	r.Game.Lock()
	stateBase := map[string]any{
		"status":                 r.Game.State.Status,
		"current_turn":           r.Game.State.CurrentTurn,
		"table_card":             r.Game.State.TableCard,
		"last_player":            r.Game.State.LastPlayer,
		"last_played_cnt":        r.Game.State.LastPlayedCnt,
		"logs":                   r.Game.State.Logs,
		"deadline":               r.Game.State.Deadline,
		"pause_deadline":         r.Game.State.PauseDeadline,
		"paused_player":          r.Game.State.PausedPlayer,
		"remaining_turn_seconds": r.Game.State.RemainingTurnSeconds,
		"winner":                 r.Game.State.Winner,
		"room_code":              r.Game.State.RoomCode,
		"game_mode":              r.Game.State.GameMode,
		"max_players":            r.Game.State.MaxPlayers,
		"double_damage":          r.Game.State.DoubleDamage,
	}
	rawPlayers := make([]*model.Player, len(r.Game.State.Players))
	copy(rawPlayers, r.Game.State.Players)
	r.Game.Unlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.Clients {
		safePlayers := make([]model.SafePlayer, len(rawPlayers))
		for i, p := range rawPlayers {
			safePlayers[i] = p.ToSafe(client.ID)
		}

		state := make(map[string]any, len(stateBase)+1)
		for k, v := range stateBase {
			state[k] = v
		}
		state["players"] = safePlayers

		payload, err := json.Marshal(map[string]any{
			"type": "game_state",
			"data": state,
		})
		if err != nil {
			continue
		}

		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (r *Room) BroadcastEvent(eventType string, data any) {
	payload, err := json.Marshal(model.EventMessage{
		Type: eventType,
		Data: data,
	})
	if err != nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.Clients {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	delete(r.Clients, client)
	r.mu.Unlock()

	r.Game.Lock()
	removedIdx := -1
	for i, p := range r.Game.State.Players {
		if p.ID == client.ID {
			removedIdx = i
			break
		}
	}

	if removedIdx != -1 {
		p := r.Game.State.Players[removedIdx]
		p.ClientRef = nil

		if r.Game.State.Status == model.StatusWaiting || r.Game.State.Status == model.StatusGameOver || p.IsSpectator {
			r.Game.Log(fmt.Sprintf("👋 %s 离开了房间 / 👋 %s left the room", p.Nickname, p.Nickname))
			wasHost := p.IsHost
			r.Game.State.Players = append(r.Game.State.Players[:removedIdx], r.Game.State.Players[removedIdx+1:]...)
			if wasHost && len(r.Game.State.Players) > 0 {
				for _, pp := range r.Game.State.Players {
					if !pp.IsSpectator {
						pp.IsHost = true
						r.Game.Log(fmt.Sprintf("👑 %s 成为新房主 / 👑 %s is the new host", pp.Nickname, pp.Nickname))
						break
					}
				}
			}
		} else if r.Game.State.Status == model.StatusPlaying || r.Game.State.Status == model.StatusPaused {
			if p.IsAlive && !p.IsSpectator {
				if p.DisconnectGraceRemaining > 0 {
					// 消耗断线保护时间池
					r.Game.PauseGame(p)
				} else {
					// 保护时间已耗尽：直接判定淘汰出局
					p.IsAlive = false
					r.Game.Log(fmt.Sprintf("💀 玩家 %s 断线保护时间已耗尽，已被判定淘汰出局！ / 💀 %s has no grace time left and was eliminated!", p.Nickname, p.Nickname))

					if r.Game.State.Status == model.StatusPaused && (r.Game.State.PausedPlayerID == p.ID || (r.Game.State.PausedPlayerID == "" && r.Game.State.PausedPlayer == p.Nickname)) {
						r.Game.State.Status = model.StatusPlaying
						r.Game.State.PausedPlayer = ""
						r.Game.State.PausedPlayerID = ""
						r.Game.State.PauseDeadline = 0
					}

					r.Game.AdvanceToAlive()

					aliveCount := 0
					var lastAlive *model.Player
					for _, pp := range r.Game.State.Players {
						if !pp.IsSpectator && pp.IsAlive {
							aliveCount++
							lastAlive = pp
						}
					}
					if aliveCount <= 1 && lastAlive != nil {
						r.Game.State.Status = model.StatusGameOver
						r.Game.State.Winner = lastAlive.Nickname
						r.Game.Log(fmt.Sprintf("🏆 %s 获胜！ / 🏆 %s wins!", lastAlive.Nickname, lastAlive.Nickname))
					} else {
						r.Game.ResetTimer()
					}
				}
			}
		}
	}
	r.Game.Unlock()

	client.Close()
	r.Broadcast()
}

func (r *Room) HandleClientMessage(c *Client, msg model.WSMessage) {
	if msg.Action == "ping" {
		var req model.PingPayload
		_ = json.Unmarshal(msg.Payload, &req)
		respPayload, _ := json.Marshal(map[string]any{
			"type": "pong",
			"data": map[string]any{
				"client_time": req.ClientTime,
				"server_time": time.Now().UnixMilli(),
			},
		})
		select {
		case c.Send <- respPayload:
		default:
		}
		return
	}

	g := r.Game
	g.Lock()

	switch msg.Action {
	case "ready":
		if g.State.Status == model.StatusWaiting {
			for _, p := range g.State.Players {
				if p.ID == c.ID && !p.IsSpectator {
					p.IsReady = !p.IsReady
					status := "未准备"
					en := "is not ready"
					if p.IsReady {
						status = "已准备"
						en = "is ready"
					}
					g.Log(fmt.Sprintf("✅ %s %s / ✅ %s %s", p.Nickname, status, p.Nickname, en))
					break
				}
			}
		}

	case "remove_player":
		if g.State.Status == model.StatusWaiting || g.State.Status == model.StatusGameOver {
			var req model.RemovePlayerPayload
			_ = json.Unmarshal(msg.Payload, &req)

			callerIsHost := false
			for _, p := range g.State.Players {
				if p.ID == c.ID && p.IsHost {
					callerIsHost = true
					break
				}
			}
			if callerIsHost {
				kickIdx := -1
				for i, p := range g.State.Players {
					if p.ID == req.TargetID {
						kickIdx = i
						break
					}
				}
				if kickIdx != -1 {
					kickedPlayer := g.State.Players[kickIdx]
					g.Log(fmt.Sprintf("👢 %s 被房主移出房间 / 👢 %s was removed by host", kickedPlayer.Nickname, kickedPlayer.Nickname))
					if targetClient, ok := kickedPlayer.ClientRef.(*Client); ok && targetClient != nil {
						kickedPlayer.ClientRef = nil
						targetClient.Close()
					}
					g.State.Players = append(g.State.Players[:kickIdx], g.State.Players[kickIdx+1:]...)
				}
			}
		}

	case "start":
		if g.State.Status == model.StatusWaiting {
			callerIsHost := false
			for _, p := range g.State.Players {
				if p.ID == c.ID && p.IsHost {
					callerIsHost = true
					break
				}
			}
			if !callerIsHost {
				g.Unlock()
				return
			}

			allReady := true
			playerCount := 0
			for _, p := range g.State.Players {
				if !p.IsSpectator {
					playerCount++
					if !p.IsReady || p.ClientRef == nil {
						allReady = false
					}
				}
			}

			if callerIsHost && playerCount >= 2 && allReady {
				g.State.Status = model.StatusPlaying
				g.StartRound()
			} else if !allReady {
				g.Log("⏳ 等待所有玩家准备 / ⏳ Waiting for all players to ready up")
			}
		}

	case "reset":
		if g.State.Status == model.StatusGameOver {
			callerIsHost := false
			for _, p := range g.State.Players {
				if p.ID == c.ID && p.IsHost {
					callerIsHost = true
					break
				}
			}
			if callerIsHost {
				g.ResetGame()
			}
		}

	case "play", "play_cards":
		if g.State.Status == model.StatusPlaying && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currPlayer := g.State.Players[g.State.CurrentTurn]
			if currPlayer.ID == c.ID && currPlayer.IsAlive && !currPlayer.IsSpectator {
				var req model.PlayCardsPayload
				_ = json.Unmarshal(msg.Payload, &req)
				_ = g.PlayCards(g.State.CurrentTurn, req.Cards)
			}
		}

	case "call_liar":
		if g.State.Status == model.StatusPlaying && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currPlayer := g.State.Players[g.State.CurrentTurn]
			if currPlayer.ID == c.ID && currPlayer.IsAlive && !currPlayer.IsSpectator && g.State.LastPlayer != -1 && g.State.LastPlayer != g.State.CurrentTurn {
				g.CallLiar(
					g.State.CurrentTurn,
					g.State.LastPlayer,
					func(caller, accused string) {
						r.BroadcastEvent("liar_call", model.LiarCallEvent{Caller: caller, Accused: accused})
					},
					func(caller, accused string, cards []model.Card) {
						r.BroadcastEvent("reveal", model.RevealEvent{Caller: caller, Accused: accused, Cards: cards})
					},
					func(target string, fatal bool, doubleShot bool, armorBlocked bool) {
						r.BroadcastEvent("shot", model.ShotEvent{
							Target:       target,
							Fatal:        fatal,
							DoubleShot:   doubleShot,
							ArmorBlocked: armorBlocked,
						})
					},
				)
			}
		}

	case "use_item":
		if g.State.Status == model.StatusPlaying && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currPlayer := g.State.Players[g.State.CurrentTurn]
			if currPlayer.ID == c.ID && currPlayer.IsAlive && !currPlayer.IsSpectator {
				var req model.UseItemPayload
				if err := json.Unmarshal(msg.Payload, &req); err == nil {
					res, err := g.UseItem(g.State.CurrentTurn, req.Item)
					if err == nil {
						r.BroadcastEvent("item_used", model.ItemUsedEvent{
							PlayerID:   c.ID,
							Nickname:   c.Nickname,
							Item:       req.Item,
							ItemName:   string(req.Item),
							TargetCard: g.State.TableCard,
						})
						// 鹰眼透镜私密下发给调用者
						if inspected, ok := res["inspected_card"].(model.Card); ok {
							privData, _ := json.Marshal(map[string]any{
								"type": "eagle_eye_result",
								"data": map[string]any{
									"card": inspected,
								},
							})
							select {
							case c.Send <- privData:
							default:
							}
						}
					} else {
						errData, _ := json.Marshal(map[string]string{
							"error": err.Error(),
						})
						select {
						case c.Send <- errData:
						default:
						}
					}
				}
			}
		}
	}

	g.Unlock()
	r.Broadcast()
}

func (r *Room) Watchdog() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	emptySeconds := 0

	for range ticker.C {
		r.mu.Lock()
		clientCount := len(r.Clients)
		r.mu.Unlock()

		if clientCount == 0 {
			emptySeconds++
			if emptySeconds >= 300 { // 房间空置超过5分钟自动注销
				if r.Hub != nil {
					r.Hub.RemoveRoom(r.Code)
				}
				return
			}
		} else {
			emptySeconds = 0
		}

		r.Game.Lock()

		// 1. 处理 30 秒断线暂停超时
		if r.Game.State.Status == model.StatusPaused {
			if time.Now().Unix() >= r.Game.State.PauseDeadline {
				r.Game.HandlePauseTimeout()
				r.Game.Unlock()
				r.Broadcast()
				continue
			}
			r.Game.Unlock()
			continue
		}

		// 2. 正常游戏中回合操作超时
		if r.Game.State.Status == model.StatusPlaying && time.Now().Unix() > r.Game.State.Deadline {
			currIdx := r.Game.State.CurrentTurn
			if currIdx < 0 || currIdx >= len(r.Game.State.Players) {
				r.Game.Unlock()
				continue
			}

			p := r.Game.State.Players[currIdx]
			if p.IsSpectator {
				r.Game.NextTurn()
				r.Game.Unlock()
				r.Broadcast()
				continue
			}

			r.Game.Log(fmt.Sprintf("⏱️ %s 操作超时！系统代管... / ⏱️ %s timed out! Auto-playing...", p.Nickname, p.Nickname))

			if r.Game.State.LastPlayer == -1 || r.Game.State.LastPlayer == currIdx {
				// 本轮首次出牌或同玩家连续操作：必须强制出牌
				if len(p.Hand) == 0 {
					r.Game.Log(fmt.Sprintf("%s 没有手牌，跳过 / %s has no cards, skipping", p.Nickname, p.Nickname))
					r.Game.NextTurn()
				} else {
					card := p.Hand[0]
					p.Hand = p.Hand[1:]
					r.Game.HiddenCards = []model.Card{card}
					r.Game.State.LastPlayedCnt = 1
					r.Game.State.LastPlayer = currIdx
					r.Game.Log(fmt.Sprintf("%s 强制打出了 1 张牌 / %s auto-played 1 card", p.Nickname, p.Nickname))
					if len(p.Hand) == 0 {
						r.Game.Log(fmt.Sprintf("%s 打完了所有手牌！下家必须质疑 / %s emptied their hand! Next player must call liar", p.Nickname, p.Nickname))
					}
					r.Game.NextTurn()
				}
			} else {
				if p.ClientRef == nil {
					r.Game.Log(fmt.Sprintf("%s 已断线，跳过 / %s disconnected, skipping", p.Nickname, p.Nickname))
					r.Game.NextTurn()
				} else {
					r.Game.CallLiar(
						currIdx,
						r.Game.State.LastPlayer,
						func(caller, accused string) {
							r.BroadcastEvent("liar_call", model.LiarCallEvent{Caller: caller, Accused: accused})
						},
						func(caller, accused string, cards []model.Card) {
							r.BroadcastEvent("reveal", model.RevealEvent{Caller: caller, Accused: accused, Cards: cards})
						},
						func(target string, fatal bool, doubleShot bool, armorBlocked bool) {
							r.BroadcastEvent("shot", model.ShotEvent{
								Target:       target,
								Fatal:        fatal,
								DoubleShot:   doubleShot,
								ArmorBlocked: armorBlocked,
							})
						},
					)
				}
			}
			r.Game.Unlock()
			r.Broadcast()
			continue
		}

		r.Game.Unlock()
	}
}
