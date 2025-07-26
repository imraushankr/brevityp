package services

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/pkg/storage"
)

type userService struct {
	repo    interfaces.UserRepository
	auth    *auth.Auth
	email   *email.EmailService
	cfg     *configs.Config
	storage storage.Storage
	log     logger.Logger
}

func NewUserService(
	repo interfaces.UserRepository,
	auth *auth.Auth,
	email *email.EmailService,
	cfg *configs.Config,
	storage storage.Storage,
	log logger.Logger,
) interfaces.UserService {
	return &userService{
		repo:    repo,
		auth:    auth,
		email:   email,
		cfg:     cfg,
		storage: storage,
		log:     log,
	}
}

func (s *userService) FindUser(ctx context.Context, identifier string) (*models.User, error) {
	if identifier == "" {
		return nil, models.ErrInvalidInput
	}

	// First try to find by ID (for JWT authenticated requests)
	user, err := s.repo.FindUserByID(ctx, identifier)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, models.ErrUserNotFound) {
		return nil, err
	}

	// Fall back to email/username lookup
	user, err = s.repo.FindUserByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.log.Warn("User not found", logger.String("identifier", identifier))
		} else {
			s.log.Error("Failed to find user",
				logger.NamedError("error", err),
				logger.String("identifier", identifier))
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		s.log.Error("User validation failed", logger.NamedError("error", err))
		return err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.log.Error("Failed to update user",
			logger.NamedError("error", err),
			logger.String("user_id", user.ID))
		return err
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return models.ErrInvalidInput
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("Failed to delete user",
			logger.NamedError("error", err),
			logger.String("user_id", id))
		return err
	}

	return nil
}

func (s *userService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if userID == "" {
		return "", models.ErrInvalidInput
	}

	// Validate file
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

	// Upload to storage
	avatarURL, err := s.storage.UploadAvatar(ctx, userID, file, header)
	if err != nil {
		s.log.Error("Failed to upload avatar",
			logger.NamedError("error", err),
			logger.String("user_id", userID))
		return "", err
	}

	// Update user record
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		s.log.Error("Failed to update avatar URL",
			logger.NamedError("error", err),
			logger.String("user_id", userID))
		return "", err
	}

	return avatarURL, nil
}
