-- +goose Up
-- Содержимое pg_dump (схема базы данных)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- Команда для отката миграции (например, удаление таблиц)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;