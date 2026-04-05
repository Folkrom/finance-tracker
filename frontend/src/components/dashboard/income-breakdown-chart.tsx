"use client";

import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { CategorySumData } from "@/types";

const COLORS = [
  "#6366f1", "#22c55e", "#f59e0b", "#ef4444", "#3b82f6",
  "#ec4899", "#14b8a6", "#f97316", "#8b5cf6", "#84cc16",
];

interface Props {
  data: CategorySumData[];
}

export function IncomeBreakdownChart({ data }: Props) {
  const chartData = data.map((d) => ({
    name: d.category_name,
    value: Number(d.total),
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
          {chartData.map((_, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip formatter={(value) => typeof value === "number" ? value.toLocaleString() : value} />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
}
