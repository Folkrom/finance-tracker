"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Expense, Category, PaymentMethod } from "@/types";
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

const expenseSchema = z.object({
  name: z.string().min(1, "Name is required"),
  amount: z.string().min(1, "Amount is required"),
  currency: z.string().min(1, "Currency is required"),
  category_id: z.string().optional(),
  payment_method_id: z.string().optional(),
  type: z.enum(["expense", "saving", "investment"]),
  date: z.string().min(1, "Date is required"),
});

type ExpenseFormValues = z.infer<typeof expenseSchema>;

interface ExpenseFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: ExpenseFormValues) => Promise<void>;
  categories: Category[];
  paymentMethods: PaymentMethod[];
  defaultValues?: Expense;
}

export function ExpenseForm({
  open,
  onClose,
  onSubmit,
  categories,
  paymentMethods,
  defaultValues,
}: ExpenseFormProps) {
  const t = useTranslations("expenses");
  const tCommon = useTranslations("common");
  const params = useParams();
  const year = params.year as string;

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ExpenseFormValues>({
    resolver: zodResolver(expenseSchema),
    defaultValues: {
      name: "",
      amount: "",
      currency: "MXN",
      category_id: undefined,
      payment_method_id: undefined,
      type: "expense",
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
          type: defaultValues.type,
          date: defaultValues.date,
        });
      } else {
        reset({
          name: "",
          amount: "",
          currency: "MXN",
          category_id: undefined,
          payment_method_id: undefined,
          type: "expense",
          date: new Date().toISOString().split("T")[0],
        });
      }
    }
  }, [open, defaultValues, reset]);

  const expenseCategories = categories.filter((c) => c.domain === "expense");

  const typeItems = [
    { value: "expense", label: t("typeExpense") },
    { value: "saving", label: t("typeSaving") },
    { value: "investment", label: t("typeInvestment") },
  ];
  const categoryItems = [
    { value: "", label: "— None —" },
    ...expenseCategories.map((cat) => ({ value: cat.id, label: cat.name })),
  ];
  const paymentMethodItems = [
    { value: "", label: "— None —" },
    ...paymentMethods.map((pm) => ({ value: pm.id, label: pm.name })),
  ];

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editExpense") : t("addExpense")}
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
            <Label>{t("type")}</Label>
            <Controller
              name="type"
              control={control}
              render={({ field }) => (
                <Select
                  items={typeItems}
                  value={field.value}
                  onValueChange={(val) => {
                    if (val !== null) field.onChange(val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {typeItems.map((item) => (
                      <SelectItem key={item.value} value={item.value}>
                        {item.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-2">
            <Label>{t("category")}</Label>
            <Controller
              name="category_id"
              control={control}
              render={({ field }) => (
                <Select
                  items={categoryItems}
                  value={field.value ?? ""}
                  onValueChange={(val) => {
                    field.onChange(val === "" || val === null ? undefined : val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("category")} />
                  </SelectTrigger>
                  <SelectContent>
                    {categoryItems.map((item) => (
                      <SelectItem key={item.value} value={item.value}>
                        {item.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          <div className="space-y-2">
            <Label>{t("paymentMethod")}</Label>
            {paymentMethods.length === 0 ? (
              <p className="text-sm text-muted-foreground rounded-lg border border-dashed px-3 py-2">
                {tCommon("noPaymentMethods")}{" "}
                <Link
                  href={`/${year}/settings`}
                  className="underline underline-offset-2 hover:text-foreground"
                >
                  {tCommon("addInSettings")}
                </Link>
              </p>
            ) : (
              <Controller
                name="payment_method_id"
                control={control}
                render={({ field }) => (
                  <Select
                    items={paymentMethodItems}
                    value={field.value ?? ""}
                    onValueChange={(val) => {
                      field.onChange(val === "" || val === null ? undefined : val);
                    }}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("paymentMethod")} />
                    </SelectTrigger>
                    <SelectContent>
                      {paymentMethodItems.map((item) => (
                        <SelectItem key={item.value} value={item.value}>
                          {item.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
            )}
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
