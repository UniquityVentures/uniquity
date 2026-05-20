-- +goose Up
CREATE TABLE IF NOT EXISTS invoice_preferences (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    account_receivable_id   BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_revenue_id      BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id              BIGINT REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_invoice_preferences_deleted_at ON invoice_preferences (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoice_preferences_account_receivable_id ON invoice_preferences (account_receivable_id);
CREATE INDEX IF NOT EXISTS idx_invoice_preferences_account_revenue_id ON invoice_preferences (account_revenue_id);
CREATE INDEX IF NOT EXISTS idx_invoice_preferences_journal_id ON invoice_preferences (journal_id);

-- +goose Down
DROP TABLE IF EXISTS invoice_preferences;
