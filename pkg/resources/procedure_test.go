package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestProcedure(t *testing.T) {
	r := require.New(t)
	err := resources.Procedure().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestProcedureCreate(t *testing.T) {
	r := require.New(t)
	//{map[string]interface{}{"name": "data", "type": "varchar"},map[string]interface{}{"name": "event_dt", "type": "date"}}
	in := map[string]interface{}{
		"name":            "my_proc",
		"database":        "my_db",
		"schema":          "my_schema",
		"arguments":       []interface{}{},
		"return_type":     "varchar",
		"return_behavior": "IMMUTABLE",
		"statement":       "var message = DATA + DATA;return message",
	}
	d := procedure(t, "my_db|my_schema|my_proc|VARCHAR-DATE", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE OR REPLACE PROCEDURE "my_db"."my_schema"."my_proc"() RETURNS VARCHAR LANGUAGE javascript CALLED ON NULL INPUT IMMUTABLE COMMENT = 'user-defined function' EXECUTE AS OWNER AS XX` + "\n" + `var message = DATA + DATA;return message` + "\nXX").WillReturnResult(sqlmock.NewResult(1, 1))
		expectProcedureRead(mock)
		err := resources.CreateProcedure(d, db)
		r.NoError(err)
		r.Equal("my_proc", d.Get("name").(string))
	})
}

func expectProcedureRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure"}).
		AddRow("now", "my_proc", "my_schema", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "my_db", "N", "N", "N")
	mock.ExpectQuery(`SHOW PROCEDURES LIKE 'my_proc' IN SCHEMA "my_db"."my_schema"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"property", "value"}).
		AddRow("signature", "(data VARCHAR)").
		AddRow("returns", "VARCHAR(123456789)"). // This is how return type is stored in Snowflake DB
		AddRow("language", "JAVASCRIPT").
		AddRow("null handling", "CALLED ON NULL INPUT").
		AddRow("volatility", "IMMUTABLE").
		AddRow("execute as", "CALLER").
		AddRow("body", "\nvar message = DATA + DATA;return message\n")

	mock.ExpectQuery(`DESCRIBE PROCEDURE "my_db"."my_schema"."my_proc"\(varchar\)`).WillReturnRows(describeRows)
}
