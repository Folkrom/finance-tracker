"use client";

import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { TypeSumData } from "@/types";

const TYPE_COLORS: Record<string, string> = {
  expense: "#ef4444",
  saving: "#22c55e",
  investment: "#3b82f6",
};

const TYPE_LABELS: Record<string, string> = {
  expense: "Expense",
  saving: "Saving",
  investment: "Investment",
};

interface Props {
  data: TypeSumData[];
}

export function ExpenseBreakdownChart({ data }: Props) {
  const chartData = data.map((d) => ({
    name: TYPE_LABELS[d.type] ?? d.type,
    value: Number(d.total),
    color: TYPE_COLORS[d.type] ?? "#94a3b8",
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={chartData}
          cx="50%"
          cy="45%"
          innerRadius={60}
          outerRadius={80}
          dataKey="value"
        >
          {chartData.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color} />
          ))}
        </Pie>
        <Tooltip formatter={(value) => typeof value === "number" ? value.toLocaleString() : value} />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
}
