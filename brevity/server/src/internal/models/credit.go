// package models

// import (
// 	"time"

// 	"github.com/teris-io/shortid"
// 	"gorm.io/gorm"
// )

// var (
// 	creditSid, _ = shortid.New(1, shortid.DefaultABC, 4564)
// )

// // CreditType defines the type of credit
// type CreditType string

// const (
// 	CreditTypeFree     CreditType = "free"     // Free credits given to all users
// 	CreditTypeTrial    CreditType = "trial"    // Trial credits
// 	CreditTypePaid     CreditType = "paid"     // Paid credits from subscription
// 	CreditTypePromo    CreditType = "promo"    // Promotional credits
// 	CreditTypeReferral CreditType = "referral" // Credits from referrals
// )

// // Credit represents a user's credit balance
// type Credit struct {
// 	ID          string     `json:"id" gorm:"primaryKey;type:varchar(20)"`
// 	UserID      string     `json:"user_id" gorm:"type:varchar(20);index;not null"`
// 	User        User       `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// 	Type        CreditType `json:"type" gorm:"type:varchar(20);not null"`
// 	Amount      int        `json:"amount" gorm:"not null;default:0"`
// 	Remaining   int        `json:"remaining" gorm:"not null;default:0"`
// 	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // nil means never expires
// 	Description string     `json:"description,omitempty"`
// 	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
// }

// // CreditUsage tracks how credits are consumed
// type CreditUsage struct {
// 	ID        string    `json:"id" gorm:"primaryKey;type:varchar(20)"`
// 	UserID    string    `json:"user_id" gorm:"type:varchar(20);index;not null"`
// 	User      User      `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// 	CreditID  string    `json:"credit_id" gorm:"type:varchar(20);index;not null"`
// 	Credit    Credit    `json:"-" gorm:"foreignKey:CreditID;constraint:OnDelete:CASCADE"`
// 	URLID     string    `json:"url_id,omitempty" gorm:"type:varchar(20);index"` // Optional, tracks URL creation usage
// 	URL       URL       `json:"-" gorm:"foreignKey:URLID;constraint:OnDelete:SET NULL"`
// 	Amount    int       `json:"amount" gorm:"not null;default:1"`  // Usually 1 per URL
// 	Operation string    `json:"operation" gorm:"type:varchar(50)"` // e.g., "url_creation", "api_call"
// 	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
// }

// // CreditBalanceResponse shows user's credit status
// type CreditBalanceResponse struct {
// 	TotalCredits int  `json:"total_credits"`
// 	UsedCredits  int  `json:"used_credits"`
// 	Remaining    int  `json:"remaining_credits"`
// 	FreeLimit    int  `json:"free_limit"` // Max free URLs allowed
// 	PaidLimit    int  `json:"paid_limit"` // Max paid URLs allowed
// 	CanCreate    bool `json:"can_create"` // Whether user can create more URLs
// }

// // ApplyPromoCodeRequest for applying promo codes
// type ApplyPromoCodeRequest struct {
// 	Code string `json:"code" validate:"required"`
// }

// // BeforeCreate generates a short ID for Credit
// func (c *Credit) BeforeCreate(tx *gorm.DB) error {
// 	id, err := creditSid.Generate()
// 	if err != nil {
// 		return err
// 	}
// 	c.ID = id
// 	return nil
// }

// // BeforeCreate generates a short ID for CreditUsage
// func (cu *CreditUsage) BeforeCreate(tx *gorm.DB) error {
// 	id, err := creditSid.Generate()
// 	if err != nil {
// 		return err
// 	}
// 	cu.ID = id
// 	return nil
// }

package models

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

var (
	creditSid, _ = shortid.New(1, shortid.DefaultABC, 4564)
)

// CreditType defines the type of credit
type CreditType string

const (
	CreditTypeFree     CreditType = "free"     // Free credits given to all users
	CreditTypeTrial    CreditType = "trial"    // Trial credits
	CreditTypePaid     CreditType = "paid"     // Paid credits from subscription
	CreditTypePromo    CreditType = "promo"    // Promotional credits
	CreditTypeReferral CreditType = "referral" // Credits from referrals
)

// Credit represents a user's credit balance
type Credit struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(20)"`
	UserID      string     `json:"user_id" gorm:"type:varchar(20);index;not null"`
	User        User       `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Type        CreditType `json:"type" gorm:"type:varchar(20);not null"`
	Amount      int        `json:"amount" gorm:"not null;default:0"`
	Remaining   int        `json:"remaining" gorm:"not null;default:0"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // nil means never expires
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CreditUsage tracks how credits are consumed
type CreditUsage struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(20)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(20);index;not null"`
	User      User      `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreditID  string    `json:"credit_id" gorm:"type:varchar(20);index;not null"`
	Credit    Credit    `json:"-" gorm:"foreignKey:CreditID;constraint:OnDelete:CASCADE"`
	URLID     string    `json:"url_id,omitempty" gorm:"type:varchar(20);index"` // Optional, tracks URL creation usage
	URL       URL       `json:"-" gorm:"foreignKey:URLID;constraint:OnDelete:SET NULL"`
	Amount    int       `json:"amount" gorm:"not null;default:1"`  // Usually 1 per URL
	Operation string    `json:"operation" gorm:"type:varchar(50)"` // e.g., "url_creation", "api_call"
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// CreditBalanceResponse shows user's credit status
type CreditBalanceResponse struct {
	TotalCredits int  `json:"total_credits"`
	UsedCredits  int  `json:"used_credits"`
	Remaining    int  `json:"remaining_credits"`
	FreeLimit    int  `json:"free_limit"` // Max free URLs allowed
	UsedFree     int  `json:"used_free"`  // Number of free URLs used
	CanCreate    bool `json:"can_create"` // Whether user can create more URLs
}

// ApplyPromoCodeRequest for applying promo codes
type ApplyPromoCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

// BeforeCreate generates a short ID for Credit
func (c *Credit) BeforeCreate(tx *gorm.DB) error {
	id, err := creditSid.Generate()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

// BeforeCreate generates a short ID for CreditUsage
func (cu *CreditUsage) BeforeCreate(tx *gorm.DB) error {
	id, err := creditSid.Generate()
	if err != nil {
		return err
	}
	cu.ID = id
	return nil
}