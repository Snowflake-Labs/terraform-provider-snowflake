package sdk

import (
	"context"
	"errors"
)

type ExternalTables interface {
	// Create creates an external table with computed partitions.
	Create(ctx context.Context, id AccountObjectIdentifier) error
	// TODO: Create2 creates an external table where partitions are added and removed manually.
	Create2(ctx context.Context, id AccountObjectIdentifier) error
	// TODO: Create3 creates a delta lake external table.
	Create3(ctx context.Context, id AccountObjectIdentifier) error
	// Alter modifies an existing external table.
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterExternalTableOptions) error
	// Alter modifies an existing external table's partitions.
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
	Columns             bool
	CloudProviderParams *CloudProviderParams
	PartitionBy         []string `ddl:"list,parentheses" sql:"PARTITION BY"`
	Location            Location
	RefreshOnCreate     *bool   `ddl:"keword" sql:"REFRESH_ON_CREATE = TRUE"` // TODO: ?
	AutoRefresh         *bool   `ddl:"keword" sql:"AUTO_REFRESH = TRUE"`      // TODO: ?
	Pattern             *string `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          ExternalTableFileFormat
	AwsSnsTopic         *string `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	CopyGrants          *bool   `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy     *RowAccessPolicy
	Tag                 []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
	Comment             *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type CloudProviderParams struct {
	// One of
	GoogleCloudStorage *GoogleCloudStorageParams
	MicrosoftAzure     *MicrosoftAzure
}

func (cpp *CloudProviderParams) validate() error {
	if anyValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) && exactlyOneValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) {
		return errors.New("Only one cloud provider can be specified at once")
	}
	return nil
}

type GoogleCloudStorageParams struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type MicrosoftAzure struct {
	Integration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type Location struct{}

type ExternalTableFileFormat struct {
	// One of
	Name *string
	Type *string // TODO
}

type RowAccessPolicy struct{}

func (v *externalTables) Create(ctx context.Context, id AccountObjectIdentifier) error {
	return nil
}

// TODO: ChangeName
func (v *externalTables) Create2(ctx context.Context, id AccountObjectIdentifier) error {
	return nil
}

// TODO: ChangeName
func (v *externalTables) Create3(ctx context.Context, id AccountObjectIdentifier) error {
	return nil
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
