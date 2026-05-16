-- +goose Up
-- Recreate trigger with EXECUTE PROCEDURE (works on PostgreSQL 11+). EXECUTE FUNCTION is PG14+ only.
-- +goose StatementBegin
DROP TRIGGER IF EXISTS accounts_enforce_parent_balance_type_biud ON accounts;

CREATE TRIGGER accounts_enforce_parent_balance_type_biud
  BEFORE INSERT OR UPDATE ON accounts
  FOR EACH ROW EXECUTE PROCEDURE accounts_enforce_parent_balance_type();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS accounts_enforce_parent_balance_type_biud ON accounts;

CREATE TRIGGER accounts_enforce_parent_balance_type_biud
  BEFORE INSERT OR UPDATE ON accounts
  FOR EACH ROW EXECUTE PROCEDURE accounts_enforce_parent_balance_type();
-- +goose StatementEnd
