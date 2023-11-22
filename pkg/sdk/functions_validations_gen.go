package sdk

var (
	_ validatable = new(CreateFunctionForJavaFunctionOptions)
	_ validatable = new(CreateFunctionForJavascriptFunctionOptions)
	_ validatable = new(CreateFunctionForPythonFunctionOptions)
	_ validatable = new(CreateFunctionForScalaFunctionOptions)
	_ validatable = new(CreateFunctionForSQLFunctionOptions)
	_ validatable = new(AlterFunctionOptions)
	_ validatable = new(DropFunctionOptions)
	_ validatable = new(ShowFunctionOptions)
	_ validatable = new(DescribeFunctionOptions)
)

func (v *FunctionReturns) validate() error {
	if v == nil {
		return ErrNilOptions
	}
	var errs []error
	if ok := exactlyOneValueSet(v.ResultDataType, v.Table); !ok {
		errs = append(errs, errOneOf("Returns.ResultDataType", "Returns.Table"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateFunctionForJavaFunctionOptions) validate() error {
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

func (opts *CreateFunctionForJavascriptFunctionOptions) validate() error {
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

func (opts *CreateFunctionForPythonFunctionOptions) validate() error {
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

func (opts *CreateFunctionForScalaFunctionOptions) validate() error {
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

func (opts *CreateFunctionForSQLFunctionOptions) validate() error {
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

func (opts *AlterFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DropFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
