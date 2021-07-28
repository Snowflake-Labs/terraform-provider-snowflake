package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestRowAccessPolicy(t *testing.T) {
	r := require.New(t)
	err := resources.RowAccessPolicy().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestRowAccessPolicyCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                  "policy_name",
		"database":              "database_name",
		"schema":                "schema_name",
		"comment":               "great comment",
		"signature":             map[string]interface{}{"n": "VARCHAR", "v": "VARCHAR"},
		"row_access_expression": "case when current_role() in ('ANALYST') then true else false end",
	}

	d := rowAccessPolicy(t, "database_name|schema_name|policy_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE ROW ACCESS POLICY "database_name"."schema_name"."policy_name" AS \(n VARCHAR, v VARCHAR\) RETURNS BOOLEAN -> case when current_role\(\) in \('ANALYST'\) then true else false end COMMENT = \'great comment\'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRowAccessPolicy(mock)
		err := resources.CreateRowAccessPolicy(d, db)
		r.NoError(err)
		r.Equal("policy_name", d.Get("name").(string))
	})
}

func expectReadRowAccessPolicy(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "kind", "owner", "comment",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "policy_name", "database_name", "schema_name", "ROW_ACCESS_POLICY", "test", "this is a comment",
	)
	mock.ExpectQuery(`^SHOW ROW ACCESS POLICIES LIKE 'policy_name' IN SCHEMA "database_name"."schema_name"$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"name", "signature", "return_type", "body",
	}).AddRow(
		"policy_name", "(n VARCHAR, v VARCHAR)", "BOOLEAN", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end",
	)
	mock.ExpectQuery(`^DESCRIBE ROW ACCESS POLICY "database_name"."schema_name"."policy_name"$`).WillReturnRows(descRows)
}

func TestRowAccessPolicyDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                  "policy_name",
		"database":              "database_name",
		"schema":                "schema_name",
		"comment":               "great comment",
		"signature":             map[string]interface{}{"n": "VARCHAR", "v": "VARCHAR"},
		"row_access_expression": "case when current_role() in ('ANALYST') then true else false end",
	}

	d := rowAccessPolicy(t, "database_name|schema_name|policy_name", in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP ROW ACCESS POLICY "database_name"."schema_name"."policy_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteRowAccessPolicy(d, db)
		r.NoError(err)
	})
}
