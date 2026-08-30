package handler

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"pdnode.com/play/liarsbar-web/internal/game"
	"pdnode.com/play/liarsbar-web/internal/model"
	"pdnode.com/play/liarsbar-web/internal/room"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		u, err := url.Parse(origin)
		if err != nil {
			return false
		}
		// 校验同源 (Origin host 与 r.Host 一致) 或允许本地开发端口 (localhost / 127.0.0.1)
		if strings.EqualFold(u.Host, r.Host) {
			return true
		}
		hostName := u.Hostname()
		if hostName == "localhost" || hostName == "127.0.0.1" {
			return true
		}
		return false
	},
}

type WSHandler struct {
	Hub *room.Hub
}

func NewWSHandler(hub *room.Hub) *WSHandler {
	return &WSHandler{Hub: hub}
}

func secureRandomString(prefix string) string {
	n, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(9000000000000000))
	return fmt.Sprintf("%s%d", prefix, n.Int64()+1000000000000000)
}

func (h *WSHandler) HandleWebSocket(c echo.Context) error {
	action := c.QueryParam("action")
	code := c.QueryParam("code")
	nickname := c.QueryParam("name")
	token := c.QueryParam("token")

	if nickname == "" {
		n, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(900))
		nickname = fmt.Sprintf("Player%d", n.Int64()+100)
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	clientID := secureRandomString("")
	if action == "reconnect" && token != "" {
		clientID = token
	}

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
				if r.Game.State.Status == model.StatusPaused && r.Game.State.PausedPlayer == p.Nickname {
					r.Game.ResumeGame(p)
				} else {
					r.Game.Log(fmt.Sprintf("🔗 %s 重新连接 / 🔗 %s reconnected", nickname, nickname))
				}
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
		r.Game.Log(fmt.Sprintf("👋 %s 加入了房间 / 👋 %s joined room", nickname, nickname))
	} else {
		// 作为观战者加入
		r.Game.State.Players = append(r.Game.State.Players, &model.Player{
			ID:          client.ID,
			Nickname:    nickname,
			Hand:        []model.Card{},
			Bullets:     0,
			IsAlive:     false,
			IsHost:      false,
			IsSpectator: true,
			ClientRef:   client,
		})
		r.Game.Log(fmt.Sprintf("👀 %s 正在观战 / 👀 %s is spectating", nickname, nickname))
	}
	r.Game.Unlock()

	go client.WritePump()
	go client.ReadPump()
	r.Broadcast()
	return nil
}
