package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExternalFunction(t *testing.T) {
	r := require.New(t)
	err := resources.ExternalFunction().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestExternalFunctionCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                      "my_test_function",
		"database":                  "database_name",
		"schema":                    "schema_name",
		"arg":                       []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}},
		"return_type":               "varchar",
		"return_behavior":           "IMMUTABLE",
		"api_integration":           "test_api_integration_01",
		"header":                    []interface{}{map[string]interface{}{"name": "x-custom-header", "value": "snowflake"}},
		"context_headers":           []interface{}{"current_timestamp"},
		"url_of_proxy_and_resource": "https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function",
	}
	d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE EXTERNAL FUNCTION "database_name"."schema_name"."my_test_function" \(data varchar\) RETURNS varchar NULL CALLED ON NULL INPUT IMMUTABLE COMMENT = 'user-defined function' API_INTEGRATION = 'test_api_integration_01' HEADERS = \('x-custom-header' = 'snowflake'\) CONTEXT_HEADERS = \(current_timestamp\) COMPRESSION = 'AUTO' AS 'https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function'`).WillReturnResult(sqlmock.NewResult(1, 1))

		expectExternalFunctionRead(mock)
		err := resources.CreateExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
	})
}

func expectExternalFunctionRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure", "is_external_function", "language"}).AddRow("now", "my_test_function", "schema_name", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "database_name", "N", "N", "N", "Y", "EXTERNAL")
	mock.ExpectQuery(`SHOW EXTERNAL FUNCTIONS LIKE 'my_test_function' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"property", "value"}).
		AddRow("returns", "VARCHAR(123456789)"). // This is how return type is stored in Snowflake DB
		AddRow("null handling", "CALLED ON NULL INPUT").
		AddRow("volatility", "IMMUTABLE").
		AddRow("body", "https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function").
		AddRow("headers", "{\"x-custom-header\":\"snowflake\"").
		AddRow("context_headers", "[\"CURRENT_TIMESTAMP\"]").
		AddRow("max_batch_rows", "not set").
		AddRow("compression", "AUTO")

	mock.ExpectQuery(`DESCRIBE FUNCTION "database_name"."schema_name"."my_test_function" \(varchar\)`).WillReturnRows(describeRows)
}

func expectExternalFunctionReadVariant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure", "is_external_function", "language"}).AddRow("now", "my_test_function", "schema_name", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "database_name", "N", "N", "N", "Y", "EXTERNAL")
	mock.ExpectQuery(`SHOW EXTERNAL FUNCTIONS LIKE 'my_test_function' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)

	describeRows := sqlmock.NewRows([]string{"property", "value"}).
		AddRow("returns", "VARIANT"). // VARIANTs are different format
		AddRow("null handling", "CALLED ON NULL INPUT").
		AddRow("volatility", "IMMUTABLE").
		AddRow("body", "https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function").
		AddRow("headers", "{\"x-custom-header\":\"snowflake\"").
		AddRow("context_headers", "[\"CURRENT_TIMESTAMP\"]").
		AddRow("max_batch_rows", "not set").
		AddRow("compression", "AUTO")

	mock.ExpectQuery(`DESCRIBE FUNCTION "database_name"."schema_name"."my_test_function" \(varchar\)`).WillReturnRows(describeRows)
}

func TestExternalFunctionRead(t *testing.T) {
	r := require.New(t)

	d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", map[string]interface{}{"name": "my_test_function", "arg": []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}}, "return_type": "varchar", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectExternalFunctionRead(mock)

		err := resources.ReadExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal("VARCHAR", d.Get("return_type").(string))

		args := d.Get("arg").([]interface{})
		r.Len(args, 1)
		test_func_args := args[0].(map[string]interface{})
		r.Len(test_func_args, 2)
		r.Equal("data", test_func_args["name"].(string))
		r.Equal("varchar", test_func_args["type"].(string))

		headers := d.Get("header").(*schema.Set).List()
		r.Len(headers, 1)
		test_func_headers := headers[0].(map[string]interface{})
		r.Len(test_func_headers, 2)
		r.Equal("x-custom-header", test_func_headers["name"].(string))
		r.Equal("snowflake", test_func_headers["value"].(string))

		context_headers := d.Get("context_headers").([]interface{})
		r.Len(context_headers, 1)
		test_func_context_headers := expandStringList(context_headers)
		r.Len(test_func_context_headers, 1)
		r.Equal("CURRENT_TIMESTAMP", test_func_context_headers[0])
	})
}
func TestExternalFunctionReadReturnTypeVariant(t *testing.T) {
	r := require.New(t)

	d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", map[string]interface{}{"name": "my_test_function", "arg": []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}}, "return_type": "variant", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectExternalFunctionReadVariant(mock)

		err := resources.ReadExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal("VARIANT", d.Get("return_type").(string))

		args := d.Get("arg").([]interface{})
		r.Len(args, 1)
		test_func_args := args[0].(map[string]interface{})
		r.Len(test_func_args, 2)
		r.Equal("data", test_func_args["name"].(string))
		r.Equal("varchar", test_func_args["type"].(string))

		headers := d.Get("header").(*schema.Set).List()
		r.Len(headers, 1)
		test_func_headers := headers[0].(map[string]interface{})
		r.Len(test_func_headers, 2)
		r.Equal("x-custom-header", test_func_headers["name"].(string))
		r.Equal("snowflake", test_func_headers["value"].(string))

		context_headers := d.Get("context_headers").([]interface{})
		r.Len(context_headers, 1)
		test_func_context_headers := expandStringList(context_headers)
		r.Len(test_func_context_headers, 1)
		r.Equal("CURRENT_TIMESTAMP", test_func_context_headers[0])
	})
}

func TestExternalFunctionDelete(t *testing.T) {
	r := require.New(t)

	d := externalFunction(t, "database_name|schema_name|drop_it|", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP FUNCTION "database_name"."schema_name"."drop_it" ()`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteExternalFunction(d, db)
		r.NoError(err)
	})
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}
