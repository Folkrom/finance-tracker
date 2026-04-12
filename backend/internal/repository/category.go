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
	if err != nil {
		return nil, err
	}
	setGlobalFlag(cats)
	return cats, nil
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
	setGlobalFlagSingle(&cat)
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
	setGlobalFlagSingle(&cat)
	return &cat, nil
}

// GetGlobalByID returns a category only if it's global (user_id IS NULL).
func (r *CategoryRepository) GetGlobalByID(id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("id = ? AND user_id IS NULL", id).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	setGlobalFlagSingle(&cat)
	return &cat, nil
}

// CreateGlobal creates a global category (user_id = NULL).
func (r *CategoryRepository) CreateGlobal(cat *model.Category) error {
	cat.UserID = nil
	return r.db.Create(cat).Error
}

// ReassignAndDelete reassigns all FK references to replacementID then deletes categoryID, in a transaction.
func (r *CategoryRepository) ReassignAndDelete(categoryID, replacementID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tables := []string{"incomes", "expenses", "debts", "budgets", "wishlist_items"}
		for _, table := range tables {
			if err := tx.Exec(
				"UPDATE "+table+" SET category_id = ? WHERE category_id = ?",
				replacementID, categoryID,
			).Error; err != nil {
				return err
			}
		}
		return tx.Where("id = ?", categoryID).Delete(&model.Category{}).Error
	})
}

func setGlobalFlag(cats []model.Category) {
	for i := range cats {
		cats[i].IsGlobal = cats[i].UserID == nil
	}
}

func setGlobalFlagSingle(cat *model.Category) {
	cat.IsGlobal = cat.UserID == nil
}
