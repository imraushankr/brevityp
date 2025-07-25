// package services

// import (
// 	"context"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// )

// type creditService struct {
// 	creditRepo interfaces.CreditRepository
// 	log    logger.Logger
// }

// func NewCreditService(creditRepo interfaces.CreditRepository, log logger.Logger) interfaces.CreditService {
// 	return &creditService{
// 		creditRepo: creditRepo,
// 		log:     log,
// 	}
// }

// func (s *creditService) GetCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error) {
// 	return s.creditRepo.GetUserCreditBalance(ctx, userID)
// }

// func (s *creditService) ApplyPromoCode(ctx context.Context, userID, code string) (*models.Credit, error) {
// 	// Implement promo code validation and application logic
// 	// This is a simplified version
// 	credit := &models.Credit{
// 		UserID: userID,
// 		Amount: 10, // Example: fixed amount for any promo code
// 		Reason: "promo_code:" + code,
// 	}

// 	if err := s.creditRepo.AddCredits(ctx, credit); err != nil {
// 		s.log.Error("failed to apply promo code",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID),
// 			logger.String("code", code))
// 		return nil, err
// 	}

// 	return credit, nil
// }

// func (s *creditService) GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error) {
// 	return s.creditRepo.GetCreditUsage(ctx, userID)
// }


package services

import (
	"context"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

type creditService struct {
	creditRepo interfaces.CreditRepository
	log        logger.Logger
}

func NewCreditService(creditRepo interfaces.CreditRepository, log logger.Logger) interfaces.CreditService {
	return &creditService{
		creditRepo: creditRepo,
		log:        log,
	}
}

func (s *creditService) GetCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error) {
	balance, err := s.creditRepo.GetUserCreditBalance(ctx, userID)
	if err != nil {
		s.log.Error("failed to get credit balance",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return balance, nil
}

func (s *creditService) ApplyPromoCode(ctx context.Context, userID, code string) (*models.Credit, error) {
	// Validate promo code (in a real app, you'd check against a database of valid codes)
	if code == "" || len(code) < 5 {
		return nil, models.ErrPromoCodeInvalid
	}

	// Check if user already used this promo code
	credits, err := s.creditRepo.GetUserCredits(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, c := range credits {
		if c.Type == models.CreditTypePromo && c.Description == "Promo code: "+code {
			return nil, models.ErrPromoCodeAlreadyUsed
		}
	}

	// Create credit record
	credit := &models.Credit{
		UserID:      userID,
		Type:        models.CreditTypePromo,
		Amount:      10, // Example fixed amount
		Remaining:   10,
		Description: "Promo code: " + code,
	}

	if err := s.creditRepo.AddCredits(ctx, credit); err != nil {
		s.log.Error("failed to apply promo code",
			logger.ErrorField(err),
			logger.String("userID", userID),
			logger.String("code", code))
		return nil, err
	}

	return credit, nil
}

func (s *creditService) GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error) {
	usages, err := s.creditRepo.GetCreditUsage(ctx, userID)
	if err != nil {
		s.log.Error("failed to get credit usage",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}
	return usages, nil
}