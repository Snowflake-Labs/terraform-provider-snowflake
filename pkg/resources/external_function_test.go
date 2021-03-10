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
		"args":                      []interface{}{map[string]interface{}{"name": "data", "type": "varchar"}},
		"return_type":               "varchar",
		"return_behavior":           "IMMUTABLE",
		"api_integration":           "test_api_integration_01",
		"url_of_proxy_and_resource": "https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function",
	}
	d := externalFunction(t, "database_name|schema_name|my_test_function|varchar", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE EXTERNAL FUNCTION "database_name"."schema_name"."my_test_function" \(data varchar\) RETURNS varchar NULL IMMUTABLE API_INTEGRATION = 'test_api_integration_01' AS 'https://123456.execute-api.us-west-2.amazonaws.com/prod/my_test_function'`).WillReturnResult(sqlmock.NewResult(1, 1))

		expectExternalFunctionRead(mock)
		err := resources.CreateExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
	})
}

func expectExternalFunctionRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "schema_name", "is_builtin", "is_aggregate", "is_ansi", "min_num_arguments", "max_num_arguments", "arguments", "description", "catalog_name", "is_table_function", "valid_for_clustering", "is_secure", "is_external_function", "language"}).AddRow("now", "my_test_function", "schema_name", "N", "N", "N", "1", "1", "MY_TEST_FUNCTION(VARCHAR) RETURN VARCHAR", "mock comment", "database_name", "N", "N", "N", "Y", "EXTERNAL")
	mock.ExpectQuery(`SHOW EXTERNAL FUNCTIONS LIKE 'my_test_function' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
}

func TestExternalFunctionRead(t *testing.T) {
	r := require.New(t)

	d := externalFunction(t, "database_name|schema_name|my_test_function|", map[string]interface{}{"name": "my_test_function", "comment": "mock comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectExternalFunctionRead(mock)

		err := resources.ReadExternalFunction(d, db)
		r.NoError(err)
		r.Equal("my_test_function", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
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
