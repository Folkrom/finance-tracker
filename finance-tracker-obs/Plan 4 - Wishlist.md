# Plan 4: Wishlist Module

> **For agentic workers:** Use subagent-driven-development to implement this plan task-by-task.

**Goal:** Build the Wishlist module with Kanban statuses, 3 views (Gallery, Table, Board), image support (URL-based), up to 5 purchase links, separate category list, and priority sorting. Wishlist is NOT year-scoped — items persist across years.

**Depends on:** Plan 3 complete (all 7 modules except Wishlist working).

**Key design decisions:**
- Wishlist uses its own category list (domain: "wishlist"), separate from Expenses/Debt
- Not year-scoped — lives at `/wishlist` (already linked in sidebar footer)
- Kanban statuses with 3 groups: To-do, In Progress, Complete
- Image via URL (no upload for now — simplifies storage)
- Up to 5 purchase links stored as JSON array

---

## Phase A: Backend

### Task 1: Wishlist Migration

**Files:**
- Create: `backend/migrations/000008_create_wishlist_items.up.sql`
- Create: `backend/migrations/000008_create_wishlist_items.down.sql`

```sql
-- 000008_create_wishlist_items.up.sql
CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_url TEXT,
    price DECIMAL(12, 2),
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    links TEXT[] DEFAULT '{}',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    status VARCHAR(30) NOT NULL DEFAULT 'interested' CHECK (status IN (
        'interested',
        'saving_for', 'waiting_for_sale', 'ordered',
        'purchased', 'received', 'cancelled'
    )),
    target_date DATE,
    monthly_contribution DECIMAL(12, 2),
    contribution_currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wishlist_user_id ON wishlist_items(user_id);
CREATE INDEX idx_wishlist_status ON wishlist_items(user_id, status);
CREATE INDEX idx_wishlist_category ON wishlist_items(user_id, category_id);
CREATE INDEX idx_wishlist_priority ON wishlist_items(user_id, priority);
```

```sql
-- 000008_create_wishlist_items.down.sql
DROP TABLE IF EXISTS wishlist_items;
```

Notes:
- `links` is a Postgres text array — stores up to 5 URLs. GORM handles this with `pq.StringArray`.
- `status` maps to the Kanban groups:
  - **To-do:** interested
  - **In Progress:** saving_for, waiting_for_sale, ordered
  - **Complete:** purchased, received, cancelled
- `priority`: low, medium, high

- [ ] **Step 1: Create migration files**
- [ ] **Step 2: Commit**

---

### Task 2: Wishlist Model

**Files:**
- Create: `backend/internal/model/wishlist_item.go`

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type WishlistPriority string

const (
	WishlistPriorityLow    WishlistPriority = "low"
	WishlistPriorityMedium WishlistPriority = "medium"
	WishlistPriorityHigh   WishlistPriority = "high"
)

type WishlistStatus string

const (
	// To-do
	WishlistStatusInterested WishlistStatus = "interested"
	// In Progress
	WishlistStatusSavingFor      WishlistStatus = "saving_for"
	WishlistStatusWaitingForSale WishlistStatus = "waiting_for_sale"
	WishlistStatusOrdered        WishlistStatus = "ordered"
	// Complete
	WishlistStatusPurchased WishlistStatus = "purchased"
	WishlistStatusReceived  WishlistStatus = "received"
	WishlistStatusCancelled WishlistStatus = "cancelled"
)

type WishlistItem struct {
	Base
	Name                 string           `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL             *string          `gorm:"type:text" json:"image_url,omitempty"`
	Price                *decimal.Decimal `gorm:"type:decimal(12,2)" json:"price,omitempty"`
	Currency             string           `gorm:"type:varchar(3);not null;default:MXN" json:"currency"`
	Links                pq.StringArray   `gorm:"type:text[]" json:"links"`
	CategoryID           *uuid.UUID       `gorm:"type:uuid" json:"category_id,omitempty"`
	Category             *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Priority             WishlistPriority `gorm:"type:varchar(10);not null;default:medium" json:"priority"`
	Status               WishlistStatus   `gorm:"type:varchar(30);not null;default:interested" json:"status"`
	TargetDate           *time.Time       `gorm:"type:date" json:"target_date,omitempty"`
	MonthlyContribution  *decimal.Decimal `gorm:"type:decimal(12,2)" json:"monthly_contribution,omitempty"`
	ContributionCurrency string           `gorm:"type:varchar(3);not null;default:MXN" json:"contribution_currency"`
}

func (WishlistItem) TableName() string { return "wishlist_items" }
```

**Important:** Add `github.com/lib/pq` dependency: `cd backend && go get github.com/lib/pq`

- [ ] **Step 1: Install pq dependency**
- [ ] **Step 2: Create model file**
- [ ] **Step 3: Verify compiles**
- [ ] **Step 4: Commit**

---

### Task 3: Wishlist CRUD (Repo + Service + Handler + Tests)

**Files:**
- Create: `backend/internal/repository/wishlist_item.go`
- Create: `backend/internal/repository/wishlist_item_test.go`
- Create: `backend/internal/service/wishlist_item.go`
- Create: `backend/internal/handler/wishlist_item.go`
- Create: `backend/internal/handler/wishlist_item_test.go`

**Repository methods:**
- `Create(item *WishlistItem) error`
- `ListByUser(userID uuid.UUID) ([]WishlistItem, error)` — preloads Category, ordered by priority DESC then name
- `ListByStatus(userID uuid.UUID, statuses []WishlistStatus) ([]WishlistItem, error)` — filter by multiple statuses, preloads Category
- `ListByCategory(userID uuid.UUID, categoryID uuid.UUID) ([]WishlistItem, error)` — for Board view
- `GetByID(userID, id uuid.UUID) (*WishlistItem, error)` — preloads Category
- `Update(item *WishlistItem) error`
- `UpdateStatus(userID, id uuid.UUID, status WishlistStatus) error` — lightweight status-only update for drag-and-drop
- `Delete(userID, id uuid.UUID) error`

**Service:**
- Delegates to repo
- Validates: links max 5 items, priority must be valid, status must be valid
- On create: default currency MXN, default priority medium, default status interested

**Handler request struct:**
```go
type createWishlistItemRequest struct {
	Name                 string   `json:"name"`
	ImageURL             *string  `json:"image_url"`
	Price                *string  `json:"price"`
	Currency             string   `json:"currency"`
	Links                []string `json:"links"`
	CategoryID           *string  `json:"category_id"`
	Priority             string   `json:"priority"`
	Status               string   `json:"status"`
	TargetDate           *string  `json:"target_date"`
	MonthlyContribution  *string  `json:"monthly_contribution"`
	ContributionCurrency string   `json:"contribution_currency"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}
```

**Handler endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/wishlist` | Create item |
| GET | `/api/v1/wishlist` | List items (optional `?status=` filter, comma-separated) |
| GET | `/api/v1/wishlist/:id` | Get single item |
| PUT | `/api/v1/wishlist/:id` | Update item |
| PATCH | `/api/v1/wishlist/:id/status` | Update status only (for drag-and-drop) |
| DELETE | `/api/v1/wishlist/:id` | Delete item |

**Tests:**
- Repository: Create, ListByUser, ListByStatus, UpdateStatus, Delete
- Handler: CreateAndList, UpdateStatus

- [ ] **Step 1: Create repository with tests**
- [ ] **Step 2: Create service**
- [ ] **Step 3: Create handler with tests**
- [ ] **Step 4: Verify tests pass**
- [ ] **Step 5: Commit**

---

### Task 4: Seed Wishlist Categories + Wire Routes

**Files:**
- Modify: `backend/internal/service/category.go` — add wishlist defaults to SeedDefaults
- Modify: `backend/internal/router/router.go` — add wishlist routes
- Modify: `backend/cmd/server/main.go` — wire wishlist repo/service/handler

**Wishlist default categories to seed:**
```go
wishlistCategories := []string{
	"Electronics", "Clothing", "Home & Kitchen", "Books & Media",
	"Sports & Outdoors", "Beauty & Personal Care", "Toys & Games", "Other",
}
```

Add to `SeedDefaults` with domain `CategoryDomainWishlist`.

**Routes:**
```go
// Wishlist
wishlist := api.Group("/wishlist")
wishlist.Post("/", wishlistHandler.Create)
wishlist.Get("/", wishlistHandler.List)
wishlist.Get("/:id", wishlistHandler.GetByID)
wishlist.Put("/:id", wishlistHandler.Update)
wishlist.Patch("/:id/status", wishlistHandler.UpdateStatus)
wishlist.Delete("/:id", wishlistHandler.Delete)
```

- [ ] **Step 1: Add wishlist categories to SeedDefaults**
- [ ] **Step 2: Update router.go**
- [ ] **Step 3: Update main.go**
- [ ] **Step 4: Verify backend compiles**
- [ ] **Step 5: Commit**

---

## Phase B: Frontend

### Task 5: TypeScript Types + i18n

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/lib/i18n/en.json`
- Modify: `frontend/src/lib/i18n/es.json`

**Types:**
```ts
export type WishlistPriority = "low" | "medium" | "high";
export type WishlistStatus =
  | "interested"
  | "saving_for"
  | "waiting_for_sale"
  | "ordered"
  | "purchased"
  | "received"
  | "cancelled";

export interface WishlistItem {
  id: string;
  user_id: string;
  name: string;
  image_url?: string;
  price?: string;
  currency: string;
  links: string[];
  category_id?: string;
  category?: Category;
  priority: WishlistPriority;
  status: WishlistStatus;
  target_date?: string;
  monthly_contribution?: string;
  contribution_currency: string;
  created_at: string;
  updated_at: string;
}

export const WISHLIST_STATUS_GROUPS = {
  todo: ["interested"] as WishlistStatus[],
  in_progress: ["saving_for", "waiting_for_sale", "ordered"] as WishlistStatus[],
  complete: ["purchased", "received", "cancelled"] as WishlistStatus[],
};
```

**EN translations:**
```json
"wishlist": {
  "title": "Wishlist",
  "name": "Item Name",
  "imageUrl": "Image URL",
  "price": "Price",
  "links": "Links to Buy",
  "addLink": "Add Link",
  "category": "Category",
  "priority": "Priority",
  "status": "Status",
  "targetDate": "Target Purchase Date",
  "monthlyContribution": "Monthly Contribution",
  "addItem": "Add Item",
  "editItem": "Edit Item",
  "noItems": "Your wishlist is empty",
  "priorityLow": "Low",
  "priorityMedium": "Medium",
  "priorityHigh": "High",
  "statusInterested": "Interested",
  "statusSavingFor": "Saving For",
  "statusWaitingForSale": "Waiting for Sale",
  "statusOrdered": "Ordered",
  "statusPurchased": "Purchased",
  "statusReceived": "Received",
  "statusCancelled": "Cancelled",
  "groupTodo": "To-Do",
  "groupInProgress": "In Progress",
  "groupComplete": "Complete",
  "viewGallery": "Gallery",
  "viewTable": "Table",
  "viewBoard": "Board"
}
```

**ES translations:**
```json
"wishlist": {
  "title": "Lista de Deseos",
  "name": "Nombre del Artículo",
  "imageUrl": "URL de Imagen",
  "price": "Precio",
  "links": "Enlaces de Compra",
  "addLink": "Agregar Enlace",
  "category": "Categoría",
  "priority": "Prioridad",
  "status": "Estado",
  "targetDate": "Fecha Objetivo de Compra",
  "monthlyContribution": "Contribución Mensual",
  "addItem": "Agregar Artículo",
  "editItem": "Editar Artículo",
  "noItems": "Tu lista de deseos está vacía",
  "priorityLow": "Baja",
  "priorityMedium": "Media",
  "priorityHigh": "Alta",
  "statusInterested": "Interesado",
  "statusSavingFor": "Ahorrando",
  "statusWaitingForSale": "Esperando Oferta",
  "statusOrdered": "Ordenado",
  "statusPurchased": "Comprado",
  "statusReceived": "Recibido",
  "statusCancelled": "Cancelado",
  "groupTodo": "Por Hacer",
  "groupInProgress": "En Progreso",
  "groupComplete": "Completado",
  "viewGallery": "Galería",
  "viewTable": "Tabla",
  "viewBoard": "Tablero"
}
```

- [ ] **Step 1: Add types**
- [ ] **Step 2: Add EN translations**
- [ ] **Step 3: Add ES translations**
- [ ] **Step 4: Commit**

---

### Task 6: Wishlist Form Component

**Files:**
- Create: `frontend/src/components/wishlist/wishlist-form.tsx`

A dialog form with React Hook Form + Zod. This form is more complex than others due to the number of fields.

**Fields:**
- name (text, required)
- image_url (URL input, optional)
- price (number, optional)
- currency (MXN/USD select)
- links (dynamic list of up to 5 URL inputs — add/remove buttons)
- category_id (select, domain=wishlist categories)
- priority (low/medium/high select)
- status (7-option select)
- target_date (date input, optional)
- monthly_contribution (number, optional)
- contribution_currency (MXN/USD select)

**Zod schema:**
```ts
const wishlistSchema = z.object({
  name: z.string().min(1, "Required"),
  image_url: z.string().optional(),
  price: z.string().optional(),
  currency: z.string().min(1),
  links: z.array(z.string().url("Must be a valid URL")).max(5).default([]),
  category_id: z.string().optional(),
  priority: z.string().min(1),
  status: z.string().min(1),
  target_date: z.string().optional(),
  monthly_contribution: z.string().optional(),
  contribution_currency: z.string().min(1),
});
```

**Dynamic links management:**
- Show existing links as inputs with a remove (×) button
- "Add Link" button (disabled when 5 links reached)
- Use `useFieldArray` from react-hook-form or manage manually with state

Props: `open, onClose, onSubmit, categories: Category[], defaultValues?: WishlistItem`

- [ ] **Step 1: Create wishlist form component**
- [ ] **Step 2: Verify build**
- [ ] **Step 3: Commit**

---

### Task 7: Gallery View Component

**Files:**
- Create: `frontend/src/components/wishlist/wishlist-gallery.tsx`

Large cards in a responsive grid (`grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4`).

Each card (shadcn Card):
- **Cover image** (if image_url set): `<img>` with aspect-ratio-video, object-cover. Placeholder div if no image.
- **Body:** Item name (bold), Price + currency, Category badge, Priority badge (color-coded: high=red, medium=yellow, low=green), Status badge
- **Footer:** Edit/Delete buttons
- Filtered to active items (To-do + In Progress statuses) by default, with option to show all
- Sorted by priority (high first)

Props: `items: WishlistItem[], onEdit, onDelete, onStatusChange`

- [ ] **Step 1: Create gallery component**
- [ ] **Step 2: Commit**

---

### Task 8: Table View Component

**Files:**
- Create: `frontend/src/components/wishlist/wishlist-table.tsx`

Full shadcn Table with ALL properties visible.

Columns: Name, Image (thumbnail), Price, Category, Priority (badge), Status (badge), Target Date, Monthly Contribution, Links (count), Actions (edit/delete)

Sorted by status group (To-do first, then In Progress, then Complete), then by priority within each group.

Props: `items: WishlistItem[], onEdit, onDelete, onStatusChange`

- [ ] **Step 1: Create table component**
- [ ] **Step 2: Commit**

---

### Task 9: Board View Component

**Files:**
- Create: `frontend/src/components/wishlist/wishlist-board.tsx`

Kanban-style board with columns grouped by **status group** (To-do, In Progress, Complete).

Each column:
- Header with group name and item count
- Cards showing: name, price, priority badge, category badge
- Status dropdown on each card to change status (triggers PATCH /api/v1/wishlist/:id/status)

Layout: 3 columns side-by-side, scrollable horizontally on small screens.

Within each column, items sorted by priority (high first).

Note: We're NOT implementing drag-and-drop for now — just dropdown status change. Drag-and-drop can be added later with a library like @dnd-kit.

Props: `items: WishlistItem[], onEdit, onDelete, onStatusChange`

- [ ] **Step 1: Create board component**
- [ ] **Step 2: Commit**

---

### Task 10: Wishlist Page + Sidebar Update

**Files:**
- Create: `frontend/src/app/wishlist/page.tsx`
- Create: `frontend/src/app/wishlist/layout.tsx`
- Modify: `frontend/src/components/layout/sidebar.tsx`

**Wishlist layout:** Since wishlist is NOT year-scoped, it needs its own layout. Reuse the same Sidebar + Header pattern but without the year context.

`frontend/src/app/wishlist/layout.tsx`:
```tsx
import { Sidebar } from "@/components/layout/sidebar";

export default function WishlistLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <header className="h-14 border-b bg-white flex items-center px-6">
          <h1 className="text-lg font-semibold">Wishlist</h1>
        </header>
        <main className="flex-1 p-6 bg-gray-50">{children}</main>
      </div>
    </div>
  );
}
```

**Wishlist page (`/wishlist/page.tsx`):**
- View toggle at top: Gallery | Table | Board (3 buttons/tabs, persisted in state)
- Loads all wishlist items: `apiGet<ListResponse<WishlistItem>>("/api/v1/wishlist")`
- Loads wishlist categories: `apiGet<ListResponse<Category>>("/api/v1/categories?domain=wishlist")`
- "Add Item" button opens WishlistForm
- Renders the selected view component
- CRUD: create (apiPost), update (apiPut), delete (apiDelete)
- Status change: calls `apiPatch("/api/v1/wishlist/${id}/status", { status })` — NOTE: need to add `apiPatch` to `lib/api.ts`

**Add `apiPatch` to `frontend/src/lib/api.ts`:**
```ts
export async function apiPatch<T>(path: string, body: unknown): Promise<T> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, {
    method: "PATCH",
    headers,
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "API error");
  }
  return res.json();
}
```

**Sidebar update:** Change wishlist link from `/wishlist` to use proper active state detection:
```tsx
// In sidebar.tsx, update the wishlist link:
const isWishlistActive = pathname.startsWith("/wishlist");
```

- [ ] **Step 1: Add apiPatch to api.ts**
- [ ] **Step 2: Create wishlist layout**
- [ ] **Step 3: Create wishlist page with view toggle**
- [ ] **Step 4: Update sidebar wishlist link with active state**
- [ ] **Step 5: Verify build**
- [ ] **Step 6: Commit**

---

## Phase C: Verification

### Task 11: End-to-End Verification

- [ ] **Step 1: Run all backend tests**: `cd backend && go test ./... -v -count=1`
- [ ] **Step 2: Verify backend compiles**: `cd backend && go build ./cmd/server`
- [ ] **Step 3: Verify frontend builds**: `cd frontend && npm run build`
- [ ] **Step 4: Verify all routes registered** (should include `/wishlist`)
- [ ] **Step 5: Final commit if needed**

---

## API Endpoint Summary (New in Plan 4)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/wishlist` | Create wishlist item |
| GET | `/api/v1/wishlist` | List items (optional `?status=` filter) |
| GET | `/api/v1/wishlist/:id` | Get single item |
| PUT | `/api/v1/wishlist/:id` | Update item |
| PATCH | `/api/v1/wishlist/:id/status` | Update status only |
| DELETE | `/api/v1/wishlist/:id` | Delete item |

## Frontend Routes (New in Plan 4)

| Route | Description |
|-------|-------------|
| `/wishlist` | Wishlist with Gallery/Table/Board views (not year-scoped) |
