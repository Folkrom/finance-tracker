# Admin Middleware & Routes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add admin role detection from JWT claims, admin middleware, admin CRUD for global categories, and a basic stats endpoint.

**Architecture:** Auth middleware stores full JWT claims in Locals. New admin middleware reads `app_metadata.role` from claims. AdminHandler delegates to CategoryService (global ops) and AdminService (stats). Transactional delete reassigns FK references to "Other" before removing.

**Tech Stack:** Go/Fiber, GORM, PostgreSQL, golang-jwt/jwt/v5

---

### Task 1: Auth Middleware — Store full claims in Locals

**Files:**
- Modify: `backend/internal/middleware/auth.go`
- Modify: `backend/internal/middleware/auth_test.go`

- [ ] **Step 1: Update auth middleware to store claims**

In `backend/internal/middleware/auth.go`, add one line after `c.Locals("user_id", userID)`:

```go
c.Locals("claims", claims)
```

The full function already has `claims` as `jwt.MapClaims`. This just exposes it to downstream middleware.

- [ ] **Step 2: Add GetClaims helper**

Add this function after the existing `GetUserID` helper:

```go
// GetClaims extracts the full JWT claims from the Fiber context.
func GetClaims(c *fiber.Ctx) jwt.MapClaims {
	return c.Locals("claims").(jwt.MapClaims)
}
```

Add `"github.com/golang-jwt/jwt/v5"` to the import block (it's already there for `jwt.Keyfunc`).

- [ ] **Step 3: Update auth test to verify claims are stored**

In `backend/internal/middleware/auth_test.go`, update `TestAuthMiddleware_ValidToken` to also verify claims are stored. Replace the handler inside the test:

```go
app.Get("/test", func(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	claims := c.Locals("claims")
	return c.JSON(fiber.Map{"user_id": uid, "has_claims": claims != nil})
})
```

- [ ] **Step 4: Run tests**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/middleware/... -v
```

Expected: 3/3 pass.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/middleware/auth.go backend/internal/middleware/auth_test.go
git commit -m "feat: auth middleware stores full JWT claims in Locals"
```

---

### Task 2: Admin Middleware — Check app_metadata.role

**Files:**
- Create: `backend/internal/middleware/admin.go`

- [ ] **Step 1: Create admin middleware**

Create `backend/internal/middleware/admin.go`:

```go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// NewAdminMiddleware checks that the authenticated user has admin role
// in their JWT app_metadata. Must run after NewAuthMiddleware.
func NewAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.MapClaims)

		if !isAdmin(claims) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}

		return c.Next()
	}
}

// IsAdmin checks if the current request is from an admin user.
func IsAdmin(c *fiber.Ctx) bool {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return false
	}
	return isAdmin(claims)
}

func isAdmin(claims jwt.MapClaims) bool {
	appMeta, ok := claims["app_metadata"].(map[string]interface{})
	if !ok {
		return false
	}
	role, ok := appMeta["role"].(string)
	if !ok {
		return false
	}
	return role == "admin"
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/middleware/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/admin.go
git commit -m "feat: admin middleware — JWT app_metadata.role check"
```

---

### Task 3: AdminStats Model

**Files:**
- Create: `backend/internal/model/admin.go`

- [ ] **Step 1: Create AdminStats struct**

Create `backend/internal/model/admin.go`:

```go
package model

// AdminStats holds platform-level statistics. Not a DB model.
type AdminStats struct {
	Users            int `json:"users"`
	CategoriesGlobal int `json:"categories_global"`
	CategoriesUser   int `json:"categories_user"`
	Profiles         int `json:"profiles"`
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/model/admin.go
git commit -m "feat: AdminStats response model"
```

---

### Task 4: Category Repository — Admin operations

**Files:**
- Modify: `backend/internal/repository/category.go`

- [ ] **Step 1: Add GetGlobalByID method**

Add to `backend/internal/repository/category.go`:

```go
// GetGlobalByID returns a category only if it's global (user_id IS NULL).
func (r *CategoryRepository) GetGlobalByID(id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("id = ? AND user_id IS NULL", id).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	setGlobalFlagSingle(&cat)
	return &cat, nil
}
```

- [ ] **Step 2: Add CreateGlobal method**

```go
// CreateGlobal creates a global category (user_id = NULL).
func (r *CategoryRepository) CreateGlobal(cat *model.Category) error {
	cat.UserID = nil
	return r.db.Create(cat).Error
}
```

- [ ] **Step 3: Add ReassignAndDelete method**

```go
// ReassignAndDelete reassigns all FK references from one category to another,
// then deletes the original. Runs in a transaction.
func (r *CategoryRepository) ReassignAndDelete(categoryID, replacementID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tables := []string{"incomes", "expenses", "debts", "budgets", "wishlist_items"}
		for _, table := range tables {
			if err := tx.Exec(
				"UPDATE "+table+" SET category_id = ? WHERE category_id = ?",
				replacementID, categoryID,
			).Error; err != nil {
				return err
			}
		}
		return tx.Where("id = ?", categoryID).Delete(&model.Category{}).Error
	})
}
```

- [ ] **Step 4: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/repository/...
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/repository/category.go
git commit -m "feat: category repo — GetGlobalByID, CreateGlobal, ReassignAndDelete"
```

---

### Task 5: Admin Repository — Stats queries

**Files:**
- Create: `backend/internal/repository/admin.go`

- [ ] **Step 1: Create admin repository**

Create `backend/internal/repository/admin.go`:

```go
package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetStats() (*model.AdminStats, error) {
	var stats model.AdminStats

	// Count distinct users (from profiles, since auto-created)
	if err := r.db.Raw("SELECT COUNT(*) FROM profiles").Scan(&stats.Profiles).Error; err != nil {
		return nil, err
	}
	stats.Users = stats.Profiles // profiles are auto-created, so users ≈ profiles

	// Count global categories
	if err := r.db.Raw("SELECT COUNT(*) FROM categories WHERE user_id IS NULL").Scan(&stats.CategoriesGlobal).Error; err != nil {
		return nil, err
	}

	// Count user categories
	if err := r.db.Raw("SELECT COUNT(*) FROM categories WHERE user_id IS NOT NULL").Scan(&stats.CategoriesUser).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/repository/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/repository/admin.go
git commit -m "feat: admin repository — platform stats queries"
```

---

### Task 6: Category Service — Global admin operations

**Files:**
- Modify: `backend/internal/service/category.go`

- [ ] **Step 1: Add CreateGlobal method**

Add to `backend/internal/service/category.go`:

```go
func (s *CategoryService) CreateGlobal(name string, domain model.CategoryDomain, color *string, sortOrder int) (*model.Category, error) {
	cat := &model.Category{
		Name:      name,
		Domain:    domain,
		Color:     color,
		SortOrder: sortOrder,
	}
	if err := s.repo.CreateGlobal(cat); err != nil {
		return nil, err
	}
	return cat, nil
}
```

- [ ] **Step 2: Add UpdateGlobal method**

```go
func (s *CategoryService) UpdateGlobal(id uuid.UUID, name string, color *string, sortOrder *int) (*model.Category, error) {
	cat, err := s.repo.GetGlobalByID(id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		cat.Name = name
	}
	if color != nil {
		cat.Color = color
	}
	if sortOrder != nil {
		cat.SortOrder = *sortOrder
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}
```

- [ ] **Step 3: Add DeleteGlobal method**

```go
func (s *CategoryService) DeleteGlobal(id uuid.UUID) error {
	cat, err := s.repo.GetGlobalByID(id)
	if err != nil {
		return err
	}
	if cat.IsSystem {
		return ErrSystemCategoryProtected
	}
	other, err := s.repo.GetOtherCategory(cat.Domain)
	if err != nil {
		return err
	}
	return s.repo.ReassignAndDelete(cat.ID, other.ID)
}
```

- [ ] **Step 4: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/service/...
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/category.go
git commit -m "feat: category service — CreateGlobal, UpdateGlobal, DeleteGlobal"
```

---

### Task 7: Admin Service — Stats

**Files:**
- Create: `backend/internal/service/admin.go`

- [ ] **Step 1: Create admin service**

Create `backend/internal/service/admin.go`:

```go
package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetStats() (*model.AdminStats, error) {
	return s.repo.GetStats()
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/service/admin.go
git commit -m "feat: admin service — GetStats"
```

---

### Task 8: Admin Handler

**Files:**
- Create: `backend/internal/handler/admin.go`

- [ ] **Step 1: Create admin handler**

Create `backend/internal/handler/admin.go`:

```go
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
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/handler/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/admin.go
git commit -m "feat: admin handler — global category CRUD and stats"
```

---

### Task 9: Router and main.go — Wire admin routes

**Files:**
- Modify: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Update router.go**

Add `adminHandler *handler.AdminHandler` as the last parameter of the `Setup` function.

Add admin routes at the end of the function, before the closing `}`:

```go
	// Admin
	admin := api.Group("/admin", middleware.NewAdminMiddleware())
	adminCats := admin.Group("/categories")
	adminCats.Post("/", adminHandler.CreateCategory)
	adminCats.Put("/:id", adminHandler.UpdateCategory)
	adminCats.Delete("/:id", adminHandler.DeleteCategory)
	admin.Get("/stats", adminHandler.GetStats)
```

- [ ] **Step 2: Update main.go**

In Repositories section, add:
```go
adminRepo := repository.NewAdminRepository(db)
```

In Services section, add:
```go
adminSvc := service.NewAdminService(adminRepo)
```

In Handlers section, add:
```go
adminHandler := handler.NewAdminHandler(categorySvc, adminSvc)
```

Update `router.Setup` call to add `adminHandler` at the end:
```go
router.Setup(app, jwks.Keyfunc, profileRepo, categoryHandler, paymentMethodHandler, incomeHandler, expenseHandler, debtHandler, budgetHandler, dashboardHandler, cardHandler, wishlistHandler, profileHandler, adminHandler)
```

- [ ] **Step 3: Verify full build and tests**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./... && go test ./internal/middleware/... -v
```

Expected: compiles, 3/3 tests pass.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/router/router.go backend/cmd/server/main.go
git commit -m "feat: wire admin handler, middleware, and routes into router"
```

---

### Task 10: Integration Test — Verify admin routes

**Files:** None (manual verification)

- [ ] **Step 1: Start backend**

```bash
cd /home/folkrom/projects/finance-tracker && mise run dev-backend
```

- [ ] **Step 2: Verify non-admin gets 403**

With a regular (non-admin) user token:

```bash
curl -s http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: Bearer <token>" | jq
```

Expected: 403 `"admin access required"`

- [ ] **Step 3: Set admin role in Supabase**

In Supabase dashboard: Authentication → Users → select your user → Edit → `app_metadata`: `{"role": "admin"}`. Save.

Log out and back in to get a fresh JWT with the claim.

- [ ] **Step 4: Verify admin stats**

```bash
curl -s http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: Bearer <admin-token>" | jq
```

Expected: `{ users: 1, categories_global: 28, categories_user: 0, profiles: 1 }`

- [ ] **Step 5: Test admin create category**

```bash
curl -s -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Pets", "domain": "expense", "sort_order": 15}' | jq
```

Expected: 201 with new global category.

- [ ] **Step 6: Test admin delete category**

```bash
# Get the ID of the "Pets" category just created
PETS_ID=$(curl -s http://localhost:8080/api/v1/categories?domain=expense \
  -H "Authorization: Bearer <admin-token>" | jq -r '.data[] | select(.name=="Pets") | .id')

curl -s -X DELETE http://localhost:8080/api/v1/admin/categories/$PETS_ID \
  -H "Authorization: Bearer <admin-token>" -w "\n%{http_code}\n"
```

Expected: 204 No Content.

- [ ] **Step 7: Verify system category cannot be deleted**

```bash
OTHER_ID=$(curl -s http://localhost:8080/api/v1/categories?domain=expense \
  -H "Authorization: Bearer <admin-token>" | jq -r '.data[] | select(.name=="Other") | .id')

curl -s -X DELETE http://localhost:8080/api/v1/admin/categories/$OTHER_ID \
  -H "Authorization: Bearer <admin-token>" | jq
```

Expected: 403 `"cannot delete system categories"`
