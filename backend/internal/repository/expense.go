package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *model.Expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseRepository) ListByYear(userID uuid.UUID, year int) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.
		Where("user_id = ? AND year = ?", userID, year).
		Preload("Category").
		Preload("PaymentMethod").
		Order("date DESC").
		Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) GetByID(userID, id uuid.UUID) (*model.Expense, error) {
	var expense model.Expense
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		Preload("PaymentMethod").
		First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) Update(expense *model.Expense) error {
	return r.db.Save(expense).Error
}

func (r *ExpenseRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Expense{}).Error
}

func (r *ExpenseRepository) SumByCategoryMonth(userID, categoryID uuid.UUID, month, year int) (decimal.Decimal, error) {
	var result decimal.Decimal
	err := r.db.
		Model(&model.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND category_id = ? AND EXTRACT(MONTH FROM date) = ? AND year = ?", userID, categoryID, month, year).
		Scan(&result).Error
	return result, err
}
