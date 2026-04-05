package service

import (
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DebtService struct {
	repo *repository.DebtRepository
}

func NewDebtService(repo *repository.DebtRepository) *DebtService {
	return &DebtService{repo: repo}
}

func (s *DebtService) Create(userID uuid.UUID, name string, amount decimal.Decimal, currency string, paymentMethodID *uuid.UUID, categoryID *uuid.UUID, date time.Time) (*model.Debt, error) {
	if currency == "" {
		currency = "MXN"
	}
	debt := &model.Debt{
		Base:            model.Base{UserID: userID},
		Name:            name,
		Amount:          amount,
		Currency:        currency,
		PaymentMethodID: paymentMethodID,
		CategoryID:      categoryID,
		Date:            date,
		Year:            date.Year(),
	}
	if err := s.repo.Create(debt); err != nil {
		return nil, err
	}
	// Re-fetch to preload Category and PaymentMethod
	return s.repo.GetByID(userID, debt.ID)
}

func (s *DebtService) ListByYear(userID uuid.UUID, year int) ([]model.Debt, error) {
	return s.repo.ListByYear(userID, year)
}

func (s *DebtService) GetByID(userID, id uuid.UUID) (*model.Debt, error) {
	return s.repo.GetByID(userID, id)
}

func (s *DebtService) Update(userID, id uuid.UUID, name string, amount decimal.Decimal, currency string, paymentMethodID *uuid.UUID, categoryID *uuid.UUID, date time.Time) (*model.Debt, error) {
	debt, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		currency = "MXN"
	}
	debt.Name = name
	debt.Amount = amount
	debt.Currency = currency
	debt.PaymentMethodID = paymentMethodID
	debt.CategoryID = categoryID
	debt.Date = date
	debt.Year = date.Year()

	if err := s.repo.Update(debt); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, debt.ID)
}

func (s *DebtService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
