package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *model.Card) error {
	return r.db.Create(card).Error
}

func (r *CardRepository) ListByUser(userID uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := r.db.
		Where("user_id = ?", userID).
		Preload("PaymentMethod").
		Order("bank ASC").
		Find(&cards).Error
	return cards, err
}

func (r *CardRepository) GetByID(userID, id uuid.UUID) (*model.Card, error) {
	var card model.Card
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("PaymentMethod").
		First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) Update(card *model.Card) error {
	return r.db.Save(card).Error
}

func (r *CardRepository) Delete(userID, id uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Card{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
