package sdk

var (
	_ validatable = new(CreateSequenceOptions)
	_ validatable = new(AlterSequenceOptions)
	_ validatable = new(ShowSequenceOptions)
	_ validatable = new(DescribeSequenceOptions)
	_ validatable = new(DropSequenceOptions)
)

func (opts *CreateSequenceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSequenceOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterSequenceOptions) validate() error {
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
	if !exactlyOneValueSet(opts.RenameTo, opts.SetIncrement, opts.Set, opts.UnsetComment) {
		errs = append(errs, errExactlyOneOf("AlterSequenceOptions", "RenameTo", "SetIncrement", "Set", "UnsetComment"))
	}
	return JoinErrors(errs...)
}

func (opts *ShowSequenceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeSequenceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DropSequenceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Constraint) {
		if !exactlyOneValueSet(opts.Constraint.Cascade, opts.Constraint.Restrict) {
			errs = append(errs, errExactlyOneOf("DropSequenceOptions.Constraint", "Cascade", "Restrict"))
		}
	}
	return JoinErrors(errs...)
}
