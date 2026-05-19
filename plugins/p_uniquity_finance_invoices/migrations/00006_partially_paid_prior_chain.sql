-- +goose Up
-- 00005 used CREATE TABLE IF NOT EXISTS; existing DBs may lack the payment chain column.
ALTER TABLE partially_paid_invoices
    ADD COLUMN IF NOT EXISTS prior_partially_paid_invoice_id BIGINT REFERENCES partially_paid_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE paid_invoices
    ADD COLUMN IF NOT EXISTS prior_partially_paid_invoice_id BIGINT REFERENCES partially_paid_invoices (id) ON UPDATE CASCADE ON DELETE RESTRICT;

-- +goose Down
ALTER TABLE paid_invoices
    DROP COLUMN IF EXISTS prior_partially_paid_invoice_id;

ALTER TABLE partially_paid_invoices
    DROP COLUMN IF EXISTS prior_partially_paid_invoice_id;
