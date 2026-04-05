package handler

import (
	"strconv"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type IncomeHandler struct {
	svc *service.IncomeService
}

func NewIncomeHandler(svc *service.IncomeService) *IncomeHandler {
	return &IncomeHandler{svc: svc}
}

type createIncomeRequest struct {
	Source     string  `json:"source"`
	Amount     string  `json:"amount"`
	Currency   string  `json:"currency"`
	CategoryID *string `json:"category_id"`
	Date       string  `json:"date"`
}

type updateIncomeRequest struct {
	Source     string  `json:"source"`
	Amount     string  `json:"amount"`
	Currency   string  `json:"currency"`
	CategoryID *string `json:"category_id"`
	Date       string  `json:"date"`
}

func parseIncomeRequest(amountStr string, categoryIDStr *string, dateStr string) (
	decimal.Decimal, *uuid.UUID, time.Time, error,
) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Zero, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid amount")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return decimal.Zero, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
	}

	var categoryID *uuid.UUID
	if categoryIDStr != nil && *categoryIDStr != "" {
		parsed, err := uuid.Parse(*categoryIDStr)
		if err != nil {
			return decimal.Zero, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid category_id")
		}
		categoryID = &parsed
	}

	return amount, categoryID, date, nil
}

func (h *IncomeHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createIncomeRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Source == "" {
		return respondError(c, fiber.StatusBadRequest, "source is required")
	}
	if req.Amount == "" {
		return respondError(c, fiber.StatusBadRequest, "amount is required")
	}
	if req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "date is required")
	}

	amount, categoryID, date, err := parseIncomeRequest(req.Amount, req.CategoryID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	income, err := h.svc.Create(userID, req.Source, amount, req.Currency, categoryID, date)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create income")
	}
	return respondCreated(c, income)
}

func (h *IncomeHandler) ListByYear(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	incomes, err := h.svc.ListByYear(userID, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list incomes")
	}
	return respondList(c, incomes)
}

func (h *IncomeHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	income, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "income not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get income")
	}
	return respondJSON(c, income)
}

func (h *IncomeHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateIncomeRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Source == "" {
		return respondError(c, fiber.StatusBadRequest, "source is required")
	}
	if req.Amount == "" {
		return respondError(c, fiber.StatusBadRequest, "amount is required")
	}
	if req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "date is required")
	}

	amount, categoryID, date, err := parseIncomeRequest(req.Amount, req.CategoryID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	income, err := h.svc.Update(userID, id, req.Source, amount, req.Currency, categoryID, date)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "income not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update income")
	}
	return respondJSON(c, income)
}

func (h *IncomeHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete income")
	}
	return respondNoContent(c)
}
