package blackjack

import (
	"fmt"
	"math/rand"
	"time"
	"math"
)

// Dealer :
type Dealer struct {
	// Balence float64
	HandsOfDealer []Hand
	DealerGameState string // 標記bj dragon bust...用
	PointOfDealer int // 非bj dragon bust時用
	ResultOfDealer string
}

// Room :define 房型 最小注 最大注 金流 遊戲要用的牌 莊家 賭徒
type Room struct {
	RoomID string
	MinBet float64
	MaxBet float64
	// MaxSeat int
	// UsersInRoom int
	CashIn []float64 // 弄成array 可以統計玩任意有限局數
	CashOut []float64 // 弄成array 可以統計玩任意有限局數
	Deck *Deck
	Dealer *Dealer
	Player *Player
}

/* // InitialDealer 初始化莊家
func (dealer *Dealer) InitialDealer() *Dealer {
	// dealer := &Dealer{}
	dealer.HandsOfDealer = make([]Hand, 1) // 莊家只有ㄧ手
	dealer.DealerGameState = "" // 一開始預設空字串
	dealer.PointOfDealer = 0
	dealer.ResultOfDealer = ""
	return dealer
} */

/* // InitialNewRoom 初始化房間
func (room *Room) InitialNewRoom() *Room {
	// room := &Room{}
	room.RoomID = "" //待改
	room.MinBet = float64(3) //待改
	room.MaxBet = float64(300) //待改
	room.CashIn = make([]float64, 10000)
	room.CashOut = make([]float64, 10000)
	room.Deck = InitializeDeck(8) // 8副牌
	room.Dealer = room.Dealer.InitialDealer()
	room.Player = room.Player.InitialPlayer()

	return room
} */

// IsBetIllegle :檢查『單人』投注值有沒有落在所屬房間的最小到最大注之間(含)
func (room *Room) IsBetIllegle(bet float64) bool {
	if bet < room.MinBet || bet > room.MaxBet {
		return false
	}
	return true
}

// 假投注產生 local端驗證時 需要的
func fakeBetGenerator() float64 {
	rand.Seed(time.Now().UnixNano())
	x := float64(rand.Intn(100))
	return x
}

// 假動作產生器 從action array挑要做的決策 local端驗證時only
func fakeDecisionGeneraor(player *Player) string {
	indx := rand.Intn(len(player.Action))
	act := player.Action[indx]
	return act
}

// 假保險決定產生器 從action array挑要做的決策 local端驗證時only
func fakeInsuranceGeneraor(player *Player) bool {
	indx := rand.Intn(2)
	if indx == 0 {
		player.InsuranceDecision[0] = true
	} else {
		player.InsuranceDecision[0] = false
	}
	return player.InsuranceDecision[0]
}

// decisionIsHit :簡化遊戲流程用 特別抽出來寫
// 此時玩家手上至少有兩張牌 indexOfHand 抽出來方便分牌的case也能用 自行指定index
// pointer寫入info
func (room *Room) decisionIsHit(indexOfHand int) {
	room.Player.HandsOfPlayer[indexOfHand].Cards = append(room.Player.HandsOfPlayer[indexOfHand].Cards, room.Deck.Draw(1)[0])
}

// decisionIsDouble :簡化遊戲流程用 特別抽出來寫
// 此時玩家手上至少有兩張牌 indexOfHand 抽出來方便分牌的case也能用 自行指定index
func (room *Room) decisionIsDouble(indexOfHand int) {
	room.Player.HandsOfPlayer[indexOfHand].Cards = append(room.Player.HandsOfPlayer[indexOfHand].Cards, room.Deck.Draw(1)[0])
}

// ExecuteGame :實際流程的邏輯寫在這
// return項目再計  initial另找地方寫
// 用endgame變數來標示流程結束否
func (room *Room) ExecuteGame() {
	
	// 每局重新洗牌  學開封...開元
	room.Deck = InitializeDeck(8)
	room.Deck.ShuffleDeck()
	
	previousBalence := room.Player.Balence // 遊戲開始前的餘額 後面計算會用到 因為遊戲中間餘額已有變
	
	// 先收bet 再用上前面的bet check  balence也要扣
	// 玩家不能下超過自己身家的投注檢查先寫在這 用判斷式 之後獨立寫 TBA
	baseBet := fakeBetGenerator()
	// baseBet := float64(5)
	if room.IsBetIllegle(baseBet) && baseBet <= room.Player.Balence {
		room.Player.Bets[0] = baseBet
		room.Player.Balence -= baseBet
	} else {
		room.Player.Bets[0] = room.MinBet // 沒有就預設最低注
		room.Player.Balence -= room.MinBet
	}
	 
	// 先依序發前兩張
	room.Dealer.HandsOfDealer[0].Cards = append(room.Dealer.HandsOfDealer[0].Cards, room.Deck.Draw(1)[0])
	room.Player.HandsOfPlayer[0].Cards = append(room.Player.HandsOfPlayer[0].Cards, room.Deck.Draw(1)[0])

	// 莊家第二張 到時請前端做成暗牌
	room.Dealer.HandsOfDealer[0].Cards = append(room.Dealer.HandsOfDealer[0].Cards, room.Deck.Draw(1)[0])
	room.Player.HandsOfPlayer[0].Cards = append(room.Player.HandsOfPlayer[0].Cards, room.Deck.Draw(1)[0])

	endgame := false // 用來標示處理玩家的程序 
	cardcount := 2  // 用來抓一直hit時 牌發了五張沒
	cardcountOfSplit := []int{2, 2} // 原理類上
	cardcountOfDealer := 2 // 原理類上
	
	for endgame == false {
		// 先看莊家明牌有無A 有就先問保險 問完檢視bj否
		if room.Dealer.HandsOfDealer[0].Cards[0].Symbol == "A" {
			// 有bj不用保 先拿標籤
			if room.Player.HandsOfPlayer[0].IsBlackJack() {
				room.Player.PlayerGameState[0] = "blackjack"
				endgame = true
				break
			}
			// 先篩保險資格
			if room.Player.InsuranceFilterByBalence() {
				// 之後也許要加個接收保險訊號過程

				// 模擬玩家要執行保險
				if fakeInsuranceGeneraor(room.Player) {
					// 錢先扣起來  下面decision那個待改
					room.Player.InsuranceRecord(true)
				}
			
			}
			// 看看dealer bj否 非bj才往下
			if room.Dealer.HandsOfDealer[0].IsBlackJack() {
				room.Dealer.DealerGameState = "blackjack"
				if room.Player.InsuranceDecision[0] != true {
					room.Player.PlayerGameState[0] = "needPoint" // 給沒買保險的人
				}
				// 暫定出去後 直接用InsuranceDecision []bool去看 出去再結算
				endgame = true
				break
			}
		}
		// 開始看玩家 應先check bj 
		if room.Player.HandsOfPlayer[0].IsBlackJack() {
			room.Player.PlayerGameState[0] = "blackjack"
			endgame = true
			break
		}
		// 先寫入可選動作 才能生成決策
		room.Player.ActionFilterByBalence()

		// 暫定依序設計分牌 賭倍 停 hit
		// 接收訊息的部分先用假的生成 並返回訊息部分(TBA)
		fakeDec := fakeDecisionGeneraor(room.Player)
		
		// 分牌
		if fakeDec == "split" {
			// bet balance 分成兩手處理
			// room.Player.PlayerGameState = "split"
			room.Player.DecisionRecordForSplitandDouble("split") // 處理hand bet
			// 每手分別發牌 711改append
			room.Player.PlayerGameState = append(room.Player.PlayerGameState, "")
			room.Player.PointOfPlayer = append(room.Player.PointOfPlayer, 0)
			room.Player.ResultOfPlayer = append(room.Player.ResultOfPlayer, "")
			
			room.Player.HandsOfPlayer[0].Cards = append(room.Player.HandsOfPlayer[0].Cards, room.Deck.Draw(1)[0])
			room.Player.HandsOfPlayer[1].Cards = append(room.Player.HandsOfPlayer[1].Cards, room.Deck.Draw(1)[0])

			// 之後每手只能double hit stand
			room.Player.ActionFilterAfterSplit()
			// info action TBA
			// 分兩手看
			for i := 0; i < 2; i++ {
				fakeDecAfterSplit := fakeDecisionGeneraor(room.Player)

				switch fakeDecAfterSplit {
				case "double":
					room.Player.PlayerGameState[i] = "double"
					room.Player.DecisionRecordForSplitandDouble("double")
					// 只發一張就停了
					room.Player.HandsOfPlayer[i].Cards = append(room.Player.HandsOfPlayer[i].Cards, room.Deck.Draw(1)[0])
					// 先看bust否 可以先label
					if room.Player.HandsOfPlayer[i].IsBust() {
						room.Player.PlayerGameState[i] = "bust"
					} else {
						room.Player.PlayerGameState[i] = "doubleNeedPoint"
					}
				case "stand":
					room.Player.PlayerGameState[i] = "needPoint"
					continue
				case "hit":
					splitHitStop := false // 老招 先用吧
			
					for splitHitStop == false {
						room.decisionIsHit(i) // 沒分牌 第i手加牌
						cardcountOfSplit[i]++
						if cardcountOfSplit[i] == 5 {
							if room.Player.HandsOfPlayer[i].IsNonBustAfterFive() {
								// 五小龍 另外給label
								room.Player.PlayerGameState[i] = "dragon"
								splitHitStop = true
								break
							}
							// 非五小龍表bust 出去認賠
							room.Player.PlayerGameState[i] = "bust"
							splitHitStop = true
							break
						}
						// 如果bust當然要跳出
						if room.Player.HandsOfPlayer[i].IsBust() {
							room.Player.PlayerGameState[i] = "bust"
							splitHitStop = true
							break
						}
						// 計數 local是出去做 
						room.Player.ActionFilterAfterHit()
						fakDecAfterSPHit := fakeDecisionGeneraor(room.Player)
						if fakDecAfterSPHit == "stand" {
							room.Player.PlayerGameState[i] = "needPoint"
							splitHitStop = true
							break
						}
					}
				}

			}
			endgame = true
			break
		}

		// 賭倍 按規則只再加發一張
		if fakeDec == "double" {
			room.Player.PlayerGameState[0] = "double"
			// 處理balance bet
			room.Player.DecisionRecordForSplitandDouble("double")
			// 只發一張就停了
			room.Player.HandsOfPlayer[0].Cards = append(room.Player.HandsOfPlayer[0].Cards, room.Deck.Draw(1)[0])
			// 先看bust否 可先給label
			if room.Player.HandsOfPlayer[0].IsBust() {
				room.Player.PlayerGameState[0] = "bust"
			} else {
				room.Player.PlayerGameState[0] = "doubleNeedPoint"
			}
			// 計數待補 然後跳出去和莊家比
			endgame = true
			break
		}

		// 選停牌case
		if fakeDec == "stand" {
			room.Player.PlayerGameState[0] = "needPoint"
			// 計數待捕 然後跳出去和莊家比
			endgame = true
			break
		}

		if fakeDec == "hit" {
			hitStop := false // 老招 先用吧
			
			for hitStop == false {
				room.decisionIsHit(0) // 沒分牌 第0手加牌
				cardcount++
				if cardcount == 5 {
					if room.Player.HandsOfPlayer[0].IsNonBustAfterFive() {
						// 五小龍達成給label
						room.Player.PlayerGameState[0] = "dragon"
						hitStop = true
						break
					}
					// bust出去認賠
					room.Player.PlayerGameState[0] = "bust"
					hitStop = true
					break
				}
				// 如果bust當然要跳出 TBA
				if room.Player.HandsOfPlayer[0].IsBust() {
					room.Player.PlayerGameState[0] = "bust"
					hitStop = true
					break
				}
				// 計數 TBA 
				room.Player.ActionFilterAfterHit()
				fakDecAfterHit := fakeDecisionGeneraor(room.Player)
				if fakDecAfterHit == "stand" {
					room.Player.PlayerGameState[0] = "needPoint"
					hitStop = true
					break
				}
			}
			
			endgame = true
			break
		}

	}
	// 跳出來後 才計算金額 輸贏
	// 開始處理莊家 當然頭兩張也不可能爆
	// 檢查可能先T 再A 或是底牌A 但觀眾bj先跳出 而因此沒抓到的bj
	if room.Dealer.HandsOfDealer[0].IsBlackJack() {
		room.Dealer.DealerGameState = "blackjack"
	}

	for room.Dealer.HandsOfDealer[0].FinalPoint() < 17 && cardcountOfDealer < 5 && !room.Dealer.HandsOfDealer[0].IsBlackJack() {
		room.Dealer.HandsOfDealer[0].Cards = append(room.Dealer.HandsOfDealer[0].Cards, room.Deck.Draw(1)[0])
		cardcountOfDealer++
	}
	
	// 回頭看莊家 小心被覆蓋到
	if !room.Dealer.HandsOfDealer[0].IsBlackJack() {
		if room.Dealer.HandsOfDealer[0].IsBust() {
			room.Dealer.DealerGameState = "bust"
		} else if room.Dealer.HandsOfDealer[0].IsNonBustAfterFive() {
			room.Dealer.DealerGameState = "dragon"
		} else {
			room.Dealer.DealerGameState = "needPoint"
		}
	}
	
	// 再來針對莊家和玩家 只要標needPoint的手牌 計算點數 請看下面
	// 注意到PlayerGameState這array有多長 表示玩家會有相應的手數 莊家當然只會有一手
	if room.Dealer.DealerGameState == "needPoint" {
		room.Dealer.PointOfDealer = room.Dealer.HandsOfDealer[0].FinalPoint()
	}

	// 玩家每一手算點
	for i := range room.Player.PlayerGameState {
		if room.Player.PlayerGameState[i] == "needPoint" || room.Player.PlayerGameState[i] == "doubleNeedPoint" {
			room.Player.PointOfPlayer[i] = room.Player.HandsOfPlayer[i].FinalPoint()
		}
	}

	room.InsuranceProcessing() // 莊bj 玩家有保才會成立
	room.JudgeWinPushLoss() // 開始決鬥惹  先定輸贏派標籤 然後金流 從標籤分類一一決定結果
	
	// 有贏錢的寫進balaence屬性
	room.Player.ResultToPay(previousBalence)

	fmt.Println("==============================")
	// 秀雙方彼此的手牌 local測print出來即可 實際要改回傳的form
	//fmt.Println("莊家手牌:", room.Dealer.HandsOfDealer) //先mark
	printdealerv := make([]string, 1)
	for i := range room.Dealer.HandsOfDealer[0].Cards {
		printdealera := room.Dealer.HandsOfDealer[0].Cards[i].Symbol
		printdealerb := room.Dealer.HandsOfDealer[0].Cards[i].Suit
		printdealer := printdealera + printdealerb
		printdealerv = append(printdealerv, printdealer)

		// printdealerc := room.Dealer.HandsOfDealer[0].Cards[i].Point
		
	}
	fmt.Println("文字化莊家手牌", printdealerv)
	
	//fmt.Println("玩家手牌:", room.Player.HandsOfPlayer) //先mark
	printplayerv := make([]string, 1)
	for i := range room.Player.HandsOfPlayer[0].Cards {
		printplayera := room.Player.HandsOfPlayer[0].Cards[i].Symbol
		printplayerb := room.Player.HandsOfPlayer[0].Cards[i].Suit
		printplayer := printplayera + printplayerb
		printplayerv = append(printplayerv, printplayer)

		// printdealerc := room.Dealer.HandsOfDealer[0].Cards[i].Point
		
	}
	fmt.Println("文字化玩家手牌", printplayerv)

	// 如果有分牌
	if len(room.Player.HandsOfPlayer) == 2 {
		printplayerw := make([]string, 1)
	
		for i := range room.Player.HandsOfPlayer[1].Cards {
			printplayerc := room.Player.HandsOfPlayer[1].Cards[i].Symbol
			printplayerd := room.Player.HandsOfPlayer[1].Cards[i].Suit
			printplayersp := printplayerc + printplayerd
			printplayerw = append(printplayerw, printplayersp)

		// printdealerc := room.Dealer.HandsOfDealer[0].Cards[i].Point
		
		}
		fmt.Println("文字化玩家分牌後另一手牌", printplayerw)
	}

	fmt.Println("玩家投注記錄:", room.Player.Bets)
	// 莊家 玩家輸贏紀錄 保險與否
	// fmt.Println("莊家輸贏", room.Dealer.DealerGameState) // 待改 result那個沒用到 可在下面改
	if room.Dealer.DealerGameState == "needPoint" {
		fmt.Println("莊家輸贏", room.Dealer.PointOfDealer)
	} else {
		fmt.Println("莊家輸贏", room.Dealer.DealerGameState)
	}

	if len(room.Player.ResultOfPlayer) == 2 {
		fmt.Println("玩家輸贏", room.Player.ResultOfPlayer[0], room.Player.ResultOfPlayer[1])
	} else {
		fmt.Println("玩家輸贏", room.Player.ResultOfPlayer[0])
	}
	
	if room.Player.InsuranceDecision[0] == true {
		if room.Player.InsuranceResult {
			fmt.Println("insurance Success")
		} else {
			fmt.Println("insurance fail!")
		}
	}
	
	// show金流紀錄給觀眾看
	room.Player.ShowMoneyToPlayer(previousBalence)
	room.CashFlowstatistic(previousBalence)
	fmt.Println("==============================")

	// 玩家破產通知和玩的次數(在main.go那弄)
	// 每單回合結束 清 player dealer一些暫時用的紀錄
	room.Dealer.CleanDealerForNextRound()
	room.Player.CleanPlayerForNextRound()
}

// InsuranceProcessing :單獨處理遊戲中 有玩家選擇買保險的後續處理
func (room *Room) InsuranceProcessing() {
	if room.Player.InsuranceDecision[0] && room.Dealer.DealerGameState == "blackjack" {
		room.Player.InsuranceResult = true // 預設false 只有保成功 才true
	}
}

// JudgeWinPushLoss :根據標籤分類決鬥 保險在上個func處理不在此弄
func (room *Room) JudgeWinPushLoss() {
	switch room.Dealer.DealerGameState {
	case "blackjack":
		for i := range room.Player.PlayerGameState {
			if room.Player.PlayerGameState[i] == "blackjack" {
				room.Player.ResultOfPlayer[i] = "push"
			} else {
				room.Player.ResultOfPlayer[i] = "loss"
			}	
		}
	case "dragon":
		for i := range room.Player.PlayerGameState {
			if room.Player.PlayerGameState[0] == "blackjack" {
				room.Player.ResultOfPlayer[0] = "bjWin" // 最初那手才可能bj
			} else if room.Player.PlayerGameState[i] == "dragon" {
				room.Player.ResultOfPlayer[i] = "push"
			} else {
				room.Player.ResultOfPlayer[i] = "loss"
			}	
		}
	case "bust":
		for i := range room.Player.PlayerGameState {
			if room.Player.PlayerGameState[i] == "bust" {
				room.Player.ResultOfPlayer[i] = "loss"
			} else if room.Player.PlayerGameState[i] == "dragon" {
				room.Player.ResultOfPlayer[i] = "dragonWin" // 賠率3:2
			} else if room.Player.PlayerGameState[0] == "blackjack" {
				room.Player.ResultOfPlayer[0] = "bjWin" // 最初那手才可能bj
			} else if room.Player.PlayerGameState[i] == "doubleNeedPoint" {
				room.Player.ResultOfPlayer[i] = "doubleWin"
			} else {
				room.Player.ResultOfPlayer[i] = "win"
			}	
		}
	case "needPoint":
		for i := range room.Player.PlayerGameState {
			if room.Player.PlayerGameState[0] == "blackjack" {
				room.Player.ResultOfPlayer[0] = "bjWin" // 最初那手才可能bj
			} else if room.Player.PlayerGameState[i] == "dragon" {
				room.Player.ResultOfPlayer[i] = "dragonWin" // 賠率3:2
			} else if room.Player.PlayerGameState[i] == "bust" {
				room.Player.ResultOfPlayer[i] = "loss"
			} else if room.Player.PlayerGameState[i] == "doubleNeedPoint" {
				if room.Player.PointOfPlayer[i] == room.Dealer.PointOfDealer {
					room.Player.ResultOfPlayer[i] = "doublePush"
				} else if room.Player.PointOfPlayer[i] > room.Dealer.PointOfDealer {
					room.Player.ResultOfPlayer[i] = "doubleWin"
				} else {
					room.Player.ResultOfPlayer[i] = "doubleLoss"
				}
			} else if room.Player.PlayerGameState[i] == "needPoint" {
				// 大家來比點數
				if room.Player.PointOfPlayer[i] == room.Dealer.PointOfDealer {
					room.Player.ResultOfPlayer[i] = "push"
				} else if room.Player.PointOfPlayer[i] > room.Dealer.PointOfDealer {
					room.Player.ResultOfPlayer[i] = "win"
				} else {
					room.Player.ResultOfPlayer[i] = "loss"
				}
			}	
		}
	}
}

// CashFlowstatistic :
// cashin cashout這邊用玩家的觀點看  方便算rtp
func (room *Room) CashFlowstatistic(previousBalence float64) {
	var cashIN float64
	var cashOut float64
	
	for _, v := range room.Player.Bets {
		cashIN += v
	}
	
	tmpPremium := room.Player.Balence - previousBalence
	if tmpPremium >= 0 {
		cashOut = tmpPremium + cashIN // push和贏錢的場合 含原本下注的總本金
	} else {
		if math.Abs(tmpPremium) == cashIN {
			cashOut = 0
		} else {
			cashOut = math.Abs(math.Abs(tmpPremium)- cashIN)
		}
	}
	
	fmt.Println("cash in:", cashIN)
	fmt.Println("cash out:", cashOut)
	room.CashIn = append(room.CashIn, cashIN)
	room.CashOut = append(room.CashOut, cashOut)
}

// ShowRTP :計算RTP  玩若干回合後使用
// main.go使用
func (room *Room) ShowRTP() {
	var moneyIn float64
	var moneyOut float64
	var rtp float64

	for _, v := range room.CashIn {
		moneyIn += v
	}
	fmt.Println("total cash in:", moneyIn)
	
	for _, v := range room.CashOut {
		moneyOut += v
	}
	fmt.Println("total cash out:", moneyOut)

	rtp = moneyOut / moneyIn
	fmt.Println("RTP is :", rtp)
}

// CleanDealerForNextRound :下一回合前清除計算用的暫存data
func (dealer *Dealer) CleanDealerForNextRound() {
	dealer.HandsOfDealer = make([]Hand, 1) // 莊家只有ㄧ手
	dealer.DealerGameState = "" // 一開始預設空字串
	dealer.PointOfDealer = 0
	dealer.ResultOfDealer = ""
}
