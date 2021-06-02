package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// NotificationIntegration returns a pointer to a Builder that abstracts the DDL operations for a notification integration.
//
// Supported DDL operations are:
//   - CREATE NOTIFICATION INTEGRATION
//   - ALTER NOTIFICATION INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#notification-integrations)
func NotificationIntegration(name string) *Builder {
	return &Builder{
		entityType: NotificationIntegrationType,
		name:       name,
	}
}

type notificationIntegration struct {
	Name      sql.NullString `db:"name"`
	Category  sql.NullString `db:"category"`
	Type      sql.NullString `db:"type"`
	CreatedOn sql.NullString `db:"created_on"`
	Enabled   sql.NullBool   `db:"enabled"`
}

func ScanNotificationIntegration(row *sqlx.Row) (*notificationIntegration, error) {
	r := &notificationIntegration{}
	err := row.StructScan(r)
	return r, err
}
