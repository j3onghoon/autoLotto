package pkg

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v3"
	"time"
)

var jwtAccessSecret = []byte("your_secret_key")  // 비밀 키 설정
var jwtRefreshSecret = []byte("your_secret_key") // 비밀 키 설정

// 토큰 생성 함수
func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAccessSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtRefreshSecret)
}

// 리프레시 토큰을 이용해 새로운 액세스 토큰 발급
func RefreshAccessToken(c fiber.Ctx) string {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return ""
	}

	// 리프레시 토큰 검증
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshSecret, nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	// 새로운 액세스 토큰 발급
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)
	newAccessToken, _ := GenerateAccessToken(userID)

	accessCookie := fiber.Cookie{
		Name:   "access_token",
		Value:  newAccessToken,
		MaxAge: 900,
	}
	c.Cookie(&accessCookie)
	return newAccessToken
}

// 액세스 토큰 검증 미들웨어
func AuthMiddleware(c fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	refreshToken := c.Cookies("refresh_token")

	if accessToken == "" && refreshToken == "" {
		c.Redirect().Status(fiber.StatusSeeOther).To("/?alert=auth_fail")
		return nil
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtAccessSecret, nil
	})

	// 액세스 토큰이 만료되었거나 유효하지 않은 경우
	if err != nil || !token.Valid {
		// 리프레시 토큰을 이용해 새 액세스 토큰 발급
		accessToken = RefreshAccessToken(c)
		if accessToken == "" {
			c.Redirect().Status(fiber.StatusSeeOther).To("/?alert=auth_fail")
			return nil
		} else {
			token, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				return jwtAccessSecret, nil
			})
		}
	}
	if token == nil {
		c.Redirect().Status(fiber.StatusSeeOther).To("/?alert=auth_fail")
		return nil
	}

	var userID string
	// 클레임에서 userId 추출
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, _ = claims["userId"].(string)
	}

	if err != nil {
		return err
	}
	dbUser, err := GetUserByID(context.Background(), userID)
	lotteryUser, err := SetLottoUserInfo(dbUser.ID, dbUser.Password)
	if err != nil {
		ResetUser(c)
	}
	c.Locals("lotteryUser", *lotteryUser)

	// 다음 핸들러로 요청 전달
	c.Next()
	return nil
}
