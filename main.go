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

// --- 游戏核心模型 ---

type Card string

const (
	King  Card = "K"
	Queen Card = "Q"
	Ace   Card = "A"
	Two   Card = "2"
)

type Player struct {
	ID       string   `json:"id"`
	Nickname string   `json:"nickname"`
	Hand     []Card   `json:"hand"`
	Revolver []string `json:"-"`       // 枪管里的子弹，不传给前端
	Bullets  int      `json:"bullets"` // 剩余开枪次数
	IsAlive  bool     `json:"is_alive"`
	Client   *Client  `json:"-"`
}

type GameState struct {
	Status        string    `json:"status"` // waiting, playing, game_over
	Players       []*Player `json:"players"`
	CurrentTurn   int       `json:"current_turn"`    // 当前玩家的索引
	TableCard     Card      `json:"table_card"`      // 本轮真牌
	LastPlayer    int       `json:"last_player"`     // 上一次出牌的玩家索引
	LastPlayedCnt int       `json:"last_played_cnt"` // 上一次出了几张牌
	Logs          []string  `json:"logs"`            // 游戏日志
	Deadline      int64     `json:"deadline"`        // 回合截止时间戳
	Winner        string    `json:"winner,omitempty"`
}

type Game struct {
	State       *GameState
	HiddenCards []Card // 当前中央牌堆的实际卡牌（背面朝上）
	mu          sync.Mutex
}

// --- 房间与 WS 模型 ---

type Room struct {
	Hub     *Hub
	ID      string
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

// --- 游戏逻辑核心 ---

func NewGame() *Game {
	return &Game{
		State: &GameState{
			Status:      "waiting",
			Players:     make([]*Player, 0),
			CurrentTurn: -1,
			LastPlayer:  -1,
		},
	}
}

func (g *Game) Log(msg string) {
	g.State.Logs = append(g.State.Logs, time.Now().Format("15:04:05")+" "+msg)
	if len(g.State.Logs) > 10 {
		g.State.Logs = g.State.Logs[len(g.State.Logs)-10:]
	}
}

// 初始化/重置一轮
func (g *Game) StartRound() {
	// 构建并洗牌 24 张
	deck := []Card{King, King, King, King, King, King, Queen, Queen, Queen, Queen, Queen, Queen, Ace, Ace, Ace, Ace, Ace, Ace, Two, Two, Two, Two, Two, Two}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	// 翻开顶牌作为真牌
	g.State.TableCard = deck[0]
	deck = deck[1:]
	g.Log(fmt.Sprintf("新的一轮开始！本轮真牌是: %s", g.State.TableCard))

	// 为存活玩家发牌 (每人 5 张)
	aliveCount := 0
	for _, p := range g.State.Players {
		if p.IsAlive {
			p.Hand = deck[:5]
			deck = deck[5:]
			aliveCount++
		}
	}

	if aliveCount <= 1 {
		g.State.Status = "game_over"
		for _, p := range g.State.Players {
			if p.IsAlive {
				g.State.Winner = p.Nickname
			}
		}
		g.Log("游戏结束！")
		return
	}

	g.HiddenCards = []Card{}
	g.State.LastPlayedCnt = 0
	g.State.LastPlayer = -1

	// 如果这是整局游戏的第一轮，随机选择第一个人
	if g.State.CurrentTurn == -1 {
		g.State.CurrentTurn = rand.Intn(len(g.State.Players))
	}
	// 确保当前回合的人是活着的
	for !g.State.Players[g.State.CurrentTurn].IsAlive {
		g.State.CurrentTurn = (g.State.CurrentTurn + 1) % len(g.State.Players)
	}

	g.ResetTimer()
}

func (g *Game) ResetTimer() {
	// 给予30秒时间
	g.State.Deadline = time.Now().Add(30 * time.Second).Unix()
}

// 下一个玩家
func (g *Game) NextTurn() {
	for {
		g.State.CurrentTurn = (g.State.CurrentTurn + 1) % len(g.State.Players)
		if g.State.Players[g.State.CurrentTurn].IsAlive {
			break
		}
	}
	g.ResetTimer()
}

// 开枪结算
func (g *Game) FireGun(playerIdx int) {
	p := g.State.Players[playerIdx]
	bullet := p.Revolver[0]
	p.Revolver = p.Revolver[1:]
	p.Bullets = len(p.Revolver)

	if bullet == "Fatal" {
		p.IsAlive = false
		g.Log(fmt.Sprintf("💥 砰！%s 抽中致命子弹，被淘汰出局！", p.Nickname))
	} else {
		g.Log(fmt.Sprintf("💨 咔哒。%s 抽中空包弹，逃过一劫。", p.Nickname))
	}

	// 结算后，该玩家(或存活的下一位)成为下轮首位
	g.State.CurrentTurn = playerIdx
	for !g.State.Players[g.State.CurrentTurn].IsAlive {
		g.State.CurrentTurn = (g.State.CurrentTurn + 1) % len(g.State.Players)
	}

	// 重新开始新一轮
	g.StartRound()
}

func (r *Room) Broadcast() {
	r.Game.mu.Lock()
	b, _ := json.Marshal(map[string]interface{}{
		"type": "game_state",
		"data": r.Game.State,
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

// --- WS 路由与处理器 ---

func main() {
	rand.Seed(time.Now().UnixNano())
	e := echo.New()

	// 提供静态前端文件
	e.Static("/", "public")

	e.GET("/ws", func(c echo.Context) error {
		roomID := c.QueryParam("room")
		nickname := c.QueryParam("name")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		hub.mu.Lock()
		room, exists := hub.Rooms[roomID]
		if !exists {
			room = &Room{ID: roomID, Clients: make(map[*Client]bool), Game: NewGame()}
			hub.Rooms[roomID] = room
			// 启动房间定时器检查 (处理 30s 超时)
			go room.Watchdog()
		}
		hub.mu.Unlock()

		client := &Client{ID: fmt.Sprintf("%d", rand.Int()), Nickname: nickname, Room: room, Conn: conn, Send: make(chan []byte, 256)}

		room.mu.Lock()
		room.Clients[client] = true
		room.mu.Unlock()

		go client.WritePump()
		go client.ReadPump()

		// 玩家加入逻辑
		room.Game.mu.Lock()
		if room.Game.State.Status == "waiting" && len(room.Game.State.Players) < 4 {
			revolver := []string{"Blank", "Blank", "Blank", "Blank", "Blank", "Fatal"}
			rand.Shuffle(len(revolver), func(i, j int) { revolver[i], revolver[j] = revolver[j], revolver[i] })

			room.Game.State.Players = append(room.Game.State.Players, &Player{
				ID: client.ID, Nickname: nickname, Hand: []Card{},
				Revolver: revolver, Bullets: 6, IsAlive: true, Client: client,
			})
			room.Game.Log(nickname + " 加入了房间")
		}
		room.Game.mu.Unlock()
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
		r.Game.Log(fmt.Sprintf("👋 %s 离开了房间", p.Nickname))

		if r.Game.State.Status == "waiting" {
			// 还没开始，直接从列表移除
			r.Game.State.Players = append(r.Game.State.Players[:removedIdx], r.Game.State.Players[removedIdx+1:]...)
		} else {
			// 游戏中：标记淘汰，检查是否只剩一人
			aliveCount := 0
			for _, p := range r.Game.State.Players {
				if p.IsAlive {
					aliveCount++
				}
			}
			if aliveCount <= 1 {
				r.Game.State.Status = "game_over"
				for _, p := range r.Game.State.Players {
					if p.IsAlive {
						r.Game.State.Winner = p.Nickname
					}
				}
				r.Game.Log("游戏结束！其他玩家已离开")
			} else if r.Game.State.CurrentTurn == removedIdx {
				// 断开的是当前回合的人，跳到下一个
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
	defer func() {
		c.Room.RemoveClient(c)
	}()
	for {
		_, text, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		json.Unmarshal(text, &msg)

		g := c.Room.Game
		g.mu.Lock()

		if msg.Action == "start" && g.State.Status == "waiting" && len(g.State.Players) >= 2 {
			g.State.Status = "playing"
			g.StartRound()
		}

		if g.State.Status == "playing" && g.State.Players[g.State.CurrentTurn].ID == c.ID {
			if msg.Action == "play_cards" {
				var req struct {
					Cards []Card `json:"cards"`
				}
				json.Unmarshal(msg.Payload, &req)

				if len(req.Cards) >= 1 && len(req.Cards) <= 3 {
					// 从玩家手中移除牌
					p := g.State.Players[g.State.CurrentTurn]
					newHand := []Card{}
					// 简单过滤逻辑
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
					g.Log(fmt.Sprintf("%s 宣称打出了 %d 张真牌", p.Nickname, len(req.Cards)))
					g.NextTurn()
				}
			}

			if msg.Action == "call_liar" && g.State.LastPlayer != -1 {
				g.CallLiar(g.State.CurrentTurn, g.State.LastPlayer)
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

// 看门狗：检查超时
func (r *Room) Watchdog() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		r.Game.mu.Lock()
		if r.Game.State.Status == "playing" && time.Now().Unix() > r.Game.State.Deadline {
			currIdx := r.Game.State.CurrentTurn
			r.Game.Log(fmt.Sprintf("⏱️ %s 操作超时！系统代管...", r.Game.State.Players[currIdx].Nickname))

			// 超时逻辑：如果是首回合，强制随机出一张牌；否则强制质疑
			if r.Game.State.LastPlayer == -1 {
				p := r.Game.State.Players[currIdx]
				card := p.Hand[0]
				p.Hand = p.Hand[1:]
				r.Game.HiddenCards = []Card{card}
				r.Game.State.LastPlayedCnt = 1
				r.Game.State.LastPlayer = currIdx
				r.Game.Log(fmt.Sprintf("%s 强制打出了 1 张牌", p.Nickname))
				r.Game.NextTurn()
			} else {
				r.Game.CallLiar(currIdx, r.Game.State.LastPlayer)
			}
		}
		r.Game.mu.Unlock()
		// 每秒广播一次更新前端定时器
		r.Broadcast()
	}
}
