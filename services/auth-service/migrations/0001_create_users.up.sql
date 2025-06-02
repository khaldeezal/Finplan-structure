-- Расширение pgcrypto для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Таблицу юзеров
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
    );