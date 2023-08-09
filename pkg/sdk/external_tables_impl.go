package sdk

import "context"

var _ ExternalTables = (*externalTables)(nil)

type externalTables struct {
	client *Client
}

func (v *externalTables) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateExternalTableOpts) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) CreateWithManualPartitioning(ctx context.Context, id AccountObjectIdentifier, opts *CreateWithManualPartitioningExternalTableOpts) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) CreateDeltaLake(ctx context.Context, id AccountObjectIdentifier, opts *CreateDeltaLakeExternalTableOpts) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTableOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) AlterPartitions(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTablePartitionOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropExternalTableOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalTables) Show(ctx context.Context, opts *ShowExternalTableOptions) ([]ExternalTable, error) {
	opts = createIfNil(opts)
	rows, err := validateAndQuery[externalTableRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	externalTables := make([]ExternalTable, len(rows))
	for i, row := range rows {
		externalTables[i] = row.ToExternalTable()
	}

	return externalTables, err
}

func (v *externalTables) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalTable, error) {
	if !validObjectidentifier(id) {
		return nil, ErrInvalidObjectIdentifier
	}

	externalTables, err := v.client.ExternalTables.Show(ctx, &ShowExternalTableOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, t := range externalTables {
		if t.ID() == id {
			return &t, nil
		}
	}

	return nil, ErrObjectNotExistOrAuthorized
}

func (v *externalTables) DescribeColumns(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableColumnDetails, error) {
	rows, err := validateAndQuery[externalTableColumnDetailsRow](v.client, ctx, &describeExternalTableColumns{
		name: id,
	})
	if err != nil {
		return nil, err
	}

	var result []ExternalTableColumnDetails
	for _, r := range rows {
		result = append(result, r.toExternalTableColumnDetails())
	}
	return result, nil
}

func (v *externalTables) DescribeStage(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableStageDetails, error) {
	rows, err := validateAndQuery[externalTableStageDetailsRow](v.client, ctx, &describeExternalTableStage{
		name: id,
	})
	if err != nil {
		return nil, err
	}

	var result []ExternalTableStageDetails
	for _, r := range rows {
		result = append(result, r.toExternalTableStageDetails())
	}

	return result, nil
}
