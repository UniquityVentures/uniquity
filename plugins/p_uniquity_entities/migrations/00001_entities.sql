-- +goose Up
CREATE TABLE IF NOT EXISTS entities (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       VARCHAR(255) NOT NULL,
    address    TEXT,
    mobile1    VARCHAR(50),
    mobile2    VARCHAR(50),
    email      VARCHAR(255),
    website    VARCHAR(512)
);

CREATE INDEX IF NOT EXISTS idx_entities_deleted_at ON entities (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS entities;
