package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewRoleBuilder(name string) *Builder {
	return &Builder{
		entityType: RoleType,
		name:       name,
	}
}

type Role struct {
	Name    sql.NullString `db:"name"`
	Comment sql.NullString `db:"comment"`
	Owner   sql.NullString `db:"owner"`
}

func ScanRole(row *sqlx.Row) (*Role, error) {
	r := &Role{}
	err := row.StructScan(r)
	return r, err
}

func ListRoles(db *sql.DB, rolePattern string) ([]*Role, error) {
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

	roles := []*Role{}
	if err := sqlx.StructScan(rows, &roles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no roles found")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan stmt = %v err = %w", stmt, err)
	}
	return roles, nil
}
