-- +goose Up
CREATE TABLE IF NOT EXISTS source_docs (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    source_doc_type  TEXT NOT NULL,
    source_doc_id    BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_source_docs_deleted_at ON source_docs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_source_docs_type_id ON source_docs (source_doc_type, source_doc_id);

-- +goose Down
DROP TABLE IF EXISTS source_docs;
