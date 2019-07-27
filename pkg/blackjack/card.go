package blackjack

// Card :suit為花色 symbol為輸出牌面值用 
// point為真正計算點數用 (A先默認為1)
// alternativePoint為拿到A時 軟點牌型計點用 (A=11 sometimes)
// 約定數字(symbol)在前 花色(suit)在後string小寫
// Symbol {A 2 3 4 5 6 7 8 9 T J Q K}
type Card struct {
	Symbol string
	Suit string
	Point int
	AlternativePoint int
}
