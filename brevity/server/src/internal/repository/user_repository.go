package repository

import (
	"context"
	"errors"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
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
	// r.log.Debug("Creating user", logger.String("email", user.Email))
	if err := user.Validate(); err != nil {
		r.log.Error("User validation failed", logger.NamedError("error", err))
		return err
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.log.Error("Failed to create user", logger.NamedError("error", err))
		return err
	}

	return err
}

func (r *userRespository) FindUser(ctx context.Context, identifier string) (*models.User, error) {
	// r.log.Debug("Finding user by identifier", logger.String("identifier", identifier))
	var user models.User
	query := r.db.WithContext(ctx).Where("email = ? OR username = ? OR id = ?", identifier, identifier, identifier)
	err := query.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// r.log.Debug("User not found", logger.String("identifier", identifier))
		return nil, models.ErrUserNotFound
	}

	if err != nil {
		r.log.Error("Failed to find user", logger.NamedError("error", err), logger.String("identifier", identifier))
	}

	return &user, err
}

func (r *userRespository) FindByID(ctx context.Context, id string) (*models.User, error) {
	return r.FindUser(ctx, id)
}

func (r *userRespository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.FindUser(ctx, email)
}

func (r *userRespository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.FindUser(ctx, username)
}

func (r *userRespository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (r *userRespository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *userRespository) SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error {
	return nil
}

func (r *userRespository) VerifyUser(ctx context.Context, token string) error {
	return nil
}

func (r *userRespository) SaveResetToken(ctx context.Context, email, token string, expires time.Time) error {
	return nil
}

func (r *userRespository) ResetPassword(ctx context.Context, token, newPassword string) error {
	return nil
}

func (r *userRespository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	return nil
}
