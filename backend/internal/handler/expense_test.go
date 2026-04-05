package handler_test

import (
	"bytes"
	"encoding/json"
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

func TestExpenseHandler_CreateAndList(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "categories")

	expenseRepo := repository.NewExpenseRepository(db)
	expenseSvc := service.NewExpenseService(expenseRepo)
	h := handler.NewExpenseHandler(expenseSvc)

	userID := uuid.New()

	app := fiber.New()
	app.Use(fakeAuth(userID))
	app.Post("/api/v1/years/:year/expenses", h.Create)
	app.Get("/api/v1/years/:year/expenses", h.ListByYear)

	body := map[string]any{
		"name":     "Groceries",
		"amount":   "2500.00",
		"currency": "MXN",
		"date":     "2026-01-15",
		"type":     "expense",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/years/2026/expenses", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var created model.Expense
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, "Groceries", created.Name)
	assert.Equal(t, 2026, created.Year)
	assert.Equal(t, "MXN", created.Currency)
	assert.Equal(t, model.ExpenseTypeExpense, created.Type)

	req2 := httptest.NewRequest("GET", "/api/v1/years/2026/expenses", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var list handler.ListResponse[model.Expense]
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&list))
	assert.Equal(t, 1, list.Total)
	assert.Equal(t, "Groceries", list.Data[0].Name)
}
