-- +goose Up
CREATE TABLE IF NOT EXISTS invoice_taxes (
    invoice_id  BIGINT NOT NULL REFERENCES invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id      BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (invoice_id, tax_id)
);

CREATE INDEX IF NOT EXISTS idx_invoice_taxes_tax_id ON invoice_taxes (tax_id);

-- +goose Down
DROP TABLE IF EXISTS invoice_taxes;
