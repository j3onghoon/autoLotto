package pkg

import (
	"github.com/gofiber/fiber/v3"
	"strconv"
	"strings"
)

func ResetUser(c fiber.Ctx) {
	refreshCookie := fiber.Cookie{
		Name:   "refresh_token",
		Value:  "",
		MaxAge: -1,
	}
	accessCookie := fiber.Cookie{
		Name:   "access_token",
		Value:  "",
		MaxAge: -1,
	}
	c.Cookie(&refreshCookie)
	c.Cookie(&accessCookie)
	c.Locals("lotteryUser", "")
}

func FormatWithCommas(numberStr string) string {
	// 문자열을 정수로 변환
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return ""
	}

	// 정수를 문자열로 포맷팅하여 1,000 단위로 콤마 삽입
	result := strconv.Itoa(number)
	if len(result) <= 3 {
		return result
	}

	// 빌더를 사용해 결과를 만들어냄
	var builder strings.Builder
	for i, v := range result {
		if (len(result)-i)%3 == 0 && i != 0 {
			builder.WriteString(",")
		}
		builder.WriteRune(v)
	}

	return builder.String()
}
