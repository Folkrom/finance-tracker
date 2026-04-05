export interface Category {
  id: string;
  user_id: string;
  name: string;
  domain: "income" | "expense" | "wishlist";
  color?: string;
  sort_order: number;
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

export interface ListResponse<T> {
  data: T[];
  total: number;
}

export interface ErrorResponse {
  error: string;
}
