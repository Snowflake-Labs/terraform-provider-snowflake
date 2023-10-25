package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SamlIntegration returns a pointer to a Builder that abstracts the DDL operations for a SAML2 integration.
//
// Supported DDL operations are:
//   - CREATE SECURITY INTEGRATION
//   - ALTER SECURITY INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#security-integrations)
func NewSamlIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: SecurityIntegrationType,
		name:       name,
	}
}

type SamlIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Enabled         sql.NullBool   `db:"enabled"`
}

func ScanSamlIntegration(row *sqlx.Row) (*SamlIntegration, error) {
	r := &SamlIntegration{}
	if err := row.StructScan(r); err != nil {
		return r, fmt.Errorf("error scanning struct err = %w", err)
	}
	return r, nil
}
