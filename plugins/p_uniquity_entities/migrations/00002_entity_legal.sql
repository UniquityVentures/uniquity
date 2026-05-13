-- +goose Up
ALTER TABLE entities ADD COLUMN IF NOT EXISTS slug VARCHAR(255) NOT NULL DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_entities_slug ON entities (slug) WHERE slug <> '' AND deleted_at IS NULL;

ALTER TABLE entities ADD COLUMN IF NOT EXISTS tax_id VARCHAR(50) NOT NULL DEFAULT '';

ALTER TABLE entities ADD COLUMN IF NOT EXISTS currency_id BIGINT NULL
    REFERENCES currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE entities ADD COLUMN IF NOT EXISTS logo_path VARCHAR(512) NOT NULL DEFAULT '';

ALTER TABLE entities ADD COLUMN IF NOT EXISTS phone VARCHAR(50) NOT NULL DEFAULT '';

UPDATE entities
SET currency_id = (SELECT id FROM currencies WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 1)
WHERE currency_id IS NULL
  AND EXISTS (SELECT 1 FROM currencies WHERE deleted_at IS NULL LIMIT 1);

-- Fail fast if there is no currency row to attach; seed currencies before migrating.
ALTER TABLE entities ALTER COLUMN currency_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_entities_currency_id ON entities (currency_id);

-- +goose Down
DROP INDEX IF EXISTS idx_entities_currency_id;
ALTER TABLE entities DROP CONSTRAINT IF EXISTS entities_currency_id_fkey;
DROP INDEX IF EXISTS idx_entities_slug;
ALTER TABLE entities DROP COLUMN IF EXISTS phone;
ALTER TABLE entities DROP COLUMN IF EXISTS logo_path;
ALTER TABLE entities DROP COLUMN IF EXISTS currency_id;
ALTER TABLE entities DROP COLUMN IF EXISTS tax_id;
ALTER TABLE entities DROP COLUMN IF EXISTS slug;
