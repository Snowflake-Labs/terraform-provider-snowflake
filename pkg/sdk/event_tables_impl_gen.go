package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ EventTables = (*eventTables)(nil)

type eventTables struct {
	client *Client
}

func (v *eventTables) Create(ctx context.Context, request *CreateEventTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *eventTables) Show(ctx context.Context, request *ShowEventTableRequest) ([]EventTable, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[eventTableRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[eventTableRow, EventTable](dbRows)
	return resultList, nil
}

func (v *eventTables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*EventTable, error) {
	request := NewShowEventTableRequest().WithLike(id.Name())
	functions, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(functions, func(r EventTable) bool { return r.Name == id.Name() })
}

func (v *eventTables) Describe(ctx context.Context, request *DescribeEventTableRequest) (*EventTableDetails, error) {
	opts := request.toOpts()
	result, err := validateAndQueryOne[eventTableDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (v *eventTables) Alter(ctx context.Context, request *AlterEventTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateEventTableRequest) toOpts() *CreateEventTableOptions {
	opts := &CreateEventTableOptions{
		OrReplace:                  r.OrReplace,
		IfNotExists:                r.IfNotExists,
		name:                       r.name,
		ClusterBy:                  r.ClusterBy,
		DataRetentionTimeInDays:    r.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: r.MaxDataExtensionTimeInDays,
		ChangeTracking:             r.ChangeTracking,
		DefaultDdlCollation:        r.DefaultDdlCollation,
		CopyGrants:                 r.CopyGrants,
		Comment:                    r.Comment,
		RowAccessPolicy:            r.RowAccessPolicy,
		Tag:                        r.Tag,
	}
	return opts
}

func (r *ShowEventTableRequest) toOpts() *ShowEventTableOptions {
	opts := &ShowEventTableOptions{
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
		From:       r.From,
	}
	return opts
}

func (r eventTableRow) convert() *EventTable {
	return &EventTable{
		CreatedOn:      r.CreatedOn,
		Name:           r.Name,
		DatabaseName:   r.DatabaseName,
		SchemaName:     r.SchemaName,
		Owner:          r.Owner,
		Comment:        r.Comment,
		OwnerRoleType:  r.OwnerRoleType,
		ChangeTracking: r.ChangeTracking == "ON",
	}
}

func (r *DescribeEventTableRequest) toOpts() *DescribeEventTableOptions {
	opts := &DescribeEventTableOptions{
		name: r.name,
	}
	return opts
}

func (r eventTableDetailsRow) convert() *EventTableDetails {
	return &EventTableDetails{
		Name:    r.Name,
		Kind:    r.Kind,
		Comment: r.Comment,
	}
}

func (r *AlterEventTableRequest) toOpts() *AlterEventTableOptions {
	opts := &AlterEventTableOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,

		AddRowAccessPolicy: r.AddRowAccessPolicy,

		DropAllRowAccessPolicies: r.DropAllRowAccessPolicies,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
		RenameTo:  r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &EventTableSet{
			DataRetentionTimeInDays:    r.Set.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Set.MaxDataExtensionTimeInDays,
			ChangeTracking:             r.Set.ChangeTracking,
			Comment:                    r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &EventTableUnset{
			DataRetentionTimeInDays:    r.Unset.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Unset.MaxDataExtensionTimeInDays,
			ChangeTracking:             r.Unset.ChangeTracking,
			Comment:                    r.Unset.Comment,
		}
	}
	if r.DropRowAccessPolicy != nil {
		opts.DropRowAccessPolicy = &EventTableDropRowAccessPolicy{
			Name: r.DropRowAccessPolicy.Name,
		}
	}
	if r.ClusteringAction != nil {
		opts.ClusteringAction = &EventTableClusteringAction{
			ClusterBy:         r.ClusteringAction.ClusterBy,
			SuspendRecluster:  r.ClusteringAction.SuspendRecluster,
			ResumeRecluster:   r.ClusteringAction.ResumeRecluster,
			DropClusteringKey: r.ClusteringAction.DropClusteringKey,
		}
	}
	if r.SearchOptimizationAction != nil {
		opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{}
		if r.SearchOptimizationAction.Add != nil {
			opts.SearchOptimizationAction.Add = &SearchOptimization{
				On: r.SearchOptimizationAction.Drop.On,
			}
		}
		if r.SearchOptimizationAction.Drop != nil {
			opts.SearchOptimizationAction.Drop = &SearchOptimization{
				On: r.SearchOptimizationAction.Drop.On,
			}
		}
	}
	return opts
}
