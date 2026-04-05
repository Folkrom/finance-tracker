package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Income struct {
	Base
	Source     string          `gorm:"type:varchar(255);not null" json:"source"`
	Amount     decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency   string          `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category   *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Date       time.Time       `gorm:"type:date;not null" json:"date"`
	Year       int             `gorm:"not null;index" json:"year"`
}

func (Income) TableName() string {
	return "incomes"
}
