package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Functions = (*functions)(nil)

type functions struct {
	client *Client
}

func (v *functions) CreateForJava(ctx context.Context, request *CreateForJavaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForJavascript(ctx context.Context, request *CreateForJavascriptFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForPython(ctx context.Context, request *CreateForPythonFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForScala(ctx context.Context, request *CreateForScalaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForSQL(ctx context.Context, request *CreateForSQLFunctionRequest) error {
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
	request := NewShowFunctionRequest().WithIn(&In{Database: NewAccountObjectIdentifier(id.DatabaseName())}).WithLike(&Like{String(id.Name())})
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

func (r *CreateForJavaFunctionRequest) toOpts() *CreateForJavaFunctionOptions {
	opts := &CreateForJavaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		FunctionDefinition:         r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:      v.ArgName,
				ArgDataType:  v.ArgDataType,
				DefaultValue: v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataType: r.Returns.ResultDataType.ResultDataType,
		}
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
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForJavascriptFunctionRequest) toOpts() *CreateForJavascriptFunctionOptions {
	opts := &CreateForJavascriptFunctionOptions{
		OrReplace: r.OrReplace,
		Temporary: r.Temporary,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Comment:               r.Comment,
		FunctionDefinition:    r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:      v.ArgName,
				ArgDataType:  v.ArgDataType,
				DefaultValue: v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataType: r.Returns.ResultDataType.ResultDataType,
		}
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
	return opts
}

func (r *CreateForPythonFunctionRequest) toOpts() *CreateForPythonFunctionOptions {
	opts := &CreateForPythonFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		FunctionDefinition:         r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:      v.ArgName,
				ArgDataType:  v.ArgDataType,
				DefaultValue: v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataType: r.Returns.ResultDataType.ResultDataType,
		}
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
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForScalaFunctionRequest) toOpts() *CreateForScalaFunctionOptions {
	opts := &CreateForScalaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants:            r.CopyGrants,
		ResultDataType:        r.ResultDataType,
		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:            r.Handler,
		TargetPath:         r.TargetPath,
		FunctionDefinition: r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:      v.ArgName,
				ArgDataType:  v.ArgDataType,
				DefaultValue: v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForSQLFunctionRequest) toOpts() *CreateForSQLFunctionOptions {
	opts := &CreateForSQLFunctionOptions{
		OrReplace: r.OrReplace,
		Temporary: r.Temporary,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Memoizable:            r.Memoizable,
		Comment:               r.Comment,
		FunctionDefinition:    r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:      v.ArgName,
				ArgDataType:  v.ArgDataType,
				DefaultValue: v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataType: r.Returns.ResultDataType.ResultDataType,
		}
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
	return opts
}

func (r *AlterFunctionRequest) toOpts() *AlterFunctionOptions {
	opts := &AlterFunctionOptions{
		IfExists:          r.IfExists,
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
		RenameTo:          r.RenameTo,
		SetComment:        r.SetComment,
		SetLogLevel:       r.SetLogLevel,
		SetTraceLevel:     r.SetTraceLevel,
		SetSecure:         r.SetSecure,
		UnsetSecure:       r.UnsetSecure,
		UnsetLogLevel:     r.UnsetLogLevel,
		UnsetTraceLevel:   r.UnsetTraceLevel,
		UnsetComment:      r.UnsetComment,
		SetTags:           r.SetTags,
		UnsetTags:         r.UnsetTags,
	}
	return opts
}

func (r *DropFunctionRequest) toOpts() *DropFunctionOptions {
	opts := &DropFunctionOptions{
		IfExists:          r.IfExists,
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
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
	e := &Function{
		CreatedOn:          r.CreatedOn,
		Name:               r.Name,
		SchemaName:         r.SchemaName,
		IsBuiltin:          r.IsBuiltin == "Y",
		IsAggregate:        r.IsAggregate == "Y",
		IsAnsi:             r.IsAnsi == "Y",
		MinNumArguments:    r.MinNumArguments,
		MaxNumArguments:    r.MaxNumArguments,
		Arguments:          r.Arguments,
		Description:        r.Description,
		CatalogName:        r.CatalogName,
		IsTableFunction:    r.IsTableFunction == "Y",
		ValidForClustering: r.ValidForClustering == "Y",
		IsExternalFunction: r.IsExternalFunction == "Y",
		Language:           r.Language,
	}
	if r.IsSecure.Valid {
		e.IsSecure = r.IsSecure.String == "Y"
	}
	if r.IsMemoizable.Valid {
		e.IsMemoizable = r.IsMemoizable.String == "Y"
	}
	return e
}

func (r *DescribeFunctionRequest) toOpts() *DescribeFunctionOptions {
	opts := &DescribeFunctionOptions{
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
	}
	return opts
}

func (r functionDetailRow) convert() *FunctionDetail {
	return &FunctionDetail{
		Property: r.Property,
		Value:    r.Value,
	}
}
