package repository

import "github.com/shopspring/decimal"

type MonthSum struct {
	Month int             `json:"month"`
	Total decimal.Decimal `json:"total"`
}

type CategorySumRow struct {
	CategoryID   string          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Total        decimal.Decimal `json:"total"`
}

type TypeSumRow struct {
	Type  string          `json:"type"`
	Total decimal.Decimal `json:"total"`
}

type DaySumRow struct {
	Date  string          `json:"date"`
	Total decimal.Decimal `json:"total"`
}

type PaymentMethodSumRow struct {
	PaymentMethodID string          `json:"payment_method_id"`
	Total           decimal.Decimal `json:"total"`
}
