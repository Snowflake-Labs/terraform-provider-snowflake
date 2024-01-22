package sdk

import (
	"context"

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
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegesToDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *GrantPrivilegesToDatabaseRoleOptions) error {
	if opts == nil {
		opts = &GrantPrivilegesToDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegesFromDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *RevokePrivilegesFromDatabaseRoleOptions) error {
	if opts == nil {
		opts = &RevokePrivilegesFromDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegeToShare(ctx context.Context, privilege ObjectPrivilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error {
	opts := &grantPrivilegeToShareOptions{
		privilege: privilege,
		On:        on,
		to:        to,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegeFromShare(ctx context.Context, privilege ObjectPrivilege, on *RevokePrivilegeFromShareOn, id AccountObjectIdentifier) error {
	opts := &revokePrivilegeFromShareOptions{
		privilege: privilege,
		On:        on,
		from:      id,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantOwnership(ctx context.Context, on OwnershipGrantOn, to OwnershipGrantTo, opts *GrantOwnershipOptions) error {
	if opts == nil {
		opts = &GrantOwnershipOptions{}
	}
	opts.On = on
	opts.To = to
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
