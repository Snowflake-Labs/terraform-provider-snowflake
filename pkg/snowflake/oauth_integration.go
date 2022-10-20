package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// OAuthIntegration returns a pointer to a Builder that abstracts the DDL operations for an api integration.
//
// Supported DDL operations are:
//   - CREATE SECURITY INTEGRATION
//   - ALTER SECURITY INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#security-integrations)
func OAuthIntegration(name string) *Builder {
	return &Builder{
		entityType: SecurityIntegrationType,
		name:       name,
	}
}

type oauthIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	Enabled         sql.NullBool   `db:"enabled"`
	Comment         sql.NullString `db:"comment"`
	CreatedOn       sql.NullString `db:"created_on"`
}

func ScanOAuthIntegration(row *sqlx.Row) (*oauthIntegration, error) {
	r := &oauthIntegration{}
	return r, errors.Wrap(row.StructScan(r), "error scanning struct")
}

func ListIntegrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW INTEGRATIONS")
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		var integration oauthIntegration
		err := rows.Scan(&integration)
		if err != nil {
			return nil, err
		}
		names = append(names, integration.Name.String)
	}
	return names, nil
}

func DropIntegration(db *sql.DB, name string) error {
	stmt := OAuthIntegration(name).Drop()
	return Exec(db, stmt)
}
