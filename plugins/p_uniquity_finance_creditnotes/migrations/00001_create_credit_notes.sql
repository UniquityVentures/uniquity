-- +goose Up
CREATE TABLE credit_notes (
    id                          BIGSERIAL PRIMARY KEY,
    created_at                  TIMESTAMPTZ,
    updated_at                  TIMESTAMPTZ,
    deleted_at                  TIMESTAMPTZ,
    datetime                    TIMESTAMPTZ NOT NULL,
    reason                      TEXT,
    journal_entry_id            BIGINT NOT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    reversed_journal_entry_id  BIGINT NOT NULL REFERENCES journal_entries (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX idx_credit_notes_deleted_at ON credit_notes (deleted_at);
CREATE INDEX idx_credit_notes_journal_entry_id ON credit_notes (journal_entry_id);
CREATE INDEX idx_credit_notes_reversed_journal_entry_id ON credit_notes (reversed_journal_entry_id);

-- +goose Down
DROP TABLE IF EXISTS credit_notes;
