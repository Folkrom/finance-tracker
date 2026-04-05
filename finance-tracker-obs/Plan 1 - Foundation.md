# Plan 1: Foundation — Project Setup, Auth, Shared Entities & Income Module

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver a working finance tracker with authentication, shared entities (categories, payment methods), and a complete Income module — end to end.

**Architecture:** Go/Fiber REST API with handler → service → repository layers. Next.js 14 App Router frontend with Tailwind + shadcn/ui. Supabase for PostgreSQL hosting and auth (JWT). Multi-tenant via `user_id` scoping at middleware level.

**Tech Stack:** Go 1.22+, Fiber v2, GORM, Zap, golang-migrate, Wire | Next.js 14, TypeScript, Tailwind CSS, shadcn/ui, next-intl, @supabase/ssr, React Hook Form, Zod | PostgreSQL (Supabase), Docker Compose for local dev.

**Spec reference:** [[Finance Tracker - Design Spec]]

---

## File Structure

```
finance-tracker/
├── docker-compose.yml              # Local Postgres + API + Frontend
├── mise.toml                       # Dev tasks (migrate, dev, test, lint)
├── .gitignore
│
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   ├── Makefile
│   ├── cmd/
│   │   └── server/
│   │       ├── main.go             # Entry point
│   │       ├── wire.go             # Wire injector definition
│   │       └── wire_gen.go         # Wire generated (auto)
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go           # Env config struct
│   │   ├── database/
│   │   │   └── database.go         # GORM connection + logger
│   │   ├── middleware/
│   │   │   ├── auth.go             # Supabase JWT validation
│   │   │   └── auth_test.go
│   │   ├── model/
│   │   │   ├── base.go             # Shared base model (ID, timestamps, user_id)
│   │   │   ├── category.go
│   │   │   ├── payment_method.go
│   │   │   └── income.go
│   │   ├── repository/
│   │   │   ├── category.go
│   │   │   ├── category_test.go
│   │   │   ├── payment_method.go
│   │   │   ├── payment_method_test.go
│   │   │   ├── income.go
│   │   │   └── income_test.go
│   │   ├── service/
│   │   │   ├── category.go
│   │   │   ├── category_test.go
│   │   │   ├── payment_method.go
│   │   │   ├── payment_method_test.go
│   │   │   ├── income.go
│   │   │   └── income_test.go
│   │   ├── handler/
│   │   │   ├── category.go
│   │   │   ├── category_test.go
│   │   │   ├── payment_method.go
│   │   │   ├── payment_method_test.go
│   │   │   ├── income.go
│   │   │   ├── income_test.go
│   │   │   └── response.go         # Shared response helpers
│   │   └── router/
│   │       └── router.go           # All route registration
│   ├── migrations/
│   │   ├── 000001_create_categories.up.sql
│   │   ├── 000001_create_categories.down.sql
│   │   ├── 000002_create_payment_methods.up.sql
│   │   ├── 000002_create_payment_methods.down.sql
│   │   ├── 000003_create_incomes.up.sql
│   │   └── 000003_create_incomes.down.sql
│   └── testutil/
│       └── testutil.go             # Test DB setup, fixtures
│
├── frontend/
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.ts
│   ├── tailwind.config.ts
│   ├── .env.local.example
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx          # Root layout (providers, i18n)
│   │   │   ├── page.tsx            # Redirect to /[year]/dashboard
│   │   │   ├── login/
│   │   │   │   └── page.tsx        # Login page
│   │   │   └── [year]/
│   │   │       ├── layout.tsx      # Year-scoped layout (sidebar + header)
│   │   │       ├── income/
│   │   │       │   └── page.tsx    # Income list + CRUD
│   │   │       └── settings/
│   │   │           └── page.tsx    # Categories + Payment Methods
│   │   ├── components/
│   │   │   ├── layout/
│   │   │   │   ├── sidebar.tsx
│   │   │   │   ├── header.tsx
│   │   │   │   └── year-switcher.tsx
│   │   │   ├── income/
│   │   │   │   ├── income-table.tsx
│   │   │   │   └── income-form.tsx
│   │   │   └── settings/
│   │   │       ├── category-manager.tsx
│   │   │       └── payment-method-manager.tsx
│   │   ├── lib/
│   │   │   ├── api.ts              # Fetch wrapper with auth
│   │   │   ├── supabase/
│   │   │   │   ├── client.ts       # Browser client
│   │   │   │   ├── server.ts       # Server client
│   │   │   │   └── middleware.ts   # Auth middleware for Next.js
│   │   │   └── i18n/
│   │   │       ├── config.ts
│   │   │       ├── en.json
│   │   │       └── es.json
│   │   └── types/
│   │       └── index.ts            # Shared TypeScript types
│   └── middleware.ts               # Next.js middleware (auth redirect)
│
└── finance-tracker-obs/            # Existing Obsidian vault (don't touch)
```

---

## Phase A: Infrastructure

### Task 1: Root Project Setup

**Files:**
- Create: `docker-compose.yml`
- Create: `mise.toml`
- Create: `.gitignore`

- [ ] **Step 1: Initialize git repository**

```bash
cd /home/folkrom/projects/finance-tracker
git init
```

- [ ] **Step 2: Create .gitignore**

```gitignore
# Go
backend/tmp/
backend/.env

# Node
frontend/node_modules/
frontend/.next/
frontend/.env.local

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store

# Obsidian (keep tracked but ignore workspace)
finance-tracker-obs/.obsidian/workspace.json
```

- [ ] **Step 3: Create docker-compose.yml**

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: finance
      POSTGRES_PASSWORD: finance_dev
      POSTGRES_DB: finance_tracker
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

> Note: In production, Supabase manages Postgres. This is for local dev only.

- [ ] **Step 4: Create mise.toml**

```toml
[tools]
go = "1.24"
node = "20"

[tasks.dev-backend]
run = "cd backend && go run ./cmd/server"
description = "Run backend dev server"

[tasks.dev-frontend]
run = "cd frontend && npm run dev"
description = "Run frontend dev server"

[tasks.migrate-up]
run = "cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=migrations -database \"$DATABASE_URL\" up"
description = "Run database migrations"

[tasks.migrate-down]
run = "cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=migrations -database \"$DATABASE_URL\" down 1"
description = "Rollback last migration"

[tasks.test-backend]
run = "cd backend && go test ./... -v"
description = "Run backend tests"

[tasks.test-frontend]
run = "cd frontend && npm test"
description = "Run frontend tests"
```

- [ ] **Step 5: Commit**

```bash
git add .gitignore docker-compose.yml mise.toml
git commit -m "chore: root project setup with docker-compose and mise"
```

---

### Task 2: Backend Scaffolding

**Files:**
- Create: `backend/go.mod`
- Create: `backend/.env.example`
- Create: `backend/Makefile`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/database/database.go`
- Create: `backend/cmd/server/main.go`

- [ ] **Step 1: Initialize Go module and install dependencies**

```bash
cd /home/folkrom/projects/finance-tracker
mkdir -p backend/cmd/server backend/internal/{config,database,middleware,model,repository,service,handler,router} backend/migrations backend/testutil
cd backend
go mod init github.com/folkrom/finance-tracker/backend
go get github.com/gofiber/fiber/v2@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get go.uber.org/zap@latest
go get github.com/golang-jwt/jwt/v5@latest
go get github.com/google/uuid@latest
go get github.com/stretchr/testify@latest
go get github.com/caarlos0/env/v11@latest
```

- [ ] **Step 2: Create .env.example**

```env
# Database
DATABASE_URL=postgres://finance:finance_dev@localhost:5432/finance_tracker?sslmode=disable

# Supabase Auth
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_JWT_SECRET=your-jwt-secret

# Server
PORT=8080
ENVIRONMENT=development
```

- [ ] **Step 3: Create config/config.go**

```go
package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL,required"`
	SupabaseURL      string `env:"SUPABASE_URL,required"`
	SupabaseJWTSecret string `env:"SUPABASE_JWT_SECRET,required"`
	Port             string `env:"PORT" envDefault:"8080"`
	Environment      string `env:"ENVIRONMENT" envDefault:"development"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
```

- [ ] **Step 4: Create database/database.go**

```go
package database

import (
	"github.com/folkrom/finance-tracker/backend/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	logger.Info("database connected")
	return db, nil
}
```

- [ ] **Step 5: Create cmd/server/main.go (minimal, no routes yet)**

```go
package main

import (
	"log"
	"os"

	"github.com/folkrom/finance-tracker/backend/internal/config"
	"github.com/folkrom/finance-tracker/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	if os.Getenv("ENVIRONMENT") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	_, err = database.New(cfg, logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
```

- [ ] **Step 6: Create Makefile**

```makefile
.PHONY: dev test lint

dev:
	go run ./cmd/server

test:
	go test ./... -v -count=1

lint:
	go vet ./...
```

- [ ] **Step 7: Verify it compiles**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go build ./cmd/server`
Expected: No errors.

- [ ] **Step 8: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/
git commit -m "feat: backend scaffolding with Fiber, GORM, Zap, and config"
```

---

### Task 3: Database Migrations

**Files:**
- Create: `backend/migrations/000001_create_categories.up.sql`
- Create: `backend/migrations/000001_create_categories.down.sql`
- Create: `backend/migrations/000002_create_payment_methods.up.sql`
- Create: `backend/migrations/000002_create_payment_methods.down.sql`
- Create: `backend/migrations/000003_create_incomes.up.sql`
- Create: `backend/migrations/000003_create_incomes.down.sql`

- [ ] **Step 1: Create categories migration**

`000001_create_categories.up.sql`:
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(50) NOT NULL CHECK (domain IN ('income', 'expense', 'wishlist')),
    color VARCHAR(7),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, domain)
);

CREATE INDEX idx_categories_user_id ON categories(user_id);
CREATE INDEX idx_categories_domain ON categories(user_id, domain);
```

`000001_create_categories.down.sql`:
```sql
DROP TABLE IF EXISTS categories;
```

- [ ] **Step 2: Create payment_methods migration**

`000002_create_payment_methods.up.sql`:
```sql
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('cash', 'debit_card', 'credit_card', 'digital_wallet', 'crypto')),
    details VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE INDEX idx_payment_methods_user_id ON payment_methods(user_id);
CREATE INDEX idx_payment_methods_type ON payment_methods(user_id, type);
```

`000002_create_payment_methods.down.sql`:
```sql
DROP TABLE IF EXISTS payment_methods;
```

- [ ] **Step 3: Create incomes migration**

`000003_create_incomes.up.sql`:
```sql
CREATE TABLE incomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    source VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    year INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_incomes_user_id ON incomes(user_id);
CREATE INDEX idx_incomes_year ON incomes(user_id, year);
CREATE INDEX idx_incomes_date ON incomes(user_id, date);
```

`000003_create_incomes.down.sql`:
```sql
DROP TABLE IF EXISTS incomes;
```

- [ ] **Step 4: Start Postgres and run migrations**

```bash
cd /home/folkrom/projects/finance-tracker
docker compose up -d postgres
export DATABASE_URL="postgres://finance:finance_dev@localhost:5432/finance_tracker?sslmode=disable"
cd backend
go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=migrations -database "$DATABASE_URL" up
```

Expected: `1/u create_categories`, `2/u create_payment_methods`, `3/u create_incomes` — no errors.

- [ ] **Step 5: Verify tables exist**

```bash
docker compose exec postgres psql -U finance -d finance_tracker -c "\dt"
```

Expected: Tables `categories`, `payment_methods`, `incomes` listed.

- [ ] **Step 6: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/migrations/
git commit -m "feat: database migrations for categories, payment_methods, and incomes"
```

---

### Task 4: GORM Models

**Files:**
- Create: `backend/internal/model/base.go`
- Create: `backend/internal/model/category.go`
- Create: `backend/internal/model/payment_method.go`
- Create: `backend/internal/model/income.go`

- [ ] **Step 1: Create base model**

```go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

- [ ] **Step 2: Create category model**

```go
package model

type CategoryDomain string

const (
	CategoryDomainIncome   CategoryDomain = "income"
	CategoryDomainExpense  CategoryDomain = "expense"
	CategoryDomainWishlist CategoryDomain = "wishlist"
)

type Category struct {
	Base
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Domain    CategoryDomain `gorm:"type:varchar(50);not null" json:"domain"`
	Color     *string        `gorm:"type:varchar(7)" json:"color,omitempty"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
}

func (Category) TableName() string {
	return "categories"
}
```

- [ ] **Step 3: Create payment method model**

```go
package model

type PaymentMethodType string

const (
	PaymentMethodCash          PaymentMethodType = "cash"
	PaymentMethodDebitCard     PaymentMethodType = "debit_card"
	PaymentMethodCreditCard    PaymentMethodType = "credit_card"
	PaymentMethodDigitalWallet PaymentMethodType = "digital_wallet"
	PaymentMethodCrypto        PaymentMethodType = "crypto"
)

type PaymentMethod struct {
	Base
	Name    string            `gorm:"type:varchar(255);not null" json:"name"`
	Type    PaymentMethodType `gorm:"type:varchar(50);not null" json:"type"`
	Details *string           `gorm:"type:varchar(255)" json:"details,omitempty"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
```

- [ ] **Step 4: Create income model**

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Income struct {
	Base
	Source     string          `gorm:"type:varchar(255);not null" json:"source"`
	Amount    decimal.Decimal  `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency  string           `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category   *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Date      time.Time        `gorm:"type:date;not null" json:"date"`
	Year      int              `gorm:"not null;index" json:"year"`
}

func (Income) TableName() string {
	return "incomes"
}
```

> Note: Install shopspring/decimal for precise money handling:
> `go get github.com/shopspring/decimal@latest`

- [ ] **Step 5: Verify models compile**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go build ./...`
Expected: No errors.

- [ ] **Step 6: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/model/
git commit -m "feat: GORM models for category, payment_method, and income"
```

---

### Task 5: Auth Middleware

**Files:**
- Create: `backend/internal/middleware/auth.go`
- Create: `backend/internal/middleware/auth_test.go`

- [ ] **Step 1: Write failing test for auth middleware**

```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestJWT(secret string, userID string, expired bool) string {
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp.Unix(),
		"aud": "authenticated",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret"
	userID := "550e8400-e29b-41d4-a716-446655440000"
	token := generateTestJWT(secret, userID, false)

	app := fiber.New()
	app.Use(NewAuthMiddleware(secret))
	app.Get("/test", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": uid})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	app := fiber.New()
	app.Use(NewAuthMiddleware("test-secret"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	token := generateTestJWT(secret, "some-user", true)

	app := fiber.New()
	app.Use(NewAuthMiddleware(secret))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/middleware/ -v`
Expected: FAIL — `NewAuthMiddleware` undefined.

- [ ] **Step 3: Implement auth middleware**

```go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing sub claim",
			})
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user id",
			})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}

// GetUserID extracts the authenticated user's UUID from the Fiber context.
func GetUserID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("user_id").(uuid.UUID)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/middleware/ -v`
Expected: All 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/middleware/
git commit -m "feat: JWT auth middleware with Supabase token validation"
```

---

## Phase B: Backend API

### Task 6: Shared Response Helpers

**Files:**
- Create: `backend/internal/handler/response.go`

- [ ] **Step 1: Create response helpers**

```go
package handler

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func respondError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ErrorResponse{Error: msg})
}

func respondList[T any](c *fiber.Ctx, data []T) error {
	if data == nil {
		data = []T{}
	}
	return c.JSON(ListResponse[T]{Data: data, Total: len(data)})
}

func respondJSON(c *fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

func respondCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func respondNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
```

- [ ] **Step 2: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/handler/response.go
git commit -m "feat: shared API response helpers"
```

---

### Task 7: Categories CRUD API

**Files:**
- Create: `backend/internal/repository/category.go`
- Create: `backend/internal/repository/category_test.go`
- Create: `backend/internal/service/category.go`
- Create: `backend/internal/service/category_test.go`
- Create: `backend/internal/handler/category.go`
- Create: `backend/internal/handler/category_test.go`
- Create: `backend/testutil/testutil.go`

- [ ] **Step 1: Create test utilities**

```go
package testutil

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://finance:finance_dev@localhost:5432/finance_tracker?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db
}

func CleanTable(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	db.Exec("DELETE FROM " + table)
}
```

- [ ] **Step 2: Write failing test for category repository**

```go
package repository

import (
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")
	repo := NewCategoryRepository(db)

	userID := uuid.New()
	cat := &model.Category{
		Base:   model.Base{UserID: userID},
		Name:   "Groceries",
		Domain: model.CategoryDomainExpense,
	}

	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cat.ID)
}

func TestCategoryRepository_ListByDomain(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")
	repo := NewCategoryRepository(db)

	userID := uuid.New()
	repo.Create(&model.Category{Base: model.Base{UserID: userID}, Name: "Salary", Domain: model.CategoryDomainIncome})
	repo.Create(&model.Category{Base: model.Base{UserID: userID}, Name: "Groceries", Domain: model.CategoryDomainExpense})

	incomeCategories, err := repo.ListByDomain(userID, model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, incomeCategories, 1)
	assert.Equal(t, "Salary", incomeCategories[0].Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")
	repo := NewCategoryRepository(db)

	userID := uuid.New()
	cat := &model.Category{Base: model.Base{UserID: userID}, Name: "Old Name", Domain: model.CategoryDomainExpense}
	repo.Create(cat)

	cat.Name = "New Name"
	err := repo.Update(cat)
	require.NoError(t, err)

	found, err := repo.GetByID(userID, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", found.Name)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")
	repo := NewCategoryRepository(db)

	userID := uuid.New()
	cat := &model.Category{Base: model.Base{UserID: userID}, Name: "ToDelete", Domain: model.CategoryDomainExpense}
	repo.Create(cat)

	err := repo.Delete(userID, cat.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(userID, cat.ID)
	assert.Error(t, err)
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestCategory`
Expected: FAIL — `NewCategoryRepository` undefined.

- [ ] **Step 4: Implement category repository**

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

func (r *CategoryRepository) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("user_id = ? AND domain = ?", userID, domain).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*model.Category, error) {
	var cat model.Category
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&cat).Error
	return &cat, err
}

func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
```

- [ ] **Step 5: Run repository tests to verify they pass**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestCategory`
Expected: All 4 tests PASS.

- [ ] **Step 6: Implement category service**

```go
package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID uuid.UUID, name string, domain model.CategoryDomain, color *string) (*model.Category, error) {
	cat := &model.Category{
		Base:   model.Base{UserID: userID},
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

func (s *CategoryService) GetByID(userID uuid.UUID, id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(userID, id)
}

func (s *CategoryService) Update(userID uuid.UUID, id uuid.UUID, name string, color *string) (*model.Category, error) {
	cat, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	cat.Name = name
	cat.Color = color
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}

// SeedDefaults creates the default categories for a new user.
func (s *CategoryService) SeedDefaults(userID uuid.UUID) error {
	incomeDefaults := []string{"Salary", "Bonus", "Freelance", "Dividends", "Interest", "Side Hustle"}
	expenseDefaults := []string{
		"Home Expenses", "Eating Out", "Self Care", "Coffee/Drink",
		"Entertainment", "Transportation", "Groceries", "Utilities",
		"Clothes", "Other", "Card Payments", "Savings/Investment",
		"Alcohol", "Drugs", "Taxes", "Knowledge", "Tech",
	}

	for i, name := range incomeDefaults {
		cat := &model.Category{
			Base:      model.Base{UserID: userID},
			Name:      name,
			Domain:    model.CategoryDomainIncome,
			SortOrder: i,
		}
		if err := s.repo.Create(cat); err != nil {
			return err
		}
	}
	for i, name := range expenseDefaults {
		cat := &model.Category{
			Base:      model.Base{UserID: userID},
			Name:      name,
			Domain:    model.CategoryDomainExpense,
			SortOrder: i,
		}
		if err := s.repo.Create(cat); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 7: Implement category handler**

```go
package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

type createCategoryRequest struct {
	Name   string               `json:"name"`
	Domain model.CategoryDomain `json:"domain"`
	Color  *string              `json:"color,omitempty"`
}

type updateCategoryRequest struct {
	Name  string  `json:"name"`
	Color *string `json:"color,omitempty"`
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req createCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}
	if req.Domain != model.CategoryDomainIncome && req.Domain != model.CategoryDomainExpense && req.Domain != model.CategoryDomainWishlist {
		return respondError(c, fiber.StatusBadRequest, "domain must be one of: income, expense, wishlist")
	}

	cat, err := h.svc.Create(userID, req.Name, req.Domain, req.Color)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create category")
	}
	return respondCreated(c, cat)
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	domain := model.CategoryDomain(c.Query("domain"))
	if domain == "" {
		return respondError(c, fiber.StatusBadRequest, "domain query param is required")
	}

	categories, err := h.svc.ListByDomain(userID, domain)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list categories")
	}
	return respondList(c, categories)
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	cat, err := h.svc.GetByID(userID, id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "category not found")
	}
	return respondJSON(c, cat)
}

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
		return respondError(c, fiber.StatusNotFound, "category not found")
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
		return respondError(c, fiber.StatusNotFound, "category not found")
	}
	return respondNoContent(c)
}

func (h *CategoryHandler) SeedDefaults(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if err := h.svc.SeedDefaults(userID); err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to seed defaults")
	}
	return respondJSON(c, fiber.Map{"message": "defaults seeded"})
}
```

- [ ] **Step 8: Write handler test**

```go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeAuth injects a user_id into context for testing without JWT.
func fakeAuth(userID uuid.UUID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func setupCategoryApp(t *testing.T, userID uuid.UUID) (*fiber.App, *CategoryHandler) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	h := NewCategoryHandler(svc)

	app := fiber.New()
	api := app.Group("/api/v1", fakeAuth(userID))
	api.Post("/categories", h.Create)
	api.Get("/categories", h.List)
	api.Get("/categories/:id", h.GetByID)
	api.Put("/categories/:id", h.Update)
	api.Delete("/categories/:id", h.Delete)
	api.Post("/categories/seed", h.SeedDefaults)
	return app, h
}

func TestCategoryHandler_CreateAndList(t *testing.T) {
	userID := uuid.New()
	app, _ := setupCategoryApp(t, userID)

	// Create
	body, _ := json.Marshal(createCategoryRequest{Name: "Groceries", Domain: "expense"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// List
	req = httptest.NewRequest(http.MethodGet, "/api/v1/categories?domain=expense", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp ListResponse[map[string]interface{}]
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, 1, listResp.Total)
	assert.Equal(t, "Groceries", listResp.Data[0]["name"])
}
```

- [ ] **Step 9: Run handler tests**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/handler/ -v -run TestCategory`
Expected: PASS.

- [ ] **Step 10: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/repository/category.go backend/internal/repository/category_test.go backend/internal/service/category.go backend/internal/handler/category.go backend/internal/handler/category_test.go backend/testutil/
git commit -m "feat: categories CRUD API (repository, service, handler with tests)"
```

---

### Task 8: Payment Methods CRUD API

**Files:**
- Create: `backend/internal/repository/payment_method.go`
- Create: `backend/internal/repository/payment_method_test.go`
- Create: `backend/internal/service/payment_method.go`
- Create: `backend/internal/handler/payment_method.go`
- Create: `backend/internal/handler/payment_method_test.go`

- [ ] **Step 1: Write failing test for payment method repository**

```go
package repository

import (
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentMethodRepository_CRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")
	repo := NewPaymentMethodRepository(db)

	userID := uuid.New()
	details := "****1234"
	pm := &model.PaymentMethod{
		Base:    model.Base{UserID: userID},
		Name:    "BBVA Debit",
		Type:    model.PaymentMethodDebitCard,
		Details: &details,
	}

	// Create
	err := repo.Create(pm)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, pm.ID)

	// List
	list, err := repo.ListByUser(userID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	pm.Name = "BBVA Platinum"
	err = repo.Update(pm)
	require.NoError(t, err)

	found, _ := repo.GetByID(userID, pm.ID)
	assert.Equal(t, "BBVA Platinum", found.Name)

	// Delete
	err = repo.Delete(userID, pm.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(userID, pm.ID)
	assert.Error(t, err)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestPaymentMethod`
Expected: FAIL — `NewPaymentMethodRepository` undefined.

- [ ] **Step 3: Implement payment method repository**

```go
package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) Create(pm *model.PaymentMethod) error {
	return r.db.Create(pm).Error
}

func (r *PaymentMethodRepository) ListByUser(userID uuid.UUID) ([]model.PaymentMethod, error) {
	var pms []model.PaymentMethod
	err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&pms).Error
	return pms, err
}

func (r *PaymentMethodRepository) ListByType(userID uuid.UUID, pmType model.PaymentMethodType) ([]model.PaymentMethod, error) {
	var pms []model.PaymentMethod
	err := r.db.Where("user_id = ? AND type = ?", userID, pmType).Order("name ASC").Find(&pms).Error
	return pms, err
}

func (r *PaymentMethodRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*model.PaymentMethod, error) {
	var pm model.PaymentMethod
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&pm).Error
	return &pm, err
}

func (r *PaymentMethodRepository) Update(pm *model.PaymentMethod) error {
	return r.db.Save(pm).Error
}

func (r *PaymentMethodRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.PaymentMethod{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
```

- [ ] **Step 4: Run repository tests**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestPaymentMethod`
Expected: PASS.

- [ ] **Step 5: Implement payment method service**

```go
package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

type PaymentMethodService struct {
	repo *repository.PaymentMethodRepository
}

func NewPaymentMethodService(repo *repository.PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{repo: repo}
}

func (s *PaymentMethodService) Create(userID uuid.UUID, name string, pmType model.PaymentMethodType, details *string) (*model.PaymentMethod, error) {
	pm := &model.PaymentMethod{
		Base:    model.Base{UserID: userID},
		Name:    name,
		Type:    pmType,
		Details: details,
	}
	if err := s.repo.Create(pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *PaymentMethodService) ListByUser(userID uuid.UUID) ([]model.PaymentMethod, error) {
	return s.repo.ListByUser(userID)
}

func (s *PaymentMethodService) ListByType(userID uuid.UUID, pmType model.PaymentMethodType) ([]model.PaymentMethod, error) {
	return s.repo.ListByType(userID, pmType)
}

func (s *PaymentMethodService) GetByID(userID uuid.UUID, id uuid.UUID) (*model.PaymentMethod, error) {
	return s.repo.GetByID(userID, id)
}

func (s *PaymentMethodService) Update(userID uuid.UUID, id uuid.UUID, name string, pmType model.PaymentMethodType, details *string) (*model.PaymentMethod, error) {
	pm, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	pm.Name = name
	pm.Type = pmType
	pm.Details = details
	if err := s.repo.Update(pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (s *PaymentMethodService) Delete(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
```

- [ ] **Step 6: Implement payment method handler**

```go
package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentMethodHandler struct {
	svc *service.PaymentMethodService
}

func NewPaymentMethodHandler(svc *service.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{svc: svc}
}

type createPaymentMethodRequest struct {
	Name    string                  `json:"name"`
	Type    model.PaymentMethodType `json:"type"`
	Details *string                 `json:"details,omitempty"`
}

type updatePaymentMethodRequest struct {
	Name    string                  `json:"name"`
	Type    model.PaymentMethodType `json:"type"`
	Details *string                 `json:"details,omitempty"`
}

func (h *PaymentMethodHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req createPaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" || req.Type == "" {
		return respondError(c, fiber.StatusBadRequest, "name and type are required")
	}

	pm, err := h.svc.Create(userID, req.Name, req.Type, req.Details)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create payment method")
	}
	return respondCreated(c, pm)
}

func (h *PaymentMethodHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pmType := model.PaymentMethodType(c.Query("type"))

	var pms []model.PaymentMethod
	var err error
	if pmType != "" {
		pms, err = h.svc.ListByType(userID, pmType)
	} else {
		pms, err = h.svc.ListByUser(userID)
	}
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list payment methods")
	}
	return respondList(c, pms)
}

func (h *PaymentMethodHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	pm, err := h.svc.GetByID(userID, id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "payment method not found")
	}
	return respondJSON(c, pm)
}

func (h *PaymentMethodHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updatePaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Name == "" || req.Type == "" {
		return respondError(c, fiber.StatusBadRequest, "name and type are required")
	}

	pm, err := h.svc.Update(userID, id, req.Name, req.Type, req.Details)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "payment method not found")
	}
	return respondJSON(c, pm)
}

func (h *PaymentMethodHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusNotFound, "payment method not found")
	}
	return respondNoContent(c)
}
```

- [ ] **Step 7: Write handler test**

```go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPaymentMethodApp(t *testing.T, userID uuid.UUID) *fiber.App {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")
	repo := repository.NewPaymentMethodRepository(db)
	svc := service.NewPaymentMethodService(repo)
	h := NewPaymentMethodHandler(svc)

	app := fiber.New()
	api := app.Group("/api/v1", fakeAuth(userID))
	api.Post("/payment-methods", h.Create)
	api.Get("/payment-methods", h.List)
	api.Get("/payment-methods/:id", h.GetByID)
	api.Put("/payment-methods/:id", h.Update)
	api.Delete("/payment-methods/:id", h.Delete)
	return app
}

func TestPaymentMethodHandler_CreateAndList(t *testing.T) {
	userID := uuid.New()
	app := setupPaymentMethodApp(t, userID)

	body, _ := json.Marshal(createPaymentMethodRequest{Name: "Nu Credit", Type: "credit_card"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment-methods", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/payment-methods", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp ListResponse[map[string]interface{}]
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, 1, listResp.Total)
	assert.Equal(t, "Nu Credit", listResp.Data[0]["name"])
}
```

- [ ] **Step 8: Run all tests**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./... -v`
Expected: All tests PASS.

- [ ] **Step 9: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/repository/payment_method.go backend/internal/repository/payment_method_test.go backend/internal/service/payment_method.go backend/internal/handler/payment_method.go backend/internal/handler/payment_method_test.go
git commit -m "feat: payment methods CRUD API (repository, service, handler with tests)"
```

---

### Task 9: Income CRUD API

**Files:**
- Create: `backend/internal/repository/income.go`
- Create: `backend/internal/repository/income_test.go`
- Create: `backend/internal/service/income.go`
- Create: `backend/internal/handler/income.go`
- Create: `backend/internal/handler/income_test.go`

- [ ] **Step 1: Write failing test for income repository**

```go
package repository

import (
	"testing"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncomeRepository_CRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")
	repo := NewIncomeRepository(db)
	catRepo := NewCategoryRepository(db)

	userID := uuid.New()

	// Create a category for the income
	cat := &model.Category{Base: model.Base{UserID: userID}, Name: "Salary", Domain: model.CategoryDomainIncome}
	catRepo.Create(cat)

	income := &model.Income{
		Base:       model.Base{UserID: userID},
		Source:     "Company ABC",
		Amount:     decimal.NewFromFloat(25000.50),
		Currency:   "MXN",
		CategoryID: &cat.ID,
		Date:       time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		Year:       2026,
	}

	// Create
	err := repo.Create(income)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, income.ID)

	// List by year
	list, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Company ABC", list[0].Source)
	assert.NotNil(t, list[0].Category) // Preloaded

	// Get by ID
	found, err := repo.GetByID(userID, income.ID)
	require.NoError(t, err)
	assert.Equal(t, "Company ABC", found.Source)
	assert.True(t, found.Amount.Equal(decimal.NewFromFloat(25000.50)))

	// Update
	income.Source = "Company XYZ"
	err = repo.Update(income)
	require.NoError(t, err)
	found, _ = repo.GetByID(userID, income.ID)
	assert.Equal(t, "Company XYZ", found.Source)

	// Delete
	err = repo.Delete(userID, income.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(userID, income.ID)
	assert.Error(t, err)
}

func TestIncomeRepository_ListByYear_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	repo := NewIncomeRepository(db)

	list, err := repo.ListByYear(uuid.New(), 2026)
	require.NoError(t, err)
	assert.Empty(t, list)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestIncome`
Expected: FAIL — `NewIncomeRepository` undefined.

- [ ] **Step 3: Implement income repository**

```go
package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) Create(income *model.Income) error {
	return r.db.Create(income).Error
}

func (r *IncomeRepository) ListByYear(userID uuid.UUID, year int) ([]model.Income, error) {
	var incomes []model.Income
	err := r.db.Preload("Category").
		Where("user_id = ? AND year = ?", userID, year).
		Order("date DESC").
		Find(&incomes).Error
	return incomes, err
}

func (r *IncomeRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*model.Income, error) {
	var income model.Income
	err := r.db.Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&income).Error
	return &income, err
}

func (r *IncomeRepository) Update(income *model.Income) error {
	return r.db.Save(income).Error
}

func (r *IncomeRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Income{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
```

- [ ] **Step 4: Run repository tests**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/repository/ -v -run TestIncome`
Expected: All tests PASS.

- [ ] **Step 5: Implement income service**

```go
package service

import (
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IncomeService struct {
	repo *repository.IncomeRepository
}

func NewIncomeService(repo *repository.IncomeRepository) *IncomeService {
	return &IncomeService{repo: repo}
}

func (s *IncomeService) Create(userID uuid.UUID, source string, amount decimal.Decimal, currency string, categoryID *uuid.UUID, date time.Time) (*model.Income, error) {
	income := &model.Income{
		Base:       model.Base{UserID: userID},
		Source:     source,
		Amount:     amount,
		Currency:   currency,
		CategoryID: categoryID,
		Date:       date,
		Year:       date.Year(),
	}
	if err := s.repo.Create(income); err != nil {
		return nil, err
	}
	// Re-fetch to preload Category
	return s.repo.GetByID(userID, income.ID)
}

func (s *IncomeService) ListByYear(userID uuid.UUID, year int) ([]model.Income, error) {
	return s.repo.ListByYear(userID, year)
}

func (s *IncomeService) GetByID(userID uuid.UUID, id uuid.UUID) (*model.Income, error) {
	return s.repo.GetByID(userID, id)
}

func (s *IncomeService) Update(userID uuid.UUID, id uuid.UUID, source string, amount decimal.Decimal, currency string, categoryID *uuid.UUID, date time.Time) (*model.Income, error) {
	income, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	income.Source = source
	income.Amount = amount
	income.Currency = currency
	income.CategoryID = categoryID
	income.Date = date
	income.Year = date.Year()
	if err := s.repo.Update(income); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, income.ID)
}

func (s *IncomeService) Delete(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
```

- [ ] **Step 6: Implement income handler**

```go
package handler

import (
	"strconv"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IncomeHandler struct {
	svc *service.IncomeService
}

func NewIncomeHandler(svc *service.IncomeService) *IncomeHandler {
	return &IncomeHandler{svc: svc}
}

type createIncomeRequest struct {
	Source     string  `json:"source"`
	Amount    string  `json:"amount"`
	Currency  string  `json:"currency"`
	CategoryID *string `json:"category_id,omitempty"`
	Date      string  `json:"date"` // YYYY-MM-DD
}

type updateIncomeRequest struct {
	Source     string  `json:"source"`
	Amount    string  `json:"amount"`
	Currency  string  `json:"currency"`
	CategoryID *string `json:"category_id,omitempty"`
	Date      string  `json:"date"`
}

func (h *IncomeHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req createIncomeRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Source == "" || req.Amount == "" || req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "source, amount, and date are required")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid amount format")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	currency := req.Currency
	if currency == "" {
		currency = "MXN"
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		parsed, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid category_id")
		}
		categoryID = &parsed
	}

	income, err := h.svc.Create(userID, req.Source, amount, currency, categoryID, date)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create income")
	}
	return respondCreated(c, income)
}

func (h *IncomeHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	yearStr := c.Params("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid year")
	}

	incomes, err := h.svc.ListByYear(userID, year)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to list incomes")
	}
	return respondList(c, incomes)
}

func (h *IncomeHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	income, err := h.svc.GetByID(userID, id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "income not found")
	}
	return respondJSON(c, income)
}

func (h *IncomeHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	var req updateIncomeRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if req.Source == "" || req.Amount == "" || req.Date == "" {
		return respondError(c, fiber.StatusBadRequest, "source, amount, and date are required")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid amount format")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	currency := req.Currency
	if currency == "" {
		currency = "MXN"
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		parsed, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid category_id")
		}
		categoryID = &parsed
	}

	income, err := h.svc.Update(userID, id, req.Source, amount, currency, categoryID, date)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "income not found")
	}
	return respondJSON(c, income)
}

func (h *IncomeHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(userID, id); err != nil {
		return respondError(c, fiber.StatusNotFound, "income not found")
	}
	return respondNoContent(c)
}
```

- [ ] **Step 7: Write handler test**

```go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIncomeApp(t *testing.T, userID uuid.UUID) *fiber.App {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")
	repo := repository.NewIncomeRepository(db)
	svc := service.NewIncomeService(repo)
	h := NewIncomeHandler(svc)

	app := fiber.New()
	api := app.Group("/api/v1", fakeAuth(userID))
	api.Post("/years/:year/incomes", h.Create)
	api.Get("/years/:year/incomes", h.List)
	api.Get("/incomes/:id", h.GetByID)
	api.Put("/incomes/:id", h.Update)
	api.Delete("/incomes/:id", h.Delete)
	return app
}

func TestIncomeHandler_CreateAndList(t *testing.T) {
	userID := uuid.New()
	app := setupIncomeApp(t, userID)

	body, _ := json.Marshal(createIncomeRequest{
		Source:   "Company ABC",
		Amount:  "25000.50",
		Currency: "MXN",
		Date:    "2026-03-15",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/years/2026/incomes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, "Company ABC", created["source"])
	assert.Equal(t, "25000.50", created["amount"])

	// List
	req = httptest.NewRequest(http.MethodGet, "/api/v1/years/2026/incomes", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp ListResponse[map[string]interface{}]
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, 1, listResp.Total)
}

func TestIncomeHandler_CreateValidation(t *testing.T) {
	userID := uuid.New()
	app := setupIncomeApp(t, userID)

	// Missing required fields
	body, _ := json.Marshal(createIncomeRequest{Source: "", Amount: "", Date: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/years/2026/incomes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
```

- [ ] **Step 8: Run all tests**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./... -v`
Expected: All tests PASS.

- [ ] **Step 9: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/repository/income.go backend/internal/repository/income_test.go backend/internal/service/income.go backend/internal/handler/income.go backend/internal/handler/income_test.go
git commit -m "feat: income CRUD API (repository, service, handler with tests)"
```

---

### Task 10: Router and Wired Main

**Files:**
- Create: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Create router with all routes**

```go
package router

import (
	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	jwtSecret string,
	categoryHandler *handler.CategoryHandler,
	paymentMethodHandler *handler.PaymentMethodHandler,
	incomeHandler *handler.IncomeHandler,
) {
	// Health check (no auth)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Authenticated routes
	api := app.Group("/api/v1", middleware.NewAuthMiddleware(jwtSecret))

	// Categories
	categories := api.Group("/categories")
	categories.Post("/", categoryHandler.Create)
	categories.Get("/", categoryHandler.List)
	categories.Post("/seed", categoryHandler.SeedDefaults)
	categories.Get("/:id", categoryHandler.GetByID)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)

	// Payment Methods
	paymentMethods := api.Group("/payment-methods")
	paymentMethods.Post("/", paymentMethodHandler.Create)
	paymentMethods.Get("/", paymentMethodHandler.List)
	paymentMethods.Get("/:id", paymentMethodHandler.GetByID)
	paymentMethods.Put("/:id", paymentMethodHandler.Update)
	paymentMethods.Delete("/:id", paymentMethodHandler.Delete)

	// Income (year-scoped)
	api.Post("/years/:year/incomes", incomeHandler.Create)
	api.Get("/years/:year/incomes", incomeHandler.List)
	api.Get("/incomes/:id", incomeHandler.GetByID)
	api.Put("/incomes/:id", incomeHandler.Update)
	api.Delete("/incomes/:id", incomeHandler.Delete)
}
```

- [ ] **Step 2: Update main.go to wire everything together**

```go
package main

import (
	"log"
	"os"

	"github.com/folkrom/finance-tracker/backend/internal/config"
	"github.com/folkrom/finance-tracker/backend/internal/database"
	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/router"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	if os.Getenv("ENVIRONMENT") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.New(cfg, logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)

	// Services
	categorySvc := service.NewCategoryService(categoryRepo)
	paymentMethodSvc := service.NewPaymentMethodService(paymentMethodRepo)
	incomeSvc := service.NewIncomeService(incomeRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodSvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)

	// Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	router.Setup(app, cfg.SupabaseJWTSecret, categoryHandler, paymentMethodHandler, incomeHandler)

	logger.Info("server starting", zap.String("port", cfg.Port))
	log.Fatal(app.Listen(":" + cfg.Port))
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go build ./cmd/server`
Expected: No errors.

- [ ] **Step 4: Run all tests one final time**

Run: `cd /home/folkrom/projects/finance-tracker/backend && go test ./... -v -count=1`
Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add backend/internal/router/ backend/cmd/server/main.go
git commit -m "feat: router setup and wired main.go with all Plan 1 endpoints"
```

---

## Phase C: Frontend Foundation

### Task 11: Next.js Scaffolding

**Files:**
- Create: `frontend/` (via create-next-app)
- Modify: `frontend/next.config.ts`
- Create: `frontend/.env.local.example`
- Create: `frontend/src/types/index.ts`

- [ ] **Step 1: Create Next.js project**

```bash
cd /home/folkrom/projects/finance-tracker
npx create-next-app@latest frontend --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --no-turbopack
```

- [ ] **Step 2: Install dependencies**

```bash
cd /home/folkrom/projects/finance-tracker/frontend
npm install @supabase/supabase-js @supabase/ssr next-intl react-hook-form @hookform/resolvers zod
npx shadcn@latest init -d
npx shadcn@latest add button input label select dialog table card badge dropdown-menu sheet separator toast
```

- [ ] **Step 3: Create .env.local.example**

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

- [ ] **Step 4: Create shared TypeScript types**

```typescript
// src/types/index.ts

export interface Category {
  id: string;
  user_id: string;
  name: string;
  domain: "income" | "expense" | "wishlist";
  color?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface PaymentMethod {
  id: string;
  user_id: string;
  name: string;
  type: "cash" | "debit_card" | "credit_card" | "digital_wallet" | "crypto";
  details?: string;
  created_at: string;
  updated_at: string;
}

export interface Income {
  id: string;
  user_id: string;
  source: string;
  amount: string; // decimal as string from API
  currency: string;
  category_id?: string;
  category?: Category;
  date: string;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface ListResponse<T> {
  data: T[];
  total: number;
}

export interface ErrorResponse {
  error: string;
}
```

- [ ] **Step 5: Update next.config.ts for i18n plugin**

```typescript
// next.config.ts
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/lib/i18n/request.ts");

const nextConfig = {};

export default withNextIntl(nextConfig);
```

- [ ] **Step 6: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/
git commit -m "feat: Next.js scaffolding with Tailwind, shadcn/ui, Supabase, and i18n deps"
```

---

### Task 12: Supabase Auth + API Client

**Files:**
- Create: `frontend/src/lib/supabase/client.ts`
- Create: `frontend/src/lib/supabase/server.ts`
- Create: `frontend/src/lib/supabase/middleware.ts`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/middleware.ts`

- [ ] **Step 1: Create browser Supabase client**

```typescript
// src/lib/supabase/client.ts
import { createBrowserClient } from "@supabase/ssr";

export function createClient() {
  return createBrowserClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
  );
}
```

- [ ] **Step 2: Create server Supabase client**

```typescript
// src/lib/supabase/server.ts
import { createServerClient } from "@supabase/ssr";
import { cookies } from "next/headers";

export async function createClient() {
  const cookieStore = await cookies();

  return createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll() {
          return cookieStore.getAll();
        },
        setAll(cookiesToSet) {
          try {
            cookiesToSet.forEach(({ name, value, options }) =>
              cookieStore.set(name, value, options)
            );
          } catch {
            // Server Component — ignore
          }
        },
      },
    }
  );
}
```

- [ ] **Step 3: Create Supabase middleware helper**

```typescript
// src/lib/supabase/middleware.ts
import { createServerClient } from "@supabase/ssr";
import { NextResponse, type NextRequest } from "next/server";

export async function updateSession(request: NextRequest) {
  let supabaseResponse = NextResponse.next({ request });

  const supabase = createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll();
        },
        setAll(cookiesToSet) {
          cookiesToSet.forEach(({ name, value }) =>
            request.cookies.set(name, value)
          );
          supabaseResponse = NextResponse.next({ request });
          cookiesToSet.forEach(({ name, value, options }) =>
            supabaseResponse.cookies.set(name, value, options)
          );
        },
      },
    }
  );

  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (
    !user &&
    !request.nextUrl.pathname.startsWith("/login") &&
    !request.nextUrl.pathname.startsWith("/auth")
  ) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return supabaseResponse;
}
```

- [ ] **Step 4: Create Next.js root middleware**

```typescript
// frontend/middleware.ts
import { type NextRequest } from "next/server";
import { updateSession } from "@/lib/supabase/middleware";

export async function middleware(request: NextRequest) {
  return await updateSession(request);
}

export const config = {
  matcher: [
    "/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
```

- [ ] **Step 5: Create API fetch wrapper**

```typescript
// src/lib/api.ts
import { createClient } from "@/lib/supabase/client";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function getAuthHeaders(): Promise<Record<string, string>> {
  const supabase = createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session?.access_token) {
    throw new Error("Not authenticated");
  }

  return {
    Authorization: `Bearer ${session.access_token}`,
    "Content-Type": "application/json",
  };
}

export async function apiGet<T>(path: string): Promise<T> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, { headers });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "API error");
  }
  return res.json();
}

export async function apiPost<T>(path: string, body: unknown): Promise<T> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "API error");
  }
  return res.json();
}

export async function apiPut<T>(path: string, body: unknown): Promise<T> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, {
    method: "PUT",
    headers,
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "API error");
  }
  return res.json();
}

export async function apiDelete(path: string): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, {
    method: "DELETE",
    headers,
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "API error");
  }
}
```

- [ ] **Step 6: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/lib/ frontend/middleware.ts
git commit -m "feat: Supabase auth integration and API client wrapper"
```

---

### Task 13: i18n Setup

**Files:**
- Create: `frontend/src/lib/i18n/config.ts`
- Create: `frontend/src/lib/i18n/request.ts`
- Create: `frontend/src/lib/i18n/en.json`
- Create: `frontend/src/lib/i18n/es.json`

- [ ] **Step 1: Create i18n config**

```typescript
// src/lib/i18n/config.ts
export const locales = ["en", "es"] as const;
export type Locale = (typeof locales)[number];
export const defaultLocale: Locale = "en";
```

- [ ] **Step 2: Create request config for next-intl**

```typescript
// src/lib/i18n/request.ts
import { getRequestConfig } from "next-intl/server";
import { defaultLocale } from "./config";

export default getRequestConfig(async () => {
  // For now, use cookie or default. Locale switching will read from cookie.
  const locale = defaultLocale;

  return {
    locale,
    messages: (await import(`./${locale}.json`)).default,
  };
});
```

- [ ] **Step 3: Create English translations**

```json
{
  "common": {
    "save": "Save",
    "cancel": "Cancel",
    "delete": "Delete",
    "edit": "Edit",
    "create": "Create",
    "search": "Search",
    "loading": "Loading...",
    "noResults": "No results found",
    "confirm": "Are you sure?",
    "currency": "Currency"
  },
  "nav": {
    "dashboard": "Dashboard",
    "income": "Income",
    "expenses": "Expenses",
    "debt": "Debt",
    "budget": "Budget",
    "cards": "Cards",
    "wishlist": "Wishlist",
    "settings": "Settings"
  },
  "income": {
    "title": "Income",
    "source": "Source",
    "amount": "Amount",
    "category": "Category",
    "date": "Date",
    "addIncome": "Add Income",
    "editIncome": "Edit Income",
    "noIncome": "No income records yet"
  },
  "settings": {
    "title": "Settings",
    "categories": "Categories",
    "paymentMethods": "Payment Methods",
    "addCategory": "Add Category",
    "addPaymentMethod": "Add Payment Method",
    "seedDefaults": "Load Default Categories",
    "language": "Language"
  },
  "auth": {
    "login": "Log In",
    "signup": "Sign Up",
    "email": "Email",
    "password": "Password",
    "logout": "Log Out"
  }
}
```

- [ ] **Step 4: Create Spanish translations**

```json
{
  "common": {
    "save": "Guardar",
    "cancel": "Cancelar",
    "delete": "Eliminar",
    "edit": "Editar",
    "create": "Crear",
    "search": "Buscar",
    "loading": "Cargando...",
    "noResults": "No se encontraron resultados",
    "confirm": "¿Estás seguro?",
    "currency": "Moneda"
  },
  "nav": {
    "dashboard": "Inicio",
    "income": "Ingresos",
    "expenses": "Gastos",
    "debt": "Deudas",
    "budget": "Presupuesto",
    "cards": "Tarjetas",
    "wishlist": "Lista de Deseos",
    "settings": "Configuración"
  },
  "income": {
    "title": "Ingresos",
    "source": "Fuente",
    "amount": "Monto",
    "category": "Categoría",
    "date": "Fecha",
    "addIncome": "Agregar Ingreso",
    "editIncome": "Editar Ingreso",
    "noIncome": "No hay registros de ingresos"
  },
  "settings": {
    "title": "Configuración",
    "categories": "Categorías",
    "paymentMethods": "Métodos de Pago",
    "addCategory": "Agregar Categoría",
    "addPaymentMethod": "Agregar Método de Pago",
    "seedDefaults": "Cargar Categorías Predeterminadas",
    "language": "Idioma"
  },
  "auth": {
    "login": "Iniciar Sesión",
    "signup": "Registrarse",
    "email": "Correo Electrónico",
    "password": "Contraseña",
    "logout": "Cerrar Sesión"
  }
}
```

- [ ] **Step 5: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/lib/i18n/
git commit -m "feat: i18n setup with English and Spanish translations"
```

---

### Task 14: Login Page

**Files:**
- Create: `frontend/src/app/login/page.tsx`

- [ ] **Step 1: Create login page**

```tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTranslations } from "next-intl";

export default function LoginPage() {
  const t = useTranslations("auth");
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSignUp, setIsSignUp] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    const supabase = createClient();

    const { error: authError } = isSignUp
      ? await supabase.auth.signUp({ email, password })
      : await supabase.auth.signInWithPassword({ email, password });

    setLoading(false);

    if (authError) {
      setError(authError.message);
      return;
    }

    const currentYear = new Date().getFullYear();
    router.push(`/${currentYear}/income`);
    router.refresh();
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-center">
            Finance Tracker
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">{t("email")}</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">{t("password")}</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={6}
              />
            </div>
            {error && (
              <p className="text-sm text-red-600">{error}</p>
            )}
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "..." : isSignUp ? t("signup") : t("login")}
            </Button>
            <Button
              type="button"
              variant="ghost"
              className="w-full"
              onClick={() => setIsSignUp(!isSignUp)}
            >
              {isSignUp ? t("login") : t("signup")}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /home/folkrom/projects/finance-tracker/frontend && npm run build`
Expected: Build succeeds (or minor config issues to resolve).

- [ ] **Step 3: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/app/login/
git commit -m "feat: login page with Supabase auth (email/password + signup)"
```

---

## Phase D: Frontend App Shell & Income Module

### Task 15: App Layout with Sidebar

**Files:**
- Create: `frontend/src/components/layout/sidebar.tsx`
- Create: `frontend/src/components/layout/header.tsx`
- Create: `frontend/src/components/layout/year-switcher.tsx`
- Create: `frontend/src/app/[year]/layout.tsx`
- Modify: `frontend/src/app/layout.tsx`
- Modify: `frontend/src/app/page.tsx`

- [ ] **Step 1: Create sidebar**

```tsx
"use client";

import Link from "next/link";
import { useParams, usePathname } from "next/navigation";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";

const navItems = [
  { key: "income", path: "income", icon: "↑" },
  { key: "expenses", path: "expenses", icon: "↓" },
  { key: "debt", path: "debt", icon: "↗" },
  { key: "budget", path: "budget", icon: "◎" },
  { key: "cards", path: "cards", icon: "▭" },
  { key: "settings", path: "settings", icon: "⚙" },
] as const;

export function Sidebar() {
  const t = useTranslations("nav");
  const params = useParams();
  const pathname = usePathname();
  const year = params.year as string;

  return (
    <aside className="w-64 border-r bg-white h-screen sticky top-0 flex flex-col">
      <div className="p-6">
        <h1 className="text-xl font-bold">Finance Tracker</h1>
      </div>
      <nav className="flex-1 px-3">
        {navItems.map((item) => {
          const href = `/${year}/${item.path}`;
          const isActive = pathname.startsWith(href);
          return (
            <Link
              key={item.key}
              href={href}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors mb-1",
                isActive
                  ? "bg-gray-100 text-gray-900 font-medium"
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900"
              )}
            >
              <span className="text-lg">{item.icon}</span>
              {t(item.key)}
            </Link>
          );
        })}
      </nav>
      <div className="p-3 border-t">
        <Link
          href="/wishlist"
          className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-gray-600 hover:bg-gray-50"
        >
          <span className="text-lg">★</span>
          {t("wishlist")}
        </Link>
      </div>
    </aside>
  );
}
```

- [ ] **Step 2: Create header with year switcher**

```tsx
// src/components/layout/year-switcher.tsx
"use client";

import { useParams, useRouter, usePathname } from "next/navigation";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function YearSwitcher() {
  const params = useParams();
  const router = useRouter();
  const pathname = usePathname();
  const currentYear = params.year as string;

  const thisYear = new Date().getFullYear();
  const years = Array.from({ length: 5 }, (_, i) => thisYear - 2 + i);

  function handleChange(year: string) {
    const newPath = pathname.replace(`/${currentYear}/`, `/${year}/`);
    router.push(newPath);
  }

  return (
    <Select value={currentYear} onValueChange={handleChange}>
      <SelectTrigger className="w-28">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {years.map((y) => (
          <SelectItem key={y} value={String(y)}>
            {y}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
```

```tsx
// src/components/layout/header.tsx
"use client";

import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { YearSwitcher } from "./year-switcher";
import { useTranslations } from "next-intl";

export function Header() {
  const t = useTranslations("auth");
  const router = useRouter();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <YearSwitcher />
      <Button variant="ghost" size="sm" onClick={handleLogout}>
        {t("logout")}
      </Button>
    </header>
  );
}
```

- [ ] **Step 3: Create year-scoped layout**

```tsx
// src/app/[year]/layout.tsx
import { Sidebar } from "@/components/layout/sidebar";
import { Header } from "@/components/layout/header";

export default function YearLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Header />
        <main className="flex-1 p-6 bg-gray-50">{children}</main>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: Update root layout**

```tsx
// src/app/layout.tsx
import type { Metadata } from "next";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import "./globals.css";

export const metadata: Metadata = {
  title: "Finance Tracker",
  description: "Personal finance tracker",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html lang={locale}>
      <body>
        <NextIntlClientProvider messages={messages}>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
```

- [ ] **Step 5: Create root page redirect**

```tsx
// src/app/page.tsx
import { redirect } from "next/navigation";

export default function Home() {
  const currentYear = new Date().getFullYear();
  redirect(`/${currentYear}/income`);
}
```

- [ ] **Step 6: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/app/ frontend/src/components/layout/
git commit -m "feat: app shell with sidebar, header, year switcher, and year-scoped layout"
```

---

### Task 16: Income Module UI

**Files:**
- Create: `frontend/src/app/[year]/income/page.tsx`
- Create: `frontend/src/components/income/income-table.tsx`
- Create: `frontend/src/components/income/income-form.tsx`

- [ ] **Step 1: Create income table component**

```tsx
// src/components/income/income-table.tsx
"use client";

import { useTranslations } from "next-intl";
import { Income } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface IncomeTableProps {
  incomes: Income[];
  onEdit: (income: Income) => void;
  onDelete: (id: string) => void;
}

export function IncomeTable({ incomes, onEdit, onDelete }: IncomeTableProps) {
  const t = useTranslations("income");
  const tc = useTranslations("common");

  if (incomes.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        {t("noIncome")}
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{t("source")}</TableHead>
          <TableHead>{t("amount")}</TableHead>
          <TableHead>{t("category")}</TableHead>
          <TableHead>{t("date")}</TableHead>
          <TableHead className="w-24"></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {incomes.map((income) => (
          <TableRow key={income.id}>
            <TableCell className="font-medium">{income.source}</TableCell>
            <TableCell>
              ${Number(income.amount).toLocaleString()} {income.currency}
            </TableCell>
            <TableCell>
              {income.category ? (
                <Badge variant="secondary">{income.category.name}</Badge>
              ) : (
                "—"
              )}
            </TableCell>
            <TableCell>{new Date(income.date).toLocaleDateString()}</TableCell>
            <TableCell>
              <div className="flex gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onEdit(income)}
                >
                  {tc("edit")}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-red-600"
                  onClick={() => onDelete(income.id)}
                >
                  {tc("delete")}
                </Button>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
```

- [ ] **Step 2: Create income form component**

```tsx
// src/components/income/income-form.tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslations } from "next-intl";
import { Category, Income } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

const incomeSchema = z.object({
  source: z.string().min(1, "Source is required"),
  amount: z.string().min(1, "Amount is required"),
  currency: z.string().default("MXN"),
  category_id: z.string().optional(),
  date: z.string().min(1, "Date is required"),
});

type IncomeFormData = z.infer<typeof incomeSchema>;

interface IncomeFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: IncomeFormData) => Promise<void>;
  categories: Category[];
  defaultValues?: Income;
}

export function IncomeForm({
  open,
  onClose,
  onSubmit,
  categories,
  defaultValues,
}: IncomeFormProps) {
  const t = useTranslations("income");
  const tc = useTranslations("common");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<IncomeFormData>({
    resolver: zodResolver(incomeSchema),
    defaultValues: defaultValues
      ? {
          source: defaultValues.source,
          amount: defaultValues.amount,
          currency: defaultValues.currency,
          category_id: defaultValues.category_id || undefined,
          date: defaultValues.date.split("T")[0],
        }
      : { currency: "MXN" },
  });

  async function handleFormSubmit(data: IncomeFormData) {
    await onSubmit(data);
    reset();
    onClose();
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editIncome") : t("addIncome")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label>{t("source")}</Label>
            <Input {...register("source")} />
            {errors.source && (
              <p className="text-sm text-red-600">{errors.source.message}</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>{t("amount")}</Label>
              <Input {...register("amount")} type="number" step="0.01" />
              {errors.amount && (
                <p className="text-sm text-red-600">{errors.amount.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label>{tc("currency")}</Label>
              <Select
                value={watch("currency")}
                onValueChange={(v) => setValue("currency", v)}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="MXN">MXN</SelectItem>
                  <SelectItem value="USD">USD</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label>{t("category")}</Label>
            <Select
              value={watch("category_id") || "none"}
              onValueChange={(v) =>
                setValue("category_id", v === "none" ? undefined : v)
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">—</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>{t("date")}</Label>
            <Input {...register("date")} type="date" />
            {errors.date && (
              <p className="text-sm text-red-600">{errors.date.message}</p>
            )}
          </div>

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose}>
              {tc("cancel")}
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? tc("loading") : tc("save")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 3: Create income page**

```tsx
// src/app/[year]/income/page.tsx
"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Income, Category, ListResponse } from "@/types";
import { IncomeTable } from "@/components/income/income-table";
import { IncomeForm } from "@/components/income/income-form";
import { Button } from "@/components/ui/button";

export default function IncomePage() {
  const t = useTranslations("income");
  const params = useParams();
  const year = params.year as string;

  const [incomes, setIncomes] = useState<Income[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Income | undefined>();
  const [loading, setLoading] = useState(true);

  async function loadData() {
    setLoading(true);
    try {
      const [incomeRes, categoryRes] = await Promise.all([
        apiGet<ListResponse<Income>>(`/api/v1/years/${year}/incomes`),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=income"),
      ]);
      setIncomes(incomeRes.data);
      setCategories(categoryRes.data);
    } catch (err) {
      console.error("Failed to load data:", err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, [year]);

  async function handleCreate(data: {
    source: string;
    amount: string;
    currency: string;
    category_id?: string;
    date: string;
  }) {
    await apiPost(`/api/v1/years/${year}/incomes`, data);
    await loadData();
  }

  async function handleUpdate(data: {
    source: string;
    amount: string;
    currency: string;
    category_id?: string;
    date: string;
  }) {
    if (!editing) return;
    await apiPut(`/api/v1/incomes/${editing.id}`, data);
    setEditing(undefined);
    await loadData();
  }

  async function handleDelete(id: string) {
    await apiDelete(`/api/v1/incomes/${id}`);
    await loadData();
  }

  function handleEdit(income: Income) {
    setEditing(income);
    setFormOpen(true);
  }

  function handleCloseForm() {
    setFormOpen(false);
    setEditing(undefined);
  }

  if (loading) {
    return <div className="p-6">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">{t("title")}</h2>
        <Button onClick={() => setFormOpen(true)}>{t("addIncome")}</Button>
      </div>

      <IncomeTable
        incomes={incomes}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />

      <IncomeForm
        open={formOpen}
        onClose={handleCloseForm}
        onSubmit={editing ? handleUpdate : handleCreate}
        categories={categories}
        defaultValues={editing}
      />
    </div>
  );
}
```

- [ ] **Step 4: Verify frontend compiles**

Run: `cd /home/folkrom/projects/finance-tracker/frontend && npm run build`
Expected: Build succeeds.

- [ ] **Step 5: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/app/[year]/income/ frontend/src/components/income/
git commit -m "feat: income module UI with table, form, and CRUD operations"
```

---

### Task 17: Settings Page (Categories + Payment Methods)

**Files:**
- Create: `frontend/src/app/[year]/settings/page.tsx`
- Create: `frontend/src/components/settings/category-manager.tsx`
- Create: `frontend/src/components/settings/payment-method-manager.tsx`

- [ ] **Step 1: Create category manager component**

```tsx
// src/components/settings/category-manager.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Category } from "@/types";
import { apiPost, apiPut, apiDelete } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface CategoryManagerProps {
  categories: Category[];
  onRefresh: () => void;
}

export function CategoryManager({
  categories,
  onRefresh,
}: CategoryManagerProps) {
  const t = useTranslations("settings");
  const tc = useTranslations("common");
  const [name, setName] = useState("");
  const [domain, setDomain] = useState<string>("expense");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");

  async function handleAdd() {
    if (!name.trim()) return;
    await apiPost("/api/v1/categories", { name, domain });
    setName("");
    onRefresh();
  }

  async function handleUpdate(id: string) {
    await apiPut(`/api/v1/categories/${id}`, { name: editName });
    setEditingId(null);
    onRefresh();
  }

  async function handleDelete(id: string) {
    await apiDelete(`/api/v1/categories/${id}`);
    onRefresh();
  }

  async function handleSeedDefaults() {
    await apiPost("/api/v1/categories/seed", {});
    onRefresh();
  }

  const grouped = {
    income: categories.filter((c) => c.domain === "income"),
    expense: categories.filter((c) => c.domain === "expense"),
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">{t("categories")}</h3>
        <Button variant="outline" size="sm" onClick={handleSeedDefaults}>
          {t("seedDefaults")}
        </Button>
      </div>

      <div className="flex gap-2">
        <Input
          placeholder={t("addCategory")}
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleAdd()}
        />
        <Select value={domain} onValueChange={setDomain}>
          <SelectTrigger className="w-36">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="income">Income</SelectItem>
            <SelectItem value="expense">Expense</SelectItem>
            <SelectItem value="wishlist">Wishlist</SelectItem>
          </SelectContent>
        </Select>
        <Button onClick={handleAdd}>{tc("create")}</Button>
      </div>

      {Object.entries(grouped).map(([key, cats]) => (
        <div key={key} className="space-y-2">
          <h4 className="text-sm font-medium text-gray-500 uppercase">{key}</h4>
          <div className="flex flex-wrap gap-2">
            {cats.map((cat) => (
              <div key={cat.id} className="flex items-center gap-1">
                {editingId === cat.id ? (
                  <>
                    <Input
                      className="h-7 w-32 text-sm"
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      onKeyDown={(e) =>
                        e.key === "Enter" && handleUpdate(cat.id)
                      }
                    />
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => handleUpdate(cat.id)}
                    >
                      {tc("save")}
                    </Button>
                  </>
                ) : (
                  <Badge
                    variant="secondary"
                    className="cursor-pointer"
                    onClick={() => {
                      setEditingId(cat.id);
                      setEditName(cat.name);
                    }}
                  >
                    {cat.name}
                    <button
                      className="ml-1 text-gray-400 hover:text-red-500"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDelete(cat.id);
                      }}
                    >
                      ×
                    </button>
                  </Badge>
                )}
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
```

- [ ] **Step 2: Create payment method manager component**

```tsx
// src/components/settings/payment-method-manager.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { PaymentMethod } from "@/types";
import { apiPost, apiDelete } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const typeLabels: Record<string, string> = {
  cash: "Cash",
  debit_card: "Debit Card",
  credit_card: "Credit Card",
  digital_wallet: "Digital Wallet",
  crypto: "Crypto",
};

interface PaymentMethodManagerProps {
  paymentMethods: PaymentMethod[];
  onRefresh: () => void;
}

export function PaymentMethodManager({
  paymentMethods,
  onRefresh,
}: PaymentMethodManagerProps) {
  const t = useTranslations("settings");
  const tc = useTranslations("common");
  const [name, setName] = useState("");
  const [type, setType] = useState("debit_card");
  const [details, setDetails] = useState("");

  async function handleAdd() {
    if (!name.trim()) return;
    await apiPost("/api/v1/payment-methods", {
      name,
      type,
      details: details || undefined,
    });
    setName("");
    setDetails("");
    onRefresh();
  }

  async function handleDelete(id: string) {
    await apiDelete(`/api/v1/payment-methods/${id}`);
    onRefresh();
  }

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-semibold">{t("paymentMethods")}</h3>

      <div className="flex gap-2">
        <Input
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <Select value={type} onValueChange={setType}>
          <SelectTrigger className="w-40">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {Object.entries(typeLabels).map(([key, label]) => (
              <SelectItem key={key} value={key}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Input
          placeholder="Details (optional)"
          value={details}
          onChange={(e) => setDetails(e.target.value)}
          className="w-36"
        />
        <Button onClick={handleAdd}>{tc("create")}</Button>
      </div>

      <div className="space-y-2">
        {paymentMethods.map((pm) => (
          <div
            key={pm.id}
            className="flex items-center justify-between p-3 bg-white rounded-lg border"
          >
            <div className="flex items-center gap-3">
              <span className="font-medium">{pm.name}</span>
              <Badge variant="outline">{typeLabels[pm.type] || pm.type}</Badge>
              {pm.details && (
                <span className="text-sm text-gray-500">{pm.details}</span>
              )}
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="text-red-600"
              onClick={() => handleDelete(pm.id)}
            >
              {tc("delete")}
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Create settings page**

```tsx
// src/app/[year]/settings/page.tsx
"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { apiGet } from "@/lib/api";
import { Category, PaymentMethod, ListResponse } from "@/types";
import { CategoryManager } from "@/components/settings/category-manager";
import { PaymentMethodManager } from "@/components/settings/payment-method-manager";
import { Separator } from "@/components/ui/separator";

export default function SettingsPage() {
  const t = useTranslations("settings");
  const [categories, setCategories] = useState<Category[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);

  async function loadCategories() {
    const [incomeRes, expenseRes] = await Promise.all([
      apiGet<ListResponse<Category>>("/api/v1/categories?domain=income"),
      apiGet<ListResponse<Category>>("/api/v1/categories?domain=expense"),
    ]);
    setCategories([...incomeRes.data, ...expenseRes.data]);
  }

  async function loadPaymentMethods() {
    const res = await apiGet<ListResponse<PaymentMethod>>(
      "/api/v1/payment-methods"
    );
    setPaymentMethods(res.data);
  }

  useEffect(() => {
    loadCategories();
    loadPaymentMethods();
  }, []);

  return (
    <div className="max-w-3xl space-y-8">
      <h2 className="text-2xl font-bold">{t("title")}</h2>
      <CategoryManager categories={categories} onRefresh={loadCategories} />
      <Separator />
      <PaymentMethodManager
        paymentMethods={paymentMethods}
        onRefresh={loadPaymentMethods}
      />
    </div>
  );
}
```

- [ ] **Step 4: Verify frontend compiles**

Run: `cd /home/folkrom/projects/finance-tracker/frontend && npm run build`
Expected: Build succeeds.

- [ ] **Step 5: Commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add frontend/src/app/[year]/settings/ frontend/src/components/settings/
git commit -m "feat: settings page with category and payment method management"
```

---

## Final Verification

### Task 18: End-to-End Smoke Test

- [ ] **Step 1: Start all services**

```bash
cd /home/folkrom/projects/finance-tracker
docker compose up -d postgres
# In terminal 1:
cd backend && cp .env.example .env  # Fill in real values
go run ./cmd/server
# In terminal 2:
cd frontend && cp .env.local.example .env.local  # Fill in real values
npm run dev
```

- [ ] **Step 2: Verify backend health**

```bash
curl http://localhost:8080/health
```
Expected: `{"status":"ok"}`

- [ ] **Step 3: Run all backend tests**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go test ./... -v -count=1
```
Expected: All tests PASS.

- [ ] **Step 4: Run frontend build check**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npm run build
```
Expected: Build succeeds.

- [ ] **Step 5: Final commit**

```bash
cd /home/folkrom/projects/finance-tracker
git add -A
git commit -m "chore: Plan 1 complete — foundation, auth, shared entities, income module"
```

---

## API Endpoint Summary

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/api/v1/categories` | Create category |
| GET | `/api/v1/categories?domain=` | List categories by domain |
| POST | `/api/v1/categories/seed` | Seed default categories |
| GET | `/api/v1/categories/:id` | Get category |
| PUT | `/api/v1/categories/:id` | Update category |
| DELETE | `/api/v1/categories/:id` | Delete category |
| POST | `/api/v1/payment-methods` | Create payment method |
| GET | `/api/v1/payment-methods` | List payment methods |
| GET | `/api/v1/payment-methods/:id` | Get payment method |
| PUT | `/api/v1/payment-methods/:id` | Update payment method |
| DELETE | `/api/v1/payment-methods/:id` | Delete payment method |
| POST | `/api/v1/years/:year/incomes` | Create income |
| GET | `/api/v1/years/:year/incomes` | List incomes by year |
| GET | `/api/v1/incomes/:id` | Get income |
| PUT | `/api/v1/incomes/:id` | Update income |
| DELETE | `/api/v1/incomes/:id` | Delete income |
