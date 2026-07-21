-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    token_hash varchar NOT NULL,
    is_blocked bool NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS sessions;