// package repository

// import (
// 	"context"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"gorm.io/gorm"
// )

// type creditRepository struct {
// 	db     *gorm.DB
// 	log logger.Logger
// }

// func NewCreditRepository(db *gorm.DB, log logger.Logger) interfaces.CreditRepository {
// 	return &creditRepository{
// 		db:     db,
// 		log: log,
// 	}
// }

// func (r *creditRepository) GetUserCredits(ctx context.Context, userID string) ([]*models.Credit, error) {
// 	var credits []*models.Credit
// 	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&credits).Error
// 	if err != nil {
// 		r.log.Error("failed to get user credits",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}
// 	return credits, nil
// }

// func (r *creditRepository) GetUserCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error) {
// 	// Implementation depends on your business logic
// 	// This is a simplified version
// 	var balance int64
// 	err := r.db.WithContext(ctx).Model(&models.Credit{}).
// 		Select("COALESCE(SUM(amount), 0)").
// 		Where("user_id = ?", userID).
// 		Scan(&balance).Error
// 	if err != nil {
// 		r.log.Error("failed to get credit balance",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}

// 	return &models.CreditBalanceResponse{
// 		Balance:    balance,
// 		CanCreate:  balance > 0,
// 	}, nil
// }

// func (r *creditRepository) AddCredits(ctx context.Context, credit *models.Credit) error {
// 	err := r.db.WithContext(ctx).Create(credit).Error
// 	if err != nil {
// 		r.log.Error("failed to add credits",
// 			logger.ErrorField(err),
// 			logger.Any("credit", credit))
// 		return err
// 	}
// 	return nil
// }

// func (r *creditRepository) UseCredits(ctx context.Context, usage *models.CreditUsage) error {
// 	err := r.db.WithContext(ctx).Create(usage).Error
// 	if err != nil {
// 		r.log.Error("failed to record credit usage",
// 			logger.ErrorField(err),
// 			logger.Any("usage", usage))
// 		return err
// 	}
// 	return nil
// }

// func (r *creditRepository) GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error) {
// 	var usages []*models.CreditUsage
// 	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&usages).Error
// 	if err != nil {
// 		r.log.Error("failed to get credit usage",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}
// 	return usages, nil
// }


package repository

import (
	"context"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"gorm.io/gorm"
)

type creditRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCreditRepository(db *gorm.DB, log logger.Logger) interfaces.CreditRepository {
	return &creditRepository{
		db:  db,
		log: log,
	}
}

func (r *creditRepository) GetUserCredits(ctx context.Context, userID string) ([]*models.Credit, error) {
	var credits []*models.Credit
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&credits).Error
	if err != nil {
		r.log.Error("failed to get user credits",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return credits, nil
}

func (r *creditRepository) GetUserCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error) {
	var totalCredits int64
	var remainingCredits int64

	// Get total credits
	err := r.db.WithContext(ctx).Model(&models.Credit{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&totalCredits).Error
	if err != nil {
		r.log.Error("failed to get total credits",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}

	// Get remaining credits
	err = r.db.WithContext(ctx).Model(&models.Credit{}).
		Select("COALESCE(SUM(remaining), 0)").
		Where("user_id = ?", userID).
		Scan(&remainingCredits).Error
	if err != nil {
		r.log.Error("failed to get remaining credits",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}

	// Calculate used credits
	usedCredits := totalCredits - remainingCredits

	return &models.CreditBalanceResponse{
		TotalCredits: int(totalCredits),
		UsedCredits:  int(usedCredits),
		Remaining:    int(remainingCredits),
		CanCreate:    remainingCredits > 0,
	}, nil
}

func (r *creditRepository) AddCredits(ctx context.Context, credit *models.Credit) error {
	err := r.db.WithContext(ctx).Create(credit).Error
	if err != nil {
		r.log.Error("failed to add credits",
			logger.ErrorField(err),
			logger.Any("credit", credit))
		return err
	}
	return nil
}

func (r *creditRepository) UseCredits(ctx context.Context, usage *models.CreditUsage) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First, check if credit exists and has enough remaining
	var credit models.Credit
	if err := tx.Where("id = ? AND user_id = ?", usage.CreditID, usage.UserID).First(&credit).Error; err != nil {
		tx.Rollback()
		return err
	}

	if credit.Remaining < usage.Amount {
		tx.Rollback()
		return models.ErrInsufficientCredits
	}

	// Update credit remaining
	if err := tx.Model(&models.Credit{}).
		Where("id = ?", usage.CreditID).
		Update("remaining", gorm.Expr("remaining - ?", usage.Amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Record the usage
	if err := tx.Create(usage).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *creditRepository) GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error) {
	var usages []*models.CreditUsage
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&usages).Error
	if err != nil {
		r.log.Error("failed to get credit usage",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return usages, nil
}