package handlers

import (
	"autoLotto/templates"
	"github.com/gofiber/fiber/v3"
)

func AuthHandler(c fiber.Ctx) error {
	accessToken := c.Cookies("accessToken")
	refreshToken := c.Cookies("refreshToken")

	if accessToken == "" && refreshToken == "" {
		templates.Render(c, templates.Index())
		//c.HTML(http.StatusOK, "", templates.Index())
	} else {
		c.Redirect().Status(fiber.StatusSeeOther).To("/mypage")
		//c.JSON(http.StatusFound, gin.H{"redirect": "/mypage"})
		//c.Redirect(http.StatusFound, "/mypage")
	}
	return nil
}
