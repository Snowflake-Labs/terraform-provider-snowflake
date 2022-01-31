package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	r := require.New(t)
	err := resources.Pipe().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestPipeCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_pipe",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Pipe().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE PIPE "test_db"."test_schema"."test_pipe" COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadPipe(mock)
		err := resources.CreatePipe(d, db)
		r.NoError(err)

		r.Empty(d.Get("error_integration"), "Null string must be treated as empty")
	})
}

func TestPipeRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_pipe",
		"database": "test_db",
		"schema":   "test_schema",
	}

	d := pipe(t, "test_db|test_schema|test_pipe", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Pipe("test_pipe", "test_db", "test_schema").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadPipe(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}

func expectReadPipe(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "definition", "owner", "notification_channel", "comment", "error_integration"},
	).AddRow("2019-12-23 17:20:50.088 +0000", "test_pipe", "test_db", "test_schema", "test definition", "N", "test", "great comment", "null")
	mock.ExpectQuery(`^SHOW PIPES LIKE 'test_pipe' IN SCHEMA "test_db"."test_schema"$`).WillReturnRows(rows)
}
