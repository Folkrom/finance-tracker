"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { Trash2 } from "lucide-react";
import { PaymentMethod } from "@/types";
import { apiPost, apiDelete } from "@/lib/api";
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

interface PaymentMethodManagerProps {
  paymentMethods: PaymentMethod[];
  onRefresh: () => void;
}

type PMType = PaymentMethod["type"];

const PM_TYPES: PMType[] = [
  "cash",
  "debit_card",
  "credit_card",
  "digital_wallet",
  "crypto",
];

const PM_TYPE_LABELS: Record<PMType, string> = {
  cash: "Cash",
  debit_card: "Debit Card",
  credit_card: "Credit Card",
  digital_wallet: "Digital Wallet",
  crypto: "Crypto",
};

export function PaymentMethodManager({
  paymentMethods,
  onRefresh,
}: PaymentMethodManagerProps) {
  const t = useTranslations("settings");
  const tCommon = useTranslations("common");

  const [newName, setNewName] = useState("");
  const [newType, setNewType] = useState<PMType>("cash");
  const [newDetails, setNewDetails] = useState("");
  const [creating, setCreating] = useState(false);

  const handleCreate = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    try {
      await apiPost<PaymentMethod>("/api/v1/payment-methods", {
        name: newName.trim(),
        type: newType,
        details: newDetails.trim() || undefined,
      });
      setNewName("");
      setNewDetails("");
      toast.success("Payment method created");
      onRefresh();
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to create payment method"
      );
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiDelete(`/api/v1/payment-methods/${id}`);
      toast.success("Payment method deleted");
      onRefresh();
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to delete payment method"
      );
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold mb-4">{t("paymentMethods")}</h2>

        {/* Add form */}
        <div className="flex gap-2 items-end flex-wrap">
          <div className="space-y-1">
            <Label htmlFor="pm-name">Name</Label>
            <Input
              id="pm-name"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder="e.g. HSBC Debit"
              className="w-44"
            />
          </div>

          <div className="space-y-1">
            <Label>Type</Label>
            <Select
              value={newType}
              onValueChange={(val) => {
                if (val !== null) setNewType(val as PMType);
              }}
            >
              <SelectTrigger className="w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {PM_TYPES.map((type) => (
                  <SelectItem key={type} value={type}>
                    {PM_TYPE_LABELS[type]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1">
            <Label htmlFor="pm-details">Details (optional)</Label>
            <Input
              id="pm-details"
              value={newDetails}
              onChange={(e) => setNewDetails(e.target.value)}
              placeholder="Last 4 digits, etc."
              className="w-44"
            />
          </div>

          <Button onClick={handleCreate} disabled={creating || !newName.trim()}>
            {tCommon("create")}
          </Button>
        </div>
      </div>

      {/* List */}
      <div className="space-y-2">
        {paymentMethods.length === 0 ? (
          <p className="text-sm text-muted-foreground">{tCommon("noResults")}</p>
        ) : (
          paymentMethods.map((pm) => (
            <div
              key={pm.id}
              className="flex items-center justify-between rounded-lg border bg-white px-4 py-3"
            >
              <div className="flex items-center gap-3">
                <span className="font-medium text-sm">{pm.name}</span>
                <Badge variant="outline">{PM_TYPE_LABELS[pm.type]}</Badge>
                {pm.details && (
                  <span className="text-sm text-muted-foreground">{pm.details}</span>
                )}
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleDelete(pm.id)}
                aria-label={tCommon("delete")}
                className="text-destructive hover:text-destructive"
              >
                <Trash2 className="size-4" />
              </Button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
