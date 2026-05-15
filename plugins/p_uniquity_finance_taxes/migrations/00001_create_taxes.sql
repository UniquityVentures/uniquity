-- +goose Up
CREATE TABLE IF NOT EXISTS taxes (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    name        TEXT NOT NULL,
    percentage  NUMERIC(19, 6) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_taxes_deleted_at ON taxes (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS taxes;
