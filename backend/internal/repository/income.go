package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type IncomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) Create(income *model.Income) error {
	return r.db.Create(income).Error
}

func (r *IncomeRepository) ListByYear(userID uuid.UUID, year int) ([]model.Income, error) {
	var incomes []model.Income
	err := r.db.
		Where("user_id = ? AND year = ?", userID, year).
		Preload("Category").
		Order("date DESC").
		Find(&incomes).Error
	return incomes, err
}

func (r *IncomeRepository) GetByID(userID, id uuid.UUID) (*model.Income, error) {
	var income model.Income
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		First(&income).Error
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (r *IncomeRepository) Update(income *model.Income) error {
	return r.db.Save(income).Error
}

func (r *IncomeRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Income{}).Error
}

func (r *IncomeRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error) {
	var results []MonthSum
	err := r.db.Model(&model.Income{}).
		Select("EXTRACT(MONTH FROM date)::int AS month, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND year = ?", userID, year).
		Group("month").Order("month").
		Scan(&results).Error
	return results, err
}

func (r *IncomeRepository) SumByCategory(userID uuid.UUID, year int) ([]CategorySumRow, error) {
	var results []CategorySumRow
	err := r.db.Model(&model.Income{}).
		Select("c.id AS category_id, c.name AS category_name, COALESCE(SUM(incomes.amount), 0) AS total").
		Joins("LEFT JOIN categories c ON incomes.category_id = c.id").
		Where("incomes.user_id = ? AND incomes.year = ?", userID, year).
		Group("c.id, c.name").Order("total DESC").
		Scan(&results).Error
	return results, err
}

func (r *IncomeRepository) TotalByYear(userID uuid.UUID, year int) (decimal.Decimal, error) {
	var result decimal.Decimal
	err := r.db.Model(&model.Income{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND year = ?", userID, year).
		Scan(&result).Error
	return result, err
}
