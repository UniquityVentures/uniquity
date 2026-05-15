-- +goose Up
CREATE TABLE journal_entry_items (
    id                 BIGSERIAL PRIMARY KEY,
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ,
    deleted_at         TIMESTAMPTZ,
    datetime           TIMESTAMPTZ NOT NULL,
    account_id         BIGINT NOT NULL REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    amount             NUMERIC(19, 6) NOT NULL,
    journal_entry_id   BIGINT NOT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX idx_journal_entry_items_deleted_at ON journal_entry_items (deleted_at);
CREATE INDEX idx_journal_entry_items_account_id ON journal_entry_items (account_id);
CREATE INDEX idx_journal_entry_items_journal_entry_id ON journal_entry_items (journal_entry_id);

-- +goose Down
DROP TABLE IF EXISTS journal_entry_items;
