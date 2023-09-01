package sdk

//go:generate go run ./dto-builder-generator/main.go

type RowAccessPolicyRequest struct {
	Name SchemaObjectIdentifier // required
	On   []string
}

// TODO Check subtypes
// TODO Option types should be unexported

type CreateExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []ExternalTableColumn
	cloudProviderParams *CloudProviderParams
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	fileFormat          []ExternalTableFileFormat // required
	awsSnsTopic         *string
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []TagAssociation
}

type CreateWithManualPartitioningExternalTableRequest struct {
	orReplace                  *bool
	ifNotExists                *bool
	name                       AccountObjectIdentifier // required
	columns                    []ExternalTableColumn
	cloudProviderParams        *CloudProviderParams
	partitionBy                []string
	location                   string // required
	userSpecifiedPartitionType *bool
	fileFormat                 []ExternalTableFileFormat // required
	copyGrants                 *bool
	comment                    *string
	rowAccessPolicy            *RowAccessPolicy
	tag                        []TagAssociation
}

type CreateDeltaLakeExternalTableRequest struct {
	orReplace                  *bool
	ifNotExists                *bool
	name                       AccountObjectIdentifier // required
	columns                    []ExternalTableColumn
	cloudProviderParams        *CloudProviderParams
	partitionBy                []string
	location                   string // required
	userSpecifiedPartitionType *bool
	fileFormat                 []ExternalTableFileFormat // required
	deltaTableFormat           *bool
	copyGrants                 *bool
	comment                    *string
	rowAccessPolicy            *RowAccessPolicy
	tag                        []TagAssociation
}

type CreateExternalTableUsingTemplateRequest struct {
	orReplace           *bool
	name                AccountObjectIdentifier // required
	copyGrants          *bool
	query               string
	cloudProviderParams *CloudProviderParams
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	fileFormat          []ExternalTableFileFormat // required
	awsSnsTopic         *string
	comment             *string
	rowAccessPolicy     *RowAccessPolicy
	tag                 []TagAssociation
}

type AlterExternalTableRequest struct {
	ifExists    *bool
	name        AccountObjectIdentifier // required
	refresh     *RefreshExternalTable
	addFiles    []ExternalTableFile
	removeFiles []ExternalTableFile
	set         *ExternalTableSet
	unset       *ExternalTableUnset
}
type AlterExternalTablePartitionRequest struct {
	ifExists      *bool
	name          AccountObjectIdentifier // required
	addPartitions []Partition
	dropPartition *bool
	location      string
}

type DropExternalTableRequest struct {
	ifExists   *bool
	name       AccountObjectIdentifier // required
	dropOption *ExternalTableDropOption
}

type ShowExternalTableRequest struct {
	terse      *bool
	like       *Like
	in         *In
	startsWith *string
	limitFrom  *LimitFrom
}

type ShowExternalTableByIDRequest struct {
	id AccountObjectIdentifier // required
}

type DescribeExternalTableColumnsRequest struct {
	id AccountObjectIdentifier // required
}

type DescribeExternalTableStageRequest struct {
	id AccountObjectIdentifier // required
}
