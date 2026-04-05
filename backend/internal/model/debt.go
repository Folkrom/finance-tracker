package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Debt struct {
	Base
	Name            string          `gorm:"type:varchar(255);not null" json:"name"`
	Amount          decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string          `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	Date            time.Time       `gorm:"type:date;not null" json:"date"`
	Year            int             `gorm:"not null;index" json:"year"`
	PaymentMethodID *uuid.UUID      `gorm:"type:uuid" json:"payment_method_id,omitempty"`
	PaymentMethod   *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	CategoryID      *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category        *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Debt) TableName() string { return "debts" }
