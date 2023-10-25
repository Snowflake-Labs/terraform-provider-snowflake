// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestMaterializedView(t *testing.T) {
	r := require.New(t)
	err := resources.MaterializedView().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestMaterializedViewCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":      "good_name",
		"database":  "test_db",
		"schema":    "test_schema",
		"warehouse": "test_wh",
		"comment":   "great comment",
		"statement": "SELECT * FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'",
		"is_secure": true,
	}
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec(
			`^USE WAREHOUSE test_wh;$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`CREATE SECURE MATERIALIZED VIEW "test_db"."test_schema"."good_name" COMMENT = 'great comment' AS SELECT \* FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		expectReadMaterializedView(mock)
		err := resources.CreateMaterializedView(d, db)
		r.NoError(err)
	})
}

func TestMaterializedViewCreateOrReplace(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":       "good_name",
		"database":   "test_db",
		"schema":     "test_schema",
		"warehouse":  "test_wh",
		"comment":    "great comment",
		"statement":  "SELECT * FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'",
		"is_secure":  true,
		"or_replace": true,
	}
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec(
			`^USE WAREHOUSE test_wh;$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^CREATE OR REPLACE SECURE MATERIALIZED VIEW "test_db"."test_schema"."good_name" COMMENT = 'great comment' AS SELECT \* FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		expectReadMaterializedView(mock)
		err := resources.CreateMaterializedView(d, db)
		r.NoError(err)
	})
}

func TestMaterializedViewCreateAmpersand(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":      "good_name",
		"database":  "test_db",
		"schema":    "test_schema",
		"warehouse": "test_wh",
		"comment":   "great comment",
		"statement": "SELECT * FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id LIKE 'bob%'",
		"is_secure": true,
	}
	d := schema.TestResourceDataRaw(t, resources.MaterializedView().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec(
			`^USE WAREHOUSE test_wh;$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^CREATE SECURE MATERIALIZED VIEW "test_db"."test_schema"."good_name" COMMENT = 'great comment' AS SELECT \* FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id LIKE 'bob%'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		expectReadMaterializedView(mock)
		err := resources.CreateMaterializedView(d, db)
		r.NoError(err)
	})
}

func expectReadMaterializedView(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "reserved", "database_name", "schema_name", "owner", "comment", "text", "is_secure", "is_materialized",
	},
	).AddRow("2019-05-19 16:55:36.530 -0700", "good_name", "", "test_db", "GREAT_SCHEMA", "admin", "great comment", "SELECT * FROM test_db.GREAT_SCHEMA.GREAT_TABLE WHERE account_id = 'bobs-account-id'", true, true)
	mock.ExpectQuery(`^SHOW MATERIALIZED VIEWS LIKE 'good_name' IN DATABASE "test_db"$`).WillReturnRows(rows)
}

func TestMaterializedViewRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
		"schema":   "test_schema",
	}

	d := materializedView(t, "test_db|schema_name|good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.NewMaterializedViewBuilder("good_name").WithDB("test_db").WithSchema("test_schema").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadMaterializedView(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}
