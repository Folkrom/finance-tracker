package model

type CategoryDomain string

const (
	CategoryDomainIncome   CategoryDomain = "income"
	CategoryDomainExpense  CategoryDomain = "expense"
	CategoryDomainWishlist CategoryDomain = "wishlist"
)

type Category struct {
	Base
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Domain    CategoryDomain `gorm:"type:varchar(50);not null" json:"domain"`
	Color     *string        `gorm:"type:varchar(7)" json:"color,omitempty"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
}

func (Category) TableName() string {
	return "categories"
}
