"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api";
import { Card, CardSummary, PaymentMethod, ListResponse } from "@/types";
import { CardHealthCard } from "@/components/cards/card-health-card";
import { CardForm } from "@/components/cards/card-form";
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

interface CardFormValues {
  payment_method_id: string;
  bank: string;
  card_limit: string;
  recommended_max_pct: string;
  manual_usage_override?: string;
  level?: string;
}

export default function CardsPage() {
  const params = useParams();
  const year = Number(params.year as string);
  const t = useTranslations("cards");
  const tCommon = useTranslations("common");

  const currentMonth = new Date().getMonth() + 1;

  const [selectedMonth, setSelectedMonth] = useState(currentMonth);
  const [summaries, setSummaries] = useState<CardSummary[]>([]);
  const [creditCardMethods, setCreditCardMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editingCard, setEditingCard] = useState<Card | undefined>(undefined);

  const loadSummaries = useCallback(async () => {
    try {
      const res = await apiGet<{ data: CardSummary[] }>(
        `/api/v1/cards?month=${selectedMonth}&year=${year}`
      );
      setSummaries(res.data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load cards");
    }
  }, [selectedMonth, year]);

  const loadPaymentMethods = useCallback(async () => {
    try {
      const res = await apiGet<ListResponse<PaymentMethod>>("/api/v1/payment-methods");
      setCreditCardMethods(res.data.filter((pm) => pm.type === "credit_card"));
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to load payment methods");
    }
  }, []);

  useEffect(() => {
    const fetchAll = async () => {
      setLoading(true);
      await Promise.all([loadSummaries(), loadPaymentMethods()]);
      setLoading(false);
    };
    fetchAll();
  }, [loadSummaries, loadPaymentMethods]);

  const handleOpenCreate = () => {
    setEditingCard(undefined);
    setFormOpen(true);
  };

  const handleOpenEdit = (card: Card) => {
    setEditingCard(card);
    setFormOpen(true);
  };

  const handleClose = () => {
    setFormOpen(false);
    setEditingCard(undefined);
  };

  const handleSubmit = async (values: CardFormValues) => {
    const payload = {
      payment_method_id: values.payment_method_id,
      bank: values.bank,
      card_limit: values.card_limit,
      recommended_max_pct: values.recommended_max_pct,
      manual_usage_override: values.manual_usage_override || undefined,
      level: values.level || undefined,
    };

    try {
      if (editingCard) {
        await apiPut<Card>(`/api/v1/cards/${editingCard.id}`, payload);
        toast.success("Card updated");
      } else {
        await apiPost<Card>("/api/v1/cards", payload);
        toast.success("Card created");
      }
      handleClose();
      await loadSummaries();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save card");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/cards/${id}`);
      toast.success("Card deleted");
      await loadSummaries();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete card");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("title")}</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="size-4 mr-2" />
          {t("addCard")}
        </Button>
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm font-medium">{MONTH_NAMES[selectedMonth - 1]}</span>
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
      ) : summaries.length === 0 ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          {t("noCards")}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {summaries.map((summary) => (
            <CardHealthCard
              key={summary.card.id}
              summary={summary}
              onEdit={handleOpenEdit}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      <CardForm
        open={formOpen}
        onClose={handleClose}
        onSubmit={handleSubmit}
        creditCardPaymentMethods={creditCardMethods}
        defaultValues={editingCard}
      />
    </div>
  );
}
