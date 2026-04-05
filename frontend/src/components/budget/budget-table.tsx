"use client";

import { useTranslations } from "next-intl";
import { Pencil, Trash2 } from "lucide-react";
import { BudgetLine, Budget } from "@/types";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface BudgetTableProps {
  lines: BudgetLine[];
  onEdit: (budget: Budget) => void;
  onDelete: (id: string) => void;
}

export function BudgetTable({ lines, onEdit, onDelete }: BudgetTableProps) {
  const t = useTranslations("budget");
  const tCommon = useTranslations("common");

  if (lines.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {t("noBudgets")}
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{t("category")}</TableHead>
          <TableHead>{t("limit")}</TableHead>
          <TableHead>{t("spent")}</TableHead>
          <TableHead>{t("remaining")}</TableHead>
          <TableHead className="w-[100px]" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {lines.map((line) => {
          const limit = parseFloat(line.budget.monthly_limit);
          const spent = parseFloat(line.spent);
          const remaining = parseFloat(line.remaining);
          const pct = limit > 0 ? Math.min(100, Math.round((spent / limit) * 100)) : 0;

          return (
            <TableRow key={line.budget.id}>
              <TableCell className="font-medium">
                {line.budget.category?.name ?? line.budget.category_id}
              </TableCell>
              <TableCell>
                {limit.toLocaleString("en-US", {
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                })}
              </TableCell>
              <TableCell>
                <div className="space-y-1">
                  <span>
                    {spent.toLocaleString("en-US", {
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2,
                    })}
                  </span>
                  <div className="h-1.5 w-full rounded-full bg-muted overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${pct >= 100 ? "bg-red-500" : "bg-primary"}`}
                      style={{ width: `${pct}%` }}
                    />
                  </div>
                </div>
              </TableCell>
              <TableCell
                className={remaining >= 0 ? "text-green-600" : "text-red-600"}
              >
                {remaining.toLocaleString("en-US", {
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                })}
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onEdit(line.budget)}
                    aria-label={tCommon("edit")}
                  >
                    <Pencil className="size-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onDelete(line.budget.id)}
                    aria-label={tCommon("delete")}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
