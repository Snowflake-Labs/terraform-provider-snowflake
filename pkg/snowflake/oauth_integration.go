package snowflake

import (
	"database/sql"
	"log"

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

func ListIntegrations(db *sql.DB) ([]oauthIntegration, error) {
	rows, err := db.Query("SHOW INTEGRATIONS")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	r := []oauthIntegration{}
	err = sqlx.StructScan(rows, &r)
	if err == sql.ErrNoRows {
		log.Println("[DEBUG] no integrations found")
		return nil, nil
	}
	return r, nil
}

func DropIntegration(db *sql.DB, name string) error {
	stmt := OAuthIntegration(name).Drop()
	return Exec(db, stmt)
}
