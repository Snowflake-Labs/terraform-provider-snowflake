package sdk

import (
	"context"
	"database/sql"
)

type Views interface {
	Create(ctx context.Context, request *CreateViewRequest) error
	Alter(ctx context.Context, request *AlterViewRequest) error
	Drop(ctx context.Context, request *DropViewRequest) error
	Show(ctx context.Context, request *ShowViewRequest) ([]View, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*View, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]ViewDetails, error)
}

// CreateViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-view.
type CreateViewOptions struct {
	create                 bool                      `ddl:"static" sql:"CREATE"`
	OrReplace              *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                 *bool                     `ddl:"keyword" sql:"SECURE"`
	Temporary              *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Recursive              *bool                     `ddl:"keyword" sql:"RECURSIVE"`
	view                   bool                      `ddl:"static" sql:"VIEW"`
	IfNotExists            *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                   SchemaObjectIdentifier    `ddl:"identifier"`
	Columns                []ViewColumn              `ddl:"list,parentheses"`
	ColumnsMaskingPolicies []ViewColumnMaskingPolicy `ddl:"list,no_parentheses,no_equals"`
	CopyGrants             *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Comment                *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy        *ViewRowAccessPolicy      `ddl:"keyword"`
	Tag                    []TagAssociation          `ddl:"keyword,parentheses" sql:"TAG"`
	as                     bool                      `ddl:"static" sql:"AS"`
	sql                    string                    `ddl:"keyword,no_quotes"`
}

type ViewColumn struct {
	Name    string  `ddl:"keyword,double_quotes"`
	Comment *string `ddl:"parameter,single_quotes,no_equals" sql:"COMMENT"`
}

type ViewColumnMaskingPolicy struct {
	Name          string                 `ddl:"keyword"`
	MaskingPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
	Using         []string               `ddl:"list,parentheses" sql:"USING"`
	Tag           []TagAssociation       `ddl:"keyword,parentheses" sql:"TAG"`
}

type ViewRowAccessPolicy struct {
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []string               `ddl:"list,parentheses" sql:"ON"`
}

// AlterViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-view.
type AlterViewOptions struct {
	alter                      bool                           `ddl:"static" sql:"ALTER"`
	view                       bool                           `ddl:"static" sql:"VIEW"`
	IfExists                   *bool                          `ddl:"keyword" sql:"IF EXISTS"`
	name                       SchemaObjectIdentifier         `ddl:"identifier"`
	RenameTo                   *SchemaObjectIdentifier        `ddl:"identifier" sql:"RENAME TO"`
	SetComment                 *string                        `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment               *bool                          `ddl:"keyword" sql:"UNSET COMMENT"`
	SetSecure                  *bool                          `ddl:"keyword" sql:"SET SECURE"`
	SetChangeTracking          *bool                          `ddl:"parameter" sql:"SET CHANGE_TRACKING"`
	UnsetSecure                *bool                          `ddl:"keyword" sql:"UNSET SECURE"`
	SetTags                    []TagAssociation               `ddl:"keyword" sql:"SET TAG"`
	UnsetTags                  []ObjectIdentifier             `ddl:"keyword" sql:"UNSET TAG"`
	AddRowAccessPolicy         *ViewAddRowAccessPolicy        `ddl:"keyword"`
	DropRowAccessPolicy        *ViewDropRowAccessPolicy       `ddl:"keyword"`
	DropAndAddRowAccessPolicy  *ViewDropAndAddRowAccessPolicy `ddl:"list,no_parentheses"`
	DropAllRowAccessPolicies   *bool                          `ddl:"keyword" sql:"DROP ALL ROW ACCESS POLICIES"`
	SetMaskingPolicyOnColumn   *ViewSetColumnMaskingPolicy    `ddl:"keyword"`
	UnsetMaskingPolicyOnColumn *ViewUnsetColumnMaskingPolicy  `ddl:"keyword"`
	SetTagsOnColumn            *ViewSetColumnTags             `ddl:"keyword"`
	UnsetTagsOnColumn          *ViewUnsetColumnTags           `ddl:"keyword"`
}

type ViewAddRowAccessPolicy struct {
	add             bool                   `ddl:"static" sql:"ADD"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []string               `ddl:"list,parentheses" sql:"ON"`
}

type ViewDropRowAccessPolicy struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
}

type ViewDropAndAddRowAccessPolicy struct {
	Drop ViewDropRowAccessPolicy `ddl:"keyword"`
	Add  ViewAddRowAccessPolicy  `ddl:"keyword"`
}

type ViewSetColumnMaskingPolicy struct {
	alter         bool                   `ddl:"static" sql:"ALTER"`
	column        bool                   `ddl:"static" sql:"COLUMN"`
	Name          string                 `ddl:"keyword"`
	set           bool                   `ddl:"static" sql:"SET"`
	MaskingPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
	Using         []string               `ddl:"list,parentheses" sql:"USING"`
	Force         *bool                  `ddl:"keyword" sql:"FORCE"`
}

type ViewUnsetColumnMaskingPolicy struct {
	alter         bool   `ddl:"static" sql:"ALTER"`
	column        bool   `ddl:"static" sql:"COLUMN"`
	Name          string `ddl:"keyword"`
	unset         bool   `ddl:"static" sql:"UNSET"`
	maskingPolicy bool   `ddl:"static" sql:"MASKING POLICY"`
}

type ViewSetColumnTags struct {
	alter   bool             `ddl:"static" sql:"ALTER"`
	column  bool             `ddl:"static" sql:"COLUMN"`
	Name    string           `ddl:"keyword"`
	SetTags []TagAssociation `ddl:"keyword" sql:"SET TAG"`
}

type ViewUnsetColumnTags struct {
	alter     bool               `ddl:"static" sql:"ALTER"`
	column    bool               `ddl:"static" sql:"COLUMN"`
	Name      string             `ddl:"keyword"`
	UnsetTags []ObjectIdentifier `ddl:"keyword" sql:"UNSET TAG"`
}

// DropViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-view.
type DropViewOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	view     bool                   `ddl:"static" sql:"VIEW"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-views.
type ShowViewOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	views      bool       `ddl:"static" sql:"VIEWS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type viewDBRow struct {
	CreatedOn      string         `db:"created_on"`
	Name           string         `db:"name"`
	Kind           sql.NullString `db:"kind"`
	Reserved       sql.NullString `db:"reserved"`
	DatabaseName   string         `db:"database_name"`
	SchemaName     string         `db:"schema_name"`
	Owner          sql.NullString `db:"owner"`
	Comment        sql.NullString `db:"comment"`
	Text           sql.NullString `db:"text"`
	IsSecure       sql.NullBool   `db:"is_secure"`
	IsMaterialized sql.NullBool   `db:"is_materialized"`
	OwnerRoleType  sql.NullString `db:"owner_role_type"`
	ChangeTracking sql.NullString `db:"change_tracking"`
}

type View struct {
	CreatedOn      string
	Name           string
	Kind           string
	Reserved       string
	DatabaseName   string
	SchemaName     string
	Owner          string
	Comment        string
	Text           string
	IsSecure       bool
	IsMaterialized bool
	OwnerRoleType  string
	ChangeTracking string
}

func (v *View) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

// DescribeViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-view.
type DescribeViewOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	view     bool                   `ddl:"static" sql:"VIEW"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// TODO: extract common type for describe
// viewDetailsRow is a copy of externalTableColumnDetailsRow.
type viewDetailsRow struct {
	Name       string         `db:"name"`
	Type       DataType       `db:"type"`
	Kind       string         `db:"kind"`
	IsNullable string         `db:"null?"`
	Default    sql.NullString `db:"default"`
	IsPrimary  string         `db:"primary key"`
	IsUnique   string         `db:"unique key"`
	Check      sql.NullString `db:"check"`
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
	PolicyName sql.NullString `db:"policy name"`
}

// ViewDetails is a copy of ExternalTableColumnDetails.
type ViewDetails struct {
	Name       string
	Type       DataType
	Kind       string
	IsNullable bool
	Default    *string
	IsPrimary  bool
	IsUnique   bool
	Check      *bool
	Expression *string
	Comment    *string
	PolicyName *string
}
