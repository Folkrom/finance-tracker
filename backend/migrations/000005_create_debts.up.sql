CREATE TABLE debts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    date DATE NOT NULL,
    year INT NOT NULL,
    payment_method_id UUID REFERENCES payment_methods(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_debts_user_id ON debts(user_id);
CREATE INDEX idx_debts_year ON debts(user_id, year);
CREATE INDEX idx_debts_date ON debts(user_id, date);
CREATE INDEX idx_debts_category ON debts(user_id, category_id);
CREATE INDEX idx_debts_payment_method ON debts(user_id, payment_method_id);
