-- +goose Up
CREATE TABLE IF NOT EXISTS payment_terms (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    payment_term_type   TEXT NOT NULL,
    backing_id          BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_terms_type_backing
    ON payment_terms (payment_term_type, backing_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_payment_terms_deleted_at ON payment_terms (deleted_at);

INSERT INTO payment_terms (created_at, updated_at, payment_term_type, backing_id)
SELECT MIN(i.created_at),
       MAX(i.updated_at),
       i.payment_term_type,
       i.payment_term_id
FROM invoices i
GROUP BY i.payment_term_type, i.payment_term_id;

ALTER TABLE invoices ADD COLUMN payment_term_ref_id BIGINT;

UPDATE invoices i
SET payment_term_ref_id = pt.id
FROM payment_terms pt
WHERE i.payment_term_type = pt.payment_term_type
  AND i.payment_term_id = pt.backing_id;

ALTER TABLE invoices DROP COLUMN payment_term_type;
ALTER TABLE invoices DROP COLUMN payment_term_id;

ALTER TABLE invoices RENAME COLUMN payment_term_ref_id TO payment_term_id;

ALTER TABLE invoices
    ALTER COLUMN payment_term_id SET NOT NULL;

ALTER TABLE invoices
    ADD CONSTRAINT invoices_payment_term_id_fkey
    FOREIGN KEY (payment_term_id) REFERENCES payment_terms (id) ON UPDATE CASCADE ON DELETE RESTRICT;

DROP INDEX IF EXISTS idx_invoices_payment_term;

CREATE INDEX IF NOT EXISTS idx_invoices_payment_term_id ON invoices (payment_term_id);

-- +goose Down
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_payment_term_id_fkey;

DROP INDEX IF EXISTS idx_invoices_payment_term_id;

ALTER TABLE invoices RENAME COLUMN payment_term_id TO payment_term_ref_id;

ALTER TABLE invoices ADD COLUMN payment_term_type TEXT NOT NULL DEFAULT '';
ALTER TABLE invoices ADD COLUMN payment_term_id BIGINT NOT NULL DEFAULT 0;

UPDATE invoices i
SET payment_term_type = pt.payment_term_type,
    payment_term_id   = pt.backing_id
FROM payment_terms pt
WHERE i.payment_term_ref_id = pt.id;

ALTER TABLE invoices ALTER COLUMN payment_term_type DROP DEFAULT;
ALTER TABLE invoices ALTER COLUMN payment_term_id DROP DEFAULT;

ALTER TABLE invoices DROP COLUMN IF EXISTS payment_term_ref_id;

CREATE INDEX IF NOT EXISTS idx_invoices_payment_term ON invoices (payment_term_type, payment_term_id);

DROP INDEX IF EXISTS idx_payment_terms_type_backing;
DROP INDEX IF EXISTS idx_payment_terms_deleted_at;

DROP TABLE IF EXISTS payment_terms;
