// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// NewManagedAccountBuilder returns a pointer to a Builder that abstracts the DDL
// operations for a reader account.
//
// Supported DDL operations are:
//   - CREATE MANAGED ACCOUNT
//   - DROP MANAGED ACCOUNT
//   - SHOW MANAGED ACCOUNTS
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/data-sharing-reader-create.html)
func NewManagedAccountBuilder(name string) *Builder {
	return &Builder{
		entityType: ManagedAccountType,
		name:       name,
	}
}

type ManagedAccount struct {
	Name      sql.NullString `db:"name"`
	Cloud     sql.NullString `db:"cloud"`
	Region    sql.NullString `db:"region"`
	Locator   sql.NullString `db:"locator"`
	CreatedOn sql.NullString `db:"created_on"`
	URL       sql.NullString `db:"url"`
	Comment   sql.NullString `db:"comment"`
	IsReader  bool           `db:"is_reader"`
}

func ScanManagedAccount(row *sqlx.Row) (*ManagedAccount, error) {
	a := &ManagedAccount{}
	e := row.StructScan(a)
	return a, e
}
