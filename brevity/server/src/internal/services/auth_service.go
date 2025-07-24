package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

type authService struct {
	repo  interfaces.AuthRepository
	auth  *auth.Auth
	email *email.EmailService
	cfg   *configs.Config
	log   logger.Logger
}

func NewAuthService(
	repo interfaces.AuthRepository,
	auth *auth.Auth,
	email *email.EmailService,
	cfg *configs.Config,
	log logger.Logger,
) interfaces.AuthService {
	return &authService{
		repo:  repo,
		auth:  auth,
		email: email,
		cfg:   cfg,
		log:   log,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	existingUser, err := s.repo.FindUserByIdentifier(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, models.ErrEmailAlreadyExists
	}

	existingUser, err = s.repo.FindUserByIdentifier(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, models.ErrUsernameAlreadyExists
	}

	hashedPassword, err := auth.EncryptPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.RoleUser,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.auth.GenerateVerificationToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.repo.SaveVerificationToken(ctx, user.Email, token, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save verification token: %w", err)
	}

	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.cfg.App.BaseURL, token)
	if err := s.email.SendVerificationEmail(user.Email, verificationLink); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.FindUserByIdentifier(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsVerified {
		return nil, models.ErrUserNotVerified
	}

	if err := auth.IsPasswordCorrect(user.Password, req.Password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	tokens, err := s.auth.GenerateTokens(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.LoginResponse{
		User:        *user,
		AccessToken: tokens.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// In JWT, logout is handled client-side by discarding the token
	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	if err := s.repo.VerifyUser(ctx, token); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (s *authService) InitiatePasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.FindUserByIdentifier(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil // Don't reveal if user exists
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	token, err := s.auth.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	if err := s.repo.SaveResetToken(ctx, email, token, expiresAt); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/api/v1/auth/reset-password/%s", s.cfg.App.BaseURL, token)
	return s.email.SendPasswordResetEmail(email, resetLink)
}

func (s *authService) CompletePasswordReset(ctx context.Context, token, newPassword string) error {
	hashedPassword, err := auth.EncryptPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.ResetPassword(ctx, token, hashedPassword); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	claims, err := s.auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, models.ErrInvalidToken
	}

	tokens, err := s.auth.GenerateTokens(claims.UserId, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return &models.RefreshTokenResponse{
		AccessToken: tokens.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	user, err := s.repo.FindUserByIdentifier(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := auth.IsPasswordCorrect(user.Password, req.CurrentPassword); err != nil {
		return models.ErrInvalidCredentials
	}

	hashedPassword, err := auth.EncryptPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, userID, hashedPassword)
}