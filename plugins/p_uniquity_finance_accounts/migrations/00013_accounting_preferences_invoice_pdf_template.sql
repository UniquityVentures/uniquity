-- +goose Up
ALTER TABLE accounting_preferences
    ADD COLUMN IF NOT EXISTS invoice_pdf_template TEXT;

-- +goose Down
ALTER TABLE accounting_preferences
    DROP COLUMN IF EXISTS invoice_pdf_template;
