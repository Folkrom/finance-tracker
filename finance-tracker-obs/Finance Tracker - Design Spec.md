# Finance Tracker - Design Spec

> **Status:** Reviewed — pending final user sign-off
> **Date:** 2026-04-03
> **Approach:** C — Go API + Next.js SSR Frontend

---

## Overview

Personal finance tracker replacing the current Notion-based system ("Personal Finance Tracker - Plantilla"). Goal is to build an extensible, multi-tenant product that can be shared with friends/family and potentially commercialized.

---

## Architecture

```
┌─────────────────────┐     ┌─────────────────────┐
│   Next.js Frontend  │────▶│   Go/Fiber REST API  │
│  (SSR + CSR)        │◀────│  (GORM, Wire, Zap)   │
└─────────────────────┘     └──────────┬──────────┘
                                       │
                            ┌──────────▼──────────┐
                            │  PostgreSQL/Supabase │
                            │  + Supabase Auth     │
                            └─────────────────────┘
```

### Tech Stack
- **Backend:** Golang + Fiber, GORM, Wire (DI), Zap (logging)
- **Frontend:** Next.js + TypeScript
- **Database:** PostgreSQL via Supabase
- **Auth:** Supabase Auth (JWT) + Fiber middleware for authorization
- **Dev tooling:** mise-en-place
- **Multi-tenancy:** Row-level `user_id` on all tables + Supabase RLS as safety net
- **i18n:** English (default) + Spanish, with language toggle. Extensible for more languages.
- **Currency:** Configurable per-field (USD / MXN for now, extensible)

---

## Key Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Cuentas por Cobrar | Dropped | Not needed in new app |
| PPR Teresa | Dropped | Not needed in new app |
| Categories | Single shared list across Expenses & Debt | Budget math requires exact matching |
| Wishlist Categories | Separate list from Expenses/Debt | Different domain (product types vs spending types) |
| Time scoping | Year-scoped workspaces, cross-year queryable | Matches existing Notion workflow |
| Language | i18n (English + Spanish) with toggle | User mixes both; commercialization needs it |
| Payment methods | User-defined + typed (Cash, Debit Card, Credit Card, Digital Wallet, Crypto) | Flexible + enables auto-linking to Cards |
| "Daily cash expenses" chart | All expenses by day, not just cash method | "Cash" = cash-flow |
| Card usage calculation | Auto from Debt + manual override | Flexibility for untracked charges |
| Card health thresholds | Green 0-20%, Yellow 21-30%, Orange 31-70%, Red 71-100%+ | User-defined based on personal finance goals |
| Auth & multi-user | Multi-tenant from day one | Future commercialization |
| Data migration from Notion | Later concern | Build first, migrate later |
| Budget recurrence | Recurring defaults + per-month overrides | Set once, tweak as needed |
| "Debt" category tag | Renamed to "Card Payments" | Avoid naming conflict with Debt module |
| "Alcohol/Drugs" category | Split into "Alcohol" and "Drugs" | Separate tracking |

---

## Modules

Each module = a left-menu navigation item in the UI.

### 1. Income
Where salary, bonus, and extra income are registered.

| Property | Type | Notes |
|---|---|---|
| Source | string | Title/name of the income entry |
| Amount | decimal | Currency configurable (MXN default) |
| Tag/Category | select | Salary, Bonus, Freelance, Dividends, Interest, Side Hustle |
| Date | date | When it was received |
| Month | derived | Auto-derived from Date, links to year-workspace |

### 2. Expenses/Outs
Where cash, bank account, and liquid money expenses are registered.

| Property | Type | Notes |
|---|---|---|
| Name | string | Title/description of the expense |
| Amount | decimal | Currency configurable |
| Date | date | |
| Month | derived | Auto-derived from Date |
| Payment Method | relation | Links to user-defined payment methods (typed) |
| Tags/Categories | select | **Shared category list** with Debt |
| Type | select | Expense / Saving / Investment |

### 3. Debt
Where credit card balance and other debt are registered (credit card charges, loans, etc.).

| Property | Type | Notes |
|---|---|---|
| Name | string | Title/description |
| Amount | decimal | Currency configurable |
| Date | date | |
| Month | derived | Auto-derived from Date |
| Payment Method/Creditor | relation | Links to user-defined payment methods |
| Tags/Categories | select | **Same shared list** as Expenses |

### Shared Category List (Expenses + Debt)
Home Expenses, Eating Out, Self Care, Coffee/Drink, Entertainment, Transportation, Groceries, Utilities, Clothes, Other, Card Payments, Savings/Investment, Alcohol, Drugs, Taxes, Knowledge, Tech

> User-defined: users can add/edit/remove categories.

### 4. Budget
Monthly spending limits applied to the combined total of Expenses + Debt per category.

| Property | Type | Notes |
|---|---|---|
| Category | relation | Links to shared category list |
| Monthly Limit | decimal | Max spend for this category per month |
| Period | month/year | Which month the budget applies to |

**Recurrence model:** Users define default budgets (recurring template). Each month auto-inherits defaults. Individual months can be overridden without affecting the template.

**Budget calculation:** `remaining = limit - (sum of Expenses in category for month) - (sum of Debt in category for month)`

> Example: $100 budget for Coffee in May. $40 expense (cash) + $50 debt (credit card) = $10 remaining.

### 5. Dashboards (Home Page)

All charts are configurable to show yearly or monthly data.

#### 5a. Net Savings Over Year — Line Chart
- Cumulative line chart showing monthly net (income - expenses) over the year
- Smooth line, data labels, sorted chronologically
- Replicated from Notion's "2026 Savings" chart

#### 5b. Income Breakdown — Donut Chart
- Income grouped by tag (Salary, Bonus, Freelance, etc.)
- Aggregated by sum of amount
- Replicated from Notion's "Income Breakdown" chart

#### 5c. Expenses Breakdown — Donut Chart
- Expenses grouped by Type (Expense / Saving / Investment)
- Aggregated by sum of amount
- Replicated from Notion's "Expenses Breakdown" chart

#### 5d. Yearly Income vs Expenses — Donut Chart
- Total yearly income vs total yearly expenses (excludes Debt)
- Scoped to current year-workspace
- New chart (from draft)

#### 5e. Daily Expenses — Vertical Bar Chart
- All expenses summed by day, sorted chronologically
- Clicking a bar previews the list of expenses for that day
- Clicking an individual expense navigates to it
- Replicated from Notion's "Cash" chart + draft interaction behavior

#### 5f. Daily Debt (TDCs) — Vertical Bar Chart
- Debt entries summed by day, sorted chronologically
- Filterable by month
- Same click-to-preview/navigate interaction as 5e
- Replicated from Notion's "TDCs" chart + draft interaction behavior

### 6. Cards (Optional but desired)

| Property | Type | Notes |
|---|---|---|
| Bank | string | Issuing bank |
| Card Limit | decimal | Credit limit |
| Recommended Max Usage | computed | 30% of limit (configurable) |
| % Usage This Month | computed | Auto-calculated from Debt entries for this card + manual override |
| Health Indicator | computed | Based on usage % (see thresholds below) |
| Level | string | Optional (e.g., Gold, Platinum) |

**Auto-calculation:** Sum all Debt entries where Payment Method matches this card in the current month, divided by card limit. Manual override available for untracked charges.

**Health Indicator Thresholds:**

| Color | Usage % | Meaning |
|---|---|---|
| Green | 0–20% | Healthy |
| Yellow | 21–30% | Recommended max |
| Orange | 31–70% | Warning |
| Red | 71–100%+ | Danger |

**Auto-linking:** Credit Card type payment methods automatically associate with a Card entry.

### 7. Wishlist

| Property | Type | Notes |
|---|---|---|
| Item/Name | string | Title |
| Image | file | Upload with URL fallback |
| Price | decimal | Currency configurable (MXN/USD) |
| Links to Buy | up to 5 URLs | Purchase options |
| Category | select | Separate list (see below) |
| Priority | select | Low / Medium / High |
| Status | status (kanban) | Grouped (see below) |
| Target Purchase Date | date | |
| Monthly Contribution | decimal | Currency configurable |

**Wishlist Categories (separate from Expenses/Debt):**
Electronics, Clothing, Home & Kitchen, Books & Media, Sports & Outdoors, Beauty & Personal Care, Toys & Games, Other

**Status Groups (Kanban):**
- **To-do:** Interested
- **In Progress:** Saving For, Waiting for Sale, Ordered
- **Complete:** Purchased, Received, Cancelled

**Views (replicated from Notion):**
1. **Gallery View** — Large cards with Image as cover. Shows Item, Price, Category, Priority, Status. Filtered to active items (To-do + In Progress). Sorted by Priority.
2. **Main Table** — All properties visible. Sorted by Status, then Priority.
3. **Board (By Category)** — Kanban-style columns grouped by Category.

---

## Shared Entities

### Payment Methods (User-defined + Typed)

| Property | Type | Notes |
|---|---|---|
| Name | string | e.g., "BBVA Debit", "Nu Credit" |
| Type | select | Cash / Debit Card / Credit Card / Digital Wallet / Crypto |
| Details | string | Optional (last 4 digits, account name, etc.) |

Credit Card type payment methods auto-link to the Cards module.

### Categories — Expenses/Debt (Shared)
Single list used by both Expenses and Debt. User-defined, seeded with defaults on first setup.

### Categories — Wishlist (Separate)
Separate list for Wishlist items. User-defined, seeded with defaults on first setup.

### Year Workspaces
Logical grouping by year. All data belongs to a year. Cross-year queries supported via API filters.

---

## Multi-Tenancy & Auth

- Supabase Auth handles registration, login (email/password, social, magic link)
- JWT tokens validated by Fiber middleware
- All database tables have `user_id` column
- Go middleware enforces `user_id` scoping on every query
- Supabase RLS as defense-in-depth

---

## Existing Notion Data (Reference)

### Payment Methods in Use
**Debit/Cash:** Cash, Cuenta MercadoPago, Debit Card/BBVA, Debit Card/Nu
**Credit/Other:** Credit Card/Nu, Credit Card/Mercado Pago, Credit Card/AMEX, Credit Card/Plata, Bitso

### Income Tags
Salary, Bonus, Freelance, Dividends, Interest, Side Hustle

### Notion Chart Configurations (Reference)
| Chart | Type | Data Source | Grouping | Aggregation |
|---|---|---|---|---|
| 2026 Savings | Cumulative Line | Total Savings | By month (Name) | Monthly Net formula |
| Income Breakdown | Donut | Income | By Tags | Sum of Amount |
| Expenses Breakdown | Donut | Expenses | By Type | Sum of Amount |
| Cash (Daily Expenses) | Column | Expenses | By Date (day) | Sum of Amount |
| TDCs (Daily Debt) | Column | Debt | By Date (day) | Sum of Amount |
