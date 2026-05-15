-- +goose Up
CREATE TABLE journal_entries (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    datetime        TIMESTAMPTZ NOT NULL,
    source_doc_id   BIGINT NOT NULL REFERENCES source_docs (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_id      BIGINT NOT NULL REFERENCES journals (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX idx_journal_entries_deleted_at ON journal_entries (deleted_at);
CREATE INDEX idx_journal_entries_source_doc_id ON journal_entries (source_doc_id);
CREATE INDEX idx_journal_entries_journal_id ON journal_entries (journal_id);

-- +goose Down
DROP TABLE IF EXISTS journal_entries;
