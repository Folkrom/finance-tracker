CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    payment_method_id UUID NOT NULL REFERENCES payment_methods(id) ON DELETE CASCADE,
    bank VARCHAR(255) NOT NULL,
    card_limit DECIMAL(12, 2) NOT NULL,
    recommended_max_pct DECIMAL(5, 2) NOT NULL DEFAULT 30.00,
    manual_usage_override DECIMAL(12, 2),
    level VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, payment_method_id)
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
