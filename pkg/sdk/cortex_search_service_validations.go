package sdk

var (
	_ validatable = new(createCortexSearchServiceOptions)
	_ validatable = new(alterCortexSearchServiceOptions)
	_ validatable = new(dropCortexSearchServiceOptions)
	_ validatable = new(showCortexSearchServiceOptions)
	_ validatable = new(describeCortexSearchServiceOptions)
	_ validatable = new(CortexSearchServiceSet)
)

func (opts *createCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.warehouse) {
		errs = append(errs, errInvalidIdentifier("createCortexSearchServiceOptions", "warehouse"))
	}
	return JoinErrors(errs...)
}

func (dts *CortexSearchServiceSet) validate() error {
	var errs []error
	if dts.Warehouse != nil && !ValidObjectIdentifier(*dts.Warehouse) {
		errs = append(errs, errInvalidIdentifier("CortexSearchServiceSet", "Warehouse"))
	}
	return JoinErrors(errs...)
}

func (opts *alterCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Set); !ok {
		errs = append(errs, errExactlyOneOf("alterCortexSearchServiceOptions", "Set"))
	}
	return JoinErrors(errs...)
}

func (opts *showCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("showCortexSearchServiceOptions.In", "Account", "Database", "Schema"))
	}
	return JoinErrors(errs...)
}

func (opts *dropCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !ValidObjectIdentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *describeCortexSearchServiceOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !ValidObjectIdentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}
