package room

import (
	"math/rand"
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

func RandomCode(n int) string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
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

