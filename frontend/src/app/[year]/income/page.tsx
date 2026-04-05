"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Income, Category, ListResponse } from "@/types";
import { IncomeTable } from "@/components/income/income-table";
import { IncomeForm } from "@/components/income/income-form";
import { Button } from "@/components/ui/button";

export default function IncomePage() {
  const params = useParams();
  const year = params.year as string;
  const t = useTranslations("income");
  const tCommon = useTranslations("common");

  const [incomes, setIncomes] = useState<Income[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingIncome, setEditingIncome] = useState<Income | undefined>(undefined);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [incomesRes, categoriesRes] = await Promise.all([
        apiGet<ListResponse<Income>>(`/api/v1/years/${year}/incomes`),
        apiGet<ListResponse<Category>>(`/api/v1/categories?domain=income`),
      ]);
      setIncomes(incomesRes.data);
      setCategories(categoriesRes.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load data");
    } finally {
      setLoading(false);
    }
  }, [year]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleOpenCreate = () => {
    setEditingIncome(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (income: Income) => {
    setEditingIncome(income);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingIncome(undefined);
  };

  const handleSubmit = async (values: {
    source: string;
    amount: string;
    currency: string;
    category_id?: string;
    date: string;
  }) => {
    try {
      if (editingIncome) {
        await apiPut<Income>(
          `/api/v1/years/${year}/incomes/${editingIncome.id}`,
          values
        );
        toast.success("Income updated");
      } else {
        await apiPost<Income>(`/api/v1/years/${year}/incomes`, values);
        toast.success("Income created");
      }
      handleClose();
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save income");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/years/${year}/incomes/${id}`);
      toast.success("Income deleted");
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete income");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="size-4 mr-2" />
          {t("addIncome")}
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {tCommon("loading")}
        </div>
      ) : (
        <div className="rounded-lg border bg-white">
          <IncomeTable
            incomes={incomes}
            onEdit={handleOpenEdit}
            onDelete={handleDelete}
          />
        </div>
      )}

      <IncomeForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        categories={categories}
        defaultValues={editingIncome}
      />
    </div>
  );
}
