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
	if valueSet(opts.DebugMode) && !valueSet(opts.Version) {
		errs = append(errs, NewError("CreateApplicationOptions.DebugMode can be set only when CreateApplicationOptions.Version is set"))
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
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.Upgrade, opts.UpgradeVersion, opts.UnsetReferences, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApplicationOptions", "Set", "Unset", "Upgrade", "UpgradeVersion", "UnsetReferences", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Unset) {
		if !exactlyOneValueSet(opts.Unset.Comment, opts.Unset.ShareEventsWithProvider, opts.Unset.DebugMode) {
			errs = append(errs, errExactlyOneOf("AlterApplicationOptions.Unset", "Comment", "ShareEventsWithProvider", "DebugMode"))
		}
	}
	if valueSet(opts.UpgradeVersion) {
		if !exactlyOneValueSet(opts.UpgradeVersion.VersionDirectory, opts.UpgradeVersion.VersionAndPatch) {
			errs = append(errs, errExactlyOneOf("AlterApplicationOptions.UpgradeVersion", "VersionDirectory", "VersionAndPatch"))
		}
	}
	if valueSet(opts.IfExists) {
		if !valueSet(opts.Set) && !valueSet(opts.Unset) {
			errs = append(errs, NewError("AlterApplicationOptions.IfExists can be set only when AlterApplicationOptions.Set or AlterApplicationOptions.Unset is set"))
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
