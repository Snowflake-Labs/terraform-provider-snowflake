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
func NewAPIIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: APIIntegrationType,
		name:       name,
	}
}

type APIIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Comment         sql.NullString `db:"comment"`
	Enabled         sql.NullBool   `db:"enabled"`
}

func ScanAPIIntegration(row *sqlx.Row) (*APIIntegration, error) {
	r := &APIIntegration{}
	err := row.StructScan(r)
	return r, err
}
