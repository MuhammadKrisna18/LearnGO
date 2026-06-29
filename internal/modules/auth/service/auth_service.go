package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"siakad-pro/config"
	"siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/shared/apperrors"
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

	var pendingEmail *string
	pendingReq, err := s.repo.GetPendingEmailRequestByUserID(ctx, id)
	if err == nil && pendingReq != nil {
		pendingEmail = &pendingReq.NewEmail
	}

	return &domain.UserProfileResponse{
		ID:             user.ID,
		Name:           user.Name,
		Nickname:       user.Nickname,
		NID:            user.NID,
		Email:          user.Email,
		Role:           user.Role,
		ProgramStudiID: user.ProgramStudiID,
		ProgramStudi:   user.ProgramStudi,
		PendingEmail:   pendingEmail,
		PhotoURL:       user.PhotoURL,
		CreatedAt:      user.CreatedAt,
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

	nidStr := fmt.Sprintf("%05d", rand.Intn(90000)+10000)

	newUser := &domain.User{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Email:          email,
		NID:            &nidStr,
		Password:       string(hashedPassword),
		Role:           "dosen",
		ProgramStudiID: &req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create dosen account")
	}

	return &domain.UserProfileResponse{
		ID:             newUser.ID,
		Name:           newUser.Name,
		Nickname:       newUser.Nickname,
		NID:            newUser.NID,
		Email:          newUser.Email,
		Role:           newUser.Role,
		ProgramStudiID: newUser.ProgramStudiID,
		CreatedAt:      newUser.CreatedAt,
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
			ID:             u.ID,
			Name:           u.Name,
			Nickname:       u.Nickname,
			NID:            u.NID,
			Email:          u.Email,
			Role:           u.Role,
			ProgramStudiID: u.ProgramStudiID,
			ProgramStudi:   u.ProgramStudi,
			PhotoURL:       u.PhotoURL,
			CreatedAt:      u.CreatedAt,
		})
	}
	return res, nil
}

func (s *authService) DeleteDosen(ctx context.Context, id string) error {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Akun dosen tidak ditemukan", err.Error())
	}
	if user.Role != "dosen" {
		return &apperrors.AppError{Code: 400, Message: "Akun ini bukan dosen"}
	}
	return s.repo.DeleteUser(ctx, id)
}

func (s *authService) UpdateProfile(ctx context.Context, id string, req domain.UpdateProfileRequest) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Nickname != "" {
		user.Nickname = &req.Nickname
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate profil", err.Error())
	}

	return s.GetProfile(ctx, id)
}

func (s *authService) UpdateProfilePhoto(ctx context.Context, id string, photoURL string) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	user.PhotoURL = &photoURL

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate foto profil", err.Error())
	}

	return s.GetProfile(ctx, id)
}

func (s *authService) RequestEmailChange(ctx context.Context, userID string, req domain.EmailChangeRequestPayload) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	if user.Email == req.NewEmail {
		return apperrors.NewBadRequest("Email baru tidak boleh sama dengan email saat ini")
	}

	existingUser, err := s.repo.GetByEmail(ctx, req.NewEmail)
	if err == nil && existingUser.ID != user.ID {
		return apperrors.NewBadRequest("Email sudah digunakan oleh pengguna lain")
	}

	existingReq, err := s.repo.GetPendingEmailRequestByUserID(ctx, userID)
	if err == nil && existingReq != nil {

		existingReq.NewEmail = req.NewEmail
		return s.repo.UpdateEmailChangeRequest(ctx, existingReq)
	}

	newReq := &domain.EmailChangeRequest{
		ID:       uuid.New().String(),
		UserID:   userID,
		NewEmail: req.NewEmail,
		Status:   "pending",
	}

	return s.repo.CreateEmailChangeRequest(ctx, newReq)
}

func (s *authService) GetPendingEmailRequests(ctx context.Context) ([]*domain.EmailChangeRequest, error) {
	return s.repo.GetAllPendingEmailRequests(ctx)
}

func (s *authService) ReviewEmailRequest(ctx context.Context, requestID string, approve bool) error {
	req, err := s.repo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return apperrors.NewNotFound("Request tidak ditemukan", err.Error())
	}
	if req.Status != "pending" {
		return apperrors.NewBadRequest("Request sudah diproses")
	}

	if approve {
		req.Status = "approved"
		
		user, err := s.repo.GetByID(ctx, req.UserID)
		if err != nil {
			return apperrors.NewNotFound("User tidak ditemukan", err.Error())
		}
		
		user.Email = req.NewEmail
		if err := s.repo.Update(ctx, user); err != nil {
			return apperrors.NewInternal("Gagal mengupdate email user", err.Error())
		}
	} else {
		req.Status = "rejected"
	}

	return s.repo.UpdateEmailChangeRequest(ctx, req)
}
