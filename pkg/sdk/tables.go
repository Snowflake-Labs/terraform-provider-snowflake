package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

var _ convertibleRow[Table] = new(tableDBRow)

type Tables interface {
	Create(ctx context.Context, req *CreateTableRequest) error
	CreateAsSelect(ctx context.Context, req *CreateTableAsSelectRequest) error
	CreateUsingTemplate(ctx context.Context, req *CreateTableUsingTemplateRequest) error
	CreateLike(ctx context.Context, req *CreateTableLikeRequest) error
	CreateClone(ctx context.Context, req *CreateTableCloneRequest) error
	Alter(ctx context.Context, req *AlterTableRequest) error
	Drop(ctx context.Context, req *DropTableRequest) error
	Show(ctx context.Context, req *ShowTableRequest) ([]Table, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Table, error)
}

type createTableAsSelectOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table           bool                   `ddl:"static" sql:"TABLE"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	leftParen       bool                   `ddl:"static" sql:"("`
	Columns         []TableAsSelectColumn  `ddl:"keyword"`
	rightParen      bool                   `ddl:"static" sql:")"`
	ClusterBy       []string               `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	CopyGrants      *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy *RowAccessPolicy       `ddl:"keyword"`
	Query           *string                `ddl:"parameter,no_equals" sql:"AS SELECT"`
}

type TableAsSelectColumn struct {
	Name          string                            `ddl:"keyword"`
	Type          *DataType                         `ddl:"keyword"`
	MaskingPolicy *TableAsSelectColumnMaskingPolicy `ddl:"keyword"`
}

type TableAsSelectColumnMaskingPolicy struct {
	With          *bool                  `ddl:"keyword" sql:"WITH"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	Name          SchemaObjectIdentifier `ddl:"identifier"`
}

type createTableUsingTemplateOptions struct {
	create     bool                   `ddl:"static" sql:"CREATE"`
	OrReplace  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table      bool                   `ddl:"static" sql:"TABLE"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Query      []string               `ddl:"parameter,no_equals,parentheses" sql:"USING TEMPLATE"`
}

type createTableLikeOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table       bool                   `ddl:"static" sql:"TABLE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	like        bool                   `ddl:"static" sql:"LIKE"`
	SourceTable SchemaObjectIdentifier `ddl:"identifier"`
	ClusterBy   []string               `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	CopyGrants  *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type createTableCloneOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table       bool                   `ddl:"static" sql:"TABLE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	clone       bool                   `ddl:"static" sql:"CLONE"`
	SourceTable SchemaObjectIdentifier `ddl:"identifier"`
	ClonePoint  *ClonePoint            `ddl:"keyword"`
	CopyGrants  *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type ClonePoint struct {
	Moment CloneMoment `ddl:"parameter,no_equals"`
	At     TimeTravel  `ddl:"list,parentheses,no_comma"`
}

type CloneMoment string

const (
	CloneMomentAt     CloneMoment = "AT"
	CloneMomentBefore CloneMoment = "BEFORE"
)

type createTableOptions struct {
	create                     bool                       `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                      `ddl:"keyword" sql:"OR REPLACE"`
	Scope                      *TableScope                `ddl:"keyword"`
	Kind                       *TableKind                 `ddl:"keyword"`
	table                      bool                       `ddl:"static" sql:"TABLE"`
	IfNotExists                *bool                      `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier     `ddl:"identifier"`
	leftParen                  bool                       `ddl:"static" sql:"("`
	Columns                    []TableColumn              `ddl:"keyword"`
	OutOfLineConstraint        *CreateOutOfLineConstraint `ddl:"keyword"`
	rightParen                 bool                       `ddl:"static" sql:")"`
	ClusterBy                  []string                   `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	EnableSchemaEvolution      *bool                      `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            []StageFileFormat          `ddl:"parameter,equals,parentheses" sql:"STAGE_FILE_FORMAT"`
	StageCopyOptions           []StageCopyOption          `ddl:"parameter,equals,parentheses" sql:"STAGE_COPY_OPTIONS"`
	DataRetentionTimeInDays    *int                       `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int                       `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool                      `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string                    `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	CopyGrants                 *bool                      `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy            *RowAccessPolicy           `ddl:"keyword"`
	Tags                       []TagAssociation           `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                    `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type TableScope string

const (
	GlobalTableScope TableScope = "GLOBAL"
	LocalTableScope  TableScope = "LOCAL"
)

type TableKind string

const (
	TemporaryTableKind TableKind = "TEMPORARY"
	VolatileTableKind  TableKind = "VOLATILE"
	TransientTableKind TableKind = "TRANSIENT"
)

type TableColumn struct {
	Name             string                  `ddl:"keyword"`
	Type             DataType                `ddl:"keyword"`
	Collate          *string                 `ddl:"parameter,no_equals,single_quotes" sql:"COLLATE"`
	Comment          *string                 `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
	DefaultValue     *ColumnDefaultValue     `ddl:"keyword"`
	NotNull          *bool                   `ddl:"keyword" sql:"NOT NULL"`
	MaskingPolicy    *ColumnMaskingPolicy    `ddl:"keyword"`
	With             *bool                   `ddl:"keyword" sql:"WITH"`
	Tags             []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
	InlineConstraint *ColumnInlineConstraint `ddl:"keyword"`
}

type ColumnDefaultValue struct {
	// One of
	Expression *string         `ddl:"parameter,no_equals" sql:"DEFAULT"`
	Identity   *ColumnIdentity `ddl:"keyword" sql:"IDENTITY"`
}
type ColumnIdentity struct {
	Start     int `ddl:"parameter,no_quotes,no_equals" sql:"START"`
	Increment int `ddl:"parameter,no_quotes,no_equals" sql:"INCREMENT"`
}

type ColumnMaskingPolicy struct {
	With          *bool                  `ddl:"keyword" sql:"WITH"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	Name          SchemaObjectIdentifier `ddl:"identifier"`
	Using         []string               `ddl:"keyword,parentheses" sql:"USING"`
}

type CreateOutOfLineConstraint struct {
	Name       string               `ddl:"parameter,no_equals" sql:", CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	Columns    []string             `ddl:"keyword,parentheses"`
	ForeignKey *OutOfLineForeignKey `ddl:"keyword"`

	// optional
	Enforced           *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced        *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Deferrable         *bool `ddl:"keyword" sql:"DEFERRABLE"`
	NotDeferrable      *bool `ddl:"keyword" sql:"NOT DEFERRABLE"`
	InitiallyDeferred  *bool `ddl:"keyword" sql:"INITIALLY DEFERRED"`
	InitiallyImmediate *bool `ddl:"keyword" sql:"INITIALLY IMMEDIATE"`
	Enable             *bool `ddl:"keyword" sql:"ENABLE"`
	Disable            *bool `ddl:"keyword" sql:"DISABLE"`
	Validate           *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate         *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely               *bool `ddl:"keyword" sql:"RELY"`
	NoRely             *bool `ddl:"keyword" sql:"NORELY"`
}

type AlterOutOfLineConstraint struct {
	Name       string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	Columns    []string             `ddl:"keyword,parentheses"`
	ForeignKey *OutOfLineForeignKey `ddl:"keyword"`

	// optional
	Enforced           *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced        *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Deferrable         *bool `ddl:"keyword" sql:"DEFERRABLE"`
	NotDeferrable      *bool `ddl:"keyword" sql:"NOT DEFERRABLE"`
	InitiallyDeferred  *bool `ddl:"keyword" sql:"INITIALLY DEFERRED"`
	InitiallyImmediate *bool `ddl:"keyword" sql:"INITIALLY IMMEDIATE"`
	Enable             *bool `ddl:"keyword" sql:"ENABLE"`
	Disable            *bool `ddl:"keyword" sql:"DISABLE"`
	Validate           *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate         *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely               *bool `ddl:"keyword" sql:"RELY"`
	NoRely             *bool `ddl:"keyword" sql:"NORELY"`
}

type OutOfLineForeignKey struct {
	references  bool                   `ddl:"static" sql:"REFERENCES"`
	TableName   SchemaObjectIdentifier `ddl:"identifier"`
	ColumnNames []string               `ddl:"parameter,no_equals,parentheses"`
	Match       *MatchType             `ddl:"parameter,no_equals" sql:"MATCH"`
	On          *ForeignKeyOnAction    `ddl:"keyword"`
}

type StageCopyOption struct {
	InnerValue StageCopyOptionsInnerValue `ddl:"keyword"`
}

type StageCopyOptionsInnerValue struct {
	OnError           *StageCopyOnErrorOptions  `ddl:"parameter" sql:"ON_ERROR"`
	SizeLimit         *int                      `ddl:"parameter" sql:"SIZE_LIMIT"`
	Purge             *bool                     `ddl:"parameter" sql:"PURGE"`
	ReturnFailedOnly  *bool                     `ddl:"parameter" sql:"RETURN_FAILED_ONLY"`
	MatchByColumnName *StageCopyColumnMapOption `ddl:"parameter" sql:"MATCH_BY_COLUMN_NAME"`
	EnforceLength     *bool                     `ddl:"parameter" sql:"ENFORCE_LENGTH"`
	TruncateColumns   *bool                     `ddl:"parameter" sql:"TRUNCATECOLUMNS"`
	Force             *bool                     `ddl:"parameter" sql:"FORCE"`
}

type alterTableOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	table    bool                   `ddl:"static" sql:"TABLE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	NewName                   *SchemaObjectIdentifier        `ddl:"identifier" sql:"RENAME TO"`
	SwapWith                  *SchemaObjectIdentifier        `ddl:"identifier" sql:"SWAP WITH"`
	ClusteringAction          *TableClusteringAction         `ddl:"keyword"`
	ColumnAction              *TableColumnAction             `ddl:"keyword"`
	ConstraintAction          *TableConstraintAction         `ddl:"keyword"`
	ExternalTableAction       *TableExternalTableAction      `ddl:"keyword"`
	SearchOptimizationAction  *TableSearchOptimizationAction `ddl:"keyword"`
	Set                       *TableSet                      `ddl:"keyword" sql:"SET"`
	SetTags                   []TagAssociation               `ddl:"parameter,no_equals" sql:"SET TAG"`
	UnsetTags                 []ObjectIdentifier             `ddl:"keyword" sql:"UNSET TAG"`
	Unset                     *TableUnset                    `ddl:"keyword" sql:"UNSET"`
	AddRowAccessPolicy        *AddRowAccessPolicy            `ddl:"keyword"`
	DropRowAccessPolicy       *string                        `ddl:"parameter,no_equals" sql:"DROP ROW ACCESS POLICY"`
	DropAndAddRowAccessPolicy *DropAndAddRowAccessPolicy     `ddl:"keyword"`
	DropAllAccessRowPolicies  *bool                          `ddl:"keyword" sql:"DROP ALL ROW ACCESS POLICIES"`
}

type TableClusteringAction struct {
	// one of
	ClusterBy            []string                   `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	Recluster            *TableReclusterAction      `ddl:"keyword" sql:"RECLUSTER"`
	ChangeReclusterState *TableReclusterChangeState `ddl:"keyword"`
	DropClusteringKey    *bool                      `ddl:"keyword" sql:"DROP CLUSTERING KEY"`
}

type TableReclusterAction struct {
	MaxSize   *int    `ddl:"parameter" sql:"MAX_SIZE"`
	Condition *string `ddl:"parameter,no_equals" sql:"WHERE"`
}

type TableReclusterChangeState struct {
	State     *ReclusterState `ddl:"keyword"`
	recluster bool            `ddl:"static" sql:"RECLUSTER"`
}

type ReclusterState string

const (
	ReclusterStateResume  ReclusterState = "RESUME"
	ReclusterStateSuspend ReclusterState = "SUSPEND"
)

type TableColumnAction struct {
	// One of
	Add                *TableColumnAddAction                     `ddl:"keyword" sql:"ADD"`
	Rename             *TableColumnRenameAction                  `ddl:"keyword"`
	Alter              []TableColumnAlterAction                  `ddl:"keyword" sql:"ALTER"`
	SetMaskingPolicy   *TableColumnAlterSetMaskingPolicyAction   `ddl:"keyword"`
	UnsetMaskingPolicy *TableColumnAlterUnsetMaskingPolicyAction `ddl:"keyword"`
	SetTags            *TableColumnAlterSetTagsAction            `ddl:"keyword"`
	UnsetTags          *TableColumnAlterUnsetTagsAction          `ddl:"keyword"`
	DropColumns        *TableColumnAlterDropColumns              `ddl:"keyword"`
}

type TableColumnAddAction struct {
	Column           *bool                           `ddl:"keyword" sql:"COLUMN"`
	Name             string                          `ddl:"keyword"`
	Type             DataType                        `ddl:"keyword"`
	DefaultValue     *ColumnDefaultValue             `ddl:"keyword"`
	InlineConstraint *TableColumnAddInlineConstraint `ddl:"keyword"`
	MaskingPolicy    *ColumnMaskingPolicy            `ddl:"keyword"`
	With             *bool                           `ddl:"keyword" sql:"WITH"`
	Tags             []TagAssociation                `ddl:"keyword,parentheses" sql:"TAG"`
}

type TableColumnAddInlineConstraint struct {
	NotNull    *bool                `ddl:"keyword" sql:"NOT NULL"`
	Name       string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *ColumnAddForeignKey `ddl:"keyword"`
}

type ColumnAddForeignKey struct {
	TableName  string `ddl:"keyword" sql:"REFERENCES"`
	ColumnName string `ddl:"keyword,parentheses"`
}

type TableColumnRenameAction struct {
	OldName string `ddl:"parameter,no_equals" sql:"RENAME COLUMN"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableColumnAlterAction struct {
	Column *bool  `ddl:"keyword" sql:"COLUMN"`
	Name   string `ddl:"keyword"`

	// One of
	DropDefault       *bool         `ddl:"keyword" sql:"DROP DEFAULT"`
	SetDefault        *SequenceName `ddl:"parameter,no_equals" sql:"SET DEFAULT"`
	NotNullConstraint *TableColumnNotNullConstraint
	Type              *DataType `ddl:"parameter,no_equals" sql:"SET DATA TYPE"`
	Comment           *string   `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
	UnsetComment      *bool     `ddl:"keyword" sql:"UNSET COMMENT"`
}

type TableColumnAlterSetMaskingPolicyAction struct {
	alter             bool                   `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName        string                 `ddl:"keyword"`
	setMaskingPolicy  bool                   `ddl:"static" sql:"SET MASKING POLICY"`
	MaskingPolicyName SchemaObjectIdentifier `ddl:"identifier"`
	Using             []string               `ddl:"keyword,parentheses" sql:"USING"`
	Force             *bool                  `ddl:"keyword" sql:"FORCE"`
}

type TableColumnAlterUnsetMaskingPolicyAction struct {
	alter            bool   `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName       string `ddl:"keyword"`
	setMaskingPolicy bool   `ddl:"static" sql:"UNSET MASKING POLICY"`
}

type TableColumnAlterSetTagsAction struct {
	alter      bool             `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName string           `ddl:"keyword"`
	set        bool             `ddl:"static" sql:"SET"`
	Tags       []TagAssociation `ddl:"keyword" sql:"TAG"`
}

type TableColumnAlterUnsetTagsAction struct {
	alter      bool               `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName string             `ddl:"keyword"`
	unset      bool               `ddl:"static" sql:"UNSET"`
	Tags       []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

type TableColumnAlterDropColumns struct {
	dropColumn bool     `ddl:"static" sql:"DROP COLUMN"`
	Columns    []string `ddl:"keyword"`
}

type TableColumnAlterSequenceName interface {
	String() string
}

type SequenceName string

func (sn SequenceName) String() string {
	return fmt.Sprintf("%s.NEXTVAL", string(sn))
}

type TableColumnNotNullConstraint struct {
	Set  *bool `ddl:"keyword" sql:"SET NOT NULL"`
	Drop *bool `ddl:"keyword" sql:"DROP NOT NULL"`
}

type TableConstraintAction struct {
	Add    *AlterOutOfLineConstraint    `ddl:"keyword" sql:"ADD"`
	Rename *TableConstraintRenameAction `ddl:"keyword" sql:"RENAME CONSTRAINT"`
	Alter  *TableConstraintAlterAction  `ddl:"keyword" sql:"ALTER"`
	Drop   *TableConstraintDropAction   `ddl:"keyword" sql:"DROP"`
}

type TableConstraintRenameAction struct {
	OldName string `ddl:"keyword"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableConstraintAlterAction struct {
	// One of
	ConstraintName *string  `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	PrimaryKey     *bool    `ddl:"keyword" sql:"PRIMARY KEY"`
	Unique         *bool    `ddl:"keyword" sql:"UNIQUE"`
	ForeignKey     *bool    `ddl:"keyword" sql:"FOREIGN KEY"`
	Columns        []string `ddl:"keyword,parentheses"`

	// Optional
	Enforced    *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Validate    *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate  *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely        *bool `ddl:"keyword" sql:"RELY"`
	NoRely      *bool `ddl:"keyword" sql:"NORELY"`
}

type TableConstraintDropAction struct {
	// One of
	ConstraintName *string  `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	PrimaryKey     *bool    `ddl:"keyword" sql:"PRIMARY KEY"`
	Unique         *bool    `ddl:"keyword" sql:"UNIQUE"`
	ForeignKey     *bool    `ddl:"keyword" sql:"FOREIGN KEY"`
	Columns        []string `ddl:"keyword,parentheses"`

	// Optional
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
}

type TableUnsetTags struct {
	Tag []ObjectIdentifier `ddl:"keyword"`
}

type TableExternalTableAction struct {
	// One of
	Add    *TableExternalTableColumnAddAction    `ddl:"keyword"`
	Rename *TableExternalTableColumnRenameAction `ddl:"keyword"`
	Drop   *TableExternalTableColumnDropAction   `ddl:"keyword"`
}

type TableExternalTableColumnAddAction struct {
	addColumn  bool     `ddl:"static" sql:"ADD COLUMN"`
	Name       string   `ddl:"keyword"`
	Type       DataType `ddl:"keyword"`
	Expression []string `ddl:"parameter,no_equals,parentheses" sql:"AS"`
}

type TableExternalTableColumnRenameAction struct {
	OldName string `ddl:"parameter,no_equals" sql:"RENAME COLUMN"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableExternalTableColumnDropAction struct {
	Columns []string `ddl:"keyword" sql:"DROP COLUMN"`
}

type TableSearchOptimizationAction struct {
	// One of
	Add  *AddSearchOptimization  `ddl:"keyword"`
	Drop *DropSearchOptimization `ddl:"keyword"`
}

type AddSearchOptimization struct {
	addSearchOptimization bool `ddl:"static" sql:"ADD SEARCH OPTIMIZATION"`
	// Optional
	On []string `ddl:"keyword" sql:"ON"`
}

type DropSearchOptimization struct {
	dropSearchOptimization bool `ddl:"static" sql:"DROP SEARCH OPTIMIZATION"`
	// Optional
	On []string `ddl:"keyword" sql:"ON"`
}

type TableSet struct {
	// Optional
	EnableSchemaEvolution      *bool             `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            []StageFileFormat `ddl:"parameter,equals,parentheses" sql:"STAGE_FILE_FORMAT"`
	StageCopyOptions           []StageCopyOption `ddl:"parameter,equals,parentheses" sql:"STAGE_COPY_OPTIONS"`
	DataRetentionTimeInDays    *int              `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int              `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool             `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string           `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type TableUnset struct {
	DataRetentionTimeInDays    *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool `ddl:"keyword" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	EnableSchemaEvolution      *bool `ddl:"keyword" sql:"ENABLE_SCHEMA_EVOLUTION"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
}

type AddRowAccessPolicy struct {
	PolicyName  string   `ddl:"parameter,no_equals" sql:"ADD ROW ACCESS POLICY"`
	ColumnNames []string `ddl:"parameter,no_equals,parentheses" sql:"ON"`
}

type DropAndAddRowAccessPolicy struct {
	DroppedPolicyName string              `ddl:"parameter,no_equals" sql:"DROP ROW ACCESS POLICY"`
	comma             bool                `ddl:"static" sql:","`
	AddedPolicy       *AddRowAccessPolicy `ddl:"keyword"`
}

// dropTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-table
type dropTableOptions struct {
	drop         bool                   `ddl:"static" sql:"DROP"`
	databaseRole bool                   `ddl:"static" sql:"TABLE"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
}

type showTableOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	tables     bool       `ddl:"static" sql:"TABLES"`
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type tableDBRow struct {
	CreatedOn                  string         `db:"created_on"`
	Name                       string         `db:"name"`
	SchemaName                 string         `db:"schema_name"`
	DatabaseName               string         `db:"database_name"`
	Kind                       string         `db:"kind"`
	Comment                    sql.NullString `db:"comment"`
	ClusterBy                  sql.NullString `db:"cluster_by"`
	Rows                       int            `db:"rows"`
	Owner                      string         `db:"owner"`
	RetentionTime              int            `db:"retention_time"`
	AutomaticClustering        sql.NullString `db:"automatic_clustering"`
	ChangeTracking             sql.NullString `db:"change_tracking"`
	SearchOptimization         sql.NullString `db:"search_optimization"`
	SearchOptimizationProgress sql.NullString `db:"search_optimization_progress"`
	IsExternal                 sql.NullString `db:"is_external"`
	EnableSchemaEvolution      sql.NullString `db:"enable_schema_evolution"`
	OwnerRoleType              sql.NullString `db:"owner_role_type"`
	IsEvent                    sql.NullString `db:"is_event"`
}

type Table struct {
	CreatedOn                  string
	Name                       string
	DatabaseName               string
	SchemaName                 string
	Kind                       string
	Comment                    string
	ClusterBy                  string
	Rows                       int
	Owner                      string
	RetentionTime              int
	AutomaticClustering        bool
	ChangeTracking             bool
	SearchOptimization         bool
	SearchOptimizationProgress string
	IsExternal                 bool
	EnableSchemaEvolution      bool
	OwnerRoleType              string
	IsEvent                    bool
}

func (row tableDBRow) convert() *Table {
	databaseRole := Table{
		CreatedOn:     row.CreatedOn,
		Name:          row.Name,
		SchemaName:    row.SchemaName,
		DatabaseName:  row.DatabaseName,
		Rows:          row.Rows,
		Owner:         row.Owner,
		Kind:          row.Kind,
		RetentionTime: row.RetentionTime,
	}
	if row.AutomaticClustering.Valid {
		databaseRole.AutomaticClustering = row.AutomaticClustering.String == "ON"
	}
	if row.ChangeTracking.Valid {
		databaseRole.ChangeTracking = row.ChangeTracking.String == "ON"
	}
	if row.SearchOptimization.Valid {
		databaseRole.SearchOptimization = row.SearchOptimization.String == "ON"
	}
	if row.SearchOptimizationProgress.Valid {
		databaseRole.SearchOptimizationProgress = row.SearchOptimizationProgress.String
	}
	if row.IsExternal.Valid {
		databaseRole.IsExternal = row.IsExternal.String == "Y"
	}
	if row.IsEvent.Valid {
		databaseRole.IsEvent = row.IsEvent.String == "Y"
	}
	if row.EnableSchemaEvolution.Valid {
		databaseRole.EnableSchemaEvolution = row.EnableSchemaEvolution.String == "Y"
	}
	if row.Comment.Valid {
		databaseRole.Comment = row.Comment.String
	}
	if row.ClusterBy.Valid {
		databaseRole.ClusterBy = row.ClusterBy.String
	}
	if row.OwnerRoleType.Valid {
		databaseRole.OwnerRoleType = row.OwnerRoleType.String
	}
	return &databaseRole
}

func (v *Table) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Table) ObjectType() ObjectType {
	return ObjectTypeTable
}
