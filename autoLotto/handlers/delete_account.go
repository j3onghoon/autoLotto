package handlers

import (
	"autoLotto/pkg"
	"github.com/gofiber/fiber/v3"
)

func DeleteAccountHandler(c fiber.Ctx) error {
	lotteryUser := c.Locals("lotteryUser").(pkg.LotteryUser)
	err := pkg.DeleteUser(c.UserContext(), lotteryUser.User)
	pkg.ResetUser(c)
	return err
}
