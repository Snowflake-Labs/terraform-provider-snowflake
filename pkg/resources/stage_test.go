package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestStage(t *testing.T) {
	r := require.New(t)
	err := resources.Stage().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestInternalStageCreate(t *testing.T) {
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

func TestExternalStageCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_stage",
		"database": "test_db",
		"url":      "s3://com.example.bucket/prefix",
		"schema":   "test_schema",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Stage().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STAGE "test_db"."test_schema"."test_stage" URL = 's3://com.example.bucket/prefix' COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadStage(mock)
		expectReadStageShow(mock)
		err := resources.CreateStage(d, db)
		r.NoError(err)
	})
}

func expectReadStage(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"parent_property", "property", "property_type", "property_value", "property_default",
	},
	).AddRow("STAGE_LOCATION", "URL", "string", `["s3://load/test/"]`, "").
		AddRow("STAGE_CREDENTIALS", "AWS_EXTERNAL_ID", "string", "test", "").
		AddRow("STAGE_FILE_FORMAT", "FORMAT_NAME", "string", "CSV", "").
		AddRow("DIRECTORY", "ENABLED", "Boolean", true, false)
	mock.ExpectQuery(`^DESCRIBE STAGE "test_db"."test_schema"."test_stage"$`).WillReturnRows(rows)
}

func expectReadStageShow(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "url", "has_credentials", "has_encryption_key", "owner", "comment", "region", "type", "cloud", "notification_channel", "storage_integration",
	},
	).AddRow("2019-12-23 17:20:50.088 +0000", "test_stage", "test_db", "test_schema", "s3://load/test/", "N", "Y", "test", "great comment", "us-east-1", "EXTERNAL", "AWS", "NULL", "NULL")
	mock.ExpectQuery(`^SHOW STAGES LIKE 'test_stage' IN SCHEMA "test_db"."test_schema"$`).WillReturnRows(rows)
}

func TestStageRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_stage",
		"database": "test_db",
		"schema":   "test_schema",
	}
	d := stage(t, "test_db|test_schema|test_stage", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.NewStageBuilder("test_stage", "test_db", "test_schema").Describe()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadStage(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}

func TestStageUpdateWithSIAndURL(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                "test_stage",
		"database":            "test_db",
		"schema":              "test_schema",
		"url":                 "s3://changed_url",
		"storage_integration": "changed_integration",
	}

	d := stage(t, "test_db|test_schema|test_stage", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER STAGE "test_db"."test_schema"."test_stage" SET STORAGE_INTEGRATION = "changed_integration" URL = 's3://changed_url'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStage(mock)
		expectReadStageShow(mock)
		err := resources.UpdateStage(d, db)
		r.NoError(err)
	})
}

func TestStageUpdateWithJustURL(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test_stage",
		"database": "test_db",
		"schema":   "test_schema",
		"url":      "s3://changed_url",
	}

	d := stage(t, "test_db|test_schema|test_stage", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER STAGE "test_db"."test_schema"."test_stage" SET URL = 's3://changed_url'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStage(mock)
		expectReadStageShow(mock)
		err := resources.UpdateStage(d, db)
		r.NoError(err)
	})
}

func TestStageUpdateWithJustSI(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                "test_stage",
		"database":            "test_db",
		"schema":              "test_schema",
		"storage_integration": "changed_integration",
	}

	d := stage(t, "test_db|test_schema|test_stage", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER STAGE "test_db"."test_schema"."test_stage" SET STORAGE_INTEGRATION = "changed_integration"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStage(mock)
		expectReadStageShow(mock)
		err := resources.UpdateStage(d, db)
		r.NoError(err)
	})
}
