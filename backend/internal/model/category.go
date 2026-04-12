package model

import "github.com/google/uuid"

type CategoryDomain string

const (
	CategoryDomainIncome   CategoryDomain = "income"
	CategoryDomainExpense  CategoryDomain = "expense"
	CategoryDomainWishlist CategoryDomain = "wishlist"
)

type Category struct {
	Base
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Domain    CategoryDomain `gorm:"type:varchar(50);not null" json:"domain"`
	Color     *string        `gorm:"type:varchar(7)" json:"color,omitempty"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsSystem  bool           `gorm:"not null;default:false" json:"is_system"`
}

func (Category) TableName() string {
	return "categories"
}

func (c Category) IsGlobal() bool {
	return c.UserID == nil
}
