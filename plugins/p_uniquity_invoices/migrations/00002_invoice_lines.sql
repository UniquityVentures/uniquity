-- +goose Up
CREATE TABLE IF NOT EXISTS invoice_lines (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    invoice_id      BIGINT NOT NULL REFERENCES invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id      BIGINT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE SET NULL,
    label           VARCHAR(200) NOT NULL DEFAULT '',
    quantity        NUMERIC(12, 2) NOT NULL DEFAULT 1,
    price_unit      NUMERIC(12, 2) NOT NULL DEFAULT 0,
    discount        NUMERIC(5, 2) NOT NULL DEFAULT 0,
    account_id      BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    price_subtotal  NUMERIC(16, 2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_invoice_lines_deleted_at ON invoice_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_invoice_id ON invoice_lines (invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_product_id ON invoice_lines (product_id);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_account_id ON invoice_lines (account_id);

CREATE TABLE IF NOT EXISTS invoice_line_tax_rates (
    invoice_line_id BIGINT NOT NULL REFERENCES invoice_lines (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_rate_id     BIGINT NOT NULL REFERENCES tax_rates (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (invoice_line_id, tax_rate_id)
);

-- +goose Down
DROP TABLE IF EXISTS invoice_line_tax_rates;
DROP INDEX IF EXISTS idx_invoice_lines_account_id;
DROP INDEX IF EXISTS idx_invoice_lines_product_id;
DROP INDEX IF EXISTS idx_invoice_lines_invoice_id;
DROP INDEX IF EXISTS idx_invoice_lines_deleted_at;
DROP TABLE IF EXISTS invoice_lines;
