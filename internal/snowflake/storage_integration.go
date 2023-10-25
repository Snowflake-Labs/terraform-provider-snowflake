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

// StorageIntegration returns a pointer to a Builder that abstracts the DDL operations for a storage integration.
//
// Supported DDL operations are:
//   - CREATE STORAGE INTEGRATION
//   - ALTER STORAGE INTEGRATION
//   - DROP INTEGRATION
//   - SHOW INTEGRATIONS
//   - DESCRIBE INTEGRATION
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-user-security.html#storage-integrations)
func NewStorageIntegrationBuilder(name string) *Builder {
	return &Builder{
		entityType: StorageIntegrationType,
		name:       name,
	}
}

type StorageIntegration struct {
	Name            sql.NullString `db:"name"`
	Category        sql.NullString `db:"category"`
	IntegrationType sql.NullString `db:"type"`
	CreatedOn       sql.NullString `db:"created_on"`
	Enabled         sql.NullBool   `db:"enabled"`
	Comment         sql.NullString `db:"comment"`
}

func ScanStorageIntegration(row *sqlx.Row) (*StorageIntegration, error) {
	r := &StorageIntegration{}
	err := row.StructScan(r)
	return r, err
}

func ListStorageIntegrations(db *sql.DB) ([]StorageIntegration, error) {
	stmt := "SHOW STORAGE INTEGRATIONS"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []StorageIntegration{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no resource monitors found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
