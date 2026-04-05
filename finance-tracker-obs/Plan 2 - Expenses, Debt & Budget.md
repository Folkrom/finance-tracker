# Plan 2: Expenses, Debt & Budget Modules

> **For agentic workers:** Use subagent-driven-development to implement this plan task-by-task.

**Goal:** Add Expenses, Debt, and Budget modules — complete backend CRUD + frontend UI. Expenses and Debt share the same category list (domain: "expense"). Budget tracks monthly spending limits per category with recurring defaults and per-month overrides. Budget calculation: `remaining = limit - expenses_sum - debt_sum`.

**Depends on:** Plan 1 complete (categories, payment methods, income, auth, app shell).

**Patterns to follow:** Mirror the Income module's architecture (repo → service → handler, year-scoped routes, frontend page with table + form dialog). See existing code in `backend/internal/` and `frontend/src/`.

---

## Phase A: Backend — Expenses & Debt

### Task 1: Expense & Debt Migrations

**Files:**
- Create: `backend/migrations/000004_create_expenses.up.sql`
- Create: `backend/migrations/000004_create_expenses.down.sql`
- Create: `backend/migrations/000005_create_debts.up.sql`
- Create: `backend/migrations/000005_create_debts.down.sql`

- [ ] **Step 1: Create expenses migration**

`000004_create_expenses.up.sql`:
```sql
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    date DATE NOT NULL,
    year INT NOT NULL,
    payment_method_id UUID REFERENCES payment_methods(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'expense' CHECK (type IN ('expense', 'saving', 'investment')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_year ON expenses(user_id, year);
CREATE INDEX idx_expenses_date ON expenses(user_id, date);
CREATE INDEX idx_expenses_category ON expenses(user_id, category_id);
```

`000004_create_expenses.down.sql`:
```sql
DROP TABLE IF EXISTS expenses;
```

- [ ] **Step 2: Create debts migration**

`000005_create_debts.up.sql`:
```sql
CREATE TABLE debts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    date DATE NOT NULL,
    year INT NOT NULL,
    payment_method_id UUID REFERENCES payment_methods(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_debts_user_id ON debts(user_id);
CREATE INDEX idx_debts_year ON debts(user_id, year);
CREATE INDEX idx_debts_date ON debts(user_id, date);
CREATE INDEX idx_debts_category ON debts(user_id, category_id);
CREATE INDEX idx_debts_payment_method ON debts(user_id, payment_method_id);
```

`000005_create_debts.down.sql`:
```sql
DROP TABLE IF EXISTS debts;
```

- [ ] **Step 3: Commit**

---

### Task 2: Expense & Debt Models

**Files:**
- Create: `backend/internal/model/expense.go`
- Create: `backend/internal/model/debt.go`

- [ ] **Step 1: Create expense model**

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseType string

const (
	ExpenseTypeExpense    ExpenseType = "expense"
	ExpenseTypeSaving     ExpenseType = "saving"
	ExpenseTypeInvestment ExpenseType = "investment"
)

type Expense struct {
	Base
	Name            string          `gorm:"type:varchar(255);not null" json:"name"`
	Amount          decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string          `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	Date            time.Time       `gorm:"type:date;not null" json:"date"`
	Year            int             `gorm:"not null;index" json:"year"`
	PaymentMethodID *uuid.UUID      `gorm:"type:uuid" json:"payment_method_id,omitempty"`
	PaymentMethod   *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	CategoryID      *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category        *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Type            ExpenseType     `gorm:"type:varchar(20);not null;default:expense" json:"type"`
}

func (Expense) TableName() string {
	return "expenses"
}
```

- [ ] **Step 2: Create debt model**

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Debt struct {
	Base
	Name            string          `gorm:"type:varchar(255);not null" json:"name"`
	Amount          decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string          `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	Date            time.Time       `gorm:"type:date;not null" json:"date"`
	Year            int             `gorm:"not null;index" json:"year"`
	PaymentMethodID *uuid.UUID      `gorm:"type:uuid" json:"payment_method_id,omitempty"`
	PaymentMethod   *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	CategoryID      *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category        *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Debt) TableName() string {
	return "debts"
}
```

- [ ] **Step 3: Verify models compile**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go build ./...`

- [ ] **Step 4: Commit**

---

### Task 3: Expense CRUD (Repo + Service + Handler + Tests)

**Files:**
- Create: `backend/internal/repository/expense.go`
- Create: `backend/internal/repository/expense_test.go`
- Create: `backend/internal/service/expense.go`
- Create: `backend/internal/handler/expense.go`
- Create: `backend/internal/handler/expense_test.go`

Follow the exact same patterns as the Income module. Key differences from Income:

- **Expense has:** `Name` (not `Source`), `PaymentMethodID`, `Type` (expense/saving/investment)
- **Repository:** `ListByYear` preloads both `Category` and `PaymentMethod`
- **Handler request struct:**
  ```go
  type createExpenseRequest struct {
      Name            string  `json:"name"`
      Amount          string  `json:"amount"`
      Currency        string  `json:"currency"`
      Date            string  `json:"date"`
      PaymentMethodID *string `json:"payment_method_id"`
      CategoryID      *string `json:"category_id"`
      Type            string  `json:"type"`
  }
  ```
- **Validation:** `type` must be one of "expense", "saving", "investment" — default to "expense" if empty
- **Parse helper:** Reuse `parseIncomeRequest` pattern but add `paymentMethodID` parsing. Create a shared `parseCommonFields` in `response.go` or a new `parse.go`:
  ```go
  func parseMoneyFields(amountStr string, categoryIDStr *string, paymentMethodIDStr *string, dateStr string) (
      decimal.Decimal, *uuid.UUID, *uuid.UUID, time.Time, error,
  )
  ```

- [ ] **Step 1: Write failing repository tests** (Create, ListByYear, Update, Delete)
- [ ] **Step 2: Implement expense repository**
- [ ] **Step 3: Run repository tests — all pass**
- [ ] **Step 4: Implement expense service** (year auto-derived from date, defaults currency to MXN, defaults type to "expense")
- [ ] **Step 5: Implement expense handler** (Create, ListByYear, GetByID, Update, Delete)
- [ ] **Step 6: Write handler test** (CreateAndList)
- [ ] **Step 7: Verify all tests pass**
- [ ] **Step 8: Commit**

---

### Task 4: Debt CRUD (Repo + Service + Handler + Tests)

**Files:**
- Create: `backend/internal/repository/debt.go`
- Create: `backend/internal/repository/debt_test.go`
- Create: `backend/internal/service/debt.go`
- Create: `backend/internal/handler/debt.go`
- Create: `backend/internal/handler/debt_test.go`

Debt is very similar to Expense but without the `Type` field. Same shared parse helper.

- [ ] **Step 1: Write failing repository tests**
- [ ] **Step 2: Implement debt repository** (preloads Category + PaymentMethod)
- [ ] **Step 3: Run repository tests — all pass**
- [ ] **Step 4: Implement debt service**
- [ ] **Step 5: Implement debt handler**
- [ ] **Step 6: Write handler test**
- [ ] **Step 7: Verify all tests pass**
- [ ] **Step 8: Commit**

---

## Phase B: Backend — Budget

### Task 5: Budget Migration

**Files:**
- Create: `backend/migrations/000006_create_budgets.up.sql`
- Create: `backend/migrations/000006_create_budgets.down.sql`

```sql
-- 000006_create_budgets.up.sql
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    monthly_limit DECIMAL(12, 2) NOT NULL,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INT NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, category_id, month, year)
);

CREATE INDEX idx_budgets_user_year ON budgets(user_id, year);
CREATE INDEX idx_budgets_user_category ON budgets(user_id, category_id);
```

`is_recurring = true` marks a row as a template. When loading budgets for a specific month, the logic is:
1. Look for a specific override (is_recurring=false) for that month/year
2. If none, fall back to the recurring template (is_recurring=true) for that category

- [ ] **Step 1: Create migration files**
- [ ] **Step 2: Commit**

---

### Task 6: Budget Model + CRUD + Calculation

**Files:**
- Create: `backend/internal/model/budget.go`
- Create: `backend/internal/repository/budget.go`
- Create: `backend/internal/repository/budget_test.go`
- Create: `backend/internal/service/budget.go`
- Create: `backend/internal/handler/budget.go`
- Create: `backend/internal/handler/budget_test.go`

**Model:**
```go
type Budget struct {
    Base
    CategoryID   uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
    Category     *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    MonthlyLimit decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"monthly_limit"`
    Month        int             `gorm:"not null" json:"month"`
    Year         int             `gorm:"not null" json:"year"`
    IsRecurring  bool            `gorm:"not null;default:false" json:"is_recurring"`
}
```

**Repository methods:**
- `Create(budget *Budget) error`
- `ListByMonthYear(userID uuid.UUID, month, year int) ([]Budget, error)` — returns both overrides and recurring, preloads Category
- `ListRecurring(userID uuid.UUID) ([]Budget, error)` — only is_recurring=true
- `GetByID(userID, id uuid.UUID) (*Budget, error)`
- `Update(budget *Budget) error`
- `Delete(userID, id uuid.UUID) error`

**Service — key logic:**

`GetBudgetSummary(userID, month, year)` returns a list of `BudgetLine`:
```go
type BudgetLine struct {
    Budget    Budget          `json:"budget"`
    Spent     decimal.Decimal `json:"spent"`
    Remaining decimal.Decimal `json:"remaining"`
}
```

Calculation for each category:
1. Get all budgets for the month/year (overrides take precedence over recurring)
2. For each budget's category, sum Expenses + Debts for that user/category/month/year
3. `remaining = monthly_limit - spent`

The service needs access to the Expense and Debt repositories for the sum queries. Add sum methods:
- `ExpenseRepository.SumByCategoryMonth(userID, categoryID, month, year) (decimal.Decimal, error)`
- `DebtRepository.SumByCategoryMonth(userID, categoryID, month, year) (decimal.Decimal, error)`

**Handler endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/budgets` | Create budget (override or recurring) |
| GET | `/api/v1/budgets?month=&year=` | Get budget summary for month/year |
| GET | `/api/v1/budgets/recurring` | List recurring templates |
| PUT | `/api/v1/budgets/:id` | Update budget |
| DELETE | `/api/v1/budgets/:id` | Delete budget |

- [ ] **Step 1: Create budget model**
- [ ] **Step 2: Write failing repository tests**
- [ ] **Step 3: Implement budget repository**
- [ ] **Step 4: Add `SumByCategoryMonth` to expense and debt repositories**
- [ ] **Step 5: Implement budget service with GetBudgetSummary**
- [ ] **Step 6: Implement budget handler**
- [ ] **Step 7: Write handler test**
- [ ] **Step 8: Verify all tests pass**
- [ ] **Step 9: Commit**

---

### Task 7: Wire Routes + Update main.go

**Files:**
- Modify: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`

Add routes:
```go
// Expenses (year-scoped)
api.Post("/years/:year/expenses", expenseHandler.Create)
api.Get("/years/:year/expenses", expenseHandler.ListByYear)
api.Get("/expenses/:id", expenseHandler.GetByID)
api.Put("/expenses/:id", expenseHandler.Update)
api.Delete("/expenses/:id", expenseHandler.Delete)

// Debts (year-scoped)
api.Post("/years/:year/debts", debtHandler.Create)
api.Get("/years/:year/debts", debtHandler.ListByYear)
api.Get("/debts/:id", debtHandler.GetByID)
api.Put("/debts/:id", debtHandler.Update)
api.Delete("/debts/:id", debtHandler.Delete)

// Budgets
budgets := api.Group("/budgets")
budgets.Post("/", budgetHandler.Create)
budgets.Get("/", budgetHandler.GetSummary)
budgets.Get("/recurring", budgetHandler.ListRecurring)
budgets.Put("/:id", budgetHandler.Update)
budgets.Delete("/:id", budgetHandler.Delete)
```

Wire DI in `main.go` — add expense/debt/budget repos, services, handlers. Pass them to `router.Setup()`.

- [ ] **Step 1: Update router.go** (add new handler params + routes)
- [ ] **Step 2: Update main.go** (wire new repos → services → handlers)
- [ ] **Step 3: Verify backend compiles**: `go build ./cmd/server`
- [ ] **Step 4: Run all tests**: `go test ./... -v`
- [ ] **Step 5: Commit**

---

## Phase C: Frontend — Expenses & Debt UI

### Task 8: TypeScript Types + i18n Updates

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/lib/i18n/en.json`
- Modify: `frontend/src/lib/i18n/es.json`

Add types:
```ts
export interface Expense {
  id: string;
  user_id: string;
  name: string;
  amount: string;
  currency: string;
  date: string;
  year: number;
  payment_method_id?: string;
  payment_method?: PaymentMethod;
  category_id?: string;
  category?: Category;
  type: "expense" | "saving" | "investment";
  created_at: string;
  updated_at: string;
}

export interface Debt {
  id: string;
  user_id: string;
  name: string;
  amount: string;
  currency: string;
  date: string;
  year: number;
  payment_method_id?: string;
  payment_method?: PaymentMethod;
  category_id?: string;
  category?: Category;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  user_id: string;
  category_id: string;
  category?: Category;
  monthly_limit: string;
  month: number;
  year: number;
  is_recurring: boolean;
  created_at: string;
  updated_at: string;
}

export interface BudgetLine {
  budget: Budget;
  spent: string;
  remaining: string;
}
```

Add i18n keys (EN):
```json
"expenses": {
  "title": "Expenses",
  "name": "Name",
  "amount": "Amount",
  "category": "Category",
  "paymentMethod": "Payment Method",
  "type": "Type",
  "date": "Date",
  "addExpense": "Add Expense",
  "editExpense": "Edit Expense",
  "noExpenses": "No expense records yet",
  "typeExpense": "Expense",
  "typeSaving": "Saving",
  "typeInvestment": "Investment"
},
"debt": {
  "title": "Debt",
  "name": "Name",
  "amount": "Amount",
  "category": "Category",
  "paymentMethod": "Payment Method",
  "date": "Date",
  "addDebt": "Add Debt",
  "editDebt": "Edit Debt",
  "noDebts": "No debt records yet"
},
"budget": {
  "title": "Budget",
  "category": "Category",
  "limit": "Monthly Limit",
  "spent": "Spent",
  "remaining": "Remaining",
  "addBudget": "Add Budget",
  "editBudget": "Edit Budget",
  "noBudgets": "No budgets set for this month",
  "recurring": "Recurring Defaults",
  "override": "Monthly Override",
  "month": "Month",
  "isRecurring": "Set as recurring default"
}
```

Add i18n keys (ES):
```json
"expenses": {
  "title": "Gastos",
  "name": "Nombre",
  "amount": "Monto",
  "category": "Categoría",
  "paymentMethod": "Método de Pago",
  "type": "Tipo",
  "date": "Fecha",
  "addExpense": "Agregar Gasto",
  "editExpense": "Editar Gasto",
  "noExpenses": "No hay registros de gastos",
  "typeExpense": "Gasto",
  "typeSaving": "Ahorro",
  "typeInvestment": "Inversión"
},
"debt": {
  "title": "Deudas",
  "name": "Nombre",
  "amount": "Monto",
  "category": "Categoría",
  "paymentMethod": "Método de Pago",
  "date": "Fecha",
  "addDebt": "Agregar Deuda",
  "editDebt": "Editar Deuda",
  "noDebts": "No hay registros de deudas"
},
"budget": {
  "title": "Presupuesto",
  "category": "Categoría",
  "limit": "Límite Mensual",
  "spent": "Gastado",
  "remaining": "Restante",
  "addBudget": "Agregar Presupuesto",
  "editBudget": "Editar Presupuesto",
  "noBudgets": "No hay presupuestos para este mes",
  "recurring": "Valores Recurrentes",
  "override": "Ajuste Mensual",
  "month": "Mes",
  "isRecurring": "Establecer como valor recurrente"
}
```

- [ ] **Step 1: Add types**
- [ ] **Step 2: Add EN translations**
- [ ] **Step 3: Add ES translations**
- [ ] **Step 4: Commit**

---

### Task 9: Expenses Page UI

**Files:**
- Create: `frontend/src/app/[year]/expenses/page.tsx`
- Create: `frontend/src/components/expenses/expense-table.tsx`
- Create: `frontend/src/components/expenses/expense-form.tsx`

Mirror the Income module pattern but with these differences:
- Table columns: Name, Amount, Category, Payment Method, Type (badge), Date, Actions
- Form fields: name, amount, currency (MXN/USD), category_id (select, domain=expense), payment_method_id (select), type (expense/saving/investment), date
- API paths: `/api/v1/years/${year}/expenses`, `/api/v1/categories?domain=expense`, `/api/v1/payment-methods`
- Load categories (domain=expense) + payment methods on mount

- [ ] **Step 1: Create expense table component**
- [ ] **Step 2: Create expense form dialog**
- [ ] **Step 3: Create expenses page**
- [ ] **Step 4: Verify build**
- [ ] **Step 5: Commit**

---

### Task 10: Debt Page UI

**Files:**
- Create: `frontend/src/app/[year]/debt/page.tsx`
- Create: `frontend/src/components/debt/debt-table.tsx`
- Create: `frontend/src/components/debt/debt-form.tsx`

Same as Expenses but without the `Type` field.
- Table columns: Name, Amount, Category, Payment Method, Date, Actions
- Form fields: name, amount, currency, category_id (domain=expense — shared!), payment_method_id, date
- API paths: `/api/v1/years/${year}/debts`

- [ ] **Step 1: Create debt table component**
- [ ] **Step 2: Create debt form dialog**
- [ ] **Step 3: Create debt page**
- [ ] **Step 4: Verify build**
- [ ] **Step 5: Commit**

---

### Task 11: Budget Page UI

**Files:**
- Create: `frontend/src/app/[year]/budget/page.tsx`
- Create: `frontend/src/components/budget/budget-table.tsx`
- Create: `frontend/src/components/budget/budget-form.tsx`

The budget page is different from the others:
- Shows a **month selector** (1-12 dropdown) at the top alongside the year from URL
- Loads `GET /api/v1/budgets?month=M&year=Y` which returns `BudgetLine[]` with spent/remaining calculated
- Table columns: Category, Monthly Limit, Spent, Remaining (color-coded: green if positive, red if negative)
- A "Manage Recurring" section/toggle that shows `GET /api/v1/budgets/recurring` templates
- Form dialog: category_id (select, domain=expense), monthly_limit, is_recurring checkbox
- When creating: if is_recurring, it's a template; otherwise it's an override for the selected month/year

**Budget table component:**
```tsx
// Columns: Category | Limit | Spent | Remaining | Actions
// Remaining cell: green text if >= 0, red text if < 0
// Show progress bar (spent / limit) with color coding
```

- [ ] **Step 1: Create budget table component** (shows summary with progress bars)
- [ ] **Step 2: Create budget form dialog** (category select, amount, recurring checkbox)
- [ ] **Step 3: Create budget page** (month selector + table + manage recurring)
- [ ] **Step 4: Verify build**
- [ ] **Step 5: Commit**

---

## Phase D: Verification

### Task 12: End-to-End Verification

- [ ] **Step 1: Run all backend tests**: `cd backend && go test ./... -v -count=1`
- [ ] **Step 2: Verify backend compiles**: `cd backend && go build ./cmd/server`
- [ ] **Step 3: Verify frontend builds**: `cd frontend && npm run build`
- [ ] **Step 4: Final commit if needed**

---

## API Endpoint Summary (New in Plan 2)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/years/:year/expenses` | Create expense |
| GET | `/api/v1/years/:year/expenses` | List expenses by year |
| GET | `/api/v1/expenses/:id` | Get expense |
| PUT | `/api/v1/expenses/:id` | Update expense |
| DELETE | `/api/v1/expenses/:id` | Delete expense |
| POST | `/api/v1/years/:year/debts` | Create debt |
| GET | `/api/v1/years/:year/debts` | List debts by year |
| GET | `/api/v1/debts/:id` | Get debt |
| PUT | `/api/v1/debts/:id` | Update debt |
| DELETE | `/api/v1/debts/:id` | Delete debt |
| POST | `/api/v1/budgets` | Create budget |
| GET | `/api/v1/budgets?month=&year=` | Get budget summary |
| GET | `/api/v1/budgets/recurring` | List recurring templates |
| PUT | `/api/v1/budgets/:id` | Update budget |
| DELETE | `/api/v1/budgets/:id` | Delete budget |
