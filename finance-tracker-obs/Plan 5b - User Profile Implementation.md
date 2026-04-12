# User Profile Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add user profile with currency and language preferences, auto-created on first authenticated request.

**Architecture:** New `profiles` table with auto-creation middleware. Profile middleware runs after auth middleware, ensures a profile row exists for every authenticated user, and caches the check in a sync.Map. Two API endpoints (GET/PUT) for reading and updating preferences. Frontend adds a Profile section to the existing Settings page.

**Tech Stack:** Go/Fiber, GORM, PostgreSQL, golang-migrate, Next.js/React, TypeScript, shadcn/ui

---

### Task 1: Migration — Create profiles table

**Files:**
- Create: `backend/migrations/000010_create_profiles.up.sql`
- Create: `backend/migrations/000010_create_profiles.down.sql`

- [ ] **Step 1: Write the up migration**

Create `backend/migrations/000010_create_profiles.up.sql`:

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    language VARCHAR(5) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);
```

- [ ] **Step 2: Write the down migration**

Create `backend/migrations/000010_create_profiles.down.sql`:

```sql
DROP TABLE IF EXISTS profiles;
```

- [ ] **Step 3: Commit**

```bash
git add backend/migrations/000010_create_profiles.up.sql backend/migrations/000010_create_profiles.down.sql
git commit -m "feat: migration 000010 — create profiles table"
```

---

### Task 2: Model — Profile struct

**Files:**
- Create: `backend/internal/model/profile.go`

- [ ] **Step 1: Create the Profile model**

Create `backend/internal/model/profile.go`:

```go
package model

type Profile struct {
	Base
	Currency string `gorm:"type:varchar(3);not null;default:'MXN'" json:"currency"`
	Language string `gorm:"type:varchar(5);not null;default:'en'" json:"language"`
}

func (Profile) TableName() string {
	return "profiles"
}
```

Uses `Base` which provides `ID`, `UserID` (non-nullable uuid.UUID), `CreatedAt`, `UpdatedAt`.

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/model/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/model/profile.go
git commit -m "feat: Profile model — currency and language preferences"
```

---

### Task 3: Repository — Profile CRUD

**Files:**
- Create: `backend/internal/repository/profile.go`

- [ ] **Step 1: Create the repository**

Create `backend/internal/repository/profile.go`:

```go
package repository

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByUserID(userID uuid.UUID) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) Create(profile *model.Profile) error {
	return r.db.Create(profile).Error
}

func (r *ProfileRepository) Update(profile *model.Profile) error {
	return r.db.Save(profile).Error
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/repository/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/repository/profile.go
git commit -m "feat: profile repository — GetByUserID, Create, Update"
```

---

### Task 4: Service — GetOrCreate and Update with validation

**Files:**
- Create: `backend/internal/service/profile.go`

- [ ] **Step 1: Create the service**

Create `backend/internal/service/profile.go`:

```go
package service

import (
	"errors"
	"slices"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	allowedCurrencies = []string{"MXN", "USD", "EUR", "GBP", "BRL", "COP", "ARS"}
	allowedLanguages  = []string{"en", "es"}

	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidLanguage = errors.New("invalid language")
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetOrCreate(userID uuid.UUID) (*model.Profile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	profile = &model.Profile{
		Base:     model.Base{UserID: userID},
		Currency: "MXN",
		Language: "en",
	}
	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *ProfileService) Update(userID uuid.UUID, currency, language string) (*model.Profile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if currency != "" {
		if !slices.Contains(allowedCurrencies, currency) {
			return nil, ErrInvalidCurrency
		}
		profile.Currency = currency
	}

	if language != "" {
		if !slices.Contains(allowedLanguages, language) {
			return nil, ErrInvalidLanguage
		}
		profile.Language = language
	}

	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}
	return profile, nil
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/service/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/profile.go
git commit -m "feat: profile service — GetOrCreate with defaults, Update with validation"
```

---

### Task 5: Handler — GET and PUT endpoints

**Files:**
- Create: `backend/internal/handler/profile.go`

- [ ] **Step 1: Create the handler**

Create `backend/internal/handler/profile.go`:

```go
package handler

import (
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	svc *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

type updateProfileRequest struct {
	Currency string `json:"currency"`
	Language string `json:"language"`
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	profile, err := h.svc.GetOrCreate(userID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to get profile")
	}
	return respondJSON(c, profile)
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request body")
	}

	profile, err := h.svc.Update(userID, req.Currency, req.Language)
	if err != nil {
		if err == service.ErrInvalidCurrency {
			return respondError(c, fiber.StatusBadRequest, "invalid currency — allowed: MXN, USD, EUR, GBP, BRL, COP, ARS")
		}
		if err == service.ErrInvalidLanguage {
			return respondError(c, fiber.StatusBadRequest, "invalid language — allowed: en, es")
		}
		return respondError(c, fiber.StatusInternalServerError, "failed to update profile")
	}
	return respondJSON(c, profile)
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/handler/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/profile.go
git commit -m "feat: profile handler — GET and PUT with validation errors"
```

---

### Task 6: Profile Middleware — Auto-create on first request

**Files:**
- Create: `backend/internal/middleware/profile.go`

- [ ] **Step 1: Create the middleware**

Create `backend/internal/middleware/profile.go`:

```go
package middleware

import (
	"sync"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NewProfileMiddleware ensures a profile exists for every authenticated user.
// Uses a sync.Map to cache user IDs that have been verified, avoiding a DB
// lookup on every request after the first.
func NewProfileMiddleware(profileRepo *repository.ProfileRepository) fiber.Handler {
	var seen sync.Map

	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		userIDStr := userID.String()

		// Fast path: already verified this user has a profile
		if _, ok := seen.Load(userIDStr); ok {
			return c.Next()
		}

		// Check DB for existing profile
		_, err := profileRepo.GetByUserID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Auto-create default profile
				profile := &model.Profile{
					Base:     model.Base{UserID: userID},
					Currency: "MXN",
					Language: "en",
				}
				if createErr := profileRepo.Create(profile); createErr != nil {
					// Another request might have created it concurrently — check again
					if _, retryErr := profileRepo.GetByUserID(userID); retryErr != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "failed to create profile",
						})
					}
				}
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to check profile",
				})
			}
		}

		// Cache this user so we skip the DB check next time
		seen.Store(userIDStr, true)
		return c.Next()
	}
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./internal/middleware/...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/profile.go
git commit -m "feat: profile middleware — auto-create profile on first authenticated request"
```

---

### Task 7: Router and main.go — Wire everything together

**Files:**
- Modify: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Update router.go**

Add `profileHandler *handler.ProfileHandler` and `profileRepo *repository.ProfileRepository` parameters to the `Setup` function signature. Add profile middleware and routes.

The updated function signature:

```go
func Setup(
	app *fiber.App,
	keyfunc jwt.Keyfunc,
	profileRepo *repository.ProfileRepository,
	categoryHandler *handler.CategoryHandler,
	paymentMethodHandler *handler.PaymentMethodHandler,
	incomeHandler *handler.IncomeHandler,
	expenseHandler *handler.ExpenseHandler,
	debtHandler *handler.DebtHandler,
	budgetHandler *handler.BudgetHandler,
	dashboardHandler *handler.DashboardHandler,
	cardHandler *handler.CardHandler,
	wishlistHandler *handler.WishlistItemHandler,
	profileHandler *handler.ProfileHandler,
) {
```

Update the `api` group to include profile middleware:

```go
api := app.Group("/api/v1", middleware.NewAuthMiddleware(keyfunc), middleware.NewProfileMiddleware(profileRepo))
```

Add profile routes after the existing wishlist routes:

```go
	// Profile
	profile := api.Group("/profile")
	profile.Get("/", profileHandler.Get)
	profile.Put("/", profileHandler.Update)
```

Add the `repository` import to the import block:

```go
import (
	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)
```

- [ ] **Step 2: Update main.go**

Add profile repo, service, and handler creation. In the Repositories section, add:

```go
profileRepo := repository.NewProfileRepository(db)
```

In the Services section, add:

```go
profileSvc := service.NewProfileService(profileRepo)
```

In the Handlers section, add:

```go
profileHandler := handler.NewProfileHandler(profileSvc)
```

Update the `router.Setup` call to include `profileRepo` and `profileHandler`:

```go
router.Setup(app, jwks.Keyfunc, profileRepo, categoryHandler, paymentMethodHandler, incomeHandler, expenseHandler, debtHandler, budgetHandler, dashboardHandler, cardHandler, wishlistHandler, profileHandler)
```

- [ ] **Step 3: Verify full backend compiles**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go build ./...
```

- [ ] **Step 4: Run auth middleware tests**

```bash
cd /home/folkrom/projects/finance-tracker/backend && go test ./internal/middleware/... -v
```

Expected: 3/3 pass.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/router/router.go backend/cmd/server/main.go
git commit -m "feat: wire profile repo, service, handler, middleware into router"
```

---

### Task 8: Frontend — Profile section in Settings page

**Files:**
- Modify: `frontend/src/types/index.ts`
- Create: `frontend/src/components/settings/profile-manager.tsx`
- Modify: `frontend/src/app/[year]/settings/page.tsx`

- [ ] **Step 1: Add Profile type**

In `frontend/src/types/index.ts`, add at the end of the file:

```typescript
export interface Profile {
  id: string;
  user_id: string;
  currency: string;
  language: string;
  created_at: string;
  updated_at: string;
}
```

- [ ] **Step 2: Create ProfileManager component**

Create `frontend/src/components/settings/profile-manager.tsx`:

```tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Profile } from "@/types";
import { apiPut } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface ProfileManagerProps {
  profile: Profile;
  onRefresh: () => void;
}

const CURRENCIES = [
  { value: "MXN", label: "MXN — Mexican Peso" },
  { value: "USD", label: "USD — US Dollar" },
  { value: "EUR", label: "EUR — Euro" },
  { value: "GBP", label: "GBP — British Pound" },
  { value: "BRL", label: "BRL — Brazilian Real" },
  { value: "COP", label: "COP — Colombian Peso" },
  { value: "ARS", label: "ARS — Argentine Peso" },
];

const LANGUAGES = [
  { value: "en", label: "English" },
  { value: "es", label: "Español" },
];

export function ProfileManager({ profile, onRefresh }: ProfileManagerProps) {
  const t = useTranslations("settings");
  const [currency, setCurrency] = useState(profile.currency);
  const [language, setLanguage] = useState(profile.language);
  const [saving, setSaving] = useState(false);

  const hasChanges = currency !== profile.currency || language !== profile.language;

  const handleSave = async () => {
    setSaving(true);
    try {
      await apiPut<Profile>("/api/v1/profile", { currency, language });
      toast.success("Profile updated");
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update profile");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold">{t("profile") || "Profile"}</h2>

      <div className="flex gap-4 items-end flex-wrap">
        <div className="space-y-1">
          <Label>Currency</Label>
          <Select value={currency} onValueChange={setCurrency}>
            <SelectTrigger className="w-56">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {CURRENCIES.map((c) => (
                <SelectItem key={c.value} value={c.value}>
                  {c.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-1">
          <Label>Language</Label>
          <Select value={language} onValueChange={setLanguage}>
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {LANGUAGES.map((l) => (
                <SelectItem key={l.value} value={l.value}>
                  {l.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Button onClick={handleSave} disabled={saving || !hasChanges}>
          {saving ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Update Settings page to include ProfileManager**

Replace the contents of `frontend/src/app/[year]/settings/page.tsx`:

```tsx
"use client";

import { useState, useEffect, useCallback } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { Category, PaymentMethod, Profile, ListResponse } from "@/types";
import { ProfileManager } from "@/components/settings/profile-manager";
import { CategoryManager } from "@/components/settings/category-manager";
import { PaymentMethodManager } from "@/components/settings/payment-method-manager";
import { Separator } from "@/components/ui/separator";

export default function SettingsPage() {
  const t = useTranslations("settings");
  const tCommon = useTranslations("common");

  const [categories, setCategories] = useState<Category[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [incomeRes, expenseRes, pmRes, profileRes] = await Promise.all([
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=income"),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=expense"),
        apiGet<ListResponse<PaymentMethod>>("/api/v1/payment-methods"),
        apiGet<Profile>("/api/v1/profile"),
      ]);
      setCategories([...incomeRes.data, ...expenseRes.data]);
      setPaymentMethods(pmRes.data);
      setProfile(profileRes);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load settings");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {tCommon("loading")}
      </div>
    );
  }

  return (
    <div className="space-y-8 max-w-3xl">
      <h1 className="text-2xl font-bold">{t("title")}</h1>

      {profile && <ProfileManager profile={profile} onRefresh={loadData} />}

      <Separator />

      <CategoryManager categories={categories} onRefresh={loadData} />

      <Separator />

      <PaymentMethodManager paymentMethods={paymentMethods} onRefresh={loadData} />
    </div>
  );
}
```

- [ ] **Step 4: Verify frontend compiles**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/types/index.ts frontend/src/components/settings/profile-manager.tsx frontend/src/app/\\[year\\]/settings/page.tsx
git commit -m "feat: profile section in settings — currency and language selectors"
```

---

### Task 9: Integration Test — Verify end-to-end

**Files:** None (manual verification)

- [ ] **Step 1: Run the migration**

```bash
cd /home/folkrom/projects/finance-tracker && mise run migrate-up
```

Expected: migration 000010 applied.

- [ ] **Step 2: Start backend and verify auto-creation**

```bash
mise run dev-backend
```

Then with a valid token:

```bash
curl -s http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <token>" | jq
```

Expected: profile auto-created with `currency: "MXN"`, `language: "en"`.

- [ ] **Step 3: Test update with valid values**

```bash
curl -s -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"currency": "USD", "language": "es"}' | jq
```

Expected: updated profile returned.

- [ ] **Step 4: Test update with invalid values**

```bash
curl -s -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"currency": "INVALID"}' | jq
```

Expected: 400 with "invalid currency" message.

- [ ] **Step 5: Verify Settings page in browser**

Open http://localhost:3000/2026/settings:
- Profile section appears at the top with currency and language selectors
- Changing a value enables the Save button
- Saving shows success toast and persists across page reload

- [ ] **Step 6: Final commit (if any adjustments needed)**

Only if fixes were needed during testing.
