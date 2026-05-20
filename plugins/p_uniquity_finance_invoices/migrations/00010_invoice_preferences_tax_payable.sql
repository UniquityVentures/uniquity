-- +goose Up
ALTER TABLE invoice_preferences
    ADD COLUMN IF NOT EXISTS account_tax_payable_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_invoice_preferences_account_tax_payable_id ON invoice_preferences (account_tax_payable_id);

ALTER TABLE draft_invoices
    DROP COLUMN IF EXISTS account_tax_payable_id;

-- +goose Down
ALTER TABLE draft_invoices
    ADD COLUMN IF NOT EXISTS account_tax_payable_id BIGINT NOT NULL DEFAULT 0 REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE draft_invoices
    ALTER COLUMN account_tax_payable_id DROP DEFAULT;

ALTER TABLE invoice_preferences
    DROP COLUMN IF EXISTS account_tax_payable_id;
