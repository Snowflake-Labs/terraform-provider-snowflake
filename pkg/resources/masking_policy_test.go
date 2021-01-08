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

func TestMaskingPolicy(t *testing.T) {
	r := require.New(t)
	err := resources.MaskingPolicy().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestMaskingPolicyCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":               "policy_name",
		"database":           "database_name",
		"schema":             "schema_name",
		"comment":            "great comment",
		"value_data_type":    "string",
		"masking_expression": "case when current_role() in ('ANALYST') then val else sha2(val, 512) end",
		"return_data_type":   "string",
	}

	d := maskingPolicy(t, "database_name|schema_name|policy_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE MASKING POLICY "database_name"."schema_name"."policy_name" AS \(VAL string\) RETURNS string -> case when current_role\(\) in \('ANALYST'\) then val else sha2\(val, 512\) end COMMENT = \'great comment\'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadMaskingPolicy(mock)
		err := resources.CreateMaskingPolicy(d, db)
		r.NoError(err)
		r.Equal("policy_name", d.Get("name").(string))
	})
}

func expectReadMaskingPolicy(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "kind", "owner", "comment",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "policy_name", "database_name", "schema_name", "MASKING_POLICY", "test", "this is a comment",
	)
	mock.ExpectQuery(`^SHOW MASKING POLICIES LIKE 'policy_name' IN SCHEMA "database_name"."schema_name"$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"name", "signature", "return_type", "body",
	}).AddRow(
		"policy_name", "(VAL VARCHAR)", "VARCHAR(16777216)", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end",
	)
	mock.ExpectQuery(`^DESCRIBE MASKING POLICY "database_name"."schema_name"."policy_name"$`).WillReturnRows(descRows)
}

func TestMaskingPolicyDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":               "policy_name",
		"database":           "database_name",
		"schema":             "schema_name",
		"comment":            "great comment",
		"value_data_type":    "string",
		"masking_expression": "case when current_role() in ('ANALYST') then val else sha2(val, 512) end",
		"return_data_type":   "string",
	}

	d := maskingPolicy(t, "database_name|schema_name|policy_name", in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^DROP MASKING POLICY "database_name"."schema_name"."policy_name"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteMaskingPolicy(d, db)
		r.NoError(err)
	})
}
