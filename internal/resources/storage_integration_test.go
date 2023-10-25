// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		"storage_aws_object_acl":    "bucket-owner-full-control",
	}

	in2 := map[string]interface{}{
		"name":                      "test_storage_integration_with_s3gov",
		"comment":                   "great comment",
		"storage_allowed_locations": []interface{}{"s3://great-bucket/great-path/"},
		"storage_provider":          "S3GOV",
		"storage_aws_role_arn":      "we-should-probably-validate-this-string",
		"storage_aws_object_acl":    "bucket-owner-full-control",
	}

	d := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, in)
	d2 := schema.TestResourceDataRaw(t, resources.StorageIntegration().Schema, in2)

	r.NotNil(d)
	r.NotNil(d2)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STORAGE INTEGRATION "test_storage_integration" COMMENT='great comment' STORAGE_AWS_OBJECT_ACL='bucket-owner-full-control' STORAGE_AWS_ROLE_ARN='we-should-probably-validate-this-string' STORAGE_PROVIDER='S3' TYPE='EXTERNAL_STAGE' STORAGE_ALLOWED_LOCATIONS=\('s3://great-bucket/great-path/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStorageIntegration(mock)

		err := resources.CreateStorageIntegration(d, db)
		r.NoError(err)
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE STORAGE INTEGRATION "test_storage_integration_with_s3gov" COMMENT='great comment' STORAGE_AWS_OBJECT_ACL='bucket-owner-full-control' STORAGE_AWS_ROLE_ARN='we-should-probably-validate-this-string' STORAGE_PROVIDER='S3GOV' TYPE='EXTERNAL_STAGE' STORAGE_ALLOWED_LOCATIONS=\('s3://great-bucket/great-path/'\) ENABLED=true$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStorageIntegrationWithS3GOV(mock)

		err := resources.CreateStorageIntegration(d2, db)
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

func TestStorageIntegrationReadEmpty(t *testing.T) {
	r := require.New(t)

	d := storageIntegration(t, "test_storage_integration", map[string]interface{}{"name": "not_existing_storage_integration"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStorageIntegrationEmpty(mock)

		err := resources.ReadStorageIntegration(d, db)
		r.Nil(err)
	})
}

func TestStorageIntegrationUpdate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                   "test_storage_integration_acl",
		"storage_aws_object_acl": "bucket-owner-full-control",
	}

	d := storageIntegration(t, "test_storage_integration_acl", in)

	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on",
	},
	).AddRow("test_storage_integration_acl", "EXTERNAL_STAGE", "STORAGE", true, "now")

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).AddRow("ENABLED", "Boolean", true, false).
		AddRow("STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", nil)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^ALTER STORAGE INTEGRATION "test_storage_integration_acl" SET STORAGE_AWS_OBJECT_ACL = 'bucket-owner-full-control'`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER STORAGE INTEGRATION "test_storage_integration_acl" SET ENABLED=true`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS LIKE 'test_storage_integration_acl'$`).WillReturnRows(showRows)
		mock.ExpectQuery(`DESCRIBE STORAGE INTEGRATION "test_storage_integration_acl"$`).WillReturnRows(descRows)

		err := resources.UpdateStorageIntegration(d, db)
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
		"name", "type", "category", "enabled", "created_on",
	},
	).AddRow("test_storage_integration", "EXTERNAL_STAGE", "STORAGE", true, "now")
	mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS LIKE 'test_storage_integration'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).
		AddRow("ENABLED", "Boolean", true, false).
		AddRow("STORAGE_PROVIDER", "String", "S3", nil).
		AddRow("STORAGE_ALLOWED_LOCATIONS", "List", "s3://bucket-a/path-a/,s3://bucket-b/", nil).
		AddRow("STORAGE_BLOCKED_LOCATIONS", "List", "s3://bucket-c/path-c/,s3://bucket-d/", nil).
		AddRow("STORAGE_AWS_IAM_USER_ARN", "String", "arn:aws:iam::000000000000:/user/test", nil).
		AddRow("STORAGE_AWS_ROLE_ARN", "String", "arn:aws:iam::000000000001:/role/test", nil).
		AddRow("STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", nil).
		AddRow("STORAGE_AWS_EXTERNAL_ID", "String", "AGreatExternalID", nil)

	mock.ExpectQuery(`DESCRIBE STORAGE INTEGRATION "test_storage_integration"$`).WillReturnRows(descRows)
}

func expectReadStorageIntegrationWithS3GOV(mock sqlmock.Sqlmock) {
	showRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on",
	},
	).AddRow("test_storage_integration_with_s3gov", "EXTERNAL_STAGE", "STORAGE", true, "now")
	mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS LIKE 'test_storage_integration_with_s3gov'$`).WillReturnRows(showRows)

	descRows := sqlmock.NewRows([]string{
		"property", "property_type", "property_value", "property_default",
	}).
		AddRow("ENABLED", "Boolean", true, false).
		AddRow("STORAGE_PROVIDER", "String", "S3GOV", nil).
		AddRow("STORAGE_ALLOWED_LOCATIONS", "List", "s3://bucket-a/path-a/,s3://bucket-b/", nil).
		AddRow("STORAGE_BLOCKED_LOCATIONS", "List", "s3://bucket-c/path-c/,s3://bucket-d/", nil).
		AddRow("STORAGE_AWS_IAM_USER_ARN", "String", "arn:aws:iam::000000000000:/user/test", nil).
		AddRow("STORAGE_AWS_ROLE_ARN", "String", "arn:aws:iam::000000000001:/role/test", nil).
		AddRow("STORAGE_AWS_OBJECT_ACL", "String", "bucket-owner-full-control", nil).
		AddRow("STORAGE_AWS_EXTERNAL_ID", "String", "AGreatExternalID", nil)

	mock.ExpectQuery(`DESCRIBE STORAGE INTEGRATION "test_storage_integration_with_s3gov"$`).WillReturnRows(descRows)
}

func expectReadStorageIntegrationEmpty(mock sqlmock.Sqlmock) {
	noRows := sqlmock.NewRows([]string{
		"name", "type", "category", "enabled", "created_on",
	},
	)
	mock.ExpectQuery(`^SHOW STORAGE INTEGRATIONS.*`).WillReturnRows(noRows)
}
