-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id           BIGSERIAL PRIMARY KEY,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    name         TEXT NOT NULL,
    base_cost    NUMERIC(19, 6) NOT NULL,
    sales_price  NUMERIC(19, 6) NOT NULL,
    hsn_code     BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products (deleted_at);

CREATE TABLE IF NOT EXISTS product_taxes (
    product_id BIGINT NOT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id     BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (product_id, tax_id)
);

-- +goose Down
DROP TABLE IF EXISTS product_taxes;
DROP TABLE IF EXISTS products;
