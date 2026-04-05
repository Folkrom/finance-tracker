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

func fakeAuth(userID uuid.UUID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func TestCategoryHandler_CreateAndList(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	userID := uuid.New()

	app := fiber.New()
	app.Use(fakeAuth(userID))
	app.Post("/api/v1/categories", h.Create)
	app.Get("/api/v1/categories", h.List)

	// Create a category
	body := map[string]any{
		"name":   "Salary",
		"domain": "income",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var created model.Category
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, "Salary", created.Name)
	assert.Equal(t, model.CategoryDomainIncome, created.Domain)

	// List categories
	req2 := httptest.NewRequest("GET", "/api/v1/categories?domain=income", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var list handler.ListResponse[model.Category]
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&list))
	assert.Equal(t, 1, list.Total)
	assert.Equal(t, "Salary", list.Data[0].Name)
}
