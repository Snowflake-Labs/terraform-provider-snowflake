package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Procedures = (*procedures)(nil)

type procedures struct {
	client *Client
}

func (v *procedures) CreateProcedureForJava(ctx context.Context, request *CreateProcedureForJavaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateProcedureForJavaScript(ctx context.Context, request *CreateProcedureForJavaScriptProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateProcedureForPython(ctx context.Context, request *CreateProcedureForPythonProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateProcedureForScala(ctx context.Context, request *CreateProcedureForScalaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateProcedureForSQL(ctx context.Context, request *CreateProcedureForSQLProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) Alter(ctx context.Context, request *AlterProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) Drop(ctx context.Context, request *DropProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) Show(ctx context.Context, request *ShowProcedureRequest) ([]Procedure, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[procedureRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[procedureRow, Procedure](dbRows)
	return resultList, nil
}

func (v *procedures) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Procedure, error) {
	request := NewShowProcedureRequest().WithLike(id.Name())
	procedures, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(procedures, func(r Procedure) bool { return r.Name == id.Name() })
}

func (v *procedures) Describe(ctx context.Context, request *DescribeProcedureRequest) ([]ProcedureDetail, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[procedureDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[procedureDetailRow, ProcedureDetail](rows), nil
}

func (r *CreateProcedureForJavaProcedureRequest) toOpts() *CreateProcedureForJavaProcedureOptions {
	opts := &CreateProcedureForJavaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,

		TargetPath: r.TargetPath,

		Comment: r.Comment,

		As: r.As,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &ProcedureReturns{}
		if r.Returns.ResultDataType != nil {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
				ResultDataType: r.Returns.ResultDataType.ResultDataType,
				Null:           r.Returns.ResultDataType.Null,
				NotNull:        r.Returns.ResultDataType.NotNull,
			}
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &ProcedureReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = ProcedureColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Secrets != nil {
		s := make([]ProcedureSecret, len(r.Secrets))
		for i, v := range r.Secrets {
			s[i] = ProcedureSecret{
				SecretVariableName: v.SecretVariableName,
				SecretName:         v.SecretName,
			}
		}
		opts.Secrets = s
	}
	if r.StrictOrNot != nil {
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict:            r.StrictOrNot.Strict,
			CalledOnNullInput: r.StrictOrNot.CalledOnNullInput,
		}
	}
	if r.VolatileOrNot != nil {
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile:  r.VolatileOrNot.Volatile,
			Immutable: r.VolatileOrNot.Immutable,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *CreateProcedureForJavaScriptProcedureRequest) toOpts() *CreateProcedureForJavaScriptProcedureOptions {
	opts := &CreateProcedureForJavaScriptProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		Comment: r.Comment,

		As: r.As,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &ProcedureReturns2{
			ResultDataType: r.Returns.ResultDataType,
			NotNull:        r.Returns.NotNull,
		}
	}
	if r.StrictOrNot != nil {
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict:            r.StrictOrNot.Strict,
			CalledOnNullInput: r.StrictOrNot.CalledOnNullInput,
		}
	}
	if r.VolatileOrNot != nil {
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile:  r.VolatileOrNot.Volatile,
			Immutable: r.VolatileOrNot.Immutable,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *CreateProcedureForPythonProcedureRequest) toOpts() *CreateProcedureForPythonProcedureOptions {
	opts := &CreateProcedureForPythonProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,

		Comment: r.Comment,

		As: r.As,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &ProcedureReturns{}
		if r.Returns.ResultDataType != nil {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
				ResultDataType: r.Returns.ResultDataType.ResultDataType,
				Null:           r.Returns.ResultDataType.Null,
				NotNull:        r.Returns.ResultDataType.NotNull,
			}
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &ProcedureReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = ProcedureColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Secrets != nil {
		s := make([]ProcedureSecret, len(r.Secrets))
		for i, v := range r.Secrets {
			s[i] = ProcedureSecret{
				SecretVariableName: v.SecretVariableName,
				SecretName:         v.SecretName,
			}
		}
		opts.Secrets = s
	}
	if r.StrictOrNot != nil {
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict:            r.StrictOrNot.Strict,
			CalledOnNullInput: r.StrictOrNot.CalledOnNullInput,
		}
	}
	if r.VolatileOrNot != nil {
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile:  r.VolatileOrNot.Volatile,
			Immutable: r.VolatileOrNot.Immutable,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *CreateProcedureForScalaProcedureRequest) toOpts() *CreateProcedureForScalaProcedureOptions {
	opts := &CreateProcedureForScalaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:    r.Handler,
		TargetPath: r.TargetPath,

		Comment: r.Comment,

		As: r.As,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &ProcedureReturns{}
		if r.Returns.ResultDataType != nil {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
				ResultDataType: r.Returns.ResultDataType.ResultDataType,
				Null:           r.Returns.ResultDataType.Null,
				NotNull:        r.Returns.ResultDataType.NotNull,
			}
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &ProcedureReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = ProcedureColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.StrictOrNot != nil {
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict:            r.StrictOrNot.Strict,
			CalledOnNullInput: r.StrictOrNot.CalledOnNullInput,
		}
	}
	if r.VolatileOrNot != nil {
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile:  r.VolatileOrNot.Volatile,
			Immutable: r.VolatileOrNot.Immutable,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *CreateProcedureForSQLProcedureRequest) toOpts() *CreateProcedureForSQLProcedureOptions {
	opts := &CreateProcedureForSQLProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		Comment: r.Comment,

		As: r.As,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:     v.ArgName,
				ArgDataType: v.ArgDataType,
			}
		}
		opts.Arguments = s
	}
	if r.Returns != nil {
		opts.Returns = &ProcedureReturns3{
			NotNull: r.Returns.NotNull,
		}
		if r.Returns.ResultDataType != nil {
			opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
				ResultDataType: r.Returns.ResultDataType.ResultDataType,
			}
		}
		if r.Returns.Table != nil {
			opts.Returns.Table = &ProcedureReturnsTable{}
			if r.Returns.Table.Columns != nil {
				s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
				for i, v := range r.Returns.Table.Columns {
					s[i] = ProcedureColumn{
						ColumnName:     v.ColumnName,
						ColumnDataType: v.ColumnDataType,
					}
				}
				opts.Returns.Table.Columns = s
			}
		}
	}
	if r.StrictOrNot != nil {
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict:            r.StrictOrNot.Strict,
			CalledOnNullInput: r.StrictOrNot.CalledOnNullInput,
		}
	}
	if r.VolatileOrNot != nil {
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile:  r.VolatileOrNot.Volatile,
			Immutable: r.VolatileOrNot.Immutable,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *AlterProcedureRequest) toOpts() *AlterProcedureOptions {
	opts := &AlterProcedureOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo:  r.RenameTo,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.ArgumentTypes != nil {
		s := make([]ProcedureArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = ProcedureArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	if r.Set != nil {
		opts.Set = &ProcedureSet{
			LogLevel:   r.Set.LogLevel,
			TraceLevel: r.Set.TraceLevel,
			Comment:    r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ProcedureUnset{
			Comment: r.Unset.Comment,
		}
	}
	if r.ExecuteAs != nil {
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: r.ExecuteAs.Caller,
			Owner:  r.ExecuteAs.Owner,
		}
	}
	return opts
}

func (r *DropProcedureRequest) toOpts() *DropProcedureOptions {
	opts := &DropProcedureOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.ArgumentTypes != nil {
		s := make([]ProcedureArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = ProcedureArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	return opts
}

func (r *ShowProcedureRequest) toOpts() *ShowProcedureOptions {
	opts := &ShowProcedureOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r procedureRow) convert() *Procedure {
	return &Procedure{
		CreatedOn:       r.CreatedOn,
		Name:            r.Name,
		SchemaName:      r.SchemaName,
		MinNumArguments: r.MinNumArguments,
		MaxNumArguments: r.MaxNumArguments,
		Arguments:       r.Arguments,
		IsTableFunction: r.IsTableFunction,
	}
}

func (r *DescribeProcedureRequest) toOpts() *DescribeProcedureOptions {
	opts := &DescribeProcedureOptions{
		name: r.name,
	}
	if r.ArgumentTypes != nil {
		s := make([]ProcedureArgumentType, len(r.ArgumentTypes))
		for i, v := range r.ArgumentTypes {
			s[i] = ProcedureArgumentType{
				ArgDataType: v.ArgDataType,
			}
		}
		opts.ArgumentTypes = s
	}
	return opts
}

func (r procedureDetailRow) convert() *ProcedureDetail {
	return &ProcedureDetail{
		Property: r.Property,
		Value:    r.Value,
	}
}
