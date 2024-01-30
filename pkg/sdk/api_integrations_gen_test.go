package sdk

import "testing"

const (
	awsAllowedPrefix    = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	azureAllowedPrefix  = "https://apim-hello-world.azure-api.net/"
	googleAllowedPrefix = "https://gateway-id-123456.uc.gateway.dev/"

	apiAwsRoleArn        = "arn:aws:iam::000000000001:/role/test"
	azureTenantId        = "00000000-0000-0000-0000-000000000000"
	azureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	googleAudience       = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
)

func TestApiIntegrations_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateApiIntegrationOptions for AWS
	defaultOptsAws := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			AwsApiProviderParams: &AwsApiParams{
				ApiProvider:   ApiIntegrationAwsApiGateway,
				ApiAwsRoleArn: apiAwsRoleArn,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}},
			Enabled:            true,
		}
	}

	// Minimal valid CreateApiIntegrationOptions for Azure
	defaultOptsAzure := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			AzureApiProviderParams: &AzureApiParams{
				AzureTenantId:        azureTenantId,
				AzureAdApplicationId: azureAdApplicationId,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}},
			Enabled:            true,
		}
	}

	// Minimal valid CreateApiIntegrationOptions for Google
	defaultOptsGoogle := func() *CreateApiIntegrationOptions {
		return &CreateApiIntegrationOptions{
			name: id,
			GoogleApiProviderParams: &GoogleApiParams{
				GoogleAudience: googleAudience,
			},
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}},
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION %s API_PROVIDER = aws_api_gateway API_AWS_ROLE_ARN = '%s' API_ALLOWED_PREFIXES = ('%s') ENABLED = true`, id.FullyQualifiedName(), apiAwsRoleArn, awsAllowedPrefix)
	})

	t.Run("all options - aws", func(t *testing.T) {
		opts := defaultOptsAws()
		opts.IfNotExists = Bool(true)
		opts.AwsApiProviderParams.ApiProvider = ApiIntegrationAwsPrivateApiGateway
		opts.AwsApiProviderParams.ApiKey = String("key")
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}, {Path: azureAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = aws_private_api_gateway API_AWS_ROLE_ARN = '%s' API_KEY = 'key' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), apiAwsRoleArn, awsAllowedPrefix, googleAllowedPrefix, azureAllowedPrefix)
	})

	t.Run("all options - azure", func(t *testing.T) {
		opts := defaultOptsAzure()
		opts.IfNotExists = Bool(true)
		opts.AzureApiProviderParams.ApiKey = String("key")
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: googleAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = azure_api_management AZURE_TENANT_ID = '%s' AZURE_AD_APPLICATION_ID = '%s' API_KEY = 'key' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), azureTenantId, azureAdApplicationId, azureAllowedPrefix, awsAllowedPrefix, googleAllowedPrefix)
	})

	t.Run("all options - google", func(t *testing.T) {
		opts := defaultOptsGoogle()
		opts.IfNotExists = Bool(true)
		opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: azureAllowedPrefix}}
		opts.Enabled = false
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE API INTEGRATION IF NOT EXISTS %s API_PROVIDER = google_api_gateway GOOGLE_AUDIENCE = '%s' API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') ENABLED = false COMMENT = 'some comment'`, id.FullyQualifiedName(), googleAudience, googleAllowedPrefix, awsAllowedPrefix, azureAllowedPrefix)
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
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.SetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterApiIntegrationOptions", "IfExists", "SetTags"))
	})

	t.Run("validation: conflicting fields for [opts.IfExists opts.UnsetTags]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("one"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterApiIntegrationOptions", "IfExists", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			Enabled: Bool(true),
		}
		opts.Unset = &ApiIntegrationUnset{
			Enabled: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterApiIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: conflicting fields for [opts.Set.AwsParams opts.Set.AzureParams opts.Set.GoogleParams]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			AwsParams:   &SetAwsApiParams{ApiKey: String("key")},
			AzureParams: &SetAzureApiParams{ApiKey: String("key")},
		}
		assertOptsInvalidJoinedErrors(t, opts, errMoreThanOneOf("AlterApiIntegrationOptions.Set", "AwsParams", "AzureParams", "GoogleParams"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AwsParams opts.Set.AzureParams opts.Set.GoogleParams opts.Set.Enabled opts.Set.ApiAllowedPrefixes opts.Set.ApiBlockedPrefixes opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set", "AwsParams", "AzureParams", "GoogleParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AwsParams.ApiAwsRoleArn opts.Set.AwsParams.ApiKey] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			AwsParams: &SetAwsApiParams{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AwsParams", "ApiAwsRoleArn", "ApiKey"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AzureParams.AzureTenantId opts.Set.AzureParams.AzureAdApplicationId opts.Set.AzureParams.ApiKey] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			AzureParams: &SetAzureApiParams{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AzureParams", "AzureTenantId", "AzureAdApplicationId", "ApiKey"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ApiKey opts.Unset.Enabled opts.Unset.ApiBlockedPrefixes opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterApiIntegrationOptions.Unset", "ApiKey", "Enabled", "ApiBlockedPrefixes", "Comment"))
	})

	t.Run("set - aws", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			AwsParams: &SetAwsApiParams{
				ApiAwsRoleArn: String("new-aws-role-arn"),
				ApiKey:        String("key"),
			},
			Enabled:            Bool(true),
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}},
			ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}, {Path: googleAllowedPrefix}},
			Comment:            String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER API INTEGRATION %s SET API_AWS_ROLE_ARN = 'new-aws-role-arn' API_KEY = 'key' ENABLED = true API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') COMMENT = 'comment'", id.FullyQualifiedName(), awsAllowedPrefix, azureAllowedPrefix, googleAllowedPrefix)
	})

	t.Run("set - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			AzureParams: &SetAzureApiParams{
				AzureAdApplicationId: String("new-azure-ad-application-id"),
				ApiKey:               String("key"),
			},
			Enabled:            Bool(true),
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}},
			ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: googleAllowedPrefix}},
			Comment:            String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER API INTEGRATION %s SET AZURE_AD_APPLICATION_ID = 'new-azure-ad-application-id' API_KEY = 'key' ENABLED = true API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') COMMENT = 'comment'", id.FullyQualifiedName(), azureAllowedPrefix, awsAllowedPrefix, googleAllowedPrefix)
	})

	t.Run("set - google", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ApiIntegrationSet{
			Enabled:            Bool(true),
			ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}},
			ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: azureAllowedPrefix}},
			Comment:            String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER API INTEGRATION %s SET ENABLED = true API_ALLOWED_PREFIXES = ('%s') API_BLOCKED_PREFIXES = ('%s', '%s') COMMENT = 'comment'", id.FullyQualifiedName(), googleAllowedPrefix, awsAllowedPrefix, azureAllowedPrefix)
	})

	t.Run("unset single", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiIntegrationUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER API INTEGRATION %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("unset multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ApiIntegrationUnset{
			ApiKey:             Bool(true),
			Enabled:            Bool(true),
			ApiBlockedPrefixes: Bool(true),
			Comment:            Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER API INTEGRATION %s UNSET API_KEY, ENABLED, API_BLOCKED_PREFIXES, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER API INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER API INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
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
