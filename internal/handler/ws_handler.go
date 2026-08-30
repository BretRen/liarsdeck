package handler

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"pdnode.com/play/liarsbar-web/internal/game"
	"pdnode.com/play/liarsbar-web/internal/model"
	"pdnode.com/play/liarsbar-web/internal/room"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	Hub *room.Hub
}

func NewWSHandler(hub *room.Hub) *WSHandler {
	return &WSHandler{Hub: hub}
}

func (h *WSHandler) HandleWebSocket(c echo.Context) error {
	action := c.QueryParam("action")
	code := c.QueryParam("code")
	nickname := c.QueryParam("name")
	token := c.QueryParam("token")

	if nickname == "" {
		nickname = fmt.Sprintf("Player%d", rand.Intn(900)+100)
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	clientID := fmt.Sprintf("%d", rand.Int63())

	if action == "create" {
		r := h.Hub.CreateRoom()
		client := room.NewClient(clientID, nickname, r, conn)

		r.AddClient(client)

		r.Game.Lock()
		r.Game.State.Players = append(r.Game.State.Players, &model.Player{
			ID:          client.ID,
			Nickname:    nickname,
			Hand:        []model.Card{},
			Revolver:    game.NewRevolver(),
			Bullets:     6,
			IsAlive:     true,
			IsHost:      true,
			IsSpectator: false,
			ClientRef:   client,
		})
		r.Game.Log(fmt.Sprintf("🏠 %s 创建了房间 %s / 🏠 %s created room %s", nickname, r.Code, nickname, r.Code))
		r.Game.Unlock()

		go client.WritePump()
		go client.ReadPump()
		r.Broadcast()
		return nil
	}

	r, exists := h.Hub.GetRoom(code)
	if !exists {
		_ = conn.WriteJSON(map[string]string{"error": "房间不存在 / Room does not exist"})
		_ = conn.Close()
		return nil
	}

	client := room.NewClient(clientID, nickname, r, conn)
	r.AddClient(client)

	r.Game.Lock()
	isSpectator := action == "spectate"

	// ---- Reconnect 重连逻辑 ----
	if action == "reconnect" {
		for _, p := range r.Game.State.Players {
			if p.ID == token && p.ClientRef == nil {
				p.ClientRef = client
				r.Game.Log(fmt.Sprintf("🔗 %s 重新连接 / 🔗 %s reconnected", nickname, nickname))
				r.Game.Unlock()
				go client.WritePump()
				go client.ReadPump()
				r.Broadcast()
				return nil
			}
		}
		r.Game.Unlock()
		_ = conn.WriteJSON(map[string]string{"error": "无法重连 / Reconnect failed"})
		_ = conn.Close()
		return nil
	}

	// 统计非观战玩家数
	playerCount := 0
	for _, pp := range r.Game.State.Players {
		if !pp.IsSpectator {
			playerCount++
		}
	}

	if !isSpectator && r.Game.State.Status == model.StatusWaiting && playerCount < 4 {
		r.Game.State.Players = append(r.Game.State.Players, &model.Player{
			ID:          client.ID,
			Nickname:    nickname,
			Hand:        []model.Card{},
			Revolver:    game.NewRevolver(),
			Bullets:     6,
			IsAlive:     true,
			IsHost:      false,
			IsSpectator: false,
			ClientRef:   client,
		})
		r.Game.Log(nickname + " 加入了房间 / " + nickname + " joined the room")
	} else if isSpectator {
		r.Game.State.Players = append(r.Game.State.Players, &model.Player{
			ID:          client.ID,
			Nickname:    nickname,
			Hand:        []model.Card{},
			Revolver:    []string{},
			Bullets:     0,
			IsAlive:     true,
			IsHost:      false,
			IsSpectator: true,
			ClientRef:   client,
		})
		r.Game.Log(fmt.Sprintf("👀 %s 以观众身份加入 / 👀 %s is spectating", nickname, nickname))
	} else if !isSpectator && r.Game.State.Status == model.StatusWaiting {
		r.Game.Log(fmt.Sprintf("⚠️ %s 尝试加入但房间已满 / ⚠️ %s tried to join but room is full", nickname, nickname))
	} else if !isSpectator {
		r.Game.Log(fmt.Sprintf("👋 %s 尝试加入进行中的游戏，转为观战 / 👋 %s joined as spectator (game in progress)", nickname, nickname))
		r.Game.State.Players = append(r.Game.State.Players, &model.Player{
			ID:          client.ID,
			Nickname:    nickname,
			Hand:        []model.Card{},
			Revolver:    []string{},
			Bullets:     0,
			IsAlive:     true,
			IsHost:      false,
			IsSpectator: true,
			ClientRef:   client,
		})
	}
	r.Game.Unlock()

	go client.WritePump()
	go client.ReadPump()
	r.Broadcast()
	return nil
}
