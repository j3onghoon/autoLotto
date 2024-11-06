package handlers

import (
	"autoLotto/pkg"
	"github.com/gofiber/fiber/v3"
)

func LogoutHandler(c fiber.Ctx) error {
	pkg.ResetUser(c)
	return nil
}
