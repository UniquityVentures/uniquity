-- +goose Up
CREATE TABLE IF NOT EXISTS accounting_preferences (
    id                     BIGSERIAL PRIMARY KEY,
    created_at             TIMESTAMPTZ,
    updated_at             TIMESTAMPTZ,
    deleted_at             TIMESTAMPTZ,
    invoice_number_format  TEXT
);

CREATE INDEX IF NOT EXISTS idx_accounting_preferences_deleted_at ON accounting_preferences (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS accounting_preferences;
