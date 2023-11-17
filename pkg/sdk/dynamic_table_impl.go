package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ DynamicTables = (*dynamicTables)(nil)

type dynamicTables struct {
	client *Client
}

func (v *dynamicTables) Create(ctx context.Context, request *CreateDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) Alter(ctx context.Context, request *AlterDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) Drop(ctx context.Context, request *DropDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) Describe(ctx context.Context, request *DescribeDynamicTableRequest) (*DynamicTableDetails, error) {
	opts := request.toOpts()
	row, err := validateAndQueryOne[dynamicTableDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return row.convert(), nil
}

func (v *dynamicTables) Show(ctx context.Context, request *ShowDynamicTableRequest) ([]DynamicTable, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[dynamicTableRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	result := convertRows[dynamicTableRow, DynamicTable](rows)
	return result, nil
}

func (v *dynamicTables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*DynamicTable, error) {
	request := NewShowDynamicTableRequest().WithIn(&In{Schema: NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())}).WithLike(&Like{String(id.Name())})
	dynamicTables, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(dynamicTables, func(r DynamicTable) bool { return r.Name == id.Name() })
}

func (s *CreateDynamicTableRequest) toOpts() *createDynamicTableOptions {
	return &createDynamicTableOptions{
		OrReplace: Bool(s.orReplace),
		name:      s.name,
		warehouse: s.warehouse,
		targetLag: s.targetLag,
		query:     s.query,
		Comment:   s.comment,
	}
}

func (s *AlterDynamicTableRequest) toOpts() *alterDynamicTableOptions {
	opts := alterDynamicTableOptions{
		name: s.name,
	}
	if s.suspend != nil {
		opts.Suspend = s.suspend
	}
	if s.resume != nil {
		opts.Resume = s.resume
	}
	if s.refresh != nil {
		opts.Refresh = s.refresh
	}
	if s.set != nil {
		opts.Set = &DynamicTableSet{s.set.targetLag, s.set.warehourse}
	}
	return &opts
}

func (s *DropDynamicTableRequest) toOpts() *dropDynamicTableOptions {
	return &dropDynamicTableOptions{
		name: s.name,
	}
}

func (s *DescribeDynamicTableRequest) toOpts() *describeDynamicTableOptions {
	return &describeDynamicTableOptions{
		name: s.name,
	}
}

func (s *ShowDynamicTableRequest) toOpts() *showDynamicTableOptions {
	opts := showDynamicTableOptions{}
	if s.like != nil {
		opts.Like = s.like
	}
	if s.in != nil {
		opts.In = s.in
	}
	if s.startsWith != nil {
		opts.StartsWith = s.startsWith
	}
	if s.limit != nil {
		opts.Limit = s.limit
	}
	return &opts
}
