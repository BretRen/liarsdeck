package game

import (
	"math/rand"

	"pdnode.com/play/liarsbar-web/internal/model"
)

func NewDeck() []model.Card {
	deck := []model.Card{
		model.King, model.King, model.King, model.King, model.King, model.King,
		model.Queen, model.Queen, model.Queen, model.Queen, model.Queen, model.Queen,
		model.Ace, model.Ace, model.Ace, model.Ace, model.Ace, model.Ace,
		model.Two, model.Two, model.Two, model.Two, model.Two, model.Two,
	}
	ShuffleCards(deck)
	return deck
}

func ShuffleCards(cards []model.Card) {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func NewRevolver() []string {
	revolver := []string{
		model.BulletBlank,
		model.BulletBlank,
		model.BulletBlank,
		model.BulletBlank,
		model.BulletBlank,
		model.BulletFatal,
	}
	rand.Shuffle(len(revolver), func(i, j int) {
		revolver[i], revolver[j] = revolver[j], revolver[i]
	})
	return revolver
}

func DrawTableCard(deck []model.Card) (model.Card, []model.Card) {
	if len(deck) == 0 {
		return model.King, deck
	}
	tableCard := deck[0]
	deck = deck[1:]
	// 2 作为万能牌不能作为真牌
	for tableCard == model.Two && len(deck) > 0 {
		tableCard = deck[0]
		deck = deck[1:]
	}
	if tableCard == model.Two {
		tableCard = model.King
	}
	return tableCard, deck
}
