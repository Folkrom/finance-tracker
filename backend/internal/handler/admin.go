package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminHandler struct {
	categorySvc *service.CategoryService
	adminSvc    *service.AdminService
}

func NewAdminHandler(categorySvc *service.CategoryService, adminSvc *service.AdminService) *AdminHandler {
	return &AdminHandler{
		categorySvc: categorySvc,
		adminSvc:    adminSvc,
	}
}

type createGlobalCategoryRequest struct {
	Name      string               `json:"name"`
	Domain    model.CategoryDomain `json:"domain"`
	Color     *string              `json:"color"`
	SortOrder int                  `json:"sort_order"`
}

type updateGlobalCategoryRequest struct {
	Name      string  `json:"name"`
	Color     *string `json:"color"`
	SortOrder *int    `json:"sort_order"`
}

func (h *AdminHandler) CreateCategory(c *fiber.Ctx) error {
	var req createGlobalCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Domain == "" {
		return respondError(c, fiber.StatusBadRequest, "domain is required")
	}

	cat, err := h.categorySvc.CreateGlobal(req.Name, req.Domain, req.Color, req.SortOrder)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create global category")
	}
	return respondCreated(c, cat)
}

func (h *AdminHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateGlobalCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}

	cat, err := h.categorySvc.UpdateGlobal(id, req.Name, req.Color, req.SortOrder)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "global category not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update global category")
	}
	return respondJSON(c, cat)
}

func (h *AdminHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.categorySvc.DeleteGlobal(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "global category not found")
		}
		if err == service.ErrSystemCategoryProtected {
			return respondError(c, fiber.StatusForbidden, "cannot delete system categories")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to delete global category")
	}
	return respondNoContent(c)
}

func (h *AdminHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.adminSvc.GetStats()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to get stats")
	}
	return respondJSON(c, stats)
}
