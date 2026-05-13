-- +goose Up
CREATE TABLE IF NOT EXISTS contacts (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_contacts_deleted_at ON contacts (deleted_at);

CREATE TABLE IF NOT EXISTS payment_terms (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_payment_terms_deleted_at ON payment_terms (deleted_at);

CREATE TABLE IF NOT EXISTS invoices (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    entity_id         BIGINT NOT NULL REFERENCES entities (id) ON UPDATE CASCADE ON DELETE CASCADE,
    number            VARCHAR(50) NULL,
    partner_id        BIGINT NOT NULL REFERENCES contacts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id        BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    invoice_type      VARCHAR(20) NOT NULL,
    state             VARCHAR(20) NOT NULL DEFAULT 'draft',
    reference         VARCHAR(100) NOT NULL DEFAULT '',
    invoice_date      DATE NOT NULL,
    payment_term_id   BIGINT NULL REFERENCES payment_terms (id) ON UPDATE CASCADE ON DELETE SET NULL,
    due_date          DATE NULL,
    currency_id       BIGINT NOT NULL REFERENCES currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    move_id           BIGINT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE SET NULL,
    amount_untaxed    NUMERIC(20, 2) NOT NULL DEFAULT 0,
    amount_tax        NUMERIC(20, 2) NOT NULL DEFAULT 0,
    amount_total      NUMERIC(20, 2) NOT NULL DEFAULT 0,
    amount_residual   NUMERIC(20, 2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_invoices_deleted_at ON invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoices_entity_id ON invoices (entity_id);
CREATE INDEX IF NOT EXISTS idx_invoices_partner_id ON invoices (partner_id);
CREATE INDEX IF NOT EXISTS idx_invoices_journal_id ON invoices (journal_id);
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_type ON invoices (invoice_type);
CREATE INDEX IF NOT EXISTS idx_invoices_number ON invoices (number);
CREATE UNIQUE INDEX IF NOT EXISTS idx_invoices_move_id ON invoices (move_id) WHERE move_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_invoices_move_id;
DROP INDEX IF EXISTS idx_invoices_number;
DROP INDEX IF EXISTS idx_invoices_invoice_type;
DROP INDEX IF EXISTS idx_invoices_journal_id;
DROP INDEX IF EXISTS idx_invoices_partner_id;
DROP INDEX IF EXISTS idx_invoices_entity_id;
DROP INDEX IF EXISTS idx_invoices_deleted_at;
DROP TABLE IF EXISTS invoices;
DROP INDEX IF EXISTS idx_payment_terms_deleted_at;
DROP TABLE IF EXISTS payment_terms;
DROP INDEX IF EXISTS idx_contacts_deleted_at;
DROP TABLE IF EXISTS contacts;
