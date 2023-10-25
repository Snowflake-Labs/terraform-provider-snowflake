package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// NewSCIMIntegrationBuilder returns a pointer to a Builder that abstracts the DDL operations for an api integration.
//
// Supported DDL operations are:
//   - CREATE SECURITY INTEGRATION
//   - ALTER SECURITY INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#security-integrations)
func NewSCIMIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: SecurityIntegrationType,
		name:       name,
	}
}

type SCIMIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
}

func ScanScimIntegration(row *sqlx.Row) (*SCIMIntegration, error) {
	r := &SCIMIntegration{}
	if err := row.StructScan(r); err != nil {
		return r, fmt.Errorf("error scanning struct err = %w", err)
	}
	return r, nil
}
