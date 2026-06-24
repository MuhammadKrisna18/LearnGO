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

// Me returns the profile of the currently logged-in user.
// @Summary Get User Profile
// @Description Returns the profile of the user based on the provided JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse{data=domain.UserProfileResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Extract userID from Locals (injected by JWT middleware)
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
