# Category Revamp Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace per-user category seeding with global system categories (nullable user_id) while preserving user-created custom categories.

**Architecture:** Add `is_system` boolean and make `user_id` nullable on categories table. Migration seeds global defaults and deduplicates existing user data. Repository queries change to `WHERE (user_id IS NULL OR user_id = ?)`. Handler guards prevent modifying/deleting global categories.

**Tech Stack:** Go/Fiber, GORM, PostgreSQL, golang-migrate, Next.js/React, TypeScript

---

### Task 1: Migration — Schema Changes and Global Category Seeding

**Files:**
- Create: `backend/migrations/000009_category_revamp.up.sql`
- Create: `backend/migrations/000009_category_revamp.down.sql`

- [ ] **Step 1: Write the up migration**

Create `backend/migrations/000009_category_revamp.up.sql`:

```sql
-- Step 1: Add is_system column
ALTER TABLE categories ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT false;

-- Step 2: Make user_id nullable
ALTER TABLE categories ALTER COLUMN user_id DROP NOT NULL;

-- Step 3: Drop old unique constraint
ALTER TABLE categories DROP CONSTRAINT categories_user_id_name_domain_key;

-- Step 4: Insert global categories (user_id = NULL)
-- Income
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Salary', 'income', 0, false),
  ('Bonus', 'income', 1, false),
  ('Freelance', 'income', 2, false),
  ('Dividends', 'income', 3, false),
  ('Interest', 'income', 4, false),
  ('Side Hustle', 'income', 5, false),
  ('Other', 'income', 6, true);

-- Expense
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Home Expenses', 'expense', 0, false),
  ('Eating Out', 'expense', 1, false),
  ('Self Care', 'expense', 2, false),
  ('Coffee/Drink', 'expense', 3, false),
  ('Entertainment', 'expense', 4, false),
  ('Transportation', 'expense', 5, false),
  ('Groceries', 'expense', 6, false),
  ('Utilities', 'expense', 7, false),
  ('Clothes', 'expense', 8, false),
  ('Card Payments', 'expense', 9, false),
  ('Savings/Investment', 'expense', 10, false),
  ('Taxes', 'expense', 11, false),
  ('Knowledge', 'expense', 12, false),
  ('Tech', 'expense', 13, false),
  ('Other', 'expense', 14, true);

-- Wishlist
INSERT INTO categories (name, domain, sort_order, is_system) VALUES
  ('Electronics', 'wishlist', 0, false),
  ('Clothing', 'wishlist', 1, false),
  ('Home & Kitchen', 'wishlist', 2, false),
  ('Books & Media', 'wishlist', 3, false),
  ('Sports & Outdoors', 'wishlist', 4, false),
  ('Other', 'wishlist', 5, true);

-- Step 5: Reassign FK references from user categories to matching globals, then delete user dupes
-- For each table with category_id, update rows where the user's category matches a global by name+domain
DO $$
DECLARE
  user_cat RECORD;
  global_id UUID;
BEGIN
  FOR user_cat IN
    SELECT uc.id AS user_cat_id, uc.name, uc.domain, uc.user_id
    FROM categories uc
    WHERE uc.user_id IS NOT NULL
      AND EXISTS (
        SELECT 1 FROM categories gc
        WHERE gc.user_id IS NULL
          AND gc.name = uc.name
          AND gc.domain = uc.domain
      )
  LOOP
    SELECT gc.id INTO global_id
    FROM categories gc
    WHERE gc.user_id IS NULL
      AND gc.name = user_cat.name
      AND gc.domain = user_cat.domain;

    -- Update all FK references
    UPDATE incomes SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE expenses SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE debts SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE budgets SET category_id = global_id WHERE category_id = user_cat.user_cat_id;
    UPDATE wishlist_items SET category_id = global_id WHERE category_id = user_cat.user_cat_id;

    -- Delete the user's duplicate
    DELETE FROM categories WHERE id = user_cat.user_cat_id;
  END LOOP;
END $$;

-- Step 6: Create partial unique indexes
CREATE UNIQUE INDEX idx_categories_global_unique
  ON categories (name, domain) WHERE user_id IS NULL;
CREATE UNIQUE INDEX idx_categories_user_unique
  ON categories (user_id, name, domain) WHERE user_id IS NOT NULL;
```

- [ ] **Step 2: Write the down migration**

Create `backend/migrations/000009_category_revamp.down.sql`:

```sql
-- Remove partial unique indexes
DROP INDEX IF EXISTS idx_categories_global_unique;
DROP INDEX IF EXISTS idx_categories_user_unique;

-- Delete all global categories
DELETE FROM categories WHERE user_id IS NULL;

-- Restore NOT NULL on user_id
ALTER TABLE categories ALTER COLUMN user_id SET NOT NULL;

-- Restore original unique constraint
ALTER TABLE categories ADD CONSTRAINT categories_user_id_name_domain_key UNIQUE (user_id, name, domain);

-- Drop is_system column
ALTER TABLE categories DROP COLUMN is_system;
```

- [ ] **Step 3: Run the migration**

```bash
cd /home/folkrom/projects/finance-tracker && mise run migrate-up
```

Expected: migration 000009 applied successfully.

- [ ] **Step 4: Verify migration**

```bash
docker exec -it $(docker ps -q -f name=postgres) psql -U finance -d finance_tracker -c "
  SELECT count(*) AS global_count FROM categories WHERE user_id IS NULL;
  SELECT count(*) AS user_count FROM categories WHERE user_id IS NOT NULL;
  SELECT name, domain, is_system FROM categories WHERE user_id IS NULL ORDER BY domain, sort_order;
"
```

Expected: 28 global categories (7 income + 15 expense + 6 wishlist), 0 user dupes of global names.

- [ ] **Step 5: Commit**

```bash
git add backend/migrations/000009_category_revamp.up.sql backend/migrations/000009_category_revamp.down.sql
git commit -m "feat: migration 000009 — global categories with nullable user_id"
```

---

### Task 2: Model — Add IsSystem Field, Make UserID Nullable on Category

**Files:**
- Modify: `backend/internal/model/category.go`

- [ ] **Step 1: Update Category model**

Add `IsSystem` field and override `UserID` as pointer to avoid changing the shared `Base` struct:

Replace the full `Category` struct in `backend/internal/model/category.go`:

```go
type Category struct {
	Base
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Domain    CategoryDomain `gorm:"type:varchar(50);not null" json:"domain"`
	Color     *string        `gorm:"type:varchar(7)" json:"color,omitempty"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsSystem  bool           `gorm:"not null;default:false" json:"is_system"`
}
```

Note: `UserID` on Category overrides `Base.UserID`. GORM uses the most specific field. The `Base.UserID` (non-pointer) remains unchanged for all other models.

Add import for `uuid` at the top of the file:

```go
import "github.com/google/uuid"
```

Add a computed JSON field helper:

```go
func (c Category) IsGlobal() bool {
	return c.UserID == nil
}
```

- [ ] **Step 2: Verify it compiles**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./...
```

Expected: compilation succeeds. There will be type errors in repository/service/handler — those are fixed in the next tasks.

Note: if `go build ./...` fails due to downstream code expecting `uuid.UUID` instead of `*uuid.UUID`, that's expected. We fix those in Tasks 3-5. Just confirm the model file itself is valid by checking:

```bash
go build ./internal/model/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/model/category.go
git commit -m "feat: Category model — nullable UserID, IsSystem field"
```

---

### Task 3: Repository — Update Queries for Global + User Categories

**Files:**
- Modify: `backend/internal/repository/category.go`
- Modify: `backend/internal/repository/category_test.go`

- [ ] **Step 1: Update repository methods**

Replace the contents of `backend/internal/repository/category.go`:

```go
package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

// ListByDomain returns global categories (user_id IS NULL) and user's own categories.
func (r *CategoryRepository) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	var cats []model.Category
	err := r.db.
		Where("(user_id IS NULL OR user_id = ?) AND domain = ?", userID, domain).
		Order("sort_order ASC, name ASC").
		Find(&cats).Error
	return cats, err
}

// GetByID returns a category if it's global or owned by the user.
func (r *CategoryRepository) GetByID(userID, id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", id, userID).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(userID, id uuid.UUID) error {
	return r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Category{}).Error
}

// GetOtherCategory returns the "Other" system category for a given domain.
// Used when admin-deleting a global category to reassign references.
func (r *CategoryRepository) GetOtherCategory(domain model.CategoryDomain) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("user_id IS NULL AND domain = ? AND is_system = true AND name = 'Other'", domain).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}
```

Key changes:
- `ListByDomain`: `user_id = ?` → `(user_id IS NULL OR user_id = ?)`
- `GetByID`: `user_id = ?` → `(user_id IS NULL OR user_id = ?)`
- `Delete`: keeps `user_id = ?` (not IS NULL) — users can only delete their own
- New: `GetOtherCategory` for future admin delete reassignment

- [ ] **Step 2: Update repository tests**

Replace the contents of `backend/internal/repository/category_test.go`:

```go
package repository_test

import (
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		Base:   model.Base{},
		UserID: &userID,
		Name:   "Salary",
		Domain: model.CategoryDomainIncome,
	}

	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cat.ID)
}

func TestCategoryRepository_ListByDomain_IncludesGlobal(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	// Create a global category (user_id = nil)
	global := &model.Category{
		Name:   "Global Salary",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(global))

	// Create a user category
	user := &model.Category{
		UserID: &userID,
		Name:   "My Side Gig",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(user))

	// User should see both global and their own
	cats, err := repo.ListByDomain(userID, model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, cats, 2)

	// Other user should see only global
	other, err := repo.ListByDomain(uuid.New(), model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, other, 1)
	assert.Equal(t, "Global Salary", other[0].Name)
}

func TestCategoryRepository_GetByID_GlobalCategory(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)

	global := &model.Category{
		Name:   "Global Expense",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, repo.Create(global))

	// Any user can read a global category
	cat, err := repo.GetByID(uuid.New(), global.ID)
	require.NoError(t, err)
	assert.Equal(t, "Global Expense", cat.Name)
}

func TestCategoryRepository_Delete_CannotDeleteGlobal(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)

	global := &model.Category{
		Name:   "Undeletable",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, repo.Create(global))

	// Attempting to delete a global category with any user_id should not delete it
	// because Delete requires user_id = ? (not IS NULL)
	err := repo.Delete(uuid.New(), global.ID)
	require.NoError(t, err) // no error, but 0 rows affected

	// Global category should still exist
	cat, err := repo.GetByID(uuid.New(), global.ID)
	require.NoError(t, err)
	assert.Equal(t, "Undeletable", cat.Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		UserID: &userID,
		Name:   "Old Name",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(cat))

	cat.Name = "New Name"
	err := repo.Update(cat)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/repository/...
```

Expected: compiles. (Tests won't run without DB, but they should compile.)

- [ ] **Step 4: Commit**

```bash
git add backend/internal/repository/category.go backend/internal/repository/category_test.go
git commit -m "feat: category repo — query global+user categories, guard global deletes"
```

---

### Task 4: Service — Guard Global Category Mutations, Deprecate Seed

**Files:**
- Modify: `backend/internal/service/category.go`

- [ ] **Step 1: Update service with guards**

Replace the contents of `backend/internal/service/category.go`:

```go
package service

import (
	"errors"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrGlobalCategoryReadOnly = errors.New("cannot modify global categories")
	ErrSystemCategoryProtected = errors.New("cannot delete system categories")
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID uuid.UUID, name string, domain model.CategoryDomain, color *string) (*model.Category, error) {
	cat := &model.Category{
		UserID: &userID,
		Name:   name,
		Domain: domain,
		Color:  color,
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	return s.repo.ListByDomain(userID, domain)
}

func (s *CategoryService) GetByID(userID, id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(userID, id)
}

func (s *CategoryService) Update(userID, id uuid.UUID, name string, color *string) (*model.Category, error) {
	cat, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if cat.IsGlobal() {
		return nil, ErrGlobalCategoryReadOnly
	}
	cat.Name = name
	cat.Color = color
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(userID, id uuid.UUID) error {
	cat, err := s.repo.GetByID(userID, id)
	if err != nil {
		return err
	}
	if cat.IsSystem {
		return ErrSystemCategoryProtected
	}
	if cat.IsGlobal() {
		return ErrGlobalCategoryReadOnly
	}
	return s.repo.Delete(userID, id)
}

// SeedDefaults is deprecated — global categories exist from migration.
func (s *CategoryService) SeedDefaults(_ uuid.UUID) error {
	return nil
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/service/...
```

Expected: compiles. The handler still references old service API, but service itself should be clean since the method signatures didn't change.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/category.go
git commit -m "feat: category service — guard global/system mutations, deprecate seed"
```

---

### Task 5: Handler — Return 403 for Global Category Mutations

**Files:**
- Modify: `backend/internal/handler/category.go`

- [ ] **Step 1: Update handler with error mapping**

Replace the `Update`, `Delete`, and `SeedDefaults` methods in `backend/internal/handler/category.go`:

```go
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
```

Update the import block at the top of `category.go` to include `service`:

```go
import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
```

- [ ] **Step 2: Verify full backend compiles**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./...
```

Expected: compiles with no errors.

- [ ] **Step 3: Run auth middleware tests (the ones that don't need DB)**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/middleware/... -v
```

Expected: 3/3 pass.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/handler/category.go
git commit -m "feat: category handler — 403 on global mutations, deprecate seed endpoint"
```

---

### Task 6: Frontend — Update Category Type and Remove Seed Button

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/components/settings/category-manager.tsx`

- [ ] **Step 1: Update the Category TypeScript type**

In `frontend/src/types/index.ts`, update the Category interface:

```typescript
export interface Category {
  id: string;
  user_id: string | null;
  name: string;
  domain: "income" | "expense" | "wishlist";
  color?: string;
  sort_order: number;
  is_system: boolean;
  is_global: boolean;
  created_at: string;
  updated_at: string;
}
```

Note: `user_id` is now `string | null` (null for globals). `is_system` and `is_global` are new boolean fields.

Wait — the backend doesn't currently serialize `is_global`. We need to add it to the JSON response. Two options:
- Add a custom JSON marshaler on Category
- Add `IsGlobal` as a computed GORM field

Simpler: add it in the handler response. But actually, the frontend can derive `is_global` from `user_id === null`. Let's keep it simple:

```typescript
export interface Category {
  id: string;
  user_id: string | null;
  name: string;
  domain: "income" | "expense" | "wishlist";
  color?: string;
  sort_order: number;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}
```

The frontend derives `isGlobal` with: `category.user_id === null`.

- [ ] **Step 2: Update category-manager.tsx — remove Seed button, add global guards**

In `frontend/src/components/settings/category-manager.tsx`, make these changes:

**Remove the seed button and its handler.** Delete `handleSeedDefaults`, the `seeding` state, and the Seed Defaults `<Button>`.

**Add global category detection to the category list rendering.** For each category badge:
- If `cat.user_id === null`: show a lock icon, no edit/delete buttons
- If `cat.user_id !== null`: show edit/delete buttons as today

Replace the category badge rendering section. The exact code depends on the current JSX structure, but the logic is:

```tsx
{categories
  .filter((c) => c.domain === domain)
  .map((cat) => {
    const isGlobal = cat.user_id === null;

    if (editingId === cat.id && !isGlobal) {
      // Edit mode (same as today)
      return (/* existing edit mode JSX */);
    }

    return (
      <Badge key={cat.id} variant={isGlobal ? "secondary" : "outline"} className="gap-1">
        {isGlobal && <Lock className="h-3 w-3" />}
        {!isGlobal && (
          <button onClick={() => handleStartEdit(cat)}>
            <Pencil className="h-3 w-3" />
          </button>
        )}
        {cat.name}
        {!isGlobal && (
          <button onClick={() => handleDelete(cat.id)}>
            <X className="h-3 w-3" />
          </button>
        )}
      </Badge>
    );
  })}
```

Import `Lock` from lucide-react alongside existing icon imports.

- [ ] **Step 3: Verify frontend compiles**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

Expected: no type errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/types/index.ts frontend/src/components/settings/category-manager.tsx
git commit -m "feat: frontend — global category badges with lock icon, remove seed button"
```

---

### Task 7: Backend — Add is_global to JSON Response

**Files:**
- Modify: `backend/internal/model/category.go`

The `IsGlobal()` method exists but isn't serialized to JSON. We need the frontend to receive it. The cleanest approach is a custom `MarshalJSON` or simply a JSON struct tag on a computed field. Since GORM ignores fields with `gorm:"-"`, we can add:

- [ ] **Step 1: Add IsGlobal computed JSON field**

In `backend/internal/model/category.go`, add a method that the JSON encoder will pick up. The simplest approach is to use a response wrapper in the handler, but to keep it DRY, let's use a custom JSON marshal:

Actually, the simplest approach: add the field with `gorm:"-"` and populate it in the repository after fetch.

Add to the Category struct:

```go
IsGlobal  bool           `gorm:"-" json:"is_global"`
```

Then in `backend/internal/repository/category.go`, add a helper and call it after every query that returns categories:

```go
func setGlobalFlag(cats []model.Category) {
	for i := range cats {
		cats[i].IsGlobal = cats[i].UserID == nil
	}
}

func setGlobalFlagSingle(cat *model.Category) {
	cat.IsGlobal = cat.UserID == nil
}
```

Update `ListByDomain`:
```go
func (r *CategoryRepository) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	var cats []model.Category
	err := r.db.
		Where("(user_id IS NULL OR user_id = ?) AND domain = ?", userID, domain).
		Order("sort_order ASC, name ASC").
		Find(&cats).Error
	if err != nil {
		return nil, err
	}
	setGlobalFlag(cats)
	return cats, nil
}
```

Update `GetByID`:
```go
func (r *CategoryRepository) GetByID(userID, id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", id, userID).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	setGlobalFlagSingle(&cat)
	return &cat, nil
}
```

Update `GetOtherCategory`:
```go
func (r *CategoryRepository) GetOtherCategory(domain model.CategoryDomain) (*model.Category, error) {
	var cat model.Category
	err := r.db.
		Where("user_id IS NULL AND domain = ? AND is_system = true AND name = 'Other'", domain).
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	setGlobalFlagSingle(&cat)
	return &cat, nil
}
```

- [ ] **Step 2: Verify full build**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./...
```

Expected: compiles.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/model/category.go backend/internal/repository/category.go
git commit -m "feat: serialize is_global in category JSON responses"
```

---

### Task 8: Integration Test — Verify End-to-End

**Files:** None (manual verification)

- [ ] **Step 1: Start the stack**

```bash
cd /home/folkrom/projects/finance-tracker
docker compose up -d
mise run migrate-up
mise run dev-backend
```

- [ ] **Step 2: Verify global categories exist via API**

Log in at http://localhost:3000/login, then:

```bash
# Get a valid token from browser dev tools, then:
curl -s http://localhost:8080/api/v1/categories?domain=income \
  -H "Authorization: Bearer <token>" | jq '.data[] | {name, is_global, is_system, user_id}'
```

Expected: 7 income categories with `is_global: true`, `user_id: null`.

- [ ] **Step 3: Verify seed endpoint is a no-op**

```bash
curl -s -X POST http://localhost:8080/api/v1/categories/seed \
  -H "Authorization: Bearer <token>" | jq
```

Expected: `{"message": "global categories are available by default"}`

- [ ] **Step 4: Verify creating a user category works**

```bash
curl -s -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Custom", "domain": "income"}' | jq '{name, is_global, user_id}'
```

Expected: `is_global: false`, `user_id: "<your-uuid>"`.

- [ ] **Step 5: Verify global category cannot be modified**

```bash
# Get a global category ID
GLOBAL_ID=$(curl -s http://localhost:8080/api/v1/categories?domain=income \
  -H "Authorization: Bearer <token>" | jq -r '.data[] | select(.name=="Salary") | .id')

curl -s -X PUT http://localhost:8080/api/v1/categories/$GLOBAL_ID \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Hacked"}' | jq
```

Expected: 403 `"cannot modify global categories"`

- [ ] **Step 6: Verify Settings page in browser**

Open http://localhost:3000/2026/settings:
- Global categories should appear with lock icons
- No "Seed Defaults" button
- User can still create/edit/delete their own categories
- Global categories have no edit/delete buttons

- [ ] **Step 7: Final commit (if any adjustments needed)**

```bash
git add -A && git commit -m "fix: adjustments from integration testing"
```

Only if fixes were needed. Skip if everything worked.
