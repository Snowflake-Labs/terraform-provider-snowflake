package resources_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

		expectExternalFunctionRead(mock, nil, map[string]string{"returns": "VARCHAR(123456789)"})
		err := resources.CreateExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
	})
}

func TestExternalFunctionRead(t *testing.T) {
	r := require.New(t)

	d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", map[string]interface{}{"name": "my_test_function", "arg": []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}}, "return_type": "varchar", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectExternalFunctionRead(mock, nil, map[string]string{"returns": "VARCHAR(123456789)"})

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
		expectExternalFunctionRead(mock, nil, nil)

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

func TestExternalFunctionUpdateComment(t *testing.T) {
	// r := require.New(t)

	// d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", map[string]interface{}{"name": "my_test_function", "arg": []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}}, "return_type": "variant", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		providers := providers()
		providers["snowflake"].ConfigureFunc = func(s *schema.ResourceData) (interface{}, error) {
			fmt.Println("HERE")
			panic("foobar")
			return db, nil
		}

		accName := "accName"

		resource.Test(t, resource.TestCase{
			Providers:  providers,
			IsUnitTest: true, // we've mocked the DB
			Steps: []resource.TestStep{
				{
					Config: externalFunctionConfig(accName, []string{"https://123456.execute-api.us-west-2.amazonaws.com/prod/"}, "https://123456.execute-api.us-west-2.amazonaws.com/prod/test_func"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("snowflake_external_function.test_func", "name", accName),
						resource.TestCheckResourceAttr("snowflake_external_function.test_func", "comment", "Terraform acceptance test"),
						resource.TestCheckResourceAttrSet("snowflake_external_function.test_func", "created_on"),
					),
				},
			},
		})
	})
}

// helpers
// ------------------------------
func expectExternalFunctionRead(mock sqlmock.Sqlmock, overrideShow map[string]string, overrideDescribe map[string]string) {
	// order matters
	columns := []string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure", "is_external_function", "language"}
	defaultRowFields := []string{"now", "my_test_function", "schema_name", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "database_name", "N", "N", "N", "Y", "EXTERNAL"}
	defaultRow := zip(columns, defaultRowFields)

	r := getRow(columns, defaultRow, overrideShow)

	fmt.Printf("col %d, row %d", len(columns), len(r))
	rows := sqlmock.NewRows(columns).AddRow(r...)
	mock.ExpectQuery(`SHOW EXTERNAL FUNCTIONS LIKE 'my_test_function' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)

	defaultDescribeRows := map[string]string{
		"returns":         "VARIANT",
		"null_handling":   "CALLED ON NULL INPUT",
		"volatility":      "IMMUTABLE",
		"body":            "https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function",
		"headers":         "{\"x-custom-header\":\"snowflake\"",
		"context_headers": "[\"CURRENT_TIMESTAMP\"]",
		"max_batch_rows":  "not set",
		"compression":     "AUTO",
	}

	describeRows := sqlmock.NewRows([]string{"property", "value"})
	for property, _ := range defaultDescribeRows {
		describeRows.AddRow(property, getDefault(property, defaultDescribeRows, overrideDescribe))
	}
	mock.ExpectQuery(`DESCRIBE FUNCTION "database_name"."schema_name"."my_test_function" \(varchar\)`).WillReturnRows(describeRows)
}

func getRow(columns []string, defaultRow map[string]string, override map[string]string) []driver.Value {
	ret := []driver.Value{}

	// order matters
	for _, col := range columns {
		ret = append(ret, getDefault(col, defaultRow, override))
	}

	return ret
}

func getDefault(key string, defaults map[string]string, overrides map[string]string) string {
	val, ok := overrides[key]
	if ok {
		return val
	}

	return defaults[key]
}

func zip(keys []string, values []string) map[string]string {
	if len(keys) != len(values) {
		panic("keys and values must have same length")
	}

	ret := map[string]string{}

	for i := 0; i < len(keys); i++ {
		ret[keys[i]] = values[i]
	}
	return ret
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
