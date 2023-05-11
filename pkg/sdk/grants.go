package sdk

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type Grants interface {
	GrantPrivilegeToShare(ctx context.Context, objectPrivilege Privilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error
	RevokePrivilegeFromShare(ctx context.Context, objectPrivilege Privilege, on *RevokePrivilegeFromShareOn, from AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowGrantsOptions) ([]*Grant, error)
}

var _ Grants = (*grants)(nil)

type grants struct {
	client  *Client
	builder *sqlBuilder
}

type Grant struct {
	CreatedOn   time.Time
	Privilege   Privilege
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
		Privilege:   Privilege(row.Privilege),
		GrantedOn:   ObjectType(row.GrantedOn),
		GrantedTo:   grantedTo,
		Name:        NewAccountObjectIdentifier(strings.Trim(row.Name, "\"")),
		GranteeName: granteeName,
		GrantOption: row.GrantOption,
		GrantedBy:   NewAccountObjectIdentifier(row.GrantedBy),
	}
	return grant, nil
}

type grantPrivilegeToShareOptions struct {
	grant           *bool                    `ddl:"static" db:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	objectPrivilege Privilege                `ddl:"keyword"`
	On              *GrantPrivilegeToShareOn `ddl:"keyword" db:"ON"`
	to              AccountObjectIdentifier  `ddl:"identifier" db:"TO SHARE"`
}

func (opts *grantPrivilegeToShareOptions) validate() error {
	if opts.to.FullyQualifiedName() == "" {
		return fmt.Errorf("to is required")
	}
	if opts.On == nil || opts.objectPrivilege == "" {
		return fmt.Errorf("on and objectPrivilege are required")
	}
	return nil
}

type GrantPrivilegeToShareOn struct {
	Database AccountObjectIdentifier `ddl:"identifier" db:"DATABASE"`
	Schema   SchemaIdentifier        `ddl:"identifier" db:"SCHEMA"`
	Function SchemaObjectIdentifier  `ddl:"identifier" db:"FUNCTION"`
	Table    *OnTable                `ddl:"keyword"`
	View     SchemaObjectIdentifier  `ddl:"identifier" db:"VIEW"`
}

type OnTable struct {
	Name        SchemaObjectIdentifier `ddl:"identifier" db:"TABLE"`
	AllInSchema SchemaIdentifier       `ddl:"identifier" db:"ALL TABLES IN SCHEMA"`
}

func (v *grants) GrantPrivilegeToShare(ctx context.Context, objectPrivilege Privilege, on *GrantPrivilegeToShareOn, to AccountObjectIdentifier) error {
	opts := &grantPrivilegeToShareOptions{
		objectPrivilege: objectPrivilege,
		On:              on,
		to:              to,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	sql := v.builder.sql(clauses...)
	_, err = v.client.exec(ctx, sql)
	return err
}

type revokePrivilegeFromShareOptions struct {
	revoke          *bool                       `ddl:"static" db:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	objectPrivilege Privilege                   `ddl:"keyword"`
	On              *RevokePrivilegeFromShareOn `ddl:"keyword" db:"ON"`
	from            AccountObjectIdentifier     `ddl:"identifier" db:"FROM SHARE"`
}

func (opts *revokePrivilegeFromShareOptions) validate() error {
	if opts.from.FullyQualifiedName() == "" {
		return fmt.Errorf("from is required")
	}
	if opts.On == nil || opts.objectPrivilege == "" {
		return fmt.Errorf("on and objectPrivilege are required")
	}
	return nil
}

type RevokePrivilegeFromShareOn struct {
	Database AccountObjectIdentifier `ddl:"identifier" db:"DATABASE"`
	Schema   SchemaIdentifier        `ddl:"identifier" db:"SCHEMA"`
	Table    *OnTable                `ddl:"keyword"`
	View     *OnView                 `ddl:"keyword"`
}

type OnView struct {
	Name        SchemaObjectIdentifier `ddl:"identifier" db:"VIEW"`
	AllInSchema SchemaIdentifier       `ddl:"identifier" db:"ALL VIEWS IN SCHEMA"`
}

func (v *grants) RevokePrivilegeFromShare(ctx context.Context, objectPrivilege Privilege, on *RevokePrivilegeFromShareOn, id AccountObjectIdentifier) error {
	opts := &revokePrivilegeFromShareOptions{
		objectPrivilege: objectPrivilege,
		On:              on,
		from:            id,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return err
	}
	sql := v.builder.sql(clauses...)
	log.Printf("sql: %s", sql)
	_, err = v.client.exec(ctx, sql)
	return err
}

type ShowGrantsOptions struct {
	show   *bool         `ddl:"static" db:"SHOW"`   //lint:ignore U1000 This is used in the ddl tag
	grants *bool         `ddl:"static" db:"GRANTS"` //lint:ignore U1000 This is used in the ddl tag
	On     *ShowGrantsOn `ddl:"keyword" db:"ON"`
	To     *ShowGrantsTo `ddl:"keyword" db:"TO"`
	Of     *ShowGrantsOf `ddl:"keyword" db:"OF"`
}

func (opts *ShowGrantsOptions) validate() error {
	if opts.On == nil && opts.To == nil && opts.Of == nil {
		return fmt.Errorf("at least one of on, to, or of is required")
	}
	return nil
}

type ShowGrantsOn struct {
	Account *bool `ddl:"keyword" db:"ACCOUNT"`
	Object  *Object
}

type ShowGrantsTo struct {
	Role  AccountObjectIdentifier `ddl:"identifier" db:"ROLE"`
	User  AccountObjectIdentifier `ddl:"identifier" db:"USER"`
	Share AccountObjectIdentifier `ddl:"identifier" db:"SHARE"`
}

type ShowGrantsOf struct {
	Role  AccountObjectIdentifier `ddl:"identifier" db:"ROLE"`
	Share AccountObjectIdentifier `ddl:"identifier" db:"SHARE"`
}

func (v *grants) Show(ctx context.Context, opts *ShowGrantsOptions) ([]*Grant, error) {
	if opts == nil {
		opts = &ShowGrantsOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	clauses, err := v.builder.parseStruct(opts)
	if err != nil {
		return nil, err
	}
	sql := v.builder.sql(clauses...)
	var rows []grantRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, decodeDriverError(err)
	}
	var grants []*Grant
	for _, row := range rows {
		grant, err := row.toGrant()
		if err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}
	return grants, nil
}
