"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { apiGet } from "@/lib/api";
import { DashboardData } from "@/types";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { NetSavingsChart } from "@/components/dashboard/net-savings-chart";
import { IncomeBreakdownChart } from "@/components/dashboard/income-breakdown-chart";
import { ExpenseBreakdownChart } from "@/components/dashboard/expense-breakdown-chart";
import { IncomeVsExpensesChart } from "@/components/dashboard/income-vs-expenses-chart";
import { DailyExpensesChart } from "@/components/dashboard/daily-expenses-chart";
import { DailyDebtsChart } from "@/components/dashboard/daily-debts-chart";

export default function DashboardPage() {
  const params = useParams();
  const year = params.year as string;
  const t = useTranslations("dashboard");
  const tCommon = useTranslations("common");

  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadData() {
      setLoading(true);
      try {
        const result = await apiGet<DashboardData>(
          `/api/v1/years/${year}/dashboard`
        );
        setData(result);
      } catch (err) {
        toast.error(
          err instanceof Error ? err.message : "Failed to load dashboard"
        );
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, [year]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {tCommon("loading")}
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        {t("noData")}
      </div>
    );
  }

  const hasNetSavings = data.net_savings.length > 0;
  const hasIncomeBreakdown = data.income_breakdown.length > 0;
  const hasExpenseBreakdown = data.expense_breakdown.length > 0;
  const hasDailyExpenses = data.daily_expenses.length > 0;
  const hasDailyDebts = data.daily_debts.length > 0;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t("title")}</h1>

      {/* Row 1: Net Savings (2/3) + Income vs Expenses (1/3) */}
      <div className="grid grid-cols-3 gap-6">
        <Card className="col-span-2">
          <CardHeader>
            <CardTitle>{t("netSavings")}</CardTitle>
          </CardHeader>
          <CardContent>
            {hasNetSavings ? (
              <NetSavingsChart data={data.net_savings} />
            ) : (
              <div className="flex items-center justify-center h-[300px] text-muted-foreground text-sm">
                {t("noData")}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="col-span-1">
          <CardHeader>
            <CardTitle>{t("incomeVsExpenses")}</CardTitle>
          </CardHeader>
          <CardContent>
            <IncomeVsExpensesChart data={data.income_vs_expenses} />
          </CardContent>
        </Card>
      </div>

      {/* Row 2: Income Breakdown + Expense Breakdown */}
      <div className="grid grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>{t("incomeBreakdown")}</CardTitle>
          </CardHeader>
          <CardContent>
            {hasIncomeBreakdown ? (
              <IncomeBreakdownChart data={data.income_breakdown} />
            ) : (
              <div className="flex items-center justify-center h-[300px] text-muted-foreground text-sm">
                {t("noData")}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("expenseBreakdown")}</CardTitle>
          </CardHeader>
          <CardContent>
            {hasExpenseBreakdown ? (
              <ExpenseBreakdownChart data={data.expense_breakdown} />
            ) : (
              <div className="flex items-center justify-center h-[300px] text-muted-foreground text-sm">
                {t("noData")}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Row 3: Daily Expenses (full width) */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dailyExpenses")}</CardTitle>
        </CardHeader>
        <CardContent>
          {hasDailyExpenses ? (
            <DailyExpensesChart data={data.daily_expenses} />
          ) : (
            <div className="flex items-center justify-center h-[300px] text-muted-foreground text-sm">
              {t("noData")}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Row 4: Daily Debts (full width) */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dailyDebts")}</CardTitle>
        </CardHeader>
        <CardContent>
          {hasDailyDebts ? (
            <DailyDebtsChart data={data.daily_debts} />
          ) : (
            <div className="flex items-center justify-center h-[300px] text-muted-foreground text-sm">
              {t("noData")}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
