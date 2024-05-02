package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

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

func (v *streams) CreateOnDirectoryTable(ctx context.Context, request *CreateOnDirectoryTableStreamRequest) error {
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

func (v *streams) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Stream, error) {
	streams, err := v.Show(ctx, NewShowStreamRequest().
		WithIn(&In{
			Schema: NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()),
		}).
		WithLike(&Like{Pattern: String(id.Name())}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(streams, func(r Stream) bool { return r.Name == id.Name() })
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
				Timestamp: r.On.Statement.Timestamp,
				Offset:    r.On.Statement.Offset,
				Statement: r.On.Statement.Statement,
				Stream:    r.On.Statement.Stream,
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
				Timestamp: r.On.Statement.Timestamp,
				Offset:    r.On.Statement.Offset,
				Statement: r.On.Statement.Statement,
				Stream:    r.On.Statement.Stream,
			},
		}
	}
	return opts
}

func (r *CreateOnDirectoryTableStreamRequest) toOpts() *CreateOnDirectoryTableStreamOptions {
	opts := &CreateOnDirectoryTableStreamOptions{
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
				Timestamp: r.On.Statement.Timestamp,
				Offset:    r.On.Statement.Offset,
				Statement: r.On.Statement.Statement,
				Stream:    r.On.Statement.Stream,
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
	s := &Stream{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.StaleAfter.Valid {
		s.StaleAfter = &r.StaleAfter.Time
	}
	if r.TableOn.Valid {
		s.TableOn = &r.TableOn.String
	}
	if r.Owner.Valid {
		s.Owner = &r.Owner.String
	}
	if r.Comment.Valid {
		s.Comment = &r.Comment.String
	}
	if r.TableName.Valid {
		s.TableName = &r.TableName.String
	}
	if r.SourceType.Valid {
		s.SourceType = &r.SourceType.String
	}
	if r.BaseTables.Valid {
		s.BaseTables = &r.BaseTables.String
	}
	if r.Type.Valid {
		s.Type = &r.Type.String
	}
	if r.Stale.Valid {
		s.Stale = &r.Stale.String
	}
	if r.Mode.Valid {
		s.Mode = &r.Mode.String
	}
	if r.InvalidReason.Valid {
		s.InvalidReason = &r.InvalidReason.String
	}
	if r.OwnerRoleType.Valid {
		s.OwnerRoleType = &r.OwnerRoleType.String
	}
	return s
}

func (r *DescribeStreamRequest) toOpts() *DescribeStreamOptions {
	opts := &DescribeStreamOptions{
		name: r.name,
	}
	return opts
}
