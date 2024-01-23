package sdk

import "testing"

func TestApiIntegrations_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateApiIntegrationOptions
	defaultOpts := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApiIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.S3ApiProviderParams opts.AzureApiProviderParams opts.GCSApiProviderParams] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateApiIntegrationOptions", "S3ApiProviderParams", "AzureApiProviderParams", "GCSApiProviderParams"))
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

func TestApiIntegrations_Alter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid AlterApiIntegrationOptions
	defaultOpts := func() *AlterApiIntegrationOptions {
		return &AlterApiIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.SetTags]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterApiIntegrationOptions", "IfExists", "SetTags"))
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterApiIntegrationOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: conflicting fields for [opts.Set.S3Params opts.Set.AzureParams]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterApiIntegrationOptions.Set", "S3Params", "AzureParams"))
	})

	t.Run("validation: at least one of the fields [opts.Set.S3Params opts.Set.AzureParams opts.Set.Enabled opts.Set.ApiAllowedPrefixes opts.Set.ApiBlockedPrefixes opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set", "S3Params", "AzureParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Set.S3Params.ApiAwsRoleArn opts.Set.S3Params.ApiKey] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set.S3Params", "ApiAwsRoleArn", "ApiKey"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AzureParams.AzureAdApplicationId opts.Set.AzureParams.ApiKey] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AzureParams", "AzureAdApplicationId", "ApiKey"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ApiKey opts.Unset.Enabled opts.Unset.ApiBlockedPrefixes opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Unset", "ApiKey", "Enabled", "ApiBlockedPrefixes", "Comment"))
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

func TestApiIntegrations_Drop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DropApiIntegrationOptions
	defaultOpts := func() *DropApiIntegrationOptions {
		return &DropApiIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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

func TestApiIntegrations_Show(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid ShowApiIntegrationOptions
	defaultOpts := func() *ShowApiIntegrationOptions {
		return &ShowApiIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
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

func TestApiIntegrations_Describe(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DescribeApiIntegrationOptions
	defaultOpts := func() *DescribeApiIntegrationOptions {
		return &DescribeApiIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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
