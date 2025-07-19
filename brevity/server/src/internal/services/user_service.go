package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/pkg/storage"
	"github.com/imraushankr/bervity/server/src/internal/repository"
	"gorm.io/gorm"
)

type userService struct {
	userRepo repository.UserRepository
	auth     *auth.Auth
	email    *email.EmailService
	cfg      *configs.Config
	storage  storage.Storage
	log      logger.Logger
}

func NewUserService(
	userRepo repository.UserRepository,
	auth *auth.Auth,
	email *email.EmailService,
	cfg *configs.Config,
	storage storage.Storage,
	log logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		auth:     auth,
		email:    email,
		cfg:      cfg,
		storage:  storage,
		log:      log,
	}
}

func (s *userService) Register(ctx context.Context, user *models.User) error {
	s.log.Info("Registering user", logger.String("email", user.Email), logger.String("username", user.Username))

	existingUser, err := s.userRepo.FindUser(ctx, user)
	if err == nil {
		if existingUser.Email == user.Email {
			return models.ErrEmailAlreadyExists
		}
		return models.ErrUsernameAlreadyExists
	} else if !errors.Is(err, models.ErrUserNotFound) {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := auth.EncryptPassword(user.Password)
	if err != nil {
		s.log.Error("Failed to hash password", logger.NamedError("error", err),
			logger.String("email", user.Email))
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("Failed to create user record", logger.NamedError("error", err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.auth.GenerateVerificationToken(user.ID)
	if err != nil {
		s.log.Error("Failed to generate verification token", logger.NamedError("error", err),
			logger.String("userID", user.ID))
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	expireAt := time.Now().Add(24 * time.Hour)
	if err := s.userRepo.SaveVerificationToken(ctx, user.Email, token, expireAt); err != nil {
		s.log.Error("Failed to save verification token", logger.NamedError("error", err),
			logger.String("email", user.Email))
		return fmt.Errorf("failed to save verification token: %w", err)
	}

	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.cfg.App.BaseURL, token)
	if err := s.email.SendVerificationEmail(user.Email, verificationLink); err != nil {
		s.log.Error("Failed to send verification email", logger.NamedError("error", err),
			logger.String("email", user.Email))
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	s.log.Info("User registered successfully", logger.String("email", user.Email),
		logger.String("userID", user.ID))
	return nil
}

func (s *userService) Login(ctx context.Context, userId, password string) (*models.User, string, error) {
	s.log.Info("Logging attempt", logger.String("userId", userId))

	if userId == "" || password == "" {
		return nil, "", models.ErrInvalidInput
	}

	user := &models.User{Email: userId, Username: userId}
	foundUser, err := s.userRepo.FindUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, "", models.ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}

	if !foundUser.IsVerified {
		return nil, "", models.ErrUserNotVerified
	}

	if err := auth.IsPasswordCorrect(foundUser.Password, password); err != nil {
		return nil, "", models.ErrInvalidCredentials
	}

	token, err := s.auth.GenerateTokens(foundUser.ID, string(foundUser.Role))
	if err != nil {
		s.log.Error("Failed to generate token", logger.NamedError("error", err),
			logger.String("userID", foundUser.ID))
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.log.Info("User logged in successfully", logger.String("userID", foundUser.ID))
	return foundUser, token.AccessToken, nil
}

func (s *userService) FindUser(ctx context.Context, identifier string) (*models.User, error) {
	s.log.Info("Finding user", logger.String("identifier", identifier))

	if identifier == "" {
		return nil, models.ErrInvalidInput
	}

	user := &models.User{Email: identifier, Username: identifier}
	foundUser, err := s.userRepo.FindUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return foundUser, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	s.log.Info("Updating user", logger.String("userID", user.ID))

	if err := user.Validate(); err != nil {
		s.log.Error("User validation failed", logger.NamedError("error", err))
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	s.log.Info("Deleting user", logger.String("userID", id))

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, token string) error {
	s.log.Info("Verifying email with token")

	if token == "" {
		return models.ErrInvalidToken
	}

	if err := s.userRepo.VerifyUser(ctx, token); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrInvalidVerificationToken
		}
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (s *userService) InitiatePasswordReset(ctx context.Context, email string) error {
	s.log.Info("Initiating password reset", logger.String("email", email))

	if email == "" {
		return models.ErrInvalidInput
	}

	user := &models.User{Email: email}
	foundUser, err := s.userRepo.FindUser(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	token, err := s.auth.GeneratePasswordResetToken(foundUser.ID)
	if err != nil {
		s.log.Error("Failed to generate reset token", logger.NamedError("error", err),
			logger.String("userID", foundUser.ID))
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expireAt := time.Now().Add(1 * time.Hour)
	if err := s.userRepo.SaveResetToken(ctx, email, token, expireAt); err != nil {
		s.log.Error("Failed to save reset token", logger.NamedError("error", err),
			logger.String("email", email))
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", s.cfg.App.BaseURL, token)
	if err := s.email.SendPasswordResetEmail(email, resetLink); err != nil {
		s.log.Error("Failed to send reset email", logger.NamedError("error", err),
			logger.String("email", email))
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

func (s *userService) CompletePasswordReset(ctx context.Context, token, newPassword string) error {
	s.log.Info("Completing password reset")

	if token == "" || newPassword == "" {
		return models.ErrInvalidInput
	}

	if err := s.userRepo.ResetPassword(ctx, token, newPassword); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrInvalidResetToken
		}
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	s.log.Info("Refreshing token")

	if refreshToken == "" {
		return "", models.ErrInvalidToken
	}

	claims, err := s.auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", models.ErrInvalidToken
	}

	newToken, err := s.auth.GenerateTokens(claims.UserId, claims.Role)
	if err != nil {
		s.log.Error("Failed to generate new token", logger.NamedError("error", err))
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	return newToken.AccessToken, nil
}

func (s *userService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	s.log.Info("Uploading avatar", logger.String("userID", userID))

	defer file.Close()

	if userID == "" {
		return "", models.ErrInvalidInput
	}

	if header.Size > s.cfg.Storage.MaxAvatarSize {
		return "", models.ErrFileTooLarge
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[header.Header.Get("Content-Type")] {
		return "", models.ErrInvalidFileType
	}

	avatarURL, err := s.storage.UploadAvatar(ctx, userID, file, header)
	if err != nil {
		s.log.Error("Failed to upload avatar", logger.NamedError("error", err),
			logger.String("userID", userID))
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	if err := s.userRepo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		s.log.Error("Failed to update avatar URL", logger.NamedError("error", err),
			logger.String("userID", userID))
		return "", fmt.Errorf("failed to update avatar URL: %w", err)
	}

	return avatarURL, nil
}
