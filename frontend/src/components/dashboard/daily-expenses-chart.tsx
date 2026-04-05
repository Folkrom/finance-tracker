"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { DailySumData } from "@/types";

interface Props {
  data: DailySumData[];
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

export function DailyExpensesChart({ data }: Props) {
  const chartData = data.map((d) => ({
    date: formatDate(d.date),
    total: Number(d.total),
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={chartData} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
        <XAxis dataKey="date" tick={{ fontSize: 11 }} interval="preserveStartEnd" />
        <YAxis tick={{ fontSize: 12 }} />
        <Tooltip />
        <Bar dataKey="total" fill="#f97316" name="Amount" />
      </BarChart>
    </ResponsiveContainer>
  );
}
