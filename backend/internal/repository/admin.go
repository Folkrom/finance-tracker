package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetStats() (*model.AdminStats, error) {
	var stats model.AdminStats

	if err := r.db.Raw("SELECT COUNT(*) FROM profiles").Scan(&stats.Profiles).Error; err != nil {
		return nil, err
	}
	stats.Users = stats.Profiles

	if err := r.db.Raw("SELECT COUNT(*) FROM categories WHERE user_id IS NULL").Scan(&stats.CategoriesGlobal).Error; err != nil {
		return nil, err
	}

	if err := r.db.Raw("SELECT COUNT(*) FROM categories WHERE user_id IS NOT NULL").Scan(&stats.CategoriesUser).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
