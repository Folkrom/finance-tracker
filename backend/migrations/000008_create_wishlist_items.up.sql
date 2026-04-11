CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_url TEXT,
    price DECIMAL(12, 2),
    currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    links TEXT[] DEFAULT '{}',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    status VARCHAR(30) NOT NULL DEFAULT 'interested' CHECK (status IN (
        'interested',
        'saving_for', 'waiting_for_sale', 'ordered',
        'purchased', 'received', 'cancelled'
    )),
    target_date DATE,
    monthly_contribution DECIMAL(12, 2),
    contribution_currency VARCHAR(3) NOT NULL DEFAULT 'MXN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wishlist_user_id ON wishlist_items(user_id);
CREATE INDEX idx_wishlist_status ON wishlist_items(user_id, status);
CREATE INDEX idx_wishlist_category ON wishlist_items(user_id, category_id);
CREATE INDEX idx_wishlist_priority ON wishlist_items(user_id, priority);
