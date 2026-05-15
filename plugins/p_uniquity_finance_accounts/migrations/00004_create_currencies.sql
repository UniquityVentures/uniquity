-- +goose Up
CREATE TABLE IF NOT EXISTS currencies (
    id           BIGSERIAL PRIMARY KEY,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    code         INTEGER     NOT NULL UNIQUE,
    name         TEXT        NOT NULL,
    symbol       TEXT        NOT NULL DEFAULT '',
    minor_unit   INTEGER     NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_currencies_deleted_at ON currencies (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS currencies;
