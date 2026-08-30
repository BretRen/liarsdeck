package game

import (
	"testing"
	"time"

	"pdnode.com/play/liarsbar-web/internal/model"
)

func TestNewDeckAndDrawTableCard(t *testing.T) {
	deck := NewDeck()
	if len(deck) != 24 {
		t.Fatalf("expected 24 cards, got %d", len(deck))
	}

	tableCard, remaining := DrawTableCard(deck)
	if tableCard == model.Two {
		t.Fatalf("table card should never be 2 (wild)")
	}
	if len(remaining) != 23 {
		t.Fatalf("expected 23 remaining cards, got %d", len(remaining))
	}
}

func TestNewRevolver(t *testing.T) {
	revolver := NewRevolver()
	if len(revolver) != 6 {
		t.Fatalf("expected 6 bullets in revolver, got %d", len(revolver))
	}

	fatalCount := 0
	blankCount := 0
	for _, b := range revolver {
		if b == model.BulletFatal {
			fatalCount++
		} else if b == model.BulletBlank {
			blankCount++
		}
	}
	if fatalCount != 1 || blankCount != 5 {
		t.Fatalf("expected 1 fatal and 5 blank, got fatal=%d, blank=%d", fatalCount, blankCount)
	}
}

func TestGameStartRoundAndPlayCards(t *testing.T) {
	g := NewGame("TEST01")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", IsAlive: true, IsReady: true, Bullets: 6, Revolver: NewRevolver()}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", IsAlive: true, IsReady: true, Bullets: 6, Revolver: NewRevolver()}
	g.State.Players = []*model.Player{p1, p2}

	g.StartRound()

	if len(p1.Hand) != 5 || len(p2.Hand) != 5 {
		t.Fatalf("each alive player should have 5 cards")
	}
	if g.State.TableCard == model.Two || g.State.TableCard == "" {
		t.Fatalf("invalid table card: %s", g.State.TableCard)
	}

	curr := g.State.CurrentTurn
	currPlayer := g.State.Players[curr]
	playedCard := currPlayer.Hand[0]

	success := g.PlayCards(curr, []model.Card{playedCard})
	if !success {
		t.Fatalf("play cards failed")
	}
	if len(currPlayer.Hand) != 4 {
		t.Fatalf("expected 4 cards in hand after play, got %d", len(currPlayer.Hand))
	}
	if g.State.LastPlayer != curr {
		t.Fatalf("expected last player %d, got %d", curr, g.State.LastPlayer)
	}
	if g.State.CurrentTurn == curr {
		t.Fatalf("expected turn to advance")
	}
}

func TestCallLiarHonest(t *testing.T) {
	g := NewGame("TEST02")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", Hand: []model.Card{model.King}, IsAlive: true, Bullets: 6, Revolver: []string{model.BulletBlank, model.BulletFatal}}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", Hand: []model.Card{model.Queen}, IsAlive: true, Bullets: 6, Revolver: []string{model.BulletBlank, model.BulletFatal}}
	g.State.Players = []*model.Player{p1, p2}
	g.State.TableCard = model.King
	g.HiddenCards = []model.Card{model.King, model.Two} // Honest!

	var shotTarget string
	var shotFatal bool

	// p2 calls p1 a liar, but p1 played King and Two (Honest). So p2 (caller) gets shot.
	g.CallLiar(1, 0, nil, nil, func(target string, fatal bool) {
		shotTarget = target
		shotFatal = fatal
	})

	if shotTarget != "Bob" {
		t.Fatalf("expected Bob (caller) to get shot for false challenge, got %s", shotTarget)
	}
	if shotFatal != false {
		t.Fatalf("first bullet in mock was blank")
	}
}

func TestCallLiarBluff(t *testing.T) {
	g := NewGame("TEST03")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", IsAlive: true, Bullets: 6, Revolver: []string{model.BulletFatal}}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", IsAlive: true, Bullets: 6, Revolver: []string{model.BulletFatal}}
	g.State.Players = []*model.Player{p1, p2}
	g.State.TableCard = model.King
	g.HiddenCards = []model.Card{model.Queen} // Bluff!

	var shotTarget string
	var shotFatal bool

	// p2 calls p1 a liar, p1 played Queen (Liar). p1 gets shot.
	g.CallLiar(1, 0, nil, nil, func(target string, fatal bool) {
		shotTarget = target
		shotFatal = fatal
	})

	if shotTarget != "Alice" {
		t.Fatalf("expected Alice (accused) to get shot for bluffing, got %s", shotTarget)
	}
	if !shotFatal {
		t.Fatalf("Alice should have received fatal bullet")
	}
	if p1.IsAlive {
		t.Fatalf("Alice should be eliminated")
	}
	if g.State.Status != model.StatusGameOver || g.State.Winner != "Bob" {
		t.Fatalf("Bob should be declared winner, status=%s, winner=%s", g.State.Status, g.State.Winner)
	}
}

func TestCallLiarEmptyHandHonestVictory(t *testing.T) {
	g := NewGame("TEST04")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", Hand: []model.Card{}, IsAlive: true, Bullets: 6, Revolver: []string{model.BulletBlank}}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", Hand: []model.Card{model.Queen}, IsAlive: true, Bullets: 6, Revolver: []string{model.BulletBlank}}
	g.State.Players = []*model.Player{p1, p2}
	g.State.TableCard = model.Ace
	g.HiddenCards = []model.Card{model.Ace} // Alice played last card Ace (Honest)

	var shotTarget string
	var shotFatal bool

	g.CallLiar(1, 0, nil, nil, func(target string, fatal bool) {
		shotTarget = target
		shotFatal = fatal
	})

	if shotTarget != "Bob" {
		t.Fatalf("expected Bob (caller) to get shot first, got %s", shotTarget)
	}
	if shotFatal != false {
		t.Fatalf("expected blank shot for Bob")
	}

	if g.State.Status != model.StatusGameOver {
		t.Fatalf("expected game_over on empty hand honest win, got %s", g.State.Status)
	}
	if g.State.Winner != "Alice" {
		t.Fatalf("expected Alice to win, got %s", g.State.Winner)
	}
}

func TestPauseAndResumeGame(t *testing.T) {
	g := NewGame("TEST05")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", IsAlive: true, Bullets: 6}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", IsAlive: true, Bullets: 6}
	g.State.Players = []*model.Player{p1, p2}
	g.State.Status = model.StatusPlaying
	g.State.Deadline = time.Now().Add(25 * time.Second).Unix()

	// Alice disconnects -> Game pauses
	g.PauseGame(p1)

	if g.State.Status != model.StatusPaused {
		t.Fatalf("expected status paused, got %s", g.State.Status)
	}
	if g.State.PausedPlayer != "Alice" {
		t.Fatalf("expected PausedPlayer Alice, got %s", g.State.PausedPlayer)
	}
	if g.State.RemainingTurnSeconds < 20 || g.State.RemainingTurnSeconds > 26 {
		t.Fatalf("expected remaining seconds around 25, got %d", g.State.RemainingTurnSeconds)
	}

	// Alice reconnects within 30s -> Game resumes
	g.ResumeGame(p1)

	if g.State.Status != model.StatusPlaying {
		t.Fatalf("expected status playing after resume, got %s", g.State.Status)
	}
	if g.State.PausedPlayer != "" {
		t.Fatalf("expected empty PausedPlayer, got %s", g.State.PausedPlayer)
	}
	if g.State.Deadline <= time.Now().Unix() {
		t.Fatalf("expected deadline in future")
	}
}

func TestPauseTimeoutKill(t *testing.T) {
	g := NewGame("TEST06")
	p1 := &model.Player{ID: "p1", Nickname: "Alice", IsAlive: true, Bullets: 6, ClientRef: nil}
	p2 := &model.Player{ID: "p2", Nickname: "Bob", IsAlive: true, Bullets: 6, ClientRef: "connected"}
	g.State.Players = []*model.Player{p1, p2}
	g.State.Status = model.StatusPlaying
	g.State.CurrentTurn = 0

	g.PauseGame(p1)
	g.State.PauseDeadline = time.Now().Add(-1 * time.Second).Unix() // Expired

	g.HandlePauseTimeout()

	if p1.IsAlive {
		t.Fatalf("Alice should be eliminated after timeout")
	}
	if g.State.Status != model.StatusGameOver {
		t.Fatalf("expected game over with only Bob alive, got %s", g.State.Status)
	}
	if g.State.Winner != "Bob" {
		t.Fatalf("expected Bob to win, got %s", g.State.Winner)
	}
}


