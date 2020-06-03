package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTask(t *testing.T) {
	r := require.New(t)
	err := resources.Task().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTaskCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"enabled":       true,
		"name":          "test_task",
		"database":      "test_db",
		"schema":        "test_schema",
		"warehouse":     "much_warehouse",
		"sql_statement": "select hi from hello",
		"comment":       "wow comment",
	}

	d := schema.TestResourceDataRaw(t, resources.Task().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "much_warehouse" COMMENT = 'wow comment' AS select hi from hello$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(
			`^ALTER TASK "test_db"."test_schema"."test_task" RESUME$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadTask(mock)
		expectReadTaskParams(mock)
		err := resources.CreateTask(d, db)
		r.NoError(err)
	})
}

func expectReadTask(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "owner", "comment", "warehouse", "schedule", "predecessors", "state", "definition", "condition"},
	).AddRow("2020-05-14 17:20:50.088 +0000", "test_task", "test_db", "test_schema", "ACCOUNTADMIN", "wow comment", "", "", "", "started", "select hi from hello", "")
	mock.ExpectQuery(`^SHOW TASKS LIKE 'test_task' IN DATABASE "test_db"$`).WillReturnRows(rows)
}

func expectReadTaskParams(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"key", "value", "default", "level", "description", "type"},
	).AddRow("ABORT_DETACHED_QUERY", "false", "false", "", "wow desc", "BOOLEAN")
	mock.ExpectQuery(`^SHOW PARAMETERS IN TASK "test_db"."test_schema"."test_task"$`).WillReturnRows(rows)
}
