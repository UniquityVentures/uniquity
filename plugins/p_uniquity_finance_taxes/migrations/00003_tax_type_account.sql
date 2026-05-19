-- +goose Up
CREATE TYPE "TaxKind" AS ENUM ('levied', 'withholding');

ALTER TABLE taxes
    ADD COLUMN tax_type "TaxKind" NOT NULL DEFAULT 'levied',
    ADD COLUMN account_id BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_taxes_account_id ON taxes (account_id);

-- +goose Down
DROP INDEX IF EXISTS idx_taxes_account_id;

ALTER TABLE taxes
    DROP COLUMN IF EXISTS account_id,
    DROP COLUMN IF EXISTS tax_type;

DROP TYPE IF EXISTS "TaxKind";
