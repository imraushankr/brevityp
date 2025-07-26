package repository

import (
	"context"
	"errors"

	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewUserRepository(db *gorm.DB, log logger.Logger) interfaces.UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) FindUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		r.log.Error("Failed to find user",
			logger.NamedError("error", err),
			logger.String("identifier", identifier))
		return nil, err
	}

	return &user, nil
}

// repository/user_repository.go
func (r *userRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		r.log.Error("Failed to find user by ID",
			logger.NamedError("error", err),
			logger.String("user_id", id))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		r.log.Error("User validation failed",
			logger.NamedError("error", err),
			logger.String("user_id", user.ID))
		return err
	}

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		r.log.Error("Failed to update user",
			logger.NamedError("error", err),
			logger.String("user_id", user.ID))
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return models.ErrInvalidInput
	}

	err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
	if err != nil {
		r.log.Error("Failed to delete user",
			logger.NamedError("error", err),
			logger.String("user_id", id))
		return err
	}

	return nil
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	if userID == "" || avatarURL == "" {
		return models.ErrInvalidInput
	}

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar", avatarURL).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn("User not found when updating avatar",
				logger.String("user_id", userID))
			return models.ErrUserNotFound
		}
		r.log.Error("Failed to update avatar",
			logger.NamedError("error", err),
			logger.String("user_id", userID))
		return err
	}

	return nil
}