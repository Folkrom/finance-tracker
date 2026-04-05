package service

import (
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseService struct {
	repo *repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) Create(userID uuid.UUID, name string, amount decimal.Decimal, currency string, paymentMethodID *uuid.UUID, categoryID *uuid.UUID, date time.Time, expenseType model.ExpenseType) (*model.Expense, error) {
	if currency == "" {
		currency = "MXN"
	}
	if expenseType == "" {
		expenseType = model.ExpenseTypeExpense
	}
	expense := &model.Expense{
		Base:            model.Base{UserID: userID},
		Name:            name,
		Amount:          amount,
		Currency:        currency,
		PaymentMethodID: paymentMethodID,
		CategoryID:      categoryID,
		Date:            date,
		Year:            date.Year(),
		Type:            expenseType,
	}
	if err := s.repo.Create(expense); err != nil {
		return nil, err
	}
	// Re-fetch to preload Category and PaymentMethod
	return s.repo.GetByID(userID, expense.ID)
}

func (s *ExpenseService) ListByYear(userID uuid.UUID, year int) ([]model.Expense, error) {
	return s.repo.ListByYear(userID, year)
}

func (s *ExpenseService) GetByID(userID, id uuid.UUID) (*model.Expense, error) {
	return s.repo.GetByID(userID, id)
}

func (s *ExpenseService) Update(userID, id uuid.UUID, name string, amount decimal.Decimal, currency string, paymentMethodID *uuid.UUID, categoryID *uuid.UUID, date time.Time, expenseType model.ExpenseType) (*model.Expense, error) {
	expense, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		currency = "MXN"
	}
	if expenseType == "" {
		expenseType = model.ExpenseTypeExpense
	}
	expense.Name = name
	expense.Amount = amount
	expense.Currency = currency
	expense.PaymentMethodID = paymentMethodID
	expense.CategoryID = categoryID
	expense.Date = date
	expense.Year = date.Year()
	expense.Type = expenseType

	if err := s.repo.Update(expense); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, expense.ID)
}

func (s *ExpenseService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
