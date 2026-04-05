"use client";

import { useTranslations } from "next-intl";
import { Pencil, Trash2 } from "lucide-react";
import { Card, CardSummary } from "@/types";
import {
  Card as UICard,
  CardHeader,
  CardTitle,
  CardAction,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface CardHealthCardProps {
  summary: CardSummary;
  onEdit: (card: Card) => void;
  onDelete: (id: string) => void;
}

const healthColorBar: Record<string, string> = {
  green: "bg-green-500",
  yellow: "bg-yellow-500",
  orange: "bg-orange-500",
  red: "bg-red-500",
};

const healthBadgeStyle: Record<string, string> = {
  green: "bg-green-100 text-green-800 border-green-200",
  yellow: "bg-yellow-100 text-yellow-800 border-yellow-200",
  orange: "bg-orange-100 text-orange-800 border-orange-200",
  red: "bg-red-100 text-red-800 border-red-200",
};

export function CardHealthCard({ summary, onEdit, onDelete }: CardHealthCardProps) {
  const t = useTranslations("cards");

  const { card, auto_usage, manual_override, total_usage, usage_percent, health_color } = summary;

  const limit = parseFloat(card.card_limit);
  const usedPct = Math.min(100, usage_percent);

  const healthLabel: Record<string, string> = {
    green: t("healthy"),
    yellow: t("recommended"),
    orange: t("warning"),
    red: t("danger"),
  };

  const fmt = (val: string | undefined) =>
    val ? `$${parseFloat(val).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : "$0.00";

  return (
    <UICard>
      <CardHeader>
        <div className="flex flex-col gap-0.5">
          <div className="flex items-center gap-2">
            <CardTitle>{card.bank}</CardTitle>
            {card.level && (
              <Badge variant="outline" className="text-xs">
                {card.level}
              </Badge>
            )}
          </div>
          <p className="text-sm text-muted-foreground">
            {card.payment_method?.name ?? "—"}
          </p>
        </div>
        <CardAction>
          <div className="flex gap-1">
            <Button
              size="sm"
              variant="ghost"
              onClick={() => onEdit(card)}
              aria-label="Edit card"
            >
              <Pencil className="size-4" />
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => onDelete(card.id)}
              aria-label="Delete card"
              className="text-destructive hover:text-destructive"
            >
              <Trash2 className="size-4" />
            </Button>
          </div>
        </CardAction>
      </CardHeader>

      <CardContent className="space-y-3">
        {/* Credit limit */}
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">{t("cardLimit")}</span>
          <span className="font-medium">{fmt(card.card_limit)}</span>
        </div>

        {/* Usage progress bar */}
        <div className="space-y-1">
          <div className="h-2.5 w-full overflow-hidden rounded-full bg-muted">
            <div
              className={`h-full rounded-full transition-all ${healthColorBar[health_color] ?? "bg-green-500"}`}
              style={{ width: `${usedPct}%` }}
            />
          </div>
          <p className="text-xs text-right text-muted-foreground">
            {usedPct.toFixed(1)}% {t("usageThisMonth")}
          </p>
        </div>

        {/* Usage breakdown */}
        <div className="text-xs text-muted-foreground space-y-0.5">
          <div className="flex justify-between">
            <span>{t("autoUsage")}</span>
            <span>{fmt(auto_usage)}</span>
          </div>
          {manual_override && (
            <div className="flex justify-between">
              <span>{t("manualOverride")}</span>
              <span>{fmt(manual_override)}</span>
            </div>
          )}
          <div className="flex justify-between font-medium text-foreground border-t pt-1 mt-1">
            <span>{t("totalUsage")}</span>
            <span>{fmt(total_usage)}</span>
          </div>
        </div>

        {/* Recommended max */}
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">{t("recommendedMax")}</span>
          <span className="text-muted-foreground">
            {card.recommended_max_pct}% ({fmt(summary.recommended_max)})
          </span>
        </div>
      </CardContent>

      <CardFooter className="justify-between">
        <span className="text-xs text-muted-foreground">{t("healthIndicator")}</span>
        <span
          className={`inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium ${healthBadgeStyle[health_color] ?? ""}`}
        >
          {healthLabel[health_color] ?? health_color}
        </span>
      </CardFooter>
    </UICard>
  );
}
