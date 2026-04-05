package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetLine struct {
	Budget    model.Budget    `json:"budget"`
	Spent     decimal.Decimal `json:"spent"`
	Remaining decimal.Decimal `json:"remaining"`
}

type BudgetService struct {
	repo        *repository.BudgetRepository
	expenseRepo *repository.ExpenseRepository
	debtRepo    *repository.DebtRepository
}

func NewBudgetService(repo *repository.BudgetRepository, expenseRepo *repository.ExpenseRepository, debtRepo *repository.DebtRepository) *BudgetService {
	return &BudgetService{repo: repo, expenseRepo: expenseRepo, debtRepo: debtRepo}
}

// GetSummary returns budget lines for the given month/year, deduplicating overrides vs. recurring templates.
func (s *BudgetService) GetSummary(userID uuid.UUID, month, year int) ([]BudgetLine, error) {
	budgets, err := s.repo.ListByMonthYear(userID, month, year)
	if err != nil {
		return nil, err
	}

	// Deduplicate: if a category has a specific override (is_recurring=false) AND a recurring
	// template, keep only the override. Build a map: categoryID -> override budget, then
	// separately track which categories have overrides.
	overrides := make(map[uuid.UUID]model.Budget)
	recurring := make(map[uuid.UUID]model.Budget)

	for _, b := range budgets {
		if !b.IsRecurring {
			overrides[b.CategoryID] = b
		} else {
			recurring[b.CategoryID] = b
		}
	}

	// Merged: overrides take priority over recurring templates
	merged := make(map[uuid.UUID]model.Budget)
	for catID, b := range recurring {
		merged[catID] = b
	}
	for catID, b := range overrides {
		merged[catID] = b
	}

	lines := make([]BudgetLine, 0, len(merged))
	for _, budget := range merged {
		expenseSum, err := s.expenseRepo.SumByCategoryMonth(userID, budget.CategoryID, month, year)
		if err != nil {
			return nil, err
		}
		debtSum, err := s.debtRepo.SumByCategoryMonth(userID, budget.CategoryID, month, year)
		if err != nil {
			return nil, err
		}
		spent := expenseSum.Add(debtSum)
		remaining := budget.MonthlyLimit.Sub(spent)
		lines = append(lines, BudgetLine{
			Budget:    budget,
			Spent:     spent,
			Remaining: remaining,
		})
	}

	return lines, nil
}

func (s *BudgetService) Create(userID uuid.UUID, categoryID uuid.UUID, monthlyLimit decimal.Decimal, month, year int, isRecurring bool) (*model.Budget, error) {
	budget := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   categoryID,
		MonthlyLimit: monthlyLimit,
		Month:        month,
		Year:         year,
		IsRecurring:  isRecurring,
	}
	if err := s.repo.Create(budget); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, budget.ID)
}

func (s *BudgetService) ListRecurring(userID uuid.UUID) ([]model.Budget, error) {
	return s.repo.ListRecurring(userID)
}

func (s *BudgetService) Update(userID, id uuid.UUID, categoryID uuid.UUID, monthlyLimit decimal.Decimal, month, year int, isRecurring bool) (*model.Budget, error) {
	budget, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	budget.CategoryID = categoryID
	budget.MonthlyLimit = monthlyLimit
	budget.Month = month
	budget.Year = year
	budget.IsRecurring = isRecurring

	if err := s.repo.Update(budget); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, budget.ID)
}

func (s *BudgetService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
