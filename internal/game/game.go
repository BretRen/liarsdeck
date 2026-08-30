package game

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"pdnode.com/play/liarsbar-web/internal/model"
)

type Game struct {
	State       *model.GameState
	HiddenCards []model.Card
	mu          sync.Mutex
}

func NewGame(code string) *Game {
	return &Game{
		State: &model.GameState{
			Status:      model.StatusWaiting,
			Players:     make([]*model.Player, 0),
			CurrentTurn: -1,
			LastPlayer:  -1,
			RoomCode:    code,
			Logs:        make([]string, 0),
		},
		HiddenCards: make([]model.Card, 0),
	}
}

func (g *Game) Lock()   { g.mu.Lock() }
func (g *Game) Unlock() { g.mu.Unlock() }

func (g *Game) Log(msg string) {
	entry := time.Now().Format("15:04:05") + " " + msg
	g.State.Logs = append(g.State.Logs, entry)
	if len(g.State.Logs) > 15 {
		g.State.Logs = g.State.Logs[len(g.State.Logs)-15:]
	}
}

func (g *Game) ResetTimer() {
	g.State.Deadline = time.Now().Add(30 * time.Second).Unix()
}

func (g *Game) PauseGame(disconnectedPlayer *model.Player) {
	remaining := g.State.Deadline - time.Now().Unix()
	if remaining < 1 {
		remaining = 1
	}
	g.State.RemainingTurnSeconds = remaining
	g.State.Status = model.StatusPaused
	g.State.PausedPlayer = disconnectedPlayer.Nickname
	g.State.PauseDeadline = time.Now().Add(30 * time.Second).Unix()
	g.Log(fmt.Sprintf("⏸️ 玩家 %s 断线，游戏暂停 30 秒等待重连... / ⏸️ %s disconnected, game paused for 30s...", disconnectedPlayer.Nickname, disconnectedPlayer.Nickname))
}

func (g *Game) ResumeGame(reconnectedPlayer *model.Player) {
	g.State.Status = model.StatusPlaying
	g.State.PausedPlayer = ""
	g.State.PauseDeadline = 0

	remaining := g.State.RemainingTurnSeconds
	if remaining < 3 {
		remaining = 3
	}
	g.State.Deadline = time.Now().Add(time.Duration(remaining) * time.Second).Unix()
	g.Log(fmt.Sprintf("▶️ 玩家 %s 重新连接成功，游戏继续！(剩余操作时间: %d秒) / ▶️ %s reconnected, game resumed!", reconnectedPlayer.Nickname, remaining, reconnectedPlayer.Nickname))
}

func (g *Game) HandlePauseTimeout() {
	var timedOutPlayer *model.Player
	for _, p := range g.State.Players {
		if p.Nickname == g.State.PausedPlayer {
			timedOutPlayer = p
			break
		}
	}

	if timedOutPlayer != nil && timedOutPlayer.IsAlive && timedOutPlayer.ClientRef == nil {
		timedOutPlayer.IsAlive = false
		g.Log(fmt.Sprintf("⌛ 玩家 %s 30秒断线重连超时，已被判定出局！ / ⌛ %s failed to reconnect within 30s and was eliminated!", timedOutPlayer.Nickname, timedOutPlayer.Nickname))
	}

	g.State.Status = model.StatusPlaying
	g.State.PausedPlayer = ""
	g.State.PauseDeadline = 0

	g.AdvanceToAlive()

	aliveCount := 0
	var lastAlive *model.Player
	for _, pp := range g.State.Players {
		if !pp.IsSpectator && pp.IsAlive {
			aliveCount++
			lastAlive = pp
		}
	}
	if aliveCount <= 1 && lastAlive != nil {
		g.State.Status = model.StatusGameOver
		g.State.Winner = lastAlive.Nickname
		g.Log(fmt.Sprintf("🏆 %s 获胜！ / 🏆 %s wins!", lastAlive.Nickname, lastAlive.Nickname))
		return
	}

	g.ResetTimer()
}

func (g *Game) StartRound() {
	deck := NewDeck()
	var tableCard model.Card
	tableCard, deck = DrawTableCard(deck)
	g.State.TableCard = tableCard

	g.Log(fmt.Sprintf("新的一轮开始！本轮真牌是: %s / New round! True card: %s", tableCard, tableCard))

	alive := make([]*model.Player, 0)
	others := make([]*model.Player, 0)
	for _, p := range g.State.Players {
		if !p.IsSpectator && p.IsAlive {
			alive = append(alive, p)
		} else {
			others = append(others, p)
		}
	}

	CryptoShuffle(len(alive), func(i, j int) {
		alive[i], alive[j] = alive[j], alive[i]
	})
	g.State.Players = append(alive, others...)

	aliveCount := 0
	for _, p := range g.State.Players {
		if !p.IsSpectator && p.IsAlive {
			if len(deck) >= 5 {
				p.Hand = make([]model.Card, 5)
				copy(p.Hand, deck[:5])
				deck = deck[5:]
			} else {
				p.Hand = make([]model.Card, len(deck))
				copy(p.Hand, deck)
				deck = nil
			}
			aliveCount++
		}
	}

	if aliveCount <= 1 {
		g.State.Status = model.StatusGameOver
		for _, p := range g.State.Players {
			if !p.IsSpectator && p.IsAlive {
				g.State.Winner = p.Nickname
			}
		}
		g.Log("游戏结束！ / Game over!")
		return
	}

	g.State.Status = model.StatusPlaying
	g.HiddenCards = []model.Card{}
	g.State.LastPlayedCnt = 0
	g.State.LastPlayer = -1

	if g.State.CurrentTurn == -1 && len(g.State.Players) > 0 {
		nBig, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(g.State.Players))))
		if err == nil {
			g.State.CurrentTurn = int(nBig.Int64())
		} else {
			g.State.CurrentTurn = 0
		}
	}
	g.AdvanceToAlive()
	g.ResetTimer()
}

func (g *Game) AdvanceToAlive() {
	n := len(g.State.Players)
	if n == 0 {
		return
	}
	start := g.State.CurrentTurn
	if start < 0 {
		start = 0
	}
	for i := 0; i < n; i++ {
		idx := (start + i) % n
		if g.State.Players[idx].IsAlive && !g.State.Players[idx].IsSpectator {
			g.State.CurrentTurn = idx
			return
		}
	}

	g.State.Status = model.StatusGameOver
	g.State.Winner = "无人存活"
	g.Log("游戏结束！所有玩家都已淘汰 / All players eliminated")
}

func (g *Game) NextTurn() {
	start := g.State.CurrentTurn
	n := len(g.State.Players)
	for i := 1; i <= n; i++ {
		idx := (start + i) % n
		if g.State.Players[idx].IsAlive && !g.State.Players[idx].IsSpectator {
			g.State.CurrentTurn = idx
			g.ResetTimer()
			return
		}
	}
	g.State.Status = model.StatusGameOver
	g.State.Winner = "无人存活"
	g.Log("游戏结束！所有玩家都已淘汰 / All players eliminated")
}

func (g *Game) PlayCards(playerIdx int, cards []model.Card) bool {
	if len(cards) < 1 || len(cards) > 3 {
		return false
	}
	if playerIdx < 0 || playerIdx >= len(g.State.Players) {
		return false
	}

	// 如果上家手牌已经出空，下家必须发起质疑，禁止继续出牌
	if g.State.LastPlayer >= 0 && g.State.LastPlayer < len(g.State.Players) {
		lastP := g.State.Players[g.State.LastPlayer]
		if len(lastP.Hand) == 0 {
			return false
		}
	}

	p := g.State.Players[playerIdx]
	if !p.IsAlive || p.IsSpectator {
		return false
	}

	newHand := make([]model.Card, 0, len(p.Hand))
	used := make([]bool, len(cards))
	for _, hc := range p.Hand {
		removed := false
		for i, rc := range cards {
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
	for _, u := range used {
		if !u {
			return false
		}
	}

	p.Hand = newHand
	g.HiddenCards = cards
	g.State.LastPlayedCnt = len(cards)
	g.State.LastPlayer = playerIdx

	g.Log(fmt.Sprintf("%s 打出了 %d 张暗牌 / %s played %d cards", p.Nickname, len(cards), p.Nickname, len(cards)))

	if len(p.Hand) == 0 {
		g.Log(fmt.Sprintf("%s 手牌已清空！下家必须质疑！ / %s has 0 cards left! Next player must call liar!", p.Nickname, p.Nickname))
	}

	g.NextTurn()
	return true
}

func (g *Game) CallLiar(
	callerIdx, accusedIdx int,
	onLiarCall func(caller, accused string),
	onReveal func(caller, accused string, cards []model.Card),
	onShot func(target string, fatal bool),
) {
	if callerIdx < 0 || callerIdx >= len(g.State.Players) ||
		accusedIdx < 0 || accusedIdx >= len(g.State.Players) ||
		callerIdx == accusedIdx || len(g.HiddenCards) == 0 {
		return
	}

	caller := g.State.Players[callerIdx]
	accused := g.State.Players[accusedIdx]

	g.Log(fmt.Sprintf("🚨 %s 质疑 %s 说谎！ / 🚨 %s calls out %s!", caller.Nickname, accused.Nickname, caller.Nickname, accused.Nickname))

	if onLiarCall != nil {
		onLiarCall(caller.Nickname, accused.Nickname)
	}
	if onReveal != nil {
		onReveal(caller.Nickname, accused.Nickname, g.HiddenCards)
	}

	isLiar := false
	revealMsg := fmt.Sprintf("%s 的底牌是: / %s's cards: ", accused.Nickname, accused.Nickname)
	for _, c := range g.HiddenCards {
		revealMsg += string(c) + " "
		if c != g.State.TableCard && c != model.Two {
			isLiar = true
		}
	}
	g.Log(revealMsg)

	if isLiar {
		g.Log("👉 质疑成功！出牌者说谎！ / ✅ Liar! The bluffer was caught!")
		g.FireGun(accusedIdx, onShot)
	} else {
		g.Log("❌ 质疑失败！出牌者是清白的！ / ❌ Challenge failed! The cards were honest!")
		// 1. 质疑者（诬告者）必须先开枪受罚！
		g.ShootPlayer(callerIdx, onShot)

		// 2. 检查剩余存活人数
		g.State.CurrentTurn = callerIdx
		g.AdvanceToAlive()
		if g.State.Status == model.StatusGameOver {
			return
		}

		aliveCount := 0
		var lastAlive *model.Player
		for _, pp := range g.State.Players {
			if !pp.IsSpectator && pp.IsAlive {
				aliveCount++
				lastAlive = pp
			}
		}

		// 唯一胜利条件：全场仅剩最后一名存活者
		if aliveCount <= 1 && lastAlive != nil {
			g.State.Status = model.StatusGameOver
			g.State.Winner = lastAlive.Nickname
			g.Log(fmt.Sprintf("🏆 %s 成为唯一存活者，获得整场胜利！ / 🏆 %s is the sole survivor and wins!", lastAlive.Nickname, lastAlive.Nickname))
			return
		}

		if len(accused.Hand) == 0 {
			g.Log(fmt.Sprintf("%s 成功出清全部手牌，本轮结束，重新发牌！ / %s emptied all cards, dealing new round!", accused.Nickname, accused.Nickname))
		}

		// 存活人数 > 1，重新洗牌发牌开启新一轮
		g.StartRound()
	}
}

func (g *Game) ShootPlayer(playerIdx int, onShot func(target string, fatal bool)) bool {
	if playerIdx < 0 || playerIdx >= len(g.State.Players) {
		return false
	}
	p := g.State.Players[playerIdx]

	if len(p.Revolver) == 0 {
		p.Revolver = NewRevolver()
		p.Bullets = 6
	}

	bullet := p.Revolver[0]
	p.Revolver = p.Revolver[1:]
	p.Bullets = len(p.Revolver)

	isFatal := bullet == model.BulletFatal

	if onShot != nil {
		onShot(p.Nickname, isFatal)
	}

	if isFatal {
		p.IsAlive = false
		g.Log(fmt.Sprintf("💥 砰！%s 抽中致命子弹，被淘汰出局！ / 💥 BANG! %s was shot fatally!", p.Nickname, p.Nickname))
	} else {
		g.Log(fmt.Sprintf("💨 咔哒。%s 抽中空包弹，逃过一劫。 / 💨 Click. %s survived!", p.Nickname, p.Nickname))
	}

	return isFatal
}

func (g *Game) FireGun(playerIdx int, onShot func(target string, fatal bool)) {
	g.ShootPlayer(playerIdx, onShot)

	g.State.CurrentTurn = playerIdx
	g.AdvanceToAlive()
	if g.State.Status == model.StatusGameOver {
		return
	}

	aliveCount := 0
	var lastAlive *model.Player
	for _, pp := range g.State.Players {
		if !pp.IsSpectator && pp.IsAlive {
			aliveCount++
			lastAlive = pp
		}
	}
	if aliveCount <= 1 && lastAlive != nil {
		g.State.Status = model.StatusGameOver
		g.State.Winner = lastAlive.Nickname
		g.Log(fmt.Sprintf("🏆 %s 获胜！ / 🏆 %s wins!", lastAlive.Nickname, lastAlive.Nickname))
		return
	}

	g.StartRound()
}

func (g *Game) ResetGame() {
	g.State.Status = model.StatusWaiting
	g.State.TableCard = ""
	g.State.LastPlayedCnt = 0
	g.State.LastPlayer = -1
	g.State.CurrentTurn = -1
	g.State.Winner = ""
	g.State.PauseDeadline = 0
	g.State.PausedPlayer = ""
	g.State.RemainingTurnSeconds = 0
	g.HiddenCards = []model.Card{}

	// 仅保留当前处于在线连接状态的玩家，自动清理断线玩家
	connectedPlayers := make([]*model.Player, 0, len(g.State.Players))
	for _, p := range g.State.Players {
		if p.ClientRef != nil {
			p.Hand = []model.Card{}
			p.IsReady = false
			p.IsAlive = true
			p.Bullets = 6
			p.Revolver = NewRevolver()
			p.HasUsedDisconnectGrace = false
			connectedPlayers = append(connectedPlayers, p)
		} else {
			g.Log(fmt.Sprintf("🧹 清理断线玩家: %s / 🧹 Removed disconnected player: %s", p.Nickname, p.Nickname))
		}
	}
	g.State.Players = connectedPlayers

	// 如果原房主已断线被清理，确保首位在线玩家继承房主身份
	hasHost := false
	for _, p := range g.State.Players {
		if p.IsHost {
			hasHost = true
			break
		}
	}
	if !hasHost && len(g.State.Players) > 0 {
		for _, p := range g.State.Players {
			if !p.IsSpectator {
				p.IsHost = true
				g.Log(fmt.Sprintf("👑 %s 成为新房主 / 👑 %s is the new host", p.Nickname, p.Nickname))
				break
			}
		}
	}

	g.Log("🔄 游戏已重置，等待玩家准备 / Game reset")
}
