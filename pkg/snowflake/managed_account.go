package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ManagedAccount returns a pointer to a Builder that abstracts the DDL
// operations for a reader account.
//
// Supported DDL operations are:
//   - CREATE MANAGED ACCOUNT
//   - DROP MANAGED ACCOUNT
//   - SHOW MANAGED ACCOUNTS
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/data-sharing-reader-create.html)
func ManagedAccount(name string) *Builder {
	return &Builder{
		entityType: ManagedAccountType,
		name:       name,
	}
}

type managedAccount struct {
	Name      sql.NullString `db:"name"`
	Cloud     sql.NullString `db:"cloud"`
	Region    sql.NullString `db:"region"`
	Locator   sql.NullString `db:"locator"`
	CreatedOn sql.NullString `db:"created_on"`
	Url       sql.NullString `db:"url"`
	Comment   sql.NullString `db:"comment"`
	IsReader  bool           `db:"is_reader"`
}

func ScanManagedAccount(row *sqlx.Row) (*managedAccount, error) {
	a := &managedAccount{}
	e := row.StructScan(a)
	return a, e
}
