package sdk

var (
	_ validatable = new(CreateForJavaProcedureOptions)
	_ validatable = new(CreateForJavaScriptProcedureOptions)
	_ validatable = new(CreateForPythonProcedureOptions)
	_ validatable = new(CreateForScalaProcedureOptions)
	_ validatable = new(CreateForSQLProcedureOptions)
	_ validatable = new(AlterProcedureOptions)
	_ validatable = new(DropProcedureOptions)
	_ validatable = new(ShowProcedureOptions)
	_ validatable = new(DescribeProcedureOptions)
)

func (opts *CreateForJavaProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.RuntimeVersion) {
		errs = append(errs, errNotSet("CreateForJavaProcedureOptions", "RuntimeVersion"))
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForJavaProcedureOptions", "Handler"))
	}
	if !valueSet(opts.Packages) {
		errs = append(errs, errNotSet("CreateForJavaProcedureOptions", "Packages"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForJavaProcedureOptions.Returns", "ResultDataType", "Table"))
		}
	}
	if opts.ProcedureDefinition == nil && opts.TargetPath != nil {
		errs = append(errs, NewError("TARGET_PATH must be nil when AS is nil"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateForJavaScriptProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.ProcedureDefinition) {
		errs = append(errs, errNotSet("CreateForJavaScriptProcedureOptions", "ProcedureDefinition"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *CreateForPythonProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.RuntimeVersion) {
		errs = append(errs, errNotSet("CreateForPythonProcedureOptions", "RuntimeVersion"))
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForPythonProcedureOptions", "Handler"))
	}
	if !valueSet(opts.Packages) {
		errs = append(errs, errNotSet("CreateForPythonProcedureOptions", "Packages"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForPythonProcedureOptions.Returns", "ResultDataType", "Table"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateForScalaProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.RuntimeVersion) {
		errs = append(errs, errNotSet("CreateForScalaProcedureOptions", "RuntimeVersion"))
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForScalaProcedureOptions", "Handler"))
	}
	if !valueSet(opts.Packages) {
		errs = append(errs, errNotSet("CreateForScalaProcedureOptions", "Packages"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForScalaProcedureOptions.Returns", "ResultDataType", "Table"))
		}
	}
	if opts.ProcedureDefinition == nil && opts.TargetPath != nil {
		errs = append(errs, NewError("TARGET_PATH must be nil when AS is nil"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateForSQLProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.ProcedureDefinition) {
		errs = append(errs, errNotSet("CreateForSQLProcedureOptions", "ProcedureDefinition"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForSQLProcedureOptions.Returns", "ResultDataType", "Table"))
		}
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
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.SetLogLevel, opts.SetTraceLevel, opts.UnsetComment, opts.SetTags, opts.UnsetTags, opts.ExecuteAs) {
		errs = append(errs, errExactlyOneOf("AlterProcedureOptions", "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "UnsetComment", "SetTags", "UnsetTags", "ExecuteAs"))
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

func (opts *ShowProcedureOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
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
