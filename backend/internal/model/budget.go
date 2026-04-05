package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Budget struct {
	Base
	CategoryID   uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Category     *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	MonthlyLimit decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"monthly_limit"`
	Month        int             `gorm:"not null" json:"month"`
	Year         int             `gorm:"not null" json:"year"`
	IsRecurring  bool            `gorm:"not null;default:false" json:"is_recurring"`
}

func (Budget) TableName() string { return "budgets" }
