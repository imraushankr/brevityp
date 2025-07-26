package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/imraushankr/bervity/server/src/internal/models"
)

type UserRepository interface {
	FindUserByIdentifier(ctx context.Context, identifier string) (*models.User, error)
	FindUserByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
}

type UserService interface {
	FindUser(ctx context.Context, identifier string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
}
