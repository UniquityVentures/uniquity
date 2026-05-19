-- +goose Up
CREATE TABLE IF NOT EXISTS payment_taxes (
    payment_id BIGINT NOT NULL REFERENCES payments (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id     BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (payment_id, tax_id)
);

CREATE INDEX IF NOT EXISTS idx_payment_taxes_tax_id ON payment_taxes (tax_id);

-- +goose Down
DROP TABLE IF EXISTS payment_taxes;
