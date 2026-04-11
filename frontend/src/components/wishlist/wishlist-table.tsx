"use client";

import { useTranslations } from "next-intl";
import { Pencil, Trash2 } from "lucide-react";
import { WishlistItem, WISHLIST_STATUS_GROUPS } from "@/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface WishlistTableProps {
  items: WishlistItem[];
  onEdit: (item: WishlistItem) => void;
  onDelete: (id: string) => void;
}

const statusGroupOrder = (status: string): number => {
  if (WISHLIST_STATUS_GROUPS.todo.includes(status as any)) return 0;
  if (WISHLIST_STATUS_GROUPS.in_progress.includes(status as any)) return 1;
  return 2;
};

const priorityOrder: Record<string, number> = { high: 0, medium: 1, low: 2 };

const statusKey: Record<string, string> = {
  interested: "statusInterested",
  saving_for: "statusSavingFor",
  waiting_for_sale: "statusWaitingForSale",
  ordered: "statusOrdered",
  purchased: "statusPurchased",
  received: "statusReceived",
  cancelled: "statusCancelled",
};

const priorityKey: Record<string, string> = {
  low: "priorityLow",
  medium: "priorityMedium",
  high: "priorityHigh",
};

const priorityBadgeClass: Record<string, string> = {
  high: "bg-red-100 text-red-800",
  medium: "bg-yellow-100 text-yellow-800",
  low: "bg-green-100 text-green-800",
};

export function WishlistTable({ items, onEdit, onDelete }: WishlistTableProps) {
  const t = useTranslations("wishlist");
  const tCommon = useTranslations("common");

  if (items.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {t("noItems")}
      </div>
    );
  }

  const sorted = [...items].sort((a, b) => {
    const groupDiff = statusGroupOrder(a.status) - statusGroupOrder(b.status);
    if (groupDiff !== 0) return groupDiff;
    return (priorityOrder[a.priority] ?? 1) - (priorityOrder[b.priority] ?? 1);
  });

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{t("name")}</TableHead>
          <TableHead>{t("image")}</TableHead>
          <TableHead>{t("price")}</TableHead>
          <TableHead>{t("category")}</TableHead>
          <TableHead>{t("priority")}</TableHead>
          <TableHead>{t("status")}</TableHead>
          <TableHead>{t("targetDate")}</TableHead>
          <TableHead>{t("monthlyContribution")}</TableHead>
          <TableHead>{t("links")}</TableHead>
          <TableHead className="w-[100px]" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {sorted.map((item) => (
          <TableRow key={item.id}>
            <TableCell className="font-medium">{item.name}</TableCell>
            <TableCell>
              {item.image_url ? (
                <img
                  src={item.image_url}
                  alt={item.name}
                  className="w-10 h-10 rounded object-cover"
                />
              ) : (
                <span className="text-muted-foreground">—</span>
              )}
            </TableCell>
            <TableCell>
              {item.price ? (
                `$${item.price} ${item.currency}`
              ) : (
                <span className="text-muted-foreground">—</span>
              )}
            </TableCell>
            <TableCell>
              {item.category ? (
                <Badge variant="secondary">{item.category.name}</Badge>
              ) : (
                <span className="text-muted-foreground text-sm">—</span>
              )}
            </TableCell>
            <TableCell>
              <Badge
                className={priorityBadgeClass[item.priority] ?? ""}
                variant="outline"
              >
                {priorityKey[item.priority] ? t(priorityKey[item.priority]) : item.priority}
              </Badge>
            </TableCell>
            <TableCell>
              <Badge variant="secondary">
                {statusKey[item.status] ? t(statusKey[item.status]) : item.status}
              </Badge>
            </TableCell>
            <TableCell>
              {item.target_date ? (
                new Date(item.target_date).toLocaleDateString()
              ) : (
                <span className="text-muted-foreground">—</span>
              )}
            </TableCell>
            <TableCell>
              {item.monthly_contribution ? (
                `$${item.monthly_contribution} ${item.contribution_currency}`
              ) : (
                <span className="text-muted-foreground">—</span>
              )}
            </TableCell>
            <TableCell>{item.links?.length || 0}</TableCell>
            <TableCell>
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onEdit(item)}
                  aria-label={tCommon("edit")}
                >
                  <Pencil className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onDelete(item.id)}
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
