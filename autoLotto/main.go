package main

import (
	"autoLotto/handlers"
	"autoLotto/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	app := fiber.New()

	var err error
	err = pkg.ConnectPostgres()
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	app.Get("/", func(c fiber.Ctx) error {
		return handlers.AuthHandler(c)
		//return templates.Render(c, templates.Index())
	})
	app.Post("/", func(c fiber.Ctx) error {
		return handlers.LoginHandler(c)
	})

	app.Use(pkg.AuthMiddleware)

	app.Get("/mypage", func(c fiber.Ctx) error {
		return handlers.MyPageHandler(c)
	})

	app.Post("/mypage", func(c fiber.Ctx) error {
		return handlers.UserInfoHandler(c)
	})

	app.Post("/logout", func(c fiber.Ctx) error {
		return handlers.LogoutHandler(c)
	})

	app.Post("/delete_account", func(c fiber.Ctx) error {
		return handlers.DeleteAccountHandler(c)
	})

	log.Fatal(app.Listen(":8080"))

	//if lotteryUser.Deposit == "0" {
	//	//	충전
	//}
	// 1. 동행복권 user 정보를 저장할 데이터 구조 만들기 - id, password 입력받고 로그인되는지 확인하기
	// 2. 구매 예약 정보 - 몇개살지. 매주구매여부.
	// 3. 당첨 정보 확인 - 당첨되면 카톡이나 문자로 알림가기.
	// 4. 하루 전에 예치금 정보 확인 후 부족하면 알림 보내기

	//  로또를 일주일에 사용자가 5천원 구매
	//  번호는 자동구매로 진행
	//  0. 동행복권 로그인
	//  1. 잔액확인
	//  2. 돈 충전
	//  3. 구매
	//  4. 로그아웃

	//  스케쥴러는 다시 생각
	//  unbuntu에서 할 때는 cron으로 진행, 아니라면 다른 거 찾아야함
	//  일주일에 1번씩 매주 월요일 9시에 구매
	//	추가
	//	당첨됐을 때 카카오로 알려주는 기능 추가 <- 유료화 월 1000원
	//  앱푸시로 기본 알림으로 하기
	//  안드로이드 개발로 하기
	//  앱 시작시에 예치금 부족 알림 혹은, 당첨 사실 알려주기
}
