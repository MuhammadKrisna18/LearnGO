package http

import (
	"github.com/gofiber/fiber/v2"
	"modular-monolith/internal/modules/auth/domain"
	"modular-monolith/internal/shared/response"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.Login)
}

// Login authenticates a user and returns a JWT token.
// @Summary User Login
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login Credentials"
// @Success 200 {object} response.APIResponse{data=domain.LoginResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "email and password are required", nil)
	}

	res, err := h.service.Login(c.UserContext(), req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}
		return response.Error(c, fiber.StatusInternalServerError, "login failed", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "login successful", res)
}
