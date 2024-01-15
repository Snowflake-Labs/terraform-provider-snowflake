package sdk

var (
	_ validatable = new(CreateApplicationOptions)
	_ validatable = new(DropApplicationOptions)
	_ validatable = new(AlterApplicationOptions)
	_ validatable = new(ShowApplicationOptions)
	_ validatable = new(DescribeApplicationOptions)
)

func (opts *CreateApplicationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.PackageName) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Version) {
		if !exactlyOneValueSet(opts.Version.VersionDirectory, opts.Version.VersionAndPatch) {
			errs = append(errs, errExactlyOneOf("CreateApplicationOptions.Version", "VersionDirectory", "VersionAndPatch"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropApplicationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterApplicationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.UnsetComment, opts.UnsetShareEventsWithProvider, opts.UnsetDebugMode, opts.Upgrade, opts.UpgradeVersion, opts.UnsetReferences, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApplicationOptions", "Set", "UnsetComment", "UnsetShareEventsWithProvider", "UnsetDebugMode", "Upgrade", "UpgradeVersion", "UnsetReferences", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.UpgradeVersion) {
		if !exactlyOneValueSet(opts.UpgradeVersion.VersionDirectory, opts.UpgradeVersion.VersionAndPatch) {
			errs = append(errs, errExactlyOneOf("AlterApplicationOptions.UpgradeVersion", "VersionDirectory", "VersionAndPatch"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowApplicationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeApplicationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
