-- +goose Up
CREATE TABLE IF NOT EXISTS employees (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id    BIGINT NOT NULL UNIQUE REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees (deleted_at);

CREATE TABLE IF NOT EXISTS points_transactions (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    points          NUMERIC(19, 2) NOT NULL,
    from_user_id    BIGINT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    to_employee_id  BIGINT NOT NULL REFERENCES employees (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_points_transactions_deleted_at ON points_transactions (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS points_transactions;
DROP TABLE IF EXISTS employees;
