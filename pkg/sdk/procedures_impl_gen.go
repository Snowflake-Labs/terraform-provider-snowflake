package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Procedures = (*procedures)(nil)

type procedures struct {
	client *Client
}

func (v *procedures) CreateForJava(ctx context.Context, request *CreateForJavaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForJavaScript(ctx context.Context, request *CreateForJavaScriptProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForPython(ctx context.Context, request *CreateForPythonProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForScala(ctx context.Context, request *CreateForScalaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForSQL(ctx context.Context, request *CreateForSQLProcedureRequest) error {
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
	request := NewShowProcedureRequest().WithIn(&In{Database: NewAccountObjectIdentifier(id.DatabaseName())}).WithLike(&Like{String(id.Name())})
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

func (r *CreateForJavaProcedureRequest) toOpts() *CreateForJavaProcedureOptions {
	opts := &CreateForJavaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		NullInputBehavior:          r.NullInputBehavior,
		Comment:                    r.Comment,
		ExecuteAs:                  r.ExecuteAs,
		ProcedureDefinition:        r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument(v)
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
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
				s[i] = ProcedureColumn(v)
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage(v)
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport(v)
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForJavaScriptProcedureRequest) toOpts() *CreateForJavaScriptProcedureOptions {
	opts := &CreateForJavaScriptProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants:          r.CopyGrants,
		ResultDataType:      r.ResultDataType,
		NotNull:             r.NotNull,
		NullInputBehavior:   r.NullInputBehavior,
		Comment:             r.Comment,
		ExecuteAs:           r.ExecuteAs,
		ProcedureDefinition: r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument(v)
		}
		opts.Arguments = s
	}
	return opts
}

func (r *CreateForPythonProcedureRequest) toOpts() *CreateForPythonProcedureOptions {
	opts := &CreateForPythonProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		NullInputBehavior:          r.NullInputBehavior,
		Comment:                    r.Comment,
		ExecuteAs:                  r.ExecuteAs,
		ProcedureDefinition:        r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument(v)
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
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
				s[i] = ProcedureColumn(v)
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage(v)
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport(v)
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForScalaProcedureRequest) toOpts() *CreateForScalaProcedureOptions {
	opts := &CreateForScalaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:             r.Handler,
		TargetPath:          r.TargetPath,
		NullInputBehavior:   r.NullInputBehavior,
		Comment:             r.Comment,
		ExecuteAs:           r.ExecuteAs,
		ProcedureDefinition: r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument(v)
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
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
				s[i] = ProcedureColumn(v)
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage(v)
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport(v)
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForSQLProcedureRequest) toOpts() *CreateForSQLProcedureOptions {
	opts := &CreateForSQLProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		NullInputBehavior:   r.NullInputBehavior,
		Comment:             r.Comment,
		ExecuteAs:           r.ExecuteAs,
		ProcedureDefinition: r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument(v)
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureSQLReturns{
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
				s[i] = ProcedureColumn(v)
			}
			opts.Returns.Table.Columns = s
		}
	}
	return opts
}

func (r *AlterProcedureRequest) toOpts() *AlterProcedureOptions {
	opts := &AlterProcedureOptions{
		IfExists:          r.IfExists,
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
		RenameTo:          r.RenameTo,
		SetComment:        r.SetComment,
		SetLogLevel:       r.SetLogLevel,
		SetTraceLevel:     r.SetTraceLevel,
		UnsetComment:      r.UnsetComment,
		SetTags:           r.SetTags,
		UnsetTags:         r.UnsetTags,
		ExecuteAs:         r.ExecuteAs,
	}
	return opts
}

func (r *DropProcedureRequest) toOpts() *DropProcedureOptions {
	opts := &DropProcedureOptions{
		IfExists:          r.IfExists,
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
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
	e := &Procedure{
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
	}
	if r.IsSecure.Valid {
		e.IsSecure = r.IsSecure.String == "Y"
	}
	return e
}

func (r *DescribeProcedureRequest) toOpts() *DescribeProcedureOptions {
	opts := &DescribeProcedureOptions{
		name:              r.name,
		ArgumentDataTypes: r.ArgumentDataTypes,
	}
	return opts
}

func (r procedureDetailRow) convert() *ProcedureDetail {
	return &ProcedureDetail{
		Property: r.Property,
		Value:    r.Value,
	}
}
