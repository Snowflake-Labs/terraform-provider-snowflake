package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Share returns a pointer to a Builder that abstracts the DDL operations for a share.
//
// Supported DDL operations are:
//   - CREATE SHARE
//   - ALTER SHARE
//   - DROP SHARE
//   - SHOW SHARES
//   - DESCRIBE SHARE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-database.html#share-management)

func NewShareBuilder(name string) *Builder {
	return &Builder{
		entityType: ShareType,
		name:       name,
	}
}

type Share struct {
	Name    sql.NullString `db:"name"`
	To      sql.NullString `db:"to"`
	Comment sql.NullString `db:"comment"`
	Kind    sql.NullString `db:"kind"`
	Owner   sql.NullString `db:"owner"`
}

func ScanShare(row *sqlx.Row) (*Share, error) {
	r := &Share{}
	err := row.StructScan(r)
	return r, err
}

func ListShares(db *sql.DB, sharePattern string) ([]*Share, error) {
	stmt := strings.Builder{}
	stmt.WriteString("SHOW SHARES")
	if sharePattern != "" {
		stmt.WriteString(fmt.Sprintf(` LIKE '%v'`, sharePattern))
	}
	rows, err := Query(db, stmt.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shares := []*Share{}
	if err := sqlx.StructScan(rows, &shares); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no shares found")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan stmt = %v err = %w", stmt, err)
	}
	return shares, nil
}
