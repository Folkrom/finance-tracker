# Finance Tracker

A personal finance management web app built to replace a Notion-based system. Multi-tenant, bilingual (EN/ES), year-scoped workspaces, with 7 modules covering the full personal finance lifecycle.

## Architecture 

```
frontend/          Next.js 16 (App Router, TypeScript, Tailwind, shadcn/ui)
backend/           Go REST API (Fiber v2, GORM, Zap)
docker-compose.yml PostgreSQL 16 (local dev)
mise.toml          Task runner & tool versions
```

**Auth:** Supabase Auth issues JWTs. The frontend uses `@supabase/ssr` for session management. The backend validates JWTs with the shared secret and extracts `user_id` from the `sub` claim. Every database row is scoped to `user_id` (multi-tenant).

**Database:** PostgreSQL via Supabase (production) or Docker (local dev). Migrations managed with `golang-migrate`.

**i18n:** English (default) + Spanish, toggled via `next-intl`.

**Currency:** Configurable per field (MXN/USD), extensible.

## Modules

| Module | Route | Description |
|---|---|---|
| Dashboard | `/:year/dashboard` | 6 Recharts charts (net savings, income/expense breakdowns, income vs expenses, daily expenses, daily debts) |
| Income | `/:year/income` | Income records with source, amount, category, date |
| Expenses | `/:year/expenses` | Expenses with type (expense/saving/investment), payment method, category |
| Debt | `/:year/debt` | Debt records (credit card charges, loans) with payment method and category |
| Budget | `/:year/budget` | Per-category monthly limits with recurring defaults and overrides. Spent = Expenses + Debt |
| Cards | `/:year/cards` | Credit card health monitoring (auto-calculated usage from debts, color-coded thresholds) |
| Wishlist | `/wishlist` | Not year-scoped. 3 views (Gallery, Table, Board/Kanban). 7 statuses, priority levels, purchase links |

## Prerequisites

- [mise](https://mise.jdx.dev/) (manages Go 1.24+ and Node 20)
- Docker & Docker Compose (for local PostgreSQL)
- A [Supabase](https://supabase.com/) project (free tier works) for authentication

## Setup

### 1. Trust mise config and install tools

```bash
cd finance-tracker
mise trust
mise install
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

This starts PostgreSQL 16 on `localhost:5466` (mapped from container's `5432` to avoid conflicts with other local postgres instances) with:
- User: `finance`
- Password: `finance_dev`
- Database: `finance_tracker`

### 3. Configure environment variables

**Backend** - create `backend/.env`:

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env`:

```env
DATABASE_URL=postgres://finance:finance_dev@localhost:5466/finance_tracker?sslmode=disable
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_JWT_SECRET=your-jwt-secret   # Settings > API > JWT Secret in Supabase dashboard
BACKEND_PORT=8080
ENVIRONMENT=development
```

> **Note:** `mise.toml` has `[env] _.file = "backend/.env"` so mise tasks automatically load these variables. You don't need to `source` anything manually.

**Frontend** - create `frontend/.env.local`:

```bash
cp frontend/.env.local.example frontend/.env.local
```

Edit `frontend/.env.local`:

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key   # Settings > API > anon/public key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 4. Run database migrations

```bash
mise run migrate-up
```

This runs all 8 migrations in order:

| # | Migration |
|---|---|
| 1 | `create_categories` |
| 2 | `create_payment_methods` |
| 3 | `create_incomes` |
| 4 | `create_expenses` |
| 5 | `create_debts` |
| 6 | `create_budgets` |
| 7 | `create_cards` |
| 8 | `create_wishlist_items` |

To rollback the last migration: `mise run migrate-down`

### 5. Install frontend dependencies

```bash
cd frontend && npm install
```

### 6. Start dev servers

In two terminals (or use tmux/split pane):

```bash
# Terminal 1 - Backend (localhost:8080)
mise run dev-backend

# Terminal 2 - Frontend (localhost:3000)
mise run dev-frontend
```

Open http://localhost:3000. You'll land on the login page.

### 7. First-time setup after login

After signing up/logging in, go to **Settings** and click **"Load Default Categories"** to seed the default category lists for Income, Expenses/Debt, and Wishlist.

## mise Tasks

| Task | Command | Description |
|---|---|---|
| `dev-backend` | `mise run dev-backend` | Start Go backend on `:8080` |
| `dev-frontend` | `mise run dev-frontend` | Start Next.js on `:3000` |
| `migrate-up` | `mise run migrate-up` | Apply all pending migrations |
| `migrate-down` | `mise run migrate-down` | Rollback last migration |
| `test-backend` | `mise run test-backend` | Run Go tests |
| `test-frontend` | `mise run test-frontend` | Run frontend tests |

## API Reference

All endpoints under `/api/v1/` require `Authorization: Bearer <supabase-jwt>`.

### Categories

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/categories` | Create category (body: `name`, `domain`, `color?`) |
| GET | `/api/v1/categories?domain=income\|expense\|wishlist` | List by domain |
| POST | `/api/v1/categories/seed` | Seed default categories |
| GET | `/api/v1/categories/:id` | Get by ID |
| PUT | `/api/v1/categories/:id` | Update |
| DELETE | `/api/v1/categories/:id` | Delete |

### Payment Methods

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/payment-methods` | Create (body: `name`, `type`: cash/debit_card/credit_card/digital_wallet/crypto) |
| GET | `/api/v1/payment-methods` | List all |
| GET | `/api/v1/payment-methods/:id` | Get by ID |
| PUT | `/api/v1/payment-methods/:id` | Update |
| DELETE | `/api/v1/payment-methods/:id` | Delete |

### Income (year-scoped)

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/years/:year/incomes` | Create |
| GET | `/api/v1/years/:year/incomes` | List by year |
| GET | `/api/v1/incomes/:id` | Get by ID |
| PUT | `/api/v1/incomes/:id` | Update |
| DELETE | `/api/v1/incomes/:id` | Delete |

### Expenses (year-scoped)

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/years/:year/expenses` | Create (body includes `type`: expense/saving/investment) |
| GET | `/api/v1/years/:year/expenses` | List by year |
| GET | `/api/v1/expenses/:id` | Get by ID |
| PUT | `/api/v1/expenses/:id` | Update |
| DELETE | `/api/v1/expenses/:id` | Delete |

### Debts (year-scoped)

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/years/:year/debts` | Create |
| GET | `/api/v1/years/:year/debts` | List by year |
| GET | `/api/v1/debts/:id` | Get by ID |
| PUT | `/api/v1/debts/:id` | Update |
| DELETE | `/api/v1/debts/:id` | Delete |

### Budgets

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/budgets` | Create (body: `category_id`, `monthly_limit`, `month`, `year`, `is_recurring`) |
| GET | `/api/v1/budgets?month=&year=` | Get summary (budget lines with spent/remaining) |
| GET | `/api/v1/budgets/recurring` | List recurring defaults |
| PUT | `/api/v1/budgets/:id` | Update |
| DELETE | `/api/v1/budgets/:id` | Delete |

### Dashboard

| Method | Path | Description |
|---|---|---|
| GET | `/api/v1/years/:year/dashboard?month=` | Aggregated dashboard data (9 parallel queries) |

### Cards

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/cards` | Create (requires credit_card payment method) |
| GET | `/api/v1/cards?month=&year=` | Get summaries with health indicators |
| GET | `/api/v1/cards/:id` | Get by ID |
| PUT | `/api/v1/cards/:id` | Update |
| DELETE | `/api/v1/cards/:id` | Delete |

Card health thresholds: Green 0-20%, Yellow 21-30%, Orange 31-70%, Red 71-100%+

### Wishlist

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/wishlist` | Create item |
| GET | `/api/v1/wishlist?status=interested,saving_for` | List (optional comma-separated status filter) |
| GET | `/api/v1/wishlist/:id` | Get by ID |
| PUT | `/api/v1/wishlist/:id` | Update |
| PATCH | `/api/v1/wishlist/:id/status` | Update status only (for Kanban drag-and-drop) |
| DELETE | `/api/v1/wishlist/:id` | Delete |

Wishlist statuses grouped for Kanban:
- **To-Do:** interested
- **In Progress:** saving_for, waiting_for_sale, ordered
- **Complete:** purchased, received, cancelled

## Project Structure

```
backend/
  cmd/server/main.go        Entry point, DI wiring
  internal/
    config/                  Env var loading
    database/                GORM connection setup
    middleware/               JWT auth middleware
    model/                   GORM models (Category, PaymentMethod, Income, Expense, Debt, Budget, Card, WishlistItem)
    repository/              Database access layer
    service/                 Business logic layer
    handler/                 HTTP handlers (Fiber)
    router/                  Route registration
  migrations/                SQL migration files (000001-000008)

frontend/
  src/
    app/
      page.tsx               Redirect to /:year/dashboard
      login/page.tsx         Supabase auth
      [year]/                Year-scoped layout + modules
        dashboard/
        income/
        expenses/
        debt/
        budget/
        cards/
        settings/
      wishlist/              Non-year-scoped (own layout)
    components/
      ui/                    shadcn/ui components (uses @base-ui/react, NOT Radix)
      layout/                Sidebar, Header, YearSwitcher
      income/                Income form + table
      expenses/              Expense form + table
      debt/                  Debt form + table
      budget/                Budget form + table
      dashboard/             6 Recharts chart components
      cards/                 Card form + health card
      wishlist/              Form, Gallery, Table, Board views
    lib/
      api.ts                 API client (apiGet, apiPost, apiPut, apiPatch, apiDelete)
      i18n/                  en.json, es.json translations
      supabase/              Supabase client setup
    types/index.ts           All TypeScript interfaces
```

## Tech Stack

### Backend
- **Go 1.22+** with Fiber v2 (HTTP framework)
- **GORM** (ORM) with PostgreSQL driver
- **golang-migrate** for schema migrations
- **shopspring/decimal** for money arithmetic
- **golang-jwt/v5** for JWT validation
- **Zap** for structured logging
- **lib/pq** for PostgreSQL array types

### Frontend
- **Next.js 16** (App Router, React 19)
- **TypeScript**
- **Tailwind CSS v4** with shadcn/ui components
- **@base-ui/react** (underlying primitive library for shadcn, NOT Radix)
- **next-intl** for i18n
- **@supabase/ssr** for auth session management
- **React Hook Form + Zod** for form validation
- **Recharts** for dashboard charts
- **Sonner** for toast notifications
- **Lucide React** for icons

## Key Design Decisions

- **Year-scoped workspaces:** All financial modules use `/:year/` URL prefix. Wishlist is the exception (persists across years).
- **Shared categories:** Expenses and Debt share the same category list (domain: "expense"). Income and Wishlist each have their own.
- **Typed payment methods:** Cash, Debit Card, Credit Card, Digital Wallet, Crypto. Cards module only works with credit_card type.
- **Budget calculation:** Spent = sum of Expenses + sum of Debts for each category/month.
- **Dashboard aggregation:** Single endpoint fires 9 parallel goroutines for all chart data.
- **Multi-tenant isolation:** Every query is scoped by `user_id` extracted from the JWT.
