package snowflake

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Enabled         sql.NullBool   `db:"enabled"`
	Comment         sql.NullString `db:"comment"`
}

func ScanStorageIntegration(row *sqlx.Row) (*storageIntegration, error) {
	r := &storageIntegration{}
	err := row.StructScan(r)
	return r, err
}

func ListStorageIntegrations(db *sql.DB) ([]storageIntegration, error) {
	stmt := "SHOW STORAGE INTEGRATIONS"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []storageIntegration{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no resouce monitors found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
