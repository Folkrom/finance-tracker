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

func (r *ExpenseRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error) {
	var results []MonthSum
	err := r.db.Model(&model.Expense{}).
		Select("EXTRACT(MONTH FROM date)::int AS month, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year).
		Group("month").Order("month").
		Scan(&results).Error
	return results, err
}

func (r *ExpenseRepository) SumByType(userID uuid.UUID, year int) ([]TypeSumRow, error) {
	var results []TypeSumRow
	err := r.db.Model(&model.Expense{}).
		Select("type, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year).
		Group("type").
		Scan(&results).Error
	return results, err
}

func (r *ExpenseRepository) TotalByYear(userID uuid.UUID, year int) (decimal.Decimal, error) {
	var result decimal.Decimal
	err := r.db.Model(&model.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND year = ?", userID, year).
		Scan(&result).Error
	return result, err
}

func (r *ExpenseRepository) SumByDay(userID uuid.UUID, year int, month *int) ([]DaySumRow, error) {
	var results []DaySumRow
	query := r.db.Model(&model.Expense{}).
		Select("date::text AS date, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year)
	if month != nil {
		query = query.Where("EXTRACT(MONTH FROM date) = ?", *month)
	}
	err := query.Group("date").Order("date").Scan(&results).Error
	return results, err
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
