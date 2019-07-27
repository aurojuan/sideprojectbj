package blackjack

import (
	"math/rand"
	"time"
)

// Deck :組成單元為Card 心裡想的是8副牌
type Deck struct {
	Cards []Card
}

// 以下是為了初始化deck做的
var (
	suitHash = map[int]string {
		0: "s",
		1: "h",
		2: "c",
		3: "d",
	}

	symbolHash = map[int]string {
		0: "A",
		1: "2",
		2: "3",
		3: "4",
		4: "5",
		5: "6",
		6: "7",
		7: "8",
		8: "9",
		9: "T",
		10: "J",
		11: "Q",
		12: "K",
	}
)

// SetPropertyOfDeck :n表要對幾副牌賦值 (我們預設是8副)
// 把card該有的Symbol Suit Point AlternativePoint在此一次設定好
func (deck *Deck) SetPropertyOfDeck(n int) *Deck {
	for i := 0; i < 52*n; i++ {
		tmp := i % 13
		deck.Cards[i].Symbol = symbolHash[tmp]
	}
	
	for i := 0; i < 52*n; i++ {
		tmp := i % 4
		deck.Cards[i].Suit = suitHash[tmp]
	}

	// J Q K都是當10
	for i := 0; i < 52*n; i++ {
		tmp := i % 13
		if tmp <= 9 {
			deck.Cards[i].Point = tmp + 1
		} else {
			deck.Cards[i].Point = 10
		}
	}

	for i := 0; i < 52*n; i++ {
		tmp := i % 13
		if tmp == 0 {
			deck.Cards[i].AlternativePoint = 11
		} else {
			deck.Cards[i].AlternativePoint = deck.Cards[i].Point
		}
	}
	return deck
}

// InitializeDeck :每張牌設置好 並洗牌 一切準備好可以開始玩的狀態
func InitializeDeck(n int) *Deck {
	deck := NewDeck(n)
	deck.SetPropertyOfDeck(n).ShuffleDeck()
	return deck
}

// NewDeck :n為想玩的副數
func NewDeck(n int) *Deck {
	deck := &Deck{}
	deck.Cards = make([]Card, 52*n) // 特設
	return deck
}

// ShuffleDeck :寫自己喜歡的洗牌
func (deck *Deck) ShuffleDeck() {
	size := len(deck.Cards)

	pos := make([]Card, size)

	for i := 0; i < size; i++ {
		pos[i] = deck.Cards[i]
	}

	for i := 0; i < size; i++ {
		rand.Seed(int64(time.Now().UnixNano()))
		luckOne := rand.Intn(size)
		luckTwo := rand.Intn(size)

		temp := pos[luckOne]
		pos[luckOne] = pos[luckTwo]
		pos[luckTwo] = temp
	}
	// 把換好的拿來用
	for i := 0; i < size; i++ {
		deck.Cards[i] = pos[i]
	}
}

// Draw :啊就發牌 動作完成後也會記下牌盒變化
func (deck *Deck) Draw(n int) []Card {
	cards := make([]Card, n)
	copy(cards, deck.Cards[:n])
	deck.Cards = deck.Cards[n:]
	return cards
}

// Remove :use it when cheat is needed
func (deck *Deck) Remove(wanted []Card) {
	for _, w := range wanted {
		for i, t := range deck.Cards {
			if w == t {
				//fmt.Println("before remove", deck.cards[i])
				copy(deck.Cards[i:], deck.Cards[i+1:])
				//fmt.Println("after remove", deck.cards[i])
			}
		}
	}
}

// IsDeckUnderPenetration :用來檢查需要洗牌重設置否 可能不一定需要這機制
// n為副數 比照賭場的穿透率 只玩68%的deck 低於qualify的量(返回true)就洗牌重置
func (deck *Deck) IsDeckUnderPenetration(n int) bool {
	return float64(len(deck.Cards)) < float64(n) * 52 * (1 - 0.68)
}
