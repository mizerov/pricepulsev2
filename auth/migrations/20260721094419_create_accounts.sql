-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id uuid PRIMARY KEY,
    email varchar UNIQUE NOT NULL,
    password_hash varchar NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS accounts;