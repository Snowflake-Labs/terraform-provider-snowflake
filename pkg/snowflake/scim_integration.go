package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ScimIntegration returns a pointer to a Builder that abstracts the DDL operations for an api integration.
//
// Supported DDL operations are:
//   - CREATE SECURITY INTEGRATION
//   - ALTER SECURITY INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#security-integrations)
func ScimIntegration(name string) *Builder {
	return &Builder{
		entityType: SecurityIntegrationType,
		name:       name,
	}
}

type scimIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
}

func ScanScimIntegration(row *sqlx.Row) (*scimIntegration, error) {
	r := &scimIntegration{}
	if err := row.StructScan(r); err != nil {
		return r, fmt.Errorf("error scanning struct err = %w", err)
	}
	return r, nil
}
