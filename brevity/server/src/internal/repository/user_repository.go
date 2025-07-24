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

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		r.log.Error("User validation failed", logger.NamedError("error", err))
		return err
	}

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		r.log.Error("Failed to update user", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
	if err != nil {
		r.log.Error("Failed to delete user", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar", avatarURL).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		r.log.Error("Failed to update avatar", logger.NamedError("error", err))
		return err
	}
	return nil
}