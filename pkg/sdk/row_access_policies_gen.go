package sdk

import (
	"context"
	"database/sql"
)

type RowAccessPolicies interface {
	Create(ctx context.Context, request *CreateRowAccessPolicyRequest) error
	Alter(ctx context.Context, request *AlterRowAccessPolicyRequest) error
	Drop(ctx context.Context, request *DropRowAccessPolicyRequest) error
	Show(ctx context.Context, request *ShowRowAccessPolicyRequest) ([]RowAccessPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*RowAccessPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*RowAccessPolicyDescription, error)
}

// CreateRowAccessPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-row-access-policy.
type CreateRowAccessPolicyOptions struct {
	create          bool                        `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                       `ddl:"keyword" sql:"OR REPLACE"`
	rowAccessPolicy bool                        `ddl:"static" sql:"ROW ACCESS POLICY"`
	IfNotExists     *bool                       `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier      `ddl:"identifier"`
	as              bool                        `ddl:"static" sql:"AS"`
	args            []CreateRowAccessPolicyArgs `ddl:"parameter,parentheses,no_equals"`
	returnsBoolean  bool                        `ddl:"static" sql:"RETURNS BOOLEAN"`
	body            string                      `ddl:"parameter,no_quotes,no_equals" sql:"->"`
	Comment         *string                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type CreateRowAccessPolicyArgs struct {
	Name string   `ddl:"keyword,no_quotes"`
	Type DataType `ddl:"keyword,no_quotes"`
}

// AlterRowAccessPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-row-access-policy.
type AlterRowAccessPolicyOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`
	rowAccessPolicy bool                    `ddl:"static" sql:"ROW ACCESS POLICY"`
	name            SchemaObjectIdentifier  `ddl:"identifier"`
	RenameTo        *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetBody         *string                 `ddl:"parameter,no_quotes,no_equals" sql:"SET BODY ->"`
	SetTags         []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags       []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
	SetComment      *string                 `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment    *bool                   `ddl:"keyword" sql:"UNSET COMMENT"`
}

// DropRowAccessPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-row-access-policy.
type DropRowAccessPolicyOptions struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"`
	IfExists        *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowRowAccessPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies.
type ShowRowAccessPolicyOptions struct {
	show              bool  `ddl:"static" sql:"SHOW"`
	rowAccessPolicies bool  `ddl:"static" sql:"ROW ACCESS POLICIES"`
	Like              *Like `ddl:"keyword" sql:"LIKE"`
	In                *In   `ddl:"keyword" sql:"IN"`
}

type rowAccessPolicyDBRow struct {
	CreatedOn     string         `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Kind          string         `db:"kind"`
	Owner         string         `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       string         `db:"options"`
	OwnerRoleType string         `db:"owner_role_type"`
}

type RowAccessPolicy struct {
	CreatedOn     string
	Name          string
	DatabaseName  string
	SchemaName    string
	Kind          string
	Owner         string
	Comment       string
	Options       string
	OwnerRoleType string
}

func (v *RowAccessPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

// DescribeRowAccessPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-row-access-policy.
type DescribeRowAccessPolicyOptions struct {
	describe        bool                   `ddl:"static" sql:"DESCRIBE"`
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
}

type describeRowAccessPolicyDBRow struct {
	Name       string `db:"name"`
	Signature  string `db:"signature"`
	ReturnType string `db:"return_type"`
	Body       string `db:"body"`
}

type RowAccessPolicyDescription struct {
	Name       string
	Signature  string
	ReturnType string
	Body       string
}
