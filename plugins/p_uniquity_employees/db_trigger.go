package p_uniquity_employees

import (
	"log/slog"

	"gorm.io/gorm"
)

// installPointsTransactionSuperuserTrigger adds a Postgres trigger so inserts
// are rejected unless from_user_id references a superuser (defense in depth
// alongside PointsTransaction.BeforeCreate).
//
// Each statement runs in its own Exec: Postgres rejects multiple commands in
// one prepared statement (SQLSTATE 42601).
func installPointsTransactionSuperuserTrigger(db *gorm.DB) {
	if db == nil || db.Dialector == nil || db.Dialector.Name() != "postgres" {
		return
	}
	sess := db.Session(&gorm.Session{PrepareStmt: false})

	createFn := `
CREATE OR REPLACE FUNCTION uniquity_points_transaction_check_from_superuser() RETURNS TRIGGER AS $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM users WHERE id = NEW.from_user_id AND is_superuser IS TRUE
  ) THEN
    RAISE EXCEPTION 'from_user_id must reference a superuser';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql`

	if err := sess.Exec(createFn).Error; err != nil {
		slog.Error("p_uniquity_employees: create points_transactions trigger function", "error", err)
		return
	}

	if err := sess.Exec(`DROP TRIGGER IF EXISTS uniquity_points_transaction_bi ON points_transactions`).Error; err != nil {
		slog.Error("p_uniquity_employees: drop points_transactions trigger", "error", err)
		return
	}

	createTrg := `
CREATE TRIGGER uniquity_points_transaction_bi
  BEFORE INSERT ON points_transactions
  FOR EACH ROW EXECUTE FUNCTION uniquity_points_transaction_check_from_superuser()`
	if err := sess.Exec(createTrg).Error; err != nil {
		slog.Error("p_uniquity_employees: create points_transactions trigger", "error", err)
	}
}
