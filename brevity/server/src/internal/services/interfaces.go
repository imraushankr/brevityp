// package services

// import (
// 	"context"
// 	"mime/multipart"

// 	"github.com/imraushankr/bervity/server/src/internal/models"
// )

// type UserService interface {
// 	Register(ctx context.Context, user *models.User) error
// 	Login(ctx context.Context, userId, password string) (*models.User, string, error)
// 	FindUser(ctx context.Context, identifier string) (*models.User, error)
// 	UpdateUser(ctx context.Context, user *models.User) error
// 	DeleteUser(ctx context.Context, id string) error
// 	VerifyEmail(ctx context.Context, token string) error
// 	InitiatePasswordReset(ctx context.Context, email string) error
// 	CompletePasswordReset(ctx context.Context, token, newPassword string) error
// 	RefreshToken(ctx context.Context, refreshToken string) (string, error)
// 	UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
// }

package services