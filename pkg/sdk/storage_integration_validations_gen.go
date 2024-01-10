package sdk

var (
	_ validatable = new(CreateStorageIntegrationOptions)
	_ validatable = new(AlterStorageIntegrationOptions)
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
	return JoinErrors(errs...)
}
