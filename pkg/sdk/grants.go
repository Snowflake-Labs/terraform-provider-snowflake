package sdk

import (
	"context"
	"strings"
	"time"
)

var _ convertibleRow[Grant] = new(grantRow)

type Grants interface {
	GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error
	RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error
	GrantPrivilegesToDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *GrantPrivilegesToDatabaseRoleOptions) error
	GrantPrivilegeToShare(ctx context.Context, privilege ObjectPrivilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error
	RevokePrivilegeFromShare(ctx context.Context, privilege ObjectPrivilege, on *RevokePrivilegeFromShareOn, from AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowGrantOptions) ([]Grant, error)
}

type GrantPrivilegesToAccountRoleOptions struct {
	grant           bool                        `ddl:"static" sql:"GRANT"`
	privileges      *AccountRoleGrantPrivileges `ddl:"-"`
	on              *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole     AccountObjectIdentifier     `ddl:"identifier" sql:"TO ROLE"`
	WithGrantOption *bool                       `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

type AccountRoleGrantPrivileges struct {
	GlobalPrivileges        []GlobalPrivilege        `ddl:"-"`
	AccountObjectPrivileges []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges        []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges  []SchemaObjectPrivilege  `ddl:"-"`
	AllPrivileges           *bool                    `ddl:"keyword" sql:"ALL PRIVILEGES"`
}

type AccountRoleGrantOn struct {
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
	Schema                  *DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	AllSchemasInDatabase    *AccountObjectIdentifier  `ddl:"identifier" sql:"ALL SCHEMAS IN DATABASE"`
	FutureSchemasInDatabase *AccountObjectIdentifier  `ddl:"identifier" sql:"FUTURE SCHEMAS IN DATABASE"`
}

type GrantOnSchemaObject struct {
	SchemaObject *Object                `ddl:"-"`
	All          *GrantOnSchemaObjectIn `ddl:"keyword" sql:"ALL"`
	Future       *GrantOnSchemaObjectIn `ddl:"keyword" sql:"FUTURE"`
}

type GrantOnSchemaObjectIn struct {
	PluralObjectType PluralObjectType          `ddl:"keyword" sql:"ALL"`
	InDatabase       *AccountObjectIdentifier  `ddl:"identifier" sql:"IN DATABASE"`
	InSchema         *DatabaseObjectIdentifier `ddl:"identifier" sql:"IN SCHEMA"`
}

type RevokePrivilegesFromAccountRoleOptions struct {
	revoke         bool                        `ddl:"static" sql:"REVOKE"`
	GrantOptionFor *bool                       `ddl:"keyword" sql:"GRANT OPTION FOR"`
	privileges     *AccountRoleGrantPrivileges `ddl:"-"`
	on             *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole    AccountObjectIdentifier     `ddl:"identifier" sql:"FROM ROLE"`
	Restrict       *bool                       `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                       `ddl:"keyword" sql:"CASCADE"`
}

type GrantPrivilegesToDatabaseRoleOptions struct {
	grant           bool                         `ddl:"static" sql:"GRANT"`
	privileges      *DatabaseRoleGrantPrivileges `ddl:"-"`
	on              *DatabaseRoleGrantOn         `ddl:"keyword" sql:"ON"`
	databaseRole    DatabaseObjectIdentifier     `ddl:"identifier" sql:"TO DATABASE ROLE"`
	WithGrantOption *bool                        `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

type DatabaseRoleGrantPrivileges struct {
	// TODO: check other values (here only { CREATE SCHEMA| MODIFY | MONITOR | USAGE } [ , ... ] in doc)
	DatabasePrivileges     []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges       []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges []SchemaObjectPrivilege  `ddl:"-"`
}

type DatabaseRoleGrantOn struct {
	Database     *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Schema       *GrantOnSchema           `ddl:"-"`
	SchemaObject *GrantOnSchemaObject     `ddl:"-"`
}

type grantPrivilegeToShareOptions struct {
	grant     bool                     `ddl:"static" sql:"GRANT"`
	privilege ObjectPrivilege          `ddl:"keyword"`
	On        *GrantPrivilegeToShareOn `ddl:"keyword" sql:"ON"`
	to        AccountObjectIdentifier  `ddl:"identifier" sql:"TO SHARE"`
}

type GrantPrivilegeToShareOn struct {
	Database AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
	Schema   DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	Function SchemaObjectIdentifier   `ddl:"identifier" sql:"FUNCTION"`
	Table    *OnTable                 `ddl:"-"`
	View     SchemaObjectIdentifier   `ddl:"identifier" sql:"VIEW"`
}

type OnTable struct {
	Name        SchemaObjectIdentifier   `ddl:"identifier" sql:"TABLE"`
	AllInSchema DatabaseObjectIdentifier `ddl:"identifier" sql:"ALL TABLES IN SCHEMA"`
}

type revokePrivilegeFromShareOptions struct {
	revoke    bool                        `ddl:"static" sql:"REVOKE"`
	privilege ObjectPrivilege             `ddl:"keyword"`
	On        *RevokePrivilegeFromShareOn `ddl:"keyword" sql:"ON"`
	from      AccountObjectIdentifier     `ddl:"identifier" sql:"FROM SHARE"`
}

type RevokePrivilegeFromShareOn struct {
	Database AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
	Schema   DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	Table    *OnTable                 `ddl:"-"`
	View     *OnView                  `ddl:"-"`
}

type OnView struct {
	Name        SchemaObjectIdentifier   `ddl:"identifier" sql:"VIEW"`
	AllInSchema DatabaseObjectIdentifier `ddl:"identifier" sql:"ALL VIEWS IN SCHEMA"`
}

type ShowGrantOptions struct {
	show   bool          `ddl:"static" sql:"SHOW"`
	Future *bool         `ddl:"keyword" sql:"FUTURE"`
	grants bool          `ddl:"static" sql:"GRANTS"`
	On     *ShowGrantsOn `ddl:"keyword" sql:"ON"`
	To     *ShowGrantsTo `ddl:"keyword" sql:"TO"`
	Of     *ShowGrantsOf `ddl:"keyword" sql:"OF"`
	In     *ShowGrantsIn `ddl:"keyword" sql:"IN"`
}

type ShowGrantsIn struct {
	Schema   *DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	Database *AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
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

type Grant struct {
	CreatedOn   time.Time
	Privilege   string
	GrantedOn   ObjectType
	GrantOn     ObjectType
	Name        ObjectIdentifier
	GrantedTo   ObjectType
	GranteeName AccountObjectIdentifier
	GrantOption bool
	GrantedBy   AccountObjectIdentifier
}

func (v *Grant) ID() ObjectIdentifier {
	return v.Name
}

func (row grantRow) convert() *Grant {
	grantedTo := ObjectType(strings.ReplaceAll(row.GrantedTo, "_", " "))
	granteeName := NewAccountObjectIdentifier(row.GranteeName)
	if grantedTo == ObjectTypeShare {
		parts := strings.Split(row.GranteeName, ".")
		name := strings.Join(parts[1:], ".")
		granteeName = NewAccountObjectIdentifier(name)
	}
	grant := &Grant{
		CreatedOn:   row.CreatedOn,
		Privilege:   row.Privilege,
		GrantedTo:   grantedTo,
		Name:        NewAccountObjectIdentifier(strings.Trim(row.Name, "\"")),
		GranteeName: granteeName,
		GrantOption: row.GrantOption,
		GrantedBy:   NewAccountObjectIdentifier(row.GrantedBy),
	}

	// true for current grants
	if row.GrantedOn != "" {
		grant.GrantedOn = ObjectType(strings.ReplaceAll(row.GrantedOn, "_", " "))
	}
	// true for future grants
	if row.GrantOn != "" {
		grant.GrantOn = ObjectType(strings.ReplaceAll(row.GrantOn, "_", " "))
	}
	return grant
}
