"use client";

import { useTranslations } from "next-intl";
import { Pencil, Trash2 } from "lucide-react";
import { Income } from "@/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface IncomeTableProps {
  incomes: Income[];
  onEdit: (income: Income) => void;
  onDelete: (id: string) => void;
}

export function IncomeTable({ incomes, onEdit, onDelete }: IncomeTableProps) {
  const t = useTranslations("income");
  const tCommon = useTranslations("common");

  if (incomes.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {t("noIncome")}
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{t("source")}</TableHead>
          <TableHead>{t("amount")}</TableHead>
          <TableHead>{t("category")}</TableHead>
          <TableHead>{t("date")}</TableHead>
          <TableHead className="w-[100px]" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {incomes.map((income) => (
          <TableRow key={income.id}>
            <TableCell className="font-medium">{income.source}</TableCell>
            <TableCell>
              {income.currency} {parseFloat(income.amount).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </TableCell>
            <TableCell>
              {income.category ? (
                <Badge variant="secondary">{income.category.name}</Badge>
              ) : (
                <span className="text-muted-foreground text-sm">—</span>
              )}
            </TableCell>
            <TableCell>{income.date}</TableCell>
            <TableCell>
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onEdit(income)}
                  aria-label={tCommon("edit")}
                >
                  <Pencil className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onDelete(income.id)}
                  aria-label={tCommon("delete")}
                  className="text-destructive hover:text-destructive"
                >
                  <Trash2 className="size-4" />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
