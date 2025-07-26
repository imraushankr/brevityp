package repository

import (
	"context"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewSubscriptionRepository(db *gorm.DB, log logger.Logger) interfaces.SubscriptionRepository {
	return &subscriptionRepository{
		db:  db,
		log: log,
	}
}

func (r *subscriptionRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check for existing active subscription
	var existing models.Subscription
	if err := tx.Where("user_id = ? AND is_active = true", subscription.UserID).First(&existing).Error; err == nil {
		tx.Rollback()
		return models.ErrActiveSubscriptionExists
	}

	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback()
		r.log.Error("failed to create subscription",
			logger.ErrorField(err),
			logger.Any("subscription", subscription))
		return err
	}

	return tx.Commit().Error
}

func (r *subscriptionRepository) GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		First(&subscription).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrSubscriptionNotActive
		}
		r.log.Error("failed to get subscription",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, subscription *models.Subscription) error {
	err := r.db.WithContext(ctx).Save(subscription).Error
	if err != nil {
		r.log.Error("failed to update subscription",
			logger.ErrorField(err),
			logger.Any("subscription", subscription))
		return err
	}
	return nil
}

func (r *subscriptionRepository) CancelSubscription(ctx context.Context, userID string) error {
	err := r.db.WithContext(ctx).
		Model(&models.Subscription{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"cancelled_at":  time.Now(),
		}).Error

	if err != nil {
		r.log.Error("failed to cancel subscription",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return err
	}
	return nil
}

func (r *subscriptionRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	err := r.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		r.log.Error("failed to create payment",
			logger.ErrorField(err),
			logger.Any("payment", payment))
		return err
	}
	return nil
}

func (r *subscriptionRepository) GetUserPayments(ctx context.Context, userID string) ([]*models.Payment, error) {
	var payments []*models.Payment
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error

	if err != nil {
		r.log.Error("failed to get payments",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return payments, nil
}