package sdk

import "context"

var _ = (*pipes)(nil)

type pipes struct {
	client *Client
}

func (v *pipes) Create(ctx context.Context, opts *PipeCreateOptions) error {
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Alter(ctx context.Context, opts *PipeAlterOptions) error {
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Drop(ctx context.Context, opts *PipeDropOptions) error {
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Show(ctx context.Context, opts *PipeShowOptions) ([]*Pipe, error) {
	dbRows, err := validateAndQuery[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList := make([]*Pipe, len(*dbRows))
	for i, row := range *dbRows {
		resultList[i] = row.toPipe()
	}

	return resultList, nil
}

func (v *pipes) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	pipes, err := v.Show(ctx, &PipeShowOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: id.SchemaIdentifier(),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, pipe := range pipes {
		if pipe.ID().name == id.Name() {
			return pipe, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

func (v *pipes) Describe(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	opts := &describePipeOptions{
		name: id,
	}
	pipeRow, err := validateAndQueryOne[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return pipeRow.toPipe(), nil
}
