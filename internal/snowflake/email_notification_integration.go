// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// EmailNotificationIntegration returns a pointer to a Builder that abstracts the DDL operations for a email notification integration.
//
// Supported DDL operations are:
//   - CREATE NOTIFICATION INTEGRATION
//   - ALTER NOTIFICATION INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-user-security.html#notification-integrations)
func NewEmailNotificationIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: NotificationIntegrationType,
		name:       name,
	}
}

type EmailNotificationIntegration struct {
	Name    sql.NullString `db:"name"`
	Type    sql.NullString `db:"type"`
	Comment sql.NullString `db:"comment"`
	Enabled sql.NullBool   `db:"enabled"`
}

func ScanEmailNotificationIntegration(row *sqlx.Row) (*EmailNotificationIntegration, error) {
	r := &EmailNotificationIntegration{}
	err := row.StructScan(r)
	return r, err
}
