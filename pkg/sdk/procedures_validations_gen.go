package sdk

var (
	_ validatable = new(CreateProcedureForJavaProcedureOptions)
	_ validatable = new(CreateProcedureForJavaScriptProcedureOptions)
	_ validatable = new(CreateProcedureForPythonProcedureOptions)
	_ validatable = new(CreateProcedureForScalaProcedureOptions)
	_ validatable = new(CreateProcedureForSQLProcedureOptions)
)

func (v *ProcedureReturns) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := exactlyOneValueSet(v.ResultDataType, v.Table); !ok {
		errs = append(errs, errOneOf("Returns.ResultDataType", "Returns.Table"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateProcedureForJavaProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if err := opts.Returns.validate(); err != nil {
		errs = append(errs, err)
	}
	return JoinErrors(errs...)
}

func (opts *CreateProcedureForJavaScriptProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *CreateProcedureForPythonProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if err := opts.Returns.validate(); err != nil {
		errs = append(errs, err)
	}
	return JoinErrors(errs...)
}

func (opts *CreateProcedureForScalaProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if err := opts.Returns.validate(); err != nil {
		errs = append(errs, err)
	}
	return JoinErrors(errs...)
}

func (opts *CreateProcedureForSQLProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *AlterProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DropProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
