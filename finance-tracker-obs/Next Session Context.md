# Next Session Context

**Last updated:** 2026-04-12
**Branch:** main (pushed to origin)

---

## What Was Completed This Session (2026-04-12)

### Plan 5 — Admin Dashboard & Category Revamp (ALL STEPS DONE)

**Step 1: Category Revamp** (Plan 5a)
- Migration `000009_category_revamp` — nullable `user_id`, `is_system` column, 28 global categories seeded, user categories deduplicated, partial unique indexes
- Category model: `UserID *uuid.UUID` (nullable), `IsSystem bool`, `IsGlobal` computed field
- Repository: `ListByDomain` returns globals + user's own, `Delete` scoped to user only, `GetOtherCategory` for fallback
- Service: `ErrGlobalCategoryReadOnly`, `ErrSystemCategoryProtected` sentinel errors, `SeedDefaults` is now a no-op
- Frontend: removed seed button, global categories show Lock icon + secondary badge, no edit/delete on globals

**Step 2: User Profile** (Plan 5b)
- Migration `000010_create_profiles` — `user_id` unique, `currency` (default MXN), `language` (default en)
- Profile model, repo, service (GetOrCreate, Update with currency/language validation)
- Profile middleware — auto-creates profile on first authenticated request, `sync.Map` cache
- Handler: GET (calls GetOrCreate), PUT (maps validation errors to 400)
- Frontend: ProfileManager component in Settings — currency selector (MXN/USD/EUR/GBP/BRL/COP/ARS), language toggle (en/es)

**Step 3: Admin Middleware & Routes** (Plan 5c)
- Auth middleware now stores full `jwt.MapClaims` in `c.Locals("claims")` + `GetClaims` helper
- Admin middleware: `NewAdminMiddleware()` checks `app_metadata.role == "admin"`, `IsAdmin(c)` helper
- `AdminStats` model (not a DB model, just response struct)
- Category repo: `GetGlobalByID`, `CreateGlobal`, `ReassignAndDelete` (transactional FK reassignment to "Other")
- Admin repo: `GetStats()` — counts for users, profiles, global/user categories
- Category service: `CreateGlobal`, `UpdateGlobal`, `DeleteGlobal`
- Admin service: `GetStats()`
- Admin handler: `CreateCategory`, `UpdateCategory`, `DeleteCategory`, `GetStats`
- Router: `/api/v1/admin/` group with admin middleware, category CRUD + stats
- main.go: wired adminRepo, adminSvc, adminHandler

**Step 4: Admin Dashboard Frontend** (Plan 5d)
- `useAdmin` hook — reads `app_metadata.role` from Supabase session JWT
- `AdminStats` type in `types/admin.ts`
- Admin layout (`/admin`) — separate sidebar (Stats, Categories), header (Back to app + logout), auth gate redirects non-admins
- `/admin/stats` page — 4 stat cards (Users, Profiles, Global Categories, User Categories)
- `/admin/categories` page — grouped by domain, full CRUD via dialogs (create/edit/delete), system categories locked
- Main header: Shield "Admin" link, visible only for admin users

### CRUD Audit (all 9 modules)
- Parallel agent audit of all frontend→backend routes across Categories, Payment Methods, Income, Expenses, Debt, Budget, Dashboard, Cards, Wishlist
- Result: 100% alignment, no gaps found

### Other Changes
- Added git remote `origin` at `git@github.com:Folkrom/finance-tracker.git`
- `.gitignore`: added `backend/server` and `skills/`
- Consolidated all docs from `docs/superpowers/` into `finance-tracker-obs/` (single source of truth)
- Fixed budget test compile errors (`model.CategoryTypeExpense` → `model.CategoryDomainExpense`)
- Fixed `IsGlobal` field/method conflict in Go (renamed method to `IsGlobalCategory()`)

---

## Current Architecture

### Backend (Go/Fiber)
- **Migrations:** `000001` through `000010` (golang-migrate, PostgreSQL)
- **Auth:** Supabase JWKS-based JWT verification via `keyfunc/v3`, `user_id` + `claims` stored in Fiber Locals
- **Middleware chain:** Auth → Profile (auto-create) → route handlers. Admin routes add Admin middleware.
- **Categories:** Hybrid global (`user_id IS NULL`) + user-scoped. `is_system=true` for "Other" per domain (undeletable). Admin can CRUD globals, users can only CRUD their own.
- **Admin routes:** `POST/PUT/DELETE /api/v1/admin/categories`, `GET /api/v1/admin/stats`
- **Profile:** auto-created on first request, currency (MXN default) + language (en default)

### Frontend (Next.js 16)
- **App router:** `[year]/` layout for main app, `/admin/` layout for admin, `/wishlist/` standalone, `/login/` standalone
- **Auth:** Supabase SSR client, `proxy.ts` for session refresh
- **API:** `lib/api.ts` — `apiGet`, `apiPost`, `apiPut`, `apiPatch`, `apiDelete` with auto auth headers
- **UI:** shadcn/ui components, lucide-react icons, sonner toasts
- **i18n:** next-intl (EN/ES)

### Database Tables
`categories`, `payment_methods`, `incomes`, `expenses`, `debts`, `budgets`, `cards`, `wishlist_items`, `profiles`

---

## What Needs Testing (Not Yet Verified in Browser)

These features were implemented but not manually tested in a running browser:

1. **Profile (Plan 5b):**
   - Profile auto-creation on first request
   - Currency/language selectors in Settings
   - Profile persists across sessions

2. **Admin (Plans 5c + 5d):**
   - Set `app_metadata.role = "admin"` in Supabase dashboard for your user
   - Log out and back in to get fresh JWT with admin claim
   - Verify: Shield "Admin" link appears in header
   - Verify: `/admin/stats` shows correct counts
   - Verify: `/admin/categories` shows all 28 global categories grouped by domain
   - Verify: Can create/edit/delete non-system categories via dialogs
   - Verify: System categories ("Other") show lock icon, can't be deleted
   - Verify: Non-admin user can't access `/admin` (redirected)
   - Verify: Non-admin gets 403 from backend admin routes

3. **Category revamp (Plan 5a):**
   - Global categories appear for all users without seeding
   - Users can create their own categories alongside globals
   - Users can't edit/delete global categories (lock icon, no controls)
   - "Seed Defaults" button removed from Settings

---

## What's Next — Potential Plan 6 Ideas

Plan 5 is fully complete. Here are potential next steps (not yet designed or prioritized):

1. **i18n completion** — admin pages use hardcoded English, add to translation files
2. **Admin user management** — `GET /api/v1/admin/users` (list), `GET /api/v1/admin/users/:id` (detail with stats). Requires Supabase Admin API integration.
3. **Audit logging** — track admin actions (category create/update/delete)
4. **Currency formatting** — use profile.currency to format amounts across the app
5. **Dark mode** — shadcn supports it, just needs theme toggle
6. **Data export** — CSV/PDF export of income, expenses, debts per year
7. **Recurring transactions** — auto-create income/expenses on a schedule
8. **Mobile responsiveness** — current layout is desktop-first
9. **Notifications/reminders** — budget threshold alerts
10. **Multi-currency support** — per-transaction currency with conversion

### Open Questions from Plan 5
- Should users be able to "hide" global categories they don't use?
- Do we need admin audit logging from day 1?
- Should currency be per-profile or per-year?

---

## How to Start Dev

```bash
# Terminal 1: Database
docker compose up -d

# Terminal 2: Migrations
cd backend && mise run migrate-up

# Terminal 3: Backend
mise run dev-backend  # localhost:8080

# Terminal 4: Frontend
cd frontend && npm run dev  # localhost:3000
```

## Key File Paths

| What | Path |
|------|------|
| Backend entry | `backend/cmd/server/main.go` |
| Router | `backend/internal/router/router.go` |
| Auth middleware | `backend/internal/middleware/auth.go` |
| Admin middleware | `backend/internal/middleware/admin.go` |
| Profile middleware | `backend/internal/middleware/profile.go` |
| Migrations | `backend/migrations/000001-000010` |
| Frontend app | `frontend/src/app/` |
| Admin pages | `frontend/src/app/admin/` |
| Admin components | `frontend/src/components/admin/` |
| API helpers | `frontend/src/lib/api.ts` |
| Types | `frontend/src/types/index.ts` |
| Supabase client | `frontend/src/lib/supabase/client.ts` |
| All specs & plans | `finance-tracker-obs/` |
