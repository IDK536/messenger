package middleware

import (
	"messanger/internal/config"

	tokenS "messanger/internal/services/token"

	"strings"

	"github.com/gofiber/fiber/v3"
)

func NewAuth(cfg *config.Config, tokenService *tokenS.TokenService) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "authorization header is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
		}

		tokenString := parts[1]

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "token is empty")
		}

		payload, err := tokenService.ValidateAccessToken(tokenString)

		if err != nil {
			return fiber.NewError(fiber.ErrUnauthorized.Code, err.Error())
		}

		c.Locals("user_id", payload.Id)
		return c.Next()
	}
}
