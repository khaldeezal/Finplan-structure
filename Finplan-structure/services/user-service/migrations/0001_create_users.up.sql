CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Создаём таблицу пользователей, ориентируясь на user-профиль
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный идентификатор пользователя
    name TEXT NOT NULL,                            -- Имя пользователя
    email TEXT UNIQUE NOT NULL,                    -- Email (уникальный)
    currency TEXT NOT NULL,                        -- Валюта, выбранная пользователем
    language TEXT NOT NULL,                        -- Язык интерфейса пользователя
    created_at TIMESTAMP DEFAULT now(),            -- Дата создания записи
    updated_at TIMESTAMP DEFAULT now()             -- Дата последнего обновления
);