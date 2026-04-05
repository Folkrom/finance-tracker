"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Debt, Category, PaymentMethod, ListResponse } from "@/types";
import { DebtTable } from "@/components/debt/debt-table";
import { DebtForm } from "@/components/debt/debt-form";
import { Button } from "@/components/ui/button";

export default function DebtPage() {
  const params = useParams();
  const year = params.year as string;
  const t = useTranslations("debt");
  const tCommon = useTranslations("common");

  const [debts, setDebts] = useState<Debt[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingDebt, setEditingDebt] = useState<Debt | undefined>(undefined);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [debtsRes, categoriesRes, paymentMethodsRes] = await Promise.all([
        apiGet<ListResponse<Debt>>(`/api/v1/years/${year}/debts`),
        apiGet<ListResponse<Category>>(`/api/v1/categories?domain=expense`),
        apiGet<ListResponse<PaymentMethod>>(`/api/v1/payment-methods`),
      ]);
      setDebts(debtsRes.data);
      setCategories(categoriesRes.data);
      setPaymentMethods(paymentMethodsRes.data);
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
    setEditingDebt(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (debt: Debt) => {
    setEditingDebt(debt);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingDebt(undefined);
  };

  const handleSubmit = async (values: {
    name: string;
    amount: string;
    currency: string;
    category_id?: string;
    payment_method_id?: string;
    date: string;
  }) => {
    try {
      if (editingDebt) {
        await apiPut<Debt>(
          `/api/v1/years/${year}/debts/${editingDebt.id}`,
          values
        );
        toast.success("Debt updated");
      } else {
        await apiPost<Debt>(`/api/v1/years/${year}/debts`, values);
        toast.success("Debt created");
      }
      handleClose();
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save debt");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/years/${year}/debts/${id}`);
      toast.success("Debt deleted");
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete debt");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="size-4 mr-2" />
          {t("addDebt")}
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {tCommon("loading")}
        </div>
      ) : (
        <div className="rounded-lg border bg-white">
          <DebtTable
            debts={debts}
            onEdit={handleOpenEdit}
            onDelete={handleDelete}
          />
        </div>
      )}

      <DebtForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        categories={categories}
        paymentMethods={paymentMethods}
        defaultValues={editingDebt}
      />
    </div>
  );
}
