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
	if !exactlyOneValueSet(opts.S3ApiProviderParams, opts.AzureApiProviderParams, opts.GCSApiProviderParams) {
		errs = append(errs, errExactlyOneOf("CreateApiIntegrationOptions", "S3ApiProviderParams", "AzureApiProviderParams", "GCSApiProviderParams"))
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
		if everyValueSet(opts.Set.S3Params, opts.Set.AzureParams) {
			errs = append(errs, errOneOf("AlterApiIntegrationOptions.Set", "S3Params", "AzureParams"))
		}
		if !anyValueSet(opts.Set.S3Params, opts.Set.AzureParams, opts.Set.Enabled, opts.Set.ApiAllowedPrefixes, opts.Set.ApiBlockedPrefixes, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set", "S3Params", "AzureParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"))
		}
		if valueSet(opts.Set.S3Params) {
			if !anyValueSet(opts.Set.S3Params.ApiAwsRoleArn, opts.Set.S3Params.ApiKey) {
				errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set.S3Params", "ApiAwsRoleArn", "ApiKey"))
			}
		}
		if valueSet(opts.Set.AzureParams) {
			if !anyValueSet(opts.Set.AzureParams.AzureAdApplicationId, opts.Set.AzureParams.ApiKey) {
				errs = append(errs, errAtLeastOneOf("AlterApiIntegrationOptions.Set.AzureParams", "AzureAdApplicationId", "ApiKey"))
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
