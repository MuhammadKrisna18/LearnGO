package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"modular-monolith/config"
	"modular-monolith/internal/modules/auth/domain"
	"modular-monolith/internal/shared/apperrors"
)

type authService struct {
	repo domain.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config) domain.AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &domain.LoginResponse{
		Token: t,
		Role:  user.Role,
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, id string) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) RegisterDosen(ctx context.Context, req domain.RegisterDosenRequest) (*domain.UserProfileResponse, error) {
	email := req.Username + "@DosenGO.id"

	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &domain.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "dosen",
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create dosen account")
	}

	return &domain.UserProfileResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

func (s *authService) GetDosenList(ctx context.Context) ([]*domain.UserProfileResponse, error) {
	users, err := s.repo.GetUsersByRole(ctx, "dosen")
	if err != nil {
		return nil, errors.New("failed to fetch dosen list")
	}

	var res []*domain.UserProfileResponse
	for _, u := range users {
		res = append(res, &domain.UserProfileResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}
	return res, nil
}

func (s *authService) DeleteDosen(ctx context.Context, id string) error {
	// Optional: verify if user exists and is a dosen
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Akun dosen tidak ditemukan", err.Error())
	}
	if user.Role != "dosen" {
		return &apperrors.AppError{Code: 400, Message: "Akun ini bukan dosen"}
	}
	return s.repo.DeleteUser(ctx, id)
}
