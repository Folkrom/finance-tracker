"use client";

import { useTranslations } from "next-intl";
import { WishlistItem } from "@/types";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

interface WishlistGalleryProps {
  items: WishlistItem[];
  onEdit: (item: WishlistItem) => void;
  onDelete: (id: string) => void;
}

const statusTranslationKey: Record<string, string> = {
  interested: "statusInterested",
  saving_for: "statusSavingFor",
  waiting_for_sale: "statusWaitingForSale",
  ordered: "statusOrdered",
  purchased: "statusPurchased",
  received: "statusReceived",
  cancelled: "statusCancelled",
};

const priorityTranslationKey: Record<string, string> = {
  low: "priorityLow",
  medium: "priorityMedium",
  high: "priorityHigh",
};

const priorityClassName: Record<string, string> = {
  high: "bg-red-100 text-red-800",
  medium: "bg-yellow-100 text-yellow-800",
  low: "bg-green-100 text-green-800",
};

export function WishlistGallery({ items, onEdit, onDelete }: WishlistGalleryProps) {
  const t = useTranslations("wishlist");
  const tCommon = useTranslations("common");

  if (items.length === 0) {
    return (
      <div className="flex justify-center py-12 text-muted-foreground">
        {t("noItems")}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {items.map((item) => (
        <Card key={item.id} className="p-0">
          {item.image_url ? (
            <img
              src={item.image_url}
              alt={item.name}
              className="w-full h-48 object-cover rounded-t-lg"
            />
          ) : (
            <div className="h-48 bg-gray-100 rounded-t-lg flex items-center justify-center">
              <span className="text-gray-400 text-sm">{t("noImage")}</span>
            </div>
          )}

          <CardContent className="flex flex-col gap-2 pt-3">
            <p className="font-semibold truncate">{item.name}</p>

            {item.price && (
              <p className="text-sm text-muted-foreground">
                ${item.price} {item.currency}
              </p>
            )}

            <div className="flex flex-wrap gap-1">
              {item.category && (
                <Badge variant="outline">{item.category.name}</Badge>
              )}

              <Badge className={priorityClassName[item.priority]}>
                {t(priorityTranslationKey[item.priority])}
              </Badge>

              <Badge variant="secondary">
                {t(statusTranslationKey[item.status])}
              </Badge>
            </div>
          </CardContent>

          <CardFooter className="gap-2">
            <Button size="sm" variant="outline" onClick={() => onEdit(item)}>
              {tCommon("edit")}
            </Button>
            <Button size="sm" variant="destructive" onClick={() => onDelete(item.id)}>
              {tCommon("delete")}
            </Button>
          </CardFooter>
        </Card>
      ))}
    </div>
  );
}
