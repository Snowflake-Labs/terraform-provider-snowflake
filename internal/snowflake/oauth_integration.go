// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
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
func NewOAuthIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: SecurityIntegrationType,
		name:       name,
	}
}

type OauthIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	Enabled         sql.NullBool   `db:"enabled"`
	Comment         sql.NullString `db:"comment"`
	CreatedOn       sql.NullString `db:"created_on"`
}

func ScanOAuthIntegration(row *sqlx.Row) (*OauthIntegration, error) {
	r := &OauthIntegration{}
	if err := row.StructScan(r); err != nil {
		return nil, fmt.Errorf("error scanning struct err = %w", err)
	}
	return r, nil
}

func ListIntegrations(db *sql.DB) ([]OauthIntegration, error) {
	stmt := "SHOW INTEGRATIONS"
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	r := []OauthIntegration{}
	if err := sqlx.StructScan(rows, &r); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no integrations found")
			return nil, nil
		}
		return r, fmt.Errorf("failed to scan row for %s err = %w", stmt, err)
	}
	return r, nil
}

func DropIntegration(db *sql.DB, name string) error {
	stmt := NewOAuthIntegrationBuilder(name).Drop()
	return Exec(db, stmt)
}
