package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Card string

const (
	King  Card = "K"
	Queen Card = "Q"
	Ace   Card = "A"
	Two   Card = "2"
)

type Player struct {
	ID          string  `json:"id"`
	Nickname    string  `json:"nickname"`
	Hand        []Card  `json:"hand"`
	Revolver    []string `json:"-"`
	Bullets     int     `json:"bullets"`
	IsAlive     bool    `json:"is_alive"`
	IsHost      bool    `json:"is_host"`
	IsSpectator bool    `json:"is_spectator"`
	Client      *Client `json:"-"`
}

type GameState struct {
	Status        string    `json:"status"`
	Players       []*Player `json:"players"`
	CurrentTurn   int       `json:"current_turn"`
	TableCard     Card      `json:"table_card"`
	LastPlayer    int       `json:"last_player"`
	LastPlayedCnt int       `json:"last_played_cnt"`
	Logs          []string  `json:"logs"`
	Deadline      int64     `json:"deadline"`
	Winner        string    `json:"winner,omitempty"`
	RoomCode      string    `json:"room_code"`
}

type Game struct {
	State       *GameState
	HiddenCards []Card
	mu          sync.Mutex
}

type Room struct {
	Hub     *Hub
	ID      string
	Code    string
	Clients map[*Client]bool
	Game    *Game
	mu      sync.Mutex
}

type Client struct {
	ID       string
	Nickname string
	Room     *Room
	Conn     *websocket.Conn
	Send     chan []byte
}

type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	hub      = &Hub{Rooms: make(map[string]*Room)}
)

type Hub struct {
	Rooms map[string]*Room
	mu    sync.Mutex
}

func randomCode(n int) string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewGame(code string) *Game {
	return &Game{
		State: &GameState{
			Status:      "waiting",
			Players:     make([]*Player, 0),
			CurrentTurn: -1,
			LastPlayer:  -1,
			RoomCode:    code,
		},
	}
}

func (g *Game) Log(msg string) {
	g.State.Logs = append(g.State.Logs, time.Now().Format("15:04:05")+" "+msg)
	if len(g.State.Logs) > 10 {
		g.State.Logs = g.State.Logs[len(g.State.Logs)-10:]
	}
}

func (g *Game) StartRound() {
	deck := []Card{King, King, King, King, King, King, Queen, Queen, Queen, Queen, Queen, Queen, Ace, Ace, Ace, Ace, Ace, Ace, Two, Two, Two, Two, Two, Two}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	g.State.TableCard = deck[0]
	deck = deck[1:]
	for g.State.TableCard == Two && len(deck) > 0 {
		g.State.TableCard = deck[0]
		deck = deck[1:]
	}
	if g.State.TableCard == Two {
		g.State.TableCard = King
	}
	g.Log(fmt.Sprintf("新的一轮开始！本轮真牌是: %s", g.State.TableCard))

	aliveCount := 0
	for _, p := range g.State.Players {
		if !p.IsSpectator && p.IsAlive {
			p.Hand = deck[:5]
			deck = deck[5:]
			aliveCount++
		}
	}

	if aliveCount <= 1 {
		g.State.Status = "game_over"
		for _, p := range g.State.Players {
			if !p.IsSpectator && p.IsAlive {
				g.State.Winner = p.Nickname
			}
		}
		g.Log("游戏结束！")
		return
	}

	g.HiddenCards = []Card{}
	g.State.LastPlayedCnt = 0
	g.State.LastPlayer = -1

	if g.State.CurrentTurn == -1 {
		g.State.CurrentTurn = rand.Intn(len(g.State.Players))
	}
	g.advanceToAlive()
	g.ResetTimer()
}

func (g *Game) advanceToAlive() {
	n := len(g.State.Players)
	if n == 0 {
		return
	}
	start := g.State.CurrentTurn
	for i := 0; i < n; i++ {
		idx := (start + i) % n
		if g.State.Players[idx].IsAlive && !g.State.Players[idx].IsSpectator {
			g.State.CurrentTurn = idx
			return
		}
	}
	// no alive players
	g.State.Status = "game_over"
	g.State.Winner = "无人存活"
	g.Log("游戏结束！所有玩家都已淘汰")
}

func (g *Game) ResetTimer() {
	g.State.Deadline = time.Now().Add(30 * time.Second).Unix()
}

func (g *Game) NextTurn() {
	start := g.State.CurrentTurn
	n := len(g.State.Players)
	for i := 1; i <= n; i++ {
		idx := (start + i) % n
		if g.State.Players[idx].IsAlive && !g.State.Players[idx].IsSpectator {
			g.State.CurrentTurn = idx
			g.ResetTimer()
			return
		}
	}
	g.State.Status = "game_over"
	g.State.Winner = "无人存活"
	g.Log("游戏结束！所有玩家都已淘汰")
}

func safeBroadcastShot(room *Room, targetNickname string, fatal bool) {
	room.BroadcastEvent("shot", map[string]interface{}{
		"target": targetNickname,
		"fatal":  fatal,
	})
}

func (g *Game) FireGun(playerIdx int) {
	p := g.State.Players[playerIdx]

	if len(p.Revolver) == 0 {
		revolver := []string{"Blank", "Blank", "Blank", "Blank", "Blank", "Fatal"}
		rand.Shuffle(len(revolver), func(i, j int) { revolver[i], revolver[j] = revolver[j], revolver[i] })
		p.Revolver = revolver
		p.Bullets = 6
	}

	bullet := p.Revolver[0]
	p.Revolver = p.Revolver[1:]
	p.Bullets = len(p.Revolver)

	if p.Client != nil {
		safeBroadcastShot(p.Client.Room, p.Nickname, bullet == "Fatal")
	}

	if bullet == "Fatal" {
		p.IsAlive = false
		g.Log(fmt.Sprintf("💥 砰！%s 抽中致命子弹，被淘汰出局！", p.Nickname))
	} else {
		g.Log(fmt.Sprintf("💨 咔哒。%s 抽中空包弹，逃过一劫。", p.Nickname))
	}

	g.State.CurrentTurn = playerIdx
	g.advanceToAlive()
	if g.State.Status == "game_over" {
		return
	}

	// check if only one alive player remains
	aliveCount := 0
	var lastAlive *Player
	for _, pp := range g.State.Players {
		if !pp.IsSpectator && pp.IsAlive {
			aliveCount++
			lastAlive = pp
		}
	}
	if aliveCount <= 1 && lastAlive != nil {
		g.State.Status = "game_over"
		g.State.Winner = lastAlive.Nickname
		g.Log(fmt.Sprintf("🏆 %s 获胜！", lastAlive.Nickname))
		return
	}

	g.StartRound()
}

func (r *Room) Broadcast() {
	r.Game.mu.Lock()
	// build a safe copy without private fields
	type safePlayer struct {
		ID          string `json:"id"`
		Nickname    string `json:"nickname"`
		Hand        []Card `json:"hand"`
		Bullets     int    `json:"bullets"`
		IsAlive     bool   `json:"is_alive"`
		IsHost      bool   `json:"is_host"`
		IsSpectator bool   `json:"is_spectator"`
	}
	players := make([]safePlayer, len(r.Game.State.Players))
	for i, p := range r.Game.State.Players {
		players[i] = safePlayer{
			ID:          p.ID,
			Nickname:    p.Nickname,
			Hand:        p.Hand,
			Bullets:     p.Bullets,
			IsAlive:     p.IsAlive,
			IsHost:      p.IsHost,
			IsSpectator: p.IsSpectator,
		}
	}
	state := map[string]interface{}{
		"status":         r.Game.State.Status,
		"players":        players,
		"current_turn":   r.Game.State.CurrentTurn,
		"table_card":     r.Game.State.TableCard,
		"last_player":    r.Game.State.LastPlayer,
		"last_played_cnt": r.Game.State.LastPlayedCnt,
		"logs":           r.Game.State.Logs,
		"deadline":       r.Game.State.Deadline,
		"winner":         r.Game.State.Winner,
		"room_code":      r.Game.State.RoomCode,
	}
	b, _ := json.Marshal(map[string]interface{}{
		"type": "game_state",
		"data": state,
	})
	r.Game.mu.Unlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.Clients {
		select {
		case client.Send <- b:
		default:
		}
	}
}

func (g *Game) CallLiar(callerIdx, accusedIdx int) {
	caller := g.State.Players[callerIdx]
	accused := g.State.Players[accusedIdx]
	g.Log(fmt.Sprintf("🚨 %s 质疑 %s 说谎！", caller.Nickname, accused.Nickname))

	r := caller.Client.Room
	r.BroadcastEvent("liar_call", map[string]interface{}{
		"caller": caller.Nickname, "accused": accused.Nickname,
	})
	r.BroadcastEvent("reveal", map[string]interface{}{
		"caller":  caller.Nickname,
		"accused": accused.Nickname,
		"cards":   g.HiddenCards,
	})

	isLiar := false
	revealMsg := fmt.Sprintf("%s 的底牌是: ", accused.Nickname)
	for _, c := range g.HiddenCards {
		revealMsg += string(c) + " "
		if c != g.State.TableCard && c != Two {
			isLiar = true
		}
	}
	g.Log(revealMsg)

	if isLiar {
		g.Log("👉 质疑成功！出牌者说谎！")
		g.FireGun(accusedIdx)
	} else {
		g.Log("❌ 质疑失败！出牌者是清白的！")
		g.FireGun(callerIdx)
	}
}

func (r *Room) BroadcastEvent(eventType string, data interface{}) {
	b, _ := json.Marshal(map[string]interface{}{
		"type": eventType,
		"data": data,
	})
	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.Clients {
		select {
		case client.Send <- b:
		default:
		}
	}
}

func (r *Room) getPlayersSummary() []map[string]interface{} {
	r.Game.mu.Lock()
	defer r.Game.mu.Unlock()
	out := make([]map[string]interface{}, 0)
	for _, p := range r.Game.State.Players {
		out = append(out, map[string]interface{}{
			"id":           p.ID,
			"nickname":     p.Nickname,
			"is_host":      p.IsHost,
			"is_spectator": p.IsSpectator,
			"is_alive":     p.IsAlive,
			"bullets":      p.Bullets,
			"hand_count":   len(p.Hand),
		})
	}
	return out
}

func main() {
	e := echo.New()
	e.Static("/", "public")

	e.GET("/ws", func(c echo.Context) error {
		action := c.QueryParam("action")
		code := c.QueryParam("code")
		nickname := c.QueryParam("name")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		hub.mu.Lock()

		if action == "create" {
			// generate unique room code
			for {
				code = randomCode(6)
				if _, exists := hub.Rooms[code]; !exists {
					break
				}
			}
			room := &Room{ID: code, Code: code, Clients: make(map[*Client]bool), Game: NewGame(code)}
			hub.Rooms[code] = room
			go room.Watchdog()
			hub.mu.Unlock()

			client := &Client{ID: fmt.Sprintf("%d", rand.Int()), Nickname: nickname, Room: room, Conn: conn, Send: make(chan []byte, 256)}
			room.mu.Lock()
			room.Clients[client] = true
			room.mu.Unlock()

			room.Game.mu.Lock()
			room.Game.State.Players = append(room.Game.State.Players, &Player{
				ID: client.ID, Nickname: nickname, Hand: []Card{},
				Revolver: []string{}, Bullets: 0, IsAlive: true, IsHost: true, IsSpectator: false, Client: client,
			})
			room.Game.Log(fmt.Sprintf("🏠 %s 创建了房间 %s", nickname, code))
			room.Game.mu.Unlock()

			go client.WritePump()
			go client.ReadPump()
			room.Broadcast()
			return nil
		}

		room, exists := hub.Rooms[code]
		if !exists {
			hub.mu.Unlock()
			conn.WriteJSON(map[string]string{"error": "房间不存在"})
			conn.Close()
			return nil
		}
		hub.mu.Unlock()

		client := &Client{ID: fmt.Sprintf("%d", rand.Int()), Nickname: nickname, Room: room, Conn: conn, Send: make(chan []byte, 256)}
		room.mu.Lock()
		room.Clients[client] = true
		room.mu.Unlock()

		room.Game.mu.Lock()
		isSpectator := action == "spectate"
		if !isSpectator && room.Game.State.Status == "waiting" && len(room.Game.State.Players) < 4 {
			room.Game.State.Players = append(room.Game.State.Players, &Player{
				ID: client.ID, Nickname: nickname, Hand: []Card{},
				Revolver: []string{}, Bullets: 0, IsAlive: true, IsHost: false, IsSpectator: false, Client: client,
			})
			room.Game.Log(nickname + " 加入了房间")
		} else if isSpectator {
			room.Game.State.Players = append(room.Game.State.Players, &Player{
				ID: client.ID, Nickname: nickname, Hand: []Card{},
				Revolver: []string{}, Bullets: 0, IsAlive: true, IsHost: false, IsSpectator: true, Client: client,
			})
			room.Game.Log(fmt.Sprintf("👀 %s 以观众身份加入", nickname))
		} else if !isSpectator {
			room.Game.State.Players = append(room.Game.State.Players, &Player{
				ID: client.ID, Nickname: nickname, Hand: []Card{},
				Revolver: []string{}, Bullets: 0, IsAlive: true, IsHost: false, IsSpectator: false, Client: client,
			})
			room.Game.Log(nickname + " 加入了游戏（观战）")
		}
		room.Game.mu.Unlock()

		go client.WritePump()
		go client.ReadPump()
		room.Broadcast()
		return nil
	})

	e.Logger.Fatal(e.Start(":8095"))
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	delete(r.Clients, client)
	r.mu.Unlock()

	r.Game.mu.Lock()
	removedIdx := -1
	for i, p := range r.Game.State.Players {
		if p.ID == client.ID {
			removedIdx = i
			p.IsAlive = false
			break
		}
	}

	if removedIdx != -1 {
		p := r.Game.State.Players[removedIdx]
		p.Client = nil
		r.Game.Log(fmt.Sprintf("👋 %s 离开了房间", p.Nickname))

		if r.Game.State.Status == "waiting" {
			wasHost := p.IsHost
			r.Game.State.Players = append(r.Game.State.Players[:removedIdx], r.Game.State.Players[removedIdx+1:]...)
			// transfer host if host left
			if wasHost && len(r.Game.State.Players) > 0 {
				for _, pp := range r.Game.State.Players {
					if !pp.IsSpectator {
						pp.IsHost = true
						r.Game.Log(fmt.Sprintf("👑 %s 成为新房主", pp.Nickname))
						break
					}
				}
			}
		} else {
			// check if host left during game
			wasHost := p.IsHost
			if wasHost {
				for _, pp := range r.Game.State.Players {
					if pp.Client != nil && !pp.IsSpectator && pp != p {
						pp.IsHost = true
						r.Game.Log(fmt.Sprintf("👑 %s 成为新房主", pp.Nickname))
						break
					}
				}
			}

			aliveCount := 0
			for _, pp := range r.Game.State.Players {
				if !pp.IsSpectator && pp.IsAlive {
					aliveCount++
				}
			}
			if aliveCount <= 1 {
				r.Game.State.Status = "game_over"
				for _, pp := range r.Game.State.Players {
					if !pp.IsSpectator && pp.IsAlive {
						r.Game.State.Winner = pp.Nickname
					}
				}
				r.Game.Log("游戏结束！其他玩家已离开")
			} else if r.Game.State.CurrentTurn == removedIdx {
				r.Game.NextTurn()
			}
		}
	}
	r.Game.mu.Unlock()
	close(client.Send)
	client.Conn.Close()
	r.Broadcast()
}

func (c *Client) ReadPump() {
	defer func() { c.Room.RemoveClient(c) }()
	for {
		_, text, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		json.Unmarshal(text, &msg)

		g := c.Room.Game
		g.mu.Lock()

		// ---- Room management actions ----
		if msg.Action == "remove_player" {
			var req struct {
				TargetID string `json:"target_id"`
			}
			json.Unmarshal(msg.Payload, &req)

			// find the caller and verify they're host
			var callerHost bool
			for _, p := range g.State.Players {
				if p.ID == c.ID && p.IsHost {
					callerHost = true
					break
				}
			}
			if callerHost {
				for i, p := range g.State.Players {
					if p.ID == req.TargetID && p.Client != nil {
						g.Log(fmt.Sprintf("👢 %s 被房主移出房间", p.Nickname))
						// disconnect the target client
						targetClient := p.Client
						p.Client = nil

						// broadcast removal first
						g.mu.Unlock()
						g.State.Players = append(g.State.Players[:i], g.State.Players[i+1:]...)
						targetClient.Room.RemoveClient(targetClient)
						g.mu.Lock()
						break
					}
				}
			}
			g.mu.Unlock()
			continue
		}

		// ---- Game actions ----
		if msg.Action == "start" && g.State.Status == "waiting" {
			// check that caller is host
			isHost := false
			for _, p := range g.State.Players {
				if p.ID == c.ID && p.IsHost {
					isHost = true
					break
				}
			}
			playerCount := 0
			for _, p := range g.State.Players {
				if !p.IsSpectator {
					playerCount++
				}
			}
			if isHost && playerCount >= 2 {
				g.State.Status = "playing"
				g.StartRound()
			}
		}

		// reset
		if msg.Action == "reset" && g.State.Status == "game_over" {
			alive := []*Player{}
			for _, p := range g.State.Players {
				if p.Client == nil {
					continue
				}
				p.IsAlive = true
				p.Hand = []Card{}
				p.Bullets = 0
				p.Revolver = []string{}
				alive = append(alive, p)
			}
			g.State.Players = alive
			g.State.Winner = ""
			g.State.Logs = []string{}
			g.State.CurrentTurn = -1
			g.State.LastPlayer = -1
			g.State.LastPlayedCnt = 0

			// ensure first non-spectator is host
			for _, p := range g.State.Players {
				p.IsHost = false
			}
			for _, p := range g.State.Players {
				if !p.IsSpectator {
					p.IsHost = true
					break
				}
			}

			playerCount := 0
			for _, p := range g.State.Players {
				if !p.IsSpectator {
					playerCount++
				}
			}
			if playerCount < 2 {
				g.State.Status = "waiting"
				g.Log("等待更多玩家加入...")
			} else {
				g.State.Status = "playing"
				g.StartRound()
			}
		}

		// current player actions
		if g.State.Status == "playing" && len(g.State.Players) > g.State.CurrentTurn && g.State.CurrentTurn >= 0 {
			currentPlayer := g.State.Players[g.State.CurrentTurn]
			if currentPlayer.ID == c.ID {
				if msg.Action == "play_cards" {
					var req struct {
						Cards []Card `json:"cards"`
					}
					json.Unmarshal(msg.Payload, &req)

					if len(req.Cards) >= 1 && len(req.Cards) <= 3 {
						p := g.State.Players[g.State.CurrentTurn]
						newHand := []Card{}
						used := make([]bool, len(req.Cards))
						for _, hc := range p.Hand {
							removed := false
							for i, rc := range req.Cards {
								if !used[i] && hc == rc {
									used[i] = true
									removed = true
									break
								}
							}
							if !removed {
								newHand = append(newHand, hc)
							}
						}
						p.Hand = newHand
						g.HiddenCards = req.Cards
						g.State.LastPlayedCnt = len(req.Cards)
						g.State.LastPlayer = g.State.CurrentTurn

						g.Log(fmt.Sprintf("%s 宣称打出了 %d 张牌", p.Nickname, len(req.Cards)))

						// check for winner: player emptied their hand
						if len(p.Hand) == 0 {
							g.State.Status = "game_over"
							g.State.Winner = p.Nickname
							g.Log(fmt.Sprintf("🏆 %s 打完了所有手牌，获胜！", p.Nickname))
							g.mu.Unlock()
							c.Room.Broadcast()
							continue
						}

						g.NextTurn()
					}
				}

				if msg.Action == "call_liar" && g.State.LastPlayer != -1 && g.State.LastPlayer != g.State.CurrentTurn {
					g.CallLiar(g.State.CurrentTurn, g.State.LastPlayer)
				}
			}
		}

		g.mu.Unlock()
		c.Room.Broadcast()
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (r *Room) Watchdog() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		r.Game.mu.Lock()
		if r.Game.State.Status == "playing" && time.Now().Unix() > r.Game.State.Deadline {
			currIdx := r.Game.State.CurrentTurn
			if currIdx < 0 || currIdx >= len(r.Game.State.Players) {
				r.Game.mu.Unlock()
				continue
			}
			p := r.Game.State.Players[currIdx]
			if p.IsSpectator {
				r.Game.NextTurn()
				r.Game.mu.Unlock()
				r.Broadcast()
				continue
			}
			r.Game.Log(fmt.Sprintf("⏱️ %s 操作超时！系统代管...", p.Nickname))
			if r.Game.State.LastPlayer == -1 || r.Game.State.LastPlayer == currIdx {
				// first turn of round or same player — must play a card
				if len(p.Hand) == 0 {
					r.Game.Log(fmt.Sprintf("%s 没有手牌，跳过", p.Nickname))
					r.Game.NextTurn()
				} else {
					card := p.Hand[0]
					p.Hand = p.Hand[1:]
					r.Game.HiddenCards = []Card{card}
					r.Game.State.LastPlayedCnt = 1
					r.Game.State.LastPlayer = currIdx
					r.Game.Log(fmt.Sprintf("%s 强制打出了 1 张牌", p.Nickname))

					if len(p.Hand) == 0 {
						r.Game.State.Status = "game_over"
						r.Game.State.Winner = p.Nickname
						r.Game.Log(fmt.Sprintf("🏆 %s 打完了所有手牌，获胜！", p.Nickname))
						r.Game.mu.Unlock()
						r.Broadcast()
						continue
					}
					r.Game.NextTurn()
				}
			} else {
				r.Game.CallLiar(currIdx, r.Game.State.LastPlayer)
			}
		}
		r.Game.mu.Unlock()
		r.Broadcast()
	}
}
