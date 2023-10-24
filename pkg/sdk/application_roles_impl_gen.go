package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ ApplicationRoles = (*applicationRoles)(nil)

type applicationRoles struct {
	client *Client
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

func (v *applicationRoles) ShowByID(ctx context.Context, request *ShowByIDApplicationRoleRequest) (*ApplicationRole, error) {
	appRoles, err := v.client.ApplicationRoles.Show(ctx, NewShowApplicationRoleRequest().WithApplicationName(request.ApplicationName))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(appRoles, func(role ApplicationRole) bool { return role.Name == request.name.Name() })
}

func (r *ShowApplicationRoleRequest) toOpts() *ShowApplicationRoleOptions {
	opts := &ShowApplicationRoleOptions{
		ApplicationName: r.ApplicationName,
	}
	if r.Limit != nil {
		opts.Limit = &LimitFrom{
			Rows: r.Limit.Rows,
			From: r.Limit.From,
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
		OwnerRoleType: r.OwnerRoleType,
	}
}
