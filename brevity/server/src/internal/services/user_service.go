package services

import (
	"context"
	"mime/multipart"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/auth"
	"github.com/imraushankr/bervity/server/src/internal/pkg/email"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/pkg/storage"
	"github.com/imraushankr/bervity/server/src/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
	auth     *auth.Auth
	email    *email.EmailService
	cfg      *configs.Config
	storage  storage.Storage
	log      logger.Logger
}

func NewUserService(
	userRepo repository.UserRepository,
	auth *auth.Auth,
	email *email.EmailService,
	cfg *configs.Config,
	storage storage.Storage,
	log logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		auth:     auth,
		email:    email,
		cfg:      cfg,
		storage:  storage,
		log:      log,
	}
}

func (s *userService) Register(ctx context.Context, user *models.User) error {
	return nil
}

func (s *userService) Login(ctx context.Context, userId, password string) (*models.User, string, error) {
	return nil, "", nil
}

func (s *userService) FindUser(ctx context.Context, identifier string) (*models.User, error) {
	return nil, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, token string) error {
	return nil
}

func (s *userService) InitiatePasswordReset(ctx context.Context, email string) error {
	return nil
}

func (s *userService) CompletePasswordReset(ctx context.Context, token, newPassword string) error {
	return nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}

func (s *userService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	return "", nil
}
