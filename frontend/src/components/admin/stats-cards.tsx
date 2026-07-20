"use client";

import { Users, Layers, UserPlus, Globe } from "lucide-react";
import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AdminStats } from "@/types";

interface StatsCardsProps {
  stats: AdminStats;
}

const statConfig = [
  { key: "users" as const, labelKey: "users" as const, icon: Users },
  { key: "profiles" as const, labelKey: "profiles" as const, icon: UserPlus },
  { key: "categories_global" as const, labelKey: "globalCategories" as const, icon: Globe },
  { key: "categories_user" as const, labelKey: "userCategories" as const, icon: Layers },
];

export function StatsCards({ stats }: StatsCardsProps) {
  const t = useTranslations("admin");
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {statConfig.map((item) => (
        <Card key={item.key}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{t(item.labelKey)}</CardTitle>
            <item.icon className="size-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats[item.key]}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
