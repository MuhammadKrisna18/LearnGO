package http

import (
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/middleware"
	"modular-monolith/internal/modules/auth/domain"
	"modular-monolith/internal/shared/apperrors"
	"modular-monolith/internal/shared/response"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router, jwtSecret string) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.Login)
	authGroup.Get("/me", middleware.Protected(jwtSecret), h.Me)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if req.Email == "" || req.Password == "" {
		return apperrors.NewBadRequest("email and password are required")
	}

	res, err := h.service.Login(c.UserContext(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return apperrors.NewUnauthorized(err.Error())
		}
		return apperrors.NewInternal("login failed", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "login successful", res)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return apperrors.NewUnauthorized("unauthorized")
	}

	profile, err := h.service.GetProfile(c.UserContext(), userID)
	if err != nil {
		return apperrors.NewInternal("failed to get profile", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "success", profile)
}
