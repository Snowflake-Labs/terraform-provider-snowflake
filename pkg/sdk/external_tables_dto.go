package sdk

//go:generate go run ./dto-builder-generator/main.go

// TODO Check subtypes (e.g. columns in create, etc.)
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
	fileFormat          ExternalTableFileFormat // required
	awsSnsTopic         *string
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicy
	tag                 []TagAssociation
}

func (v *CreateExternalTableRequest) toOpts() *CreateExternalTableOptions {
	return &CreateExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             v.columns,
		CloudProviderParams: v.cloudProviderParams,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		Pattern:             v.pattern,
		FileFormat:          v.fileFormat,
		AwsSnsTopic:         v.awsSnsTopic,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     v.rowAccessPolicy,
		Tag:                 v.tag,
	}
}

type CreateWithManualPartitioningExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []ExternalTableColumn
	cloudProviderParams *CloudProviderParams
	partitionBy         []string
	location            string                  // required
	fileFormat          ExternalTableFileFormat // required
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicy
	tag                 []TagAssociation
}

func (v *CreateWithManualPartitioningExternalTableRequest) toOpts() *CreateWithManualPartitioningExternalTableOptions {
	return &CreateWithManualPartitioningExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             v.columns,
		CloudProviderParams: v.cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		FileFormat:          v.fileFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     v.rowAccessPolicy,
		Tag:                 v.tag,
	}
}

type CreateDeltaLakeExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []ExternalTableColumn
	cloudProviderParams *CloudProviderParams
	partitionBy         []string
	location            string                  // required
	fileFormat          ExternalTableFileFormat // required
	deltaTableFormat    *bool
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicy
	tag                 []TagAssociation
}

func (v *CreateDeltaLakeExternalTableRequest) toOpts() *CreateDeltaLakeExternalTableOptions {
	return &CreateDeltaLakeExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             v.columns,
		CloudProviderParams: v.cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		FileFormat:          v.fileFormat,
		DeltaTableFormat:    v.deltaTableFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     v.rowAccessPolicy,
		Tag:                 v.tag,
	}
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
	fileFormat          ExternalTableFileFormat // required
	awsSnsTopic         *string
	comment             *string
	rowAccessPolicy     *RowAccessPolicy
	tag                 []TagAssociation
}

func (v *CreateExternalTableUsingTemplateRequest) toOpts() *CreateExternalTableUsingTemplateOptions {
	return &CreateExternalTableUsingTemplateOptions{
		OrReplace:           v.orReplace,
		name:                v.name,
		CopyGrants:          v.copyGrants,
		Query:               v.query,
		CloudProviderParams: v.cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		Pattern:             v.pattern,
		FileFormat:          v.fileFormat,
		AwsSnsTopic:         v.awsSnsTopic,
		Comment:             v.comment,
		RowAccessPolicy:     v.rowAccessPolicy,
		Tag:                 v.tag,
	}
}

type AlterExternalTableRequest struct {
	ifExists    *bool
	name        AccountObjectIdentifier // required
	refresh     *RefreshExternalTable
	addFiles    []ExternalTableFile
	removeFiles []ExternalTableFile
	autoRefresh *bool              `ddl:"parameter" sql:"AUTO_REFRESH"`
	setTag      []TagAssociation   `ddl:"keyword" sql:"TAG"`
	unsetTag    []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *AlterExternalTableRequest) toOpts() *AlterExternalTableOptions {
	return &AlterExternalTableOptions{
		IfExists:    v.ifExists,
		name:        v.name,
		Refresh:     v.refresh,
		AddFiles:    v.addFiles,
		RemoveFiles: v.removeFiles,
		AutoRefresh: v.autoRefresh,
		SetTag:      v.setTag,
		UnsetTag:    v.unsetTag,
	}
}

type AlterExternalTablePartitionRequest struct {
	ifExists      *bool
	name          AccountObjectIdentifier // required
	addPartitions []Partition
	dropPartition *bool
	location      string
}

func (v *AlterExternalTablePartitionRequest) toOpts() *AlterExternalTablePartitionOptions {
	return &AlterExternalTablePartitionOptions{
		IfExists:      v.ifExists,
		name:          v.name,
		AddPartitions: v.addPartitions,
		DropPartition: v.dropPartition,
		Location:      v.location,
	}
}

type DropExternalTableRequest struct {
	ifExists   *bool
	name       AccountObjectIdentifier // required
	dropOption *ExternalTableDropOption
}

func (v *DropExternalTableRequest) toOpts() *DropExternalTableOptions {
	return &DropExternalTableOptions{
		IfExists:   v.ifExists,
		name:       v.name,
		DropOption: v.dropOption,
	}
}

type ShowExternalTableRequest struct {
	terse      *bool
	like       *Like
	in         *In
	startsWith *string
	limitFrom  *LimitFrom
}

func (v *ShowExternalTableRequest) toOpts() *ShowExternalTableOptions {
	return &ShowExternalTableOptions{
		Terse:      v.terse,
		Like:       v.like,
		In:         v.in,
		StartsWith: v.startsWith,
		LimitFrom:  v.limitFrom,
	}
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
