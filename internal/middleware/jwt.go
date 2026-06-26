package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"modular-monolith/internal/shared/apperrors"
)

func Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return apperrors.NewUnauthorized("missing jwt")
		}

		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "bearer ")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.NewUnauthorized("invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return apperrors.NewUnauthorized("invalid or expired jwt", err.Error())
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apperrors.NewUnauthorized("invalid jwt claims")
		}

		c.Locals("userID", claims["id"])
		c.Locals("userRole", claims["role"])

		return c.Next()
	}
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok || userRole != role {
			return apperrors.NewForbidden("you do not have permission to perform this action")
		}
		return c.Next()
	}
}
