package sdk

var (
	_ validatable = new(CreateForJavaFunctionOptions)
	_ validatable = new(CreateForJavascriptFunctionOptions)
	_ validatable = new(CreateForPythonFunctionOptions)
	_ validatable = new(CreateForScalaFunctionOptions)
	_ validatable = new(CreateForSQLFunctionOptions)
	_ validatable = new(AlterFunctionOptions)
	_ validatable = new(DropFunctionOptions)
	_ validatable = new(ShowFunctionOptions)
	_ validatable = new(DescribeFunctionOptions)
)

func (opts *CreateForJavaFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForJavaFunctionOptions", "Handler"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateForJavaFunctionOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.Arguments) {
		// modified manually
		for _, arg := range opts.Arguments {
			if !exactlyOneValueSet(arg.ArgDataTypeOld, arg.ArgDataType) {
				errs = append(errs, errExactlyOneOf("CreateForJavaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
			}
		}
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForJavaFunctionOptions.Returns", "ResultDataType", "Table"))
		}
		if valueSet(opts.Returns.ResultDataType) {
			if !exactlyOneValueSet(opts.Returns.ResultDataType.ResultDataTypeOld, opts.Returns.ResultDataType.ResultDataType) {
				errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
			}
		}
		if valueSet(opts.Returns.Table) {
			if valueSet(opts.Returns.Table.Columns) {
				// modified manually
				for _, col := range opts.Returns.Table.Columns {
					if !exactlyOneValueSet(col.ColumnDataTypeOld, col.ColumnDataType) {
						errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
					}
				}
			}
		}
	}
	// added manually
	if opts.FunctionDefinition == nil {
		if opts.TargetPath != nil {
			errs = append(errs, NewError("TARGET_PATH must be nil when AS is nil"))
		}
		if len(opts.Imports) == 0 {
			errs = append(errs, NewError("IMPORTS must not be empty when AS is nil"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateForJavascriptFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.FunctionDefinition) {
		errs = append(errs, errNotSet("CreateForJavascriptFunctionOptions", "FunctionDefinition"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Arguments) {
		// modified manually
		for _, arg := range opts.Arguments {
			if !exactlyOneValueSet(arg.ArgDataTypeOld, arg.ArgDataType) {
				errs = append(errs, errExactlyOneOf("CreateForJavascriptFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
			}
		}
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForJavascriptFunctionOptions.Returns", "ResultDataType", "Table"))
		}
		if valueSet(opts.Returns.ResultDataType) {
			if !exactlyOneValueSet(opts.Returns.ResultDataType.ResultDataTypeOld, opts.Returns.ResultDataType.ResultDataType) {
				errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
			}
		}
		if valueSet(opts.Returns.Table) {
			if valueSet(opts.Returns.Table.Columns) {
				// modified manually
				for _, col := range opts.Returns.Table.Columns {
					if !exactlyOneValueSet(col.ColumnDataTypeOld, col.ColumnDataType) {
						errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
					}
				}
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateForPythonFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.RuntimeVersion) {
		errs = append(errs, errNotSet("CreateForPythonFunctionOptions", "RuntimeVersion"))
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForPythonFunctionOptions", "Handler"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateForPythonFunctionOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.Arguments) {
		// modified manually
		for _, arg := range opts.Arguments {
			if !exactlyOneValueSet(arg.ArgDataTypeOld, arg.ArgDataType) {
				errs = append(errs, errExactlyOneOf("CreateForPythonFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
			}
		}
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForPythonFunctionOptions.Returns", "ResultDataType", "Table"))
		}
		if valueSet(opts.Returns.ResultDataType) {
			if !exactlyOneValueSet(opts.Returns.ResultDataType.ResultDataTypeOld, opts.Returns.ResultDataType.ResultDataType) {
				errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
			}
		}
		if valueSet(opts.Returns.Table) {
			if valueSet(opts.Returns.Table.Columns) {
				// modified manually
				for _, col := range opts.Returns.Table.Columns {
					if !exactlyOneValueSet(col.ColumnDataTypeOld, col.ColumnDataType) {
						errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
					}
				}
			}
		}
	}
	// added manually
	if opts.FunctionDefinition == nil {
		if len(opts.Imports) == 0 {
			errs = append(errs, NewError("IMPORTS must not be empty when AS is nil"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateForScalaFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Handler) {
		errs = append(errs, errNotSet("CreateForScalaFunctionOptions", "Handler"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateForScalaFunctionOptions", "OrReplace", "IfNotExists"))
	}
	if !exactlyOneValueSet(opts.ResultDataTypeOld, opts.ResultDataType) {
		errs = append(errs, errExactlyOneOf("CreateForScalaFunctionOptions", "ResultDataTypeOld", "ResultDataType"))
	}
	if valueSet(opts.Arguments) {
		// modified manually
		for _, arg := range opts.Arguments {
			if !exactlyOneValueSet(arg.ArgDataTypeOld, arg.ArgDataType) {
				errs = append(errs, errExactlyOneOf("CreateForScalaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
			}
		}
	}
	// added manually
	if opts.FunctionDefinition == nil {
		if opts.TargetPath != nil {
			errs = append(errs, NewError("TARGET_PATH must be nil when AS is nil"))
		}
		if len(opts.Imports) == 0 {
			errs = append(errs, NewError("IMPORTS must not be empty when AS is nil"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *CreateForSQLFunctionOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !valueSet(opts.FunctionDefinition) {
		errs = append(errs, errNotSet("CreateForSQLFunctionOptions", "FunctionDefinition"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Arguments) {
		// modified manually
		for _, arg := range opts.Arguments {
			if !exactlyOneValueSet(arg.ArgDataTypeOld, arg.ArgDataType) {
				errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
			}
		}
	}
	if valueSet(opts.Returns) {
		if !exactlyOneValueSet(opts.Returns.ResultDataType, opts.Returns.Table) {
			errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns", "ResultDataType", "Table"))
		}
		if valueSet(opts.Returns.ResultDataType) {
			if !exactlyOneValueSet(opts.Returns.ResultDataType.ResultDataTypeOld, opts.Returns.ResultDataType.ResultDataType) {
				errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
			}
		}
		if valueSet(opts.Returns.Table) {
			if valueSet(opts.Returns.Table.Columns) {
				// modified manually
				for _, col := range opts.Returns.Table.Columns {
					if !exactlyOneValueSet(col.ColumnDataTypeOld, col.ColumnDataType) {
						errs = append(errs, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
					}
				}
			}
		}
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
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set, opts.Unset, opts.SetSecure, opts.UnsetSecure, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterFunctionOptions", "RenameTo", "Set", "Unset", "SetSecure", "UnsetSecure", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Comment, opts.Set.ExternalAccessIntegrations, opts.Set.SecretsList, opts.Set.EnableConsoleOutput, opts.Set.LogLevel, opts.Set.MetricLevel, opts.Set.TraceLevel) {
			errs = append(errs, errAtLeastOneOf("AlterFunctionOptions.Set", "Comment", "ExternalAccessIntegrations", "SecretsList", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Comment, opts.Unset.ExternalAccessIntegrations, opts.Unset.EnableConsoleOutput, opts.Unset.LogLevel, opts.Unset.MetricLevel, opts.Unset.TraceLevel) {
			errs = append(errs, errAtLeastOneOf("AlterFunctionOptions.Unset", "Comment", "ExternalAccessIntegrations", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"))
		}
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
