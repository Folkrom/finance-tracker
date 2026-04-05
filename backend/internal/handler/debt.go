package handler

import (
	"strconv"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DebtHandler struct {
	svc *service.DebtService
}

func NewDebtHandler(svc *service.DebtService) *DebtHandler {
	return &DebtHandler{svc: svc}
}

type createDebtRequest struct {
	Name            string  `json:"name"`
	Amount          string  `json:"amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
	PaymentMethodID *string `json:"payment_method_id"`
	CategoryID      *string `json:"category_id"`
}

type updateDebtRequest struct {
	Name            string  `json:"name"`
	Amount          string  `json:"amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
	PaymentMethodID *string `json:"payment_method_id"`
	CategoryID      *string `json:"category_id"`
}

func (h *DebtHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createDebtRequest
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

	amount, categoryID, paymentMethodID, date, err := parseMoneyFields(req.Amount, req.CategoryID, req.PaymentMethodID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	debt, err := h.svc.Create(userID, req.Name, amount, req.Currency, paymentMethodID, categoryID, date)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create debt")
	}
	return respondCreated(c, debt)
}

func (h *DebtHandler) ListByYear(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	debts, err := h.svc.ListByYear(userID, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list debts")
	}
	return respondList(c, debts)
}

func (h *DebtHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	debt, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "debt not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get debt")
	}
	return respondJSON(c, debt)
}

func (h *DebtHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateDebtRequest
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

	amount, categoryID, paymentMethodID, date, err := parseMoneyFields(req.Amount, req.CategoryID, req.PaymentMethodID, req.Date)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	debt, err := h.svc.Update(userID, id, req.Name, amount, req.Currency, paymentMethodID, categoryID, date)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "debt not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update debt")
	}
	return respondJSON(c, debt)
}

func (h *DebtHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete debt")
	}
	return respondNoContent(c)
}
