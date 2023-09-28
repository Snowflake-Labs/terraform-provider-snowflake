package example

import "context"

var _ DatabaseRoles = (*databaseRoles)(nil)

type databaseRoles struct {
	client *Client
}

func (v *databaseRoles) Create(ctx context.Context, request *CreateDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Alter(ctx context.Context, request *AlterDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateDatabaseRoleRequest) toOpts() *CreateDatabaseRoleOptions {
	opts := &CreateDatabaseRoleOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
	}

	return opts
}

func (r *AlterDatabaseRoleRequest) toOpts() *AlterDatabaseRoleOptions {
	opts := &AlterDatabaseRoleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Rename != nil {
		opts.Rename = &DatabaseRoleRename{
			Name: r.Rename.Name,
		}
	}
	if r.Set != nil {
		opts.Set = &DatabaseRoleSet{
			Comment: r.Set.Comment,
		}
		if r.Set.NestedThirdLevel != nil {
			opts.Set.NestedThirdLevel = &NestedThirdLevel{
				Field: r.Unset.NestedThirdLevel.Field,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &DatabaseRoleUnset{
			Comment: r.Unset.Comment,
		}
		if r.Unset.NestedThirdLevel != nil {
			opts.Unset.NestedThirdLevel = &NestedThirdLevel{
				Field: r.Unset.NestedThirdLevel.Field,
			}
		}
	}

	return opts
}
