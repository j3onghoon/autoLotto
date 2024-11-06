package handlers

import (
	"autoLotto/pkg"
	"github.com/gofiber/fiber/v3"
)

func LoginHandler(c fiber.Ctx) error {
	lotteryUser, err := UserHandler(c)
	if err != nil {
		c.Redirect().Status(fiber.StatusSeeOther).To("/?alert=auth_fail")
		return nil
	}
	newRefreshToken, err := pkg.GenerateRefreshToken(lotteryUser.User.ID)
	refreshCookie := fiber.Cookie{
		Name:   "refresh_token",
		Value:  newRefreshToken,
		MaxAge: 3600 * 24 * 90,
	}
	c.Cookie(&refreshCookie)
	c.Redirect().Status(fiber.StatusSeeOther).To("/mypage")
	return nil
}
