package templates

import (
	"bytes"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func Render(c fiber.Ctx, component templ.Component, options ...func(*templ.Component)) error {
	c.Set("Content-Type", "text/html")

	for _, option := range options {
		option(&component)
	}

	return component.Render(c.Context(), c.Response().BodyWriter())
}

func RenderToString(c fiber.Ctx, component templ.Component) (string, error) {
	var buf bytes.Buffer
	err := component.Render(c.Context(), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
