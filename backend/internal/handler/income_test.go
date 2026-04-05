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

func TestIncomeHandler_CreateAndList(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	incomeRepo := repository.NewIncomeRepository(db)
	incomeSvc := service.NewIncomeService(incomeRepo)
	h := handler.NewIncomeHandler(incomeSvc)

	userID := uuid.New()

	app := fiber.New()
	app.Use(fakeAuth(userID))
	app.Post("/api/v1/years/:year/incomes", h.Create)
	app.Get("/api/v1/years/:year/incomes", h.ListByYear)

	body := map[string]any{
		"source":   "Salary",
		"amount":   "15000.00",
		"currency": "MXN",
		"date":     "2026-01-15",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/years/2026/incomes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var created model.Income
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, "Salary", created.Source)
	assert.Equal(t, 2026, created.Year)
	assert.Equal(t, "MXN", created.Currency)

	req2 := httptest.NewRequest("GET", "/api/v1/years/2026/incomes", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var list handler.ListResponse[model.Income]
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&list))
	assert.Equal(t, 1, list.Total)
	assert.Equal(t, "Salary", list.Data[0].Source)
}
