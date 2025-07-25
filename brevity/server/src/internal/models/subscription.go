package models

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

var (
	paymentSid, _ = shortid.New(1, shortid.DefaultABC, 5675)
)

// SubscriptionPlan defines available subscription tiers
type SubscriptionPlan string

const (
	PlanFree       SubscriptionPlan = "free"
	PlanBasic      SubscriptionPlan = "basic"
	PlanPro        SubscriptionPlan = "pro"
	PlanEnterprise SubscriptionPlan = "enterprise"
)

type Subscription struct {
	ID          string           `json:"id" gorm:"primaryKey;type:varchar(20)"`
	UserID      string           `json:"user_id" gorm:"type:varchar(20);index;not null"`
	User        User             `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Plan        SubscriptionPlan `json:"plan" gorm:"type:varchar(20);not null"`
	StripeID    string           `json:"-" gorm:"type:varchar(255);index"`
	IsActive    bool             `json:"is_active" gorm:"default:true"`
	StartsAt    time.Time        `json:"starts_at" gorm:"not null"`
	ExpiresAt   time.Time        `json:"expires_at" gorm:"not null"`
	RenewsAt    *time.Time       `json:"renews_at,omitempty"`
	CancelledAt *time.Time       `json:"cancelled_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type Payment struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(20)"`
	UserID         string       `json:"user_id" gorm:"type:varchar(20);index;not null"`
	User           User         `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SubscriptionID string       `json:"subscription_id" gorm:"type:varchar(20);index"`
	Subscription   Subscription `json:"-" gorm:"foreignKey:SubscriptionID;constraint:OnDelete:SET NULL"`
	Amount         int          `json:"amount" gorm:"not null"`
	Currency       string       `json:"currency" gorm:"type:varchar(3);default:'usd'"`
	StripeID       string       `json:"-" gorm:"type:varchar(255);index"`
	Status         string       `json:"status" gorm:"type:varchar(20)"`
	Description    string       `json:"description,omitempty"`
	CreatedAt      time.Time    `json:"created_at" gorm:"autoCreateTime"`
	PaidAt         *time.Time   `json:"paid_at,omitempty"`
}

type SubscriptionPlanResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	URLsAllowed int      `json:"urls_allowed"`
	Features    []string `json:"features"`
	IsCurrent   bool     `json:"is_current,omitempty"`
}

type CreateSubscriptionRequest struct {
	Plan   SubscriptionPlan `json:"plan" validate:"required,oneof=basic pro enterprise"`
	Token  string           `json:"token" validate:"required"`
	Coupon string           `json:"coupon,omitempty"`
}

type UpdateSubscriptionRequest struct {
	Plan SubscriptionPlan `json:"plan" validate:"required,oneof=basic pro enterprise"`
}

type CancelSubscriptionRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=255"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	id, err := paymentSid.Generate()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	id, err := paymentSid.Generate()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}
