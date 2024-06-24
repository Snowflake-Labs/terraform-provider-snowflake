package sdk

var (
	_ validatable = new(CreateCortexSearchServiceOptions)
	_ validatable = new(AlterCortexSearchServiceOptions)
	_ validatable = new(ShowCortexSearchServiceOptions)
	_ validatable = new(DescribeCortexSearchServiceOptions)
	_ validatable = new(DropCortexSearchServiceOptions)
)

func (opts *CreateCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.On) {
		errs = append(errs, errNotSet("CreateCortexSearchServiceOptions", "On"))
	}
	if !valueSet(opts.TargetLag) {
		errs = append(errs, errNotSet("CreateCortexSearchServiceOptions", "TargetLag"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateCortexSearchServiceOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set) {
		errs = append(errs, errExactlyOneOf("AlterCortexSearchServiceOptions", "Set"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.TargetLag, opts.Set.Warehouse, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterCortexSearchServiceOptions.Set", "TargetLag", "Warehouse", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DropCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
