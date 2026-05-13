-- +goose Up
CREATE TABLE IF NOT EXISTS journals (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_journals_deleted_at ON journals (deleted_at);

CREATE TABLE IF NOT EXISTS journal_entries (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    journal_id BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_journal_entries_deleted_at ON journal_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_journal_entries_journal_id ON journal_entries (journal_id);

ALTER TABLE postings RENAME TO journal_entry_items;

ALTER INDEX idx_postings_deleted_at RENAME TO idx_journal_entry_items_deleted_at;
ALTER INDEX idx_postings_account_id RENAME TO idx_journal_entry_items_account_id;

-- Empty journal_entry_items (fresh deploy): add NOT NULL FK directly; app creates journal_entries per line on insert.
ALTER TABLE journal_entry_items
    ADD COLUMN journal_entry_id BIGINT NOT NULL
        REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_journal_entry_items_journal_entry_id ON journal_entry_items (journal_entry_id);

-- +goose Down
DROP INDEX IF EXISTS idx_journal_entry_items_journal_entry_id;

ALTER TABLE journal_entry_items DROP CONSTRAINT IF EXISTS journal_entry_items_journal_entry_id_fkey;

ALTER TABLE journal_entry_items DROP COLUMN IF EXISTS journal_entry_id;

ALTER TABLE journal_entry_items RENAME TO postings;

ALTER INDEX idx_journal_entry_items_deleted_at RENAME TO idx_postings_deleted_at;
ALTER INDEX idx_journal_entry_items_account_id RENAME TO idx_postings_account_id;

DROP TABLE IF EXISTS journal_entries;
DROP TABLE IF EXISTS journals;
