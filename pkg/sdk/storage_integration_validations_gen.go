package sdk

var (
	_ validatable = new(CreateStorageIntegrationOptions)
	_ validatable = new(AlterStorageIntegrationOptions)
	_ validatable = new(DropStorageIntegrationOptions)
	_ validatable = new(ShowStorageIntegrationOptions)
	_ validatable = new(DescribeStorageIntegrationOptions)
)

func (opts *CreateStorageIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateStorageIntegrationOptions", "IfNotExists", "OrReplace"))
	}
	if !exactlyOneValueSet(opts.S3StorageProviderParams, opts.GCSStorageProviderParams, opts.AzureStorageProviderParams) {
		errs = append(errs, errExactlyOneOf("CreateStorageIntegrationOptions", "S3StorageProviderParams", "GCSStorageProviderParams", "AzureStorageProviderParams"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterStorageIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfExists, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterStorageIntegrationOptions", "IfExists", "UnsetTags"))
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterStorageIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) && everyValueNil(opts.Set.S3Params, opts.Set.AzureParams, opts.Set.Enabled, opts.Set.StorageAllowedLocations, opts.Set.StorageBlockedLocations, opts.Set.Comment) {
		errs = append(errs, errAtLeastOneOf("AlterStorageIntegrationOptions.Set", "S3Params", "AzureParams", "Enabled", "StorageAllowedLocations", "StorageBlockedLocations", "Comment"))
	}
	if valueSet(opts.Unset) && everyValueNil(opts.Unset.StorageBlockedLocations, opts.Unset.StorageAwsObjectAcl, opts.Unset.Enabled, opts.Unset.Comment) {
		errs = append(errs, errAtLeastOneOf("AlterStorageIntegrationOptions.Unset", "StorageBlockedLocations", "StorageAwsObjectAcl", "Enabled", "Comment"))
	}
	return JoinErrors(errs...)
}

func (opts *DropStorageIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowStorageIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeStorageIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
