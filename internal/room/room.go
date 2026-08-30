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
	Hub     *Hub
	ID      string
	Code    string
	Clients map[*Client]bool
	Game    *game.Game
	mu      sync.Mutex
}

func NewRoom(hub *Hub, code string) *Room {
	return &Room{
		Hub:     hub,
		ID:      code,
		Code:    code,
		Clients: make(map[*Client]bool),
		Game:    game.NewGame(code),
	}
}

func (r *Room) Lock()   { r.mu.Lock() }
func (r *Room) Unlock() { r.mu.Unlock() }

func (r *Room) AddClient(c *Client) {
	r.mu.Lock()
	r.Clients[c] = true
	r.mu.Unlock()
}

func (r *Room) Broadcast() {
	r.Game.Lock()
	safePlayers := make([]model.SafePlayer, len(r.Game.State.Players))
	for i, p := range r.Game.State.Players {
		safePlayers[i] = p.ToSafe()
	}

	state := map[string]any{
		"status":          r.Game.State.Status,
		"players":         safePlayers,
		"current_turn":    r.Game.State.CurrentTurn,
		"table_card":      r.Game.State.TableCard,
		"last_player":     r.Game.State.LastPlayer,
		"last_played_cnt": r.Game.State.LastPlayedCnt,
		"logs":            r.Game.State.Logs,
		"deadline":        r.Game.State.Deadline,
		"winner":          r.Game.State.Winner,
		"room_code":       r.Game.State.RoomCode,
	}
	r.Game.Unlock()

	payload, err := json.Marshal(map[string]any{
		"type": "game_state",
		"data": state,
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

		if r.Game.State.Status == model.StatusWaiting || p.IsSpectator {
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
		} else {
			// 游戏中：保留席位以便断线重连
			r.Game.Log(fmt.Sprintf("📴 %s 断线了 / 📴 %s disconnected", p.Nickname, p.Nickname))
		}
	}
	r.Game.Unlock()

	client.Close()
	r.Broadcast()
}

func (r *Room) HandleClientMessage(c *Client, msg model.WSMessage) {
	g := r.Game
	g.Lock()

	switch msg.Action {
	case "ready":
		for _, p := range g.State.Players {
			if p.ID == c.ID {
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

	case "remove_player":
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
			for _, p := range g.State.Players {
				if p.ID == req.TargetID && p.ClientRef != nil {
					g.Log(fmt.Sprintf("👢 %s 被房主移出房间 / 👢 %s was removed by host", p.Nickname, p.Nickname))
					if targetClient, ok := p.ClientRef.(*Client); ok && targetClient != nil {
						targetClient.Close()
					}
					break
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
					if !p.IsReady {
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

	case "play_cards":
		if g.State.Status == model.StatusPlaying && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currPlayer := g.State.Players[g.State.CurrentTurn]
			if currPlayer.ID == c.ID {
				var req model.PlayCardsPayload
				_ = json.Unmarshal(msg.Payload, &req)
				_ = g.PlayCards(g.State.CurrentTurn, req.Cards)
			}
		}

	case "call_liar":
		if g.State.Status == model.StatusPlaying && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currPlayer := g.State.Players[g.State.CurrentTurn]
			if currPlayer.ID == c.ID && g.State.LastPlayer != -1 && g.State.LastPlayer != g.State.CurrentTurn {
				g.CallLiar(
					g.State.CurrentTurn,
					g.State.LastPlayer,
					func(caller, accused string) {
						r.BroadcastEvent("liar_call", model.LiarCallEvent{Caller: caller, Accused: accused})
					},
					func(caller, accused string, cards []model.Card) {
						r.BroadcastEvent("reveal", model.RevealEvent{Caller: caller, Accused: accused, Cards: cards})
					},
					func(target string, fatal bool) {
						r.BroadcastEvent("shot", model.ShotEvent{Target: target, Fatal: fatal})
					},
				)
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
						func(target string, fatal bool) {
							r.BroadcastEvent("shot", model.ShotEvent{Target: target, Fatal: fatal})
						},
					)
				}
			}
		}
		r.Game.Unlock()
		r.Broadcast()
	}
}
