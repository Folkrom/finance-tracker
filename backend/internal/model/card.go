package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Card struct {
	Base
	PaymentMethodID     uuid.UUID        `gorm:"type:uuid;not null" json:"payment_method_id"`
	PaymentMethod       *PaymentMethod   `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Bank                string           `gorm:"type:varchar(255);not null" json:"bank"`
	CardLimit           decimal.Decimal  `gorm:"type:decimal(12,2);not null" json:"card_limit"`
	RecommendedMaxPct   decimal.Decimal  `gorm:"type:decimal(5,2);not null;default:30.00" json:"recommended_max_pct"`
	ManualUsageOverride *decimal.Decimal `gorm:"type:decimal(12,2)" json:"manual_usage_override,omitempty"`
	Level               *string          `gorm:"type:varchar(50)" json:"level,omitempty"`
}

func (Card) TableName() string { return "cards" }
