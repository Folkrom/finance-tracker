"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { AdminStats } from "@/types";
import { StatsCards } from "@/components/admin/stats-cards";

export default function AdminStatsPage() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const data = await apiGet<AdminStats>("/api/v1/admin/stats");
        setStats(data);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to load stats");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        Loading...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Platform Stats</h1>
      {stats && <StatsCards stats={stats} />}
    </div>
  );
}
