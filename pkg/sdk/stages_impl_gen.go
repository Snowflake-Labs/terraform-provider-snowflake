package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Stages = (*stages)(nil)

type stages struct {
	client *Client
}

func (v *stages) CreateInternal(ctx context.Context, request *CreateInternalStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) Drop(ctx context.Context, request *DropStageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *stages) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]StageProperty, error) {
	opts := &DescribeStageOptions{
		name: id,
	}
	rows, err := validateAndQuery[stageDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[stageDescRow, StageProperty](rows), nil
}

func (v *stages) Show(ctx context.Context, request *ShowStageRequest) ([]Stage, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[stageShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[stageShowRow, Stage](dbRows)
	return resultList, nil
}

func (r *CreateInternalStageRequest) toOpts() *CreateInternalStageOptions {
	opts := &CreateInternalStageOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Comment: r.Comment,
		Tag:     r.Tag,
	}
	if r.Encryption != nil {
		opts.Encryption = &InternalStageEncryption{
			Type: r.Encryption.Type,
		}
	}
	if r.DirectoryTableOptions != nil {
		opts.DirectoryTableOptions = &InternalDirectoryTableOptions{
			Enable:          r.DirectoryTableOptions.Enable,
			RefreshOnCreate: r.DirectoryTableOptions.RefreshOnCreate,
		}
	}
	if r.FileFormat != nil {
		opts.FileFormat = &StageFileFormat{
			FormatName: r.FileFormat.FormatName,
			Type:       r.FileFormat.Type,
			TYPE:       r.FileFormat.TYPE,
		}
	}
	if r.CopyOptions != nil {
		opts.CopyOptions = &StageCopyOptions{

			SizeLimit:         r.CopyOptions.SizeLimit,
			Purge:             r.CopyOptions.Purge,
			ReturnFailedOnly:  r.CopyOptions.ReturnFailedOnly,
			MatchByColumnName: r.CopyOptions.MatchByColumnName,
			EnforceLength:     r.CopyOptions.EnforceLength,
			Truncatecolumns:   r.CopyOptions.Truncatecolumns,
			Force:             r.CopyOptions.Force,
			Type:              r.CopyOptions.Type,
			TYPE:              r.CopyOptions.TYPE,
		}
		if r.CopyOptions.OnError != nil {
			opts.CopyOptions.OnError = &StageCopyOnErrorOptions{
				Continue: r.CopyOptions.OnError.Continue,
				SkipFile: r.CopyOptions.OnError.SkipFile,
				SkipFile: r.CopyOptions.OnError.SkipFile,
			}
		}
	}
	return opts
}

func (r *DropStageRequest) toOpts() *DropStageOptions {
	opts := &DropStageOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeStageRequest) toOpts() *DescribeStageOptions {
	opts := &DescribeStageOptions{
		name: r.name,
	}
	return opts
}

func (r stageDescRow) convert() *StageProperty {
	// TODO: Mapping
	return &StageProperty{}
}

func (r *ShowStageRequest) toOpts() *ShowStageOptions {
	opts := &ShowStageOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r stageShowRow) convert() *Stage {
	// TODO: Mapping
	return &Stage{}
}
