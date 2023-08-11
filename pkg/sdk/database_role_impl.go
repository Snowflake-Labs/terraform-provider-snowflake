package sdk

import "context"

var _ DatabaseRoles = (*databaseRoles)(nil)

type databaseRoles struct {
	client *Client
}

func (v *databaseRoles) Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateDatabaseRoleOptions) error {
	opts = createIfNil[CreateDatabaseRoleOptions](opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterDatabaseRoleOptions) error {
	opts = createIfNil[AlterDatabaseRoleOptions](opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Drop(ctx context.Context, id DatabaseObjectIdentifier) error {
	opts := &DropDatabaseRoleOptions{
		name: id,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Show(ctx context.Context, opts *ShowDatabaseRoleOptions) ([]*DatabaseRole, error) {
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
	databaseRoles, err := v.Show(ctx, &ShowDatabaseRoleOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		database: NewAccountObjectIdentifier(id.DatabaseName()),
	})
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
