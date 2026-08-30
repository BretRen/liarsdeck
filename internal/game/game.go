package game

import (
	"fmt"
	"math/rand"
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

	rand.Shuffle(len(alive), func(i, j int) {
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

	g.HiddenCards = []model.Card{}
	g.State.LastPlayedCnt = 0
	g.State.LastPlayer = -1

	if g.State.CurrentTurn == -1 {
		g.State.CurrentTurn = rand.Intn(len(g.State.Players))
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
	p := g.State.Players[playerIdx]

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

	// 必须能够匹配所有打出的牌
	for _, u := range used {
		if !u {
			return false
		}
	}

	p.Hand = newHand
	g.HiddenCards = cards
	g.State.LastPlayedCnt = len(cards)
	g.State.LastPlayer = playerIdx

	g.Log(fmt.Sprintf("%s 宣称打出了 %d 张牌 / %s played %d card(s)", p.Nickname, len(cards), p.Nickname, len(cards)))

	if len(p.Hand) == 0 {
		g.Log(fmt.Sprintf("%s 打完了所有手牌！下家必须质疑 / %s emptied their hand! Next player must call liar", p.Nickname, p.Nickname))
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
		accusedIdx < 0 || accusedIdx >= len(g.State.Players) {
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
		if len(accused.Hand) == 0 {
			g.State.Status = model.StatusGameOver
			g.State.Winner = accused.Nickname
			g.Log(fmt.Sprintf("🎉 🏆 %s 成功出完手牌且未说谎，获得游戏胜利！ / 🎉 🏆 %s emptied all cards honestly and wins!", accused.Nickname, accused.Nickname))
			return
		}
		g.FireGun(callerIdx, onShot)
	}
}

func (g *Game) FireGun(playerIdx int, onShot func(target string, fatal bool)) {
	if playerIdx < 0 || playerIdx >= len(g.State.Players) {
		return
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
	alive := make([]*model.Player, 0)
	for _, p := range g.State.Players {
		if p.ClientRef == nil {
			continue
		}
		p.IsAlive = true
		p.Hand = []model.Card{}
		p.Bullets = 6
		p.Revolver = NewRevolver()
		p.IsReady = false
		alive = append(alive, p)
	}
	g.State.Players = alive
	g.State.Winner = ""
	g.State.Logs = []string{}
	g.State.CurrentTurn = -1
	g.State.LastPlayer = -1
	g.State.LastPlayedCnt = 0

	for _, p := range g.State.Players {
		p.IsHost = false
	}
	for _, p := range g.State.Players {
		if !p.IsSpectator {
			p.IsHost = true
			break
		}
	}

	g.State.Status = model.StatusWaiting
	g.Log("重新开始！请准备 / New game! Please ready up")
}
