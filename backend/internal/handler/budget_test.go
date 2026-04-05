package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgetHandler_CreateAndGetSummary(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "budgets")
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "debts")
	testutil.CleanTable(t, db, "categories")

	// Create a category first
	catRepo := repository.NewCategoryRepository(db)
	userID := uuid.New()
	cat := &model.Category{
		Base: model.Base{UserID: userID},
		Name: "Food",
		Type: model.CategoryTypeExpense,
	}
	require.NoError(t, catRepo.Create(cat))

	budgetRepo := repository.NewBudgetRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	debtRepo := repository.NewDebtRepository(db)
	budgetSvc := service.NewBudgetService(budgetRepo, expenseRepo, debtRepo)
	h := handler.NewBudgetHandler(budgetSvc)

	app := fiber.New()
	app.Use(fakeAuth(userID))
	app.Post("/api/v1/budgets", h.Create)
	app.Get("/api/v1/budgets", h.GetSummary)
	app.Get("/api/v1/budgets/recurring", h.ListRecurring)

	body := map[string]any{
		"category_id":   cat.ID.String(),
		"monthly_limit": "5000.00",
		"month":         1,
		"year":          2026,
		"is_recurring":  false,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/budgets", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var created model.Budget
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, 1, created.Month)
	assert.Equal(t, 2026, created.Year)
	assert.False(t, created.IsRecurring)

	// GetSummary
	req2 := httptest.NewRequest("GET", "/api/v1/budgets?month=1&year=2026", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var summary handler.ListResponse[service.BudgetLine]
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&summary))
	assert.Equal(t, 1, summary.Total)
	assert.Equal(t, "5000", summary.Data[0].Budget.MonthlyLimit.String())
	assert.True(t, summary.Data[0].Spent.IsZero(), fmt.Sprintf("expected spent=0, got %s", summary.Data[0].Spent))
	assert.Equal(t, "5000", summary.Data[0].Remaining.String())
}
