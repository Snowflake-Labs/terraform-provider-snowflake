package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	Owner   sql.NullString `db:"owner"`
}

func ScanRole(row *sqlx.Row) (*role, error) {
	r := &role{}
	err := row.StructScan(r)
	return r, err
}

func ListRoles(db *sql.DB, rolePattern string) ([]*role, error) {
	stmt := strings.Builder{}
	stmt.WriteString("SHOW ROLES")
	if rolePattern != "" {
		stmt.WriteString(fmt.Sprintf(` LIKE '%v'`, rolePattern))
	}
	rows, err := Query(db, stmt.String())
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
