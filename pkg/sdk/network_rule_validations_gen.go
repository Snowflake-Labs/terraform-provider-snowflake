package sdk

var (
	_ validatable = new(CreateNetworkRuleOptions)
	_ validatable = new(AlterNetworkRuleOptions)
	_ validatable = new(DropNetworkRuleOptions)
	_ validatable = new(ShowNetworkRuleOptions)
	_ validatable = new(DescribeNetworkRuleOptions)
)

func (opts *CreateNetworkRuleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterNetworkRuleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !anyValueSet(opts.Set, opts.Unset) {
		errs = append(errs, errAtLeastOneOf("AlterNetworkRuleOptions", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.ValueList, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterNetworkRuleOptions.Set", "ValueList", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.ValueList, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterNetworkRuleOptions.Unset", "ValueList", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropNetworkRuleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowNetworkRuleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeNetworkRuleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
