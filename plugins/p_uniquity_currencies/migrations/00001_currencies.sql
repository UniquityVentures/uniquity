-- +goose Up
CREATE TABLE IF NOT EXISTS currencies (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    code       VARCHAR(3) NOT NULL,
    name       VARCHAR(50) NOT NULL DEFAULT '',
    symbol     VARCHAR(10) NOT NULL DEFAULT '',
    is_active  BOOLEAN NOT NULL DEFAULT true
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_currencies_code ON currencies (code);
CREATE INDEX IF NOT EXISTS idx_currencies_deleted_at ON currencies (deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_currencies_deleted_at;
DROP INDEX IF EXISTS idx_currencies_code;
DROP TABLE IF EXISTS currencies;
