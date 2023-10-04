package sdk

import (
	"context"
)

var _ ApplicationRoles = (*applicationRoles)(nil)

type applicationRoles struct {
	client *Client
}

func (v *applicationRoles) Create(ctx context.Context, request *CreateApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Alter(ctx context.Context, request *AlterApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Drop(ctx context.Context, request *DropApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Show(ctx context.Context, request *ShowApplicationRoleRequest) ([]ApplicationRole, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[applicationRoleDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[applicationRoleDbRow, ApplicationRole](dbRows)
	return resultList, nil
}

func (v *applicationRoles) Grant(ctx context.Context, request *GrantApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Revoke(ctx context.Context, request *RevokeApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateApplicationRoleRequest) toOpts() *CreateApplicationRoleOptions {
	opts := &CreateApplicationRoleOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
	}
	return opts
}

func (r *AlterApplicationRoleRequest) toOpts() *AlterApplicationRoleOptions {
	opts := &AlterApplicationRoleOptions{
		IfExists:     r.IfExists,
		name:         r.name,
		RenameTo:     r.RenameTo,
		SetComment:   r.SetComment,
		UnsetComment: r.UnsetComment,
	}
	return opts
}

func (r *DropApplicationRoleRequest) toOpts() *DropApplicationRoleOptions {
	opts := &DropApplicationRoleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowApplicationRoleRequest) toOpts() *ShowApplicationRoleOptions {
	opts := &ShowApplicationRoleOptions{
		ApplicationName: r.ApplicationName,
	}
	if r.LimitFrom != nil {
		opts.LimitFrom = &LimitFromApplicationRole{
			Rows: r.LimitFrom.Rows,
			From: r.LimitFrom.From,
		}
	}
	return opts
}

func (r applicationRoleDbRow) convert() *ApplicationRole {
	return &ApplicationRole{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		Owner:         r.Owner,
		Comment:       r.Comment,
		OwnerRoleTYpe: r.OwnerRoleType,
	}
}

func (r *GrantApplicationRoleRequest) toOpts() *GrantApplicationRoleOptions {
	return &GrantApplicationRoleOptions{
		name: r.name,
		GrantTo: ApplicationGrantOptions{
			ParentRole:      r.GrantTo.ParentRole,
			ApplicationRole: r.GrantTo.ApplicationRole,
			Application:     r.GrantTo.Application,
		},
	}
}

func (r *RevokeApplicationRoleRequest) toOpts() *RevokeApplicationRoleOptions {
	return &RevokeApplicationRoleOptions{
		name: r.name,
		RevokeFrom: ApplicationGrantOptions{
			ParentRole:      r.RevokeFrom.ParentRole,
			ApplicationRole: r.RevokeFrom.ApplicationRole,
			Application:     r.RevokeFrom.Application,
		},
	}
}
