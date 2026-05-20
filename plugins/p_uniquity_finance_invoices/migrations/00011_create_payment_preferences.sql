-- +goose Up
CREATE TABLE IF NOT EXISTS payment_preferences (
    id                 BIGSERIAL PRIMARY KEY,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ,
    deleted_at         TIMESTAMPTZ,
    payment_account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_payment_preferences_deleted_at ON payment_preferences (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payment_preferences_payment_account_id ON payment_preferences (payment_account_id);

-- +goose Down
DROP TABLE IF EXISTS payment_preferences;
