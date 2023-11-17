package sdk

var (
	_ validatable = new(CreateInternalStageOptions)
	_ validatable = new(CreateOnS3StageOptions)
	_ validatable = new(CreateOnGCSStageOptions)
	_ validatable = new(CreateOnAzureStageOptions)
	_ validatable = new(CreateOnS3CompatibleStageOptions)
	_ validatable = new(AlterStageOptions)
	_ validatable = new(AlterInternalStageStageOptions)
	_ validatable = new(AlterExternalS3StageStageOptions)
	_ validatable = new(AlterExternalGCSStageStageOptions)
	_ validatable = new(AlterExternalAzureStageStageOptions)
	_ validatable = new(AlterDirectoryTableStageOptions)
	_ validatable = new(DropStageOptions)
	_ validatable = new(DescribeStageOptions)
	_ validatable = new(ShowStageOptions)
)

func (opts *CreateInternalStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateInternalStageOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnS3StageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOnS3StageOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.ExternalStageParams) {
		if everyValueSet(opts.ExternalStageParams.StorageIntegration, opts.ExternalStageParams.Credentials) {
			errs = append(errs, errOneOf("CreateOnS3StageOptions.ExternalStageParams", "StorageIntegration", "Credentials"))
		}
		if valueSet(opts.ExternalStageParams.Credentials) {
			if everyValueSet(opts.ExternalStageParams.Credentials.AWSKeyId, opts.ExternalStageParams.Credentials.AWSRole) {
				errs = append(errs, errOneOf("AlterExternalS3StageStageOptions.ExternalStageParams.Credentials", "AWSKeyId", "AWSRole"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnGCSStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOnGCSStageOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnAzureStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOnAzureStageOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.ExternalStageParams) {
		if everyValueSet(opts.ExternalStageParams.StorageIntegration, opts.ExternalStageParams.Credentials) {
			errs = append(errs, errOneOf("CreateOnAzureStageOptions.ExternalStageParams", "StorageIntegration", "Credentials"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateOnS3CompatibleStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOnS3CompatibleStageOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterStageOptions", "RenameTo", "SetTags", "UnsetTags"))
	}
	if everyValueSet(opts.IfExists, opts.UnsetTags) {
		errs = append(errs, errOneOf("AlterStageOptions", "IfExists", "UnsetTags"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterInternalStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterExternalS3StageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.ExternalStageParams) {
		if everyValueSet(opts.ExternalStageParams.StorageIntegration, opts.ExternalStageParams.Credentials) {
			errs = append(errs, errOneOf("AlterExternalS3StageStageOptions.ExternalStageParams", "StorageIntegration", "Credentials"))
		}
		if valueSet(opts.ExternalStageParams.Credentials) {
			if everyValueSet(opts.ExternalStageParams.Credentials.AWSKeyId, opts.ExternalStageParams.Credentials.AWSRole) {
				errs = append(errs, errOneOf("AlterExternalS3StageStageOptions.ExternalStageParams.Credentials", "AWSKeyId", "AWSRole"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterExternalGCSStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterExternalAzureStageStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.ExternalStageParams) {
		if everyValueSet(opts.ExternalStageParams.StorageIntegration, opts.ExternalStageParams.Credentials) {
			errs = append(errs, errOneOf("AlterExternalAzureStageStageOptions.ExternalStageParams", "StorageIntegration", "Credentials"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterDirectoryTableStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.SetDirectory, opts.Refresh) {
		errs = append(errs, errOneOf("AlterDirectoryTableStageOptions", "SetDirectory", "Refresh"))
	}
	return JoinErrors(errs...)
}

func (opts *DropStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowStageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
