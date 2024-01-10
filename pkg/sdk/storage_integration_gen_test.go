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
			StorageAllowedLocations: []string{"allowed-loc-1", "allowed-loc-2"},
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

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageBlockedLocations = []string{"blocked-loc-1", "blocked-loc-2"}
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterStorageIntegrationOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStorageIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
