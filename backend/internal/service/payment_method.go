package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

type PaymentMethodService struct {
	repo *repository.PaymentMethodRepository
}

func NewPaymentMethodService(repo *repository.PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{repo: repo}
}

func (s *PaymentMethodService) Create(userID uuid.UUID, name string, pmType model.PaymentMethodType, details *string) (*model.PaymentMethod, error) {
	pm := &model.PaymentMethod{
		Base:    model.Base{UserID: userID},
		Name:    name,
		Type:    pmType,
		Details: details,
	}
	if err := s.repo.Create(pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *PaymentMethodService) List(userID uuid.UUID, pmType model.PaymentMethodType) ([]model.PaymentMethod, error) {
	if pmType != "" {
		return s.repo.ListByType(userID, pmType)
	}
	return s.repo.ListByUser(userID)
}

func (s *PaymentMethodService) GetByID(userID, id uuid.UUID) (*model.PaymentMethod, error) {
	return s.repo.GetByID(userID, id)
}

func (s *PaymentMethodService) Update(userID, id uuid.UUID, name string, pmType model.PaymentMethodType, details *string) (*model.PaymentMethod, error) {
	pm, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	pm.Name = name
	pm.Type = pmType
	pm.Details = details
	if err := s.repo.Update(pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *PaymentMethodService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
