package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func SelectCurrentRole() string {
	return `SELECT CURRENT_ROLE() AS "currentRole";`
}

type CurrentRole struct {
	Role string `db:"currentRole"`
}

func ScanCurrentRole(row *sqlx.Row) (*CurrentRole, error) {
	role := &CurrentRole{}
	err := row.StructScan(role)
	return role, err
}

func ReadCurrentRole(db *sql.DB) (*CurrentRole, error) {
	row := QueryRow(db, SelectCurrentRole())
	return ScanCurrentRole(row)
}
