package handler

import (
	"strings"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WishlistItemHandler struct {
	svc *service.WishlistItemService
}

func NewWishlistItemHandler(svc *service.WishlistItemService) *WishlistItemHandler {
	return &WishlistItemHandler{svc: svc}
}

type createWishlistItemRequest struct {
	Name                 string   `json:"name"`
	ImageURL             *string  `json:"image_url"`
	Price                *string  `json:"price"`
	Currency             string   `json:"currency"`
	Links                []string `json:"links"`
	CategoryID           *string  `json:"category_id"`
	Priority             string   `json:"priority"`
	Status               string   `json:"status"`
	TargetDate           *string  `json:"target_date"`
	MonthlyContribution  *string  `json:"monthly_contribution"`
	ContributionCurrency string   `json:"contribution_currency"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func parseWishlistFields(req *createWishlistItemRequest) (
	*decimal.Decimal, *uuid.UUID, *time.Time, *decimal.Decimal, error,
) {
	var price *decimal.Decimal
	if req.Price != nil && *req.Price != "" {
		p, err := decimal.NewFromString(*req.Price)
		if err != nil {
			return nil, nil, nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid price")
		}
		price = &p
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		parsed, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, nil, nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid category_id")
		}
		categoryID = &parsed
	}

	var targetDate *time.Time
	if req.TargetDate != nil && *req.TargetDate != "" {
		t, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return nil, nil, nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid target_date, expected YYYY-MM-DD")
		}
		targetDate = &t
	}

	var monthlyContribution *decimal.Decimal
	if req.MonthlyContribution != nil && *req.MonthlyContribution != "" {
		mc, err := decimal.NewFromString(*req.MonthlyContribution)
		if err != nil {
			return nil, nil, nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid monthly_contribution")
		}
		monthlyContribution = &mc
	}

	return price, categoryID, targetDate, monthlyContribution, nil
}

func (h *WishlistItemHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createWishlistItemRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}

	price, categoryID, targetDate, monthlyContribution, err := parseWishlistFields(&req)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	item, err := h.svc.Create(
		userID, req.Name, req.ImageURL, price, req.Currency, req.Links,
		categoryID, model.WishlistPriority(req.Priority), model.WishlistStatus(req.Status),
		targetDate, monthlyContribution, req.ContributionCurrency,
	)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create wishlist item")
	}
	return respondCreated(c, item)
}

func (h *WishlistItemHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	statusFilter := c.Query("status")
	if statusFilter != "" {
		parts := strings.Split(statusFilter, ",")
		statuses := make([]model.WishlistStatus, len(parts))
		for i, s := range parts {
			statuses[i] = model.WishlistStatus(strings.TrimSpace(s))
		}
		items, err := h.svc.ListByStatus(userID, statuses)
		if err != nil {
			return respondError(c, fiber.StatusInternalServerError, "failed to list wishlist items")
		}
		return respondList(c, items)
	}

	items, err := h.svc.ListByUser(userID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list wishlist items")
	}
	return respondList(c, items)
}

func (h *WishlistItemHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	item, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "wishlist item not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get wishlist item")
	}
	return respondJSON(c, item)
}

func (h *WishlistItemHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req createWishlistItemRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}

	price, categoryID, targetDate, monthlyContribution, err := parseWishlistFields(&req)
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			return respondError(c, fe.Code, fe.Message)
		}
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	item, err := h.svc.Update(
		userID, id, req.Name, req.ImageURL, price, req.Currency, req.Links,
		categoryID, model.WishlistPriority(req.Priority), model.WishlistStatus(req.Status),
		targetDate, monthlyContribution, req.ContributionCurrency,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "wishlist item not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update wishlist item")
	}
	return respondJSON(c, item)
}

func (h *WishlistItemHandler) UpdateStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Status == "" {
		return respondError(c, fiber.StatusBadRequest, "status is required")
	}

	if err := h.svc.UpdateStatus(userID, id, model.WishlistStatus(req.Status)); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to update status")
	}
	return respondNoContent(c)
}

func (h *WishlistItemHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to delete wishlist item")
	}
	return respondNoContent(c)
}
