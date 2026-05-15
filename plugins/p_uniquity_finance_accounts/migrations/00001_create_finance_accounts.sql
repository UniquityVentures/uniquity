-- +goose Up
CREATE TYPE "BalanceType" AS ENUM ('Credit', 'Debit');

CREATE TABLE IF NOT EXISTS accounts (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    name             TEXT        NOT NULL,
    code             INTEGER     NOT NULL,
    is_group         BOOLEAN     NOT NULL DEFAULT FALSE,
    balance_type     "BalanceType" NOT NULL,
    parent_id        BIGINT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_parent_id ON accounts (parent_id);

-- +goose Down
DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS "BalanceType";
