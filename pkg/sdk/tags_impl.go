package sdk

import (
	"context"
)

var _ Tags = (*tags)(nil)

type tags struct {
	client *Client
}

func (v *tags) Create(ctx context.Context, request *CreateTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Alter(ctx context.Context, request *AlterTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Show(ctx context.Context, request *ShowTagRequest) ([]Tag, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[tagRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	result := convertRows[tagRow, Tag](rows)
	return result, nil
}

func (v *tags) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Tag, error) {
	request := NewShowTagRequest().WithLike(id.Name())
	tags, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return findOne(tags, func(r Tag) bool { return r.Name == id.Name() })
}

func (v *tags) Drop(ctx context.Context, request *DropTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tags) Undrop(ctx context.Context, request *UndropTagRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (s *CreateTagRequest) toOpts() *createTagOptions {
	return &createTagOptions{
		OrReplace:     Bool(s.orReplace),
		IfNotExists:   Bool(s.ifNotExists),
		name:          s.name,
		Comment:       s.comment,
		AllowedValues: s.allowedValues,
	}
}

func (s *AlterTagRequest) toOpts() *alterTagOptions {
	return &alterTagOptions{
		name:   s.name,
		Add:    s.add,
		Drop:   s.drop,
		Set:    s.set,
		Unset:  s.unset,
		Rename: s.rename,
	}
}

func (s *ShowTagRequest) toOpts() *showTagOptions {
	return &showTagOptions{
		Like: s.like,
		In:   s.in,
	}
}

func (s *DropTagRequest) toOpts() *dropTagOptions {
	return &dropTagOptions{
		IfExists: Bool(s.ifExists),
		name:     s.name,
	}
}

func (s *UndropTagRequest) toOpts() *undropTagOptions {
	return &undropTagOptions{
		name: s.name,
	}
}
