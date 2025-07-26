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

type urlRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewURLRepository(db *gorm.DB, logger logger.Logger) interfaces.URLRepository {
	return &urlRepository{
		db:     db,
		logger: logger,
	}
}

func (r *urlRepository) Create(ctx context.Context, url *models.URL) error {
	// Explicitly set UserID to nil if empty
	if url.UserID != nil && *url.UserID == "" {
		url.UserID = nil
	}

	err := r.db.WithContext(ctx).Create(url).Error
	if err != nil {
		r.logger.Error("failed to create URL",
			logger.ErrorField(err),
			logger.Any("url", url))
		return err
	}
	return nil
}

func (r *urlRepository) CountByIP(ctx context.Context, ip string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.URL{}).
		Where("created_by_ip = ? AND user_id IS NULL", ip).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *urlRepository) GetByID(ctx context.Context, id string) (*models.URL, error) {
	var url models.URL
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrURLNotFound
		}
		r.logger.Error("failed to get URL by ID",
			logger.ErrorField(err),
			logger.String("id", id))
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	result := r.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		First(&url)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrURLNotFound
		}
		return nil, result.Error
	}

	return &url, nil
}

func (r *urlRepository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]*models.URL, error) {
	var urls []*models.URL
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&urls).Error
	if err != nil {
		r.logger.Error("failed to get URLs by user",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return urls, nil
}

func (r *urlRepository) Update(ctx context.Context, url *models.URL) error {
	err := r.db.WithContext(ctx).Save(url).Error
	if err != nil {
		r.logger.Error("failed to update URL",
			logger.ErrorField(err),
			logger.Any("url", url))
		return err
	}
	return nil
}

func (r *urlRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&models.URL{}, "id = ?", id).Error
	if err != nil {
		r.logger.Error("failed to delete URL",
			logger.ErrorField(err),
			logger.String("id", id))
		return err
	}
	return nil
}

func (r *urlRepository) IncrementClicks(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Model(&models.URL{}).
		Where("id = ?", id).
		Update("clicks", gorm.Expr("clicks + ?", 1)).Error
	if err != nil {
		r.logger.Error("failed to increment URL clicks",
			logger.ErrorField(err),
			logger.String("id", id))
		return err
	}
	return nil
}

func (r *urlRepository) RecordClick(ctx context.Context, click *models.URLClick) error {
	err := r.db.WithContext(ctx).Create(click).Error
	if err != nil {
		r.logger.Error("failed to record URL click",
			logger.ErrorField(err),
			logger.Any("click", click))
		return err
	}
	return nil
}

func (r *urlRepository) GetClicksAnalytics(ctx context.Context, urlID string, from, to time.Time) ([]*models.URLClick, error) {
	var clicks []*models.URLClick
	err := r.db.WithContext(ctx).
		Where("url_id = ? AND created_at BETWEEN ? AND ?", urlID, from, to).
		Find(&clicks).Error
	if err != nil {
		r.logger.Error("failed to get URL click analytics",
			logger.ErrorField(err),
			logger.String("urlID", urlID),
			logger.Time("from", from),
			logger.Time("to", to))
		return nil, err
	}
	return clicks, nil
}
