package game

import (
	cryptorand "crypto/rand"
	"math/big"

	"pdnode.com/play/liarsbar-web/internal/model"
)

// CryptoShuffle 使用加密安全的 Fisher-Yates 洗牌算法，彻底避免 math/rand 默认确定性种子问题
func CryptoShuffle(n int, swap func(i, j int)) {
	for i := n - 1; i > 0; i-- {
		jBig, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		j := int(jBig.Int64())
		swap(i, j)
	}
}

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
	CryptoShuffle(len(cards), func(i, j int) {
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
	CryptoShuffle(len(revolver), func(i, j int) {
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
