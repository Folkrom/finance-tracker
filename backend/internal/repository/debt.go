package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DebtRepository struct {
	db *gorm.DB
}

func NewDebtRepository(db *gorm.DB) *DebtRepository {
	return &DebtRepository{db: db}
}

func (r *DebtRepository) Create(debt *model.Debt) error {
	return r.db.Create(debt).Error
}

func (r *DebtRepository) ListByYear(userID uuid.UUID, year int) ([]model.Debt, error) {
	var debts []model.Debt
	err := r.db.
		Where("user_id = ? AND year = ?", userID, year).
		Preload("Category").
		Preload("PaymentMethod").
		Order("date DESC").
		Find(&debts).Error
	return debts, err
}

func (r *DebtRepository) GetByID(userID, id uuid.UUID) (*model.Debt, error) {
	var debt model.Debt
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		Preload("PaymentMethod").
		First(&debt).Error
	if err != nil {
		return nil, err
	}
	return &debt, nil
}

func (r *DebtRepository) Update(debt *model.Debt) error {
	return r.db.Save(debt).Error
}

func (r *DebtRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Debt{}).Error
}

func (r *DebtRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error) {
	var results []MonthSum
	err := r.db.Model(&model.Debt{}).
		Select("EXTRACT(MONTH FROM date)::int AS month, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year).
		Group("month").Order("month").
		Scan(&results).Error
	return results, err
}

func (r *DebtRepository) SumByDay(userID uuid.UUID, year int, month *int) ([]DaySumRow, error) {
	var results []DaySumRow
	query := r.db.Model(&model.Debt{}).
		Select("date::text AS date, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year)
	if month != nil {
		query = query.Where("EXTRACT(MONTH FROM date) = ?", *month)
	}
	err := query.Group("date").Order("date").Scan(&results).Error
	return results, err
}

func (r *DebtRepository) SumByPaymentMethodMonth(userID uuid.UUID, month, year int) ([]PaymentMethodSumRow, error) {
	var results []PaymentMethodSumRow
	err := r.db.Model(&model.Debt{}).
		Select("payment_method_id::text AS payment_method_id, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND payment_method_id IS NOT NULL AND EXTRACT(MONTH FROM date) = ? AND year = ?", userID, month, year).
		Group("payment_method_id").
		Scan(&results).Error
	return results, err
}

func (r *DebtRepository) SumByCategoryMonth(userID, categoryID uuid.UUID, month, year int) (decimal.Decimal, error) {
	var result decimal.Decimal
	err := r.db.
		Model(&model.Debt{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND category_id = ? AND EXTRACT(MONTH FROM date) = ? AND year = ?", userID, categoryID, month, year).
		Scan(&result).Error
	return result, err
}
