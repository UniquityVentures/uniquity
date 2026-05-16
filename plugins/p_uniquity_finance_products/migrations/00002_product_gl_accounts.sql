-- +goose Up
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS inventory_account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    ADD COLUMN IF NOT EXISTS cost_of_sales_account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    ADD COLUMN IF NOT EXISTS input_tax_account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_products_inventory_account_id ON products (inventory_account_id);
CREATE INDEX IF NOT EXISTS idx_products_cost_of_sales_account_id ON products (cost_of_sales_account_id);
CREATE INDEX IF NOT EXISTS idx_products_input_tax_account_id ON products (input_tax_account_id);

-- +goose Down
ALTER TABLE products
    DROP COLUMN IF EXISTS input_tax_account_id,
    DROP COLUMN IF EXISTS cost_of_sales_account_id,
    DROP COLUMN IF EXISTS inventory_account_id;
