-- +goose Up
-- GORM maps time.Duration to BIGINT (nanoseconds).
ALTER TABLE payment_term_relatives
    ADD COLUMN IF NOT EXISTS duration BIGINT NOT NULL DEFAULT 0;

ALTER TABLE payment_term_relatives
    ALTER COLUMN duration DROP DEFAULT;

-- +goose Down
ALTER TABLE payment_term_relatives
    DROP COLUMN IF EXISTS duration;
