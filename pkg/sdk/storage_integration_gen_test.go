package sdk

import "testing"

func TestStorageIntegrations_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateStorageIntegrationOptions
	defaultOpts := func() *CreateStorageIntegrationOptions {
		return &CreateStorageIntegrationOptions{
			name: id,
			S3StorageProviderParams: &S3StorageParams{
				StorageAwsRoleArn: "arn:aws:iam::001234567890:role/role",
			},
			Enabled:                 true,
			StorageAllowedLocations: []StorageLocation{{Path: "allowed-loc-1"}, {Path: "allowed-loc-2"}},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateStorageIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateStorageIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.S3StorageProviderParams opts.GCSStorageProviderParams opts.AzureStorageProviderParams] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.S3StorageProviderParams = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateStorageIntegrationOptions", "S3StorageProviderParams", "GCSStorageProviderParams", "AzureStorageProviderParams"))
	})

	t.Run("validation: exactly one field from [opts.S3StorageProviderParams opts.GCSStorageProviderParams opts.AzureStorageProviderParams] should be present - two set", func(t *testing.T) {
		opts := defaultOpts()
		opts.GCSStorageProviderParams = new(GCSStorageParams)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateStorageIntegrationOptions", "S3StorageProviderParams", "GCSStorageProviderParams", "AzureStorageProviderParams"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2')`, id.FullyQualifiedName())
	})

	t.Run("all options - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.S3StorageProviderParams = &S3StorageParams{
			StorageAwsRoleArn:   "arn:aws:iam::001234567890:role/role",
			StorageAwsObjectAcl: String("bucket-owner-full-control"),
		}
		opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STORAGE INTEGRATION IF NOT EXISTS %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' STORAGE_AWS_OBJECT_ACL = 'bucket-owner-full-control' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("all options - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.S3StorageProviderParams = nil
		opts.GCSStorageProviderParams = new(GCSStorageParams)
		opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'GCS' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("all options - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.S3StorageProviderParams = nil
		opts.AzureStorageProviderParams = &AzureStorageParams{
			AzureTenantId: String("azure-tenant-id"),
		}
		opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'azure-tenant-id' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestStorageIntegrations_Alter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid AlterStorageIntegrationOptions
	defaultOpts := func() *AlterStorageIntegrationOptions {
		return &AlterStorageIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterStorageIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("one"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterStorageIntegrationOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStorageIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present - two set", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
		}
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("one"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStorageIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("set - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &StorageIntegrationSet{
			SetS3Params: &SetS3StorageParams{
				StorageAwsRoleArn:   "new-aws-role-arn",
				StorageAwsObjectAcl: String("new-aws-object-acl"),
			},
			Enabled:                 false,
			StorageAllowedLocations: []StorageLocation{{Path: "new-allowed-location"}},
			StorageBlockedLocations: []StorageLocation{{Path: "new-blocked-location"}},
			Comment:                 String("changed comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER STORAGE INTEGRATION %s SET STORAGE_AWS_ROLE_ARN = 'new-aws-role-arn' STORAGE_AWS_OBJECT_ACL = 'new-aws-object-acl' ENABLED = false STORAGE_ALLOWED_LOCATIONS = ('new-allowed-location') STORAGE_BLOCKED_LOCATIONS = ('new-blocked-location') COMMENT = 'changed comment'", id.FullyQualifiedName())
	})

	t.Run("set - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &StorageIntegrationSet{
			SetAzureParams: &SetAzureStorageParams{
				AzureTenantId: String("new-azure-tenant-id"),
			},
			Enabled:                 false,
			StorageAllowedLocations: []StorageLocation{{Path: "new-allowed-location"}},
			StorageBlockedLocations: []StorageLocation{{Path: "new-blocked-location"}},
			Comment:                 String("changed comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER STORAGE INTEGRATION %s SET AZURE_TENANT_ID = 'new-azure-tenant-id' ENABLED = false STORAGE_ALLOWED_LOCATIONS = ('new-allowed-location') STORAGE_BLOCKED_LOCATIONS = ('new-blocked-location') COMMENT = 'changed comment'", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STORAGE INTEGRATION IF EXISTS %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &StorageIntegrationUnset{
			Enabled:                 Bool(true),
			StorageBlockedLocations: Bool(true),
			Comment:                 Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER STORAGE INTEGRATION %s UNSET ENABLED, STORAGE_BLOCKED_LOCATIONS, COMMENT", id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STORAGE INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}
