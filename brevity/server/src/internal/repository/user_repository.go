// package repository

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"gorm.io/gorm"
// )

// type userRespository struct {
// 	db  *gorm.DB
// 	log logger.Logger
// }

// func NewUserRepository(db *gorm.DB, log logger.Logger) UserRepository {
// 	return &userRespository{
// 		db:  db,
// 		log: log,
// 	}
// }

// func (r *userRespository) Create(ctx context.Context, user *models.User) error {
// 	if err := user.Validate(); err != nil {
// 		r.log.Error("User validation failed", logger.NamedError("error", err))
// 		return err
// 	}

// 	err := r.db.WithContext(ctx).Create(user).Error
// 	if err != nil {
// 		r.log.Error("Failed to create user", logger.NamedError("error", err))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) FindUser(ctx context.Context, user *models.User) (*models.User, error) {
// 	var existingUser models.User
// 	query := r.db.WithContext(ctx).Where("email = ? OR username = ?", user.Email, user.Username)

// 	err := query.First(&existingUser).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, models.ErrUserNotFound
// 	}

// 	if err != nil {
// 		r.log.Error("Failed to find user", logger.NamedError("error", err),
// 			logger.String("email", user.Email), logger.String("username", user.Username))
// 		return nil, err
// 	}

// 	return &existingUser, nil
// }

// func (r *userRespository) Update(ctx context.Context, user *models.User) error {
// 	r.log.Debug("Updating user", logger.String("userID", user.ID))

// 	err := r.db.WithContext(ctx).Save(user).Error
// 	if err != nil {
// 		r.log.Error("Failed to update user", logger.NamedError("error", err),
// 			logger.String("userID", user.ID))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) Delete(ctx context.Context, id string) error {
// 	r.log.Debug("Deleting user", logger.String("userID", id))

// 	err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
// 	if err != nil {
// 		r.log.Error("Failed to delete user", logger.NamedError("error", err),
// 			logger.String("userID", id))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error {
// 	r.log.Debug("Saving verification token", logger.String("email", email))

// 	err := r.db.WithContext(ctx).Model(&models.User{}).
// 		Where("email = ?", email).
// 		Updates(map[string]interface{}{
// 			"verification_token":      token,
// 			"verification_expires": expires,
// 		}).Error

// 	if err != nil {
// 		r.log.Error("Failed to save verification token", logger.NamedError("error", err),
// 			logger.String("email", email))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) VerifyUser(ctx context.Context, token string) error {
// 	r.log.Debug("Verifying user with token")

// 	err := r.db.WithContext(ctx).Model(&models.User{}).
// 		Where("verification_token = ? AND verification_expires > ?", token, time.Now()).
// 		Updates(map[string]interface{}{
// 			"verified":                true,
// 			"verification_token":      nil,
// 			"verification_expires": nil,
// 		}).Error

// 	if err != nil {
// 		r.log.Error("Failed to verify user", logger.NamedError("error", err))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) SaveResetToken(ctx context.Context, email, token string, expires time.Time) error {
// 	r.log.Debug("Saving password reset token", logger.String("email", email))

// 	err := r.db.WithContext(ctx).Model(&models.User{}).
// 		Where("email = ?", email).
// 		Updates(map[string]interface{}{
// 			"reset_token":      token,
// 			"reset_expires_at": expires,
// 		}).Error

// 	if err != nil {
// 		r.log.Error("Failed to save reset token", logger.NamedError("error", err),
// 			logger.String("email", email))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) ResetPassword(ctx context.Context, token, newPassword string) error {
// 	r.log.Debug("Resetting password with token")

// 	hashedPassword, err := auth.EncryptPassword(newPassword)
// 	if err != nil {
// 		r.log.Error("Failed to hash password", logger.NamedError("error", err))
// 		return err
// 	}

// 	err = r.db.WithContext(ctx).Model(&models.User{}).
// 		Where("reset_token = ? AND reset_expires_at > ?", token, time.Now()).
// 		Updates(map[string]interface{}{
// 			"password":         hashedPassword,
// 			"reset_token":      nil,
// 			"reset_expires_at": nil,
// 		}).Error

// 	if err != nil {
// 		r.log.Error("Failed to reset password", logger.NamedError("error", err))
// 		return err
// 	}

// 	return nil
// }

// func (r *userRespository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
// 	r.log.Debug("Updating avatar", logger.String("userID", userID))

// 	err := r.db.WithContext(ctx).Model(&models.User{}).
// 		Where("id = ?", userID).
// 		Update("avatar_url", avatarURL).Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return models.ErrUserNotFound
// 		}
// 		r.log.Error("Failed to update avatar", logger.NamedError("error", err),
// 			logger.String("userID", userID))
// 		return err
// 	}

// 	return nil
// }

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"gorm.io/gorm"
)

type userRespository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewUserRepository(db *gorm.DB, log logger.Logger) UserRepository {
	return &userRespository{
		db:  db,
		log: log,
	}
}

func (r *userRespository) Create(ctx context.Context, user *models.User) error {
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

func (r *userRespository) FindUser(ctx context.Context, user *models.User) (*models.User, error) {
	var existingUser models.User
	query := r.db.WithContext(ctx).Where("email = ? OR username = ?", user.Email, user.Username)

	err := query.First(&existingUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, models.ErrUserNotFound
	}

	if err != nil {
		r.log.Error("Failed to find user", logger.NamedError("error", err),
			logger.String("email", user.Email), logger.String("username", user.Username))
		return nil, err
	}

	return &existingUser, nil
}

func (r *userRespository) Update(ctx context.Context, user *models.User) error {
	r.log.Debug("Updating user", logger.String("userID", user.ID))

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		r.log.Error("Failed to update user", logger.NamedError("error", err),
			logger.String("userID", user.ID))
		return err
	}

	return nil
}

func (r *userRespository) Delete(ctx context.Context, id string) error {
	r.log.Debug("Deleting user", logger.String("userID", id))

	err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
	if err != nil {
		r.log.Error("Failed to delete user", logger.NamedError("error", err),
			logger.String("userID", id))
		return err
	}

	return nil
}

func (r *userRespository) SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error {
	r.log.Debug("Saving verification token", logger.String("email", email))

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"verification_token":       token,
			"verification_expires_at": expires,
		}).Error

	if err != nil {
		r.log.Error("Failed to save verification token", logger.NamedError("error", err),
			logger.String("email", email))
		return err
	}

	return nil
}

func (r *userRespository) VerifyUser(ctx context.Context, token string) error {
	r.log.Debug("Verifying user with token")

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("verification_token = ? AND verification_expires_at > ?", token, time.Now()).
		Updates(map[string]interface{}{
			"is_verified":              true,
			"verification_token":       nil,
			"verification_expires_at":  nil,
		}).Error

	if err != nil {
		r.log.Error("Failed to verify user", logger.NamedError("error", err))
		return err
	}

	return nil
}

func (r *userRespository) SaveResetToken(ctx context.Context, email, token string, expires time.Time) error {
	r.log.Debug("Saving password reset token", logger.String("email", email))

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"reset_password_token":    token,
			"reset_password_expires_at": expires,
		}).Error

	if err != nil {
		r.log.Error("Failed to save reset token", logger.NamedError("error", err),
			logger.String("email", email))
		return err
	}

	return nil
}

func (r *userRespository) ResetPassword(ctx context.Context, token, newPassword string) error {
	r.log.Debug("Resetting password with token")

	hashedPassword, err := auth.EncryptPassword(newPassword)
	if err != nil {
		r.log.Error("Failed to hash password", logger.NamedError("error", err))
		return err
	}

	err = r.db.WithContext(ctx).Model(&models.User{}).
		Where("reset_password_token = ? AND reset_password_expires_at > ?", token, time.Now()).
		Updates(map[string]interface{}{
			"password":                   hashedPassword,
			"reset_password_token":      nil,
			"reset_password_expires_at": nil,
		}).Error

	if err != nil {
		r.log.Error("Failed to reset password", logger.NamedError("error", err))
		return err
	}

	return nil
}

func (r *userRespository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	r.log.Debug("Updating avatar", logger.String("userID", userID))

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar", avatarURL).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		r.log.Error("Failed to update avatar", logger.NamedError("error", err),
			logger.String("userID", userID))
		return err
	}

	return nil
}