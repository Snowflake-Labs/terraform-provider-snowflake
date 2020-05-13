package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// StorageIntegration returns a pointer to a Builder that abstracts the DDL operations for a storage integration.
//
// Supported DDL operations are:
//   - CREATE STORAGE INTEGRATION
//   - ALTER STORAGE INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-user-security.html#storage-integrations)
func StorageIntegration(name string) *Builder {
	return &Builder{
		entityType: StorageIntegrationType,
		name:       name,
	}
}

type storageIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"integration_type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Enabled         sql.NullBool   `db:"enabled"`
}

func ScanStorageIntegration(row *sqlx.Row) (*storageIntegration, error) {
	r := &storageIntegration{}
	err := row.StructScan(r)
	return r, err
}
