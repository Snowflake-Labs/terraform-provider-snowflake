package sdk

import (
	"context"
	"database/sql"
	"time"
)

type ExternalTables interface {
	// Create creates an external table with computed partitions.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateExternalTableOpts) error
	// CreateWithManualPartitioning creates an external table where partitions are added and removed manually.
	CreateWithManualPartitioning(ctx context.Context, id AccountObjectIdentifier, opts *CreateWithManualPartitioningExternalTableOpts) error
	// CreateDeltaLake creates a delta lake external table.
	CreateDeltaLake(ctx context.Context, id AccountObjectIdentifier, opts *CreateDeltaLakeExternalTableOpts) error
	// Alter modifies an existing external table.
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTableOptions) error
	// AlterPartitions modifies an existing external table's partitions.
	AlterPartitions(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTablePartitionOptions) error
	// Drop removes an external table.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropExternalTableOptions) error
	// Show returns a list of external tables.
	Show(ctx context.Context, opts *ShowExternalTableOptions) ([]ExternalTable, error)
	// ShowByID returns an external table by ID.
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalTable, error)
	// DescribeColumns returns the details of the external table's columns.
	DescribeColumns(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableColumnDetails, error)
	// DescribeStage returns the details of the external table's stage.
	DescribeStage(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableStageDetails, error)
}

// TODO Infer schema

type ExternalTable struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Invalid             bool
	InvalidReason       string
	Owner               string
	Comment             string
	Stage               string
	Location            string
	FileFormatName      string
	FileFormatType      string
	Cloud               string
	Region              string
	NotificationChannel string
	LastRefreshedOn     time.Time
	TableFormat         string
	LastRefreshDetails  string
	OwnerRoleType       string
}

func (v *ExternalTable) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ExternalTable) ObjectType() ObjectType {
	return ObjectTypeExternalTable
}

type CreateExternalTableOpts struct {
	create              bool                      `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                      `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists         *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier   `ddl:"identifier"`
	Columns             []ExternalTableColumn     `ddl:"list,parentheses"`
	CloudProviderParams *CloudProviderParams      // TODO Not required and used for notifications
	PartitionBy         []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location            string                    `ddl:"parameter" sql:"LOCATION"`
	RefreshOnCreate     *bool                     `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh         *bool                     `ddl:"parameter" sql:"AUTO_REFRESH"`
	Pattern             *string                   `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          []ExternalTableFileFormat `ddl:"keyword,parentheses" sql:"FILE_FORMAT ="` // TODO could be parameter ?
	AwsSnsTopic         *string                   `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	CopyGrants          *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy     *RowAccessPolicy          `ddl:"keyword" sql:"ROW ACCESS POLICY"`
	Tag                 []TagAssociation          `ddl:"keyword,parentheses" sql:"TAG"`
	Comment             *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// TODO Validate
type ExternalTableColumn struct {
	Name             string   `ddl:"keyword"`
	Type             DataType `ddl:"keyword"`
	AsExpression     string   `ddl:"parameter,parentheses,no_equals" sql:"AS"`
	InlineConstraint *ExternalTableInlineConstraint
}

// TODO common type ? + Validate
type ExternalTableInlineConstraint struct {
	NotNull    *bool                 `ddl:"keyword" sql:"NOT NULL"`
	Name       *string               `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       *ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey     `ddl:"keyword" sql:"FOREIGN KEY"`

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

var (
	ColumnConstraintTypeUnique     ColumnConstraintType = "UNIQUE"
	ColumnConstraintTypePrimaryKey ColumnConstraintType = "PRIMARY KEY"
	ColumnConstraintTypeForeignKey ColumnConstraintType = "FOREIGN KEY"
)

// TODO Common Type ? + Validate
type InlineForeignKey struct {
	TableName  string              `ddl:"keyword" sql:"REFERENCES"`
	ColumnName []string            `ddl:"keyword,parentheses"`
	Match      *MatchType          `ddl:"keyword" sql:"MATCH"`
	On         *ForeignKeyOnAction `ddl:"keyword" sql:"ON"`
}

type MatchType string

var (
	FullMatchType    MatchType = "FULL"
	SimpleMatchType  MatchType = "SIMPLE"
	PartialMatchType MatchType = "PARTIAL"
)

// TODO Validate
type ForeignKeyOnAction struct {
	OnUpdate *bool `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *bool `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

// TODO Validate + Rename ? Used mainly for notifications
type CloudProviderParams struct {
	// One of
	GoogleCloudStorage *GoogleCloudStorageParams // TODO Unwrap type ?
	MicrosoftAzure     *MicrosoftAzureParams
}

type GoogleCloudStorageParams struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type MicrosoftAzureParams struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

// TODO New type
type ExternalTableFileFormat struct {
	Name *string                      `ddl:"parameter,single_quotes" sql:"FORMAT_NAME"`
	Type *ExternalTableFileFormatType `ddl:"parameter" sql:"TYPE"`
	// TODO: Should be probably a new type because doesn't contain xml (or maybe FileFormatType should be divided into struct for every file format)
	Options *FileFormatTypeOptions
}

type ExternalTableFileFormatType string

var (
	ExternalTableFileFormatTypeCSV     ExternalTableFileFormatType = "CSV"
	ExternalTableFileFormatTypeJSON    ExternalTableFileFormatType = "JSON"
	ExternalTableFileFormatTypeAvro    ExternalTableFileFormatType = "AVRO"
	ExternalTableFileFormatTypeORC     ExternalTableFileFormatType = "ORC"
	ExternalTableFileFormatTypeParquet ExternalTableFileFormatType = "PARQUET"
)

// TODO Is it common type ?
type RowAccessPolicy struct {
	Name SchemaObjectIdentifier `ddl:"identifier"`
	On   []string               `ddl:"keyword,parentheses" sql:"ON"` // TODO What is correct (quoted values or no)
}

type CreateWithManualPartitioningExternalTableOpts struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	externalTable              bool                      `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier   `ddl:"identifier"`
	Columns                    []ExternalTableColumn     `ddl:"list,parentheses"`
	CloudProviderParams        CloudProviderParams       `ddl:"keyword"`
	PartitionBy                []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location                   string                    `ddl:"parameter" sql:"LOCATION"`
	UserSpecifiedPartitionType *bool                     `ddl:"keyword" sql:"PARTITION_TYPE = USER_SPECIFIED"`
	FileFormat                 []ExternalTableFileFormat `ddl:"keyword,parentheses" sql:"FILE_FORMAT ="` // TODO
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy            *RowAccessPolicy          `ddl:"keyword" sql:"ROW ACCESS POLICY"`
	Tag                        []TagAssociation          `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type CreateDeltaLakeExternalTableOpts struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	externalTable              bool                      `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       AccountObjectIdentifier   `ddl:"identifier"`
	Columns                    []ExternalTableColumn     `ddl:"list,parentheses"`
	CloudProviderParams        CloudProviderParams       `ddl:"keyword"`
	PartitionBy                []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location                   string                    `ddl:"parameter" sql:"LOCATION"`
	refreshOnCreate            bool                      `ddl:"static" sql:"REFRESH_ON_CREATE = FALSE"`
	autoRefresh                bool                      `ddl:"static" sql:"AUTO_REFRESH = FALSE"`
	UserSpecifiedPartitionType *bool                     `ddl:"keyword" sql:"PARTITION_TYPE = USER_SPECIFIED"`
	FileFormat                 []ExternalTableFileFormat `ddl:"keyword,parentheses" sql:"FILE_FORMAT ="` // TODO
	DeltaTableFormat           *bool                     `ddl:"keyword" sql:"TABLE_FORMAT = DELTA"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy            *RowAccessPolicy          `ddl:"keyword" sql:"ROW ACCESS POLICY"`
	Tag                        []TagAssociation          `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AlterExternalTableOptions struct {
	alterExternalTable bool                    `ddl:"static" sql:"ALTER EXTERNAL TABLE"`
	IfExists           *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	// One of
	Refresh     *RefreshExternalTable `ddl:"keyword" sql:"REFRESH"`
	AddFiles    []ExternalTableFile   `ddl:"keyword,no_quotes,parentheses" sql:"ADD FILES"`
	RemoveFiles []ExternalTableFile   `ddl:"keyword,no_quotes,parentheses" sql:"REMOVE FILES"`
	Set         *ExternalTableSet     `ddl:"keyword" sql:"SET"`
	Unset       *ExternalTableUnset   `ddl:"keyword" sql:"UNSET"`
}

type RefreshExternalTable struct {
	Path string `ddl:"parameter,no_equals,single_quotes"`
}

type ExternalTableFile struct {
	Name string `ddl:"keyword,single_quotes"`
}

// TODO Cannot set both ?
type ExternalTableSet struct {
	AutoRefresh *bool            `ddl:"parameter" sql:"AUTO_REFRESH"`
	Tag         []TagAssociation `ddl:"keyword" sql:"TAG"`
}

type ExternalTableUnset struct {
	Tag []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

type AlterExternalTablePartitionOptions struct {
	alterExternalTable bool                    `ddl:"static" sql:"ALTER EXTERNAL TABLE"`
	IfExists           *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	AddPartitions      []Partition             `ddl:"keyword,parentheses" sql:"ADD PARTITION"`
	DropPartition      *bool                   `ddl:"keyword" sql:"DROP PARTITION"`
	Location           string                  `ddl:"parameter,no_equals,single_quotes" sql:"LOCATION"`
}

type Partition struct {
	ColumnName string `ddl:"keyword"`
	Value      string `ddl:"parameter,single_quotes"`
}

type DropExternalTableOptions struct {
	dropExternalTable bool                    `ddl:"static" sql:"DROP EXTERNAL TABLE"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	DropOption        *ExternalTableDropOption
}

type ExternalTableDropOption struct {
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
}

type ShowExternalTableOptions struct {
	show           bool       `ddl:"static" sql:"SHOW"`
	Terse          *bool      `ddl:"keyword" sql:"TERSE"`
	externalTables bool       `ddl:"static" sql:"EXTERNAL TABLES"`
	Like           *Like      `ddl:"keyword" sql:"LIKE"`
	In             *In        `ddl:"keyword" sql:"IN"`
	StartsWith     *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type externalTableRow struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Invalid             bool
	InvalidReason       sql.NullString
	Owner               string
	Comment             string
	Stage               string
	Location            string
	FileFormatName      string
	FileFormatType      string
	Cloud               string
	Region              string
	NotificationChannel sql.NullString
	LastRefreshedOn     sql.NullTime
	TableFormat         string
	LastRefreshDetails  sql.NullString
	OwnerRoleType       string
}

func (e externalTableRow) ToExternalTable() ExternalTable {
	et := ExternalTable{
		CreatedOn:      e.CreatedOn,
		Name:           e.Name,
		DatabaseName:   e.DatabaseName,
		SchemaName:     e.DatabaseName,
		Invalid:        e.Invalid,
		Owner:          e.Owner,
		Stage:          e.Stage,
		Location:       e.Location,
		FileFormatName: e.FileFormatName,
		FileFormatType: e.FileFormatType,
		Cloud:          e.Cloud,
		Region:         e.Region,
		TableFormat:    e.TableFormat,
		OwnerRoleType:  e.OwnerRoleType,
		Comment:        e.Comment,
	}
	if e.InvalidReason.Valid {
		et.InvalidReason = e.InvalidReason.String
	}
	if e.NotificationChannel.Valid {
		et.NotificationChannel = e.NotificationChannel.String
	}
	if e.LastRefreshedOn.Valid {
		et.LastRefreshedOn = e.LastRefreshedOn.Time
	}
	if e.LastRefreshDetails.Valid {
		et.LastRefreshDetails = e.LastRefreshDetails.String
	}
	return et
}

type describeExternalTableColumns struct {
	describeExternalTable bool                    `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"`
	name                  AccountObjectIdentifier `ddl:"identifier"`
	columnsType           bool                    `ddl:"static" sql:"TYPE = COLUMNS"`
}

type ExternalTableColumnDetails struct {
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

type externalTableColumnDetailsRow struct {
	Name       string         `db:"name"`
	Type       DataType       `db:"type"`
	Kind       string         `db:"kind"`
	IsNullable string         `db:"null?"`
	Default    sql.NullString `db:"default"`
	IsPrimary  string         `db:"primary key"`
	IsUnique   string         `db:"unique key"`
	Check      sql.NullBool   `db:"check"` // ? Bool / String ?
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
	PolicyName sql.NullString `db:"policy name"`
}

func (r *externalTableColumnDetailsRow) toExternalTableColumnDetails() ExternalTableColumnDetails {
	details := ExternalTableColumnDetails{
		Name:       r.Name,
		Type:       r.Type,
		Kind:       r.Kind,
		IsNullable: r.IsNullable == "Y",
		IsPrimary:  r.IsPrimary == "Y",
		IsUnique:   r.IsUnique == "Y",
	}
	if r.Default.Valid {
		details.Default = String(r.Default.String)
	}
	if r.Check.Valid {
		details.Check = Bool(r.Check.Bool)
	}
	if r.Expression.Valid {
		details.Expression = String(r.Expression.String)
	}
	if r.Comment.Valid {
		details.Comment = String(r.Comment.String)
	}
	if r.PolicyName.Valid {
		details.PolicyName = String(r.PolicyName.String)
	}
	return details
}

type describeExternalTableStage struct {
	describeExternalTable bool                    `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"`
	name                  AccountObjectIdentifier `ddl:"identifier"`
	stageType             bool                    `ddl:"static" sql:"TYPE = STAGE"`
}

type ExternalTableStageDetails struct {
	ParentProperty  string
	Property        string
	PropertyType    string
	PropertyValue   string
	PropertyDefault string
}

type externalTableStageDetailsRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

func (r externalTableStageDetailsRow) toExternalTableStageDetails() ExternalTableStageDetails {
	return ExternalTableStageDetails{
		ParentProperty:  r.ParentProperty,
		Property:        r.Property,
		PropertyType:    r.PropertyType,
		PropertyValue:   r.PropertyValue,
		PropertyDefault: r.PropertyDefault,
	}
}
