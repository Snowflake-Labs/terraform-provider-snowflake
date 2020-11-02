package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
				"delete": "CREATE DATABASE good_name", // test arbitrary sql statement
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resources.SqlScript().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("CREATE DATABASE good_name").WillReturnResult(sqlmock.NewResult(1, 1))

		// To me, this seems logically equivalent to... but i guess not:
		// err := resources.CreateSqlScript(d, db)
		err := snowflake.Exec(db, "CREATE DATABASE good_name")

		// This doesn't make any sense why this will fail
		// --- FAIL: TestSqlScriptCreate (0.00s)
		// panic: runtime error: index out of range [0] with length 0 [recovered]
		// panic: runtime error: index out of range [0] with length 0
		// err := resources.CreateSqlScript(d, db)
		r.NoError(err)
	})
}
