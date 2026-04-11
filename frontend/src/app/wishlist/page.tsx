"use client";

import { useState, useEffect, useCallback } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiPatch, apiDelete } from "@/lib/api";
import {
  WishlistItem,
  WishlistStatus,
  Category,
  ListResponse,
} from "@/types";
import { WishlistForm } from "@/components/wishlist/wishlist-form";
import { WishlistGallery } from "@/components/wishlist/wishlist-gallery";
import { WishlistTable } from "@/components/wishlist/wishlist-table";
import { WishlistBoard } from "@/components/wishlist/wishlist-board";
import { Button } from "@/components/ui/button";

type ViewMode = "gallery" | "table" | "board";

export default function WishlistPage() {
  const t = useTranslations("wishlist");
  const tCommon = useTranslations("common");

  const [items, setItems] = useState<WishlistItem[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<WishlistItem | undefined>(undefined);
  const [view, setView] = useState<ViewMode>("gallery");

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [itemsRes, categoriesRes] = await Promise.all([
        apiGet<ListResponse<WishlistItem>>("/api/v1/wishlist"),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=wishlist"),
      ]);
      setItems(itemsRes.data);
      setCategories(categoriesRes.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load data");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleOpenCreate = () => {
    setEditingItem(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (item: WishlistItem) => {
    setEditingItem(item);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingItem(undefined);
  };

  const handleSubmit = async (values: Record<string, unknown>) => {
    try {
      if (editingItem) {
        await apiPut<WishlistItem>(`/api/v1/wishlist/${editingItem.id}`, values);
        toast.success("Item updated");
      } else {
        await apiPost<WishlistItem>("/api/v1/wishlist", values);
        toast.success("Item created");
      }
      handleClose();
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save item");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/wishlist/${id}`);
      toast.success("Item deleted");
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete item");
    }
  };

  const handleStatusChange = async (id: string, status: WishlistStatus) => {
    try {
      await apiPatch(`/api/v1/wishlist/${id}/status`, { status });
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update status");
    }
  };

  const viewButtons: { key: ViewMode; label: string }[] = [
    { key: "gallery", label: t("viewGallery") },
    { key: "table", label: t("viewTable") },
    { key: "board", label: t("viewBoard") },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <div className="flex items-center gap-3">
          <div className="flex rounded-lg border bg-white">
            {viewButtons.map((vb) => (
              <button
                key={vb.key}
                onClick={() => setView(vb.key)}
                className={`px-3 py-1.5 text-sm transition-colors ${
                  view === vb.key
                    ? "bg-gray-100 text-gray-900 font-medium"
                    : "text-gray-600 hover:text-gray-900"
                }`}
              >
                {vb.label}
              </button>
            ))}
          </div>
          <Button onClick={handleOpenCreate}>
            <Plus className="size-4 mr-2" />
            {t("addItem")}
          </Button>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {tCommon("loading")}
        </div>
      ) : (
        <>
          {view === "gallery" && (
            <WishlistGallery
              items={items}
              onEdit={handleOpenEdit}
              onDelete={handleDelete}
            />
          )}
          {view === "table" && (
            <div className="rounded-lg border bg-white">
              <WishlistTable
                items={items}
                onEdit={handleOpenEdit}
                onDelete={handleDelete}
              />
            </div>
          )}
          {view === "board" && (
            <WishlistBoard
              items={items}
              onEdit={handleOpenEdit}
              onDelete={handleDelete}
              onStatusChange={handleStatusChange}
            />
          )}
        </>
      )}

      <WishlistForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        categories={categories}
        defaultValues={editingItem}
      />
    </div>
  );
}
