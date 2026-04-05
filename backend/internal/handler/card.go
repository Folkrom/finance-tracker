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

type CardHandler struct {
	svc *service.CardService
}

func NewCardHandler(svc *service.CardService) *CardHandler {
	return &CardHandler{svc: svc}
}

type createCardRequest struct {
	PaymentMethodID     string  `json:"payment_method_id"`
	Bank                string  `json:"bank"`
	CardLimit           string  `json:"card_limit"`
	RecommendedMaxPct   string  `json:"recommended_max_pct"`
	ManualUsageOverride *string `json:"manual_usage_override"`
	Level               *string `json:"level"`
}

type updateCardRequest struct {
	Bank                string  `json:"bank"`
	CardLimit           string  `json:"card_limit"`
	RecommendedMaxPct   string  `json:"recommended_max_pct"`
	ManualUsageOverride *string `json:"manual_usage_override"`
	Level               *string `json:"level"`
}

func (h *CardHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req createCardRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.PaymentMethodID == "" {
		return respondError(c, fiber.StatusBadRequest, "payment_method_id is required")
	}
	if req.Bank == "" {
		return respondError(c, fiber.StatusBadRequest, "bank is required")
	}
	if req.CardLimit == "" {
		return respondError(c, fiber.StatusBadRequest, "card_limit is required")
	}

	pmID, err := uuid.Parse(req.PaymentMethodID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid payment_method_id")
	}

	cardLimit, err := decimal.NewFromString(req.CardLimit)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid card_limit")
	}

	recommendedMaxPct := decimal.NewFromFloat(30.0)
	if req.RecommendedMaxPct != "" {
		recommendedMaxPct, err = decimal.NewFromString(req.RecommendedMaxPct)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid recommended_max_pct")
		}
	}

	var manualOverride *decimal.Decimal
	if req.ManualUsageOverride != nil && *req.ManualUsageOverride != "" {
		mo, err := decimal.NewFromString(*req.ManualUsageOverride)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid manual_usage_override")
		}
		manualOverride = &mo
	}

	card, err := h.svc.Create(userID, pmID, req.Bank, cardLimit, recommendedMaxPct, manualOverride, req.Level)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	return respondCreated(c, card)
}

func (h *CardHandler) GetSummaries(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	monthStr := c.Query("month")
	yearStr := c.Query("year")
	if monthStr == "" || yearStr == "" {
		return respondError(c, fiber.StatusBadRequest, "month and year query params are required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return respondError(c, fiber.StatusBadRequest, "invalid month")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	summaries, err := h.svc.GetCardSummaries(userID, month, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load card summaries")
	}

	if summaries == nil {
		summaries = []service.CardSummary{}
	}
	return respondJSON(c, fiber.Map{"data": summaries, "total": len(summaries)})
}

func (h *CardHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	card, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "card not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get card")
	}
	return respondJSON(c, card)
}

func (h *CardHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Bank == "" {
		return respondError(c, fiber.StatusBadRequest, "bank is required")
	}
	if req.CardLimit == "" {
		return respondError(c, fiber.StatusBadRequest, "card_limit is required")
	}

	cardLimit, err := decimal.NewFromString(req.CardLimit)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid card_limit")
	}

	recommendedMaxPct := decimal.NewFromFloat(30.0)
	if req.RecommendedMaxPct != "" {
		recommendedMaxPct, err = decimal.NewFromString(req.RecommendedMaxPct)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid recommended_max_pct")
		}
	}

	var manualOverride *decimal.Decimal
	if req.ManualUsageOverride != nil && *req.ManualUsageOverride != "" {
		mo, err := decimal.NewFromString(*req.ManualUsageOverride)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid manual_usage_override")
		}
		manualOverride = &mo
	}

	card, err := h.svc.Update(userID, id, req.Bank, cardLimit, recommendedMaxPct, manualOverride, req.Level)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "card not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update card")
	}
	return respondJSON(c, card)
}

func (h *CardHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete card")
	}
	return respondNoContent(c)
}
