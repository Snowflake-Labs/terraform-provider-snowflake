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

func TestStorageIntegration(t *testing.T) {
	r := require.New(t)
	err := resources.StorageIntegration().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStorageIntegrationCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":                      "test_storage_integration",
		"comment":                   "great comment",
		"storage_allowed_locations": []interface{}{"s3://great-bucket/great-path/"},
		"storage_provider":          "S3",
		"storage_aws_role_arn":      "we-should-probably-validate-this-string",
	}
	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STORAGE INTEGRATION "test_storage_integration" COMMENT='great comment' STORAGE_AWS_ROLE_ARN='we-should-probably-validate-this-string' STORAGE_PROVIDER='S3' TYPE='EXTERNAL_STAGE' STORAGE_ALLOWED_LOCATIONS=\('s3://great-bucket/great-path/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStorageIntegration(mock)

		err := resources.CreateStorageIntegration(d, db)
		a.NoError(err)
	})
}

func TestStorageIntegrationRead(t *testing.T) {
	a := assert.New(t)

	d := storageIntegration(t, "test_storage_integration", map[string]interface{}{"name": "test_storage_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStorageIntegration(mock)

		err := resources.ReadStorageIntegration(d, db)
		a.NoError(err)
	})
}

func TestStorageIntegrationDelete(t *testing.T) {
	a := assert.New(t)

	d := storageIntegration(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP STORAGE INTEGRATION "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteStorageIntegration(d, db)
		a.NoError(err)
	})
}

func expectReadStorageIntegration(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "comment", "created_on"},
	).AddRow("test_storage_integration", "EXTERNAL_STAGE", "STORAGE", true, "great comment", "now")
	mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS LIKE 'test_storage_integration'$`).WillReturnRows(rows)
}
