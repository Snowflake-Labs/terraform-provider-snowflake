package sdk

import "context"

var _ ExternalTables = (*externalTables)(nil)

type externalTables struct {
	client *Client
}

func (v *externalTables) Create(ctx context.Context, req *CreateExternalTableRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) CreateWithManualPartitioning(ctx context.Context, req *CreateWithManualPartitioningExternalTableRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) CreateDeltaLake(ctx context.Context, req *CreateDeltaLakeExternalTableRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) CreateUsingTemplate(ctx context.Context, req *CreateExternalTableUsingTemplateRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) Alter(ctx context.Context, req *AlterExternalTableRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) AlterPartitions(ctx context.Context, req *AlterExternalTablePartitionRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) Drop(ctx context.Context, req *DropExternalTableRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *externalTables) Show(ctx context.Context, req *ShowExternalTableRequest) ([]ExternalTable, error) {
	rows, err := validateAndQuery[externalTableRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}

	externalTables := make([]ExternalTable, len(rows))
	for i, row := range rows {
		externalTables[i] = row.ToExternalTable()
	}

	return externalTables, err
}

func (v *externalTables) ShowByID(ctx context.Context, req *ShowExternalTableByIDRequest) (*ExternalTable, error) {
	if !validObjectidentifier(req.id) {
		return nil, ErrInvalidObjectIdentifier
	}

	externalTables, err := v.client.ExternalTables.Show(ctx, NewShowExternalTableRequest().WithLike(String(req.id.Name())))
	if err != nil {
		return nil, err
	}

	for _, t := range externalTables {
		if t.ID() == req.id {
			return &t, nil
		}
	}

	return nil, ErrObjectNotExistOrAuthorized
}

func (v *externalTables) DescribeColumns(ctx context.Context, req *DescribeExternalTableColumnsRequest) ([]ExternalTableColumnDetails, error) {
	rows, err := validateAndQuery[externalTableColumnDetailsRow](v.client, ctx, &describeExternalTableColumns{
		name: req.id,
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

func (v *externalTables) DescribeStage(ctx context.Context, req *DescribeExternalTableStageRequest) ([]ExternalTableStageDetails, error) {
	rows, err := validateAndQuery[externalTableStageDetailsRow](v.client, ctx, &describeExternalTableStage{
		name: req.id,
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
