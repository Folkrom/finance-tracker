"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Debt, Category, PaymentMethod } from "@/types";
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

const debtSchema = z.object({
  name: z.string().min(1, "Name is required"),
  amount: z.string().min(1, "Amount is required"),
  currency: z.string().min(1, "Currency is required"),
  category_id: z.string().optional(),
  payment_method_id: z.string().optional(),
  date: z.string().min(1, "Date is required"),
});

type DebtFormValues = z.infer<typeof debtSchema>;

interface DebtFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: DebtFormValues) => Promise<void>;
  categories: Category[];
  paymentMethods: PaymentMethod[];
  defaultValues?: Debt;
}

export function DebtForm({
  open,
  onClose,
  onSubmit,
  categories,
  paymentMethods,
  defaultValues,
}: DebtFormProps) {
  const t = useTranslations("debt");
  const tCommon = useTranslations("common");

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<DebtFormValues>({
    resolver: zodResolver(debtSchema),
    defaultValues: {
      name: "",
      amount: "",
      currency: "MXN",
      category_id: undefined,
      payment_method_id: undefined,
      date: "",
    },
  });

  useEffect(() => {
    if (open) {
      if (defaultValues) {
        reset({
          name: defaultValues.name,
          amount: defaultValues.amount,
          currency: defaultValues.currency,
          category_id: defaultValues.category_id ?? undefined,
          payment_method_id: defaultValues.payment_method_id ?? undefined,
          date: defaultValues.date,
        });
      } else {
        reset({
          name: "",
          amount: "",
          currency: "MXN",
          category_id: undefined,
          payment_method_id: undefined,
          date: new Date().toISOString().split("T")[0],
        });
      }
    }
  }, [open, defaultValues, reset]);

  const expenseCategories = categories.filter((c) => c.domain === "expense");

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editDebt") : t("addDebt")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">{t("name")}</Label>
            <Input
              id="name"
              {...register("name")}
              placeholder={t("name")}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2">
              <Label htmlFor="amount">{t("amount")}</Label>
              <Input
                id="amount"
                type="number"
                step="0.01"
                min="0"
                {...register("amount")}
                placeholder="0.00"
              />
              {errors.amount && (
                <p className="text-sm text-destructive">{errors.amount.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label>{tCommon("currency")}</Label>
              <Controller
                name="currency"
                control={control}
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={(val) => {
                      if (val !== null) field.onChange(val);
                    }}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="MXN">MXN</SelectItem>
                      <SelectItem value="USD">USD</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label>{t("category")}</Label>
            <Controller
              name="category_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value ?? ""}
                  onValueChange={(val) => {
                    field.onChange(val === "" || val === null ? undefined : val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("category")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">— None —</SelectItem>
                    {expenseCategories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-2">
            <Label>{t("paymentMethod")}</Label>
            <Controller
              name="payment_method_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value ?? ""}
                  onValueChange={(val) => {
                    field.onChange(val === "" || val === null ? undefined : val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("paymentMethod")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">— None —</SelectItem>
                    {paymentMethods.map((pm) => (
                      <SelectItem key={pm.id} value={pm.id}>
                        {pm.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="date">{t("date")}</Label>
            <Input id="date" type="date" {...register("date")} />
            {errors.date && (
              <p className="text-sm text-destructive">{errors.date.message}</p>
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
