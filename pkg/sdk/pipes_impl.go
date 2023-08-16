package sdk

import "context"

var _ Pipes = (*pipes)(nil)

type pipes struct {
	client *Client
}

func (v *pipes) Create(ctx context.Context, id SchemaObjectIdentifier, copyStatement string, opts *PipeCreateOptions) error {
	opts = createIfNil[PipeCreateOptions](opts)
	opts.name = id
	opts.copyStatement = copyStatement
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *PipeAlterOptions) error {
	opts = createIfNil[PipeAlterOptions](opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Drop(ctx context.Context, id SchemaObjectIdentifier) error {
	opts := &PipeDropOptions{
		name: id,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Show(ctx context.Context, opts *PipeShowOptions) ([]Pipe, error) {
	dbRows, err := validateAndQuery[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList := convertRows[pipeDBRow, Pipe](dbRows)

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

	return findOne(pipes, func(p Pipe) bool { return p.ID().name == id.Name() })
}

func (v *pipes) Describe(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	opts := &describePipeOptions{
		name: id,
	}
	pipeRow, err := validateAndQueryOne[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return pipeRow.convert(), nil
}
