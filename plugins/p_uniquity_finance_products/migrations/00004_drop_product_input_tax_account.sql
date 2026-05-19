-- +goose Up
DROP INDEX IF EXISTS idx_products_input_tax_account_id;
ALTER TABLE products DROP COLUMN IF EXISTS input_tax_account_id;

-- +goose Down
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS input_tax_account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_products_input_tax_account_id ON products (input_tax_account_id);
