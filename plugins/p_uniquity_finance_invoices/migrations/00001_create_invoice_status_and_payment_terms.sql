-- +goose Up
CREATE TABLE IF NOT EXISTS payment_term_due_dates (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    datetime    TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_payment_term_due_dates_deleted_at ON payment_term_due_dates (deleted_at);

CREATE TABLE IF NOT EXISTS payment_term_relatives (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_payment_term_relatives_deleted_at ON payment_term_relatives (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS payment_term_relatives;
DROP TABLE IF EXISTS payment_term_due_dates;
