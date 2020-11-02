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

func TestSqlScript(t *testing.T) {
	r := require.New(t)
	err := resources.SqlScript().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSqlScriptCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "test_sql_script",
		"lifecycle": []interface{}{
			map[string]interface{}{
				"create": "CREATE DATABASE good_name", // test arbitrary sql statement
				"delete": "DROP DATABASE good_name",   // test arbitrary sql statement
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resources.SqlScript().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("CREATE DATABASE good_name").WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.CreateSqlScript(d, db)
		// The mock works when you call snowflake.Exec directly but not via
		// resources.CreateSqlScript. Not sure why that is.
		// err := snowflake.Exec(db, "CREATE DATABASE good_name")
		r.NoError(err)
	})
}

func TestSqlScriptDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "test_sql_script",
		"lifecycle": []interface{}{
			map[string]interface{}{
				"create": "CREATE DATABASE good_name", // test arbitrary sql statement
				"delete": "DROP DATABASE good_name",   // test arbitrary sql statement
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resources.SqlScript().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("DROP DATABASE good_name").WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteSqlScript(d, db)
		// The mock works when you call snowflake.Exec directly but not via
		// resources.DeleteSqlScript. Not sure why that is.
		// err := snowflake.Exec(db, "DROP DATABASE good_name")
		r.NoError(err)
	})
}
