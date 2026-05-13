-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    code       BIGINT,
    name       TEXT NOT NULL,
    is_asset   BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts (deleted_at);

CREATE TABLE IF NOT EXISTS postings (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    amount     NUMERIC(19, 6) NOT NULL,
    account_id BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_postings_deleted_at ON postings (deleted_at);
CREATE INDEX IF NOT EXISTS idx_postings_account_id ON postings (account_id);

-- +goose Down
DROP TABLE IF EXISTS postings;
DROP TABLE IF EXISTS accounts;
