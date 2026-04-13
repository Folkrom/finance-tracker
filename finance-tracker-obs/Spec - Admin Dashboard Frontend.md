# Admin Dashboard Frontend Design

**Date:** 2026-04-12
**Status:** Approved
**Scope:** Admin frontend pages — stats overview and global category CRUD. No user management.

---

## Problem

The backend has admin routes (stats, global category CRUD) but no frontend to use them. Admins must use curl or Postman to manage global categories.

## Design

### Admin Detection (Client-Side)

Read `app_metadata.role` from the Supabase session JWT. This is a UX convenience — the backend enforces the real 403.

**Hook: `useAdmin()`**
- Reads session via `supabase.auth.getSession()`
- Extracts `session.user.app_metadata.role`
- Returns `{ isAdmin: boolean, loading: boolean }`
- Used in the header to show/hide admin link
- Used in admin layout to redirect non-admins to `/`

### Routing & Layout

Separate `/admin` route tree, outside the `[year]` layout. No year scoping needed.

```
/admin              → redirects to /admin/stats
/admin/stats        → platform stats overview
/admin/categories   → global category CRUD
```

**Admin layout** (`app/admin/layout.tsx`):
- Admin-specific sidebar with Stats and Categories nav items
- Simple header with "Back to app" link and logout button
- Uses `useAdmin()` to redirect non-admins
- Same sidebar visual style as main app (consistency)

### Header Change

Add an admin link to the existing main app header (`components/layout/header.tsx`):
- Shield icon + "Admin" text
- Only visible when `useAdmin().isAdmin === true`
- Links to `/admin`
- Positioned before the logout button

### Stats Page

**Route:** `/admin/stats`

Calls `GET /api/v1/admin/stats` which returns:
```json
{
  "users": 42,
  "categories_global": 28,
  "categories_user": 5,
  "profiles": 42
}
```

Displays four stat cards in a 2x2 grid using shadcn Card components:
- Users (from `users`)
- Profiles (from `profiles`)
- Global Categories (from `categories_global`)
- User Categories (from `categories_user`)

Each card shows: icon, label, count. Simple and informational.

### Categories Page

**Route:** `/admin/categories`

Same visual pattern as the existing Settings category manager but with full admin capabilities.

**Layout:** Table grouped by domain (income, expense, wishlist). Each domain section has:
- Domain heading with "Add Category" button
- Table rows: name, color swatch, sort order, system badge, edit/delete actions

**Create:** Dialog with fields: name, domain (pre-filled from section), color (optional), sort_order.
Calls `POST /api/v1/admin/categories`.
Returns 201 with created category.

**Edit:** Dialog pre-filled with current values: name, color, sort_order.
Calls `PUT /api/v1/admin/categories/:id`.
Domain is not editable (changing domain would break FK references).

**Delete:** Confirmation dialog that warns: "This will reassign all references to the 'Other' category for this domain."
Calls `DELETE /api/v1/admin/categories/:id`.
Returns 204.

**Guards:**
- System categories (`is_system = true`) show a lock icon, no edit/delete buttons
- Delete button disabled for system categories even if somehow rendered

### New Types

```typescript
interface AdminStats {
  users: number;
  categories_global: number;
  categories_user: number;
  profiles: number;
}
```

Added to `types/index.ts`.

### Components

| Component | File | Purpose |
|-----------|------|---------|
| AdminSidebar | `components/admin/admin-sidebar.tsx` | Nav links for admin pages |
| AdminHeader | `components/admin/admin-header.tsx` | Back to app + logout |
| StatsCards | `components/admin/stats-cards.tsx` | Stat card grid |
| CategoryAdminManager | `components/admin/category-admin-manager.tsx` | Full CRUD table for globals |
| useAdmin | `hooks/use-admin.ts` | JWT claim check hook |

### What This Does NOT Include

- User management pages (no backend routes yet)
- Audit logging UI
- Admin-specific i18n keys (use English hardcoded for now — admin is single-user dev mode)
- Role management UI (admin role set via Supabase dashboard)
