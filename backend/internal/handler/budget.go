package handler

import (
	"strconv"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetHandler struct {
	svc *service.BudgetService
}

func NewBudgetHandler(svc *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

type createBudgetRequest struct {
	CategoryID   string `json:"category_id"`
	MonthlyLimit string `json:"monthly_limit"`
	Month        int    `json:"month"`
	Year         int    `json:"year"`
	IsRecurring  bool   `json:"is_recurring"`
}

type updateBudgetRequest struct {
	CategoryID   string `json:"category_id"`
	MonthlyLimit string `json:"monthly_limit"`
	Month        int    `json:"month"`
	Year         int    `json:"year"`
	IsRecurring  bool   `json:"is_recurring"`
}

func (h *BudgetHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.CategoryID == "" {
		return respondError(c, fiber.StatusBadRequest, "category_id is required")
	}
	if req.MonthlyLimit == "" {
		return respondError(c, fiber.StatusBadRequest, "monthly_limit is required")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid category_id")
	}

	monthlyLimit, err := decimal.NewFromString(req.MonthlyLimit)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid monthly_limit")
	}

	budget, err := h.svc.Create(userID, categoryID, monthlyLimit, req.Month, req.Year, req.IsRecurring)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create budget")
	}
	return respondCreated(c, budget)
}

func (h *BudgetHandler) GetSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" {
		return respondError(c, fiber.StatusBadRequest, "month is required")
	}
	if yearStr == "" {
		return respondError(c, fiber.StatusBadRequest, "year is required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return respondError(c, fiber.StatusBadRequest, "invalid month")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	lines, err := h.svc.GetSummary(userID, month, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to get budget summary")
	}
	return respondList(c, lines)
}

func (h *BudgetHandler) ListRecurring(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	budgets, err := h.svc.ListRecurring(userID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list recurring budgets")
	}
	return respondList(c, budgets)
}

func (h *BudgetHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.CategoryID == "" {
		return respondError(c, fiber.StatusBadRequest, "category_id is required")
	}
	if req.MonthlyLimit == "" {
		return respondError(c, fiber.StatusBadRequest, "monthly_limit is required")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid category_id")
	}

	monthlyLimit, err := decimal.NewFromString(req.MonthlyLimit)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid monthly_limit")
	}

	budget, err := h.svc.Update(userID, id, categoryID, monthlyLimit, req.Month, req.Year, req.IsRecurring)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "budget not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update budget")
	}
	return respondJSON(c, budget)
}

func (h *BudgetHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete budget")
	}
	return respondNoContent(c)
}
