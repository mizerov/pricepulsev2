-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY,
    name varchar NOT NULL,
    description text,
    status varchar NOT NULL, -- e.g., 'draft', 'approved'
    owner_id uuid, -- AccountID of the user who proposed it
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_proposals (
    id uuid PRIMARY KEY,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    account_id uuid NOT NULL,
    status varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS outbox_events (
    id uuid PRIMARY KEY,
    aggregate_type varchar NOT NULL,
    aggregate_id uuid NOT NULL,
    type varchar NOT NULL,
    payload jsonb NOT NULL,
    status varchar NOT NULL DEFAULT 'pending',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS product_proposals;
DROP TABLE IF EXISTS products;
