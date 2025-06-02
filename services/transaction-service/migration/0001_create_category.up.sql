

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL, -- 'INCOME' или 'EXPENSE'
    user_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для быстрого поиска по user_id
CREATE INDEX idx_categories_user_id ON categories(user_id);