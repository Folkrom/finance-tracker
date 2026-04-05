"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Expense, Category, PaymentMethod, ListResponse } from "@/types";
import { ExpenseTable } from "@/components/expenses/expense-table";
import { ExpenseForm } from "@/components/expenses/expense-form";
import { Button } from "@/components/ui/button";

export default function ExpensesPage() {
  const params = useParams();
  const year = params.year as string;
  const t = useTranslations("expenses");
  const tCommon = useTranslations("common");

  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingExpense, setEditingExpense] = useState<Expense | undefined>(undefined);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [expensesRes, categoriesRes, paymentMethodsRes] = await Promise.all([
        apiGet<ListResponse<Expense>>(`/api/v1/years/${year}/expenses`),
        apiGet<ListResponse<Category>>(`/api/v1/categories?domain=expense`),
        apiGet<ListResponse<PaymentMethod>>(`/api/v1/payment-methods`),
      ]);
      setExpenses(expensesRes.data);
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
    setEditingExpense(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (expense: Expense) => {
    setEditingExpense(expense);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingExpense(undefined);
  };

  const handleSubmit = async (values: {
    name: string;
    amount: string;
    currency: string;
    category_id?: string;
    payment_method_id?: string;
    type: "expense" | "saving" | "investment";
    date: string;
  }) => {
    try {
      if (editingExpense) {
        await apiPut<Expense>(
          `/api/v1/years/${year}/expenses/${editingExpense.id}`,
          values
        );
        toast.success("Expense updated");
      } else {
        await apiPost<Expense>(`/api/v1/years/${year}/expenses`, values);
        toast.success("Expense created");
      }
      handleClose();
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save expense");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/years/${year}/expenses/${id}`);
      toast.success("Expense deleted");
      await loadData();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete expense");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="size-4 mr-2" />
          {t("addExpense")}
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {tCommon("loading")}
        </div>
      ) : (
        <div className="rounded-lg border bg-white">
          <ExpenseTable
            expenses={expenses}
            onEdit={handleOpenEdit}
            onDelete={handleDelete}
          />
        </div>
      )}

      <ExpenseForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        categories={categories}
        paymentMethods={paymentMethods}
        defaultValues={editingExpense}
      />
    </div>
  );
}
