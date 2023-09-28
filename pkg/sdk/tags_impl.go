package sdk

import "context"

var _ Tags = (*tags)(nil)

type tags struct {
	client *Client
}

func (v *tags) Create(ctx context.Context, request *CreateTagRequest) error {
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

func (v *tags) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Tag, error) {
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
	opts := createTagOptions{
		OrReplace:   Bool(s.orReplace),
		IfNotExists: Bool(s.ifNotExists),
		name:        s.name,
	}
	if s.comment != nil {
		opts.Comment = s.comment
	}
	if s.allowedValues != nil {
		opts.AllowedValues = s.allowedValues
	}
	return &opts
}

func (s *ShowTagRequest) toOpts() *showTagOptions {
	opts := showTagOptions{}
	if s.like != nil {
		opts.Like = s.like
	}
	if s.in != nil {
		opts.In = s.in
	}
	return &opts
}

func (s *DropTagRequest) toOpts() *dropTagOptions {
	opts := dropTagOptions{
		IfNotExists: Bool(s.ifNotExists),
		name:        s.name,
	}
	return &opts
}

func (s *UndropTagRequest) toOpts() *undropTagOptions {
	opts := undropTagOptions{
		name: s.name,
	}
	return &opts
}
