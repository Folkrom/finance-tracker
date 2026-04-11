"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { WishlistItem, Category } from "@/types";
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

const wishlistSchema = z.object({
  name: z.string().min(1, "Required"),
  image_url: z.string().optional(),
  price: z.string().optional(),
  currency: z.string().min(1),
  links: z.array(z.string()).max(5),
  category_id: z.string().optional(),
  priority: z.string().min(1),
  status: z.string().min(1),
  target_date: z.string().optional(),
  monthly_contribution: z.string().optional(),
  contribution_currency: z.string().min(1),
});

type WishlistFormValues = z.infer<typeof wishlistSchema>;

interface WishlistFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: WishlistFormValues) => Promise<void>;
  categories: Category[];
  defaultValues?: WishlistItem;
}

export function WishlistForm({
  open,
  onClose,
  onSubmit,
  categories,
  defaultValues,
}: WishlistFormProps) {
  const t = useTranslations("wishlist");
  const tCommon = useTranslations("common");

  const {
    register,
    handleSubmit,
    control,
    reset,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<WishlistFormValues>({
    resolver: zodResolver(wishlistSchema),
    defaultValues: {
      name: "",
      image_url: "",
      price: "",
      currency: "MXN",
      links: [],
      category_id: undefined,
      priority: "medium",
      status: "interested",
      target_date: "",
      monthly_contribution: "",
      contribution_currency: "MXN",
    },
  });

  useEffect(() => {
    if (open) {
      if (defaultValues) {
        reset({
          name: defaultValues.name,
          image_url: defaultValues.image_url ?? "",
          price: defaultValues.price ?? "",
          currency: defaultValues.currency,
          links: defaultValues.links ?? [],
          category_id: defaultValues.category_id ?? undefined,
          priority: defaultValues.priority,
          status: defaultValues.status,
          target_date: defaultValues.target_date ?? "",
          monthly_contribution: defaultValues.monthly_contribution ?? "",
          contribution_currency: defaultValues.contribution_currency,
        });
      } else {
        reset({
          name: "",
          image_url: "",
          price: "",
          currency: "MXN",
          links: [],
          category_id: undefined,
          priority: "medium",
          status: "interested",
          target_date: "",
          monthly_contribution: "",
          contribution_currency: "MXN",
        });
      }
    }
  }, [open, defaultValues, reset]);

  const handleFormSubmit = async (values: WishlistFormValues) => {
    await onSubmit(values);
  };

  const links = watch("links");

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) onClose(); }}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>
            {defaultValues ? t("editItem") : t("addItem")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          {/* name */}
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

          {/* image_url */}
          <div className="space-y-2">
            <Label htmlFor="image_url">{t("imageUrl")}</Label>
            <Input
              id="image_url"
              type="url"
              {...register("image_url")}
              placeholder="https://"
            />
          </div>

          {/* price + currency */}
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2">
              <Label htmlFor="price">{t("price")}</Label>
              <Input
                id="price"
                type="number"
                step="0.01"
                min="0"
                {...register("price")}
                placeholder="0.00"
              />
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

          {/* category_id */}
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
                    {categories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          {/* priority + status */}
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2">
              <Label>{t("priority")}</Label>
              <Controller
                name="priority"
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
                      <SelectItem value="low">{t("priorityLow")}</SelectItem>
                      <SelectItem value="medium">{t("priorityMedium")}</SelectItem>
                      <SelectItem value="high">{t("priorityHigh")}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>

            <div className="space-y-2">
              <Label>{t("status")}</Label>
              <Controller
                name="status"
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
                      <SelectItem value="interested">{t("statusInterested")}</SelectItem>
                      <SelectItem value="saving_for">{t("statusSavingFor")}</SelectItem>
                      <SelectItem value="waiting_for_sale">{t("statusWaitingForSale")}</SelectItem>
                      <SelectItem value="ordered">{t("statusOrdered")}</SelectItem>
                      <SelectItem value="purchased">{t("statusPurchased")}</SelectItem>
                      <SelectItem value="received">{t("statusReceived")}</SelectItem>
                      <SelectItem value="cancelled">{t("statusCancelled")}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          </div>

          {/* target_date */}
          <div className="space-y-2">
            <Label htmlFor="target_date">{t("targetDate")}</Label>
            <Input id="target_date" type="date" {...register("target_date")} />
          </div>

          {/* monthly_contribution + contribution_currency */}
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2">
              <Label htmlFor="monthly_contribution">{t("monthlyContribution")}</Label>
              <Input
                id="monthly_contribution"
                type="number"
                step="0.01"
                min="0"
                {...register("monthly_contribution")}
                placeholder="0.00"
              />
            </div>

            <div className="space-y-2">
              <Label>{tCommon("currency")}</Label>
              <Controller
                name="contribution_currency"
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

          {/* links */}
          <div className="space-y-2">
            <Label>{t("links")}</Label>
            <div className="space-y-2">
              {(links || []).map((link, i) => (
                <div key={i} className="flex gap-2">
                  <Input
                    type="url"
                    value={link}
                    placeholder="https://"
                    onChange={(e) => {
                      const updated = [...(links || [])];
                      updated[i] = e.target.value;
                      setValue("links", updated);
                    }}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      setValue("links", (links || []).filter((_, j) => j !== i));
                    }}
                  >
                    ×
                  </Button>
                </div>
              ))}
            </div>
            <Button
              type="button"
              variant="outline"
              disabled={(links || []).length >= 5}
              onClick={() => {
                setValue("links", [...(links || []), ""]);
              }}
            >
              {t("addLink")}
            </Button>
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
