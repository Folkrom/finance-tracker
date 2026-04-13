# Admin Dashboard Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add admin dashboard frontend with stats overview and global category CRUD, gated by JWT admin role detection.

**Architecture:** `useAdmin` hook reads `app_metadata.role` from Supabase session. Separate `/admin` route tree with its own layout (sidebar + header). Stats page shows platform counts. Category admin page reuses the grouped-by-domain pattern from Settings but with full CRUD on globals.

**Tech Stack:** Next.js 16, React, TypeScript, shadcn/ui, Supabase Auth, lucide-react icons

---

### File Structure

| File | Purpose |
|------|---------|
| Create: `frontend/src/hooks/use-admin.ts` | Hook to check admin role from JWT claims |
| Create: `frontend/src/types/admin.ts` | AdminStats interface |
| Modify: `frontend/src/types/index.ts` | Re-export admin types |
| Create: `frontend/src/components/admin/admin-sidebar.tsx` | Admin nav sidebar |
| Create: `frontend/src/components/admin/admin-header.tsx` | Header with back link + logout |
| Create: `frontend/src/components/admin/stats-cards.tsx` | Stat card grid |
| Create: `frontend/src/components/admin/category-admin-manager.tsx` | Global category CRUD |
| Create: `frontend/src/app/admin/layout.tsx` | Admin layout with sidebar/header + auth gate |
| Create: `frontend/src/app/admin/page.tsx` | Redirect to /admin/stats |
| Create: `frontend/src/app/admin/stats/page.tsx` | Stats page |
| Create: `frontend/src/app/admin/categories/page.tsx` | Categories page |
| Modify: `frontend/src/components/layout/header.tsx` | Add admin link for admin users |

---

### Task 1: AdminStats type + useAdmin hook

**Files:**
- Create: `frontend/src/types/admin.ts`
- Modify: `frontend/src/types/index.ts`
- Create: `frontend/src/hooks/use-admin.ts`

- [ ] **Step 1: Create AdminStats type**

Create `frontend/src/types/admin.ts`:

```typescript
export interface AdminStats {
  users: number;
  categories_global: number;
  categories_user: number;
  profiles: number;
}
```

- [ ] **Step 2: Re-export from types/index.ts**

Add to the end of `frontend/src/types/index.ts`:

```typescript
export type { AdminStats } from "./admin";
```

- [ ] **Step 3: Create useAdmin hook**

Create `frontend/src/hooks/use-admin.ts`:

```typescript
"use client";

import { useState, useEffect } from "react";
import { createClient } from "@/lib/supabase/client";

export function useAdmin() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function check() {
      const supabase = createClient();
      const { data: { session } } = await supabase.auth.getSession();
      const role = session?.user?.app_metadata?.role;
      setIsAdmin(role === "admin");
      setLoading(false);
    }
    check();
  }, []);

  return { isAdmin, loading };
}
```

- [ ] **Step 4: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/types/admin.ts frontend/src/types/index.ts frontend/src/hooks/use-admin.ts
git commit -m "feat: AdminStats type and useAdmin hook"
```

---

### Task 2: Admin sidebar and header components

**Files:**
- Create: `frontend/src/components/admin/admin-sidebar.tsx`
- Create: `frontend/src/components/admin/admin-header.tsx`

- [ ] **Step 1: Create admin sidebar**

Create `frontend/src/components/admin/admin-sidebar.tsx`:

```typescript
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { BarChart3, Tags } from "lucide-react";

const navItems = [
  { label: "Stats", path: "/admin/stats", icon: BarChart3 },
  { label: "Categories", path: "/admin/categories", icon: Tags },
];

export function AdminSidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-64 border-r bg-white h-screen sticky top-0 flex flex-col">
      <div className="p-6">
        <h1 className="text-xl font-bold">Admin Panel</h1>
      </div>
      <nav className="flex-1 px-3">
        {navItems.map((item) => {
          const isActive = pathname.startsWith(item.path);
          return (
            <Link
              key={item.path}
              href={item.path}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors mb-1",
                isActive
                  ? "bg-gray-100 text-gray-900 font-medium"
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900"
              )}
            >
              <item.icon className="size-4" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
```

- [ ] **Step 2: Create admin header**

Create `frontend/src/components/admin/admin-header.tsx`:

```typescript
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

export function AdminHeader() {
  const router = useRouter();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  const currentYear = new Date().getFullYear();

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <Link
        href={`/${currentYear}/dashboard`}
        className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="size-4" />
        Back to app
      </Link>
      <Button variant="ghost" size="sm" onClick={handleLogout}>
        Logout
      </Button>
    </header>
  );
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/admin/admin-sidebar.tsx frontend/src/components/admin/admin-header.tsx
git commit -m "feat: admin sidebar and header components"
```

---

### Task 3: Admin layout + redirect page

**Files:**
- Create: `frontend/src/app/admin/layout.tsx`
- Create: `frontend/src/app/admin/page.tsx`

- [ ] **Step 1: Create admin layout**

Create `frontend/src/app/admin/layout.tsx`:

```typescript
"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAdmin } from "@/hooks/use-admin";
import { AdminSidebar } from "@/components/admin/admin-sidebar";
import { AdminHeader } from "@/components/admin/admin-header";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { isAdmin, loading } = useAdmin();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAdmin) {
      const currentYear = new Date().getFullYear();
      router.replace(`/${currentYear}/dashboard`);
    }
  }, [loading, isAdmin, router]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen text-muted-foreground">
        Loading...
      </div>
    );
  }

  if (!isAdmin) {
    return null;
  }

  return (
    <div className="flex min-h-screen">
      <AdminSidebar />
      <div className="flex-1 flex flex-col">
        <AdminHeader />
        <main className="flex-1 p-6 bg-gray-50">{children}</main>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Create redirect page**

Create `frontend/src/app/admin/page.tsx`:

```typescript
import { redirect } from "next/navigation";

export default function AdminIndexPage() {
  redirect("/admin/stats");
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/app/admin/layout.tsx frontend/src/app/admin/page.tsx
git commit -m "feat: admin layout with auth gate and redirect"
```

---

### Task 4: Stats cards component + stats page

**Files:**
- Create: `frontend/src/components/admin/stats-cards.tsx`
- Create: `frontend/src/app/admin/stats/page.tsx`

- [ ] **Step 1: Create stats cards component**

Create `frontend/src/components/admin/stats-cards.tsx`:

```typescript
"use client";

import { Users, Layers, UserPlus, Globe } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AdminStats } from "@/types";

interface StatsCardsProps {
  stats: AdminStats;
}

const statConfig = [
  { key: "users" as const, label: "Users", icon: Users },
  { key: "profiles" as const, label: "Profiles", icon: UserPlus },
  { key: "categories_global" as const, label: "Global Categories", icon: Globe },
  { key: "categories_user" as const, label: "User Categories", icon: Layers },
];

export function StatsCards({ stats }: StatsCardsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {statConfig.map((item) => (
        <Card key={item.key}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{item.label}</CardTitle>
            <item.icon className="size-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats[item.key]}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
```

- [ ] **Step 2: Create stats page**

Create `frontend/src/app/admin/stats/page.tsx`:

```typescript
"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { AdminStats } from "@/types";
import { StatsCards } from "@/components/admin/stats-cards";

export default function AdminStatsPage() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const data = await apiGet<AdminStats>("/api/v1/admin/stats");
        setStats(data);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to load stats");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        Loading...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Platform Stats</h1>
      {stats && <StatsCards stats={stats} />}
    </div>
  );
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/admin/stats-cards.tsx frontend/src/app/admin/stats/page.tsx
git commit -m "feat: admin stats page with stat cards"
```

---

### Task 5: Category admin manager component

**Files:**
- Create: `frontend/src/components/admin/category-admin-manager.tsx`

- [ ] **Step 1: Create category admin manager**

Create `frontend/src/components/admin/category-admin-manager.tsx`:

```typescript
"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Pencil, Trash2, Lock, Plus, X, Check } from "lucide-react";
import { Category } from "@/types";
import { apiPost, apiPut, apiDelete } from "@/lib/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";

interface CategoryAdminManagerProps {
  categories: Category[];
  onRefresh: () => void;
}

type Domain = "income" | "expense" | "wishlist";

const DOMAINS: Domain[] = ["income", "expense", "wishlist"];

export function CategoryAdminManager({
  categories,
  onRefresh,
}: CategoryAdminManagerProps) {
  const [creating, setCreating] = useState<Domain | null>(null);
  const [newName, setNewName] = useState("");
  const [newColor, setNewColor] = useState("");
  const [newSortOrder, setNewSortOrder] = useState(0);

  const [editingCat, setEditingCat] = useState<Category | null>(null);
  const [editName, setEditName] = useState("");
  const [editColor, setEditColor] = useState("");
  const [editSortOrder, setEditSortOrder] = useState(0);

  const [deletingCat, setDeletingCat] = useState<Category | null>(null);

  const grouped = DOMAINS.reduce<Record<Domain, Category[]>>(
    (acc, domain) => {
      acc[domain] = categories.filter((c) => c.domain === domain);
      return acc;
    },
    { income: [], expense: [], wishlist: [] }
  );

  const handleCreate = async () => {
    if (!newName.trim() || !creating) return;
    try {
      await apiPost("/api/v1/admin/categories", {
        name: newName.trim(),
        domain: creating,
        color: newColor || undefined,
        sort_order: newSortOrder,
      });
      toast.success("Category created");
      setCreating(null);
      setNewName("");
      setNewColor("");
      setNewSortOrder(0);
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create");
    }
  };

  const handleStartEdit = (cat: Category) => {
    setEditingCat(cat);
    setEditName(cat.name);
    setEditColor(cat.color || "");
    setEditSortOrder(cat.sort_order);
  };

  const handleConfirmEdit = async () => {
    if (!editingCat || !editName.trim()) return;
    try {
      await apiPut(`/api/v1/admin/categories/${editingCat.id}`, {
        name: editName.trim(),
        color: editColor || undefined,
        sort_order: editSortOrder,
      });
      toast.success("Category updated");
      setEditingCat(null);
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update");
    }
  };

  const handleDelete = async () => {
    if (!deletingCat) return;
    try {
      await apiDelete(`/api/v1/admin/categories/${deletingCat.id}`);
      toast.success("Category deleted");
      setDeletingCat(null);
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete");
    }
  };

  return (
    <div className="space-y-6">
      {DOMAINS.map((domain) => (
        <div key={domain}>
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-medium text-muted-foreground capitalize">
              {domain}
            </h3>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setCreating(domain);
                setNewName("");
                setNewColor("");
                setNewSortOrder(
                  grouped[domain].length > 0
                    ? Math.max(...grouped[domain].map((c) => c.sort_order)) + 1
                    : 0
                );
              }}
            >
              <Plus className="size-3 mr-1" />
              Add
            </Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {grouped[domain].length === 0 ? (
              <span className="text-sm text-muted-foreground">None</span>
            ) : (
              grouped[domain].map((cat) => (
                <div key={cat.id} className="flex items-center gap-0.5">
                  <Badge variant="secondary">
                    {cat.is_system && (
                      <Lock className="size-3 mr-1 opacity-50" />
                    )}
                    {cat.color && (
                      <span
                        className="inline-block size-2 rounded-full mr-1"
                        style={{ backgroundColor: cat.color }}
                      />
                    )}
                    {cat.name}
                    <span className="ml-1 text-xs opacity-50">
                      #{cat.sort_order}
                    </span>
                  </Badge>
                  {!cat.is_system && (
                    <>
                      <button
                        onClick={() => handleStartEdit(cat)}
                        className="text-muted-foreground hover:text-foreground transition-colors"
                        aria-label="Edit"
                      >
                        <Pencil className="size-3" />
                      </button>
                      <button
                        onClick={() => setDeletingCat(cat)}
                        className="text-muted-foreground hover:text-destructive transition-colors"
                        aria-label="Delete"
                      >
                        <Trash2 className="size-3" />
                      </button>
                    </>
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      ))}

      {/* Create dialog */}
      <Dialog open={creating !== null} onOpenChange={() => setCreating(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              Create {creating} category
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <Input
              placeholder="Name"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleCreate();
              }}
              autoFocus
            />
            <Input
              placeholder="Color (#hex, optional)"
              value={newColor}
              onChange={(e) => setNewColor(e.target.value)}
            />
            <Input
              type="number"
              placeholder="Sort order"
              value={newSortOrder}
              onChange={(e) => setNewSortOrder(parseInt(e.target.value) || 0)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreating(null)}>
              Cancel
            </Button>
            <Button onClick={handleCreate} disabled={!newName.trim()}>
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit dialog */}
      <Dialog
        open={editingCat !== null}
        onOpenChange={() => setEditingCat(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit category</DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <Input
              placeholder="Name"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleConfirmEdit();
              }}
              autoFocus
            />
            <Input
              placeholder="Color (#hex, optional)"
              value={editColor}
              onChange={(e) => setEditColor(e.target.value)}
            />
            <Input
              type="number"
              placeholder="Sort order"
              value={editSortOrder}
              onChange={(e) => setEditSortOrder(parseInt(e.target.value) || 0)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingCat(null)}>
              Cancel
            </Button>
            <Button onClick={handleConfirmEdit} disabled={!editName.trim()}>
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete confirmation dialog */}
      <Dialog
        open={deletingCat !== null}
        onOpenChange={() => setDeletingCat(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete &quot;{deletingCat?.name}&quot;?</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            All references to this category will be reassigned to the
            &quot;Other&quot; category for the {deletingCat?.domain} domain.
            This cannot be undone.
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeletingCat(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/admin/category-admin-manager.tsx
git commit -m "feat: admin category manager with CRUD dialogs"
```

---

### Task 6: Admin categories page

**Files:**
- Create: `frontend/src/app/admin/categories/page.tsx`

- [ ] **Step 1: Create categories page**

Create `frontend/src/app/admin/categories/page.tsx`:

```typescript
"use client";

import { useState, useEffect, useCallback } from "react";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { Category, ListResponse } from "@/types";
import { CategoryAdminManager } from "@/components/admin/category-admin-manager";

export default function AdminCategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  const loadCategories = useCallback(async () => {
    try {
      const [income, expense, wishlist] = await Promise.all([
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=income"),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=expense"),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=wishlist"),
      ]);
      setCategories([...income.data, ...expense.data, ...wishlist.data]);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load categories");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        Loading...
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-3xl">
      <h1 className="text-2xl font-bold">Global Categories</h1>
      <CategoryAdminManager categories={categories} onRefresh={loadCategories} />
    </div>
  );
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/app/admin/categories/page.tsx
git commit -m "feat: admin categories page"
```

---

### Task 7: Add admin link to main app header

**Files:**
- Modify: `frontend/src/components/layout/header.tsx`

- [ ] **Step 1: Update header with admin link**

Replace the full content of `frontend/src/components/layout/header.tsx`:

```typescript
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { YearSwitcher } from "./year-switcher";
import { useTranslations } from "next-intl";
import { useAdmin } from "@/hooks/use-admin";
import { Shield } from "lucide-react";

export function Header() {
  const t = useTranslations("auth");
  const router = useRouter();
  const { isAdmin } = useAdmin();

  async function handleLogout() {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
    router.refresh();
  }

  return (
    <header className="h-14 border-b bg-white flex items-center justify-between px-6">
      <YearSwitcher />
      <div className="flex items-center gap-2">
        {isAdmin && (
          <Link href="/admin">
            <Button variant="ghost" size="sm" className="gap-1">
              <Shield className="size-4" />
              Admin
            </Button>
          </Link>
        )}
        <Button variant="ghost" size="sm" onClick={handleLogout}>
          {t("logout")}
        </Button>
      </div>
    </header>
  );
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd /home/folkrom/projects/finance-tracker/frontend && npx tsc --noEmit
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/layout/header.tsx
git commit -m "feat: admin link in header for admin users"
```

---

### Task 8: Visual verification

**Files:** None (manual testing)

- [ ] **Step 1: Start dev servers**

```bash
cd /home/folkrom/projects/finance-tracker && mise run dev-backend
# In another terminal:
cd /home/folkrom/projects/finance-tracker/frontend && npm run dev
```

- [ ] **Step 2: Verify non-admin user**

Log in as a regular user. Verify:
- No "Admin" button in the header
- Navigating to `/admin` redirects to `/{year}/dashboard`

- [ ] **Step 3: Set admin role**

In Supabase dashboard: Authentication → Users → select your user → Edit → `app_metadata`: `{"role": "admin"}`. Save. Log out and back in.

- [ ] **Step 4: Verify admin user**

Log in as admin. Verify:
- "Admin" shield button appears in the header
- Clicking it goes to `/admin/stats`
- Stats page shows four stat cards with counts
- Categories page shows all global categories grouped by domain
- Can create a new category via the Add button + dialog
- Can edit a non-system category via the pencil icon + dialog
- Can delete a non-system category via the trash icon + confirmation dialog
- System categories (Other) show lock icon, no edit/delete buttons
- "Back to app" link in admin header returns to main app
