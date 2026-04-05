# Plan 3: Dashboards & Cards Modules

> **For agentic workers:** Use subagent-driven-development to implement this plan task-by-task.

**Goal:** Build the Dashboard home page with 6 charts (net savings, income breakdown, expenses breakdown, income vs expenses, daily expenses, daily debt) and the Cards module for credit card health tracking with auto-calculated usage from Debt entries.

**Depends on:** Plan 2 complete (expenses, debts, budgets, all CRUD working).

**Chart library:** Recharts — lightweight, React-native, composable, good SSR support. Install via `npm install recharts`.

---

## Phase A: Backend — Dashboard Aggregation Endpoints

The dashboard needs pre-aggregated data from the backend. Rather than making the frontend fetch raw lists and compute charts client-side, we'll add dedicated aggregation endpoints.

### Task 1: Dashboard Aggregation Service + Handler

**Files:**
- Create: `backend/internal/service/dashboard.go`
- Create: `backend/internal/handler/dashboard.go`
- Create: `backend/internal/handler/dashboard_test.go`

The dashboard service needs access to Income, Expense, and Debt repositories.

```go
type DashboardService struct {
	incomeRepo  *repository.IncomeRepository
	expenseRepo *repository.ExpenseRepository
	debtRepo    *repository.DebtRepository
}
```

**Endpoint: `GET /api/v1/years/:year/dashboard`**

Returns all chart data in a single response to minimize round-trips:

```go
type DashboardData struct {
	NetSavings       []MonthlyNet       `json:"net_savings"`
	IncomeBreakdown  []CategorySum      `json:"income_breakdown"`
	ExpenseBreakdown []TypeSum          `json:"expense_breakdown"`
	IncomeVsExpenses IncomeVsExpenses   `json:"income_vs_expenses"`
	DailyExpenses    []DailySum         `json:"daily_expenses"`
	DailyDebts       []DailySum         `json:"daily_debts"`
}

type MonthlyNet struct {
	Month    int    `json:"month"`
	Income   string `json:"income"`
	Expenses string `json:"expenses"`
	Net      string `json:"net"`
}

type CategorySum struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Total        string `json:"total"`
}

type TypeSum struct {
	Type  string `json:"type"`
	Total string `json:"total"`
}

type IncomeVsExpenses struct {
	TotalIncome   string `json:"total_income"`
	TotalExpenses string `json:"total_expenses"`
}

type DailySum struct {
	Date  string `json:"date"`
	Total string `json:"total"`
}
```

**Repository methods needed (add to existing repos):**

Add to `IncomeRepository`:
```go
// SumByMonth returns income totals grouped by month for a year
func (r *IncomeRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error)

// SumByCategory returns income totals grouped by category for a year
func (r *IncomeRepository) SumByCategory(userID uuid.UUID, year int) ([]CategorySumRow, error)

// TotalByYear returns the total income for a year
func (r *IncomeRepository) TotalByYear(userID uuid.UUID, year int) (decimal.Decimal, error)
```

Add to `ExpenseRepository`:
```go
// SumByMonth returns expense totals grouped by month for a year
func (r *ExpenseRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error)

// SumByType returns expense totals grouped by type (expense/saving/investment) for a year
func (r *ExpenseRepository) SumByType(userID uuid.UUID, year int) ([]TypeSumRow, error)

// TotalByYear returns the total expenses for a year
func (r *ExpenseRepository) TotalByYear(userID uuid.UUID, year int) (decimal.Decimal, error)

// SumByDay returns expense totals grouped by day for a year (or optionally filtered by month)
func (r *ExpenseRepository) SumByDay(userID uuid.UUID, year int, month *int) ([]DaySumRow, error)
```

Add to `DebtRepository`:
```go
// SumByMonth returns debt totals grouped by month for a year
func (r *DebtRepository) SumByMonth(userID uuid.UUID, year int) ([]MonthSum, error)

// SumByDay returns debt totals grouped by day for a year (or optionally filtered by month)
func (r *DebtRepository) SumByDay(userID uuid.UUID, year int, month *int) ([]DaySumRow, error)

// SumByPaymentMethodMonth returns debt totals grouped by payment_method_id for a given month/year
// Used by Cards module for auto-calculation
func (r *DebtRepository) SumByPaymentMethodMonth(userID uuid.UUID, month, year int) ([]PaymentMethodSumRow, error)
```

Shared row types (put in a new file `backend/internal/repository/aggregates.go`):
```go
package repository

import "github.com/shopspring/decimal"

type MonthSum struct {
	Month int             `json:"month"`
	Total decimal.Decimal `json:"total"`
}

type CategorySumRow struct {
	CategoryID   string          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Total        decimal.Decimal `json:"total"`
}

type TypeSumRow struct {
	Type  string          `json:"type"`
	Total decimal.Decimal `json:"total"`
}

type DaySumRow struct {
	Date  string          `json:"date"`
	Total decimal.Decimal `json:"total"`
}

type PaymentMethodSumRow struct {
	PaymentMethodID string          `json:"payment_method_id"`
	Total           decimal.Decimal `json:"total"`
}
```

**SQL patterns for the aggregation queries:**

```sql
-- SumByMonth (income/expense/debt)
SELECT EXTRACT(MONTH FROM date)::int AS month, COALESCE(SUM(amount), 0) AS total
FROM {table} WHERE user_id = ? AND year = ? GROUP BY month ORDER BY month

-- SumByCategory (income)
SELECT c.id AS category_id, c.name AS category_name, COALESCE(SUM(i.amount), 0) AS total
FROM incomes i LEFT JOIN categories c ON i.category_id = c.id
WHERE i.user_id = ? AND i.year = ? GROUP BY c.id, c.name ORDER BY total DESC

-- SumByType (expense)
SELECT type, COALESCE(SUM(amount), 0) AS total
FROM expenses WHERE user_id = ? AND year = ? GROUP BY type

-- TotalByYear
SELECT COALESCE(SUM(amount), 0) FROM {table} WHERE user_id = ? AND year = ?

-- SumByDay (expense/debt) — optionally filtered by month
SELECT date::text AS date, COALESCE(SUM(amount), 0) AS total
FROM {table} WHERE user_id = ? AND year = ? [AND EXTRACT(MONTH FROM date) = ?]
GROUP BY date ORDER BY date

-- SumByPaymentMethodMonth (debt — for Cards)
SELECT payment_method_id::text, COALESCE(SUM(amount), 0) AS total
FROM debts WHERE user_id = ? AND payment_method_id IS NOT NULL
AND EXTRACT(MONTH FROM date) = ? AND year = ?
GROUP BY payment_method_id
```

**Dashboard service logic:**

`GetDashboard(userID, year)` calls all the aggregation repo methods in parallel (use goroutines + errgroup), then assembles `DashboardData`:

1. **net_savings**: For each month 1-12, `net = income_month - expense_month - debt_month`. Use SumByMonth from all 3 repos. Fill missing months with zero.
2. **income_breakdown**: Directly from `IncomeRepository.SumByCategory`
3. **expense_breakdown**: From `ExpenseRepository.SumByType`
4. **income_vs_expenses**: From TotalByYear for income and expenses
5. **daily_expenses**: From `ExpenseRepository.SumByDay`
6. **daily_debts**: From `DebtRepository.SumByDay`

Note: The daily endpoints also support an optional `?month=` query param for filtering.

- [ ] **Step 1: Create `repository/aggregates.go`** with shared row types
- [ ] **Step 2: Add aggregation methods to IncomeRepository** (SumByMonth, SumByCategory, TotalByYear)
- [ ] **Step 3: Add aggregation methods to ExpenseRepository** (SumByMonth, SumByType, TotalByYear, SumByDay)
- [ ] **Step 4: Add aggregation methods to DebtRepository** (SumByMonth, SumByDay, SumByPaymentMethodMonth)
- [ ] **Step 5: Create `service/dashboard.go`** with GetDashboard using errgroup
- [ ] **Step 6: Create `handler/dashboard.go`** with GetDashboard endpoint
- [ ] **Step 7: Write handler test**
- [ ] **Step 8: Verify backend compiles**
- [ ] **Step 9: Commit**

---

## Phase B: Backend — Cards Module

### Task 2: Cards Migration

**Files:**
- Create: `backend/migrations/000007_create_cards.up.sql`
- Create: `backend/migrations/000007_create_cards.down.sql`

```sql
-- 000007_create_cards.up.sql
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    payment_method_id UUID NOT NULL REFERENCES payment_methods(id) ON DELETE CASCADE,
    bank VARCHAR(255) NOT NULL,
    card_limit DECIMAL(12, 2) NOT NULL,
    recommended_max_pct DECIMAL(5, 2) NOT NULL DEFAULT 30.00,
    manual_usage_override DECIMAL(12, 2),
    level VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, payment_method_id)
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
```

```sql
-- 000007_create_cards.down.sql
DROP TABLE IF EXISTS cards;
```

Fields:
- `payment_method_id` — links to a credit_card type payment method (enforced at service level, not DB)
- `card_limit` — the credit limit
- `recommended_max_pct` — configurable (default 30%), the "safe zone" threshold
- `manual_usage_override` — if set, adds to the auto-calculated usage from Debts
- `level` — optional label (Gold, Platinum, etc.)

- [ ] **Step 1: Create migration files**
- [ ] **Step 2: Commit**

---

### Task 3: Cards Model + CRUD + Health Calculation

**Files:**
- Create: `backend/internal/model/card.go`
- Create: `backend/internal/repository/card.go`
- Create: `backend/internal/repository/card_test.go`
- Create: `backend/internal/service/card.go`
- Create: `backend/internal/handler/card.go`
- Create: `backend/internal/handler/card_test.go`

**Model:**
```go
type Card struct {
	Base
	PaymentMethodID    uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex:idx_user_pm" json:"payment_method_id"`
	PaymentMethod      *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Bank               string          `gorm:"type:varchar(255);not null" json:"bank"`
	CardLimit          decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"card_limit"`
	RecommendedMaxPct  decimal.Decimal `gorm:"type:decimal(5,2);not null;default:30.00" json:"recommended_max_pct"`
	ManualUsageOverride *decimal.Decimal `gorm:"type:decimal(12,2)" json:"manual_usage_override,omitempty"`
	Level              *string         `gorm:"type:varchar(50)" json:"level,omitempty"`
}

func (Card) TableName() string { return "cards" }
```

**Repository methods:**
- `Create(card *Card) error`
- `ListByUser(userID uuid.UUID) ([]Card, error)` — preloads PaymentMethod
- `GetByID(userID, id uuid.UUID) (*Card, error)` — preloads PaymentMethod
- `Update(card *Card) error`
- `Delete(userID, id uuid.UUID) error`

**Service — Card health summary:**

The service needs the DebtRepository for auto-calculating usage.

```go
type CardService struct {
	repo     *repository.CardRepository
	debtRepo *repository.DebtRepository
}

type CardSummary struct {
	Card            model.Card      `json:"card"`
	AutoUsage       string          `json:"auto_usage"`       // Sum of debts for this card's payment method this month
	ManualOverride  *string         `json:"manual_override"`  // From card.ManualUsageOverride
	TotalUsage      string          `json:"total_usage"`      // auto + manual
	UsagePercent    float64         `json:"usage_percent"`    // (total_usage / card_limit) * 100
	RecommendedMax  string          `json:"recommended_max"`  // card_limit * recommended_max_pct / 100
	HealthColor     string          `json:"health_color"`     // green/yellow/orange/red
}
```

`GetCardSummaries(userID, month, year)`:
1. List all cards for user
2. Get `DebtRepository.SumByPaymentMethodMonth(userID, month, year)` — returns map of payment_method_id → total
3. For each card:
   - `auto_usage` = debt sum for that card's payment_method_id (or 0)
   - `total_usage` = auto_usage + manual_usage_override (if set)
   - `usage_percent` = (total_usage / card_limit) * 100
   - `health_color`: 0-20% → green, 21-30% → yellow, 31-70% → orange, 71%+ → red

**Handler endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/cards` | Create card (validates payment_method is credit_card type) |
| GET | `/api/v1/cards?month=&year=` | Get card summaries with health for month/year |
| GET | `/api/v1/cards/:id` | Get single card |
| PUT | `/api/v1/cards/:id` | Update card (limit, recommended_max_pct, manual_override, level) |
| DELETE | `/api/v1/cards/:id` | Delete card |

**Service validation on Create:**
- Fetch the payment method by ID, verify its Type is `credit_card`
- If not, return error "payment method must be of type credit_card"

- [ ] **Step 1: Create card model**
- [ ] **Step 2: Write failing repository tests** (Create, ListByUser, Update, Delete)
- [ ] **Step 3: Implement card repository**
- [ ] **Step 4: Implement card service** with GetCardSummaries and health calculation
- [ ] **Step 5: Implement card handler** with validation
- [ ] **Step 6: Write handler test**
- [ ] **Step 7: Verify all tests pass**
- [ ] **Step 8: Commit**

---

### Task 4: Wire Dashboard + Cards Routes

**Files:**
- Modify: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`

Add routes:
```go
// Dashboard
api.Get("/years/:year/dashboard", dashboardHandler.GetDashboard)

// Cards
cards := api.Group("/cards")
cards.Post("/", cardHandler.Create)
cards.Get("/", cardHandler.GetSummaries)
cards.Get("/:id", cardHandler.GetByID)
cards.Put("/:id", cardHandler.Update)
cards.Delete("/:id", cardHandler.Delete)
```

Wire DI in `main.go`:
- DashboardService needs incomeRepo, expenseRepo, debtRepo
- CardService needs cardRepo, debtRepo

- [ ] **Step 1: Update router.go**
- [ ] **Step 2: Update main.go**
- [ ] **Step 3: Verify backend compiles**
- [ ] **Step 4: Commit**

---

## Phase C: Frontend — Install Recharts + Dashboard Page

### Task 5: Install Recharts + Dashboard Types + i18n

**Files:**
- Modify: `frontend/package.json` (via npm install)
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/lib/i18n/en.json`
- Modify: `frontend/src/lib/i18n/es.json`

Install: `cd frontend && npm install recharts`

Add types:
```ts
export interface DashboardData {
  net_savings: MonthlyNet[];
  income_breakdown: CategorySumData[];
  expense_breakdown: TypeSumData[];
  income_vs_expenses: IncomeVsExpenses;
  daily_expenses: DailySumData[];
  daily_debts: DailySumData[];
}

export interface MonthlyNet {
  month: number;
  income: string;
  expenses: string;
  net: string;
}

export interface CategorySumData {
  category_id: string;
  category_name: string;
  total: string;
}

export interface TypeSumData {
  type: string;
  total: string;
}

export interface IncomeVsExpenses {
  total_income: string;
  total_expenses: string;
}

export interface DailySumData {
  date: string;
  total: string;
}

export interface Card {
  id: string;
  user_id: string;
  payment_method_id: string;
  payment_method?: PaymentMethod;
  bank: string;
  card_limit: string;
  recommended_max_pct: string;
  manual_usage_override?: string;
  level?: string;
  created_at: string;
  updated_at: string;
}

export interface CardSummary {
  card: Card;
  auto_usage: string;
  manual_override?: string;
  total_usage: string;
  usage_percent: number;
  recommended_max: string;
  health_color: "green" | "yellow" | "orange" | "red";
}
```

Add i18n keys:

**EN:**
```json
"dashboard": {
  "title": "Dashboard",
  "netSavings": "Net Savings Over Year",
  "incomeBreakdown": "Income Breakdown",
  "expenseBreakdown": "Expenses Breakdown",
  "incomeVsExpenses": "Income vs Expenses",
  "dailyExpenses": "Daily Expenses",
  "dailyDebts": "Daily Debt (Credit Cards)",
  "income": "Income",
  "expenses": "Expenses",
  "debt": "Debt",
  "net": "Net",
  "noData": "No data for this year yet"
},
"cards": {
  "title": "Cards",
  "bank": "Bank",
  "cardLimit": "Credit Limit",
  "recommendedMax": "Recommended Max",
  "usageThisMonth": "Usage This Month",
  "healthIndicator": "Health",
  "level": "Level",
  "addCard": "Add Card",
  "editCard": "Edit Card",
  "noCards": "No credit cards registered yet",
  "manualOverride": "Manual Usage Override",
  "autoUsage": "Auto-Calculated",
  "totalUsage": "Total Usage",
  "recommendedMaxPct": "Recommended Max %",
  "selectCreditCard": "Select a credit card payment method",
  "healthy": "Healthy",
  "recommended": "At Recommended Max",
  "warning": "Warning",
  "danger": "Danger"
}
```

**ES:**
```json
"dashboard": {
  "title": "Inicio",
  "netSavings": "Ahorro Neto del Año",
  "incomeBreakdown": "Desglose de Ingresos",
  "expenseBreakdown": "Desglose de Gastos",
  "incomeVsExpenses": "Ingresos vs Gastos",
  "dailyExpenses": "Gastos Diarios",
  "dailyDebts": "Deudas Diarias (Tarjetas)",
  "income": "Ingresos",
  "expenses": "Gastos",
  "debt": "Deudas",
  "net": "Neto",
  "noData": "Sin datos para este año"
},
"cards": {
  "title": "Tarjetas",
  "bank": "Banco",
  "cardLimit": "Límite de Crédito",
  "recommendedMax": "Máximo Recomendado",
  "usageThisMonth": "Uso Este Mes",
  "healthIndicator": "Salud",
  "level": "Nivel",
  "addCard": "Agregar Tarjeta",
  "editCard": "Editar Tarjeta",
  "noCards": "No hay tarjetas de crédito registradas",
  "manualOverride": "Ajuste Manual de Uso",
  "autoUsage": "Cálculo Automático",
  "totalUsage": "Uso Total",
  "recommendedMaxPct": "% Máximo Recomendado",
  "selectCreditCard": "Selecciona un método de pago tipo tarjeta de crédito",
  "healthy": "Saludable",
  "recommended": "En Máximo Recomendado",
  "warning": "Advertencia",
  "danger": "Peligro"
}
```

- [ ] **Step 1: Install recharts**
- [ ] **Step 2: Add types**
- [ ] **Step 3: Add EN translations**
- [ ] **Step 4: Add ES translations**
- [ ] **Step 5: Commit**

---

### Task 6: Dashboard Page with 6 Charts

**Files:**
- Create: `frontend/src/app/[year]/dashboard/page.tsx`
- Create: `frontend/src/components/dashboard/net-savings-chart.tsx`
- Create: `frontend/src/components/dashboard/income-breakdown-chart.tsx`
- Create: `frontend/src/components/dashboard/expense-breakdown-chart.tsx`
- Create: `frontend/src/components/dashboard/income-vs-expenses-chart.tsx`
- Create: `frontend/src/components/dashboard/daily-expenses-chart.tsx`
- Create: `frontend/src/components/dashboard/daily-debts-chart.tsx`
- Modify: `frontend/src/components/layout/sidebar.tsx` (add dashboard nav item)

**Dashboard page layout:**
```
┌──────────────────────────────────────────────┐
│ Dashboard — 2026                             │
├──────────────────────┬───────────────────────┤
│ Net Savings (line)   │ Income vs Exp (donut) │
│ full width or 2/3    │ 1/3                   │
├──────────────────────┴───────────────────────┤
│ Income Breakdown (donut) │ Expense Breakdown  │
│ 1/2                      │ (donut) 1/2        │
├──────────────────────────┴───────────────────┤
│ Daily Expenses (bar) — full width            │
├──────────────────────────────────────────────┤
│ Daily Debt/TDCs (bar) — full width           │
└──────────────────────────────────────────────┘
```

**Chart components (all "use client"):**

#### `net-savings-chart.tsx`
- Recharts `LineChart` with `ResponsiveContainer`
- X-axis: month names (Jan-Dec)
- Lines: Income (green), Expenses (red), Net (blue, dashed)
- Tooltip with formatted values
- Props: `data: MonthlyNet[]`

#### `income-breakdown-chart.tsx`
- Recharts `PieChart` with `Pie` (donut: innerRadius=60, outerRadius=80)
- Each slice = a category, colored with a palette
- Legend below
- Props: `data: CategorySumData[]`

#### `expense-breakdown-chart.tsx`
- Recharts `PieChart` donut
- Slices: Expense, Saving, Investment (3 fixed colors)
- Props: `data: TypeSumData[]`

#### `income-vs-expenses-chart.tsx`
- Recharts `PieChart` donut
- 2 slices: Income (green), Expenses (red)
- Center label showing the difference
- Props: `data: IncomeVsExpenses`

#### `daily-expenses-chart.tsx`
- Recharts `BarChart` (vertical bars)
- X-axis: dates, Y-axis: amounts
- Color: a warm color (orange/coral)
- Props: `data: DailySumData[]`

#### `daily-debts-chart.tsx`
- Same as daily-expenses but different color (purple/red)
- Props: `data: DailySumData[]`

**Dashboard page (`[year]/dashboard/page.tsx`):**
- Fetches `apiGet<DashboardData>(/api/v1/years/${year}/dashboard)`
- Renders all 6 charts in the grid layout above
- Loading state, empty state

**Sidebar update:**
- Add `{ key: "dashboard", path: "dashboard", icon: "▣" }` as the FIRST nav item
- Update root redirect: `page.tsx` should redirect to `/${year}/dashboard` instead of `/${year}/income`

- [ ] **Step 1: Create chart components** (all 6)
- [ ] **Step 2: Create dashboard page**
- [ ] **Step 3: Update sidebar** to add dashboard as first nav item
- [ ] **Step 4: Update root page redirect** to dashboard
- [ ] **Step 5: Verify build**
- [ ] **Step 6: Commit**

---

### Task 7: Cards Page UI

**Files:**
- Create: `frontend/src/app/[year]/cards/page.tsx`
- Create: `frontend/src/components/cards/card-health-card.tsx`
- Create: `frontend/src/components/cards/card-form.tsx`

**Card health card component (`card-health-card.tsx`):**
A visual card (shadcn Card) showing:
- Bank name + level badge
- Card name (from payment_method.name)
- Credit limit
- Usage bar (total_usage / card_limit) with color = health_color
- Usage details: auto-calculated + manual override = total
- Health indicator colored badge (Green/Yellow/Orange/Red)
- Edit/Delete buttons

Props: `summary: CardSummary, onEdit, onDelete`

**Card form (`card-form.tsx`):**
Dialog form with:
- Payment method select (filtered to credit_card type only from payment methods list)
- Bank name (text)
- Card limit (number)
- Recommended max % (number, default 30)
- Manual usage override (number, optional)
- Level (text, optional)

Props: `open, onClose, onSubmit, creditCardPaymentMethods: PaymentMethod[], defaultValues?: Card`

**Cards page (`[year]/cards/page.tsx`):**
- Month selector at top (same as budget page — defaults to current month)
- Fetches `apiGet<{data: CardSummary[]}>(/api/v1/cards?month=${month}&year=${year})`
- Fetches payment methods (filtered to credit_card type for the form)
- Renders cards in a responsive grid (2-3 columns)
- "Add Card" button opens form

- [ ] **Step 1: Create card health display component**
- [ ] **Step 2: Create card form dialog**
- [ ] **Step 3: Create cards page with month selector**
- [ ] **Step 4: Verify build**
- [ ] **Step 5: Commit**

---

## Phase D: Verification

### Task 8: End-to-End Verification

- [ ] **Step 1: Run all backend tests**: `cd backend && go test ./... -v -count=1`
- [ ] **Step 2: Verify backend compiles**: `cd backend && go build ./cmd/server`
- [ ] **Step 3: Verify frontend builds**: `cd frontend && npm run build`
- [ ] **Step 4: Final commit if needed**

---

## API Endpoint Summary (New in Plan 3)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/years/:year/dashboard` | All dashboard chart data |
| POST | `/api/v1/cards` | Create card |
| GET | `/api/v1/cards?month=&year=` | Get card summaries with health |
| GET | `/api/v1/cards/:id` | Get single card |
| PUT | `/api/v1/cards/:id` | Update card |
| DELETE | `/api/v1/cards/:id` | Delete card |

## Frontend Routes (New in Plan 3)

| Route | Description |
|-------|-------------|
| `/[year]/dashboard` | Dashboard with 6 charts (new home page) |
| `/[year]/cards` | Credit card health tracking |
