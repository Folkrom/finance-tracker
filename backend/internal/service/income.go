package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type IncomeService struct {
	repo *repository.IncomeRepository
}

func NewIncomeService(repo *repository.IncomeRepository) *IncomeService {
	return &IncomeService{repo: repo}
}

func (s *IncomeService) Create(userID uuid.UUID, source string, amount decimal.Decimal, currency string, categoryID *uuid.UUID, date time.Time) (*model.Income, error) {
	if currency == "" {
		currency = "MXN"
	}
	income := &model.Income{
		Base:       model.Base{UserID: userID},
		Source:     source,
		Amount:     amount,
		Currency:   currency,
		CategoryID: categoryID,
		Date:       date,
		Year:       date.Year(),
	}
	if err := s.repo.Create(income); err != nil {
		return nil, err
	}
	// Re-fetch to preload Category
	return s.repo.GetByID(userID, income.ID)
}

func (s *IncomeService) ListByYear(userID uuid.UUID, year int) ([]model.Income, error) {
	return s.repo.ListByYear(userID, year)
}

func (s *IncomeService) GetByID(userID, id uuid.UUID) (*model.Income, error) {
	return s.repo.GetByID(userID, id)
}

func (s *IncomeService) Update(userID, id uuid.UUID, source string, amount decimal.Decimal, currency string, categoryID *uuid.UUID, date time.Time) (*model.Income, error) {
	income, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		currency = "MXN"
	}
	income.Source = source
	income.Amount = amount
	income.Currency = currency
	income.CategoryID = categoryID
	income.Date = date
	income.Year = date.Year()

	if err := s.repo.Update(income); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, income.ID)
}

func (s *IncomeService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
