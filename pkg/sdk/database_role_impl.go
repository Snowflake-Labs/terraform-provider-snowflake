package sdk

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

func (v *databaseRoles) Drop(ctx context.Context, request *DropDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Show(ctx context.Context, request *ShowDatabaseRoleRequest) ([]*DatabaseRole, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[databaseRoleDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList := make([]*DatabaseRole, len(*dbRows))
	for i, row := range *dbRows {
		resultList[i] = row.toDatabaseRole()
	}

	return resultList, nil
}

func (v *databaseRoles) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error) {
	request := NewShowDatabaseRoleRequest(NewAccountObjectIdentifier(id.DatabaseName())).WithLike(id.Name())
	databaseRoles, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, databaseRole := range databaseRoles {
		if databaseRole.Name == id.Name() {
			return databaseRole, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

func (s *CreateDatabaseRoleRequest) toOpts() *CreateDatabaseRoleOptions {
	return &CreateDatabaseRoleOptions{
		OrReplace:   Bool(s.orReplace),
		IfNotExists: Bool(s.ifNotExists),
		name:        s.name,
		Comment:     s.comment,
	}
}

func (s *AlterDatabaseRoleRequest) toOpts() *AlterDatabaseRoleOptions {
	opts := AlterDatabaseRoleOptions{
		IfExists: Bool(s.ifExists),
		name:     s.name,
	}
	if s.rename != nil {
		opts.Rename = &DatabaseRoleRename{s.rename.name}
	}
	if s.set != nil {
		opts.Set = &DatabaseRoleSet{s.set.comment}
	}
	if s.unset != nil {
		opts.Unset = &DatabaseRoleUnset{true}
	}
	return &opts
}

func (s *DropDatabaseRoleRequest) toOpts() *DropDatabaseRoleOptions {
	return &DropDatabaseRoleOptions{
		IfExists: Bool(s.ifExists),
		name:     s.name,
	}
}

func (s *ShowDatabaseRoleRequest) toOpts() *ShowDatabaseRoleOptions {
	return &ShowDatabaseRoleOptions{
		Like:     s.like,
		Database: s.database,
	}
}
