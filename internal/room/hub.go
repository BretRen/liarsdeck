package room

import (
	cryptorand "crypto/rand"
	"math/big"
	"sync"
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
