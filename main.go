package main

import (
	"fmt"
	blackjack "bjcards/pkg/blackjack"
)

// var simulateRoom blackjack.Room
// 可以fix次數看 
// 也可以for player.Balence > minBet ...看何時破產
func main() {
	var simulateRoom blackjack.Room
	// simulateRoom.InitialNewRoom()

	simulateRoom.RoomID = "" //待改
	simulateRoom.MinBet = float64(3) //待改
	simulateRoom.MaxBet = float64(300) //待改
	simulateRoom.CashIn = make([]float64, 16)
	simulateRoom.CashOut = make([]float64, 16)
	simulateRoom.Deck = blackjack.InitializeDeck(8) // 8副牌
	
	tmpPlayer := &blackjack.Player{}
	tmpPlayer.HandsOfPlayer = make([]blackjack.Hand, 1) // 最多分一次牌
	tmpPlayer.Balence = float64(10000) // 暫定
	tmpPlayer.Bets = make([]float64, 1) // 改用append 分牌又分別賭倍 最多4次
	tmpPlayer.UserID = ""
	tmpPlayer.Action = make([]string, 4)
	tmpPlayer.Decision = make([]string, 0)
	tmpPlayer.InsuranceDecision = make([]bool, 1)
	tmpPlayer.InsuranceResult = false
	tmpPlayer.PlayerGameState = make([]string, 1) // 改用append 一開始預設空字串 可能有分牌 用array
	tmpPlayer.PlayerGameState = []string{""} // 710
	tmpPlayer.PointOfPlayer = make([]int, 1)   // 改用append
	tmpPlayer.ResultOfPlayer = make([]string, 1) // 改用append 一開始預設空字串  可能有分牌 用array
	tmpPlayer.ResultOfPlayer = []string{""} // 710
	tmpPlayer.Seat = ""

	simulateRoom.Player = tmpPlayer // 奧義分段初始化再接上

	tmpDealer := &blackjack.Dealer{}
	tmpDealer.HandsOfDealer = make([]blackjack.Hand, 1) // 莊家只有ㄧ手
	tmpDealer.DealerGameState = "" // 一開始預設空字串
	tmpDealer.PointOfDealer = 0
	tmpDealer.ResultOfDealer = ""

	simulateRoom.Dealer = tmpDealer // 奧義分段初始化再接上

	playtimes := 0

	for i := 0; i < 1000; i++ {
		if simulateRoom.Player.IsPlayerBankrupt(simulateRoom.MinBet) {
			fmt.Println("玩了", playtimes, "次後破產")
			fmt.Println(simulateRoom.Player.Balence)
			break
		}
		simulateRoom.ExecuteGame()
		playtimes++
	}

	fmt.Println("玩家餘額", simulateRoom.Player.Balence)
	simulateRoom.ShowRTP()
}
