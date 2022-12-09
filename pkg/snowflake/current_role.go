package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func SelectCurrentRole() string {
	return `SELECT CURRENT_ROLE() AS "currentRole";`
}

type currentRole struct {
	Role string `db:"currentRole"`
}

func ScanCurrentRole(row *sqlx.Row) (*currentRole, error) {
	role := &currentRole{}
	err := row.StructScan(role)
	return role, err
}
func ReadCurrentRole(db *sql.DB) (*currentRole, error) {
	row := QueryRow(db, SelectCurrentRole())
	return ScanCurrentRole(row)
}
