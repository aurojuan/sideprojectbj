package blackjack

import (
	"fmt"
)

// Player :目前規劃至少有Hand 餘額 投注等資料
// HandOfPlayer []Hand 心裡想的是player可能可以下不止一注
// Action先設計放這  可以選擇做的動作
// Decision 紀錄牌局中的實際動作
// Seat 多人時用
// 暫定每回合遊戲開始前 bets HandsOfPlayer acion 都初始化淨空
type Player struct {
	HandsOfPlayer []Hand
	Balence float64
	Bets []float64
	UserID string
	Action []string
	Decision []string
	InsuranceDecision []bool // true表要保
	InsuranceResult bool
	PlayerGameState []string // 標記bj dragon bust...用
	PointOfPlayer []int      // 非bj dragon bust時用 不止一手
	ResultOfPlayer []string
	Seat string // 座位編號
}

/* // InitialPlayer 初始化玩家
func (player *Player) InitialPlayer() *Player {
	// player := &Player{}
	player.HandsOfPlayer = make([]Hand, 2) // 最多分一次牌
	player.Balence = float64(1000) // 暫定
	player.Bets = make([]float64, 20)
	player.UserID = ""
	player.Action = make([]string, 4)
	player.Decision = make([]string, 4)
	player.InsuranceDecision = make([]bool, 2)
	player.PlayerGameState = make([]string, 4) // 一開始預設空字串 可能有分牌 用array
	player.PointOfPlayer = make([]int, 4)
	player.ResultOfPlayer = make([]string, 4) // 一開始預設空字串  可能有分牌 用array
	player.Seat = ""
	return player
} */

// ActionFilterByBalence :已下注 發牌後 直接依照balance 牌型
// 對可以做的action賦值
// 資產夠: 有對子 則可分可倍 hit stand  無對子 僅可倍 hit stand
// 資產不夠: 只能hit stand
// 分牌每席只能一次 所以依seat跑審定
// 注意:玩家對自己每個位置做決定後 資產可能是動態變化的
// HandsOfPlayer[0] 表剛發完2張牌後的那一手
// Acion array 目前應該是綁定了
func (player *Player) ActionFilterByBalence() {
	if player.Balence >= player.Bets[0] {
		if player.HandsOfPlayer[0].Cards[0].Symbol ==  player.HandsOfPlayer[0].Cards[1].Symbol {
			player.Action = []string{"split", "double", "hit", "stand"} // 暫改 710 you know
			//player.Action = []string{"double", "hit", "stand"}
		} else {
			player.Action = []string{"double", "hit", "stand"}
		}	
	} else {
		player.Action = []string{"hit", "stand"}
	}
}

// ActionFilterAfterSplit :直接把分牌後能做的寫在這
func (player *Player) ActionFilterAfterSplit() {
	player.Action = []string{"double", "hit", "stand"}
}

// ActionFilterAfterHit :double split只有在手上頭兩張時可決定
// 另split每人只有一次 分完可double hit stand  第一次選hit且沒爆 之後只能hit stand
func (player *Player) ActionFilterAfterHit() {
	player.Action = []string{"hit", "stand"}
}

// DecisionRecordForSplitandDouble :每個決策完 直接紀錄變化進struct player
// 返回動作紀錄 投注紀錄 餘額變化
// 不同決策 分case處理投注 餘額  預設player.Bets[0]是發牌前的下注額
// 有新bet動作可以進行再扣款
func (player *Player) DecisionRecordForSplitandDouble(decision string) {
	player.Decision = append(player.Decision, decision)

	// 依照決定 處理balance
	switch decision {
	case "split":
		// 資產夠才能分牌 上面ActionFilterByBalence()寫了能否動作的判斷
		player.Bets = append(player.Bets, player.Bets[0])
		player.Balence -= player.Bets[0]
		player.HandsOfPlayer = append(player.HandsOfPlayer, player.HandsOfPlayer[0]) // 先複製第一手
		// 為了避免slice reference連動問題 先拉出來令變數
		// 發牌的動作 在round.go裡處理
		tmpFirstHead := player.HandsOfPlayer[0].Cards[0]
		tmpSecondHead := player.HandsOfPlayer[0].Cards[1]
		player.HandsOfPlayer[0].Cards = []Card{tmpFirstHead}
		player.HandsOfPlayer[1].Cards = []Card{tmpSecondHead}
	case "double":
		player.Bets = append(player.Bets, player.Bets[0])
		player.Balence -= player.Bets[0]
	default:
		fmt.Println("can't split or double")
	}
}

// InsuranceFilterByBalence :決定你有資格執行保險否
// 對有資格的 才出現可以按保險的選項
func (player *Player) InsuranceFilterByBalence() bool {
	if player.Balence >= 0.5 * player.Bets[0] {
		return true
	}
	return false
}

// InsuranceRecord :假設要保險的處理
// insuranceDecision == true 表玩家要保
func (player *Player) InsuranceRecord(insuranceDecision bool) {
	// 有資格能保險且決定保險
	if player.InsuranceFilterByBalence() && insuranceDecision {
		player.InsuranceDecision[0] = true // 最初手才可保險
		// player.InsuranceDecision = append(player.InsuranceDecision, insuranceDecision)
		player.Bets = append(player.Bets, 0.5 * player.Bets[0])
		player.Balence -= player.Bets[0] * 0.5
	}
}

// ResultToPay :處理round.go裡的room.Player.ResultOfPlayer[i]拿到的各種標籤 
// 當然是有贏錢的寫入他的Balence屬性
func (player *Player) ResultToPay(previousBalence float64) {
	if player.InsuranceResult {
		player.Balence += 0.5 * player.Bets[0] + player.Bets[0] // 保險要獨立在輸贏結果之外處理
	}

	if len(player.Decision) != 0 && player.Decision[0] == "split" {
		for i := range player.HandsOfPlayer {
			switch player.ResultOfPlayer[i] {
			case "push":
				player.Balence += player.Bets[0] // 拿回本金
			/*case "bjWin":
				player.Balence += player.Bets[0] + 1.5 * player.Bets[0] // 本金拿了再拿獎*/
			case "dragonWin":
				player.Balence += player.Bets[0] + 1.5 * player.Bets[0] // 本金拿了再拿獎
			case "win":
				player.Balence += player.Bets[0] + player.Bets[0] // 本金拿了再拿獎
			case "doubleWin":
				player.Balence += 2 * player.Bets[0] + 2 * player.Bets[0] // 本金拿了再拿獎
			case "doublePush":
				player.Balence += 2 * player.Bets[0] // 拿回本金
			}
		}
	} else {
		switch player.ResultOfPlayer[0] {
		case "push":
			player.Balence += player.Bets[0] // 拿回本金
		case "bjWin":
			player.Balence += player.Bets[0] + 1.5 * player.Bets[0] // 本金拿了再拿獎
		case "dragonWin":
			player.Balence += player.Bets[0] + 1.5 * player.Bets[0] // 本金拿了再拿獎
		case "win":
			player.Balence += player.Bets[0] + player.Bets[0] // 本金拿了再拿獎
		case "doubleWin":
			player.Balence += 2 * player.Bets[0] + 2 * player.Bets[0] // 本金拿了再拿獎
		case "doublePush":
			player.Balence += 2 * player.Bets[0] // 拿回本金
		}
	}

	afterBalence := player.Balence
	tmpBalence := afterBalence - previousBalence
	if tmpBalence > 0 {
		player.Balence -= 0.05 * tmpBalence // 收水錢
	}
}

// ShowMoneyToPlayer :local端先採用print出金流變化的格式
func (player *Player) ShowMoneyToPlayer(previousBalence float64) {
	if player.Balence > previousBalence {
		fmt.Println("+", player.Balence - previousBalence)
	}
	if player.Balence <= previousBalence {
		fmt.Println(player.Balence - previousBalence)
	}
}

// CleanPlayerForNextRound ::下一回合前清除計算用的暫存data
func (player *Player) CleanPlayerForNextRound() {
	player.HandsOfPlayer = make([]Hand, 1) // 最多分一次牌
	player.Bets = make([]float64, 1)
	player.Action = make([]string, 4)
	player.Decision = make([]string, 0)
	player.InsuranceDecision = make([]bool, 1)
	player.InsuranceResult = false
	player.PlayerGameState = make([]string, 1) // 一開始預設空字串 可能有分牌 用array
	player.PointOfPlayer = make([]int, 1)
	player.ResultOfPlayer = make([]string, 1) // 一開始預設空字串  可能有分牌 用array
}

// IsPlayerBankrupt :多回合遊戲進行中 量度玩家GG了沒
func (player *Player) IsPlayerBankrupt(minBet float64) bool {
	if player.Balence < minBet {
		fmt.Println("balence below minimum level !!")
		return true
	}
	return false
}
