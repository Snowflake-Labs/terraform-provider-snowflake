package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStage(t *testing.T) {
	r := require.New(t)
	err := resources.Stage().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStageCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":     "test_stage",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STAGE "test_db"."test_schema"."test_stage" COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadStage(mock)
		err := resources.CreateStage(d, db)
		a.NoError(err)
	})
}

func expectReadStage(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"parent_property", "property", "property_type", "property_value", "property_default"},
	).AddRow("STAGE_LOCATION", "URL", "string", `["s3://load/test/"]`, "").AddRow("STAGE_CREDENTIALS", "AWS_EXTERNAL_ID", "string", "test", "")
	mock.ExpectQuery(`^DESCRIBE STAGE "test_db"."test_schema"."test_stage"$`).WillReturnRows(rows)
}
