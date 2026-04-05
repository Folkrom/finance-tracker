package handler

import (
	"strconv"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	var month *int
	if m := c.Query("month"); m != "" {
		mv, err := strconv.Atoi(m)
		if err != nil || mv < 1 || mv > 12 {
			return respondError(c, fiber.StatusBadRequest, "invalid month")
		}
		month = &mv
	}

	data, err := h.svc.GetDashboard(userID, year, month)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load dashboard")
	}
	return respondJSON(c, data)
}
