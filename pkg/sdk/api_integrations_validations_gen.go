package sdk

var (
	_ validatable = new(CreateApiIntegrationOptions)
	_ validatable = new(AlterApiIntegrationOptions)
	_ validatable = new(DropApiIntegrationOptions)
	_ validatable = new(ShowApiIntegrationOptions)
	_ validatable = new(DescribeApiIntegrationOptions)
)

func (opts *CreateApiIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateApiIntegrationOptions", "IfNotExists", "OrReplace"))
	}
	if !exactlyOneValueSet(opts.AwsApiProviderParams, opts.AzureApiProviderParams, opts.GoogleApiProviderParams) {
		errs = append(errs, errExactlyOneOf("CreateApiIntegrationOptions", "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterApiIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfExists, opts.SetTags) {
		errs = append(errs, errOneOf("AlterApiIntegrationOptions", "IfExists", "SetTags"))
	}
	if everyValueSet(opts.IfExists, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterApiIntegrationOptions", "IfExists", "UnsetTags"))
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApiIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if moreThanOneValueSet(opts.Set.AwsParams, opts.Set.AzureParams, opts.Set.GoogleParams) {
			errs = append(errs, errMoreThanOneOf("AlterApiIntegrationOptions.Set", "AwsParams", "AzureParams", "GoogleParams"))
		}
		if !anyValueSet(opts.Set.AwsParams, opts.Set.AzureParams, opts.Set.GoogleParams, opts.Set.Enabled, opts.Set.ApiAllowedPrefixes, opts.Set.ApiBlockedPrefixes, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set", "AwsParams", "AzureParams", "GoogleParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"))
		}
		if valueSet(opts.Set.AwsParams) {
			if !anyValueSet(opts.Set.AwsParams.ApiAwsRoleArn, opts.Set.AwsParams.ApiKey) {
				errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AwsParams", "ApiAwsRoleArn", "ApiKey"))
			}
		}
		if valueSet(opts.Set.AzureParams) {
			if !anyValueSet(opts.Set.AzureParams.AzureTenantId, opts.Set.AzureParams.AzureAdApplicationId, opts.Set.AzureParams.ApiKey) {
				errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AzureParams", "AzureTenantId", "AzureAdApplicationId", "ApiKey"))
			}
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.ApiKey, opts.Unset.Enabled, opts.Unset.ApiBlockedPrefixes, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Unset", "ApiKey", "Enabled", "ApiBlockedPrefixes", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropApiIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowApiIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeApiIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
