-- +goose Up
ALTER TABLE accounting_preferences
    ADD COLUMN IF NOT EXISTS default_journal_id BIGINT REFERENCES journals (id) ON UPDATE CASCADE ON DELETE SET NULL;

-- +goose Down
ALTER TABLE accounting_preferences
    DROP COLUMN IF EXISTS default_journal_id;
