package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) Create(pm *model.PaymentMethod) error {
	return r.db.Create(pm).Error
}

func (r *PaymentMethodRepository) ListByUser(userID uuid.UUID) ([]model.PaymentMethod, error) {
	var pms []model.PaymentMethod
	err := r.db.
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&pms).Error
	return pms, err
}

func (r *PaymentMethodRepository) ListByType(userID uuid.UUID, pmType model.PaymentMethodType) ([]model.PaymentMethod, error) {
	var pms []model.PaymentMethod
	err := r.db.
		Where("user_id = ? AND type = ?", userID, pmType).
		Order("name ASC").
		Find(&pms).Error
	return pms, err
}

func (r *PaymentMethodRepository) GetByID(userID, id uuid.UUID) (*model.PaymentMethod, error) {
	var pm model.PaymentMethod
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&pm).Error
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func (r *PaymentMethodRepository) Update(pm *model.PaymentMethod) error {
	return r.db.Save(pm).Error
}

func (r *PaymentMethodRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.PaymentMethod{}).Error
}
