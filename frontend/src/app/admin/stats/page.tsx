"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { apiGet } from "@/lib/api";
import { AdminStats } from "@/types";
import { StatsCards } from "@/components/admin/stats-cards";

export default function AdminStatsPage() {
  const t = useTranslations("admin");
  const tCommon = useTranslations("common");
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const data = await apiGet<AdminStats>("/api/v1/admin/stats");
        setStats(data);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : t("loadStatsFailed"));
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [t]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {tCommon("loading")}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t("statsTitle")}</h1>
      {stats && <StatsCards stats={stats} />}
    </div>
  );
}
