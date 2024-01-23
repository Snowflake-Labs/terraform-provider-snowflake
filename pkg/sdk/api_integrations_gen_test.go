package sdk

import "testing"

const AwsAllowedPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
const AzureAllowedPrefix = "https://apim-hello-world.azure-api.net/"
const GoogleAllowedPrefix = "https://gateway-id-123456.uc.gateway.dev/"

const ApiAwsRoleArn = "arn:aws:iam::000000000001:/role/test"
const AzureTenantId = "00000000-0000-0000-0000-000000000000"
const AzureAdApplicationId = "11111111-1111-1111-1111-111111111111"
const GoogleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"

func TestApiIntegrations_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateApiIntegrationOptions for AWS
	defaultOptsAws := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			AwsApiProviderParams: &AwsApiParams{
				ApiProvider:   ApiIntegrationAwsApiGateway,
				ApiAwsRoleArn: ApiAwsRoleArn,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: AwsAllowedPrefix}},
			Enabled:            true,
		}
	}

	// Minimal valid CreateApiIntegrationOptions for Azure
	defaultOptsAzure := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			AzureApiProviderParams: &AzureApiParams{
				AzureTenantId:        AzureTenantId,
				AzureAdApplicationId: AzureAdApplicationId,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: AzureAllowedPrefix}},
			Enabled:            true,
		}
	}

	// Minimal valid CreateApiIntegrationOptions for Google
	defaultOptsGoogle := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			GoogleApiProviderParams: &GoogleApiParams{
				GoogleAudience: GoogleAudience,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: GoogleAllowedPrefix}},
			Enabled:            true,
		}
	}

	defaultOpts := defaultOptsAws

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApiIntegrationOptions = nil
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApiIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.AwsApiProviderParams opts.AzureApiProviderParams opts.GoogleApiProviderParams] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.AwsApiProviderParams = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateApiIntegrationOptions", "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams"))
	})

	t.Run("validation: exactly one field from [opts.AwsApiProviderParams opts.AzureApiProviderParams opts.GoogleApiProviderParams] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.AzureApiProviderParams = new(AzureApiParams)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateApiIntegrationOptions", "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION %s API_PROVIDER = aws_api_gateway API_AWS_ROLE_ARN = '%s' API_ALLOWED_PREFIXES = ('%s') ENABLED = true`, id.FullyQualifiedName(), ApiAwsRoleArn, AwsAllowedPrefix)
	})

	t.Run("all options - aws", func(t *testing.T) {
		opts := defaultOptsAws()
		opts.IfNotExists = Bool(true)
		opts.AwsApiProviderParams.ApiProvider = ApiIntegrationAwsPrivateApiGateway
		opts.AwsApiProviderParams.ApiKey = String("key")
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: GoogleAllowedPrefix}, {Path: AzureAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = aws_private_api_gateway API_AWS_ROLE_ARN = '%s' API_KEY = 'key' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), ApiAwsRoleArn, AwsAllowedPrefix, GoogleAllowedPrefix, AzureAllowedPrefix)
	})

	t.Run("all options - azure", func(t *testing.T) {
		opts := defaultOptsAzure()
		opts.IfNotExists = Bool(true)
		opts.AzureApiProviderParams.ApiKey = String("key")
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: AwsAllowedPrefix}, {Path: GoogleAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = azure_api_management AZURE_TENANT_ID = '%s' AZURE_AD_APPLICATION_ID = '%s' API_KEY = 'key' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), AzureTenantId, AzureAdApplicationId, AzureAllowedPrefix, AwsAllowedPrefix, GoogleAllowedPrefix)
	})

	t.Run("all options - google", func(t *testing.T) {
		opts := defaultOptsGoogle()
		opts.IfNotExists = Bool(true)
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: AwsAllowedPrefix}, {Path: AzureAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = google_api_gateway GOOGLE_AUDIENCE = '%s' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), GoogleAudience, GoogleAllowedPrefix, AwsAllowedPrefix, AzureAllowedPrefix)
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
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP API INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestApiIntegrations_Show(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid ShowApiIntegrationOptions
	defaultOpts := func() *ShowApiIntegrationOptions {
		return &ShowApiIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApiIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW API INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW API INTEGRATIONS LIKE '%s'", id.Name())
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
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE API INTEGRATION %s", id.FullyQualifiedName())
	})
}
