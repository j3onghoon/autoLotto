package handlers

import (
	"autoLotto/ent"
	"autoLotto/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

func UserHandler(c fiber.Ctx) (*pkg.LotteryUser, error) {
	var userForm ent.User
	if err := c.Bind().Body(&userForm); err != nil {
		c.Redirect().Status(fiber.StatusUnauthorized).To("/")
		return nil, err
	} else {
		lotteryUser, err := pkg.SetLottoUserInfo(userForm.ID, userForm.Password)
		if err != nil {
			return nil, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePublic,
				Meta: map[string]interface{}{
					"status": http.StatusUnauthorized,
				},
			}
		}
		return lotteryUser, nil
	}

	//if err := c.ShouldBind(&userForm); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return nil, err
	//} else {
	//	lotteryUser, err := pkg.SetLottoUserInfo(userForm.ID, userForm.Password)
	//	if err != nil {
	//		return nil, &gin.Error{
	//			Err:  err,
	//			Type: gin.ErrorTypePublic,
	//			Meta: map[string]interface{}{
	//				"status": http.StatusUnauthorized,
	//			},
	//		}
	//	}
	//	return lotteryUser, nil
	//}
}
