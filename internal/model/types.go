package model

import "encoding/json"

type Card string

const (
	King  Card = "K"
	Queen Card = "Q"
	Ace   Card = "A"
	Two   Card = "2" // 万能牌 Wild Card
)

const (
	BulletBlank = "Blank"
	BulletFatal = "Fatal"
)

const (
	StatusWaiting  = "waiting"
	StatusPlaying  = "playing"
	StatusPaused   = "paused"
	StatusGameOver = "game_over"
)

type Player struct {
	ID                     string   `json:"id"`
	Nickname               string   `json:"nickname"`
	Hand                   []Card   `json:"hand"`
	Revolver               []string `json:"-"`
	Bullets                int      `json:"bullets"`
	IsAlive                bool     `json:"is_alive"`
	IsHost                 bool     `json:"is_host"`
	IsSpectator            bool     `json:"is_spectator"`
	IsReady                bool     `json:"is_ready"`
	HasUsedDisconnectGrace bool     `json:"has_used_disconnect_grace"` // 是否已触发过1次30秒断线保护
	ClientRef              any      `json:"-"`
}

type SafePlayer struct {
	ID                     string `json:"id"`
	Nickname               string `json:"nickname"`
	Hand                   []Card `json:"hand"`
	Bullets                int    `json:"bullets"`
	IsAlive                bool   `json:"is_alive"`
	IsHost                 bool   `json:"is_host"`
	IsSpectator            bool   `json:"is_spectator"`
	IsReady                bool   `json:"is_ready"`
	IsConnected            bool   `json:"is_connected"`
	HasUsedDisconnectGrace bool   `json:"has_used_disconnect_grace"`
}

func (p *Player) ToSafe() SafePlayer {
	handCopy := make([]Card, len(p.Hand))
	copy(handCopy, p.Hand)
	return SafePlayer{
		ID:                     p.ID,
		Nickname:               p.Nickname,
		Hand:                   handCopy,
		Bullets:                p.Bullets,
		IsAlive:                p.IsAlive,
		IsHost:                 p.IsHost,
		IsSpectator:            p.IsSpectator,
		IsReady:                p.IsReady,
		IsConnected:            p.ClientRef != nil,
		HasUsedDisconnectGrace: p.HasUsedDisconnectGrace,
	}
}

type GameState struct {
	Status               string    `json:"status"`
	Players              []*Player `json:"players"`
	CurrentTurn          int       `json:"current_turn"`
	TableCard            Card      `json:"table_card"`
	LastPlayer           int       `json:"last_player"`
	LastPlayedCnt        int       `json:"last_played_cnt"`
	Logs                 []string  `json:"logs"`
	Deadline             int64     `json:"deadline"`
	PauseDeadline        int64     `json:"pause_deadline,omitempty"`
	PausedPlayer         string    `json:"paused_player,omitempty"`
	RemainingTurnSeconds int64     `json:"remaining_turn_seconds,omitempty"`
	Winner               string    `json:"winner,omitempty"`
	RoomCode             string    `json:"room_code"`
}

type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type PlayCardsPayload struct {
	Cards []Card `json:"cards"`
}

type RemovePlayerPayload struct {
	TargetID string `json:"target_id"`
}

type EventMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type LiarCallEvent struct {
	Caller  string `json:"caller"`
	Accused string `json:"accused"`
}

type RevealEvent struct {
	Caller  string `json:"caller"`
	Accused string `json:"accused"`
	Cards   []Card `json:"cards"`
}

type ShotEvent struct {
	Target string `json:"target"`
	Fatal  bool   `json:"fatal"`
}
