# Category Revamp: Global + User Categories

**Date:** 2026-04-12
**Status:** Approved
**Scope:** Category architecture only (admin dashboard, profiles, admin CRUD deferred to later Plan 5 phases)

---

## Problem

Every category is user-scoped. New users see nothing until they manually call "Seed Defaults." No shared baseline exists across users, and there's no way to maintain a standard set of categories that all users benefit from.

## Design

### Data Model Changes

**categories table:**

| Column | Type | Change |
|--------|------|--------|
| user_id | UUID | Now NULLABLE. NULL = global/system category |
| is_system | BOOLEAN NOT NULL DEFAULT false | NEW. True for protected categories (e.g., "Other") that cannot be deleted |

**Constraints:**
- Drop existing `UNIQUE(user_id, name, domain)`
- Add partial unique index: `(name, domain) WHERE user_id IS NULL` — prevents duplicate globals
- Add partial unique index: `(user_id, name, domain) WHERE user_id IS NOT NULL` — prevents duplicate user categories

### Global Categories (seeded in migration)

**Income (7):** Salary, Bonus, Freelance, Dividends, Interest, Side Hustle, **Other** (is_system=true)
**Expense (15):** Home Expenses, Eating Out, Self Care, Coffee/Drink, Entertainment, Transportation, Groceries, Utilities, Clothes, Card Payments, Savings/Investment, Taxes, Knowledge, Tech, **Other** (is_system=true)
**Wishlist (6):** Electronics, Clothing, Home & Kitchen, Books & Media, Sports & Outdoors, **Other** (is_system=true)

"Other" in each domain has `is_system = true` and cannot be deleted.

### Migration Strategy (000009)

Single migration, safe for dev (pre-launch):

1. `ALTER TABLE categories ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT false`
2. `ALTER TABLE categories ALTER COLUMN user_id DROP NOT NULL`
3. Insert global categories (`user_id = NULL`), mark "Other" per domain as `is_system = true`
4. For each existing user category matching a global `(name, domain)`:
   - Update FK references in `incomes`, `expenses`, `debts`, `budgets`, `wishlist_items` to point to the global category ID
   - Delete the user-owned duplicate
5. `ALTER TABLE categories DROP CONSTRAINT categories_user_id_name_domain_key`
6. `CREATE UNIQUE INDEX idx_categories_global_unique ON categories (name, domain) WHERE user_id IS NULL`
7. `CREATE UNIQUE INDEX idx_categories_user_unique ON categories (user_id, name, domain) WHERE user_id IS NOT NULL`

**Down migration:** Reverse — drop indexes, restore constraint, delete globals, make user_id NOT NULL, drop is_system.

### API Behavior

**`GET /categories?domain=X`**
- Query: `WHERE (user_id IS NULL OR user_id = ?) AND domain = ?`
- Response includes `is_system` and `is_global` (derived: `user_id IS NULL`) fields so frontend can render lock icons

**`POST /categories`**
- Unchanged: always sets `user_id` from JWT. Users create user-scoped categories.

**`PUT /categories/:id`**
- If `category.UserID == nil`: return 403 "cannot modify global categories"
- Otherwise: unchanged behavior

**`DELETE /categories/:id`**
- If `category.IsSystem`: return 403 "cannot delete system categories"
- If `category.UserID == nil`: return 403 "cannot delete global categories"
- Otherwise: unchanged behavior
- Admin delete of global categories: reassign all references to the domain's "Other" category, then delete. (Implemented when admin routes ship.)

**`POST /categories/seed`**
- Deprecated. Returns 200 with `{"message": "global categories are available by default"}`. No-op.

### Model Changes

```go
type Category struct {
    Base
    Name     string         `json:"name" gorm:"type:varchar(255);not null"`
    Domain   CategoryDomain `json:"domain" gorm:"type:varchar(50);not null"`
    Color    *string        `json:"color,omitempty" gorm:"type:varchar(7)"`
    SortOrder int           `json:"sort_order" gorm:"default:0"`
    IsSystem  bool          `json:"is_system" gorm:"not null;default:false"`
}
```

Note: `Base.UserID` changes from `uuid.UUID` to `*uuid.UUID` (pointer, nullable).

### Repository Changes

**`ListByDomain(userID uuid.UUID, domain CategoryDomain)`**
```sql
WHERE (user_id IS NULL OR user_id = ?) AND domain = ?
ORDER BY sort_order ASC, name ASC
```

**`GetByID(userID uuid.UUID, id uuid.UUID)`**
```sql
WHERE id = ? AND (user_id IS NULL OR user_id = ?)
```

This allows users to read global categories (needed for FK lookups) but write operations still check ownership.

### Handler Changes

- `Update`: fetch category, reject if `UserID == nil`
- `Delete`: fetch category, reject if `IsSystem` or `UserID == nil`
- `SeedDefaults`: return 200 no-op message

### Frontend Changes

- **Settings page:** Remove "Seed Defaults" button
- **Category lists:** Global categories show lock icon, no edit/delete buttons
- **User categories:** Edit/delete controls as today
- **"Add Category" button:** Creates user-scoped category (unchanged)
- **Category response shape:** Add `is_system` and `is_global` boolean fields for UI rendering

### Base Model Impact

`Base.UserID` becoming `*uuid.UUID` affects all models embedding `Base`. The migration only changes the categories table, but the Go struct change is shared. All other modules (income, expense, debt, etc.) continue to set `UserID` as non-nil — the pointer change is backward compatible since `&userID` works everywhere `userID` was used.

If the pointer change creates too much churn, alternative: don't change Base. Instead, override `UserID` directly on Category with a `*uuid.UUID` field (GORM respects field overrides). This isolates the change to one model.
