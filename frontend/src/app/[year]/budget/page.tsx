"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Budget, BudgetLine, Category, ListResponse } from "@/types";
import { BudgetTable } from "@/components/budget/budget-table";
import { BudgetForm } from "@/components/budget/budget-form";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

export default function BudgetPage() {
  const params = useParams();
  const year = Number(params.year as string);
  const t = useTranslations("budget");
  const tCommon = useTranslations("common");

  const currentMonth = new Date().getMonth() + 1;

  const [selectedMonth, setSelectedMonth] = useState(currentMonth);
  const [lines, setLines] = useState<BudgetLine[]>([]);
  const [recurringBudgets, setRecurringBudgets] = useState<Budget[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingBudget, setEditingBudget] = useState<Budget | undefined>(undefined);

  const loadSummary = useCallback(async () => {
    try {
      const res = await apiGet<{ data: BudgetLine[] }>(
        `/api/v1/budgets?month=${selectedMonth}&year=${year}`
      );
      setLines(res.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load budgets");
    }
  }, [selectedMonth, year]);

  const loadRecurring = useCallback(async () => {
    try {
      const res = await apiGet<ListResponse<Budget>>(`/api/v1/budgets/recurring`);
      setRecurringBudgets(res.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load recurring budgets");
    }
  }, []);

  const loadCategories = useCallback(async () => {
    try {
      const res = await apiGet<ListResponse<Category>>(`/api/v1/categories?domain=expense`);
      setCategories(res.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load categories");
    }
  }, []);

  useEffect(() => {
    const fetchAll = async () => {
      setLoading(true);
      await Promise.all([loadSummary(), loadRecurring(), loadCategories()]);
      setLoading(false);
    };
    fetchAll();
  }, [loadSummary, loadRecurring, loadCategories]);

  const handleOpenCreate = () => {
    setEditingBudget(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (budget: Budget) => {
    setEditingBudget(budget);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingBudget(undefined);
  };

  const handleSubmit = async (values: {
    category_id: string;
    monthly_limit: string;
    is_recurring: boolean;
    month: number;
    year: number;
  }) => {
    try {
      if (editingBudget) {
        await apiPut<Budget>(`/api/v1/budgets/${editingBudget.id}`, values);
        toast.success("Budget updated");
      } else {
        await apiPost<Budget>(`/api/v1/budgets`, values);
        toast.success("Budget created");
      }
      handleClose();
      await Promise.all([loadSummary(), loadRecurring()]);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save budget");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/budgets/${id}`);
      toast.success("Budget deleted");
      await Promise.all([loadSummary(), loadRecurring()]);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete budget");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="size-4 mr-2" />
          {t("addBudget")}
        </Button>
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm font-medium">{t("month")}</span>
        <Select
          value={String(selectedMonth)}
          onValueChange={(val) => {
            if (val !== null) setSelectedMonth(Number(val));
          }}
        >
          <SelectTrigger className="w-40">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {MONTH_NAMES.map((name, idx) => (
              <SelectItem key={idx + 1} value={String(idx + 1)}>
                {name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {tCommon("loading")}
        </div>
      ) : (
        <>
          <section className="space-y-3">
            <h2 className="text-lg font-semibold">{t("override")}</h2>
            <div className="rounded-lg border bg-white">
              <BudgetTable
                lines={lines}
                onEdit={handleOpenEdit}
                onDelete={handleDelete}
              />
            </div>
          </section>

          <section className="space-y-3">
            <h2 className="text-lg font-semibold">{t("recurring")}</h2>
            <div className="rounded-lg border bg-white">
              <BudgetTable
                lines={recurringBudgets.map((b) => ({
                  budget: b,
                  spent: "0",
                  remaining: b.monthly_limit,
                }))}
                onEdit={handleOpenEdit}
                onDelete={handleDelete}
              />
            </div>
          </section>
        </>
      )}

      <BudgetForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        categories={categories}
        defaultValues={editingBudget}
        currentMonth={selectedMonth}
        currentYear={year}
      />
    </div>
  );
}
