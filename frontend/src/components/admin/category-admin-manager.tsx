"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Pencil, Trash2, Lock, Plus } from "lucide-react";
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
