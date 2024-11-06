package handlers

import (
	"autoLotto/ent"
	"autoLotto/pkg"
	"github.com/gofiber/fiber/v3"
)

func UserInfoHandler(c fiber.Ctx) error {
	var userForm ent.User
	err := c.Bind().Body(&userForm)
	if err != nil {
		return err
	}

	lotteryUser := c.Locals("lotteryUser").(pkg.LotteryUser)
	lotteryUser.User.TicketCount = userForm.TicketCount
	lotteryUser.User.WeeklyPurchase = userForm.WeeklyPurchase
	err = pkg.UpdateUser(c.UserContext(), lotteryUser.User)
	if err != nil {
		return err
	}
	return nil
}
