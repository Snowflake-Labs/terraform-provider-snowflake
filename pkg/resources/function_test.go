package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

const functionBody string = "def add_py(i, j): return i+j"

func prepDummyFunctionResource(t *testing.T) *schema.ResourceData {
	t.Helper()
	argument1 := map[string]interface{}{"name": "data", "type": "varchar"}
	argument2 := map[string]interface{}{"name": "event_dt", "type": "date"}
	in := map[string]interface{}{
		"name":                "my_funct",
		"database":            "my_db",
		"schema":              "my_schema",
		"arguments":           []interface{}{argument1, argument2},
		"language":            "PYTHON",
		"null_input_behavior": "CALLED ON NULL INPUT",
		"return_behavior":     "VOLATILE",
		"runtime_version":     "3.8",
		"packages":            []interface{}{"numpy", "pandas"},
		"handler":             "add_py",
		"return_type":         "varchar",
		"statement":           functionBody, // var message = DATA + DATA;return message
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
		mock.ExpectExec(`CREATE OR REPLACE FUNCTION "my_db"."my_schema"."my_funct"\(data VARCHAR, event_dt DATE\) RETURNS VARCHAR LANGUAGE PYTHON CALLED ON NULL INPUT VOLATILE RUNTIME_VERSION = '3.8' PACKAGES = \('numpy', 'pandas'\) COMMENT = 'user-defined function' HANDLER = 'add_py' AS \$\$def add_py\(i, j\)\: return i\+j\$\$`).WillReturnResult(sqlmock.NewResult(1, 1))
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
		AddRow("language", "PYTHON").
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
		r.Equal("PYTHON", d.Get("language").(string))
		r.Equal("3.8", d.Get("runtime_version").(string))
		r.Equal("add_py", d.Get("handler").(string))
		r.Equal("CALLED ON NULL INPUT", d.Get("null_input_behavior").(string))
		r.Equal("VOLATILE", d.Get("return_behavior").(string))

		pkgs := d.Get("packages").([]interface{})
		r.Len(pkgs, 2)
		testFunctPkg1 := pkgs[0].(string)
		testFunctPkg2 := pkgs[1].(string)
		r.Equal("numpy", testFunctPkg1)
		r.Equal("pandas", testFunctPkg2)

		args := d.Get("arguments").([]interface{})
		r.Len(args, 2)
		testFunctArg1 := args[0].(map[string]interface{})
		testFunctArg2 := args[1].(map[string]interface{})
		r.Len(testFunctArg1, 2)
		r.Len(testFunctArg2, 2)
		r.Equal("data", testFunctArg1["name"].(string))
		r.Equal("VARCHAR", testFunctArg1["type"].(string))
		r.Equal("event_dt", testFunctArg2["name"].(string))
		r.Equal("DATE", testFunctArg2["type"].(string))
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
