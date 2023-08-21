package sdk

import "context"

var _ Grants = (*grants)(nil)

type grants struct {
	client *Client
}

func (v *grants) GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error {
	if opts == nil {
		opts = &GrantPrivilegesToAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return nil
}

func (v *grants) RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error {
	if opts == nil {
		opts = &RevokePrivilegesFromAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return nil
}

func (v *grants) GrantPrivilegeToShare(ctx context.Context, privilege ObjectPrivilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error {
	opts := &grantPrivilegeToShareOptions{
		privilege: privilege,
		On:        on,
		to:        to,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *grants) RevokePrivilegeFromShare(ctx context.Context, privilege ObjectPrivilege, on *RevokePrivilegeFromShareOn, id AccountObjectIdentifier) error {
	opts := &revokePrivilegeFromShareOptions{
		privilege: privilege,
		On:        on,
		from:      id,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *grants) Show(ctx context.Context, opts *ShowGrantOptions) ([]*Grant, error) {
	if opts == nil {
		opts = &ShowGrantOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []grantRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	grants := make([]*Grant, 0, len(rows))
	for _, row := range rows {
		grant, err := row.toGrant()
		if err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}
	return grants, nil
}
