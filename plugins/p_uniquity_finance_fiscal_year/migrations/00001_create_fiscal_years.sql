-- +goose Up
CREATE TABLE IF NOT EXISTS fiscal_years (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    code        TEXT NOT NULL,
    name        TEXT NOT NULL,
    starts_at   TIMESTAMPTZ NOT NULL,
    ends_at     TIMESTAMPTZ NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fiscal_years_code ON fiscal_years (code);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_deleted_at ON fiscal_years (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS fiscal_years;
