// package services

// import (
// 	"context"
// 	"github.com/imraushankr/bervity/server/src/configs"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// )

// type subscriptionService struct {
// 	subRepo    interfaces.SubscriptionRepository
// 	creditRepo interfaces.CreditRepository
// 	logger     logger.Logger
// 	cfg        *configs.Config
// }

// func NewSubscriptionService(
// 	subRepo interfaces.SubscriptionRepository,
// 	creditRepo interfaces.CreditRepository,
// 	log logger.Logger,
// 	cfg *configs.Config,
// ) interfaces.SubscriptionService {
// 	return &subscriptionService{
// 		subRepo:    subRepo,
// 		creditRepo: creditRepo,
// 		logger:     log,
// 		cfg:        cfg,
// 	}
// }

// func (s *subscriptionService) CreateSubscription(ctx context.Context, userID string, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
// 	// Implement subscription creation logic
// 	// This is a simplified version
// 	subscription := &models.Subscription{
// 		UserID:     userID,
// 		PlanID:     req.PlanID,
// 		Status:     "active",
// 		StartedAt:  time.Now(),
// 		ExpiresAt:  time.Now().AddDate(0, 1, 0), // 1 month
// 	}

// 	if err := s.subRepo.CreateSubscription(ctx, subscription); err != nil {
// 		s.logger.Error("failed to create subscription",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}

// 	// Record payment
// 	payment := &models.Payment{
// 		UserID:         userID,
// 		SubscriptionID: subscription.ID,
// 		Amount:         s.getPlanPrice(req.PlanID),
// 		Status:         "completed",
// 	}

// 	if err := s.subRepo.CreatePayment(ctx, payment); err != nil {
// 		s.logger.Error("failed to record payment",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}

// 	return subscription, nil
// }

// func (s *subscriptionService) GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error) {
// 	return s.subRepo.GetUserSubscription(ctx, userID)
// }

// func (s *subscriptionService) UpdateSubscription(ctx context.Context, userID string, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
// 	// Implement subscription update logic
// 	sub, err := s.subRepo.GetUserSubscription(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sub.PlanID = req.PlanID
// 	if err := s.subRepo.UpdateSubscription(ctx, sub); err != nil {
// 		return nil, err
// 	}

// 	return sub, nil
// }

// func (s *subscriptionService) CancelSubscription(ctx context.Context, userID string, req *models.CancelSubscriptionRequest) error {
// 	return s.subRepo.CancelSubscription(ctx, userID)
// }

// func (s *subscriptionService) GetSubscriptionPlans(ctx context.Context) ([]*models.SubscriptionPlanResponse, error) {
// 	// Return hardcoded plans for simplicity
// 	return []*models.SubscriptionPlanResponse{
// 		{
// 			ID:          "basic",
// 			Name:        "Basic",
// 			Description: "Basic plan with limited features",
// 			Price:       9.99,
// 			Features:    []string{"Feature 1", "Feature 2"},
// 		},
// 		{
// 			ID:          "pro",
// 			Name:        "Pro",
// 			Description: "Pro plan with all features",
// 			Price:       19.99,
// 			Features:    []string{"All Features", "Priority Support"},
// 		},
// 	}, nil
// }

// func (s *subscriptionService) GetPaymentHistory(ctx context.Context, userID string) ([]*models.Payment, error) {
// 	return s.subRepo.GetUserPayments(ctx, userID)
// }

// func (s *subscriptionService) getPlanPrice(planID string) float64 {
// 	// Simplified pricing logic
// 	switch planID {
// 	case "pro":
// 		return 19.99
// 	default:
// 		return 9.99
// 	}
// }



package services

import (
	"context"
	"time"
	
	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

type subscriptionService struct {
	subRepo    interfaces.SubscriptionRepository
	creditRepo interfaces.CreditRepository
	log        logger.Logger
	cfg        *configs.Config
}

func NewSubscriptionService(
	subRepo interfaces.SubscriptionRepository,
	creditRepo interfaces.CreditRepository,
	log logger.Logger,
	cfg *configs.Config,
) interfaces.SubscriptionService {
	return &subscriptionService{
		subRepo:    subRepo,
		creditRepo: creditRepo,
		log:        log,
		cfg:        cfg,
	}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, userID string, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	// Validate plan
	if !s.isValidPlan(req.Plan) {
		return nil, models.ErrInvalidPlan
	}

	// Create subscription
	subscription := &models.Subscription{
		UserID:    userID,
		Plan:      req.Plan,
		IsActive:  true,
		StartsAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(0, 1, 0), // 1 month
	}

	if err := s.subRepo.CreateSubscription(ctx, subscription); err != nil {
		s.log.Error("failed to create subscription",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}

	// Record payment
	payment := &models.Payment{
		UserID:         userID,
		SubscriptionID: subscription.ID,
		Amount:         s.getPlanPrice(req.Plan),
		Currency:       "usd",
		Status:         "paid",
		Description:    string(req.Plan) + " subscription",
		PaidAt:         timeNowPtr(),
	}

	if err := s.subRepo.CreatePayment(ctx, payment); err != nil {
		s.log.Error("failed to record payment",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, models.ErrPaymentFailed
	}

	// Add credits based on plan
	if err := s.addPlanCredits(ctx, userID, req.Plan); err != nil {
		s.log.Error("failed to add plan credits",
			logger.ErrorField(err),
			logger.String("userID", userID))
	}

	return subscription, nil
}

func (s *subscriptionService) GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error) {
	return s.subRepo.GetUserSubscription(ctx, userID)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, userID string, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	// Validate new plan
	if !s.isValidPlan(req.Plan) {
		return nil, models.ErrInvalidPlan
	}

	// Get current subscription
	sub, err := s.subRepo.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update plan
	sub.Plan = req.Plan
	if err := s.subRepo.UpdateSubscription(ctx, sub); err != nil {
		s.log.Error("failed to update subscription",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, userID string, req *models.CancelSubscriptionRequest) error {
	return s.subRepo.CancelSubscription(ctx, userID)
}

func (s *subscriptionService) GetSubscriptionPlans(ctx context.Context) ([]*models.SubscriptionPlanResponse, error) {
	return []*models.SubscriptionPlanResponse{
		{
			ID:          "basic",
			Name:        "Basic",
			Description: "Basic plan with limited features",
			Price:       999, // $9.99 in cents
			Currency:    "usd",
			URLsAllowed: 100,
			Features:    []string{"100 URLs/month", "Basic analytics"},
		},
		{
			ID:          "pro",
			Name:        "Pro",
			Description: "Professional plan with advanced features",
			Price:       1999, // $19.99 in cents
			Currency:    "usd",
			URLsAllowed: 500,
			Features:    []string{"500 URLs/month", "Advanced analytics", "Priority support"},
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise",
			Description: "Enterprise plan with unlimited features",
			Price:       4999, // $49.99 in cents
			Currency:    "usd",
			URLsAllowed: 0, // 0 means unlimited
			Features:    []string{"Unlimited URLs", "All features", "24/7 support"},
		},
	}, nil
}

func (s *subscriptionService) GetPaymentHistory(ctx context.Context, userID string) ([]*models.Payment, error) {
	return s.subRepo.GetUserPayments(ctx, userID)
}

func (s *subscriptionService) isValidPlan(plan models.SubscriptionPlan) bool {
	switch plan {
	case models.PlanBasic, models.PlanPro, models.PlanEnterprise:
		return true
	default:
		return false
	}
}

func (s *subscriptionService) getPlanPrice(plan models.SubscriptionPlan) int {
	switch plan {
	case models.PlanBasic:
		return 999
	case models.PlanPro:
		return 1999
	case models.PlanEnterprise:
		return 4999
	default:
		return 0
	}
}

func (s *subscriptionService) addPlanCredits(ctx context.Context, userID string, plan models.SubscriptionPlan) error {
	var credits int
	switch plan {
	case models.PlanBasic:
		credits = 100
	case models.PlanPro:
		credits = 500
	case models.PlanEnterprise:
		credits = 10000 // Essentially unlimited
	default:
		return nil
	}

	credit := &models.Credit{
		UserID:      userID,
		Type:        models.CreditTypePaid,
		Amount:      credits,
		Remaining:   credits,
		Description: string(plan) + " subscription credits",
	}

	return s.creditRepo.AddCredits(ctx, credit)
}

func timeNowPtr() *time.Time {
	t := time.Now()
	return &t
}