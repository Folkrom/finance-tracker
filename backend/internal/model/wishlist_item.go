package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type WishlistPriority string

const (
	WishlistPriorityLow    WishlistPriority = "low"
	WishlistPriorityMedium WishlistPriority = "medium"
	WishlistPriorityHigh   WishlistPriority = "high"
)

type WishlistStatus string

const (
	// To-do
	WishlistStatusInterested WishlistStatus = "interested"
	// In Progress
	WishlistStatusSavingFor      WishlistStatus = "saving_for"
	WishlistStatusWaitingForSale WishlistStatus = "waiting_for_sale"
	WishlistStatusOrdered        WishlistStatus = "ordered"
	// Complete
	WishlistStatusPurchased WishlistStatus = "purchased"
	WishlistStatusReceived  WishlistStatus = "received"
	WishlistStatusCancelled WishlistStatus = "cancelled"
)

type WishlistItem struct {
	Base
	Name                 string           `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL             *string          `gorm:"type:text" json:"image_url,omitempty"`
	Price                *decimal.Decimal `gorm:"type:decimal(12,2)" json:"price,omitempty"`
	Currency             string           `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	Links                pq.StringArray   `gorm:"type:text[]" json:"links"`
	CategoryID           *uuid.UUID       `gorm:"type:uuid" json:"category_id,omitempty"`
	Category             *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Priority             WishlistPriority `gorm:"type:varchar(10);not null;default:medium" json:"priority"`
	Status               WishlistStatus   `gorm:"type:varchar(30);not null;default:interested" json:"status"`
	TargetDate           *time.Time       `gorm:"type:date" json:"target_date,omitempty"`
	MonthlyContribution  *decimal.Decimal `gorm:"type:decimal(12,2)" json:"monthly_contribution,omitempty"`
	ContributionCurrency string           `gorm:"type:varchar(3);not null;default:MXN" json:"contribution_currency"`
}

func (WishlistItem) TableName() string { return "wishlist_items" }
