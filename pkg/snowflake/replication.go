package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Replication returns a pointer to a Builder that abstracts the DDL operations for a replication.
//
// Supported DDL operations are:
//   - SHOW REPLICATION DATABASES
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/database-replication-config.html)

// ReplicationBuilder is a basic builder that enables replication on databases
type ReplicationBuilder struct {
	database string
}

// DatabaseFromDatabase returns a pointer to a builder that can create a database from a source database
func Replication(database string) *ReplicationBuilder {
	return &ReplicationBuilder{
		database: database,
	}
}

type replication struct {
	Region           sql.NullString `db:"snowflake_region"`
	CreatedOn        sql.NullString `db:"created_on"`
	AccountName      sql.NullString `db:"account_name"`
	DBName           sql.NullString `db:"name"`
	Comment          sql.NullString `db:"comment"`
	IsPrimary        sql.NullBool   `db:"is_primary"`
	Primary          sql.NullString `db:"primary"`
	ReplAccounts     sql.NullString `db:"replication_allowed_to_accounts"`
	FailoverAccounts sql.NullString `db:"failover_allowed_to_accounts"`
	Org              sql.NullString `db:"organization_name"`
	AccountLocator   sql.NullString `db:"account_locator"`
}

func ScanReplication(rows *sqlx.Rows, AccName string) (*replication, error) {
	for rows.Next() {
		r := &replication{}
		err := rows.StructScan(r)
		if r.AccountName.String == AccName {
			return r, err
		}
	}
	return nil, sql.ErrNoRows
}

func (rb *ReplicationBuilder) Show() string {
	return fmt.Sprintf(`SHOW REPLICATION DATABASES LIKE '%s'`, rb.database)
}
