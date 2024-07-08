package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Streamlits = (*streamlits)(nil)

type streamlits struct {
	client *Client
}

func (v *streamlits) Create(ctx context.Context, request *CreateStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) Alter(ctx context.Context, request *AlterStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) Drop(ctx context.Context, request *DropStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) Show(ctx context.Context, request *ShowStreamlitRequest) ([]Streamlit, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[streamlitsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[streamlitsRow, Streamlit](dbRows)
	return resultList, nil
}

func (v *streamlits) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Streamlit, error) {
	// TODO: adjust request if e.g. LIKE is supported for the resource
	streamlits, err := v.Show(ctx, NewShowStreamlitRequest())
	if err != nil {
		return nil, err
	}
	return collections.FindOne(streamlits, func(r Streamlit) bool { return r.Name == id.Name() })
}

func (v *streamlits) Describe(ctx context.Context, id SchemaObjectIdentifier) (*StreamlitDetail, error) {
	opts := &DescribeStreamlitOptions{
		name: id,
	}
	result, err := validateAndQueryOne[streamlitsDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateStreamlitRequest) toOpts() *CreateStreamlitOptions {
	opts := &CreateStreamlitOptions{
		OrReplace:    r.OrReplace,
		IfNotExists:  r.IfNotExists,
		name:         r.name,
		RootLocation: r.RootLocation,
		MainFile:     r.MainFile,
		Warehouse:    r.Warehouse,

		Title:   r.Title,
		Comment: r.Comment,
	}

	if r.ExternalAccessIntegrations != nil {
		opts.ExternalAccessIntegrations = &ExternalAccessIntegrations{
			ExternalAccessIntegrations: r.ExternalAccessIntegrations.ExternalAccessIntegrations,
		}
	}

	return opts
}

func (r *AlterStreamlitRequest) toOpts() *AlterStreamlitOptions {
	opts := &AlterStreamlitOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo: r.RenameTo,
	}

	if r.Set != nil {

		opts.Set = &StreamlitSet{
			RootLocation: r.Set.RootLocation,
			MainFile:     r.Set.MainFile,
			Warehouse:    r.Set.Warehouse,

			Comment: r.Set.Comment,
			Title:   r.Set.Title,
		}

		if r.Set.ExternalAccessIntegrations != nil {
			opts.Set.ExternalAccessIntegrations = &ExternalAccessIntegrations{
				ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations.ExternalAccessIntegrations,
			}
		}

	}

	if r.Unset != nil {
		opts.Unset = &StreamlitUnset{
			QueryWarehouse: r.Unset.QueryWarehouse,
			Comment:        r.Unset.Comment,
			Title:          r.Unset.Title,
		}
	}

	return opts
}

func (r *DropStreamlitRequest) toOpts() *DropStreamlitOptions {
	opts := &DropStreamlitOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowStreamlitRequest) toOpts() *ShowStreamlitOptions {
	opts := &ShowStreamlitOptions{
		Terse: r.Terse,
		Like:  r.Like,
		In:    r.In,
		Limit: r.Limit,
	}
	return opts
}

func (r streamlitsRow) convert() *Streamlit {
	// TODO: Mapping
	return &Streamlit{}
}

func (r *DescribeStreamlitRequest) toOpts() *DescribeStreamlitOptions {
	opts := &DescribeStreamlitOptions{
		name: r.name,
	}
	return opts
}

func (r streamlitsDetailRow) convert() *StreamlitDetail {
	// TODO: Mapping
	return &StreamlitDetail{}
}
