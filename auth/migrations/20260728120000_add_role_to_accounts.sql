-- +goose Up
ALTER TABLE accounts ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';

-- +goose Down
ALTER TABLE accounts DROP COLUMN role;
