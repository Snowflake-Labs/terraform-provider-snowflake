package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Functions = (*functions)(nil)

type functions struct {
	client *Client
}

func (v *functions) CreateFunctionForJava(ctx context.Context, request *CreateFunctionForJavaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateFunctionForJavascript(ctx context.Context, request *CreateFunctionForJavascriptFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateFunctionForPython(ctx context.Context, request *CreateFunctionForPythonFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateFunctionForScala(ctx context.Context, request *CreateFunctionForScalaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateFunctionForSQL(ctx context.Context, request *CreateFunctionForSQLFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) Alter(ctx context.Context, request *AlterFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) Drop(ctx context.Context, request *DropFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) Show(ctx context.Context, request *ShowFunctionRequest) ([]Function, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[functionRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[functionRow, Function](dbRows)
	return resultList, nil
}

func (v *functions) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Function, error) {
	request := NewShowFunctionRequest().WithLike(id.Name())
	functions, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(functions, func(r Function) bool { return r.Name == id.Name() })
}

func (v *functions) Describe(ctx context.Context, request *DescribeFunctionRequest) ([]FunctionDetail, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[functionDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[functionDetailRow, FunctionDetail](rows), nil
}

func (r *CreateFunctionForJavaFunctionRequest) toOpts() *CreateFunctionForJavaFunctionOptions {
	opts := &CreateFunctionForJavaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,
		Comment:        r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,

		TargetPath: r.TargetPath,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
				Default:     v.Default,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &FunctionReturns{
			ResultDataType: r.Returns.ResultDataType,
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &FunctionReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]FunctionColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = FunctionColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.ReturnNullValues != nil {
		opts.ReturnNullValues = r.ReturnNullValues
	}
	if r.NullInputBehavior != nil {
		opts.NullInputBehavior = r.NullInputBehavior
	}
	if r.ReturnResultsBehavior != nil {
		opts.ReturnResultsBehavior = r.ReturnResultsBehavior
	}
	if r.Imports != nil {
		s := make([]FunctionImports, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImports{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackages, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackages{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Secrets != nil {
		s := make([]FunctionSecret, len(r.Secrets))
		for i, v := range r.Secrets {
			s[i] = FunctionSecret{
				SecretVariableName: v.SecretVariableName,
				SecretName:         v.SecretName,
			}
		}
		opts.Secrets = s
	}
	opts.FunctionDefinition = r.FunctionDefinition
	return opts
}

func (r *CreateFunctionForJavascriptFunctionRequest) toOpts() *CreateFunctionForJavascriptFunctionOptions {
	opts := &CreateFunctionForJavascriptFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		Comment: r.Comment,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
				Default:     v.Default,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &FunctionReturns{
			ResultDataType: r.Returns.ResultDataType,
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &FunctionReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]FunctionColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = FunctionColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.ReturnNullValues != nil {
		opts.ReturnNullValues = r.ReturnNullValues
	}
	if r.NullInputBehavior != nil {
		opts.NullInputBehavior = r.NullInputBehavior
	}
	if r.ReturnResultsBehavior != nil {
		opts.ReturnResultsBehavior = r.ReturnResultsBehavior
	}

	opts.FunctionDefinition = r.FunctionDefinition
	return opts
}

func (r *CreateFunctionForPythonFunctionRequest) toOpts() *CreateFunctionForPythonFunctionOptions {
	opts := &CreateFunctionForPythonFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,
		Comment:        r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
				Default:     v.Default,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &FunctionReturns{
			ResultDataType: r.Returns.ResultDataType,
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &FunctionReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]FunctionColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = FunctionColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.ReturnNullValues != nil {
		opts.ReturnNullValues = r.ReturnNullValues
	}
	if r.NullInputBehavior != nil {
		opts.NullInputBehavior = r.NullInputBehavior
	}
	if r.ReturnResultsBehavior != nil {
		opts.ReturnResultsBehavior = r.ReturnResultsBehavior
	}
	if r.Imports != nil {
		s := make([]FunctionImports, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImports{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackages, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackages{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Secrets != nil {
		s := make([]FunctionSecret, len(r.Secrets))
		for i, v := range r.Secrets {
			s[i] = FunctionSecret{
				SecretVariableName: v.SecretVariableName,
				SecretName:         v.SecretName,
			}
		}
		opts.Secrets = s
	}
	opts.FunctionDefinition = r.FunctionDefinition
	return opts
}

func (r *CreateFunctionForScalaFunctionRequest) toOpts() *CreateFunctionForScalaFunctionOptions {
	opts := &CreateFunctionForScalaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,
		Comment:        r.Comment,

		Handler:    r.Handler,
		TargetPath: r.TargetPath,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
				Default:     v.Default,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &FunctionReturns{
			ResultDataType: r.Returns.ResultDataType,
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &FunctionReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]FunctionColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = FunctionColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.ReturnNullValues != nil {
		opts.ReturnNullValues = r.ReturnNullValues
	}
	if r.NullInputBehavior != nil {
		opts.NullInputBehavior = r.NullInputBehavior
	}
	if r.ReturnResultsBehavior != nil {
		opts.ReturnResultsBehavior = r.ReturnResultsBehavior
	}
	if r.Imports != nil {
		s := make([]FunctionImports, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImports{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackages, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackages{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	opts.FunctionDefinition = r.FunctionDefinition
	return opts
}

func (r *CreateFunctionForSQLFunctionRequest) toOpts() *CreateFunctionForSQLFunctionOptions {
	opts := &CreateFunctionForSQLFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		Memoizable: r.Memoizable,
		Comment:    r.Comment,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
				Default:     v.Default,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &FunctionReturns{
			ResultDataType: r.Returns.ResultDataType,
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &FunctionReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]FunctionColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = FunctionColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.ReturnNullValues != nil {
		opts.ReturnNullValues = r.ReturnNullValues
	}
	if r.ReturnResultsBehavior != nil {
		opts.ReturnResultsBehavior = r.ReturnResultsBehavior
	}
	opts.FunctionDefinition = r.FunctionDefinition
	return opts
}

func (r *AlterFunctionRequest) toOpts() *AlterFunctionOptions {
	opts := &AlterFunctionOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo:  r.RenameTo,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.ArgumentTypes != nil {
		s := make([]FunctionArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = FunctionArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	if r.Set != nil {
		opts.Set = &FunctionSet{
			LogLevel:   r.Set.LogLevel,
			TraceLevel: r.Set.TraceLevel,
			Comment:    r.Set.Comment,
			Secure:     r.Set.Secure,
		}
	}
	if r.Unset != nil {
		opts.Unset = &FunctionUnset{
			Secure:     r.Unset.Secure,
			Comment:    r.Unset.Comment,
			LogLevel:   r.Unset.LogLevel,
			TraceLevel: r.Unset.TraceLevel,
		}
	}
	return opts
}

func (r *DropFunctionRequest) toOpts() *DropFunctionOptions {
	opts := &DropFunctionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.ArgumentTypes != nil {
		s := make([]FunctionArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = FunctionArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	return opts
}

func (r *ShowFunctionRequest) toOpts() *ShowFunctionOptions {
	opts := &ShowFunctionOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r functionRow) convert() *Function {
	return &Function{
		CreatedOn:          r.CreatedOn,
		Name:               r.Name,
		SchemaName:         r.SchemaName,
		IsBuiltIn:          r.IsBuiltIn == "Y",
		IsAggregate:        r.IsAggregate == "Y",
		IsAnsi:             r.IsAnsi == "Y",
		MinNumArguments:    r.MinNumArguments,
		MaxNumArguments:    r.MaxNumArguments,
		Arguments:          r.Arguments,
		Description:        r.Description,
		CatalogName:        r.CatalogName,
		IsTableFunction:    r.IsTableFunction == "Y",
		ValidForClustering: r.ValidForClustering == "Y",
		IsSecure:           r.IsSecure == "Y",
		IsExternalFunction: r.IsExternalFunction == "Y",
		Language:           r.Language,
		IsMemoizable:       r.IsMemoizable == "Y",
	}
}

func (r *DescribeFunctionRequest) toOpts() *DescribeFunctionOptions {
	opts := &DescribeFunctionOptions{
		name: r.name,
	}
	if r.ArgumentTypes != nil {
		s := make([]FunctionArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = FunctionArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	return opts
}

func (r functionDetailRow) convert() *FunctionDetail {
	return &FunctionDetail{
		Property: r.Property,
		Value:    r.Value,
	}
}
