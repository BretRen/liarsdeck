package room

import (
	"encoding/json"
	"testing"

	"pdnode.com/play/liarsbar-web/internal/model"
)

func TestHandleClientMessagePlayCards(t *testing.T) {
	hub := NewHub()
	r := hub.CreateRoom()

	c1 := &Client{ID: "p1", Nickname: "Player1", Room: r, Send: make(chan []byte, 10)}
	c2 := &Client{ID: "p2", Nickname: "Player2", Room: r, Send: make(chan []byte, 10)}

	p1 := &model.Player{ID: "p1", Nickname: "Player1", IsAlive: true, IsReady: true, IsHost: true, Hand: []model.Card{model.King, model.Queen}, Bullets: 6}
	p2 := &model.Player{ID: "p2", Nickname: "Player2", IsAlive: true, IsReady: true, IsHost: false, Hand: []model.Card{model.Ace, model.Two}, Bullets: 6}

	r.Game.State.Players = []*model.Player{p1, p2}
	r.Game.State.Status = model.StatusPlaying
	r.Game.State.CurrentTurn = 0
	r.Game.State.TableCard = model.King

	// Test action "play"
	playPayload, _ := json.Marshal(model.PlayCardsPayload{
		Cards: []model.Card{model.King},
	})
	r.HandleClientMessage(c1, model.WSMessage{
		Action:  "play",
		Payload: playPayload,
	})

	if len(p1.Hand) != 1 {
		t.Fatalf("expected p1 hand length 1 after 'play' action, got %d", len(p1.Hand))
	}
	if r.Game.State.CurrentTurn != 1 {
		t.Fatalf("expected turn to advance to 1, got %d", r.Game.State.CurrentTurn)
	}

	// Test action "play_cards"
	playCardsPayload, _ := json.Marshal(model.PlayCardsPayload{
		Cards: []model.Card{model.Two},
	})
	r.HandleClientMessage(c2, model.WSMessage{
		Action:  "play_cards",
		Payload: playCardsPayload,
	})

	if len(p2.Hand) != 1 {
		t.Fatalf("expected p2 hand length 1 after 'play_cards' action, got %d", len(p2.Hand))
	}
	if r.Game.State.CurrentTurn != 0 {
		t.Fatalf("expected turn to advance to 0, got %d", r.Game.State.CurrentTurn)
	}
}

func TestBroadcastMasking(t *testing.T) {
	hub := NewHub()
	r := hub.CreateRoom()

	c1 := &Client{ID: "p1", Nickname: "Player1", Room: r, Send: make(chan []byte, 10)}
	c2 := &Client{ID: "p2", Nickname: "Player2", Room: r, Send: make(chan []byte, 10)}
	r.AddClient(c1)
	r.AddClient(c2)

	p1 := &model.Player{ID: "p1", Nickname: "Player1", IsAlive: true, IsReady: true, IsHost: true, Hand: []model.Card{model.King, model.Queen}, Bullets: 6}
	p2 := &model.Player{ID: "p2", Nickname: "Player2", IsAlive: true, IsReady: true, IsHost: false, Hand: []model.Card{model.Ace, model.Two}, Bullets: 6}

	r.Game.State.Players = []*model.Player{p1, p2}
	r.Game.State.Status = model.StatusPlaying

	r.Broadcast()

	select {
	case msg := <-c1.Send:
		var stateMsg struct {
			Type string `json:"type"`
			Data struct {
				Players []model.SafePlayer `json:"players"`
			} `json:"data"`
		}
		if err := json.Unmarshal(msg, &stateMsg); err != nil {
			t.Fatalf("failed to unmarshal state message: %v", err)
		}
		if len(stateMsg.Data.Players) != 2 {
			t.Fatalf("expected 2 players in state, got %d", len(stateMsg.Data.Players))
		}
		// p1 should see their own cards
		if stateMsg.Data.Players[0].Hand[0] != model.King || stateMsg.Data.Players[0].Hand[1] != model.Queen {
			t.Fatalf("p1 should see their real hand, got %v", stateMsg.Data.Players[0].Hand)
		}
		// p1 should NOT see p2's real cards (masked as "?")
		if stateMsg.Data.Players[1].Hand[0] != "?" || stateMsg.Data.Players[1].Hand[1] != "?" {
			t.Fatalf("p1 should see masked hand for p2, got %v", stateMsg.Data.Players[1].Hand)
		}
	default:
		t.Fatalf("expected broadcast message in c1.Send")
	}
}

func TestRemovePlayerOffline(t *testing.T) {
	hub := NewHub()
	r := hub.CreateRoom()

	c1 := &Client{ID: "p1", Nickname: "Host", Room: r, Send: make(chan []byte, 10)}
	r.AddClient(c1)

	p1 := &model.Player{ID: "p1", Nickname: "Host", IsAlive: true, IsReady: true, IsHost: true, ClientRef: c1}
	p2 := &model.Player{ID: "p2", Nickname: "OfflineGhost", IsAlive: true, IsReady: false, IsHost: false, ClientRef: nil} // Offline

	r.Game.State.Players = []*model.Player{p1, p2}
	r.Game.State.Status = model.StatusWaiting

	// Host kicks the offline player
	reqPayload, _ := json.Marshal(model.RemovePlayerPayload{TargetID: "p2"})
	r.HandleClientMessage(c1, model.WSMessage{
		Action:  "remove_player",
		Payload: reqPayload,
	})

	if len(r.Game.State.Players) != 1 {
		t.Fatalf("expected 1 player left after kicking offline player, got %d", len(r.Game.State.Players))
	}
	if r.Game.State.Players[0].Nickname != "Host" {
		t.Fatalf("expected Host to remain, got %s", r.Game.State.Players[0].Nickname)
	}
}

func TestHubGetPublicRooms(t *testing.T) {
	hub := NewHub()
	r := hub.CreateRoom()
	c := &Client{ID: "p1", Nickname: "Alice", Room: r, Send: make(chan []byte, 10)}
	r.AddClient(c)

	p1 := &model.Player{ID: "p1", Nickname: "Alice", IsHost: true, IsAlive: true, ClientRef: c}
	r.Game.State.Players = []*model.Player{p1}

	rooms := hub.GetPublicRooms()
	if len(rooms) != 1 {
		t.Fatalf("expected 1 public room, got %d", len(rooms))
	}
	if rooms[0].RoomCode != r.Code || rooms[0].HostName != "Alice" {
		t.Fatalf("unexpected room summary: %+v", rooms[0])
	}
}

func TestPingHandling(t *testing.T) {
	hub := NewHub()
	r := hub.CreateRoom()
	c := &Client{ID: "p1", Nickname: "Alice", Room: r, Send: make(chan []byte, 10)}

	pingPayload, _ := json.Marshal(model.PingPayload{ClientTime: 1234567890})
	r.HandleClientMessage(c, model.WSMessage{
		Action:  "ping",
		Payload: pingPayload,
	})

	select {
	case msg := <-c.Send:
		var pong struct {
			Type string `json:"type"`
			Data struct {
				ClientTime int64 `json:"client_time"`
				ServerTime int64 `json:"server_time"`
			} `json:"data"`
		}
		if err := json.Unmarshal(msg, &pong); err != nil {
			t.Fatalf("failed to unmarshal pong: %v", err)
		}
		if pong.Type != "pong" || pong.Data.ClientTime != 1234567890 {
			t.Fatalf("unexpected pong payload: %+v", pong)
		}
	default:
		t.Fatalf("expected pong in c.Send")
	}
}

