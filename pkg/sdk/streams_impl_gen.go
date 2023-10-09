package sdk

import "context"

var _ Streams = (*streams)(nil)

type streams struct {
	client *Client
}

func (v *streams) CreateOnTable(ctx context.Context, request *CreateOnTableStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) CreateOnExternalTable(ctx context.Context, request *CreateOnExternalTableStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) CreateOnStage(ctx context.Context, request *CreateOnStageStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) CreateOnView(ctx context.Context, request *CreateOnViewStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) Clone(ctx context.Context, request *CloneStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) Alter(ctx context.Context, request *AlterStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) Drop(ctx context.Context, request *DropStreamRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streams) Show(ctx context.Context, request *ShowStreamRequest) ([]Stream, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showStreamsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showStreamsDbRow, Stream](dbRows)
	return resultList, nil
}

func (v *streams) ShowByID(ctx context.Context, request *ShowByIdStreamRequest) (*Stream, error) {
	// TODO: adjust request if e.g. LIKE is supported for the resource
	streams, err := v.Show(ctx, NewShowStreamRequest())
	if err != nil {
		return nil, err
	}
	return findOne(streams, func(r Stream) bool { return r.Name == request.name.Name() })
}

func (v *streams) Describe(ctx context.Context, request *DescribeStreamRequest) (*Stream, error) {
	opts := &DescribeStreamOptions{
		name: request.name,
	}
	result, err := validateAndQueryOne[showStreamsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateOnTableStreamRequest) toOpts() *CreateOnTableStreamOptions {
	opts := &CreateOnTableStreamOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		CopyGrants:  r.CopyGrants,
		TableId:     r.TableId,

		AppendOnly:      r.AppendOnly,
		ShowInitialRows: r.ShowInitialRows,
		Comment:         r.Comment,
	}
	if r.On != nil {
		opts.On = &OnStream{
			At:     r.On.At,
			Before: r.On.Before,
			Statement: OnStreamStatement{
				Timestamp: r.On.Timestamp,
				Offset:    r.On.Offset,
				Statement: r.On.Statement,
				Stream:    r.On.Stream,
			},
		}
	}
	return opts
}

func (r *CreateOnExternalTableStreamRequest) toOpts() *CreateOnExternalTableStreamOptions {
	opts := &CreateOnExternalTableStreamOptions{
		OrReplace:       r.OrReplace,
		IfNotExists:     r.IfNotExists,
		name:            r.name,
		CopyGrants:      r.CopyGrants,
		ExternalTableId: r.ExternalTableId,

		InsertOnly: r.InsertOnly,
		Comment:    r.Comment,
	}
	if r.On != nil {
		opts.On = &OnStream{
			At:     r.On.At,
			Before: r.On.Before,
			Statement: OnStreamStatement{
				Timestamp: r.On.Timestamp,
				Offset:    r.On.Offset,
				Statement: r.On.Statement,
				Stream:    r.On.Stream,
			},
		}
	}
	return opts
}

func (r *CreateOnStageStreamRequest) toOpts() *CreateOnStageStreamOptions {
	opts := &CreateOnStageStreamOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		CopyGrants:  r.CopyGrants,
		StageId:     r.StageId,
		Comment:     r.Comment,
	}
	return opts
}

func (r *CreateOnViewStreamRequest) toOpts() *CreateOnViewStreamOptions {
	opts := &CreateOnViewStreamOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		CopyGrants:  r.CopyGrants,
		ViewId:      r.ViewId,

		AppendOnly:      r.AppendOnly,
		ShowInitialRows: r.ShowInitialRows,
		Comment:         r.Comment,
	}
	if r.On != nil {
		opts.On = &OnStream{
			At:     r.On.At,
			Before: r.On.Before,
			Statement: OnStreamStatement{
				Timestamp: r.On.Timestamp,
				Offset:    r.On.Offset,
				Statement: r.On.Statement,
				Stream:    r.On.Stream,
			},
		}
	}
	return opts
}

func (r *CloneStreamRequest) toOpts() *CloneStreamOptions {
	opts := &CloneStreamOptions{
		OrReplace:    r.OrReplace,
		name:         r.name,
		sourceStream: r.sourceStream,
		CopyGrants:   r.CopyGrants,
	}
	return opts
}

func (r *AlterStreamRequest) toOpts() *AlterStreamOptions {
	opts := &AlterStreamOptions{
		IfExists:     r.IfExists,
		name:         r.name,
		SetComment:   r.SetComment,
		UnsetComment: r.UnsetComment,
		SetTags:      r.SetTags,
		UnsetTags:    r.UnsetTags,
	}
	return opts
}

func (r *DropStreamRequest) toOpts() *DropStreamOptions {
	opts := &DropStreamOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowStreamRequest) toOpts() *ShowStreamOptions {
	opts := &ShowStreamOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r showStreamsDbRow) convert() *Stream {
	// TODO: Mapping
	return &Stream{}
}

func (r *DescribeStreamRequest) toOpts() *DescribeStreamOptions {
	opts := &DescribeStreamOptions{
		name: r.name,
	}
	return opts
}
