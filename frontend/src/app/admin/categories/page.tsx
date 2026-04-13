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
