package sdk

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ CortexSearchServices = (*cortexSearchServices)(nil)

type cortexSearchServices struct {
	client *Client
}

func (v *cortexSearchServices) Create(ctx context.Context, request *CreateCortexSearchServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *cortexSearchServices) Alter(ctx context.Context, request *AlterCortexSearchServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *cortexSearchServices) Show(ctx context.Context, request *ShowCortexSearchServiceRequest) ([]CortexSearchService, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[cortexSearchServiceRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[cortexSearchServiceRow, CortexSearchService](dbRows)
	return resultList, nil
}

func (v *cortexSearchServices) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchService, error) {
	request := NewShowCortexSearchServiceRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
	cortexSearchServices, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(cortexSearchServices, func(r CortexSearchService) bool { return r.Name == id.Name() })
}

func (v *cortexSearchServices) Describe(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchServiceDetails, error) {
	opts := &DescribeCortexSearchServiceOptions{
		name: id,
	}
	result, err := validateAndQueryOne[cortexSearchServiceDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (v *cortexSearchServices) Drop(ctx context.Context, request *DropCortexSearchServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateCortexSearchServiceRequest) toOpts() *CreateCortexSearchServiceOptions {
	opts := &CreateCortexSearchServiceOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		On:          r.On,

		Warehouse:       r.Warehouse,
		TargetLag:       r.TargetLag,
		Comment:         r.Comment,
		QueryDefinition: r.QueryDefinition,
	}

	if r.Attributes != nil {
		opts.Attributes = &Attributes{
			Columns: r.Attributes.Columns,
		}
	}

	return opts
}

func (r *AlterCortexSearchServiceRequest) toOpts() *AlterCortexSearchServiceOptions {
	opts := &AlterCortexSearchServiceOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}

	if r.Set != nil {
		opts.Set = &CortexSearchServiceSet{
			TargetLag: r.Set.TargetLag,
			Warehouse: r.Set.Warehouse,
			Comment:   r.Set.Comment,
		}
	}

	return opts
}

func (r *ShowCortexSearchServiceRequest) toOpts() *ShowCortexSearchServiceOptions {
	opts := &ShowCortexSearchServiceOptions{
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r cortexSearchServiceRow) convert() *CortexSearchService {
	cortexSearchService := &CortexSearchService{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Comment.Valid {
		cortexSearchService.Comment = r.Comment.String
	}
	return cortexSearchService
}

func (r *DescribeCortexSearchServiceRequest) toOpts() *DescribeCortexSearchServiceOptions {
	opts := &DescribeCortexSearchServiceOptions{
		name: r.name,
	}
	return opts
}

func (r cortexSearchServiceDetailsRow) convert() *CortexSearchServiceDetails {
	row := &CortexSearchServiceDetails{
		CreatedOn:         r.CreatedOn,
		Name:              r.Name,
		DatabaseName:      r.DatabaseName,
		SchemaName:        r.SchemaName,
		TargetLag:         r.TargetLag,
		Warehouse:         r.Warehouse,
		ServiceQueryUrl:   r.ServiceQueryUrl,
		DataTimestamp:     r.DataTimestamp,
		SourceDataNumRows: r.SourceDataNumRows,
		IndexingState:     r.IndexingState,
	}
	if r.SearchColumn.Valid {
		row.SearchColumn = String(r.SearchColumn.String)
	}
	if r.AttributeColumns.Valid {
		row.AttributeColumns = strings.Split(r.AttributeColumns.String, ",")
	}
	if r.Columns.Valid {
		row.Columns = strings.Split(r.Columns.String, ",")
	}
	if r.Definition.Valid {
		row.Definition = String(r.Definition.String)
	}
	if r.Comment.Valid {
		row.Comment = String(r.Comment.String)
	}
	if r.IndexingError.Valid {
		row.IndexingError = String(r.IndexingError.String)
	}

	return row
}

func (r *DropCortexSearchServiceRequest) toOpts() *DropCortexSearchServiceOptions {
	opts := &DropCortexSearchServiceOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}
