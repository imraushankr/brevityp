package repository

import (
	"context"
	"errors"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"gorm.io/gorm"
)

type authRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewAuthRepository(db *gorm.DB, log logger.Logger) interfaces.AuthRepository {
	return &authRepository{db: db, log: log}
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		r.log.Error("User validation failed", logger.NamedError("error", err))
		return err
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.log.Error("Failed to create user", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *authRepository) FindUserByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("email = ? OR username = ?", identifier, identifier).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		r.log.Error("Failed to find user", logger.NamedError("error", err))
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		r.log.Error("Failed to find user by ID", logger.NamedError("error", err))
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"verification_token":      token,
			"verification_expires_at": expires,
		}).Error

	if err != nil {
		r.log.Error("Failed to save verification token", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *authRepository) VerifyUser(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("verification_token = ? AND verification_expires_at > ?", token, time.Now()).
		Updates(map[string]interface{}{
			"is_verified":             true,
			"verification_token":      nil,
			"verification_expires_at": nil,
		}).Error

	if err != nil {
		r.log.Error("Failed to verify user", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *authRepository) SaveResetToken(ctx context.Context, email, token string, expires time.Time) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"reset_password_token":      token,
			"reset_password_expires_at": expires,
		}).Error

	if err != nil {
		r.log.Error("Failed to save reset token", logger.NamedError("error", err))
		return err
	}
	return nil
}

func (r *authRepository) ResetPassword(ctx context.Context, token, newPassword string) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("reset_password_token = ? AND reset_password_expires_at > ?", token, time.Now()).
		Updates(map[string]interface{}{
			"password":                  newPassword,
			"reset_password_token":      nil,
			"reset_password_expires_at": nil,
		})

	if result.Error != nil {
		r.log.Error("Failed to reset password", logger.NamedError("error", result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrInvalidResetToken
	}
	return nil
}

func (r *authRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error

	if err != nil {
		r.log.Error("Failed to update password", logger.NamedError("error", err))
		return err
	}
	return nil
}
