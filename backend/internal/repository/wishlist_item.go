package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistItemRepository struct {
	db *gorm.DB
}

func NewWishlistItemRepository(db *gorm.DB) *WishlistItemRepository {
	return &WishlistItemRepository{db: db}
}

func (r *WishlistItemRepository) Create(item *model.WishlistItem) error {
	return r.db.Create(item).Error
}

func (r *WishlistItemRepository) ListByUser(userID uuid.UUID) ([]model.WishlistItem, error) {
	var items []model.WishlistItem
	err := r.db.
		Where("user_id = ?", userID).
		Preload("Category").
		Order("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END, name").
		Find(&items).Error
	return items, err
}

func (r *WishlistItemRepository) ListByStatus(userID uuid.UUID, statuses []model.WishlistStatus) ([]model.WishlistItem, error) {
	var items []model.WishlistItem
	err := r.db.
		Where("user_id = ? AND status IN ?", userID, statuses).
		Preload("Category").
		Order("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END, name").
		Find(&items).Error
	return items, err
}

func (r *WishlistItemRepository) ListByCategory(userID uuid.UUID, categoryID uuid.UUID) ([]model.WishlistItem, error) {
	var items []model.WishlistItem
	err := r.db.
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Preload("Category").
		Order("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END, name").
		Find(&items).Error
	return items, err
}

func (r *WishlistItemRepository) GetByID(userID, id uuid.UUID) (*model.WishlistItem, error) {
	var item model.WishlistItem
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Category").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *WishlistItemRepository) Update(item *model.WishlistItem) error {
	return r.db.Save(item).Error
}

func (r *WishlistItemRepository) UpdateStatus(userID, id uuid.UUID, status model.WishlistStatus) error {
	return r.db.Model(&model.WishlistItem{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status).Error
}

func (r *WishlistItemRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.WishlistItem{}).Error
}
