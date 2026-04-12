package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

// ListByDomain returns global categories (user_id IS NULL) and user's own categories.
func (r *CategoryRepository) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	var cats []model.Category
	err := r.db.
		Where("(user_id IS NULL OR user_id = ?) AND domain = ?", userID, domain).
		Order("sort_order ASC, name ASC").
		Find(&cats).Error
	return cats, err
}

// GetByID returns a category if it's global or owned by the user.
func (r *CategoryRepository) GetByID(userID, id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", id, userID).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Category{}).Error
}

// GetOtherCategory returns the "Other" system category for a given domain.
func (r *CategoryRepository) GetOtherCategory(domain model.CategoryDomain) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("user_id IS NULL AND domain = ? AND is_system = true AND name = 'Other'", domain).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}
