# Next Session Context

## Goal: Full CRUD & Front↔Back Integration Audit

Verify every module's CRUD operations work end-to-end: frontend forms submit correctly, backend routes respond, data persists, and the UI updates. This doc maps every frontend API call to its backend route so agents can systematically check each one.

---

## Prerequisites

Before starting verification:

1. `docker compose up -d` — Postgres on port 5466
2. `mise run migrate-up` — all 8 migrations applied
3. `mise run dev-backend` — Go server on :8080
4. `mise run dev-frontend` — Next.js on :3000
5. Log in at http://localhost:3000/login (creates Supabase session cookies)
6. Seed categories: Settings page → click "Seed Defaults" for both expense and income domains

---

## Module Verification Checklist

### 1. Categories (Settings page)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Seed defaults | `POST /api/v1/categories/seed` | `POST /categories/seed` |
| Create | `POST /api/v1/categories` | `POST /categories/` |
| List (income) | `GET /api/v1/categories?domain=income` | `GET /categories/` |
| List (expense) | `GET /api/v1/categories?domain=expense` | `GET /categories/` |
| List (wishlist) | `GET /api/v1/categories?domain=wishlist` | `GET /categories/` |
| Update | `PUT /api/v1/categories/:id` | `PUT /categories/:id` |
| Delete | `DELETE /api/v1/categories/:id` | `DELETE /categories/:id` |

**Verify:** domain query param filtering works, seeding creates both income and expense defaults, wishlist categories are separate.

### 2. Payment Methods (Settings page)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/payment-methods` | `POST /payment-methods/` |
| List | `GET /api/v1/payment-methods` | `GET /payment-methods/` |
| Get by ID | `GET /api/v1/payment-methods/:id` | `GET /payment-methods/:id` |
| Update | `PUT /api/v1/payment-methods/:id` | `PUT /payment-methods/:id` |
| Delete | `DELETE /api/v1/payment-methods/:id` | `DELETE /payment-methods/:id` |

**Verify:** type field (Cash, Debit Card, Credit Card, Digital Wallet, Crypto) persists correctly, payment methods appear in Expense and Debt dropdowns.

### 3. Income (year-scoped)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/years/:year/incomes` | `POST /years/:year/incomes` |
| List by year | `GET /api/v1/years/:year/incomes` | `GET /years/:year/incomes` |
| Update | `PUT /api/v1/years/:year/incomes/:id` | `PUT /years/:year/incomes/:id` |
| Delete | `DELETE /api/v1/years/:year/incomes/:id` | `DELETE /years/:year/incomes/:id` |

**Verify:** year scoping filters correctly (income from 2025 shouldn't appear in 2026), category association works, amount/currency persist.

### 4. Expenses (year-scoped)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/years/:year/expenses` | `POST /years/:year/expenses` |
| List by year | `GET /api/v1/years/:year/expenses` | `GET /years/:year/expenses` |
| Update | `PUT /api/v1/years/:year/expenses/:id` | `PUT /years/:year/expenses/:id` |
| Delete | `DELETE /api/v1/years/:year/expenses/:id` | `DELETE /years/:year/expenses/:id` |

**Verify:** payment method relation, category from shared list, type field (Expense/Saving/Investment), date→month derivation.

### 5. Debt (year-scoped)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/years/:year/debts` | `POST /years/:year/debts` |
| List by year | `GET /api/v1/years/:year/debts` | `GET /years/:year/debts` |
| Update | `PUT /api/v1/years/:year/debts/:id` | `PUT /years/:year/debts/:id` |
| Delete | `DELETE /api/v1/years/:year/debts/:id` | `DELETE /years/:year/debts/:id` |

**Verify:** payment method relation (creditor), category from shared list, affects Card usage calculation.

### 6. Budget

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/budgets` | `POST /budgets/` |
| Get summary | `GET /api/v1/budgets?year=:year&month=:month` | `GET /budgets/` |
| List recurring | `GET /api/v1/budgets/recurring` | `GET /budgets/recurring` |
| Update | `PUT /api/v1/budgets/:id` | `PUT /budgets/:id` |
| Delete | `DELETE /api/v1/budgets/:id` | `DELETE /budgets/:id` |

**Verify:** budget remaining = limit - expenses - debts for that category/month, recurring defaults work, per-month overrides don't affect template.

### 7. Dashboard

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Get data | `GET /api/v1/years/:year/dashboard` | `GET /years/:year/dashboard` |

**Verify:** all 6 charts render with data (Net Savings line, Income donut, Expenses donut, Income vs Expenses donut, Daily Expenses bar, Daily Debt bar). Requires seeded data across multiple months to see meaningful charts.

### 8. Cards

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/cards` | `POST /cards/` |
| List summaries | `GET /api/v1/cards?year=:year&month=:month` | `GET /cards/` |
| Get by ID | `GET /api/v1/cards/:id` | `GET /cards/:id` |
| Update | `PUT /api/v1/cards/:id` | `PUT /cards/:id` |
| Delete | `DELETE /api/v1/cards/:id` | `DELETE /cards/:id` |

**Verify:** health indicator colors match thresholds (green 0-20%, yellow 21-30%, orange 31-70%, red 71+%), usage auto-calculates from Debt entries for linked payment method, manual override works.

### 9. Wishlist (non-year-scoped)

| Operation | Frontend call | Backend route |
|-----------|--------------|---------------|
| Create | `POST /api/v1/wishlist` | `POST /wishlist/` |
| List | `GET /api/v1/wishlist` | `GET /wishlist/` |
| Update | `PUT /api/v1/wishlist/:id` | `PUT /wishlist/:id` |
| Update status | `PATCH /api/v1/wishlist/:id/status` | `PATCH /wishlist/:id/status` |
| Delete | `DELETE /api/v1/wishlist/:id` | `DELETE /wishlist/:id` |

**Verify:** three views render (gallery, table, board), status change via Kanban dropdown calls PATCH, links array (max 5), priority sorting, category from wishlist domain (not expense).

---

## Known Issues to Fix

### Pre-existing test compile errors
`backend/internal/service/budget_test.go` references `model.CategoryTypeExpense` which does not exist on the current `Category` model. Dead test code from an earlier iteration. `go test ./...` fails. Either delete the test file or update it to match the current model.

### Cross-cutting concerns to check during audit

- **CORS:** frontend on :3000 → backend on :8080. Currently `AllowOrigins: "http://localhost:3000"` in main.go. Verify no CORS errors in browser console.
- **Auth cookie refresh:** proxy.ts wires `updateSession` on every request. Verify session survives page reloads and multi-tab usage.
- **Error handling:** API errors should surface as toast messages, not unhandled rejections. Check browser console during each CRUD operation.
- **Multi-tenant isolation:** every query scopes by `user_id` from JWT `sub` claim. If testing with multiple users, verify data doesn't leak between accounts.
- **i18n:** toggle between EN/ES and verify labels on forms, toasts, and empty states.

---

## What shipped in the last session (2026-04-11)

- `frontend/src/proxy.ts` — Next.js 16 proxy wiring for Supabase session refresh (fixes auth bug)
- JWKS-based JWT validation — backend fetches public keys from Supabase instead of shared secret
- Year-scoped routes fix — income/expenses/debts get/put/delete now match frontend paths
- `go mod tidy` — `lib/pq` moved to direct dependency, `keyfunc/v3` added
- Removed `SUPABASE_JWT_SECRET` from config/env (no longer needed)
