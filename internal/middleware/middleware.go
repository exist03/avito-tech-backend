package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func RoleCheck(next fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Get("User-role") != "admin" {
			return ctx.SendStatus(http.StatusForbidden)
		}
		err := next(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
