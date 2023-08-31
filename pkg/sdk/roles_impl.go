package sdk

import (
	"context"
)

var _ Roles = (*roles)(nil)

type roles struct {
	client *Client
}

func (v *roles) Create(ctx context.Context, req *CreateRoleRequest) error {
	return validateAndExec(v.client, ctx, req.ToOpts())
}

func (v *roles) Alter(ctx context.Context, req *AlterRoleRequest) error {
	return validateAndExec(v.client, ctx, req.ToOpts())
}

func (v *roles) Drop(ctx context.Context, req *DropRoleRequest) error {
	return validateAndExec(v.client, ctx, req.ToOpts())
}

func (v *roles) Show(ctx context.Context, req *ShowRoleRequest) ([]Role, error) {
	rows, err := validateAndQuery[roleDBRow](v.client, ctx, req.ToOpts())
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

func (v *roles) Grant(ctx context.Context, req *GrantRoleRequest) error {
	return validateAndExec(v.client, ctx, req.ToOpts())
}

func (v *roles) Revoke(ctx context.Context, req *RevokeRoleRequest) error {
	return validateAndExec(v.client, ctx, req.ToOpts())
}

func (v *roles) Use(ctx context.Context, req *UseRoleRequest) error {
	return v.client.Sessions.UseRole(ctx, req.id)
}

func (v *roles) UseSecondary(ctx context.Context, req *UseSecondaryRolesRequest) error {
	return v.client.Sessions.UseSecondaryRoles(ctx, req.option)
}
