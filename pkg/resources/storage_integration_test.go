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

func TestStorageIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.StorageIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStorageIntegrationCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                      "test_storage_integration",
		"comment":                   "great comment",
		"storage_allowed_locations": []interface{}{"s3://great-bucket/great-path/"},
		"storage_provider":          "S3",
		"storage_aws_role_arn":      "we-should-probably-validate-this-string",
	}
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STORAGE INTEGRATION "test_storage_integration" COMMENT='great comment' STORAGE_AWS_ROLE_ARN='we-should-probably-validate-this-string' STORAGE_PROVIDER='S3' TYPE='EXTERNAL_STAGE' STORAGE_ALLOWED_LOCATIONS=\('s3://great-bucket/great-path/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStorageIntegration(mock)

		err := resources.CreateStorageIntegration(d, db)
		r.NoError(err)
	})
}

func TestStorageIntegrationRead(t *testing.T) {
	r := require.New(t)

	d := storageIntegration(t, "test_storage_integration", map[string]interface{}{"name": "test_storage_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStorageIntegration(mock)

		err := resources.ReadStorageIntegration(d, db)
		r.NoError(err)
	})
}

func TestStorageIntegrationDelete(t *testing.T) {
	r := require.New(t)

	d := storageIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP STORAGE INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteStorageIntegration(d, db)
		r.NoError(err)
	})
}

func expectReadStorageIntegration(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on"},
	).AddRow("test_storage_integration", "EXTERNAL_STAGE", "STORAGE", true, "now")
	mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS LIKE 'test_storage_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("STORAGE_PROVIDER", "String", "S3", nil).
		AddRow("STORAGE_ALLOWED_LOCATIONS", "List", "s3://bucket-a/path-a/,s3://bucket-b/", nil).
		AddRow("STORAGE_BLOCKED_LOCATIONS", "List", "s3://bucket-c/path-c/,s3://bucket-d/", nil).
		AddRow("STORAGE_AWS_IAM_USER_ARN", "String", "arn:aws:iam::000000000000:/user/test", nil).
		AddRow("STORAGE_AWS_ROLE_ARN", "String", "arn:aws:iam::000000000001:/role/test", nil).
		AddRow("STORAGE_AWS_EXTERNAL_ID", "String", "AGreatExternalID", nil)

	mock.ExpectQuery(`DESCRIBE STORAGE INTEGRATION "test_storage_integration"$`).WillReturnRows(descRows)
}
