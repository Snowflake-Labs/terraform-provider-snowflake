package sdk

import "fmt"

type TableCreateOptions struct {
	create                     bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace                  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Scope                      *TableScope            `ddl:"keyword"`
	Kind                       *TableKind             `ddl:"keyword"`
	table                      bool                   `ddl:"static" sql:"TABLE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists                *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier `ddl:"identifier"`
	leftParen                  bool                   `ddl:"static" sql:"("` //lint:ignore U1000 This is used in the ddl tag
	Columns                    []TableColumn          `ddl:"keyword"`
	OutOfLineConstraint        *OutOfLineConstraint   `ddl:"keyword"`
	rightParen                 bool                   `ddl:"static" sql:")"` //lint:ignore U1000 This is used in the ddl tag
	ClusterBy                  []string               `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	EnableSchemaEvolution      *bool                  `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            []StageFileFormat      `ddl:"parameter,equals,parentheses" sql:"STAGE_FILE_FORMAT"`
	StageCopyOptions           []StageCopyOptions     `ddl:"parameter,equals,parentheses" sql:"STAGE_COPY_OPTIONS"`
	DataRetentionTimeInDays    *int                   `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataRetentionTimeInDays *int                   `ddl:"parameter" sql:"MAX_DATA_RETENTION_TIME_IN_DAYS"`
	ChangeTracking             *bool                  `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string                `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	CopyGrants                 *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy            *RowAccessPolicy       `ddl:"keyword"`
	Tags                       []TagAssociation       `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type RowAccessPolicy struct {
	With            *bool                  `ddl:"keyword" sql:"WITH"`
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"` //lint:ignore U1000 This is used in the ddl tag
	Name            SchemaObjectIdentifier `ddl:"identifier"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

type TableScope string

const (
	GlobalTableScope = "GLOBAL"
	LocalTableScope  = "LOCAL"
)

type TableKind string

const (
	TempTableKind      = "TEMP"
	TemporaryTableKind = "TEMPORARY"
	VolatileTableKind  = "VOLATILE"
	TransientTableKind = "TRANSIENT"
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
	//one of
	Expression *string         `ddl:"keyword" sql:"DEFAULT"`
	Identity   *ColumnIdentity `ddl:"keyword" sql:"IDENTITY"`
}
type ColumnIdentity struct {
	Start     int `ddl:"parameter,no_quotes,no_equals" sql:"START"`
	Increment int `ddl:"parameter,no_quotes,no_equals" sql:"INCREMENT"`
}

type ColumnMaskingPolicy struct {
	With          *bool                  `ddl:"keyword" sql:"WITH"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"` //lint:ignore U1000 This is used in the ddl tag
	Name          SchemaObjectIdentifier `ddl:"identifier"`
	Using         []string               `ddl:"keyword,parentheses" sql:"USING"`
}

type ColumnInlineConstraint struct {
	Name       string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey    `ddl:"keyword" sql:"FOREIGN KEY"`

	//optional
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

type OutOfLineConstraint struct {
	Name       string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	Columns    []string             `ddl:"keyword,parentheses"`
	ForeignKey *OutOfLineForeignKey `ddl:"keyword"`

	//optional
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

type ColumnConstraintType string

const (
	ColumnConstraintTypeUnique     ColumnConstraintType = "UNIQUE"
	ColumnConstraintTypePrimaryKey ColumnConstraintType = "PRIMARY KEY"
	ColumnConstraintTypeForeignKey ColumnConstraintType = "FOREIGN KEY"
)

type InlineForeignKey struct {
	TableName  string              `ddl:"keyword" sql:"REFERENCES"`
	ColumnName []string            `ddl:"keyword,parentheses"`
	Match      *MatchType          `ddl:"keyword" sql:"MATCH"`
	On         *ForeignKeyOnAction `ddl:"keyword" sql:"ON"`
}

type OutOfLineForeignKey struct {
	TableName   SchemaObjectIdentifier `ddl:"identifier" sql:"REFERENCES"`
	ColumnNames []string               `ddl:"keyword,parentheses"`
	Match       *MatchType             `ddl:"parameter,no_equals" sql:"MATCH"`
	On          *ForeignKeyOnAction    `ddl:"keyword"`
}

type MatchType string

const (
	FullMatchType    MatchType = "FULL"
	SimpleMatchType  MatchType = "SIMPLE"
	PartialMatchType MatchType = "PARTIAL"
)

type ForeignKeyOnAction struct {
	OnUpdate *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

type ForeignKeyAction string

const (
	ForeignKeyCascadeAction    ForeignKeyAction = "CASCADE"
	ForeignKeySetNullAction    ForeignKeyAction = "SET NULL"
	ForeignKeySetDefaultAction ForeignKeyAction = "SET DEFAULT"
	ForeignKeyRestrictAction   ForeignKeyAction = "RESTRICT"
	ForeignKeyNoAction         ForeignKeyAction = "NO ACTION"
)

type StageFileFormat struct {
	InnerValue StageFileFormatInnerValue `ddl:"keyword"`
}

type StageFileFormatInnerValue struct {
	//one of
	FormatName *string         `ddl:"parameter,no_quotes" sql:"FORMAT_NAME"`
	FormatType *FileFormatType `ddl:"parameter" sql:"TYPE"`

	Options *FileFormatTypeOptions
}

type StageCopyOptions struct {
	InnerValue StageCopyOptionsInnerValue `ddl:"keyword"`
}
type StageCopyOptionsInnerValue struct {
	OnError           StageCopyOptionsOnError            `ddl:"parameter" sql:"ON_ERROR"`
	SizeLimit         *int                               `ddl:"parameter" sql:"SIZE_LIMIT"`
	Purge             *bool                              `ddl:"parameter" sql:"PURGE"`
	ReturnFailedOnly  *bool                              `ddl:"parameter" sql:"RETURN_FAILED_ONLY"`
	MatchByColumnName *StageCopyOptionsMatchByColumnName `ddl:"parameter" sql:"MATCH_BY_COLUMN_NAME"`
	EnforceLength     *bool                              `ddl:"parameter" sql:"ENFORCE_LENGTH"`
	TruncateColumns   *bool                              `ddl:"parameter" sql:"TRUNCATECOLUMNS"`
	Force             *bool                              `ddl:"parameter" sql:"FORCE"`
}

type StageCopyOptionsOnError interface {
	stageCopyOptionsOnError()
	String() string
}

type StageCopyOptionsOnErrorContinue struct{}

func (StageCopyOptionsOnErrorContinue) stageCopyOptionsOnError() {}
func (StageCopyOptionsOnErrorContinue) String() string {
	return "CONTINUE"
}

type StageCopyOptionsOnErrorSkipFile struct{}

func (StageCopyOptionsOnErrorSkipFile) stageCopyOptionsOnError() {}

func (StageCopyOptionsOnErrorSkipFile) String() string {
	return "SKIP_FILE"
}

type StageCopyOptionsOnErrorSkipFileNum struct {
	Value int 
}

func (StageCopyOptionsOnErrorSkipFileNum) stageCopyOptionsOnError() {}

func (opt StageCopyOptionsOnErrorSkipFileNum) String() string {
	return fmt.Sprintf("SKIP_FILE_%v", opt.Value)
}

type StageCopyOptionsOnErrorSkipFileNumPercentage struct {
	Value int
}

func (StageCopyOptionsOnErrorSkipFileNumPercentage) stageCopyOptionsOnError() {}

func (opt StageCopyOptionsOnErrorSkipFileNumPercentage) String() string {
	return fmt.Sprintf("'SKIP_FILE_%d%%'", opt.Value)
}

type StageCopyOptionsOnErrorAbortStatement struct{}

func (StageCopyOptionsOnErrorAbortStatement) stageCopyOptionsOnError() {}

func (StageCopyOptionsOnErrorAbortStatement) String() string {
	return "ABORT_STATEMENT"
}

type StageCopyOptionsMatchByColumnName string

const (
	CopyOptionsMatchByColumnNameCaseSensitive   StageCopyOptionsMatchByColumnName = "CASE_SENSITIVE"
	CopyOptionsMatchByColumnNameCaseInsensitive StageCopyOptionsMatchByColumnName = "CASE_INSENSITIVE"
	CopyOptionsMatchByColumnNameNone            StageCopyOptionsMatchByColumnName = "NONE"
)

type TableAlterOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	table    bool                    `ddl:"static" sql:"TABLE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier  `ddl:"identifier"`
	NewName  SchemaObjectIdentifier  `ddl:"identifier" sql:"RENAME TO"`
	SwapWith AccountObjectIdentifier `ddl:"identifier" sql:"SWAP WITH"`

	Set   *MaskingPolicySet   `ddl:"keyword" sql:"SET"`
	Unset *MaskingPolicyUnset `ddl:"keyword" sql:"UNSET"`
}

type TableClusteringAction struct {
	//one of
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
	State     ReclusterState `ddl:"keyword"`
	recluster bool           `ddl:"static" sql:"RECLUSTER"` //lint:ignore U1000 This is used in the ddl tag
}

type ReclusterState string

const (
	ReclusterStateResume  ReclusterState = "RESUME"
	ReclusterStateSuspend ReclusterState = "SUSPEND"
)

type TableColumnAction struct {
	Add *TableColumnAddAction `ddl:"keyword" sql:"ADD"`
}

type TableColumnAddAction struct {
	Column       *bool               `ddl:"keyword" sql:"COLUMN"`
	Name         string              `ddl:"keyword"`
	Type         DataType            `ddl:"keyword"`
	DefaultValue *ColumnDefaultValue `ddl:"keyword"`
}

type TableColumnAddInlineConstraint struct {
	NotNull       *bool                `ddl:"keyword" sql:"NOT NULL"`
	Name          string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type          ColumnConstraintType `ddl:"keyword"`
	ForeignKey    *ColumnAddForeignKey `ddl:"keyword"`
	MaskingPolicy *ColumnMaskingPolicy `ddl:"keyword"`
	With          *bool                `ddl:"keyword" sql:"WITH"`
	Tags          []TagAssociation     `ddl:"keyword,parentheses" sql:"TAG"`
}
type ColumnAddForeignKey struct {
	TableName  string `ddl:"keyword" sql:"REFERENCES"`
	ColumnName string `ddl:"keyword,parentheses"`
}

type TableColumnRenameAction struct {
	OldName string `ddl:"keyword" sql:"RENAME COLUMN"`
	NewName string `ddl:"keyword" sql:"TO"`
}

type TableColumnAlterAction struct {
	// alter bool `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	Column      *bool        `ddl:"keyword" sql:"COLUMN"`
	Name        string       `ddl:"keyword"`
	DropDefault *bool        `ddl:"keyword" sql:"DROP DEFAULT"`
	SetDefault  SequenceName `ddl:"keyword" sql:"SET DEFAULT"`
}

type SequenceNameInterfaceFooChangeName interface {
	String() string
}
type SequenceName string

func (sn SequenceName) String() string {
	return fmt.Sprintf("%s.NEXTVAL", string(sn))
}

type TableSet struct {
	EnableSchemaEvolution      *bool              `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            []StageFileFormat  `ddl:"parameter,equals,parentheses" sql:"STAGE_FILE_FORMAT"`
	StageCopyOptions           []StageCopyOptions `ddl:"parameter,equals,parentheses" sql:"STAGE_COPY_OPTIONS"`
	DataRetentionTimeInDays    *int               `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataRetentionTimeInDays *int               `ddl:"parameter" sql:"MAX_DATA_RETENTION_TIME_IN_DAYS"`
	ChangeTracking             *bool              `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string            `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string            `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type TableUnset struct {
}

type Table struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func (v *Table) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Table) ObjectType() ObjectType {
	return ObjectTypeTable
}
