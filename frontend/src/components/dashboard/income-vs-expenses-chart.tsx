"use client";

import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { IncomeVsExpenses } from "@/types";

interface Props {
  data: IncomeVsExpenses;
}

export function IncomeVsExpensesChart({ data }: Props) {
  const chartData = [
    { name: "Income", value: Number(data.total_income), color: "#22c55e" },
    { name: "Expenses", value: Number(data.total_expenses), color: "#ef4444" },
  ];

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
