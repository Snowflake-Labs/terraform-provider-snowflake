package snowflake

import (
	"database/sql"

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
func Share(name string) *Builder {
	return &Builder{
		entityType: ShareType,
		name:       name,
	}
}

type share struct {
	Name    sql.NullString `db:"name"`
	To      sql.NullString `db:"to"`
	Comment sql.NullString `db:"comment"`
}

func ScanShare(row *sqlx.Row) (*share, error) {
	r := &share{}
	err := row.StructScan(r)
	return r, err
}
