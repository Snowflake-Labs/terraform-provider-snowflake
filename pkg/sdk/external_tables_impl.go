package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

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
	dbRows, err := validateAndQuery[externalTableRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	resultList := convertRows[externalTableRow, ExternalTable](dbRows)
	return resultList, nil
}

func (v *externalTables) ShowByID(ctx context.Context, req *ShowExternalTableByIDRequest) (*ExternalTable, error) {
	if !ValidObjectIdentifier(req.id) {
		return nil, ErrInvalidObjectIdentifier
	}

	externalTables, err := v.client.ExternalTables.Show(ctx, NewShowExternalTableRequest().
		WithIn(NewShowExternalTableInRequest().WithSchema(NewDatabaseObjectIdentifier(req.id.DatabaseName(), req.id.SchemaName()))).
		WithLike(String(req.id.Name())))
	if err != nil {
		return nil, err
	}

	return collections.FindOne(externalTables, func(t ExternalTable) bool { return t.ID().FullyQualifiedName() == req.id.FullyQualifiedName() })
}

func (v *externalTables) DescribeColumns(ctx context.Context, req *DescribeExternalTableColumnsRequest) ([]ExternalTableColumnDetails, error) {
	rows, err := validateAndQuery[externalTableColumnDetailsRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	return convertRows[externalTableColumnDetailsRow, ExternalTableColumnDetails](rows), nil
}

func (v *externalTables) DescribeStage(ctx context.Context, req *DescribeExternalTableStageRequest) ([]ExternalTableStageDetails, error) {
	rows, err := validateAndQuery[externalTableStageDetailsRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	return convertRows[externalTableStageDetailsRow, ExternalTableStageDetails](rows), nil
}
