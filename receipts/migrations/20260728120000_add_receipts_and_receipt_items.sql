-- +goose Up
CREATE TABLE IF EXISTS receipts (
    id uuid PRIMARY KEY,
    total_sum bigint NOT NULL,
    date_time timestamptz NOT NULL,
    store_name varchar NOT NULL,
    store_inn varchar NOT NULL,
    fn varchar NOT NULL,
    fd varchar NOT NULL,
    fp varchar NOT NULL,
    created_by varchar NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE IF EXISTS receipt_items (
    id uuid PRIMARY KEY,
    receipt_id uuid NOT NULL,
    name varchar NOT NULL,
    amount bigint NOT NULL,
    price bigint NOT NULL,
    quantity int NOT NULL,
    product_id uuid
);

-- +goose Down
DROP TABLE IF EXISTS receipt_items;
DROP TABLE IF EXISTS receipts;
