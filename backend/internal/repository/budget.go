package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(budget *model.Budget) error {
	return r.db.Create(budget).Error
}

// ListByMonthYear returns budgets that apply to the given month/year:
// either a specific override (is_recurring=false, month=?, year=?) or
// recurring templates (is_recurring=true), ordered by category name.
func (r *BudgetRepository) ListByMonthYear(userID uuid.UUID, month, year int) ([]model.Budget, error) {
	var budgets []model.Budget
	err := r.db.
		Where("user_id = ? AND ((month = ? AND year = ?) OR is_recurring = true)", userID, month, year).
		Preload("Category").
		Order("category_id").
		Find(&budgets).Error
	return budgets, err
}

func (r *BudgetRepository) ListRecurring(userID uuid.UUID) ([]model.Budget, error) {
	var budgets []model.Budget
	err := r.db.
		Where("user_id = ? AND is_recurring = true", userID).
		Preload("Category").
		Find(&budgets).Error
	return budgets, err
}

func (r *BudgetRepository) GetByID(userID, id uuid.UUID) (*model.Budget, error) {
	var budget model.Budget
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) Update(budget *model.Budget) error {
	return r.db.Save(budget).Error
}

func (r *BudgetRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Budget{}).Error
}
