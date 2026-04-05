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

func TestPaymentMethodHandler_CreateAndList(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	svc := service.NewPaymentMethodService(repo)
	h := handler.NewPaymentMethodHandler(svc)

	userID := uuid.New()

	app := fiber.New()
	app.Use(fakeAuth(userID))
	app.Post("/api/v1/payment-methods", h.Create)
	app.Get("/api/v1/payment-methods", h.List)

	body := map[string]any{
		"name": "BBVA Debit",
		"type": "debit_card",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/payment-methods", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var created model.PaymentMethod
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, "BBVA Debit", created.Name)
	assert.Equal(t, model.PaymentMethodDebitCard, created.Type)

	req2 := httptest.NewRequest("GET", "/api/v1/payment-methods", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var list handler.ListResponse[model.PaymentMethod]
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&list))
	assert.Equal(t, 1, list.Total)
	assert.Equal(t, "BBVA Debit", list.Data[0].Name)
}
