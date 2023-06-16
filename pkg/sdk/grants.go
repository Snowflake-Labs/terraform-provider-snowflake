package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Grants interface {
	GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error
	RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error
	// todo: GrantPrivilegesToDatabaseRole and RevokePrivilegesFromDatabaseRole
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
	GrantOn     string    `db:"grant_on"`
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
	// true for future grants
	if row.GrantOn != "" {
		grant.GrantedOn = ObjectType(row.GrantOn)
	}
	return grant, nil
}

type GrantPrivilegesToAccountRoleOptions struct {
	grant           bool                        `ddl:"static" sql:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	privileges      *AccountRoleGrantPrivileges `ddl:"-"`
	on              *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole     AccountObjectIdentifier     `ddl:"identifier" sql:"TO ROLE"`
	WithGrantOption *bool                       `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

func (opts *GrantPrivilegesToAccountRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	return nil
}

type AccountRoleGrantPrivileges struct {
	GlobalPrivileges        []GlobalPrivilege        `ddl:"-"`
	AccountObjectPrivileges []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges        []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges  []SchemaObjectPrivilege  `ddl:"-"`
	AllPrivileges           *bool                    `ddl:"keyword" sql:"ALL PRIVILEGES"`
}

func (v *AccountRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.GlobalPrivileges, v.AccountObjectPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return fmt.Errorf("exactly one of AllPrivileges, GlobalPrivileges, AccountObjectPrivileges, SchemaPrivileges, or SchemaObjectPrivileges must be set")
	}
	return nil
}

type AccountRoleGrantOn struct {
	Account       *bool                 `ddl:"keyword" sql:"ACCOUNT"`
	AccountObject *GrantOnAccountObject `ddl:"-"`
	Schema        *GrantOnSchema        `ddl:"-"`
	SchemaObject  *GrantOnSchemaObject  `ddl:"-"`
}

func (v *AccountRoleGrantOn) validate() error {
	if !exactlyOneValueSet(v.Account, v.AccountObject, v.Schema, v.SchemaObject) {
		return fmt.Errorf("exactly one of Account, AccountObject, Schema, or SchemaObject must be set")
	}
	if valueSet(v.AccountObject) {
		if err := v.AccountObject.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			return err
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			return err
		}
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

func (v *GrantOnAccountObject) validate() error {
	if !exactlyOneValueSet(v.User, v.ResourceMonitor, v.Warehouse, v.Database, v.Integration, v.FailoverGroup, v.ReplicationGroup) {
		return fmt.Errorf("exactly one of User, ResourceMonitor, Warehouse, Database, Integration, FailoverGroup, or ReplicationGroup must be set")
	}
	return nil
}

type GrantOnSchema struct {
	Schema                  *SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	AllSchemasInDatabase    *AccountObjectIdentifier `ddl:"identifier" sql:"ALL SCHEMAS IN DATABASE"`
	FutureSchemasInDatabase *AccountObjectIdentifier `ddl:"identifier" sql:"FUTURE SCHEMAS IN DATABASE"`
}

func (v *GrantOnSchema) validate() error {
	if !exactlyOneValueSet(v.Schema, v.AllSchemasInDatabase, v.FutureSchemasInDatabase) {
		return fmt.Errorf("exactly one of Schema, AllSchemasInDatabase, or FutureSchemasInDatabase must be set")
	}
	return nil
}

type GrantOnSchemaObject struct {
	SchemaObject *Object                `ddl:"-"`
	All          *GrantOnSchemaObjectIn `ddl:"keyword" sql:"ALL"`
	Future       *GrantOnSchemaObjectIn `ddl:"keyword" sql:"FUTURE"`
}

func (v *GrantOnSchemaObject) validate() error {
	if !exactlyOneValueSet(v.SchemaObject, v.All, v.Future) {
		return fmt.Errorf("exactly one of Object, AllIn or Future must be set")
	}
	return nil
}

type GrantOnSchemaObjectIn struct {
	PluralObjectType PluralObjectType         `ddl:"keyword" sql:"ALL"`
	InDatabase       *AccountObjectIdentifier `ddl:"identifier" sql:"IN DATABASE"`
	InSchema         *SchemaIdentifier        `ddl:"identifier" sql:"IN SCHEMA"`
}

func (v *GrantOnSchemaObjectIn) validate() error {
	if !exactlyOneValueSet(v.PluralObjectType, v.InDatabase, v.InSchema) {
		return fmt.Errorf("exactly one of PluralObjectType, InDatabase, or InSchema must be set")
	}
	return nil
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

type RevokePrivilegesFromAccountRoleOptions struct {
	revoke         bool                        `ddl:"static" sql:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	GrantOptionFor *bool                       `ddl:"keyword" sql:"GRANT OPTION FOR"`
	privileges     *AccountRoleGrantPrivileges `ddl:"-"`
	on             *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole    AccountObjectIdentifier     `ddl:"identifier" sql:"FROM ROLE"`
	Restrict       *bool                       `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                       `ddl:"keyword" sql:"CASCADE"`
}

func (opts *RevokePrivilegesFromAccountRoleOptions) validate() error {
	if !valueSet(opts.privileges) {
		return fmt.Errorf("privileges must be set")
	}
	if err := opts.privileges.validate(); err != nil {
		return err
	}
	if !valueSet(opts.on) {
		return fmt.Errorf("on must be set")
	}
	if err := opts.on.validate(); err != nil {
		return err
	}
	if !validObjectidentifier(opts.accountRole) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		return fmt.Errorf("either Restrict or Cascade can be set, or neither but not both")
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
	show   bool          `ddl:"static" sql:"SHOW"` //lint:ignore U1000 This is used in the ddl tag
	Future *bool         `ddl:"keyword" sql:"FUTURE"`
	grants bool          `ddl:"static" sql:"GRANTS"` //lint:ignore U1000 This is used in the ddl tag
	On     *ShowGrantsOn `ddl:"keyword" sql:"ON"`
	To     *ShowGrantsTo `ddl:"keyword" sql:"TO"`
	Of     *ShowGrantsOf `ddl:"keyword" sql:"OF"`
	In     *ShowGrantsIn `ddl:"keyword" sql:"IN"`
}

func (opts *ShowGrantOptions) validate() error {
	if everyValueNil(opts.On, opts.To, opts.Of, opts.In) {
		return fmt.Errorf("at least one of on, to, of, or in is required")
	}
	if !exactlyOneValueSet(opts.On, opts.To, opts.Of, opts.In) {
		return fmt.Errorf("only one of on, to, of, or in can be set")
	}
	return nil
}

type ShowGrantsIn struct {
	Schema   *SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	Database *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
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
