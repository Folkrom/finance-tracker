# Admin Middleware & Routes Design

**Date:** 2026-04-12
**Status:** Approved
**Scope:** Admin JWT middleware, admin global category CRUD, basic stats endpoint. No frontend.

---

## Problem

Global categories can only be managed via SQL — there's no API for creating, updating, or deleting them. An admin role system is needed to gate these operations.

## Design

### Admin Detection

Supabase `app_metadata.role = "admin"` in JWT claims. Set via Supabase dashboard or Admin API. The JWT already contains `app_metadata` — the auth middleware parses it but currently only extracts `sub`.

### Admin Middleware

```go
func NewAdminMiddleware() fiber.Handler
```

Extracts `app_metadata` from JWT claims (stored in `c.Locals` by auth middleware). Checks `role == "admin"`. Returns 403 `"admin access required"` if not.

**Implementation detail:** The auth middleware currently parses the JWT and stores `user_id` in Locals. To avoid re-parsing the JWT, the auth middleware should also store the full `MapClaims` in `c.Locals("claims")`. The admin middleware then reads from Locals.

Helper function:
```go
func IsAdmin(c *fiber.Ctx) bool
```

### Auth Middleware Change

Store full claims in Locals so downstream middleware can access them:
```go
c.Locals("user_id", userID)
c.Locals("claims", claims) // NEW — full jwt.MapClaims
```

### Admin Routes

All under `/api/v1/admin`, protected by both auth + admin middleware.

```
POST   /api/v1/admin/categories         — create global category
PUT    /api/v1/admin/categories/:id     — update global category
DELETE /api/v1/admin/categories/:id     — delete global category
GET    /api/v1/admin/stats              — basic platform stats
```

### Admin Category CRUD

**Create:**
- Body: `{ name, domain, color?, sort_order? }`
- Inserts with `user_id = NULL` (global), `is_system = false`
- Returns 201 with created category

**Update:**
- Body: `{ name, color?, sort_order? }`
- Fetches category, rejects if `user_id IS NOT NULL` (not global) — 400
- Updates allowed fields
- Returns 200 with updated category

**Delete:**
- Rejects if `is_system = true` — 403 `"cannot delete system categories"`
- Rejects if `user_id IS NOT NULL` (not global) — 400
- Finds "Other" system category for same domain
- In a transaction: reassigns all FK references (incomes, expenses, debts, budgets, wishlist_items) from deleted category to "Other", then deletes
- Returns 204

### Stats Endpoint

**GET /api/v1/admin/stats** returns:
```json
{
  "users": 42,
  "categories_global": 28,
  "categories_user": 5,
  "profiles": 42
}
```

Simple COUNT queries on categories (split by `user_id IS NULL` vs not), profiles, and `count(distinct user_id)` from profiles (as a proxy for active users since profiles are auto-created).

### Repository Changes

**CategoryRepository — new method:**
```go
func (r *CategoryRepository) ReassignAndDelete(tx *gorm.DB, categoryID, replacementID uuid.UUID) error
```
Updates all FK references in a transaction, then deletes the category.

**CategoryRepository — new method:**
```go
func (r *CategoryRepository) GetGlobalByID(id uuid.UUID) (*model.Category, error)
```
Gets a category only if `user_id IS NULL`. For admin operations that should only affect globals.

**CategoryRepository — new method:**
```go
func (r *CategoryRepository) CreateGlobal(cat *model.Category) error
```
Same as Create but explicitly sets `user_id = nil`.

**New AdminRepository:**
```go
type AdminRepository struct { db *gorm.DB }

func (r *AdminRepository) GetStats() (*model.AdminStats, error)
```

### New Model

```go
type AdminStats struct {
    Users           int `json:"users"`
    CategoriesGlobal int `json:"categories_global"`
    CategoriesUser   int `json:"categories_user"`
    Profiles        int `json:"profiles"`
}
```

Not a DB model — just a response struct. Lives in model package for consistency.

### Service Layer

**AdminService** — thin wrapper for stats:
```go
type AdminService struct { repo *AdminRepository }
func (s *AdminService) GetStats() (*model.AdminStats, error)
```

**CategoryService** — new methods for admin operations:
```go
func (s *CategoryService) CreateGlobal(name string, domain CategoryDomain, color *string, sortOrder int) (*Category, error)
func (s *CategoryService) UpdateGlobal(id uuid.UUID, name string, color *string, sortOrder *int) (*Category, error)
func (s *CategoryService) DeleteGlobal(id uuid.UUID) error
```

### Handler

**AdminHandler** — single handler file for all admin routes:
```go
type AdminHandler struct {
    categorySvc *service.CategoryService
    adminSvc    *service.AdminService
}

func (h *AdminHandler) CreateCategory(c *fiber.Ctx) error
func (h *AdminHandler) UpdateCategory(c *fiber.Ctx) error
func (h *AdminHandler) DeleteCategory(c *fiber.Ctx) error
func (h *AdminHandler) GetStats(c *fiber.Ctx) error
```

### Router Wiring

```go
admin := api.Group("/admin", middleware.NewAdminMiddleware())
adminCats := admin.Group("/categories")
adminCats.Post("/", adminHandler.CreateCategory)
adminCats.Put("/:id", adminHandler.UpdateCategory)
adminCats.Delete("/:id", adminHandler.DeleteCategory)
admin.Get("/stats", adminHandler.GetStats)
```

### What This Does NOT Include

- Admin frontend pages (later Plan 5 phase)
- User management routes (no consumer yet)
- Audit logging (future consideration)
- Rate limiting on admin routes (single user, dev mode)
