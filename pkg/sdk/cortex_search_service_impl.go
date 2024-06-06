package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
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

func (v *cortexSearchServices) Drop(ctx context.Context, request *DropCortexSearchServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *cortexSearchServices) Describe(ctx context.Context, request *DescribeCortexSearchServiceRequest) (*CortexSearchServiceDetails, error) {
	opts := request.toOpts()
	row, err := validateAndQueryOne[cortexSearchServiceDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return row.convert(), nil
}

func (v *cortexSearchServices) Show(ctx context.Context, request *ShowCortexSearchServiceRequest) ([]CortexSearchService, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[cortexSearchServiceRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	result := convertRows[cortexSearchServiceRow, CortexSearchService](rows)
	return result, nil
}

func (v *cortexSearchServices) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchService, error) {
	request := NewShowCortexSearchServiceRequest().WithIn(&In{Schema: id.SchemaId()}).WithLike(&Like{String(id.Name())})
	cortexSearchServices, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindOne(cortexSearchServices, func(r CortexSearchService) bool { return r.Name == id.Name() })
}

func (s *CreateCortexSearchServiceRequest) toOpts() *createCortexSearchServiceOptions {
	opts := createCortexSearchServiceOptions{
		OrReplace:   Bool(s.orReplace),
		IfNotExists: Bool(s.ifNotExists),
		name:        s.name,
		on:          s.on,
		warehouse:   s.warehouse,
		targetLag:   s.targetLag,
		query:       s.query,
		Comment:     s.comment,
	}
	if s.attributes != nil {
		opts.attributes = &Attributes{columns: s.attributes}
	}
	return &opts
}

func (s *AlterCortexSearchServiceRequest) toOpts() *alterCortexSearchServiceOptions {
	opts := alterCortexSearchServiceOptions{
		name:     s.name,
		IfExists: s.IfExists,
	}
	if s.set != nil {
		opts.Set = &CortexSearchServiceSet{s.set.targetLag, s.set.warehouse}
	}
	return &opts
}

func (s *DropCortexSearchServiceRequest) toOpts() *dropCortexSearchServiceOptions {
	return &dropCortexSearchServiceOptions{
		name:     s.name,
		IfExists: s.IfExists,
	}
}

func (s *DescribeCortexSearchServiceRequest) toOpts() *describeCortexSearchServiceOptions {
	return &describeCortexSearchServiceOptions{
		name: s.name,
	}
}

func (s *ShowCortexSearchServiceRequest) toOpts() *showCortexSearchServiceOptions {
	opts := showCortexSearchServiceOptions{}
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
