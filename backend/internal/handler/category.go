package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

type createCategoryRequest struct {
	Name   string               `json:"name"`
	Domain model.CategoryDomain `json:"domain"`
	Color  *string              `json:"color"`
}

type updateCategoryRequest struct {
	Name  string  `json:"name"`
	Color *string `json:"color"`
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req createCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Domain == "" {
		return respondError(c, fiber.StatusBadRequest, "domain is required")
	}

	cat, err := h.svc.Create(userID, req.Name, req.Domain, req.Color)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create category")
	}
	return respondCreated(c, cat)
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	domain := model.CategoryDomain(c.Query("domain"))
	if domain == "" {
		return respondError(c, fiber.StatusBadRequest, "domain query param is required")
	}

	cats, err := h.svc.ListByDomain(userID, domain)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list categories")
	}
	return respondList(c, cats)
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	cat, err := h.svc.GetByID(userID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "category not found")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to get category")
	}
	return respondJSON(c, cat)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}

	cat, err := h.svc.Update(userID, id, req.Name, req.Color)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "category not found")
		}
		if err == service.ErrGlobalCategoryReadOnly {
			return respondError(c, fiber.StatusForbidden, "cannot modify global categories")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update category")
	}
	return respondJSON(c, cat)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return respondError(c, fiber.StatusNotFound, "category not found")
		}
		if err == service.ErrGlobalCategoryReadOnly || err == service.ErrSystemCategoryProtected {
			return respondError(c, fiber.StatusForbidden, err.Error())
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to delete category")
	}
	return respondNoContent(c)
}

func (h *CategoryHandler) SeedDefaults(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "global categories are available by default",
	})
}
