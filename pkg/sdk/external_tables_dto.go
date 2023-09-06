package sdk

//go:generate go run ./dto-builder-generator/main.go

type CreateExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	fileFormat          *ExternalTableFileFormatRequest // required
	awsSnsTopic         *string
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

type ExternalTableColumnRequest struct {
	Name             string   // required
	Type             DataType // required
	AsExpression     string   // required
	InlineConstraint *ColumnInlineConstraintRequest
}

func (v ExternalTableColumnRequest) toOpts() ExternalTableColumn {
	var inlineConstraint *ColumnInlineConstraint
	if v.InlineConstraint != nil {
		inlineConstraint = v.InlineConstraint.toOpts()
	}

	return ExternalTableColumn{
		Name:             v.Name,
		Type:             v.Type,
		AsExpression:     v.AsExpression,
		InlineConstraint: inlineConstraint,
	}
}

func (v *ColumnInlineConstraintRequest) toOpts() *ColumnInlineConstraint {
	return &ColumnInlineConstraint{
		NotNull:            v.NotNull,
		Name:               &v.Name,
		Type:               &v.Type,
		ForeignKey:         v.ForeignKey,
		Enforced:           v.Enforced,
		NotEnforced:        v.NotEnforced,
		Deferrable:         v.Deferrable,
		NotDeferrable:      v.NotDeferrable,
		InitiallyDeferred:  v.InitiallyDeferred,
		InitiallyImmediate: v.InitiallyImmediate,
		Enable:             v.Enable,
		Disable:            v.Disable,
		Validate:           v.Validate,
		NoValidate:         v.NoValidate,
		Rely:               v.Rely,
		NoRely:             v.NoRely,
	}
}

type ColumnInlineConstraintRequest struct {
	NotNull    *bool
	Name       string               // required
	Type       ColumnConstraintType // required
	ForeignKey *InlineForeignKey

	// optional
	Enforced           *bool
	NotEnforced        *bool
	Deferrable         *bool
	NotDeferrable      *bool
	InitiallyDeferred  *bool
	InitiallyImmediate *bool
	Enable             *bool
	Disable            *bool
	Validate           *bool
	NoValidate         *bool
	Rely               *bool
	NoRely             *bool
}

type InlineForeignKeyRequest struct {
	TableName  string // required
	ColumnName []string
	Match      *MatchType
	On         *ForeignKeyOnActionRequest
}

type ForeignKeyOnActionRequest struct {
	OnUpdate *bool
	OnDelete *bool
}

type CloudProviderParamsRequest struct {
	GoogleCloudStorageIntegration *string
	MicrosoftAzureIntegration     *string
}

func (v *CloudProviderParamsRequest) toOpts() *CloudProviderParams {
	return &CloudProviderParams{
		GoogleCloudStorageIntegration: v.GoogleCloudStorageIntegration,
		MicrosoftAzureIntegration:     v.MicrosoftAzureIntegration,
	}
}

type ExternalTableFileFormatRequest struct {
	Name    *string
	Type    *ExternalTableFileFormatType
	Options *ExternalTableFileFormatTypeOptionsRequest
}

func (v *ExternalTableFileFormatTypeOptionsRequest) toOpts() *ExternalTableFileFormatTypeOptions {
	var csvNullIf []NullString
	if v.CSVNullIf != nil {
		for _, n := range *v.CSVNullIf {
			csvNullIf = append(csvNullIf, n.toOpts())
		}
	}

	var orcNullIf []NullString
	if v.ORCNullIf != nil {
		for _, n := range *v.ORCNullIf {
			orcNullIf = append(orcNullIf, n.toOpts())
		}
	}

	return &ExternalTableFileFormatTypeOptions{
		CSVCompression:                  v.CSVCompression,
		CSVRecordDelimiter:              v.CSVRecordDelimiter,
		CSVFieldDelimiter:               v.CSVFieldDelimiter,
		CSVSkipHeader:                   v.CSVSkipHeader,
		CSVSkipBlankLines:               v.CSVSkipBlankLines,
		CSVEscapeUnenclosedField:        v.CSVEscapeUnenclosedField,
		CSVTrimSpace:                    v.CSVTrimSpace,
		CSVFieldOptionallyEnclosedBy:    v.CSVFieldOptionallyEnclosedBy,
		CSVNullIf:                       &csvNullIf,
		CSVEmptyFieldAsNull:             v.CSVEmptyFieldAsNull,
		CSVEncoding:                     v.CSVEncoding,
		JSONCompression:                 v.JSONCompression,
		JSONAllowDuplicate:              v.JSONAllowDuplicate,
		JSONStripOuterArray:             v.JSONStripOuterArray,
		JSONStripNullValues:             v.JSONStripNullValues,
		JSONReplaceInvalidCharacters:    v.JSONReplaceInvalidCharacters,
		AvroCompression:                 v.AvroCompression,
		AvroReplaceInvalidCharacters:    v.AvroReplaceInvalidCharacters,
		ORCTrimSpace:                    v.ORCTrimSpace,
		ORCReplaceInvalidCharacters:     v.ORCReplaceInvalidCharacters,
		ORCNullIf:                       &orcNullIf,
		ParquetCompression:              v.ParquetCompression,
		ParquetBinaryAsText:             v.ParquetBinaryAsText,
		ParquetReplaceInvalidCharacters: v.ParquetReplaceInvalidCharacters,
	}
}

func (v ExternalTableFileFormatRequest) toOpts() ExternalTableFileFormat {
	var options *ExternalTableFileFormatTypeOptions
	if v.Options != nil {
		options = v.Options.toOpts()
	}

	return ExternalTableFileFormat{
		Name:    v.Name,
		Type:    v.Type,
		Options: options,
	}
}

type ExternalTableFileFormatTypeOptionsRequest struct {
	// CSV type options
	CSVCompression               *ExternalTableCsvCompression
	CSVRecordDelimiter           *string
	CSVFieldDelimiter            *string
	CSVSkipHeader                *int
	CSVSkipBlankLines            *bool
	CSVEscapeUnenclosedField     *string
	CSVTrimSpace                 *bool
	CSVFieldOptionallyEnclosedBy *string
	CSVNullIf                    *[]NullStringRequest
	CSVEmptyFieldAsNull          *bool
	CSVEncoding                  *CSVEncoding

	// JSON type options
	JSONCompression              *ExternalTableJsonCompression
	JSONAllowDuplicate           *bool
	JSONStripOuterArray          *bool
	JSONStripNullValues          *bool
	JSONReplaceInvalidCharacters *bool

	// AVRO type options
	AvroCompression              *ExternalTableAvroCompression
	AvroReplaceInvalidCharacters *bool

	// ORC type options
	ORCTrimSpace                *bool
	ORCReplaceInvalidCharacters *bool
	ORCNullIf                   *[]NullStringRequest

	// PARQUET type options
	ParquetCompression              *ExternalTableParquetCompression
	ParquetBinaryAsText             *bool
	ParquetReplaceInvalidCharacters *bool
}

type NullStringRequest struct {
	str string
}

func (v NullStringRequest) toOpts() NullString {
	return NullString{
		S: v.str,
	}
}

type RowAccessPolicyRequest struct {
	Name SchemaObjectIdentifier // required
	On   []string               // required
}

func (v *RowAccessPolicyRequest) toOpts() *RowAccessPolicy {
	return nil
}

type TagAssociationRequest struct {
	Name  ObjectIdentifier // required
	Value string           // required
}

func (v TagAssociationRequest) toOpts() TagAssociation {
	return TagAssociation{
		Name:  v.Name,
		Value: v.Value,
	}
}

func (v *CreateExternalTableRequest) toOpts() *CreateExternalTableOptions {
	columns := make([]ExternalTableColumn, len(v.columns))
	if v.columns != nil {
		for i, c := range v.columns {
			columns[i] = c.toOpts()
		}
	}

	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var rowAccessPolicy *RowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		Pattern:             v.pattern,
		FileFormat:          v.fileFormat.toOpts(),
		AwsSnsTopic:         v.awsSnsTopic,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateWithManualPartitioningExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string                          // required
	fileFormat          *ExternalTableFileFormatRequest // required
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateWithManualPartitioningExternalTableRequest) toOpts() *CreateWithManualPartitioningExternalTableOptions {
	columns := make([]ExternalTableColumn, len(v.columns))
	if v.columns != nil {
		for i, c := range v.columns {
			columns[i] = c.toOpts()
		}
	}

	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = v.fileFormat.toOpts()
	}

	var rowAccessPolicy *RowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateWithManualPartitioningExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		FileFormat:          fileFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateDeltaLakeExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                AccountObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string                          // required
	fileFormat          *ExternalTableFileFormatRequest // required
	deltaTableFormat    *bool
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateDeltaLakeExternalTableRequest) toOpts() *CreateDeltaLakeExternalTableOptions {
	columns := make([]ExternalTableColumn, len(v.columns))
	if v.columns != nil {
		for i, c := range v.columns {
			columns[i] = c.toOpts()
		}
	}

	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = v.fileFormat.toOpts()
	}

	var rowAccessPolicy *RowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateDeltaLakeExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		FileFormat:          fileFormat,
		DeltaTableFormat:    v.deltaTableFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateExternalTableUsingTemplateRequest struct {
	orReplace           *bool
	name                AccountObjectIdentifier // required
	copyGrants          *bool
	query               string
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	fileFormat          *ExternalTableFileFormatRequest // required
	awsSnsTopic         *string
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateExternalTableUsingTemplateRequest) toOpts() *CreateExternalTableUsingTemplateOptions {
	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = v.fileFormat.toOpts()
	}

	var rowAccessPolicy *RowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateExternalTableUsingTemplateOptions{
		OrReplace:           v.orReplace,
		name:                v.name,
		CopyGrants:          v.copyGrants,
		Query:               v.query,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		Pattern:             v.pattern,
		FileFormat:          fileFormat,
		AwsSnsTopic:         v.awsSnsTopic,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type AlterExternalTableRequest struct {
	ifExists    *bool
	name        AccountObjectIdentifier // required
	refresh     *RefreshExternalTableRequest
	addFiles    []*ExternalTableFileRequest
	removeFiles []*ExternalTableFileRequest
	autoRefresh *bool
	setTag      []*TagAssociationRequest
	unsetTag    []ObjectIdentifier
}

type RefreshExternalTableRequest struct {
	Path string // required
}

type ExternalTableFileRequest struct {
	Name string // required
}

func (v *AlterExternalTableRequest) toOpts() *AlterExternalTableOptions {
	var refresh *RefreshExternalTable
	if v.refresh != nil {
		refresh = &RefreshExternalTable{
			Path: v.refresh.Path,
		}
	}

	addFiles := make([]ExternalTableFile, len(v.addFiles))
	if v.addFiles != nil {
		for i, f := range v.addFiles {
			addFiles[i] = ExternalTableFile{
				Name: f.Name,
			}
		}
	}

	removeFiles := make([]ExternalTableFile, len(v.removeFiles))
	if v.removeFiles != nil {
		for i, f := range v.removeFiles {
			removeFiles[i] = ExternalTableFile{
				Name: f.Name,
			}
		}
	}

	setTag := make([]TagAssociation, len(v.setTag))
	if v.setTag != nil {
		for i, t := range v.setTag {
			setTag[i] = t.toOpts()
		}
	}

	return &AlterExternalTableOptions{
		IfExists:    v.ifExists,
		name:        v.name,
		Refresh:     refresh,
		AddFiles:    addFiles,
		RemoveFiles: removeFiles,
		AutoRefresh: v.autoRefresh,
		SetTag:      setTag,
		UnsetTag:    v.unsetTag,
	}
}

type AlterExternalTablePartitionRequest struct {
	ifExists      *bool
	name          AccountObjectIdentifier // required
	addPartitions []*PartitionRequest
	dropPartition *bool
	location      string
}

type PartitionRequest struct {
	ColumnName string // required
	Value      string // required
}

func (v *AlterExternalTablePartitionRequest) toOpts() *AlterExternalTablePartitionOptions {
	addPartitions := make([]Partition, len(v.addPartitions))
	if v.addPartitions != nil {
		for i, p := range v.addPartitions {
			addPartitions[i] = Partition{
				ColumnName: p.ColumnName,
				Value:      p.Value,
			}
		}
	}

	return &AlterExternalTablePartitionOptions{
		IfExists:      v.ifExists,
		name:          v.name,
		AddPartitions: addPartitions,
		DropPartition: v.dropPartition,
		Location:      v.location,
	}
}

type DropExternalTableRequest struct {
	ifExists   *bool
	name       AccountObjectIdentifier // required
	dropOption *ExternalTableDropOptionRequest
}

type ExternalTableDropOptionRequest struct {
	Restrict *bool
	Cascade  *bool
}

func (v *ExternalTableDropOptionRequest) toOpts() *ExternalTableDropOption {
	return &ExternalTableDropOption{
		Restrict: v.Restrict,
		Cascade:  v.Cascade,
	}
}

func (v *DropExternalTableRequest) toOpts() *DropExternalTableOptions {
	var dropOption *ExternalTableDropOption
	if v.dropOption != nil {
		dropOption = v.dropOption.toOpts()
	}

	return &DropExternalTableOptions{
		IfExists:   v.ifExists,
		name:       v.name,
		DropOption: dropOption,
	}
}

type ShowExternalTableRequest struct {
	terse      *bool
	like       *string
	in         *ShowExternalTableInRequest
	startsWith *string
	limitFrom  *LimitFromRequest
}

type ShowExternalTableInRequest struct {
	Account  *bool
	Database AccountObjectIdentifier
	Schema   DatabaseObjectIdentifier
}

func (v *ShowExternalTableInRequest) toOpts() *In {
	return &In{
		Account:  v.Account,
		Database: v.Database,
		Schema:   v.Schema,
	}
}

type LimitFromRequest struct {
	Rows *int
	From *string
}

func (v *LimitFromRequest) toOpts() *LimitFrom {
	return &LimitFrom{
		Rows: v.Rows,
		From: v.From,
	}
}

func (v *ShowExternalTableRequest) toOpts() *ShowExternalTableOptions {
	var like *Like
	if v.like != nil {
		like = &Like{
			Pattern: v.like,
		}
	}

	var in *In
	if v.in != nil {
		in = v.in.toOpts()
	}

	var limitFrom *LimitFrom
	if v.limitFrom != nil {
		limitFrom = v.limitFrom.toOpts()
	}

	return &ShowExternalTableOptions{
		Terse:      v.terse,
		Like:       like,
		In:         in,
		StartsWith: v.startsWith,
		LimitFrom:  limitFrom,
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
