package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

type Storage interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string, publicID string) (string, error)
	UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, publicID string) error
}

type CloudinaryStorage struct {
	cld    *cloudinary.Cloudinary
	cfg    *configs.CloudinaryConfig
	local  *LocalStorage
	logger logger.Logger
}

type LocalStorage struct {
	uploadDir string
	baseURL   string
	logger    logger.Logger
}

func NewStorage(cfg *configs.Config) (Storage, error) {
	log := logger.Get()

	local := &LocalStorage{
		uploadDir: cfg.Storage.UploadDir,
		baseURL:   cfg.App.BaseURL,
		logger:    log,
	}

	if cfg.Cloudinary.CloudName != "" {
		cld, err := cloudinary.NewFromParams(
			cfg.Cloudinary.CloudName,
			cfg.Cloudinary.APIKey,
			cfg.Cloudinary.APISecret,
		)
		if err != nil {
			log.Error("Failed to initialize Cloudinary, falling back to local storage",
				logger.ErrorField(err))
			return local, nil
		}

		log.Info("Cloudinary storage initialized")
		return &CloudinaryStorage{
			cld:    cld,
			cfg:    &cfg.Cloudinary,
			local:  local,
			logger: log,
		}, nil
	}

	log.Info("Using local storage")
	return local, nil
}

func (cs *CloudinaryStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string, publicID string) (string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	uploadResult, err := cs.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   folder,
		PublicID: publicID,
	})
	if err != nil {
		cs.logger.Warn("Failed to upload to Cloudinary, falling back to local storage",
			logger.ErrorField(err))
		return cs.local.UploadFile(ctx, file, header, folder, publicID)
	}

	cs.logger.Info("File uploaded to Cloudinary",
		logger.String("url", uploadResult.SecureURL),
		logger.String("public_id", uploadResult.PublicID))

	return uploadResult.SecureURL, nil
}

func (cs *CloudinaryStorage) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	return cs.UploadFile(ctx, file, header, "avatars", fmt.Sprintf("avatar_%s", userID))
}

func (cs *CloudinaryStorage) DeleteFile(ctx context.Context, publicID string) error {
	_, err := cs.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudinary: %w", err)
	}
	return nil
}

func (ls *LocalStorage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string, publicID string) (string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	fullPath := filepath.Join(ls.uploadDir, folder)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	gitKeepPath := filepath.Join(fullPath, ".gitkeep")
	if _, err := os.Stat(gitKeepPath); os.IsNotExist(err) {
		if _, err := os.Create(gitKeepPath); err != nil {
			ls.logger.Warn("Failed to create .gitkeep file", logger.ErrorField(err))
		}
	}

	if publicID == "" {
		publicID = fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
	} else {
		publicID += filepath.Ext(header.Filename)
	}

	filePath := filepath.Join(fullPath, publicID)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	url := fmt.Sprintf("%s/uploads/%s/%s", strings.TrimRight(ls.baseURL, "/"), folder, publicID)

	ls.logger.Info("File saved locally",
		logger.String("path", filePath),
		logger.String("url", url))

	return url, nil
}

func (ls *LocalStorage) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	return ls.UploadFile(ctx, file, header, "avatars", fmt.Sprintf("avatar_%s", userID))
}

func (ls *LocalStorage) DeleteFile(ctx context.Context, publicID string) error {
	filePath := filepath.Join(ls.uploadDir, publicID)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete local file: %w", err)
	}
	return nil
}
