# Plan 5 — Admin Dashboard & Category Revamp

## 1. Category Architecture Revamp

### Current State
- Every category is fully user-scoped (`user_id NOT NULL`)
- Users must call `POST /categories/seed` to get defaults
- No shared/global categories exist

### Target State: Hybrid Global + User Categories

**Global (system) categories** — `user_id IS NULL`
- Seeded once at deploy/migration time (not per-user)
- Visible to all users, cannot be edited or deleted by regular users
- Admin can CRUD global categories via admin dashboard

**User-created categories** — `user_id = <uuid>`
- Only visible to the creating user
- Full CRUD by the owning user
- Can coexist with global categories (no name collision across scopes)

### Migration Changes

```sql
-- 000009_nullable_category_user_id.up.sql
ALTER TABLE categories ALTER COLUMN user_id DROP NOT NULL;

-- Update unique constraint to handle NULLs properly
ALTER TABLE categories DROP CONSTRAINT categories_user_id_name_domain_key;
CREATE UNIQUE INDEX idx_categories_global_unique 
  ON categories (name, domain) WHERE user_id IS NULL;
CREATE UNIQUE INDEX idx_categories_user_unique 
  ON categories (user_id, name, domain) WHERE user_id IS NOT NULL;

-- Seed global categories (run once)
INSERT INTO categories (name, domain, sort_order) VALUES
  -- Income
  ('Salary', 'income', 0),
  ('Bonus', 'income', 1),
  ('Freelance', 'income', 2),
  ('Dividends', 'income', 3),
  ('Interest', 'income', 4),
  ('Side Hustle', 'income', 5),
  -- Expense
  ('Home Expenses', 'expense', 0),
  ('Eating Out', 'expense', 1),
  ('Self Care', 'expense', 2),
  ('Coffee/Drink', 'expense', 3),
  ('Entertainment', 'expense', 4),
  ('Transportation', 'expense', 5),
  ('Groceries', 'expense', 6),
  ('Utilities', 'expense', 7),
  ('Clothes', 'expense', 8),
  ('Card Payments', 'expense', 9),
  ('Savings/Investment', 'expense', 10),
  ('Taxes', 'expense', 11),
  ('Knowledge', 'expense', 12),
  ('Tech', 'expense', 13),
  ('Other', 'expense', 14),
  -- Wishlist
  ('Electronics', 'wishlist', 0),
  ('Clothing', 'wishlist', 1),
  ('Home & Kitchen', 'wishlist', 2),
  ('Books & Media', 'wishlist', 3),
  ('Sports & Outdoors', 'wishlist', 4),
  ('Other', 'wishlist', 5)
ON CONFLICT DO NOTHING;
```

### Backend Changes

**Repository — `ListByDomain`:**
```go
// Before: WHERE user_id = ? AND domain = ?
// After:  WHERE (user_id IS NULL OR user_id = ?) AND domain = ?
```

**Service — `SeedDefaults`:**
- Keep for backward compat but mark as deprecated
- New users get global categories automatically (no seed call needed)

**Handler — CRUD guards:**
- `Update`/`Delete`: reject if `category.UserID == nil` (global) and user is not admin
- `Create`: always sets `user_id` to requesting user (users can't create global categories via this endpoint)

### Frontend Changes

- Remove "Seed Defaults" button from Settings (global categories exist from day 1)
- Show global categories with a lock/system icon (non-editable)
- Show user categories with edit/delete controls
- "Add Category" button creates user-scoped category

---

## 2. Admin Dashboard

### Authentication & Authorization

**Admin role detection:**
- Option A: Supabase custom claims (`app_metadata.role = "admin"`)
- Option B: `is_admin` column in a `profiles` table
- **Recommended: Option A** — no extra table, Supabase handles it, JWT `app_metadata` is already in the token

**Middleware:**
```go
func NewAdminMiddleware() fiber.Handler {
    // Extract app_metadata.role from JWT claims
    // Return 403 if role != "admin"
}
```

### Admin Routes

```
GET    /api/v1/admin/stats              — user count, category count, etc.
GET    /api/v1/admin/users              — list users (from Supabase Admin API)
GET    /api/v1/admin/users/:id          — user detail + their data stats

POST   /api/v1/admin/categories         — create global category
PUT    /api/v1/admin/categories/:id     — update global category
DELETE /api/v1/admin/categories/:id     — delete global category
GET    /api/v1/admin/categories         — list all global categories
```

### Admin Pages (Frontend)

| Page | Path | Purpose |
|------|------|---------|
| Admin Home | `/admin` | Stats overview (user count, active users, data volume) |
| Users | `/admin/users` | User list with search, last-active date |
| User Detail | `/admin/users/:id` | View user's data stats (not their actual data) |
| Global Categories | `/admin/categories` | CRUD for system-wide categories |

### Admin Layout

- Separate sidebar from main app (admin-specific nav)
- Protected by admin middleware on both frontend and backend
- Accessible via profile menu → "Admin Panel" (only visible to admins)

---

## 3. User Profile

### Profile Data (Supabase `auth.users` + optional `profiles` table)

| Field | Source | Editable |
|-------|--------|----------|
| Email | Supabase auth | No (via Supabase settings) |
| Display name | `profiles.display_name` | Yes |
| Avatar URL | `profiles.avatar_url` | Yes |
| Preferred currency | `profiles.currency` | Yes |
| Language | `profiles.language` | Yes (EN/ES) |
| Timezone | `profiles.timezone` | Yes |
| Created at | Supabase auth | No |

### Profile Routes

```
GET  /api/v1/profile          — get current user's profile
PUT  /api/v1/profile          — update profile fields
POST /api/v1/profile/avatar   — upload avatar (store in Supabase Storage)
```

### Profile Migration

```sql
-- 000010_create_profiles.up.sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    display_name VARCHAR(100),
    avatar_url TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    language VARCHAR(5) NOT NULL DEFAULT 'en',
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profiles_user_id ON profiles (user_id);
```

### Profile Page (Frontend)

- Path: `/settings/profile` (or tab within Settings)
- Avatar upload with preview
- Currency selector (affects display formatting across the app)
- Language toggle (already have i18n wiring for EN/ES)

---

## 4. Implementation Order

1. **Category revamp** — migration, repo/service/handler changes, frontend updates
2. **Profile** — migration, CRUD, settings page tab
3. **Admin middleware** — JWT claim check
4. **Admin categories** — global category CRUD via admin routes
5. **Admin dashboard** — stats page, user list

---

## 5. Open Questions

- Should users be able to "hide" global categories they don't use? (vs cluttering their lists)
- Should deleting a global category cascade or soft-fail if users have data referencing it?
- Do we need admin audit logging from day 1?
- Should currency be per-profile or per-year (some users move countries)?
