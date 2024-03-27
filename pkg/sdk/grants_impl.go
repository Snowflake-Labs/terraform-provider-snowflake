package sdk

import (
	"context"
	"errors"
	"slices"

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

	if on.Object != nil && on.Object.ObjectType == ObjectTypePipe {
		return v.grantOwnershipOnPipe(ctx, on.Object.Name.(SchemaObjectIdentifier), opts)
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

// TODO (remove ?): it was only for me to grasp what has to be done in this function, but may also leave it to guide others
// grant ownership on pipe sequence (at worst 9 operations)
// - get current role (needed to grant operate privilege later on)
// - grant operate on pipe if not granted (it will error our otherwise)
// - get pipe status (running or paused)
// - if pipe is running, stop the pipe
// - revoke operate (it may affect grant ownership call)
// - grant ownership
// - if pipe was previously running
//   - grant operate on pipe to current role
//   - unpause the pipe with system function
//   - revoke operate on pipe from current role
func (v *grants) grantOwnershipOnPipe(ctx context.Context, pipeId SchemaObjectIdentifier, opts *GrantOwnershipOptions) error {
	// To be able to call ALTER on a pipe to stop its execution,
	// the current role has to be either the owner of this pipe or be granted with OPERATE privilege.
	// The code below makes certain checks and takes care of making sure that the current role is privileged enough to act on the pipe and safely grant ownership.

	currentRole, err := v.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}
	currentRoleName := NewAccountObjectIdentifier(currentRole)

	currentGrants, err := v.client.Grants.Show(ctx, &ShowGrantOptions{
		On: &ShowGrantsOn{
			Object: &Object{
				ObjectType: ObjectTypePipe,
				Name:       pipeId,
			},
		},
	})
	if err != nil {
		return err
	}

	canOperateOnPipe := slices.ContainsFunc(currentGrants, func(grant Grant) bool {
		return grant.GranteeName == currentRoleName && (grant.Privilege == "OWNERSHIP" || grant.Privilege == SchemaObjectPrivilegeOperate.String())
	})

	var revokeOperate func() error
	if !canOperateOnPipe {
		revokeOperate, err = v.grantOperateOnPipeTemporarily(ctx, pipeId, currentRoleName)
		if err != nil {
			return err
		}
	}

	originalPipeExecutionState, err := v.client.SystemFunctions.PipeStatus(pipeId)
	if err != nil {
		return err
	}

	if originalPipeExecutionState == RunningPipeExecutionState {
		if err := v.client.Pipes.Alter(ctx, pipeId, &AlterPipeOptions{
			Set: &PipeSet{
				PipeExecutionPaused: Bool(true),
			},
		}); err != nil {
			return err
		}
	}

	if !canOperateOnPipe {
		if err := revokeOperate(); err != nil {
			return err
		}
	}

	if err := validateAndExec(v.client, ctx, opts); err != nil {
		return err
	}

	if originalPipeExecutionState == RunningPipeExecutionState {
		// We cannot resume right away and "normally" through ALTER, because:
		// 1. Insufficient privileges (the current role has to be granted with at least OPERATE privilege).
		// 2. Snowflake throws an error whenever pipes changes its owner, and you try to unpause it. The error suggests using system function to forcefully resume the pipe.
		revokeOperate, err := v.grantOperateOnPipeTemporarily(ctx, pipeId, currentRoleName)
		if err != nil {
			return err
		}

		// TODO: check if options need to be passed
		if err := v.client.SystemFunctions.PipeForceResume(pipeId, nil); err != nil {
			return err
		}

		if err := revokeOperate(); err != nil {
			return err
		}
	}

	return nil
}

func (v *grants) grantOperateOnPipeTemporarily(ctx context.Context, pipeId SchemaObjectIdentifier, currentRole AccountObjectIdentifier) (func() error, error) {
	return v.grantTemporarily(
		ctx,
		&AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{
				SchemaObjectPrivilegeOperate,
			},
		},
		&AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				SchemaObject: &Object{
					ObjectType: ObjectTypePipe,
					Name:       pipeId,
				},
			},
		},
		currentRole,
	)
}

func (v *grants) grantTemporarily(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, accountRoleName AccountObjectIdentifier) (func() error, error) {
	return func() error {
			return v.client.Grants.RevokePrivilegesFromAccountRole(
				ctx,
				privileges,
				on,
				accountRoleName,
				new(RevokePrivilegesFromAccountRoleOptions),
			)
		}, v.client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			privileges,
			on,
			accountRoleName,
			new(GrantPrivilegesToAccountRoleOptions),
		)
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

	return runOnAll(pipes, command)
}

func runOnAll[T any](collection []T, command func(T) error) error {
	var errs []error
	for _, element := range collection {
		if err := command(element); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
