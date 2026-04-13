export interface Category {
  id: string;
  user_id: string | null;
  name: string;
  domain: "income" | "expense" | "wishlist";
  color?: string;
  sort_order: number;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

export interface PaymentMethod {
  id: string;
  user_id: string;
  name: string;
  type: "cash" | "debit_card" | "credit_card" | "digital_wallet" | "crypto";
  details?: string;
  created_at: string;
  updated_at: string;
}

export interface Income {
  id: string;
  user_id: string;
  source: string;
  amount: string;
  currency: string;
  category_id?: string;
  category?: Category;
  date: string;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface Expense {
  id: string;
  user_id: string;
  name: string;
  amount: string;
  currency: string;
  date: string;
  year: number;
  payment_method_id?: string;
  payment_method?: PaymentMethod;
  category_id?: string;
  category?: Category;
  type: "expense" | "saving" | "investment";
  created_at: string;
  updated_at: string;
}

export interface Debt {
  id: string;
  user_id: string;
  name: string;
  amount: string;
  currency: string;
  date: string;
  year: number;
  payment_method_id?: string;
  payment_method?: PaymentMethod;
  category_id?: string;
  category?: Category;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  user_id: string;
  category_id: string;
  category?: Category;
  monthly_limit: string;
  month: number;
  year: number;
  is_recurring: boolean;
  created_at: string;
  updated_at: string;
}

export interface BudgetLine {
  budget: Budget;
  spent: string;
  remaining: string;
}

export interface DashboardData {
  net_savings: MonthlyNet[];
  income_breakdown: CategorySumData[];
  expense_breakdown: TypeSumData[];
  income_vs_expenses: IncomeVsExpenses;
  daily_expenses: DailySumData[];
  daily_debts: DailySumData[];
}

export interface MonthlyNet {
  month: number;
  income: string;
  expenses: string;
  net: string;
}

export interface CategorySumData {
  category_id: string;
  category_name: string;
  total: string;
}

export interface TypeSumData {
  type: string;
  total: string;
}

export interface IncomeVsExpenses {
  total_income: string;
  total_expenses: string;
}

export interface DailySumData {
  date: string;
  total: string;
}

export interface Card {
  id: string;
  user_id: string;
  payment_method_id: string;
  payment_method?: PaymentMethod;
  bank: string;
  card_limit: string;
  recommended_max_pct: string;
  manual_usage_override?: string;
  level?: string;
  created_at: string;
  updated_at: string;
}

export interface CardSummary {
  card: Card;
  auto_usage: string;
  manual_override?: string;
  total_usage: string;
  usage_percent: number;
  recommended_max: string;
  health_color: "green" | "yellow" | "orange" | "red";
}

export type WishlistPriority = "low" | "medium" | "high";
export type WishlistStatus =
  | "interested"
  | "saving_for"
  | "waiting_for_sale"
  | "ordered"
  | "purchased"
  | "received"
  | "cancelled";

export interface WishlistItem {
  id: string;
  user_id: string;
  name: string;
  image_url?: string;
  price?: string;
  currency: string;
  links: string[];
  category_id?: string;
  category?: Category;
  priority: WishlistPriority;
  status: WishlistStatus;
  target_date?: string;
  monthly_contribution?: string;
  contribution_currency: string;
  created_at: string;
  updated_at: string;
}

export const WISHLIST_STATUS_GROUPS = {
  todo: ["interested"] as WishlistStatus[],
  in_progress: ["saving_for", "waiting_for_sale", "ordered"] as WishlistStatus[],
  complete: ["purchased", "received", "cancelled"] as WishlistStatus[],
};

export interface ListResponse<T> {
  data: T[];
  total: number;
}

export interface ErrorResponse {
  error: string;
}

export interface Profile {
  id: string;
  user_id: string;
  currency: string;
  language: string;
  created_at: string;
  updated_at: string;
}

export type { AdminStats } from "./admin";
