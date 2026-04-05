"use client";

import { useTranslations } from "next-intl";
import { Pencil, Trash2 } from "lucide-react";
import { Debt } from "@/types";
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

interface DebtTableProps {
  debts: Debt[];
  onEdit: (debt: Debt) => void;
  onDelete: (id: string) => void;
}

export function DebtTable({ debts, onEdit, onDelete }: DebtTableProps) {
  const t = useTranslations("debt");
  const tCommon = useTranslations("common");

  if (debts.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {t("noDebts")}
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{t("name")}</TableHead>
          <TableHead>{t("amount")}</TableHead>
          <TableHead>{t("category")}</TableHead>
          <TableHead>{t("paymentMethod")}</TableHead>
          <TableHead>{t("date")}</TableHead>
          <TableHead className="w-[100px]" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {debts.map((debt) => (
          <TableRow key={debt.id}>
            <TableCell className="font-medium">{debt.name}</TableCell>
            <TableCell>
              {debt.currency}{" "}
              {parseFloat(debt.amount).toLocaleString("en-US", {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
              })}
            </TableCell>
            <TableCell>
              {debt.category ? (
                <Badge variant="secondary">{debt.category.name}</Badge>
              ) : (
                <span className="text-muted-foreground text-sm">—</span>
              )}
            </TableCell>
            <TableCell>
              {debt.payment_method ? (
                <Badge variant="secondary">{debt.payment_method.name}</Badge>
              ) : (
                <span className="text-muted-foreground text-sm">—</span>
              )}
            </TableCell>
            <TableCell>{debt.date}</TableCell>
            <TableCell>
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onEdit(debt)}
                  aria-label={tCommon("edit")}
                >
                  <Pencil className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onDelete(debt.id)}
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
