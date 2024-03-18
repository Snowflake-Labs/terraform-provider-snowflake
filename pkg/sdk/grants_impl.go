package sdk

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
)

var _ Grants = (*grants)(nil)

type grants struct {
	client *Client
}

func (v *grants) GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error {
	logging.DebugLogger.Printf("[DEBUG] Grant privileges to account role")
	if opts == nil {
		opts = &GrantPrivilegesToAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role
	logging.DebugLogger.Printf("[DEBUG] Grant privileges to account role: opts %+v", opts)

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantPrivilegesToAccountRole(
					ctx,
					privileges,
					&AccountRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       NewSchemaObjectIdentifier(pipe.DatabaseName, pipe.SchemaName, pipe.Name),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error {
	logging.DebugLogger.Printf("[DEBUG] Revoke privileges from account role")
	if opts == nil {
		opts = &RevokePrivilegesFromAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role
	logging.DebugLogger.Printf("[DEBUG] Revoke privileges from account role: opts %+v", opts)

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.RevokePrivilegesFromAccountRole(
					ctx,
					privileges,
					&AccountRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       NewSchemaObjectIdentifier(pipe.DatabaseName, pipe.SchemaName, pipe.Name),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegesToDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *GrantPrivilegesToDatabaseRoleOptions) error {
	if opts == nil {
		opts = &GrantPrivilegesToDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantPrivilegesToDatabaseRole(
					ctx,
					privileges,
					&DatabaseRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       NewSchemaObjectIdentifier(pipe.DatabaseName, pipe.SchemaName, pipe.Name),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegesFromDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *RevokePrivilegesFromDatabaseRoleOptions) error {
	if opts == nil {
		opts = &RevokePrivilegesFromDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.RevokePrivilegesFromDatabaseRole(
					ctx,
					privileges,
					&DatabaseRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       NewSchemaObjectIdentifier(pipe.DatabaseName, pipe.SchemaName, pipe.Name),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegeToShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, to AccountObjectIdentifier) error {
	opts := &grantPrivilegeToShareOptions{
		privileges: privileges,
		On:         on,
		to:         to,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegeFromShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, id AccountObjectIdentifier) error {
	opts := &revokePrivilegeFromShareOptions{
		privileges: privileges,
		On:         on,
		from:       id,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantOwnership(ctx context.Context, on OwnershipGrantOn, to OwnershipGrantTo, opts *GrantOwnershipOptions) error {
	if opts == nil {
		opts = &GrantOwnershipOptions{}
	}
	opts.On = on
	opts.To = to

	// TODO: Suspend / Pause Pipes / Tasks before granting ownership (and restore the state before the transfer)
	if on.Object != nil && on.Object.ObjectType == ObjectTypePipe {

	}
	if on.Object != nil && on.Object.ObjectType == ObjectTypeTask {

	}

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on.All != nil && on.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.All.InDatabase,
			on.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantOwnership(
					ctx,
					OwnershipGrantOn{
						Object: &Object{
							ObjectType: ObjectTypePipe,
							Name:       NewSchemaObjectIdentifier(pipe.DatabaseName, pipe.SchemaName, pipe.Name),
						},
					},
					to,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) Show(ctx context.Context, opts *ShowGrantOptions) ([]Grant, error) {
	logging.DebugLogger.Printf("[DEBUG] Show grants")
	if opts == nil {
		opts = &ShowGrantOptions{}
	}

	logging.DebugLogger.Printf("[DEBUG] Show grants: opts %+v", opts)
	dbRows, err := validateAndQuery[grantRow](v.client, ctx, opts)
	logging.DebugLogger.Printf("[DEBUG] Show grants: query finished err = %v", err)
	if err != nil {
		return nil, err
	}
	logging.DebugLogger.Printf("[DEBUG] Show grants: converting rows")
	resultList := convertRows[grantRow, Grant](dbRows)
	logging.DebugLogger.Printf("[DEBUG] Show grants: rows converted")
	return resultList, nil
}

func (v *grants) runOnAllPipes(ctx context.Context, inDatabase *AccountObjectIdentifier, inSchema *DatabaseObjectIdentifier, command func(Pipe) error) error {
	var in *In
	switch {
	case inDatabase != nil:
		in = &In{
			Database: *inDatabase,
		}
	case inSchema != nil:
		in = &In{
			Schema: *inSchema,
		}
	}

	pipes, err := v.client.Pipes.Show(ctx, &ShowPipeOptions{In: in})
	if err != nil {
		return err
	}

	var errs []error
	for _, pipe := range pipes {
		if err := command(pipe); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
