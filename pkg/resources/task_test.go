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

func TestTask(t *testing.T) {
	r := require.New(t)
	err := resources.Task().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTaskCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"enabled":           true,
		"name":              "test_task",
		"database":          "test_db",
		"schema":            "test_schema",
		"warehouse":         "much_warehouse",
		"sql_statement":     "select hi from hello",
		"comment":           "wow comment",
		"error_integration": "test_notification_integration",
	}

	d := schema.TestResourceDataRaw(t, resources.Task().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TASK "test_db"."test_schema"."test_task" WAREHOUSE = "much_warehouse" COMMENT = 'wow comment' ERROR_INTEGRATION = 'test_notification_integration' AS select hi from hello$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(
			`^ALTER TASK "test_db"."test_schema"."test_task" RESUME$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadTask(mock)
		expectReadTaskParams(mock)
		err := resources.CreateTask(d, db)
		r.NoError(err)

		r.Empty(d.Get("error_integration"), "Null string must be treated as empty")
	})
}

func TestTaskCreateManagedWithInitSize(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"enabled":           true,
		"name":              "test_task",
		"database":          "test_db",
		"schema":            "test_schema",
		"sql_statement":     "select hi from hello",
		"comment":           "wow comment",
		"error_integration": "test_notification_integration",
		"user_task_managed_initial_warehouse_size": "XSMALL",
	}

	d := schema.TestResourceDataRaw(t, resources.Task().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TASK "test_db"."test_schema"."test_task" USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'XSMALL' COMMENT = 'wow comment' ERROR_INTEGRATION = 'test_notification_integration' AS select hi from hello$`,
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

func TestTaskCreateManagedWithoutInitSize(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"enabled":       true,
		"name":          "test_task",
		"database":      "test_db",
		"schema":        "test_schema",
		"sql_statement": "select hi from hello",
		"comment":       "wow comment",
	}

	d := schema.TestResourceDataRaw(t, resources.Task().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TASK "test_db"."test_schema"."test_task" COMMENT = 'wow comment' AS select hi from hello$`,
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
		"created_on", "name", "database_name", "schema_name", "owner", "comment", "warehouse", "schedule", "predecessors", "state", "definition", "condition", "error_integration"},
	).AddRow("2020-05-14 17:20:50.088 +0000", "test_task", "test_db", "test_schema", "ACCOUNTADMIN", "wow comment", "", "", "", "started", "select hi from hello", "", "null")
	mock.ExpectQuery(`^SHOW TASKS LIKE 'test_task' IN SCHEMA "test_db"."test_schema"$`).WillReturnRows(rows)
}

func expectReadTaskParams(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"key", "value", "default", "level", "description", "type"},
	).AddRow("ABORT_DETACHED_QUERY", "false", "false", "", "wow desc", "BOOLEAN")
	mock.ExpectQuery(`^SHOW PARAMETERS IN TASK "test_db"."test_schema"."test_task"$`).WillReturnRows(rows)
}

func TestTaskRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_task",
		"database": "test_db",
		"schema":   "test_schema",
	}

	d := task(t, "test_db|test_schema|test_task", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Task("test_task", "test_db", "test_schema").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadTask(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}
