-- +goose Up
CREATE TABLE IF NOT EXISTS tax_rates (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    entity_id           BIGINT NOT NULL REFERENCES entities (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL,
    scope               VARCHAR(20) NOT NULL,
    amount              NUMERIC(12, 4) NOT NULL DEFAULT 0,
    account_id          BIGINT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE SET NULL,
    refund_account_id   BIGINT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_active           BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_tax_rates_deleted_at ON tax_rates (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tax_rates_entity_id ON tax_rates (entity_id);

-- +goose Down
DROP INDEX IF EXISTS idx_tax_rates_entity_id;
DROP INDEX IF EXISTS idx_tax_rates_deleted_at;
DROP TABLE IF EXISTS tax_rates;
