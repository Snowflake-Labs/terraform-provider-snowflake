package sdk

import (
	"context"
)

var _ Roles = (*roles)(nil)

type roles struct {
	client *Client
}

func (v *roles) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateRoleOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *roles) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterRoleOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *roles) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropRoleOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *roles) Show(ctx context.Context, opts *ShowRoleOptions) ([]Role, error) {
	opts = createIfNil(opts)
	rows, err := validateAndQuery[roleDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	roles := make([]Role, len(rows))
	for i, row := range rows {
		roles[i] = row.toRole()
	}

	return roles, nil
}

func (v *roles) ShowByID(ctx context.Context, req *ShowRoleByIdRequest) (*Role, error) {
	roles, err := v.client.Roles.Show(ctx, NewShowRoleRequest().WithLike(NewLikeRequest(req.id.Name())))
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.ID() == req.id {
			return &role, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

func (v *roles) Grant(ctx context.Context, id AccountObjectIdentifier, opts *GrantRoleOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *roles) Revoke(ctx context.Context, id AccountObjectIdentifier, opts *RevokeRoleOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *roles) Use(ctx context.Context, id AccountObjectIdentifier) error {
	return v.client.Sessions.UseRole(ctx, id)
}

func (v *roles) UseSecondary(ctx context.Context, opt SecondaryRoleOption) error {
	return v.client.Sessions.UseSecondaryRoles(ctx, opt)
}
