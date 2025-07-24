// package repository

// import (
// 	"context"
// 	"time"

// 	"github.com/imraushankr/bervity/server/src/internal/models"
// )

// type UserRepository interface {
// 	Create(ctx context.Context, user *models.User) error
// 	FindUser(ctx context.Context, user *models.User) (*models.User, error)
// 	Update(ctx context.Context, user *models.User) error
// 	Delete(ctx context.Context, id string) error
// 	SaveVerificationToken(ctx context.Context, email, token string, expires time.Time) error
// 	VerifyUser(ctx context.Context, token string) error
// 	SaveResetToken(ctx context.Context, email, token string, expires time.Time) error
// 	ResetPassword(ctx context.Context, token, newPassword string) error
// 	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
// }

package repository