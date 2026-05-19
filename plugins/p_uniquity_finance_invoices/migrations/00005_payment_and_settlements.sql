-- +goose Up
CREATE TABLE IF NOT EXISTS payments (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    posted_invoice_id   BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    amount              NUMERIC(19, 6) NOT NULL,
    account_id          BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    datetime            TIMESTAMPTZ NOT NULL,
    journal_entry_id    BIGINT NOT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_posted_invoice_id ON payments (posted_invoice_id);

CREATE TABLE IF NOT EXISTS partially_paid_invoices (
    id                              BIGSERIAL PRIMARY KEY,
    created_at                      TIMESTAMPTZ,
    updated_at                      TIMESTAMPTZ,
    deleted_at                      TIMESTAMPTZ,
    payment_id                      BIGINT NOT NULL REFERENCES payments (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    posted_invoice_id               BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    prior_partially_paid_invoice_id BIGINT REFERENCES partially_paid_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_partially_paid_invoices_payment_id ON partially_paid_invoices (payment_id);
CREATE INDEX IF NOT EXISTS idx_partially_paid_invoices_deleted_at ON partially_paid_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_partially_paid_invoices_posted_invoice_id ON partially_paid_invoices (posted_invoice_id);

CREATE TABLE IF NOT EXISTS paid_invoices (
    id                              BIGSERIAL PRIMARY KEY,
    created_at                      TIMESTAMPTZ,
    updated_at                      TIMESTAMPTZ,
    deleted_at                      TIMESTAMPTZ,
    payment_id                      BIGINT NOT NULL REFERENCES payments (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    posted_invoice_id               BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    prior_partially_paid_invoice_id BIGINT REFERENCES partially_paid_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_paid_invoices_payment_id ON paid_invoices (payment_id);
CREATE INDEX IF NOT EXISTS idx_paid_invoices_deleted_at ON paid_invoices (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS uix_paid_invoices_posted_invoice_active ON paid_invoices (posted_invoice_id)
WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS paid_invoices;
DROP TABLE IF EXISTS partially_paid_invoices;
DROP TABLE IF EXISTS payments;
