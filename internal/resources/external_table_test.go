// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestExternalTable(t *testing.T) {
	r := require.New(t)
	err := resources.ExternalTable().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestExternalTableCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":        "good_name",
		"database":    "database_name",
		"schema":      "schema_name",
		"comment":     "great comment",
		"column":      []interface{}{map[string]interface{}{"name": "column1", "type": "OBJECT", "as": "a"}, map[string]interface{}{"name": "column2", "type": "VARCHAR", "as": "b"}},
		"location":    "location",
		"file_format": "FORMAT_NAME = 'format'",
		"pattern":     "pattern",
	}
	d := externalTable(t, "database_name|schema_name|good_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE EXTERNAL TABLE "database_name"."schema_name"."good_name" \("column1" OBJECT AS a, "column2" VARCHAR AS b\) WITH LOCATION = location REFRESH_ON_CREATE = true AUTO_REFRESH = true PATTERN = 'pattern' FILE_FORMAT = \( FORMAT_NAME = 'format' \) COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))

		expectExternalTableRead(mock)
		err := resources.CreateExternalTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
	})
}

func expectExternalTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "type", "kind", "null?", "default", "primary key", "unique key", "check", "expression", "comment"}).AddRow("good_name", "VARCHAR()", "COLUMN", "Y", "NULL", "NULL", "N", "N", "NULL", "mock comment")
	mock.ExpectQuery(`SHOW EXTERNAL TABLES LIKE 'good_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
}

func TestExternalTableRead(t *testing.T) {
	r := require.New(t)

	d := externalTable(t, "database_name|schema_name|good_name", map[string]interface{}{"name": "good_name", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectExternalTableRead(mock)

		err := resources.ReadExternalTable(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
	})
}

func TestExternalTableDelete(t *testing.T) {
	r := require.New(t)

	d := externalTable(t, "database_name|schema_name|drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP EXTERNAL TABLE "database_name"."schema_name"."drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteExternalTable(d, db)
		r.NoError(err)
	})
}
