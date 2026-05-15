-- +goose Up
ALTER TABLE payment_term_relatives
    ADD COLUMN duration_ns BIGINT NOT NULL DEFAULT 0;

ALTER TABLE payment_term_relatives
    ALTER COLUMN duration_ns DROP DEFAULT;

-- +goose Down
ALTER TABLE payment_term_relatives DROP COLUMN IF EXISTS duration_ns;
