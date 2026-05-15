-- +goose Up
CREATE TABLE IF NOT EXISTS invoices (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    number              TEXT            NOT NULL,
    datetime            TIMESTAMPTZ     NOT NULL,
    customer_id         BIGINT          NOT NULL REFERENCES customers (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    payment_term_type   TEXT            NOT NULL,
    payment_term_id     BIGINT          NOT NULL,
    status              "InvoiceStatus" NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_invoices_deleted_at ON invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoices_customer_id ON invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_invoices_payment_term ON invoices (payment_term_type, payment_term_id);

-- +goose Down
DROP TABLE IF EXISTS invoices;
