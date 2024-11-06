package main

import (
	"autoLotto/ent/user"
	"autoLotto/pkg"
	"context"
	"github.com/gofiber/fiber/v3/log"
	"strconv"
)

const (
	Automatic = iota
)

func purchaseLotto() {
	//	자동구매가 켜져 있는 유저들 1000명씩 가져오기
	pageSize := 1000
	offset := 0
	count := 0

	for {
		users, err := pkg.DbClient.User.
			Query().
			Where(user.WeeklyPurchaseEQ(true)).
			Limit(pageSize).
			Offset(offset).
			All(context.Background())
		if err != nil {
			log.Fatalf("failed querying weekly purchase true users: %v", err)
		}
		if len(users) == 0 {
			log.Infof("purchase lotto %d times done.", count)
			break
		}
		for _, user := range users {
			lotteryUser, err := pkg.SetLottoUserInfo(user.ID, user.Password)
			if err != nil {
				log.Fatalf("failed setting lottery user: %v", err)
			}
			pkg.MakeLotto645Tickets(lotteryUser, strconv.Itoa(Automatic), "", lotteryUser.User.TicketCount)

			//err = pkg.BuyLottery(*lotteryUser)
			if err != nil {
				log.Fatalf("failed buying lottery user: %v", err)
			}
			count++
		}
		offset += pageSize
	}
}

func checkDeposit() {
	// todo
	// 예치금이 구매장수보다 부족한지 확인하고 부족하면 앱 푸시를 하게 할 것
}

func main() {
	var err error
	err = pkg.ConnectPostgres()
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	checkDeposit()
	purchaseLotto()
}
