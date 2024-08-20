package sdk

import (
	"context"
	"database/sql"
	"strings"
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
	create            bool                   `ddl:"static" sql:"CREATE"`
	OrReplace         *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure            *bool                  `ddl:"keyword" sql:"SECURE"`
	Temporary         *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Recursive         *bool                  `ddl:"keyword" sql:"RECURSIVE"`
	view              bool                   `ddl:"static" sql:"VIEW"`
	IfNotExists       *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	Columns           []ViewColumn           `ddl:"list,parentheses"`
	CopyGrants        *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Comment           *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy   *ViewRowAccessPolicy   `ddl:"keyword"`
	AggregationPolicy *ViewAggregationPolicy `ddl:"keyword"`
	Tag               []TagAssociation       `ddl:"keyword,parentheses" sql:"TAG"`
	as                bool                   `ddl:"static" sql:"AS"`
	sql               string                 `ddl:"keyword,no_quotes"`
}
type ViewColumn struct {
	Name             string                      `ddl:"keyword,double_quotes"`
	ProjectionPolicy *ViewColumnProjectionPolicy `ddl:"keyword"`
	MaskingPolicy    *ViewColumnMaskingPolicy    `ddl:"keyword"`
	Comment          *string                     `ddl:"parameter,single_quotes,no_equals" sql:"COMMENT"`
	Tag              []TagAssociation            `ddl:"keyword,parentheses" sql:"TAG"`
}
type ViewColumnProjectionPolicy struct {
	ProjectionPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"PROJECTION POLICY"`
}
type ViewColumnMaskingPolicy struct {
	MaskingPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
	Using         []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"USING"`
}
type ViewRowAccessPolicy struct {
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"ON"`
}
type ViewAggregationPolicy struct {
	AggregationPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"AGGREGATION POLICY"`
	EntityKey         []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"ENTITY KEY"`
}

// AlterViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-view.
type AlterViewOptions struct {
	alter                         bool                           `ddl:"static" sql:"ALTER"`
	view                          bool                           `ddl:"static" sql:"VIEW"`
	IfExists                      *bool                          `ddl:"keyword" sql:"IF EXISTS"`
	name                          SchemaObjectIdentifier         `ddl:"identifier"`
	RenameTo                      *SchemaObjectIdentifier        `ddl:"identifier" sql:"RENAME TO"`
	SetComment                    *string                        `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment                  *bool                          `ddl:"keyword" sql:"UNSET COMMENT"`
	SetSecure                     *bool                          `ddl:"keyword" sql:"SET SECURE"`
	SetChangeTracking             *bool                          `ddl:"parameter" sql:"SET CHANGE_TRACKING"`
	UnsetSecure                   *bool                          `ddl:"keyword" sql:"UNSET SECURE"`
	SetTags                       []TagAssociation               `ddl:"keyword" sql:"SET TAG"`
	UnsetTags                     []ObjectIdentifier             `ddl:"keyword" sql:"UNSET TAG"`
	AddDataMetricFunction         *ViewAddDataMetricFunction     `ddl:"keyword"`
	DropDataMetricFunction        *ViewDropDataMetricFunction    `ddl:"keyword"`
	SetDataMetricSchedule         *ViewSetDataMetricSchedule     `ddl:"keyword"`
	UnsetDataMetricSchedule       *ViewUnsetDataMetricSchedule   `ddl:"keyword"`
	AddRowAccessPolicy            *ViewAddRowAccessPolicy        `ddl:"keyword"`
	DropRowAccessPolicy           *ViewDropRowAccessPolicy       `ddl:"keyword"`
	DropAndAddRowAccessPolicy     *ViewDropAndAddRowAccessPolicy `ddl:"list,no_parentheses"`
	DropAllRowAccessPolicies      *bool                          `ddl:"keyword" sql:"DROP ALL ROW ACCESS POLICIES"`
	SetAggregationPolicy          *ViewSetAggregationPolicy      `ddl:"keyword"`
	UnsetAggregationPolicy        *ViewUnsetAggregationPolicy    `ddl:"keyword"`
	SetMaskingPolicyOnColumn      *ViewSetColumnMaskingPolicy    `ddl:"keyword"`
	UnsetMaskingPolicyOnColumn    *ViewUnsetColumnMaskingPolicy  `ddl:"keyword"`
	SetProjectionPolicyOnColumn   *ViewSetProjectionPolicy       `ddl:"keyword"`
	UnsetProjectionPolicyOnColumn *ViewUnsetProjectionPolicy     `ddl:"keyword"`
	SetTagsOnColumn               *ViewSetColumnTags             `ddl:"keyword"`
	UnsetTagsOnColumn             *ViewUnsetColumnTags           `ddl:"keyword"`
}
type DoubleQuotedString struct {
	Value string `ddl:"keyword,double_quotes"`
}
type ViewDataMetricFunction struct {
	DataMetricFunction SchemaObjectIdentifier `ddl:"identifier"`
	On                 []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"ON"`
}
type ViewAddDataMetricFunction struct {
	add                bool                     `ddl:"static" sql:"ADD"`
	DataMetricFunction []ViewDataMetricFunction `ddl:"parameter,no_equals" sql:"DATA METRIC FUNCTION"`
}
type ViewDropDataMetricFunction struct {
	drop               bool                     `ddl:"static" sql:"DROP"`
	DataMetricFunction []ViewDataMetricFunction `ddl:"parameter,no_equals" sql:"DATA METRIC FUNCTION"`
}
type ViewSetDataMetricSchedule struct {
	setDataMetricSchedule bool           `ddl:"static" sql:"SET DATA_METRIC_SCHEDULE ="`
	Minutes               *ViewMinute    `ddl:"keyword"`
	UsingCron             *ViewUsingCron `ddl:"keyword"`
	TriggerOnChanges      *bool          `ddl:"keyword,single_quotes" sql:"TRIGGER_ON_CHANGES"`
}
type ViewMinute struct {
	prefix  bool `ddl:"static" sql:"'"`
	Minutes int  `ddl:"keyword"`
	suffix  bool `ddl:"static" sql:"MINUTE'"`
}
type ViewUsingCron struct {
	prefix bool   `ddl:"static" sql:"'USING CRON"`
	Cron   string `ddl:"keyword"`
	suffix bool   `ddl:"static" sql:"'"`
}
type ViewTriggerOnChanges struct {
	triggerOnChanges bool `ddl:"static" sql:"TRIGGER_ON_CHANGES"`
}
type ViewUnsetDataMetricSchedule struct {
	unsetDataMetricSchedule bool `ddl:"static" sql:"UNSET DATA_METRIC_SCHEDULE"`
}
type ViewAddRowAccessPolicy struct {
	add             bool                   `ddl:"static" sql:"ADD"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"ON"`
}
type ViewDropRowAccessPolicy struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
}
type ViewDropAndAddRowAccessPolicy struct {
	Drop ViewDropRowAccessPolicy `ddl:"keyword"`
	Add  ViewAddRowAccessPolicy  `ddl:"keyword"`
}
type ViewSetAggregationPolicy struct {
	set               bool                   `ddl:"static" sql:"SET"`
	AggregationPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"AGGREGATION POLICY"`
	EntityKey         []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"ENTITY KEY"`
	Force             *bool                  `ddl:"keyword" sql:"FORCE"`
}
type ViewUnsetAggregationPolicy struct {
	unsetAggregationPolicy bool `ddl:"static" sql:"UNSET AGGREGATION POLICY"`
}
type ViewSetColumnMaskingPolicy struct {
	alter         bool                   `ddl:"static" sql:"ALTER"`
	column        bool                   `ddl:"static" sql:"COLUMN"`
	Name          string                 `ddl:"keyword,double_quotes"`
	set           bool                   `ddl:"static" sql:"SET"`
	MaskingPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
	Using         []DoubleQuotedString   `ddl:"parameter,parentheses,no_equals" sql:"USING"`
	Force         *bool                  `ddl:"keyword" sql:"FORCE"`
}
type ViewUnsetColumnMaskingPolicy struct {
	alter         bool   `ddl:"static" sql:"ALTER"`
	column        bool   `ddl:"static" sql:"COLUMN"`
	Name          string `ddl:"keyword,double_quotes"`
	unset         bool   `ddl:"static" sql:"UNSET"`
	maskingPolicy bool   `ddl:"static" sql:"MASKING POLICY"`
}
type ViewSetProjectionPolicy struct {
	alter            bool                   `ddl:"static" sql:"ALTER"`
	column           bool                   `ddl:"static" sql:"COLUMN"`
	Name             string                 `ddl:"keyword,double_quotes"`
	set              bool                   `ddl:"static" sql:"SET"`
	ProjectionPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"PROJECTION POLICY"`
	Force            *bool                  `ddl:"keyword" sql:"FORCE"`
}
type ViewUnsetProjectionPolicy struct {
	alter            bool   `ddl:"static" sql:"ALTER"`
	column           bool   `ddl:"static" sql:"COLUMN"`
	Name             string `ddl:"keyword,double_quotes"`
	unset            bool   `ddl:"static" sql:"UNSET"`
	projectionPolicy bool   `ddl:"static" sql:"PROJECTION POLICY"`
}
type ViewSetColumnTags struct {
	alter   bool             `ddl:"static" sql:"ALTER"`
	column  bool             `ddl:"static" sql:"COLUMN"`
	Name    string           `ddl:"keyword,double_quotes"`
	SetTags []TagAssociation `ddl:"keyword" sql:"SET TAG"`
}
type ViewUnsetColumnTags struct {
	alter     bool               `ddl:"static" sql:"ALTER"`
	column    bool               `ddl:"static" sql:"COLUMN"`
	Name      string             `ddl:"keyword,double_quotes"`
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
	show       bool        `ddl:"static" sql:"SHOW"`
	Terse      *bool       `ddl:"keyword" sql:"TERSE"`
	views      bool        `ddl:"static" sql:"VIEWS"`
	Like       *Like       `ddl:"keyword" sql:"LIKE"`
	In         *ExtendedIn `ddl:"keyword" sql:"IN"`
	StartsWith *string     `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
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

func (v *View) HasCopyGrants() bool {
	return strings.Contains(v.Text, " COPY GRANTS ")
}

// DescribeViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-view.
type DescribeViewOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	view     bool                   `ddl:"static" sql:"VIEW"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}
type viewDetailsRow struct {
	Name          string         `db:"name"`
	Type          DataType       `db:"type"`
	Kind          string         `db:"kind"`
	Null          string         `db:"null?"`
	Default       sql.NullString `db:"default"`
	PrimaryKey    string         `db:"primary key"`
	UniqueKey     string         `db:"unique key"`
	Check         sql.NullString `db:"check"`
	Expression    sql.NullString `db:"expression"`
	Comment       sql.NullString `db:"comment"`
	PolicyName    sql.NullString `db:"policy name"`
	PrivacyDomain sql.NullString `db:"privacy domain"`
}
type ViewDetails struct {
	Name          string
	Type          DataType
	Kind          string
	IsNullable    bool
	Default       *string
	IsPrimary     bool
	IsUnique      bool
	Check         *bool
	Expression    *string
	Comment       *string
	PolicyName    *string
	PrivacyDomain *string
}
