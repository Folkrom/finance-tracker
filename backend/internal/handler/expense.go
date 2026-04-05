package handler

import (
	"strconv"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

type createExpenseRequest struct {
	Name            string  `json:"name"`
	Amount          string  `json:"amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
	PaymentMethodID *string `json:"payment_method_id"`
	CategoryID      *string `json:"category_id"`
	Type            string  `json:"type"`
}

type updateExpenseRequest struct {
	Name            string  `json:"name"`
	Amount          string  `json:"amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
	PaymentMethodID *string `json:"payment_method_id"`
	CategoryID      *string `json:"category_id"`
	Type            string  `json:"type"`
}

func parseExpenseType(t string) (model.ExpenseType, error) {
	switch model.ExpenseType(t) {
	case model.ExpenseTypeExpense, model.ExpenseTypeSaving, model.ExpenseTypeInvestment:
		return model.ExpenseType(t), nil
	case "":
		return model.ExpenseTypeExpense, nil
	default:
		return "", fiber.NewError(fiber.StatusBadRequest, "type must be one of: expense, saving, investment")
	}
}

func (h *ExpenseHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Amount == "" {
		return respondError(c, fiber.StatusBadRequest, "amount is required")
	}
	if req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "date is required")
	}

	expenseType, err := parseExpenseType(req.Type)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	amount, categoryID, paymentMethodID, date, err := parseMoneyFields(req.Amount, req.CategoryID, req.PaymentMethodID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	expense, err := h.svc.Create(userID, req.Name, amount, req.Currency, paymentMethodID, categoryID, date, expenseType)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create expense")
	}
	return respondCreated(c, expense)
}

func (h *ExpenseHandler) ListByYear(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	expenses, err := h.svc.ListByYear(userID, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list expenses")
	}
	return respondList(c, expenses)
}

func (h *ExpenseHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	expense, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "expense not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get expense")
	}
	return respondJSON(c, expense)
}

func (h *ExpenseHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Amount == "" {
		return respondError(c, fiber.StatusBadRequest, "amount is required")
	}
	if req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "date is required")
	}

	expenseType, err := parseExpenseType(req.Type)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	amount, categoryID, paymentMethodID, date, err := parseMoneyFields(req.Amount, req.CategoryID, req.PaymentMethodID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	expense, err := h.svc.Update(userID, id, req.Name, amount, req.Currency, paymentMethodID, categoryID, date, expenseType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "expense not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update expense")
	}
	return respondJSON(c, expense)
}

func (h *ExpenseHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete expense")
	}
	return respondNoContent(c)
}
