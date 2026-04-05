"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Card, PaymentMethod } from "@/types";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const cardSchema = z.object({
  payment_method_id: z.string().min(1, "Required"),
  bank: z.string().min(1, "Required"),
  card_limit: z.string().min(1, "Required"),
  recommended_max_pct: z.string().min(1, "Required"),
  manual_usage_override: z.string().optional(),
  level: z.string().optional(),
});

type CardFormValues = z.infer<typeof cardSchema>;

interface CardFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: CardFormValues) => Promise<void>;
  creditCardPaymentMethods: PaymentMethod[];
  defaultValues?: Card;
}

export function CardForm({
  open,
  onClose,
  onSubmit,
  creditCardPaymentMethods,
  defaultValues,
}: CardFormProps) {
  const t = useTranslations("cards");
  const tCommon = useTranslations("common");

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<CardFormValues>({
    resolver: zodResolver(cardSchema),
    defaultValues: {
      payment_method_id: "",
      bank: "",
      card_limit: "",
      recommended_max_pct: "30",
      manual_usage_override: "",
      level: "",
    },
  });

  useEffect(() => {
    if (open) {
      if (defaultValues) {
        reset({
          payment_method_id: defaultValues.payment_method_id,
          bank: defaultValues.bank,
          card_limit: defaultValues.card_limit,
          recommended_max_pct: defaultValues.recommended_max_pct,
          manual_usage_override: defaultValues.manual_usage_override ?? "",
          level: defaultValues.level ?? "",
        });
      } else {
        reset({
          payment_method_id: "",
          bank: "",
          card_limit: "",
          recommended_max_pct: "30",
          manual_usage_override: "",
          level: "",
        });
      }
    }
  }, [open, defaultValues, reset]);

  const handleFormSubmit = async (values: CardFormValues) => {
    await onSubmit(values);
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editCard") : t("addCard")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          {/* Payment Method */}
          <div className="space-y-2">
            <Label>{t("selectCreditCard")}</Label>
            <Controller
              name="payment_method_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value}
                  onValueChange={(val) => {
                    if (val !== null) field.onChange(val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("selectCreditCard")} />
                  </SelectTrigger>
                  <SelectContent>
                    {creditCardPaymentMethods.map((pm) => (
                      <SelectItem key={pm.id} value={pm.id}>
                        {pm.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {errors.payment_method_id && (
              <p className="text-sm text-destructive">{errors.payment_method_id.message}</p>
            )}
          </div>

          {/* Bank */}
          <div className="space-y-2">
            <Label htmlFor="bank">{t("bank")}</Label>
            <Input
              id="bank"
              {...register("bank")}
              placeholder={t("bank")}
            />
            {errors.bank && (
              <p className="text-sm text-destructive">{errors.bank.message}</p>
            )}
          </div>

          {/* Card Limit */}
          <div className="space-y-2">
            <Label htmlFor="card_limit">{t("cardLimit")}</Label>
            <Input
              id="card_limit"
              type="number"
              step="0.01"
              min="0"
              {...register("card_limit")}
              placeholder="0.00"
            />
            {errors.card_limit && (
              <p className="text-sm text-destructive">{errors.card_limit.message}</p>
            )}
          </div>

          {/* Recommended Max % */}
          <div className="space-y-2">
            <Label htmlFor="recommended_max_pct">{t("recommendedMaxPct")}</Label>
            <Input
              id="recommended_max_pct"
              type="number"
              step="1"
              min="0"
              max="100"
              {...register("recommended_max_pct")}
              placeholder="30"
            />
            {errors.recommended_max_pct && (
              <p className="text-sm text-destructive">{errors.recommended_max_pct.message}</p>
            )}
          </div>

          {/* Manual Usage Override */}
          <div className="space-y-2">
            <Label htmlFor="manual_usage_override">{t("manualOverride")}</Label>
            <Input
              id="manual_usage_override"
              type="number"
              step="0.01"
              min="0"
              {...register("manual_usage_override")}
              placeholder="0.00 (optional)"
            />
            {errors.manual_usage_override && (
              <p className="text-sm text-destructive">{errors.manual_usage_override.message}</p>
            )}
          </div>

          {/* Level */}
          <div className="space-y-2">
            <Label htmlFor="level">{t("level")}</Label>
            <Input
              id="level"
              {...register("level")}
              placeholder="e.g. Gold, Platinum (optional)"
            />
            {errors.level && (
              <p className="text-sm text-destructive">{errors.level.message}</p>
            )}
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>
              {tCommon("cancel")}
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? tCommon("loading") : tCommon("save")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
