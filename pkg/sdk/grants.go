package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Grants interface {
	// Account Roles (grant)
	GrantPrivilegesToAccountRole(ctx context.Context, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error
	// GrantAccountObjectPrivilegesToAccountRole(ctx context.Context, opts *GrantAccountObjectPrivilegesToAccountRoleOptions) error
	// GrantSchemaPrivilegesToAccountRole(ctx context.Context, privileges []SchemaPrivilege, opts *GrantSchemaPrivilegesToRoleOptions) error
	// GrantSchemaObjectPrivilegesToAccountRole(ctx context.Context, privileges []SchemaObjectPrivilege, opts *GrantSchemaObjectPrivilegesToRoleOptions) error

	// Account Roles (revoke)
	RevokeGlobalPrivilegesFromAccountRole(ctx context.Context, opts *RevokeGlobalPrivilegesFromAccountRoleOptions) error
	RevokeAccountObjectPrivilegesFromAccountRole(ctx context.Context, opts *RevokeAccountObjectPrivilegesFromAccountRoleOptions) error

	// RevokeGlobalPrivilegesFromAccountRole(ctx context.Context, privileges []GlobalPrivilege, opts *RevokeGlobalPrivilegesFromRoleOptions) error
	// RevokeGlobalPrivilegesFromAccountRole(ctx context.Context, privileges []GlobalPrivilege, opts *RevokeGlobalPrivilegesFromRoleOptions) error
	// RevokeAccountObjectPrivilegesFromAccountRole(ctx context.Context, privileges []GlobalPrivilege, opts *GrantAccountObjectPrivilegesToRoleOptions) error
	// RevokeSchemaPrivilegesFromAccountRole(ctx context.Context, privileges []SchemaPrivilege, opts *RevokeSchemaPrivilegesFromRoleOptions) error
	// RevokeSchemaObjectPrivilegesFromAccountRole(ctx context.Context, privileges []SchemaObjectPrivilege, opts *RevokeSchemaObjectPrivilegesFromRoleOptions) error
	// Database Roles (grant)
	// RevokePrivilegesFromAccountRole(ctx context.Context, objectPrivilege ObjectPrivilege, on *RevokePrivilegeFromRoleOn, from RoleObjectIdentifier) error
	GrantPrivilegeToShare(ctx context.Context, privilege ObjectPrivilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error
	RevokePrivilegeFromShare(ctx context.Context, privilege ObjectPrivilege, on *RevokePrivilegeFromShareOn, from AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowGrantOptions) ([]*Grant, error)
}

var _ Grants = (*grants)(nil)

type grants struct {
	client *Client
}

type Grant struct {
	CreatedOn   time.Time
	Privilege   string
	GrantedOn   ObjectType
	Name        ObjectIdentifier
	GrantedTo   ObjectType
	GranteeName AccountObjectIdentifier
	GrantOption bool
	GrantedBy   AccountObjectIdentifier
}

func (v *Grant) ID() ObjectIdentifier {
	return v.Name
}

type grantRow struct {
	CreatedOn   time.Time `db:"created_on"`
	Privilege   string    `db:"privilege"`
	GrantedOn   string    `db:"granted_on"`
	Name        string    `db:"name"`
	GrantedTo   string    `db:"granted_to"`
	GranteeName string    `db:"grantee_name"`
	GrantOption bool      `db:"grant_option"`
	GrantedBy   string    `db:"granted_by"`
}

func (row *grantRow) toGrant() (*Grant, error) {
	grantedTo := ObjectType(row.GrantedTo)
	granteeName := NewAccountObjectIdentifier(row.GranteeName)
	if grantedTo == ObjectTypeShare {
		parts := strings.Split(row.GranteeName, ".")
		name := strings.Join(parts[1:], ".")
		granteeName = NewAccountObjectIdentifier(name)
	}
	grant := &Grant{
		CreatedOn:   row.CreatedOn,
		Privilege:   row.Privilege,
		GrantedOn:   ObjectType(row.GrantedOn),
		GrantedTo:   grantedTo,
		Name:        NewAccountObjectIdentifier(strings.Trim(row.Name, "\"")),
		GranteeName: granteeName,
		GrantOption: row.GrantOption,
		GrantedBy:   NewAccountObjectIdentifier(row.GrantedBy),
	}
	return grant, nil
}

type GrantPrivilegesToAccountRoleOptions struct {
	grant           bool                        `ddl:"static" sql:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	Privileges      *AccountRoleGrantPrivileges `ddl:"-"`
	AllPrivileges   *bool                       `ddl:"keyword" sql:"ALL PRIVILEGES"`
	On              *GrantOn                    `ddl:"-"`
	toRole          AccountObjectIdentifier     `ddl:"identifier" sql:"TO ROLE"`
	WithGrantOption *bool                       `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

type AccountRoleGrantPrivileges struct {
	GlobalPrivileges        []GlobalPrivilege        `ddl:"-"`
	AccountObjectPrivileges []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges        []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges  []SchemaObjectPrivilege  `ddl:"-"`
}

type GrantOn struct {
	Account       *bool                 `ddl:"keyword" sql:"ACCOUNT"`
	AccountObject *GrantOnAccountObject `ddl:"-"`
	Schema        *GrantOnSchema        `ddl:"-"`
	SchemaObject  *GrantOnSchemaObject  `ddl:"-"`
}

type GrantOnAccountObject struct {
	User             *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	ResourceMonitor  *AccountObjectIdentifier `ddl:"identifier" sql:"RESOURCE MONITOR"`
	Warehouse        *AccountObjectIdentifier `ddl:"identifier" sql:"WAREHOUSE"`
	Database         *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Integration      *AccountObjectIdentifier `ddl:"identifier" sql:"INTEGRATION"`
	FailoverGroup    *AccountObjectIdentifier `ddl:"identifier" sql:"FAILOVER GROUP"`
	ReplicationGroup *AccountObjectIdentifier `ddl:"identifier" sql:"REPLICATION GROUP"`
}

type GrantOnSchema struct {
	Schema                  *SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	AllSchemasInDatabase    *AccountObjectIdentifier `ddl:"identifier" sql:"ALL SCHEMAS IN DATABASE"`
	FutureSchemasInDatabase *AccountObjectIdentifier `ddl:"identifier" sql:"FUTURE SCHEMAS IN DATABASE"`
}

type GrantOnSchemaObject struct {
	Object *Object                   `ddl:"-"`
	AllIn  *GrantOnSchemaObjectAllIn `ddl:"-"`
}

type GrantOnSchemaObject struct {
	PluralObjectType *PluralObjectType        `ddl:"keyword" sql:"ALL"`
	InDatabase       *AccountObjectIdentifier `ddl:"identifier" sql:"IN DATABASE"`
	InSchema         *SchemaIdentifier        `ddl:"identifier" sql:"IN SCHEMA"`
}

func (opts *GrantGlobalPrivilegesToAccountRoleOptions) validate() error {
	if !exactlyOneValueSet(opts.Privileges, opts.AllPrivileges) {
		return fmt.Errorf("only one of privileges or allPrivileges can be set")
	}
	if !validObjectidentifier(opts.toRole) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *grants) GrantGlobalPrivilegesToAccountRole(ctx context.Context, opts *GrantGlobalPrivilegesToAccountRoleOptions) error {
	if opts == nil {
		return fmt.Errorf("opts cannot be nil")
	}
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

type RevokeGlobalPrivilegesFromAccountRoleOptions struct {
	revoke         bool                    `ddl:"static" sql:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	GrantOptionFor *bool                   `ddl:"keyword" sql:"GRANT OPTION FOR"`
	Privileges     []GlobalPrivilege       `ddl:"-"`
	AllPrivileges  bool                    `ddl:"keyword" sql:"ALL PRIVILEGES"`
	onAccount      bool                    `ddl:"static" sql:"ON ACCOUNT"` //lint:ignore U1000 This is used in the ddl tag
	fromRole       AccountObjectIdentifier `ddl:"identifier" sql:"FROM ROLE"`
	Restrict       *bool                   `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                   `ddl:"keyword" sql:"CASCADE"`
}

func (opts *RevokeGlobalPrivilegesFromAccountRoleOptions) validate() error {
	if !exactlyOneValueSet(opts.Privileges, opts.AllPrivileges) {
		return fmt.Errorf("only one of privileges or allPrivileges can be set")
	}
	if !validObjectidentifier(opts.fromRole) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *grants) RevokeGlobalPrivilegesFromAccountRole(ctx context.Context, opts *RevokeGlobalPrivilegesFromAccountRoleOptions) error {
	if opts == nil {
		return fmt.Errorf("opts cannot be nil")
	}
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

type GrantAccountObjectPrivilegesToAccountRoleOptions struct {
	grant           bool                     `ddl:"static" sql:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	Privileges      []AccountObjectPrivilege `ddl:"-"`
	AllPrivileges   bool                     `ddl:"keyword" sql:"ALL PRIVILEGES"`
	on              *GrantOnAccountObject    `ddl:"-" sql:"ON"`
	toRole          AccountObjectIdentifier  `ddl:"identifier" sql:"TO ROLE"`
	WithGrantOption bool                     `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

func (opts *GrantAccountObjectPrivilegesToAccountRoleOptions) validate() error {
	if !exactlyOneValueSet(opts.Privileges, opts.AllPrivileges) {
		return fmt.Errorf("only one of privileges or allPrivileges can be set")
	}
	if !validObjectidentifier(opts.toRole) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on cannot be nil")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	return nil
}

type GrantOnAccountObject struct {
	User             *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	ResourceMonitor  *AccountObjectIdentifier `ddl:"identifier" sql:"RESOURCE MONITOR"`
	Warehouse        *AccountObjectIdentifier `ddl:"identifier" sql:"WAREHOUSE"`
	Database         *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Integration      *AccountObjectIdentifier `ddl:"identifier" sql:"INTEGRATION"`
	FailoverGroup    *AccountObjectIdentifier `ddl:"identifier" sql:"FAILOVER GROUP"`
	ReplicationGroup *AccountObjectIdentifier `ddl:"identifier" sql:"REPLICATION GROUP"`
}

func (opts *GrantOnAccountObject) validate() error {
	if !exactlyOneValueSet(opts.User, opts.ResourceMonitor, opts.Warehouse, opts.Database, opts.Integration, opts.FailoverGroup, opts.ReplicationGroup) {
		return fmt.Errorf("only one of user, resourceMonitor, warehouse, database, integration, failoverGroup, replicationGroup can be set")
	}
	return nil
}

func (v *grants) GrantAccountObjectPrivilegesToAccountRole(ctx context.Context, opts *GrantAccountObjectPrivilegesToAccountRoleOptions) error {
	if opts == nil {
		return fmt.Errorf("opts cannot be nil")
	}
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

type RevokeAccountObjectPrivilegesFromAccountRoleOptions struct {
	revoke         bool                     `ddl:"static" sql:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	GrantOptionFor *bool                    `ddl:"keyword" sql:"GRANT OPTION FOR"`
	Privileges     []AccountObjectPrivilege `ddl:"-"`
	AllPrivileges  bool                     `ddl:"keyword" sql:"ALL PRIVILEGES"`
	on             *GrantOnAccountObject    `ddl:"-" sql:"ON"`
	fromRole       AccountObjectIdentifier  `ddl:"identifier" sql:"FROM ROLE"`
	Restrict       *bool                    `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                    `ddl:"keyword" sql:"CASCADE"`
}

func (opts *RevokeAccountObjectPrivilegesFromAccountRoleOptions) validate() error {
	if !exactlyOneValueSet(opts.Privileges, opts.AllPrivileges) {
		return fmt.Errorf("only one of privileges or allPrivileges can be set")
	}
	if !validObjectidentifier(opts.fromRole) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *grants) RevokeAccountObjectPrivilegesFromAccountRole(ctx context.Context, opts *RevokeAccountObjectPrivilegesFromAccountRoleOptions) error {
	if opts == nil {
		return fmt.Errorf("opts cannot be nil")
	}
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

type grantPrivilegeToShareOptions struct {
	grant     bool                     `ddl:"static" sql:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	privilege ObjectPrivilege          `ddl:"keyword"`
	On        *GrantPrivilegeToShareOn `ddl:"keyword" sql:"ON"`
	to        AccountObjectIdentifier  `ddl:"identifier" sql:"TO SHARE"`
}

func (opts *grantPrivilegeToShareOptions) validate() error {
	if !validObjectidentifier(opts.to) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		return fmt.Errorf("on and privilege are required")
	}
	if !exactlyOneValueSet(opts.On.Database, opts.On.Schema, opts.On.Function, opts.On.Table, opts.On.View) {
		return fmt.Errorf("only one of database, schema, function, table, or view can be set")
	}
	return nil
}

type GrantPrivilegeToShareOn struct {
	Database AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Schema   SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	Function SchemaObjectIdentifier  `ddl:"identifier" sql:"FUNCTION"`
	Table    *OnTable                `ddl:"-"`
	View     SchemaObjectIdentifier  `ddl:"identifier" sql:"VIEW"`
}

func (v *GrantPrivilegeToShareOn) validate() error {
	if !exactlyOneValueSet(v.Database, v.Schema, v.Function, v.Table, v.View) {
		return fmt.Errorf("only one of database, schema, function, table, or view can be set")
	}
	if valueSet(v.Table) {
		if err := v.Table.validate(); err != nil {
			return err
		}
	}
	return nil
}

type OnTable struct {
	Name        SchemaObjectIdentifier `ddl:"identifier" sql:"TABLE"`
	AllInSchema SchemaIdentifier       `ddl:"identifier" sql:"ALL TABLES IN SCHEMA"`
}

func (v *OnTable) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return fmt.Errorf("only one of name or allInSchema can be set")
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

type revokePrivilegeFromShareOptions struct {
	revoke    bool                        `ddl:"static" sql:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	privilege ObjectPrivilege             `ddl:"keyword"`
	On        *RevokePrivilegeFromShareOn `ddl:"keyword" sql:"ON"`
	from      AccountObjectIdentifier     `ddl:"identifier" sql:"FROM SHARE"`
}

func (opts *revokePrivilegeFromShareOptions) validate() error {
	if !validObjectidentifier(opts.from) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.On) || opts.privilege == "" {
		return fmt.Errorf("on and privilege are required")
	}
	if !exactlyOneValueSet(opts.On.Database, opts.On.Schema, opts.On.Table, opts.On.View) {
		return fmt.Errorf("only one of database, schema, function, table, or view can be set")
	}

	if err := opts.On.validate(); err != nil {
		return err
	}

	return nil
}

type RevokePrivilegeFromShareOn struct {
	Database AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Schema   SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	Table    *OnTable                `ddl:"-"`
	View     *OnView                 `ddl:"-"`
}

func (v *RevokePrivilegeFromShareOn) validate() error {
	if !exactlyOneValueSet(v.Database, v.Schema, v.Table, v.View) {
		return fmt.Errorf("only one of database, schema, table, or view can be set")
	}
	if valueSet(v.Table) {
		return v.Table.validate()
	}
	if valueSet(v.View) {
		return v.View.validate()
	}
	return nil
}

type OnView struct {
	Name        SchemaObjectIdentifier `ddl:"identifier" sql:"VIEW"`
	AllInSchema SchemaIdentifier       `ddl:"identifier" sql:"ALL VIEWS IN SCHEMA"`
}

func (v *OnView) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return fmt.Errorf("only one of name or allInSchema can be set")
	}
	return nil
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

type ShowGrantOptions struct {
	show   bool          `ddl:"static" sql:"SHOW"`   //lint:ignore U1000 This is used in the ddl tag
	grants bool          `ddl:"static" sql:"GRANTS"` //lint:ignore U1000 This is used in the ddl tag
	On     *ShowGrantsOn `ddl:"keyword" sql:"ON"`
	To     *ShowGrantsTo `ddl:"keyword" sql:"TO"`
	Of     *ShowGrantsOf `ddl:"keyword" sql:"OF"`
}

func (opts *ShowGrantOptions) validate() error {
	if everyValueNil(opts.On, opts.To, opts.Of) {
		return fmt.Errorf("at least one of on, to, or of is required")
	}
	if !exactlyOneValueSet(opts.On, opts.To, opts.Of) {
		return fmt.Errorf("only one of on, to, or of can be set")
	}
	return nil
}

type ShowGrantsOn struct {
	Account *bool `ddl:"keyword" sql:"ACCOUNT"`
	Object  *Object
}

type ShowGrantsTo struct {
	Role  AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
	User  AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	Share AccountObjectIdentifier `ddl:"identifier" sql:"SHARE"`
}

type ShowGrantsOf struct {
	Role  AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
	Share AccountObjectIdentifier `ddl:"identifier" sql:"SHARE"`
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
