package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByUserID(userID uuid.UUID) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) Create(profile *model.Profile) error {
	return r.db.Create(profile).Error
}

func (r *ProfileRepository) Update(profile *model.Profile) error {
	return r.db.Save(profile).Error
}
