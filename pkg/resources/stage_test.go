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

func TestStage(t *testing.T) {
	r := require.New(t)
	err := resources.Stage().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStageCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_stage",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STAGE "test_db"."test_schema"."test_stage" COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadStage(mock)
		expectReadStageShow(mock)
		err := resources.CreateStage(d, db)
		r.NoError(err)
	})
}

func expectReadStage(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"parent_property", "property", "property_type", "property_value", "property_default"},
	).AddRow("STAGE_LOCATION", "URL", "string", `["s3://load/test/"]`, "").
		AddRow("STAGE_CREDENTIALS", "AWS_EXTERNAL_ID", "string", "test", "").
		AddRow("STAGE_FILE_FORMAT", "FORMAT_NAME", "string", "CSV", "")
	mock.ExpectQuery(`^DESCRIBE STAGE "test_db"."test_schema"."test_stage"$`).WillReturnRows(rows)
}

func expectReadStageShow(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "url", "has_credentials", "has_encryption_key", "owner", "comment", "region", "type", "cloud", "notification_channel", "storage_integration"},
	).AddRow("2019-12-23 17:20:50.088 +0000", "test_stage", "test_db", "test_schema", "s3://load/test/", "N", "Y", "test", "great comment", "us-east-1", "EXTERNAL", "AWS", "NULL", "NULL")
	mock.ExpectQuery(`^SHOW STAGES LIKE 'test_stage' IN DATABASE "test_db"$`).WillReturnRows(rows)
}
