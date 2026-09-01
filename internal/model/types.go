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

const (
	ModeClassic = "classic"
	ModeItems   = "items"
)

type ItemType string

const (
	ItemEagleEye    ItemType = "eagle_eye"    // 鹰眼透镜: 窥视最后一组底牌中的任意一张
	ItemSawedOff    ItemType = "sawed_off"    // 截短猎枪: 下一次开枪判定双倍扣动扳机
	ItemHardLiquor  ItemType = "hard_liquor"  // 烈性烧酒: 弃置手牌中非台面假牌并重抽
	ItemKevlarArmor ItemType = "kevlar_armor" // 防弹护甲: 抵消一次实弹暴毙死亡
	ItemFateShift   ItemType = "fate_shift"   // 命运重抽: 强行重置台面指定花色
)

type Player struct {
	ID                       string     `json:"id"`
	Nickname                 string     `json:"nickname"`
	Hand                     []Card     `json:"hand"`
	Revolver                 []string   `json:"-"`
	Bullets                  int        `json:"bullets"`
	IsAlive                  bool       `json:"is_alive"`
	IsHost                   bool       `json:"is_host"`
	IsSpectator              bool       `json:"is_spectator"`
	IsReady                  bool       `json:"is_ready"`
	DisconnectGraceRemaining int        `json:"disconnect_grace_remaining"` // 剩余断线保护时间（秒，初始 30 秒）
	Items                    []ItemType `json:"items"`                      // 玩家拥有的道具列表（最多 2 个）
	HasArmor                 bool       `json:"has_armor"`                  // 是否装配了防弹护甲（抵消一次致命死亡）
	ClientRef                any        `json:"-"`
}

type SafePlayer struct {
	ID                       string     `json:"id"`
	Nickname                 string     `json:"nickname"`
	Hand                     []Card     `json:"hand"`
	Bullets                  int        `json:"bullets"`
	IsAlive                  bool       `json:"is_alive"`
	IsHost                   bool       `json:"is_host"`
	IsSpectator              bool       `json:"is_spectator"`
	IsReady                  bool       `json:"is_ready"`
	IsConnected              bool       `json:"is_connected"`
	DisconnectGraceRemaining int        `json:"disconnect_grace_remaining"`
	Items                    []ItemType `json:"items"`      // 私密道具：本人可见具体道具列表，对手仅见占位/数量
	ItemCount                int        `json:"item_count"` // 对手可见的道具总数量
	HasArmor                 bool       `json:"has_armor"`  // 公开护盾状态，所有人可见
}

func (p *Player) ToSafe(viewerID string) SafePlayer {
	var handCopy []Card
	var itemsCopy []ItemType

	if viewerID != "" && viewerID == p.ID {
		handCopy = make([]Card, len(p.Hand))
		copy(handCopy, p.Hand)
		itemsCopy = make([]ItemType, len(p.Items))
		copy(itemsCopy, p.Items)
	} else {
		// 隐藏对手的手牌内容与道具私密详情，仅保留数量以供 UI 渲染，防止作弊
		handCopy = make([]Card, len(p.Hand))
		for i := range handCopy {
			handCopy[i] = "?"
		}
		itemsCopy = []ItemType{}
	}

	return SafePlayer{
		ID:                       p.ID,
		Nickname:                 p.Nickname,
		Hand:                     handCopy,
		Bullets:                  p.Bullets,
		IsAlive:                  p.IsAlive,
		IsHost:                   p.IsHost,
		IsSpectator:              p.IsSpectator,
		IsReady:                  p.IsReady,
		IsConnected:              p.ClientRef != nil,
		DisconnectGraceRemaining: p.DisconnectGraceRemaining,
		Items:                    itemsCopy,
		ItemCount:                len(p.Items),
		HasArmor:                 p.HasArmor,
	}
}

type GameState struct {
	Status               string    `json:"status"`
	Players              []*Player `json:"players"`
	CurrentTurn          int       `json:"current_turn"`
	TableCard            Card      `json:"table_card"`
	LastPlayer           int       `json:"last_player"`
	LastPlayedCnt        int       `json:"last_played_cnt"`
	LastPlayedCards      []Card    `json:"-"` // 服务端保留上一组实际出牌，用于鹰眼透镜窥视
	Logs                 []string  `json:"logs"`
	Deadline             int64     `json:"deadline"`
	PauseDeadline        int64     `json:"pause_deadline,omitempty"`
	PausedPlayer         string    `json:"paused_player,omitempty"`
	PausedPlayerID       string    `json:"paused_player_id,omitempty"`
	RemainingTurnSeconds int64     `json:"remaining_turn_seconds,omitempty"`
	Winner               string    `json:"winner,omitempty"`
	RoomCode             string    `json:"room_code"`
	GameMode             string    `json:"game_mode"`              // "classic" 或 "items"
	MaxPlayers           int       `json:"max_players"`            // 2, 3 或 4
	DoubleDamage         bool      `json:"double_damage"`          // 截短猎枪双倍开枪状态
}

type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type PlayCardsPayload struct {
	Cards []Card `json:"cards"`
}

type UseItemPayload struct {
	Item ItemType `json:"item"`
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
	Target       string `json:"target"`
	Fatal        bool   `json:"fatal"`
	ArmorBlocked bool   `json:"armor_blocked,omitempty"` // 是否被防弹护甲抵消
	DoubleShot   bool   `json:"double_shot,omitempty"`   // 是否为截短双倍枪击
}

type ItemUsedEvent struct {
	PlayerID   string   `json:"player_id"`
	Nickname   string   `json:"nickname"`
	Item       ItemType `json:"item"`
	ItemName   string   `json:"item_name"`
	TargetCard Card     `json:"target_card,omitempty"` // 命运重抽时的新台面牌
}

type RoomSummary struct {
	RoomCode    string `json:"room_code"`
	HostName    string `json:"host_name"`
	PlayerCount int    `json:"player_count"`
	MaxPlayers  int    `json:"max_players"`
	GameMode    string `json:"game_mode"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

type PingPayload struct {
	ClientTime int64 `json:"client_time"`
}
