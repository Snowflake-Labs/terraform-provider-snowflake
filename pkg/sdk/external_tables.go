package sdk

import (
	"context"
	"errors"
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
	// Describe returns the details of an external table.
	Describe(ctx context.Context, id AccountObjectIdentifier, opts *DescribeExternalTableOptions) (*ExternalTableDetails, error)
}

var _ ExternalTables = (*externalTables)(nil)

type externalTables struct {
	client *Client
}

type ExternalTable struct {
	Name string
}

func (v *ExternalTable) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ExternalTable) ObjectType() ObjectType {
	return ObjectTypeExternalTable
}

type CreateExternalTableOpts struct {
	create              bool                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                    `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists         *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
	Columns             []ExternalTableColumn   `ddl:"keyword"`
	CloudProviderParams CloudProviderParams     `ddl:"parameter"`
	PartitionBy         []string                `ddl:"list,parentheses" sql:"PARTITION BY"`
	Location            string                  `ddl:"parameter"`
	RefreshOnCreate     *bool                   `ddl:"keyword" sql:"REFRESH_ON_CREATE = TRUE"`
	AutoRefresh         *bool                   `ddl:"keyword" sql:"AUTO_REFRESH = TRUE"`
	Pattern             *string                 `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          ExternalTableFileFormat `ddl:"parameter"`
	AwsSnsTopic         *string                 `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	CopyGrants          *bool                   `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy     *RowAccessPolicy        `ddl:"parameter"`
	Tag                 []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
	Comment             *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateExternalTableOpts) validate() error {
	// TODO
	return nil
}

type ExternalTableColumn struct {
	Name             string
	Type             string
	AsExpression     string
	InlineConstraint *ExternalTableInlineConstraint
}

type ExternalTableInlineConstraint struct {
	NotNull        *bool
	ConstraintName *string
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

type Location struct{}

type ExternalTableFileFormat struct {
	Name    *string
	Type    *ExternalTableFileFormatType
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
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"` //lint:ignore U1000 This is used in the ddl tag
	Name            SchemaObjectIdentifier `ddl:"identifier"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
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
	create              bool                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                    `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists         *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
	Columns             []ExternalTableColumn   `ddl:"keyword"`
	CloudProviderParams CloudProviderParams     `ddl:"parameter"`
	PartitionBy         []string                `ddl:"list,parentheses" sql:"PARTITION BY"`
	Location            Location                `ddl:"parameter"`
	RefreshOnCreate     *bool                   `ddl:"keyword" sql:"REFRESH_ON_CREATE = TRUE"`
	AutoRefresh         *bool                   `ddl:"keyword" sql:"AUTO_REFRESH = TRUE"`
	Pattern             *string                 `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          ExternalTableFileFormat `ddl:"parameter"`
	AwsSnsTopic         *string                 `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	CopyGrants          *bool                   `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy     *RowAccessPolicy        `ddl:"parameter"`
	Tag                 []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
	Comment             *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
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
	name AccountObjectIdentifier
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
	Refresh     *ExternalTableRefresh `ddl:"parameter" sql:"REFRESH"`
	AddFiles    []string              `ddl:"list,parentheses" sql:"ADD FILES"`
	RemoveFiles []string              `ddl:"list" sql:"REMOVE FILES"`
	Set         *ExternalTableSet
	Unset       *ExternalTableUnset
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
	AutoRefresh *bool
	Tag         []TagAssociation
}

type ExternalTableUnset struct {
	Tag []ObjectIdentifier
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
	AddPartitions      []Partition             `ddl:"list,parentheses" sql:"ADD PARTITION"`
	DropPartition      *string                 `ddl:"parameter,single_quotes,no_equals" sql:"DROP PARTITION LOCATION"`
}

func (opts *AlterExternalTablePartitionOptions) validate() error {
	// TODO identifier etc
	return nil
}

type Partition struct {
	ColumnName string `ddl:"parameter"`
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
	dropExternalTable bool                     `ddl:"static" sql:"DROP EXTERNAL TABLE"`
	IfExists          *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier  `ddl:"identifier"`
	DropOption        *ExternalTableDropOption `ddl:"parameter"`
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
	// TODO fill
}

func (etr externalTableRow) ToExternalTable() ExternalTable {
	// TODO
	return ExternalTable{}
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

type DescribeExternalTableOptions struct {
	describeExternalTable bool                    `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"` //lint:ignore U1000 This is used in the ddl tag
	name                  AccountObjectIdentifier `ddl:"identifier"`
	ColumnsType           *bool                   `ddl:"keyword" sql:"TYPE = COLUMNS"`
	StageType             *bool                   `ddl:"keyword" sql:"TYPE = STAGE"`
}

func (opts *DescribeExternalTableOptions) validate() error {
	if anyValueSet(opts.ColumnsType, opts.StageType) && !exactlyOneValueSet(opts.ColumnsType, opts.StageType) {
		return errors.New("") // TODO exactly one value set
	}
	if validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type ExternalTableDetails struct {
	// TODO check if different than ExternalTable
}

type externalTableDetailsRow struct {
	// TODO
	name AccountObjectIdentifier
}

func (r *externalTableDetailsRow) toExternalTableDetails() *ExternalTableDetails {
	// TODO
	return &ExternalTableDetails{}
}

func (v *externalTables) Describe(ctx context.Context, id AccountObjectIdentifier, opts *DescribeExternalTableOptions) (*ExternalTableDetails, error) {
	if opts == nil {
		opts = &DescribeExternalTableOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []externalTableDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.name == id {
			return r.toExternalTableDetails(), nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}
