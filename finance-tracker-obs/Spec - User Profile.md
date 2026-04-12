# User Profile Design

**Date:** 2026-04-12
**Status:** Approved
**Scope:** Profile table, auto-creation middleware, CRUD API, settings page section

---

## Problem

No user preferences persist across sessions. Language and currency are hardcoded or browser-default. Users have no way to set their preferred currency for amount display or lock their language preference.

## Design

### Migration (000010)

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

No FK to `auth.users` — keeps the schema portable. `user_id` comes from JWT `sub` claim, same as all other tables.

### Model

```go
type Profile struct {
    Base
    Currency string `gorm:"type:varchar(3);not null;default:'MXN'" json:"currency"`
    Language string `gorm:"type:varchar(5);not null;default:'en'" json:"language"`
}
```

Uses `Base` which includes `UserID uuid.UUID` (non-nullable, same as income/expense/etc).

### Auto-Creation Middleware

Runs after auth middleware on all `/api/v1/` routes.

**Logic:**
1. Get `user_id` from `c.Locals("user_id")` (set by auth middleware)
2. Check an in-memory `sync.Map` keyed by `user_id` string — if present, skip DB check
3. If not in cache: query `profiles` table for `user_id`
4. If no row: insert default profile (MXN, en)
5. Store `user_id` in sync.Map so subsequent requests skip the DB
6. Store profile in `c.Locals("profile")` for downstream handlers

The sync.Map is per-process and resets on server restart — acceptable for a single-server deployment. No TTL needed; profiles rarely change and the cache only gates the existence check, not the data.

**Middleware signature:**
```go
func NewProfileMiddleware(profileRepo *repository.ProfileRepository) fiber.Handler
```

Registered in router after auth middleware:
```go
api := app.Group("/api/v1", middleware.NewAuthMiddleware(keyfunc), middleware.NewProfileMiddleware(profileRepo))
```

### Repository

```go
type ProfileRepository struct { db *gorm.DB }

func (r *ProfileRepository) GetByUserID(userID uuid.UUID) (*model.Profile, error)
func (r *ProfileRepository) Create(profile *model.Profile) error
func (r *ProfileRepository) Update(profile *model.Profile) error
```

### Service

```go
type ProfileService struct { repo *repository.ProfileRepository }

func (s *ProfileService) GetOrCreate(userID uuid.UUID) (*model.Profile, error)
func (s *ProfileService) Update(userID uuid.UUID, currency, language string) (*model.Profile, error)
```

`GetOrCreate`: queries by user_id, creates default if not found.
`Update`: validates currency is in allowed list, language is "en" or "es", then updates.

### Handler

```go
type ProfileHandler struct { svc *service.ProfileService }

func (h *ProfileHandler) Get(c *fiber.Ctx) error    // GET /api/v1/profile
func (h *ProfileHandler) Update(c *fiber.Ctx) error  // PUT /api/v1/profile
```

**GET response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "currency": "MXN",
  "language": "en",
  "created_at": "...",
  "updated_at": "..."
}
```

**PUT request body:**
```json
{
  "currency": "USD",
  "language": "es"
}
```

Both fields optional — only updates provided fields. Validates:
- `currency`: must be one of MXN, USD, EUR, GBP, BRL, COP, ARS
- `language`: must be "en" or "es"

### Router Changes

```go
// Profile
profile := api.Group("/profile")
profile.Get("/", profileHandler.Get)
profile.Put("/", profileHandler.Update)
```

### Frontend

**Settings page** (`frontend/src/app/[year]/settings/page.tsx`):
- Add a "Profile" section above or below the existing Categories section
- Fetch `GET /api/v1/profile` on page load
- Currency: `<Select>` dropdown with options: MXN, USD, EUR, GBP, BRL, COP, ARS
- Language: `<Select>` dropdown with options: English (en), Spanish (es)
- Save button calls `PUT /api/v1/profile` with changed values
- Toast on success/error

**TypeScript type:**
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

### Allowed Values

**Currencies:** MXN, USD, EUR, GBP, BRL, COP, ARS
- Covers Mexico, US, Europe, UK, Brazil, Colombia, Argentina
- Can expand later without migration (just add to the validation list)

**Languages:** en, es
- Matches existing i18n translations

### What This Does NOT Include

- Avatar upload (no consumer yet)
- Display name (no consumer yet)
- Timezone (no scheduling features yet)
- Profile deletion (always exists via auto-create)
- Currency formatting integration (separate task — this just stores the preference)
