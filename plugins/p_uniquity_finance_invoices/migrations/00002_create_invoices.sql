-- +goose Up
CREATE TABLE IF NOT EXISTS payment_terms (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    type        TEXT NOT NULL,
    backing_id  BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_payment_terms_deleted_at ON payment_terms (deleted_at);

CREATE TABLE IF NOT EXISTS draft_invoices (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    number                  TEXT,
    account_receivable_id   BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_revenue_id      BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_tax_payable_id  BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id              BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    datetime                TIMESTAMPTZ NOT NULL,
    customer_id             BIGINT NOT NULL REFERENCES customers (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    payment_term_type       TEXT NOT NULL,
    payment_term_id         BIGINT NOT NULL REFERENCES payment_terms (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_draft_invoices_deleted_at ON draft_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_draft_invoices_customer_id ON draft_invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_draft_invoices_datetime ON draft_invoices (datetime);

CREATE TABLE IF NOT EXISTS draft_invoice_taxes (
    draft_invoice_id BIGINT NOT NULL REFERENCES draft_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id           BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (draft_invoice_id, tax_id)
);

CREATE TABLE IF NOT EXISTS draft_invoice_lines (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    draft_invoice_id BIGINT NOT NULL REFERENCES draft_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id       BIGINT NOT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    rate             NUMERIC(19, 6) NOT NULL,
    quantity         NUMERIC(19, 6) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_draft_invoice_lines_deleted_at ON draft_invoice_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_draft_invoice_lines_draft_invoice_id ON draft_invoice_lines (draft_invoice_id);

CREATE TABLE IF NOT EXISTS draft_invoice_line_taxes (
    draft_invoice_line_id BIGINT NOT NULL REFERENCES draft_invoice_lines (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id                BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (draft_invoice_line_id, tax_id)
);

CREATE TABLE IF NOT EXISTS posted_invoices (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    draft_invoice_id        BIGINT NOT NULL REFERENCES draft_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    posted_at               TIMESTAMPTZ,
    number                  TEXT NOT NULL,
    account_receivable_id   BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_revenue_id      BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_tax_payable_id  BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id              BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    datetime                TIMESTAMPTZ NOT NULL,
    customer_id             BIGINT NOT NULL REFERENCES customers (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    payment_term_type       TEXT NOT NULL,
    payment_term_id         BIGINT NOT NULL REFERENCES payment_terms (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_entry_id        BIGINT NOT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_posted_invoices_draft_invoice_id ON posted_invoices (draft_invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_posted_invoices_deleted_at ON posted_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_posted_invoices_journal_entry_id ON posted_invoices (journal_entry_id);

CREATE TABLE IF NOT EXISTS posted_invoice_taxes (
    posted_invoice_id BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id            BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (posted_invoice_id, tax_id)
);

CREATE TABLE IF NOT EXISTS posted_invoice_lines (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    posted_invoice_id       BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id              BIGINT NOT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    rate                    NUMERIC(19, 6) NOT NULL,
    quantity                NUMERIC(19, 6) NOT NULL,
    journal_entry_item_id   BIGINT NOT NULL REFERENCES journal_entry_items (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_posted_invoice_lines_deleted_at ON posted_invoice_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_posted_invoice_lines_posted_invoice_id ON posted_invoice_lines (posted_invoice_id);

CREATE TABLE IF NOT EXISTS posted_invoice_line_taxes (
    posted_invoice_line_id BIGINT NOT NULL REFERENCES posted_invoice_lines (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id                   BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (posted_invoice_line_id, tax_id)
);

CREATE TABLE IF NOT EXISTS cancelled_invoices (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    posted_invoice_id       BIGINT NOT NULL REFERENCES posted_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    posted_at               TIMESTAMPTZ,
    cancelled_at            TIMESTAMPTZ,
    number                  TEXT NOT NULL,
    account_receivable_id   BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_revenue_id      BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    account_tax_payable_id  BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id              BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    datetime                TIMESTAMPTZ NOT NULL,
    customer_id             BIGINT NOT NULL REFERENCES customers (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    payment_term_type       TEXT NOT NULL,
    payment_term_id         BIGINT NOT NULL REFERENCES payment_terms (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    credit_note_id          BIGINT NOT NULL REFERENCES credit_notes (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_cancelled_invoices_posted_invoice_id ON cancelled_invoices (posted_invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_cancelled_invoices_deleted_at ON cancelled_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cancelled_invoices_credit_note_id ON cancelled_invoices (credit_note_id);

CREATE TABLE IF NOT EXISTS cancelled_invoice_taxes (
    cancelled_invoice_id BIGINT NOT NULL REFERENCES cancelled_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id               BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (cancelled_invoice_id, tax_id)
);

CREATE TABLE IF NOT EXISTS cancelled_invoice_lines (
    id                      BIGSERIAL PRIMARY KEY,
    created_at              TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    cancelled_invoice_id    BIGINT NOT NULL REFERENCES cancelled_invoices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id              BIGINT NOT NULL REFERENCES products (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    rate                    NUMERIC(19, 6) NOT NULL,
    quantity                NUMERIC(19, 6) NOT NULL,
    journal_entry_item_id   BIGINT NOT NULL REFERENCES journal_entry_items (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_cancelled_invoice_lines_deleted_at ON cancelled_invoice_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cancelled_invoice_lines_cancelled_invoice_id ON cancelled_invoice_lines (cancelled_invoice_id);

CREATE TABLE IF NOT EXISTS cancelled_invoice_line_taxes (
    cancelled_invoice_line_id BIGINT NOT NULL REFERENCES cancelled_invoice_lines (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tax_id                    BIGINT NOT NULL REFERENCES taxes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    PRIMARY KEY (cancelled_invoice_line_id, tax_id)
);

-- +goose Down
DROP TABLE IF EXISTS cancelled_invoice_line_taxes;
DROP TABLE IF EXISTS cancelled_invoice_lines;
DROP TABLE IF EXISTS cancelled_invoice_taxes;
DROP TABLE IF EXISTS cancelled_invoices;
DROP TABLE IF EXISTS posted_invoice_line_taxes;
DROP TABLE IF EXISTS posted_invoice_lines;
DROP TABLE IF EXISTS posted_invoice_taxes;
DROP TABLE IF EXISTS posted_invoices;
DROP TABLE IF EXISTS draft_invoice_line_taxes;
DROP TABLE IF EXISTS draft_invoice_lines;
DROP TABLE IF EXISTS draft_invoice_taxes;
DROP TABLE IF EXISTS draft_invoices;
DROP TABLE IF EXISTS payment_terms;
