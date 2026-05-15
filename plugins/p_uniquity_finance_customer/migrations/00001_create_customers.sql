-- +goose Up
CREATE TABLE IF NOT EXISTS customers (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    name        TEXT NOT NULL,
    address     TEXT,
    gstin       TEXT,
    pan         TEXT,
    phone       TEXT,
    email       TEXT,
    website     TEXT
);

CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS customers;
