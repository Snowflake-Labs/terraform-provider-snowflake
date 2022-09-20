package snowflake

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

func Role(name string) *Builder {
	return &Builder{
		entityType: RoleType,
		name:       name,
	}
}

type role struct {
	Name    sql.NullString `db:"name"`
	Comment sql.NullString `db:"comment"`
}

func ScanRole(row *sqlx.Row) (*role, error) {
	r := &role{}
	err := row.StructScan(r)
	return r, err
}

func ListRoles(db *sql.DB) ([]*role, error) {
	stmt := "SHOW ROLES"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*role{}
	err = sqlx.StructScan(rows, &roles)
	if err == sql.ErrNoRows {
		log.Println("[DEBUG] no roles found")
		return nil, nil
	}
	return roles, nil
}
