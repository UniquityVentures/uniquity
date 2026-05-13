-- +goose Up
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS entity_id BIGINT NULL
    REFERENCES entities (id) ON UPDATE CASCADE ON DELETE CASCADE;

UPDATE accounts SET entity_id = (
    SELECT id FROM entities WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 1
) WHERE entity_id IS NULL
  AND EXISTS (SELECT 1 FROM entities WHERE deleted_at IS NULL LIMIT 1);

ALTER TABLE accounts ALTER COLUMN entity_id SET NOT NULL;

ALTER TABLE accounts ADD COLUMN IF NOT EXISTS code_str VARCHAR(20) NULL;

UPDATE accounts SET code_str = CASE WHEN code IS NULL THEN '' ELSE TRIM(BOTH FROM code::text) END;

ALTER TABLE accounts DROP COLUMN code;
ALTER TABLE accounts RENAME COLUMN code_str TO code;

ALTER TABLE accounts DROP COLUMN IF EXISTS is_asset;

ALTER TABLE accounts ADD COLUMN IF NOT EXISTS account_type VARCHAR(30) NOT NULL DEFAULT 'asset_cash';

ALTER TABLE accounts ADD COLUMN IF NOT EXISTS currency_id BIGINT NULL
    REFERENCES currencies (id) ON UPDATE CASCADE ON DELETE SET NULL;

ALTER TABLE accounts ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS is_reconcilable BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_accounts_entity_id ON accounts (entity_id);

CREATE UNIQUE INDEX IF NOT EXISTS p_accounting_account_entity_code_uniq
    ON accounts (entity_id, code) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS p_accounting_account_entity_code_uniq;
DROP INDEX IF EXISTS idx_accounts_entity_id;
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_currency_id_fkey;
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_entity_id_fkey;
ALTER TABLE accounts DROP COLUMN IF EXISTS is_reconcilable;
ALTER TABLE accounts DROP COLUMN IF EXISTS is_active;
ALTER TABLE accounts DROP COLUMN IF EXISTS currency_id;
ALTER TABLE accounts DROP COLUMN IF EXISTS account_type;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS is_asset BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS code_restore BIGINT NULL;
UPDATE accounts SET code_restore = 0 WHERE 1=1;
ALTER TABLE accounts DROP COLUMN IF EXISTS code;
ALTER TABLE accounts RENAME COLUMN code_restore TO code;
ALTER TABLE accounts DROP COLUMN IF EXISTS entity_id;
