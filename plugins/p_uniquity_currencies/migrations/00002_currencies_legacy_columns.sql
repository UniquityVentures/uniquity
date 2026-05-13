-- +goose Up
-- Align legacy currencies table created by older p_uniquity_invoices migrations.
ALTER TABLE currencies ADD COLUMN IF NOT EXISTS name VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE currencies ADD COLUMN IF NOT EXISTS symbol VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE currencies ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE currencies DROP COLUMN IF EXISTS is_active;
ALTER TABLE currencies DROP COLUMN IF EXISTS symbol;
ALTER TABLE currencies DROP COLUMN IF EXISTS name;
