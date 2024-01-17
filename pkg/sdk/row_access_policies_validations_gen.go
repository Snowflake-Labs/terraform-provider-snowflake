package sdk

var (
	_ validatable = new(CreateRowAccessPolicyOptions)
	_ validatable = new(AlterRowAccessPolicyOptions)
	_ validatable = new(DropRowAccessPolicyOptions)
	_ validatable = new(ShowRowAccessPolicyOptions)
	_ validatable = new(DescribeRowAccessPolicyOptions)
)

func (opts *CreateRowAccessPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.args) {
		errs = append(errs, errNotSet("CreateRowAccessPolicyOptions", "args"))
	}
	if !valueSet(opts.body) {
		errs = append(errs, errNotSet("CreateRowAccessPolicyOptions", "body"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateRowAccessPolicyOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterRowAccessPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.SetBody, opts.SetTags, opts.UnsetTags, opts.SetComment, opts.UnsetComment) {
		errs = append(errs, errExactlyOneOf("AlterRowAccessPolicyOptions", "RenameTo", "SetBody", "SetTags", "UnsetTags", "SetComment", "UnsetComment"))
	}
	return JoinErrors(errs...)
}

func (opts *DropRowAccessPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowRowAccessPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeRowAccessPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
