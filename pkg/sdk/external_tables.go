package sdk

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
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

var _ ExternalTables = (*externalTables)(nil)

type externalTables struct {
	client *Client
}

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
	name                AccountObjectIdentifier   `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
	Columns             []ExternalTableColumn     `ddl:"list,parentheses"`
	CloudProviderParams CloudProviderParams       // TODO Not required and used for notifications
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

func (opts *CreateExternalTableOpts) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Columns) {
		errs = append(errs, errors.New("no column provided")) // TODO message
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errors.New("no file format provided")) // TODO message
	}
	// TODO call fields validate functions
	return errors.Join(errs...)
}

type ExternalTableColumn struct {
	Name             string   `ddl:"keyword"`
	Type             DataType `ddl:"keyword"`
	AsExpression     string   `ddl:"parameter,parentheses,no_equals" sql:"AS"`
	InlineConstraint *ExternalTableInlineConstraint
}

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

type ForeignKeyOnAction struct {
	OnUpdate *bool `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *bool `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

type CloudProviderParams struct {
	// One of
	GoogleCloudStorage *GoogleCloudStorageParams
	MicrosoftAzure     *MicrosoftAzureParams
}

func (cpp *CloudProviderParams) validate() error {
	// TODO
	if anyValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) && exactlyOneValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) {
		return errors.New("")
	}
	return nil
}

type GoogleCloudStorageParams struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type MicrosoftAzureParams struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type ExternalTableFileFormat struct {
	Name *string                      `ddl:"parameter,single_quotes" sql:"FORMAT_NAME"`
	Type *ExternalTableFileFormatType `ddl:"parameter" sql:"TYPE"`
	// TODO: Should be probably a new type because doesn't contain xml (or maybe FileFormatType should be divided into struct for every file format)
	Options *FileFormatTypeOptions
}

func (opts *ExternalTableFileFormat) validate() error {
	// TODO error message
	if valueSet(opts.Name) && anyValueSet(opts.Type, opts.Options) {
		return errors.New("")
	}
	return nil
}

type ExternalTableFileFormatType string

var (
	ExternalTableFileFormatTypeCSV     ExternalTableFileFormatType = "CSV"
	ExternalTableFileFormatTypeJSON    ExternalTableFileFormatType = "JSON"
	ExternalTableFileFormatTypeAvro    ExternalTableFileFormatType = "AVRO"
	ExternalTableFileFormatTypeORC     ExternalTableFileFormatType = "ORC"
	ExternalTableFileFormatTypeParquet ExternalTableFileFormatType = "PARQUET"
)

type RowAccessPolicy struct {
	Name SchemaObjectIdentifier `ddl:"identifier"`
	On   []string               `ddl:"keyword,parentheses" sql:"ON"` // TODO What is correct (quoted values or no)
}

func (v *externalTables) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateExternalTableOpts) error {
	if opts == nil {
		opts = &CreateExternalTableOpts{}
	}
	opts.name = id
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

func (opts *CreateWithManualPartitioningExternalTableOpts) validate() error {
	return nil
}

// TODO: ChangeName
func (v *externalTables) CreateWithManualPartitioning(ctx context.Context, id AccountObjectIdentifier, opts *CreateWithManualPartitioningExternalTableOpts) error {
	if opts == nil {
		opts = &CreateWithManualPartitioningExternalTableOpts{}
	}
	opts.name = id
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

func (opts *CreateDeltaLakeExternalTableOpts) validate() error {
	return nil
}

func (v *externalTables) CreateDeltaLake(ctx context.Context, id AccountObjectIdentifier, opts *CreateDeltaLakeExternalTableOpts) error {
	if opts == nil {
		opts = &CreateDeltaLakeExternalTableOpts{}
	}
	opts.name = id
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
	Path *string `ddl:"parameter,no_equals,single_quotes"`
}

type ExternalTableFile struct {
	Name string `ddl:"keyword,single_quotes"`
}

func (opts *AlterExternalTableOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if anyValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.Set, opts.Unset) &&
		!exactlyOneValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.Set, opts.Unset) {
		return errors.New("") // TODO
	}
	return nil
}

type ExternalTableRefresh struct {
	refresh      bool    `ddl:"static" sql:"REFRESH"`
	RelativePath *string `ddl:"parameter,single_quote"`
}

type ExternalTableSet struct {
	AutoRefresh *bool            `ddl:"parameter" sql:"AUTO_REFRESH"`
	Tag         []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalTableUnset struct {
	Tag []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *externalTables) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTableOptions) error {
	if opts == nil {
		opts = &AlterExternalTableOptions{}
	}
	opts.name = id
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

type AlterExternalTablePartitionOptions struct {
	alterExternalTable bool                    `ddl:"static" sql:"ALTER EXTERNAL TABLE"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	IfExists           *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	AddPartitions      []Partition             `ddl:"keyword,parentheses" sql:"ADD PARTITION"`
	DropPartition      *bool                   `ddl:"keyword" sql:"DROP PARTITION"`
	Location           string                  `ddl:"keyword,single_quotes" sql:"LOCATION"`
}

func (opts *AlterExternalTablePartitionOptions) validate() error {
	// TODO identifier etc
	return nil
}

type Partition struct {
	ColumnName string `ddl:"keyword,double_quotes"`
	Value      string `ddl:"parameter,single_quotes"`
}

func (v *externalTables) AlterPartitions(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTablePartitionOptions) error {
	if opts == nil {
		opts = &AlterExternalTablePartitionOptions{}
	}
	opts.name = id
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

type DropExternalTableOptions struct {
	dropExternalTable bool                    `ddl:"static" sql:"DROP EXTERNAL TABLE"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	DropOption        *ExternalTableDropOption
}

func (opts *DropExternalTableOptions) validate() error {
	if valueSet(opts.DropOption) {
		if err := opts.DropOption.validate(); err != nil {
			return err
		}
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type ExternalTableDropOption struct {
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
}

func (opts *ExternalTableDropOption) validate() error {
	if anyValueSet(opts.Restrict, opts.Cascade) && !exactlyOneValueSet(opts.Restrict, opts.Cascade) {
		return errors.New("") // TODO error message
	}
	return nil
}

func (v *externalTables) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropExternalTableOptions) error {
	if opts == nil {
		opts = &DropExternalTableOptions{}
	}
	opts.name = id
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

type ShowExternalTableOptions struct {
	show           bool       `ddl:"static" sql:"SHOW"`
	Terse          *bool      `ddl:"keyword" sql:"TERSE"`
	externalTables bool       `ddl:"static" sql:"EXTERNAL TABLES"`
	Like           *Like      `ddl:"keyword" sql:"LIKE"`
	In             *In        `ddl:"keyword" sql:"IN"`
	StartsWith     *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowExternalTableOptions) validate() error {
	return nil
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

func setIfNotNil[T any](value *T, valuer driver.Valuer) {
	if v, err := valuer.Value(); err == nil {
		typedValue := v.(T)
		value = &typedValue
	}
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
	setIfNotNil(&et.InvalidReason, e.InvalidReason)
	setIfNotNil(&et.NotificationChannel, e.NotificationChannel)
	setIfNotNil(&et.LastRefreshedOn, e.LastRefreshedOn)
	setIfNotNil(&et.LastRefreshDetails, e.LastRefreshDetails)
	return et
}

func (v *externalTables) Show(ctx context.Context, opts *ShowExternalTableOptions) ([]ExternalTable, error) {
	if opts == nil {
		opts = &ShowExternalTableOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []externalTableRow
	err = v.client.query(ctx, &rows, sql)
	externalTables := make([]ExternalTable, len(rows))
	for i, row := range rows {
		externalTables[i] = row.ToExternalTable()
	}
	return externalTables, err
}

func (v *externalTables) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalTable, error) {
	if !validObjectidentifier(id) {
		return nil, ErrInvalidObjectIdentifier
	}
	externalTables, err := v.client.ExternalTables.Show(ctx, &ShowExternalTableOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, t := range externalTables {
		if t.ID() == id {
			return &t, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type describeExternalTableColumns struct {
	describeExternalTable bool                    `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"` //lint:ignore U1000 This is used in the ddl tag
	name                  AccountObjectIdentifier `ddl:"identifier"`
	columnsType           bool                    `ddl:"static" sql:"TYPE = COLUMNS"`
}

func (v *describeExternalTableColumns) validate() error {
	if validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
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
	IsNullable bool           `db:"null?"`
	Default    sql.NullString `db:"default"`
	IsPrimary  bool           `db:"primary key"`
	IsUnique   bool           `db:"unique key"`
	Check      sql.NullBool   `db:"check"` // ? BOOL ?
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
	PolicyName sql.NullString `db:"policy name"`
}

func (r *externalTableColumnDetailsRow) toExternalTableColumnDetails() ExternalTableColumnDetails {
	details := ExternalTableColumnDetails{
		Name:       r.Name,
		Type:       r.Type,
		Kind:       r.Kind,
		IsNullable: r.IsNullable,
		IsPrimary:  r.IsPrimary,
		IsUnique:   r.IsUnique,
	}
	setIfNotNil(&details.Default, r.Default)
	setIfNotNil(&details.Check, r.Check)
	setIfNotNil(&details.Expression, r.Expression)
	setIfNotNil(&details.Comment, r.Comment)
	setIfNotNil(&details.PolicyName, r.PolicyName)
	return details
}

func (v *externalTables) DescribeColumns(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableColumnDetails, error) {
	query := describeExternalTableColumns{
		name: id,
	}
	if err := query.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(query)
	if err != nil {
		return nil, err
	}

	var rows []externalTableColumnDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	var result []ExternalTableColumnDetails
	for _, r := range rows {
		result = append(result, r.toExternalTableColumnDetails())
	}
	return result, nil
}

type describeExternalTableStage struct {
	describeExternalTable bool                    `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"` //lint:ignore U1000 This is used in the ddl tag
	name                  AccountObjectIdentifier `ddl:"identifier"`
	stageType             bool                    `ddl:"static" sql:"TYPE = STAGE"`
}

func (v *describeExternalTableStage) validate() error {
	if validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type ExternalTableStageDetails struct {
	parentProperty  string
	property        string
	propertyType    string
	propertyValue   string
	propertyDefault string
}

type externalTableStageDetailsRow struct {
	parentProperty  string `db:"parent_property"`
	property        string `db:"property"`
	propertyType    string `db:"property_type"`
	propertyValue   string `db:"property_value"`
	propertyDefault string `db:"property_default"`
}

func (r externalTableStageDetailsRow) toExternalTableStageDetails() ExternalTableStageDetails {
	return ExternalTableStageDetails{
		parentProperty:  r.parentProperty,
		property:        r.property,
		propertyType:    r.propertyType,
		propertyValue:   r.propertyValue,
		propertyDefault: r.propertyDefault,
	}
}

func (v *externalTables) DescribeStage(ctx context.Context, id AccountObjectIdentifier) ([]ExternalTableStageDetails, error) {
	query := describeExternalTableStage{
		name: id,
	}
	if err := query.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(query)
	if err != nil {
		return nil, err
	}

	var rows []externalTableStageDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	var result []ExternalTableStageDetails
	for _, r := range rows {
		result = append(result, r.toExternalTableStageDetails())
	}

	return result, nil
}
