"use client";

import { useState, useEffect, useCallback } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { Category, PaymentMethod, ListResponse } from "@/types";
import { CategoryManager } from "@/components/settings/category-manager";
import { PaymentMethodManager } from "@/components/settings/payment-method-manager";
import { Separator } from "@/components/ui/separator";

export default function SettingsPage() {
  const t = useTranslations("settings");
  const tCommon = useTranslations("common");

  const [categories, setCategories] = useState<Category[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [incomeRes, expenseRes, pmRes] = await Promise.all([
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=income"),
        apiGet<ListResponse<Category>>("/api/v1/categories?domain=expense"),
        apiGet<ListResponse<PaymentMethod>>("/api/v1/payment-methods"),
      ]);
      setCategories([...incomeRes.data, ...expenseRes.data]);
      setPaymentMethods(pmRes.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load settings");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {tCommon("loading")}
      </div>
    );
  }

  return (
    <div className="space-y-8 max-w-3xl">
      <h1 className="text-2xl font-bold">{t("title")}</h1>

      <CategoryManager categories={categories} onRefresh={loadData} />

      <Separator />

      <PaymentMethodManager paymentMethods={paymentMethods} onRefresh={loadData} />
    </div>
  );
}
