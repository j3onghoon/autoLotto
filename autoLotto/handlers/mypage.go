package handlers

import (
	"autoLotto/pkg"
	"autoLotto/templates"
	"github.com/gofiber/fiber/v3"
)

func MyPageHandler(c fiber.Ctx) error {
	lotteryUser := c.Locals("lotteryUser").(pkg.LotteryUser)
	ticketCountList := []string{"1", "2", "3", "4", "5"}

	templates.Render(c, templates.MyPage(lotteryUser, ticketCountList))
	return nil
}
