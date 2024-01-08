package sdk

var (
	_ validatable = new(CreateApplicationPackageOptions)
	_ validatable = new(AlterApplicationPackageOptions)
	_ validatable = new(DropApplicationPackageOptions)
	_ validatable = new(ShowApplicationPackageOptions)
)

func (opts *CreateApplicationPackageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterApplicationPackageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.UnsetDataRetentionTimeInDays, opts.UnsetMaxDataExtensionTimeInDays, opts.UnsetDefaultDdlCollation, opts.UnsetComment, opts.UnsetDistribution, opts.ModifyReleaseDirective, opts.SetDefaultReleaseDirective, opts.SetReleaseDirective, opts.UnsetReleaseDirective, opts.AddVersion, opts.DropVersion, opts.AddPatchForVersion, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApplicationPackageOptions", "Set", "UnsetDataRetentionTimeInDays", "UnsetMaxDataExtensionTimeInDays", "UnsetDefaultDdlCollation", "UnsetComment", "UnsetDistribution", "ModifyReleaseDirective", "SetDefaultReleaseDirective", "SetReleaseDirective", "UnsetReleaseDirective", "AddVersion", "DropVersion", "AddPatchForVersion", "SetTags", "UnsetTags"))
	}
	return JoinErrors(errs...)
}

func (opts *DropApplicationPackageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowApplicationPackageOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
