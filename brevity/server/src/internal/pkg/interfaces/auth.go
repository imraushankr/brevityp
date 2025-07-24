package interfaces

import (
	"context"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByIdentifier(ctx context.Context, identifier string) (*models.User, error)
	SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error
	VerifyUser(ctx context.Context, token string) error
	SaveResetToken(ctx context.Context, email, token string, expires time.Time) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
}

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, token string) error
	InitiatePasswordReset(ctx context.Context, email string) error
	CompletePasswordReset(ctx context.Context, token, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error
}