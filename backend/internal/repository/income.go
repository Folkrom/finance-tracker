package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
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
