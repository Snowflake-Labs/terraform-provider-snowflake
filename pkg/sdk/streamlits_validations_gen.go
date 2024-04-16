package sdk

var (
	_ validatable = new(CreateStreamlitOptions)
	_ validatable = new(AlterStreamlitOptions)
	_ validatable = new(DropStreamlitOptions)
	_ validatable = new(ShowStreamlitOptions)
	_ validatable = new(DescribeStreamlitOptions)
)

func (opts *CreateStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.Warehouse != nil && !ValidObjectIdentifier(opts.Warehouse) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateStreamlitOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set) {
		errs = append(errs, errExactlyOneOf("AlterStreamlitOptions", "RenameTo", "Set"))
	}
	if valueSet(opts.Set) {
		if everyValueNil(opts.Set.RootLocation, opts.Set.MainFile, opts.Set.Warehouse, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterStreamlitOptions.Set", "RootLocation", "MainFile", "Warehouse", "Comment"))
		}

		if opts.Set.Warehouse != nil && !ValidObjectIdentifier(opts.Set.Warehouse) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeStreamlitOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
