package room

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"pdnode.com/play/liarsbar-web/internal/model"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 5120
)

type Client struct {
	ID       string
	Nickname string
	Room     *Room
	Conn     *websocket.Conn
	Send     chan []byte
	once     sync.Once
}

func NewClient(id, nickname string, room *Room, conn *websocket.Conn) *Client {
	return &Client{
		ID:       id,
		Nickname: nickname,
		Room:     room,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		if c.Room != nil {
			c.Room.RemoveClient(c)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg model.WSMessage
		if err := json.Unmarshal(text, &msg); err != nil {
			continue
		}

		if c.Room != nil {
			c.Room.HandleClientMessage(c, msg)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.Send)
		c.Conn.Close()
	})
}
