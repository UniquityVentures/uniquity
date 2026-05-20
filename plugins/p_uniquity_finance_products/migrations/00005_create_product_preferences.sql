-- +goose Up
CREATE TABLE IF NOT EXISTS product_preferences (
    id                         BIGSERIAL PRIMARY KEY,
    created_at                 TIMESTAMPTZ,
    updated_at                 TIMESTAMPTZ,
    deleted_at                 TIMESTAMPTZ,
    inventory_account_id       BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    cost_of_sales_account_id   BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_product_preferences_deleted_at ON product_preferences (deleted_at);
CREATE INDEX IF NOT EXISTS idx_product_preferences_inventory_account_id ON product_preferences (inventory_account_id);
CREATE INDEX IF NOT EXISTS idx_product_preferences_cost_of_sales_account_id ON product_preferences (cost_of_sales_account_id);

-- +goose Down
DROP TABLE IF EXISTS product_preferences;
