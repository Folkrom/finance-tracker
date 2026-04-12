"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { X, Check, Pencil, Lock } from "lucide-react";
import { Category } from "@/types";
import { apiPost, apiPut, apiDelete } from "@/lib/api";
import { Badge } from "@/components/ui/badge";
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

interface CategoryManagerProps {
  categories: Category[];
  onRefresh: () => void;
}

type Domain = "income" | "expense" | "wishlist";

const DOMAINS: Domain[] = ["income", "expense", "wishlist"];

export function CategoryManager({ categories, onRefresh }: CategoryManagerProps) {
  const t = useTranslations("settings");
  const tCommon = useTranslations("common");

  const [newName, setNewName] = useState("");
  const [newDomain, setNewDomain] = useState<Domain>("income");
  const [creating, setCreating] = useState(false);

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState("");

  const handleCreate = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    try {
      await apiPost<Category>("/api/v1/categories", {
        name: newName.trim(),
        domain: newDomain,
      });
      setNewName("");
      toast.success("Category created");
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create category");
    } finally {
      setCreating(false);
    }
  };

  const handleStartEdit = (cat: Category) => {
    setEditingId(cat.id);
    setEditingName(cat.name);
  };

  const handleConfirmEdit = async (cat: Category) => {
    if (!editingName.trim()) return;
    try {
      await apiPut<Category>(`/api/v1/categories/${cat.id}`, {
        name: editingName.trim(),
        domain: cat.domain,
      });
      toast.success("Category updated");
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update category");
    } finally {
      setEditingId(null);
      setEditingName("");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/categories/${id}`);
      toast.success("Category deleted");
      onRefresh();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete category");
    }
  };

  const grouped = DOMAINS.reduce<Record<Domain, Category[]>>(
    (acc, domain) => {
      acc[domain] = categories.filter((c) => c.domain === domain);
      return acc;
    },
    { income: [], expense: [], wishlist: [] }
  );

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold mb-4">{t("categories")}</h2>

        {/* Add form */}
        <div className="flex gap-2 items-end flex-wrap">
          <div className="space-y-1">
            <Label htmlFor="cat-name">Name</Label>
            <Input
              id="cat-name"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder="Category name"
              className="w-48"
              onKeyDown={(e) => { if (e.key === "Enter") handleCreate(); }}
            />
          </div>

          <div className="space-y-1">
            <Label>Domain</Label>
            <Select
              value={newDomain}
              onValueChange={(val) => {
                if (val !== null) setNewDomain(val as Domain);
              }}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {DOMAINS.map((d) => (
                  <SelectItem key={d} value={d}>
                    {d}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <Button onClick={handleCreate} disabled={creating || !newName.trim()}>
            {tCommon("create")}
          </Button>
        </div>
      </div>

      {/* Categories grouped by domain */}
      <div className="space-y-4">
        {DOMAINS.map((domain) => (
          <div key={domain}>
            <p className="text-sm font-medium text-muted-foreground capitalize mb-2">
              {domain}
            </p>
            <div className="flex flex-wrap gap-2">
              {grouped[domain].length === 0 ? (
                <span className="text-sm text-muted-foreground">None</span>
              ) : (
                grouped[domain].map((cat) => {
                  const isGlobal = cat.user_id === null;

                  if (isGlobal) {
                    return (
                      <div key={cat.id} className="flex items-center gap-0.5">
                        <Badge variant="secondary">
                          <Lock className="size-3 mr-1 opacity-50" />
                          {cat.name}
                        </Badge>
                      </div>
                    );
                  }

                  if (editingId === cat.id) {
                    return (
                      <div key={cat.id} className="flex items-center gap-1">
                        <Input
                          value={editingName}
                          onChange={(e) => setEditingName(e.target.value)}
                          className="h-7 w-32 text-sm"
                          onKeyDown={(e) => {
                            if (e.key === "Enter") handleConfirmEdit(cat);
                            if (e.key === "Escape") setEditingId(null);
                          }}
                          autoFocus
                        />
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0"
                          onClick={() => handleConfirmEdit(cat)}
                        >
                          <Check className="size-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0"
                          onClick={() => setEditingId(null)}
                        >
                          <X className="size-3" />
                        </Button>
                      </div>
                    );
                  }

                  return (
                    <div key={cat.id} className="flex items-center gap-0.5">
                      <Badge
                        variant="secondary"
                        className="cursor-pointer"
                        onClick={() => handleStartEdit(cat)}
                      >
                        <Pencil className="size-3 mr-1 opacity-50" />
                        {cat.name}
                      </Badge>
                      <button
                        onClick={() => handleDelete(cat.id)}
                        className="ml-0.5 text-muted-foreground hover:text-destructive transition-colors"
                        aria-label="Delete"
                      >
                        <X className="size-3" />
                      </button>
                    </div>
                  );
                })
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
