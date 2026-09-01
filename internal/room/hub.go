package room

import (
	cryptorand "crypto/rand"
	"math/big"
	"sync"

	"pdnode.com/play/liarsbar-web/internal/model"
)

type Hub struct {
	Rooms map[string]*Room
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// RandomCode 使用加密安全的 crypto/rand 生成 6 位随机大写字母/数字房间码
func RandomCode(n int) string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	for i := range b {
		num, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			b[i] = letters[0]
			continue
		}
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func (h *Hub) CreateRoom() *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	var code string
	for {
		code = RandomCode(6)
		if _, exists := h.Rooms[code]; !exists {
			break
		}
	}

	room := NewRoom(code, h)
	h.Rooms[code] = room
	go room.Watchdog()
	return room
}

func (h *Hub) Lock()   { h.mu.Lock() }
func (h *Hub) Unlock() { h.mu.Unlock() }

func (h *Hub) GetRoom(code string) (*Room, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[code]
	return room, exists
}

func (h *Hub) RemoveRoom(code string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Rooms, code)
}

// BroadcastGlobal 向全服所有房间及在线玩家推送全局广播通知
func (h *Hub) BroadcastGlobal(eventType string, data any) int {
	h.mu.Lock()
	rooms := make([]*Room, 0, len(h.Rooms))
	for _, r := range h.Rooms {
		rooms = append(rooms, r)
	}
	h.mu.Unlock()

	for _, r := range rooms {
		r.BroadcastEvent(eventType, data)
	}
	return len(rooms)
}

// GetPublicRooms 获取当前所有活跃房间的公开摘要列表
func (h *Hub) GetPublicRooms() []model.RoomSummary {
	h.mu.Lock()
	rooms := make([]*Room, 0, len(h.Rooms))
	for _, r := range h.Rooms {
		rooms = append(rooms, r)
	}
	h.mu.Unlock()

	summaries := make([]model.RoomSummary, 0, len(rooms))
	for _, r := range rooms {
		r.Game.Lock()
		hostName := "Host"
		playerCount := 0
		for _, p := range r.Game.State.Players {
			if !p.IsSpectator {
				playerCount++
			}
			if p.IsHost {
				hostName = p.Nickname
			}
		}
		status := r.Game.State.Status
		r.Game.Unlock()

		r.mu.Lock()
		clientCount := len(r.Clients)
		r.mu.Unlock()

		if clientCount > 0 {
			summaries = append(summaries, model.RoomSummary{
				RoomCode:    r.Code,
				HostName:    hostName,
				PlayerCount: playerCount,
				MaxPlayers:  4,
				Status:      status,
				CreatedAt:   r.CreatedAt,
			})
		}
	}
	return summaries
}
