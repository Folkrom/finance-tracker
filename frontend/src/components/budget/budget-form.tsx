"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Budget, Category } from "@/types";
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

const MONTHS = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];
const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

const budgetSchema = z.object({
  category_id: z.string().min(1, "Category is required"),
  monthly_limit: z.string().min(1, "Monthly limit is required"),
  is_recurring: z.boolean(),
  month: z.number().int().min(1).max(12),
  year: z.number().int(),
});

type BudgetFormValues = z.infer<typeof budgetSchema>;

interface BudgetFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: BudgetFormValues) => Promise<void>;
  categories: Category[];
  defaultValues?: Budget;
  currentMonth: number;
  currentYear: number;
}

export function BudgetForm({
  open,
  onClose,
  onSubmit,
  categories,
  defaultValues,
  currentMonth,
  currentYear,
}: BudgetFormProps) {
  const t = useTranslations("budget");
  const tCommon = useTranslations("common");

  const {
    register,
    handleSubmit,
    control,
    reset,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<BudgetFormValues>({
    resolver: zodResolver(budgetSchema),
    defaultValues: {
      category_id: "",
      monthly_limit: "",
      is_recurring: false,
      month: currentMonth,
      year: currentYear,
    },
  });

  const isRecurring = watch("is_recurring");

  useEffect(() => {
    if (open) {
      if (defaultValues) {
        reset({
          category_id: defaultValues.category_id,
          monthly_limit: defaultValues.monthly_limit,
          is_recurring: defaultValues.is_recurring,
          month: defaultValues.month,
          year: defaultValues.year,
        });
      } else {
        reset({
          category_id: "",
          monthly_limit: "",
          is_recurring: false,
          month: currentMonth,
          year: currentYear,
        });
      }
    }
  }, [open, defaultValues, reset, currentMonth, currentYear]);

  const expenseCategories = categories.filter((c) => c.domain === "expense");

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) onClose(); }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editBudget") : t("addBudget")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label>{t("category")}</Label>
            <Controller
              name="category_id"
              control={control}
              render={({ field }) => (
                <Select
                  value={field.value ?? ""}
                  onValueChange={(val) => {
                    if (val !== null) field.onChange(val);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("category")} />
                  </SelectTrigger>
                  <SelectContent>
                    {expenseCategories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {errors.category_id && (
              <p className="text-sm text-destructive">{errors.category_id.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="monthly_limit">{t("limit")}</Label>
            <Input
              id="monthly_limit"
              type="number"
              step="0.01"
              min="0"
              {...register("monthly_limit")}
              placeholder="0.00"
            />
            {errors.monthly_limit && (
              <p className="text-sm text-destructive">{errors.monthly_limit.message}</p>
            )}
          </div>

          <div className="flex items-center gap-3">
            <Controller
              name="is_recurring"
              control={control}
              render={({ field }) => (
                <input
                  id="is_recurring"
                  type="checkbox"
                  checked={field.value}
                  onChange={(e) => field.onChange(e.target.checked)}
                  className="size-4 rounded border-input accent-primary"
                />
              )}
            />
            <Label htmlFor="is_recurring" className="cursor-pointer">
              {t("isRecurring")}
            </Label>
          </div>

          {!isRecurring && (
            <div className="space-y-2">
              <Label>{t("month")}</Label>
              <Controller
                name="month"
                control={control}
                render={({ field }) => (
                  <Select
                    value={String(field.value)}
                    onValueChange={(val) => {
                      if (val !== null) field.onChange(Number(val));
                    }}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {MONTHS.map((m) => (
                        <SelectItem key={m} value={String(m)}>
                          {MONTH_NAMES[m - 1]}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          )}

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
