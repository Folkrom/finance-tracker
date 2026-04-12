package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	svc *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

type updateProfileRequest struct {
	Currency string `json:"currency"`
	Language string `json:"language"`
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	profile, err := h.svc.GetOrCreate(userID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to get profile")
	}
	return respondJSON(c, profile)
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}

	profile, err := h.svc.Update(userID, req.Currency, req.Language)
	if err != nil {
		if err == service.ErrInvalidCurrency {
			return respondError(c, fiber.StatusBadRequest, "invalid currency — allowed: MXN, USD, EUR, GBP, BRL, COP, ARS")
		}
		if err == service.ErrInvalidLanguage {
			return respondError(c, fiber.StatusBadRequest, "invalid language — allowed: en, es")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update profile")
	}
	return respondJSON(c, profile)
}
