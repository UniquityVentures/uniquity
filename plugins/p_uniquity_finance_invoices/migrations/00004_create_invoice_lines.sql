-- +goose Up
CREATE TABLE IF NOT EXISTS invoice_lines (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    invoice_id  BIGINT          NOT NULL REFERENCES invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id  BIGINT          NOT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity    NUMERIC(19, 6)  NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_invoice_lines_deleted_at ON invoice_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_invoice_id ON invoice_lines (invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_lines_product_id ON invoice_lines (product_id);

CREATE TABLE IF NOT EXISTS invoice_line_taxes (
    invoice_line_id BIGINT NOT NULL REFERENCES invoice_lines (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id          BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (invoice_line_id, tax_id)
);

CREATE INDEX IF NOT EXISTS idx_invoice_line_taxes_tax_id ON invoice_line_taxes (tax_id);

-- +goose Down
DROP TABLE IF EXISTS invoice_line_taxes;
DROP TABLE IF EXISTS invoice_lines;
