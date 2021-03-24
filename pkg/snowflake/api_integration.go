package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ApiIntegration returns a pointer to a Builder that abstracts the DDL operations for an api integration.
//
// Supported DDL operations are:
//   - CREATE API INTEGRATION
//   - ALTER API INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#api-integrations)
func ApiIntegration(name string) *Builder {
	return &Builder{
		entityType: ApiIntegrationType,
		name:       name,
	}
}

type apiIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Enabled         sql.NullBool   `db:"enabled"`
}

func ScanApiIntegration(row *sqlx.Row) (*apiIntegration, error) {
	r := &apiIntegration{}
	err := row.StructScan(r)
	return r, err
}
