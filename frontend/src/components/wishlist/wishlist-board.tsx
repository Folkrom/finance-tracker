"use client";

import { useTranslations } from "next-intl";
import { WishlistItem, WishlistStatus, WISHLIST_STATUS_GROUPS } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface WishlistBoardProps {
  items: WishlistItem[];
  onEdit: (item: WishlistItem) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: WishlistStatus) => void;
}

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

const priorityOrder: Record<string, number> = { high: 0, medium: 1, low: 2 };

const priorityBadgeClass: Record<string, string> = {
  high: "bg-red-100 text-red-800",
  medium: "bg-yellow-100 text-yellow-800",
  low: "bg-green-100 text-green-800",
};

const ALL_STATUSES: WishlistStatus[] = [
  "interested",
  "saving_for",
  "waiting_for_sale",
  "ordered",
  "purchased",
  "received",
  "cancelled",
];

export function WishlistBoard({
  items,
  onEdit,
  onDelete,
  onStatusChange,
}: WishlistBoardProps) {
  const t = useTranslations("wishlist");

  const columns = [
    { key: "todo", title: t("groupTodo"), statuses: WISHLIST_STATUS_GROUPS.todo },
    { key: "in_progress", title: t("groupInProgress"), statuses: WISHLIST_STATUS_GROUPS.in_progress },
    { key: "complete", title: t("groupComplete"), statuses: WISHLIST_STATUS_GROUPS.complete },
  ];

  return (
    <div className="flex overflow-x-auto gap-4">
      {columns.map((column) => {
        const columnItems = items
          .filter((item) => column.statuses.includes(item.status))
          .sort(
            (a, b) =>
              (priorityOrder[a.priority] ?? 1) -
              (priorityOrder[b.priority] ?? 1)
          );

        return (
          <div key={column.key} className="min-w-[300px] flex-1 bg-gray-50 rounded-lg p-3">
            <h3 className="font-semibold text-sm text-gray-500 uppercase mb-3">
              {column.title} ({columnItems.length})
            </h3>
            <div className="space-y-3">
              {columnItems.map((item) => (
                <Card key={item.id} size="sm" className="cursor-pointer" onClick={() => onEdit(item)}>
                  <CardContent className="flex flex-col gap-2">
                    <span className="font-medium">{item.name}</span>

                    {item.price && (
                      <span className="text-sm text-gray-600">
                        {item.price} {item.currency}
                      </span>
                    )}

                    <div className="flex flex-wrap gap-1">
                      <Badge className={priorityBadgeClass[item.priority]}>
                        {t(priorityKey[item.priority])}
                      </Badge>
                      {item.category && (
                        <Badge variant="outline">
                          {item.category.name}
                        </Badge>
                      )}
                    </div>

                    <div onClick={(e) => e.stopPropagation()}>
                      <Select
                        value={item.status}
                        onValueChange={(val) => {
                          if (val !== null) {
                            onStatusChange(item.id, val as WishlistStatus);
                          }
                        }}
                      >
                        <SelectTrigger size="sm" className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {ALL_STATUSES.map((status) => (
                            <SelectItem key={status} value={status}>
                              {t(statusKey[status])}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}
