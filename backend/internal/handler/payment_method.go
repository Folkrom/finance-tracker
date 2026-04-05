package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodHandler struct {
	svc *service.PaymentMethodService
}

func NewPaymentMethodHandler(svc *service.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{svc: svc}
}

type createPaymentMethodRequest struct {
	Name    string                  `json:"name"`
	Type    model.PaymentMethodType `json:"type"`
	Details *string                 `json:"details"`
}

type updatePaymentMethodRequest struct {
	Name    string                  `json:"name"`
	Type    model.PaymentMethodType `json:"type"`
	Details *string                 `json:"details"`
}

func (h *PaymentMethodHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createPaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Type == "" {
		return respondError(c, fiber.StatusBadRequest, "type is required")
	}

	pm, err := h.svc.Create(userID, req.Name, req.Type, req.Details)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create payment method")
	}
	return respondCreated(c, pm)
}

func (h *PaymentMethodHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	pmType := model.PaymentMethodType(c.Query("type"))

	pms, err := h.svc.List(userID, pmType)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list payment methods")
	}
	return respondList(c, pms)
}

func (h *PaymentMethodHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	pm, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "payment method not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get payment method")
	}
	return respondJSON(c, pm)
}

func (h *PaymentMethodHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updatePaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Type == "" {
		return respondError(c, fiber.StatusBadRequest, "type is required")
	}

	pm, err := h.svc.Update(userID, id, req.Name, req.Type, req.Details)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "payment method not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update payment method")
	}
	return respondJSON(c, pm)
}

func (h *PaymentMethodHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete payment method")
	}
	return respondNoContent(c)
}
