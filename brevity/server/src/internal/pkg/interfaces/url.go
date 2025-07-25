package interfaces

import (
	"context"
	"time"

	"github.com/imraushankr/bervity/server/src/internal/models"
)

type URLRepository interface {
	Create(ctx context.Context, url *models.URL) error
	CountByIP(ctx context.Context, ip string) (int, error)
	GetByID(ctx context.Context, id string) (*models.URL, error)
	GetByShortCode(ctx context.Context, code string) (*models.URL, error)
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]*models.URL, error)
	Update(ctx context.Context, url *models.URL) error
	Delete(ctx context.Context, id string) error
	IncrementClicks(ctx context.Context, id string) error
	RecordClick(ctx context.Context, click *models.URLClick) error
	GetClicksAnalytics(ctx context.Context, urlID string, from, to time.Time) ([]*models.URLClick, error)
}

type CreditRepository interface {
	GetUserCredits(ctx context.Context, userID string) ([]*models.Credit, error)
	GetUserCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error)
	AddCredits(ctx context.Context, credit *models.Credit) error
	UseCredits(ctx context.Context, usage *models.CreditUsage) error
	GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error)
	RecordFreeURLCreation(ctx context.Context, userID, urlID string) error
	GetFreeURLCount(ctx context.Context, userID string) (int, error)
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *models.Subscription) error
	CancelSubscription(ctx context.Context, userID string) error
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetUserPayments(ctx context.Context, userID string) ([]*models.Payment, error)
}

type AnalyticsRepository interface {
	GetDailyClicks(ctx context.Context, urlID string, days int) (map[string]int, error)
	GetReferrers(ctx context.Context, urlID string) (map[string]int, error)
	GetCountries(ctx context.Context, urlID string) (map[string]int, error)
	GetDevices(ctx context.Context, urlID string) (map[string]int, error)
}

type URLService interface {
	CreateURL(ctx context.Context, req *models.CreateURLRequest, userID string, ip string) (*models.URLResponse, error)
	GetURL(ctx context.Context, shortCode string) (*models.URL, error)
	GetUserURLs(ctx context.Context, userID string, limit, offset int) ([]*models.URLResponse, error)
	UpdateURL(ctx context.Context, url *models.URL) (*models.URLResponse, error)
	DeleteURL(ctx context.Context, id, userID string) error
	RedirectURL(ctx context.Context, shortCode string, clickData *models.URLClick) (string, error)
	GetURLAnalytics(ctx context.Context, urlID, userID string, from, to time.Time) ([]*models.URLClick, error)
}

type CreditService interface {
	GetCreditBalance(ctx context.Context, userID string) (*models.CreditBalanceResponse, error)
	ApplyPromoCode(ctx context.Context, userID, code string) (*models.Credit, error)
	GetCreditUsage(ctx context.Context, userID string) ([]*models.CreditUsage, error)
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, userID string, req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetUserSubscription(ctx context.Context, userID string) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, userID string, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, userID string, req *models.CancelSubscriptionRequest) error
	GetSubscriptionPlans(ctx context.Context) ([]*models.SubscriptionPlanResponse, error)
	GetPaymentHistory(ctx context.Context, userID string) ([]*models.Payment, error)
}