-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    entity_id  BIGINT NOT NULL REFERENCES entities (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(64) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products (deleted_at);
CREATE INDEX IF NOT EXISTS idx_products_entity_id ON products (entity_id);

-- +goose Down
DROP INDEX IF EXISTS idx_products_entity_id;
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP TABLE IF EXISTS products;
