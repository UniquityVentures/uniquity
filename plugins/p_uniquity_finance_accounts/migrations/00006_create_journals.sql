-- +goose Up
-- Extend later: ALTER TYPE "JournalType" ADD VALUE 'NewKind';
-- +goose StatementBegin
DO $$
BEGIN
  CREATE TYPE "JournalType" AS ENUM ('General');
EXCEPTION
  WHEN duplicate_object THEN
    NULL;
END;
$$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS journals (
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    name          TEXT          NOT NULL,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    currency_id   BIGINT        NOT NULL REFERENCES currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    journal_type  "JournalType" NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_journals_deleted_at ON journals (deleted_at);
CREATE INDEX IF NOT EXISTS idx_journals_currency_id ON journals (currency_id);

-- +goose Down
DROP TABLE IF EXISTS journals;
DROP TYPE IF EXISTS "JournalType";
