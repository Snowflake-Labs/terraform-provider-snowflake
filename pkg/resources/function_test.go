package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

const functionBody string = "hi"

func prepDummyFunctionResource(t *testing.T) *schema.ResourceData {
	argument1 := map[string]interface{}{"name": "data", "type": "varchar"}
	argument2 := map[string]interface{}{"name": "event_dt", "type": "date"}
	in := map[string]interface{}{
		"name":            "my_funct",
		"database":        "my_db",
		"schema":          "my_schema",
		"arguments":       []interface{}{argument1, argument2},
		"return_type":     "varchar",
		"return_behavior": "IMMUTABLE",
		"statement":       functionBody, //var message = DATA + DATA;return message
	}
	d := function(t, "my_db|my_schema|my_funct|VARCHAR-DATE", in)
	return d
}

func TestFunction(t *testing.T) {
	r := require.New(t)
	err := resources.Function().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestFunctionCreate(t *testing.T) {
	r := require.New(t)
	d := prepDummyFunctionResource(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE OR REPLACE FUNCTION "my_db"."my_schema"."my_funct"\(data VARCHAR, event_dt DATE\) RETURNS VARCHAR CALLED ON NULL INPUT IMMUTABLE COMMENT = 'user-defined function' AS \$\$hi\$\$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectFunctionRead(mock)
		err := resources.CreateFunction(d, db)
		r.NoError(err)
		r.Equal("my_funct", d.Get("name").(string))
		r.Equal("VARCHAR", d.Get("return_type").(string))
		r.Equal("user-defined function", d.Get("comment").(string))
	})
}

func expectFunctionRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure"}).
		AddRow("now", "my_funct", "my_schema", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "my_db", "N", "N", "N")
	mock.ExpectQuery(`SHOW USER FUNCTIONS LIKE 'my_funct' IN SCHEMA "my_db"."my_schema"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"property", "value"}).
		AddRow("signature", "(data VARCHAR, event_dt DATE)").
		AddRow("returns", "VARCHAR(123456789)"). // This is how return type is stored in Snowflake DB
		AddRow("language", "SQL").
		AddRow("null handling", "CALLED ON NULL INPUT").
		AddRow("volatility", "IMMUTABLE").
		AddRow("body", functionBody)

	mock.ExpectQuery(`DESCRIBE FUNCTION "my_db"."my_schema"."my_funct"\(VARCHAR, DATE\)`).WillReturnRows(describeRows)
}

func TestFunctionRead(t *testing.T) {
	r := require.New(t)

	d := prepDummyFunctionResource(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectFunctionRead(mock)

		err := resources.ReadFunction(d, db)
		r.NoError(err)
		r.Equal("my_funct", d.Get("name").(string))
		r.Equal("user-defined function", d.Get("comment").(string))
		r.Equal("VARCHAR", d.Get("return_type").(string))
		r.Equal(functionBody, d.Get("statement").(string))

		args := d.Get("arguments").([]interface{})
		r.Len(args, 2)
		test_funct_arg1 := args[0].(map[string]interface{})
		test_funct_arg2 := args[1].(map[string]interface{})
		r.Len(test_funct_arg1, 2)
		r.Len(test_funct_arg2, 2)
		r.Equal("data", test_funct_arg1["name"].(string))
		r.Equal("VARCHAR", test_funct_arg1["type"].(string))
		r.Equal("event_dt", test_funct_arg2["name"].(string))
		r.Equal("DATE", test_funct_arg2["type"].(string))
	})
}

func TestFunctionDelete(t *testing.T) {
	r := require.New(t)

	d := prepDummyFunctionResource(t)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP FUNCTION "my_db"."my_schema"."my_funct"\(VARCHAR, DATE\)`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteFunction(d, db)
		r.NoError(err)
	})
}
