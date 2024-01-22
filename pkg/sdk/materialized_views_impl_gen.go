package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ MaterializedViews = (*materializedViews)(nil)

type materializedViews struct {
	client *Client
}

func (v *materializedViews) Create(ctx context.Context, request *CreateMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) Alter(ctx context.Context, request *AlterMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) Drop(ctx context.Context, request *DropMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) Show(ctx context.Context, request *ShowMaterializedViewRequest) ([]MaterializedView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[materializedViewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[materializedViewDBRow, MaterializedView](dbRows)
	return resultList, nil
}

func (v *materializedViews) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaterializedView, error) {
	// TODO: adjust request if e.g. LIKE is supported for the resource
	materializedViews, err := v.Show(ctx, NewShowMaterializedViewRequest())
	if err != nil {
		return nil, err
	}
	return collections.FindOne(materializedViews, func(r MaterializedView) bool { return r.Name == id.Name() })
}

func (v *materializedViews) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]MaterializedViewDetails, error) {
	opts := &DescribeMaterializedViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[materializedViewDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[materializedViewDetailsRow, MaterializedViewDetails](rows), nil
}

func (r *CreateMaterializedViewRequest) toOpts() *CreateMaterializedViewOptions {
	opts := &CreateMaterializedViewOptions{
		OrReplace:   r.OrReplace,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		CopyGrants:  r.CopyGrants,

		Comment: r.Comment,

		Tag:       r.Tag,
		ClusterBy: r.ClusterBy,
		sql:       r.sql,
	}
	if r.Columns != nil {
		s := make([]MaterializedViewColumn, len(r.Columns))
		for i, v := range r.Columns {
			s[i] = MaterializedViewColumn{
				Name:    v.Name,
				Comment: v.Comment,
			}
		}
		opts.Columns = s
	}
	if r.ColumnsMaskingPolicies != nil {
		s := make([]MaterializedViewColumnMaskingPolicy, len(r.ColumnsMaskingPolicies))
		for i, v := range r.ColumnsMaskingPolicies {
			s[i] = MaterializedViewColumnMaskingPolicy{
				Name:          v.Name,
				MaskingPolicy: v.MaskingPolicy,
				Using:         v.Using,
				Tag:           v.Tag,
			}
		}
		opts.ColumnsMaskingPolicies = s
	}
	if r.RowAccessPolicy != nil {
		opts.RowAccessPolicy = &MaterializedViewRowAccessPolicy{
			RowAccessPolicy: r.RowAccessPolicy.RowAccessPolicy,
			On:              r.RowAccessPolicy.On,
		}
	}
	return opts
}

func (r *AlterMaterializedViewRequest) toOpts() *AlterMaterializedViewOptions {
	opts := &AlterMaterializedViewOptions{
		name:              r.name,
		RenameTo:          r.RenameTo,
		ClusterBy:         r.ClusterBy,
		DropClusteringKey: r.DropClusteringKey,
		SuspendRecluster:  r.SuspendRecluster,
		ResumeRecluster:   r.ResumeRecluster,
		Suspend:           r.Suspend,
		Resume:            r.Resume,
	}
	if r.Set != nil {
		opts.Set = &MaterializedViewSet{
			Secure:  r.Set.Secure,
			Comment: r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &MaterializedViewUnset{
			Secure:  r.Unset.Secure,
			Comment: r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropMaterializedViewRequest) toOpts() *DropMaterializedViewOptions {
	opts := &DropMaterializedViewOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowMaterializedViewRequest) toOpts() *ShowMaterializedViewOptions {
	opts := &ShowMaterializedViewOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r materializedViewDBRow) convert() *MaterializedView {
	// TODO: Mapping
	return &MaterializedView{}
}

func (r *DescribeMaterializedViewRequest) toOpts() *DescribeMaterializedViewOptions {
	opts := &DescribeMaterializedViewOptions{
		name: r.name,
	}
	return opts
}

func (r materializedViewDetailsRow) convert() *MaterializedViewDetails {
	// TODO: Mapping
	return &MaterializedViewDetails{}
}
