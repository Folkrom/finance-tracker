"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { MonthlyNet } from "@/types";

const MONTH_NAMES = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun",
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
];

interface Props {
  data: MonthlyNet[];
}

export function NetSavingsChart({ data }: Props) {
  const chartData = data.map((d) => ({
    month: MONTH_NAMES[d.month - 1] ?? String(d.month),
    income: Number(d.income),
    expenses: Number(d.expenses),
    net: Number(d.net),
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={chartData} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
        <XAxis dataKey="month" tick={{ fontSize: 12 }} />
        <YAxis tick={{ fontSize: 12 }} />
        <Tooltip />
        <Legend />
        <Line
          type="monotone"
          dataKey="income"
          stroke="#22c55e"
          strokeWidth={2}
          dot={false}
          name="Income"
        />
        <Line
          type="monotone"
          dataKey="expenses"
          stroke="#ef4444"
          strokeWidth={2}
          dot={false}
          name="Expenses"
        />
        <Line
          type="monotone"
          dataKey="net"
          stroke="#3b82f6"
          strokeWidth={2}
          strokeDasharray="5 5"
          dot={false}
          name="Net"
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
