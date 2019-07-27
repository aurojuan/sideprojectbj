package blackjack

/* 
分牌 賭倍等 涉及balance判定 不在這裡處理 寫在player.go
*/

// Hand :遊戲中每一手拿到的牌 也就是Deal後 每個player手上的牌型
type Hand struct {
	Cards []Card
} 

// HasCardA :軟點牌型判斷的前置工作
func (hand *Hand) HasCardA() bool {
	for _, v := range hand.Cards {
		if v.Symbol == "A" {
			return true
		}
	}
	return false
}

// PointCountWithA :軟點牌型"預"計算 也方便途中貼標籤 先把A獨立出來
// 算剩下的點數和(當中如有A 算成1) 假設為剩下總和為j
// 則此手牌為 [j+1, j+11]
// 實際遊戲時 如果PointCountWithA[1]沒爆 就取它
// 否則看PointCountWithA[0]爆了沒
func (hand *Hand) PointCountWithA() []int {
	j := 0
	for _, v := range hand.Cards {
		j += v.Point
	}
	return []int{j, j + 10}
}

// IsSoft :判斷是否軟點牌型
// 提醒一下 像A57就不是軟點 故除了A的有無 點數也要考慮
func (hand *Hand) IsSoft() bool {
	// 首先你要有妹妹(x)  要有A(O).....
	if !hand.HasCardA() {
		return false
	}
	
	// 無論是該函數返回的第一 還是第二個值 只要有超過21的
	// 表該手不再有彈性 是硬點牌 所有A的點數此後只能當1來算
	for _, v := range hand.PointCountWithA() {
		if v > 21 {
			return false
		}
	}
	return true
}

// SoftPointDetermine 軟點牌型計算 和PointCountWithA有差
// 請自行先判定軟點否
func (hand *Hand) SoftPointDetermine() []int {
	return hand.PointCountWithA()
} 

// HardPointDetermine :就是計算所有非軟點牌型的點數和
// 以及A只能當1點的case
func (hand *Hand) HardPointDetermine() int {
	j := 0
	for _, v := range hand.Cards {
		j += v.Point
	}
	return j
}

// IsBust :到時處理player和dealer能否繼續hit
func (hand *Hand) IsBust() bool {
	if hand.HardPointDetermine() > 21 {
		return true
	}
	return false
}

// IsBlackJack :判定player或dealer是否blackjack
func (hand *Hand) IsBlackJack() bool {
	if len(hand.Cards) != 2 || !hand.HasCardA() {
		return false
	}
	if hand.PointCountWithA()[1] != 21 {
		return false
	}
	return true
}

// CanSplitHand :一開始兩張牌一樣 才能給玩家有分牌這動作的選擇
// 並送相關data出去
// of course, player only :)
// 待考慮是否改放到player.go 或 round.go
// 也許這裡 也許另外別地方寫bet balance的過濾 
// 比方round可以先過濾bet level夠不夠 
// (double & insurance同理 )
func (hand *Hand) CanSplitHand() bool {
	// 只有最初剛發兩張牌時 才能選是否分牌
	if hand.Cards[0].Symbol == hand.Cards[1].Symbol {
		return true
	}
	return false
}

// IsNonBustAfterFive :五小龍判定
func (hand *Hand) IsNonBustAfterFive() bool {
	// 再檢查一次這手的張數
	if len(hand.Cards) < 5 {
		return false
	}
	return !hand.IsBust()	
}

// FinalPoint 之後round.go 最後大家比牌用 
// 自行排除bj 五小龍 bust
func (hand *Hand) FinalPoint() int {
	if hand.IsSoft() {
		return hand.SoftPointDetermine()[1]
	}
	return hand.HardPointDetermine()
}
