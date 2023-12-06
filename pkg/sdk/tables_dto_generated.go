package sdk

import (
	"time"
)

func NewCreateTableAsSelectRequest(
	name SchemaObjectIdentifier,
	columns []TableAsSelectColumnRequest,
	query string,
) *CreateTableAsSelectRequest {
	s := CreateTableAsSelectRequest{}
	s.name = name
	s.columns = columns
	s.query = query
	return &s
}

func (s *CreateTableAsSelectRequest) WithOrReplace(orReplace *bool) *CreateTableAsSelectRequest {
	s.orReplace = orReplace
	return s
}

func NewTableAsSelectColumnRequest(
	name string,
) *TableAsSelectColumnRequest {
	s := TableAsSelectColumnRequest{}
	s.name = name
	return &s
}

func (s *TableAsSelectColumnRequest) WithOrReplace(orReplace *bool) *TableAsSelectColumnRequest {
	s.orReplace = orReplace
	return s
}

func (s *TableAsSelectColumnRequest) WithType_(type_ *DataType) *TableAsSelectColumnRequest {
	s.type_ = type_
	return s
}

func (s *TableAsSelectColumnRequest) WithMaskingPolicyName(maskingPolicyName *SchemaObjectIdentifier) *TableAsSelectColumnRequest {
	s.maskingPolicyName = maskingPolicyName
	return s
}

func (s *TableAsSelectColumnRequest) WithClusterBy(clusterBy []string) *TableAsSelectColumnRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *TableAsSelectColumnRequest) WithCopyGrants(copyGrants *bool) *TableAsSelectColumnRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableUsingTemplateRequest(
	name SchemaObjectIdentifier,
	Query string,
) *CreateTableUsingTemplateRequest {
	s := CreateTableUsingTemplateRequest{}
	s.name = name
	s.Query = Query
	return &s
}

func (s *CreateTableUsingTemplateRequest) WithOrReplace(orReplace *bool) *CreateTableUsingTemplateRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableUsingTemplateRequest) WithCopyGrants(copyGrants *bool) *CreateTableUsingTemplateRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableLikeRequest(
	name SchemaObjectIdentifier,
	sourceTable SchemaObjectIdentifier,
) *CreateTableLikeRequest {
	s := CreateTableLikeRequest{}
	s.name = name
	s.sourceTable = sourceTable
	return &s
}

func (s *CreateTableLikeRequest) WithOrReplace(orReplace *bool) *CreateTableLikeRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableLikeRequest) WithClusterBy(clusterBy []string) *CreateTableLikeRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *CreateTableLikeRequest) WithCopyGrants(copyGrants *bool) *CreateTableLikeRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableCloneRequest(
	name SchemaObjectIdentifier,
	sourceTable SchemaObjectIdentifier,
) *CreateTableCloneRequest {
	s := CreateTableCloneRequest{}
	s.name = name
	s.sourceTable = sourceTable
	return &s
}

func (s *CreateTableCloneRequest) WithOrReplace(orReplace *bool) *CreateTableCloneRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableCloneRequest) WithCopyGrants(copyGrants *bool) *CreateTableCloneRequest {
	s.copyGrants = copyGrants
	return s
}

func (s *CreateTableCloneRequest) WithClonePoint(ClonePoint *ClonePointRequest) *CreateTableCloneRequest {
	s.ClonePoint = ClonePoint
	return s
}

func NewClonePointRequest() *ClonePointRequest {
	return &ClonePointRequest{}
}

func (s *ClonePointRequest) WithMoment(Moment CloneMoment) *ClonePointRequest {
	s.Moment = Moment
	return s
}

func (s *ClonePointRequest) WithAt(At TimeTravelRequest) *ClonePointRequest {
	s.At = At
	return s
}

func NewTimeTravelRequest() *TimeTravelRequest {
	return &TimeTravelRequest{}
}

func (s *TimeTravelRequest) WithTimestamp(Timestamp *time.Time) *TimeTravelRequest {
	s.Timestamp = Timestamp
	return s
}

func (s *TimeTravelRequest) WithOffset(Offset *int) *TimeTravelRequest {
	s.Offset = Offset
	return s
}

func (s *TimeTravelRequest) WithStatement(Statement *string) *TimeTravelRequest {
	s.Statement = Statement
	return s
}

func NewCreateTableRequest(
	name SchemaObjectIdentifier,
	columns []TableColumnRequest,
) *CreateTableRequest {
	s := CreateTableRequest{}
	s.name = name
	s.columns = columns
	return &s
}

func (s *CreateTableRequest) WithOrReplace(orReplace *bool) *CreateTableRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableRequest) WithIfNotExists(ifNotExists *bool) *CreateTableRequest {
	s.ifNotExists = ifNotExists
	return s
}

func (s *CreateTableRequest) WithScope(scope *TableScope) *CreateTableRequest {
	s.scope = scope
	return s
}

func (s *CreateTableRequest) WithKind(kind *TableKind) *CreateTableRequest {
	s.kind = kind
	return s
}

func (s *CreateTableRequest) WithOutOfLineConstraint(OutOfLineConstraint OutOfLineConstraintRequest) *CreateTableRequest {
	s.OutOfLineConstraints = append(s.OutOfLineConstraints, OutOfLineConstraint)
	return s
}

func (s *CreateTableRequest) WithClusterBy(clusterBy []string) *CreateTableRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *CreateTableRequest) WithEnableSchemaEvolution(enableSchemaEvolution *bool) *CreateTableRequest {
	s.enableSchemaEvolution = enableSchemaEvolution
	return s
}

func (s *CreateTableRequest) WithStageFileFormat(stageFileFormat StageFileFormatRequest) *CreateTableRequest {
	s.stageFileFormat = &stageFileFormat
	return s
}

func (s *CreateTableRequest) WithStageCopyOptions(stageCopyOptions StageCopyOptionsRequest) *CreateTableRequest {
	s.stageCopyOptions = &stageCopyOptions
	return s
}

func (s *CreateTableRequest) WithDataRetentionTimeInDays(DataRetentionTimeInDays *int) *CreateTableRequest {
	s.DataRetentionTimeInDays = DataRetentionTimeInDays
	return s
}

func (s *CreateTableRequest) WithMaxDataExtensionTimeInDays(MaxDataExtensionTimeInDays *int) *CreateTableRequest {
	s.MaxDataExtensionTimeInDays = MaxDataExtensionTimeInDays
	return s
}

func (s *CreateTableRequest) WithChangeTracking(ChangeTracking *bool) *CreateTableRequest {
	s.ChangeTracking = ChangeTracking
	return s
}

func (s *CreateTableRequest) WithDefaultDDLCollation(DefaultDDLCollation *string) *CreateTableRequest {
	s.DefaultDDLCollation = DefaultDDLCollation
	return s
}

func (s *CreateTableRequest) WithCopyGrants(CopyGrants *bool) *CreateTableRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateTableRequest) WithRowAccessPolicy(RowAccessPolicy *RowAccessPolicyRequest) *CreateTableRequest {
	s.RowAccessPolicy = RowAccessPolicy
	return s
}

func (s *CreateTableRequest) WithTags(Tags []TagAssociationRequest) *CreateTableRequest {
	s.Tags = Tags
	return s
}

func (s *CreateTableRequest) WithComment(Comment *string) *CreateTableRequest {
	s.Comment = Comment
	return s
}

func NewTableColumnRequest(
	name string,
	type_ DataType,
) *TableColumnRequest {
	s := TableColumnRequest{}
	s.name = name
	s.type_ = type_
	return &s
}

func (s *TableColumnRequest) WithCollate(collate *string) *TableColumnRequest {
	s.collate = collate
	return s
}

func (s *TableColumnRequest) WithComment(comment *string) *TableColumnRequest {
	s.comment = comment
	return s
}

func (s *TableColumnRequest) WithDefaultValue(defaultValue *ColumnDefaultValueRequest) *TableColumnRequest {
	s.defaultValue = defaultValue
	return s
}

func (s *TableColumnRequest) WithNotNull(notNull *bool) *TableColumnRequest {
	s.notNull = notNull
	return s
}

func (s *TableColumnRequest) WithMaskingPolicy(maskingPolicy *ColumnMaskingPolicyRequest) *TableColumnRequest {
	s.maskingPolicy = maskingPolicy
	return s
}

func (s *TableColumnRequest) WithTags(tags []TagAssociation) *TableColumnRequest {
	s.tags = tags
	return s
}

func (s *TableColumnRequest) WithInlineConstraint(inlineConstraint *ColumnInlineConstraintRequest) *TableColumnRequest {
	s.inlineConstraint = inlineConstraint
	return s
}

func NewColumnDefaultValueRequest() *ColumnDefaultValueRequest {
	return &ColumnDefaultValueRequest{}
}

func (s *ColumnDefaultValueRequest) WithExpression(expression *string) *ColumnDefaultValueRequest {
	s.expression = expression
	return s
}

func (s *ColumnDefaultValueRequest) WithIdentity(identity *ColumnIdentityRequest) *ColumnDefaultValueRequest {
	s.identity = identity
	return s
}

func NewColumnIdentityRequest(
	Start int,
	Increment int,
) *ColumnIdentityRequest {
	s := ColumnIdentityRequest{}
	s.Start = Start
	s.Increment = Increment
	return &s
}

func (s *ColumnIdentityRequest) WithOrder() *ColumnIdentityRequest {
	s.Order = Bool(true)
	return s
}

func (s *ColumnIdentityRequest) WithNoorder() *ColumnIdentityRequest {
	s.Noorder = Bool(true)
	return s
}

func NewColumnMaskingPolicyRequest(
	name SchemaObjectIdentifier,
) *ColumnMaskingPolicyRequest {
	s := ColumnMaskingPolicyRequest{}
	s.name = name
	return &s
}

func (s *ColumnMaskingPolicyRequest) WithWith(with *bool) *ColumnMaskingPolicyRequest {
	s.with = with
	return s
}

func (s *ColumnMaskingPolicyRequest) WithUsing(using []string) *ColumnMaskingPolicyRequest {
	s.using = using
	return s
}

func NewColumnInlineConstraintRequest(
	Name string,
	type_ ColumnConstraintType,
) *ColumnInlineConstraintRequest {
	s := ColumnInlineConstraintRequest{}
	s.Name = Name
	s.type_ = type_
	return &s
}

func (s *ColumnInlineConstraintRequest) WithNotNull(notNull *bool) *ColumnInlineConstraintRequest {
	s.notNull = notNull
	return s
}

func (s *ColumnInlineConstraintRequest) WithForeignKey(foreignKey *InlineForeignKeyRequest) *ColumnInlineConstraintRequest {
	s.foreignKey = foreignKey
	return s
}

func (s *ColumnInlineConstraintRequest) WithEnforced(enforced *bool) *ColumnInlineConstraintRequest {
	s.enforced = enforced
	return s
}

func (s *ColumnInlineConstraintRequest) WithNotEnforced(notEnforced *bool) *ColumnInlineConstraintRequest {
	s.notEnforced = notEnforced
	return s
}

func (s *ColumnInlineConstraintRequest) WithDeferrable(deferrable *bool) *ColumnInlineConstraintRequest {
	s.deferrable = deferrable
	return s
}

func (s *ColumnInlineConstraintRequest) WithNotDeferrable(notDeferrable *bool) *ColumnInlineConstraintRequest {
	s.notDeferrable = notDeferrable
	return s
}

func (s *ColumnInlineConstraintRequest) WithInitiallyDeferred(initiallyDeferred *bool) *ColumnInlineConstraintRequest {
	s.initiallyDeferred = initiallyDeferred
	return s
}

func (s *ColumnInlineConstraintRequest) WithInitiallyImmediate(initiallyImmediate *bool) *ColumnInlineConstraintRequest {
	s.initiallyImmediate = initiallyImmediate
	return s
}

func (s *ColumnInlineConstraintRequest) WithEnable(enable *bool) *ColumnInlineConstraintRequest {
	s.enable = enable
	return s
}

func (s *ColumnInlineConstraintRequest) WithDisable(disable *bool) *ColumnInlineConstraintRequest {
	s.disable = disable
	return s
}

func (s *ColumnInlineConstraintRequest) WithValidate(validate *bool) *ColumnInlineConstraintRequest {
	s.validate = validate
	return s
}

func (s *ColumnInlineConstraintRequest) WithNoValidate(noValidate *bool) *ColumnInlineConstraintRequest {
	s.noValidate = noValidate
	return s
}

func (s *ColumnInlineConstraintRequest) WithRely(rely *bool) *ColumnInlineConstraintRequest {
	s.rely = rely
	return s
}

func (s *ColumnInlineConstraintRequest) WithNoRely(noRely *bool) *ColumnInlineConstraintRequest {
	s.noRely = noRely
	return s
}

func NewOutOfLineConstraintRequest(
	Name string,
	Type ColumnConstraintType,
) *OutOfLineConstraintRequest {
	s := OutOfLineConstraintRequest{}
	s.Name = Name
	s.Type = Type
	return &s
}

func (s *OutOfLineConstraintRequest) WithColumns(Columns []string) *OutOfLineConstraintRequest {
	s.Columns = Columns
	return s
}

func (s *OutOfLineConstraintRequest) WithForeignKey(ForeignKey *OutOfLineForeignKeyRequest) *OutOfLineConstraintRequest {
	s.ForeignKey = ForeignKey
	return s
}

func (s *OutOfLineConstraintRequest) WithEnforced(Enforced *bool) *OutOfLineConstraintRequest {
	s.Enforced = Enforced
	return s
}

func (s *OutOfLineConstraintRequest) WithNotEnforced(NotEnforced *bool) *OutOfLineConstraintRequest {
	s.NotEnforced = NotEnforced
	return s
}

func (s *OutOfLineConstraintRequest) WithDeferrable(Deferrable *bool) *OutOfLineConstraintRequest {
	s.Deferrable = Deferrable
	return s
}

func (s *OutOfLineConstraintRequest) WithNotDeferrable(NotDeferrable *bool) *OutOfLineConstraintRequest {
	s.NotDeferrable = NotDeferrable
	return s
}

func (s *OutOfLineConstraintRequest) WithInitiallyDeferred(InitiallyDeferred *bool) *OutOfLineConstraintRequest {
	s.InitiallyDeferred = InitiallyDeferred
	return s
}

func (s *OutOfLineConstraintRequest) WithInitiallyImmediate(InitiallyImmediate *bool) *OutOfLineConstraintRequest {
	s.InitiallyImmediate = InitiallyImmediate
	return s
}

func (s *OutOfLineConstraintRequest) WithEnable(Enable *bool) *OutOfLineConstraintRequest {
	s.Enable = Enable
	return s
}

func (s *OutOfLineConstraintRequest) WithDisable(Disable *bool) *OutOfLineConstraintRequest {
	s.Disable = Disable
	return s
}

func (s *OutOfLineConstraintRequest) WithValidate(Validate *bool) *OutOfLineConstraintRequest {
	s.Validate = Validate
	return s
}

func (s *OutOfLineConstraintRequest) WithNoValidate(NoValidate *bool) *OutOfLineConstraintRequest {
	s.NoValidate = NoValidate
	return s
}

func (s *OutOfLineConstraintRequest) WithRely(Rely *bool) *OutOfLineConstraintRequest {
	s.Rely = Rely
	return s
}

func (s *OutOfLineConstraintRequest) WithNoRely(NoRely *bool) *OutOfLineConstraintRequest {
	s.NoRely = NoRely
	return s
}

func NewInlineForeignKeyRequest(
	TableName string,
) *InlineForeignKeyRequest {
	s := InlineForeignKeyRequest{}
	s.TableName = TableName
	return &s
}

func (s *InlineForeignKeyRequest) WithColumnName(ColumnName []string) *InlineForeignKeyRequest {
	s.ColumnName = ColumnName
	return s
}

func (s *InlineForeignKeyRequest) WithMatch(Match *MatchType) *InlineForeignKeyRequest {
	s.Match = Match
	return s
}

func (s *InlineForeignKeyRequest) WithOn(On *ForeignKeyOnAction) *InlineForeignKeyRequest {
	s.On = On
	return s
}

func NewOutOfLineForeignKeyRequest(
	TableName SchemaObjectIdentifier,
	ColumnNames []string,
) *OutOfLineForeignKeyRequest {
	s := OutOfLineForeignKeyRequest{}
	s.TableName = TableName
	s.ColumnNames = ColumnNames
	return &s
}

func (s *OutOfLineForeignKeyRequest) WithMatch(Match *MatchType) *OutOfLineForeignKeyRequest {
	s.Match = Match
	return s
}

func (s *OutOfLineForeignKeyRequest) WithOn(On *ForeignKeyOnAction) *OutOfLineForeignKeyRequest {
	s.On = On
	return s
}

func NewForeignKeyOnAction() *ForeignKeyOnAction {
	return &ForeignKeyOnAction{}
}

func (s *ForeignKeyOnAction) WithOnUpdate(OnUpdate *ForeignKeyAction) *ForeignKeyOnAction {
	s.OnUpdate = OnUpdate
	return s
}

func (s *ForeignKeyOnAction) WithOnDelete(OnDelete *ForeignKeyAction) *ForeignKeyOnAction {
	s.OnDelete = OnDelete
	return s
}

func NewAlterTableRequest(
	name SchemaObjectIdentifier,
) *AlterTableRequest {
	s := AlterTableRequest{}
	s.name = name
	return &s
}

func (s *AlterTableRequest) WithIfExists(IfExists *bool) *AlterTableRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterTableRequest) WithNewName(NewName *SchemaObjectIdentifier) *AlterTableRequest {
	s.NewName = NewName
	return s
}

func (s *AlterTableRequest) WithSwapWith(SwapWith *SchemaObjectIdentifier) *AlterTableRequest {
	s.SwapWith = SwapWith
	return s
}

func (s *AlterTableRequest) WithClusteringAction(ClusteringAction *TableClusteringActionRequest) *AlterTableRequest {
	s.ClusteringAction = ClusteringAction
	return s
}

func (s *AlterTableRequest) WithColumnAction(ColumnAction *TableColumnActionRequest) *AlterTableRequest {
	s.ColumnAction = ColumnAction
	return s
}

func (s *AlterTableRequest) WithConstraintAction(ConstraintAction *TableConstraintActionRequest) *AlterTableRequest {
	s.ConstraintAction = ConstraintAction
	return s
}

func (s *AlterTableRequest) WithExternalTableAction(ExternalTableAction *TableExternalTableActionRequest) *AlterTableRequest {
	s.ExternalTableAction = ExternalTableAction
	return s
}

func (s *AlterTableRequest) WithSearchOptimizationAction(SearchOptimizationAction *TableSearchOptimizationActionRequest) *AlterTableRequest {
	s.SearchOptimizationAction = SearchOptimizationAction
	return s
}

func (s *AlterTableRequest) WithSet(Set *TableSetRequest) *AlterTableRequest {
	s.Set = Set
	return s
}

func (s *AlterTableRequest) WithSetTags(SetTags []TagAssociationRequest) *AlterTableRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterTableRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterTableRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterTableRequest) WithUnset(Unset *TableUnsetRequest) *AlterTableRequest {
	s.Unset = Unset
	return s
}

func (s *AlterTableRequest) WithAddRowAccessPolicy(AddRowAccessPolicy *AddRowAccessPolicyRequest) *AlterTableRequest {
	s.AddRowAccessPolicy = AddRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropRowAccessPolicy(DropRowAccessPolicy *string) *AlterTableRequest {
	s.DropRowAccessPolicy = DropRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropAndAddRowAccessPolicy(DropAndAddRowAccessPolicy *DropAndAddRowAccessPolicyRequest) *AlterTableRequest {
	s.DropAndAddRowAccessPolicy = DropAndAddRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropAllAccessRowPolicies(DropAllAccessRowPolicies *bool) *AlterTableRequest {
	s.DropAllAccessRowPolicies = DropAllAccessRowPolicies
	return s
}

func NewDropTableRequest(
	Name SchemaObjectIdentifier,
) *DropTableRequest {
	s := DropTableRequest{}
	s.Name = Name
	return &s
}

func (s *DropTableRequest) WithIfExists(IfExists *bool) *DropTableRequest {
	s.IfExists = IfExists
	return s
}

func (s *DropTableRequest) WithCascade(Cascade *bool) *DropTableRequest {
	s.Cascade = Cascade
	return s
}

func (s *DropTableRequest) WithRestrict(Restrict *bool) *DropTableRequest {
	s.Restrict = Restrict
	return s
}

func NewDropAndAddRowAccessPolicyRequest(
	DroppedPolicyName string,
	AddedPolicy AddRowAccessPolicyRequest,
) *DropAndAddRowAccessPolicyRequest {
	s := DropAndAddRowAccessPolicyRequest{}
	s.DroppedPolicyName = DroppedPolicyName
	s.AddedPolicy = AddedPolicy
	return &s
}

func NewTableUnsetRequest() *TableUnsetRequest {
	return &TableUnsetRequest{}
}

func (s *TableUnsetRequest) WithDataRetentionTimeInDays(DataRetentionTimeInDays bool) *TableUnsetRequest {
	s.DataRetentionTimeInDays = DataRetentionTimeInDays
	return s
}

func (s *TableUnsetRequest) WithMaxDataExtensionTimeInDays(MaxDataExtensionTimeInDays bool) *TableUnsetRequest {
	s.MaxDataExtensionTimeInDays = MaxDataExtensionTimeInDays
	return s
}

func (s *TableUnsetRequest) WithChangeTracking(ChangeTracking bool) *TableUnsetRequest {
	s.ChangeTracking = ChangeTracking
	return s
}

func (s *TableUnsetRequest) WithDefaultDDLCollation(DefaultDDLCollation bool) *TableUnsetRequest {
	s.DefaultDDLCollation = DefaultDDLCollation
	return s
}

func (s *TableUnsetRequest) WithEnableSchemaEvolution(EnableSchemaEvolution bool) *TableUnsetRequest {
	s.EnableSchemaEvolution = EnableSchemaEvolution
	return s
}

func (s *TableUnsetRequest) WithComment(Comment bool) *TableUnsetRequest {
	s.Comment = Comment
	return s
}

func NewAddRowAccessPolicyRequest(
	PolicyName string,
	ColumnName []string,
) *AddRowAccessPolicyRequest {
	s := AddRowAccessPolicyRequest{}
	s.PolicyName = PolicyName
	s.ColumnName = ColumnName
	return &s
}

func NewTagAssociationRequest(
	Name ObjectIdentifier,
	Value string,
) *TagAssociationRequest {
	s := TagAssociationRequest{}
	s.Name = Name
	s.Value = Value
	return &s
}

func NewFileFormatTypeOptionsRequest() *FileFormatTypeOptionsRequest {
	return &FileFormatTypeOptionsRequest{}
}

func (s *FileFormatTypeOptionsRequest) WithCSVCompression(CSVCompression *CSVCompression) *FileFormatTypeOptionsRequest {
	s.CSVCompression = CSVCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVRecordDelimiter(CSVRecordDelimiter *string) *FileFormatTypeOptionsRequest {
	s.CSVRecordDelimiter = CSVRecordDelimiter
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVFieldDelimiter(CSVFieldDelimiter *string) *FileFormatTypeOptionsRequest {
	s.CSVFieldDelimiter = CSVFieldDelimiter
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVFileExtension(CSVFileExtension *string) *FileFormatTypeOptionsRequest {
	s.CSVFileExtension = CSVFileExtension
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVParseHeader(CSVParseHeader *bool) *FileFormatTypeOptionsRequest {
	s.CSVParseHeader = CSVParseHeader
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVSkipHeader(CSVSkipHeader *int) *FileFormatTypeOptionsRequest {
	s.CSVSkipHeader = CSVSkipHeader
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVSkipBlankLines(CSVSkipBlankLines *bool) *FileFormatTypeOptionsRequest {
	s.CSVSkipBlankLines = CSVSkipBlankLines
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVDateFormat(CSVDateFormat *string) *FileFormatTypeOptionsRequest {
	s.CSVDateFormat = CSVDateFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVTimeFormat(CSVTimeFormat *string) *FileFormatTypeOptionsRequest {
	s.CSVTimeFormat = CSVTimeFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVTimestampFormat(CSVTimestampFormat *string) *FileFormatTypeOptionsRequest {
	s.CSVTimestampFormat = CSVTimestampFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVBinaryFormat(CSVBinaryFormat *BinaryFormat) *FileFormatTypeOptionsRequest {
	s.CSVBinaryFormat = CSVBinaryFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVEscape(CSVEscape *string) *FileFormatTypeOptionsRequest {
	s.CSVEscape = CSVEscape
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVEscapeUnenclosedField(CSVEscapeUnenclosedField *string) *FileFormatTypeOptionsRequest {
	s.CSVEscapeUnenclosedField = CSVEscapeUnenclosedField
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVTrimSpace(CSVTrimSpace *bool) *FileFormatTypeOptionsRequest {
	s.CSVTrimSpace = CSVTrimSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVFieldOptionallyEnclosedBy(CSVFieldOptionallyEnclosedBy *string) *FileFormatTypeOptionsRequest {
	s.CSVFieldOptionallyEnclosedBy = CSVFieldOptionallyEnclosedBy
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVNullIf(CSVNullIf *[]NullString) *FileFormatTypeOptionsRequest {
	s.CSVNullIf = CSVNullIf
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVErrorOnColumnCountMismatch(CSVErrorOnColumnCountMismatch *bool) *FileFormatTypeOptionsRequest {
	s.CSVErrorOnColumnCountMismatch = CSVErrorOnColumnCountMismatch
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVReplaceInvalidCharacters(CSVReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.CSVReplaceInvalidCharacters = CSVReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVEmptyFieldAsNull(CSVEmptyFieldAsNull *bool) *FileFormatTypeOptionsRequest {
	s.CSVEmptyFieldAsNull = CSVEmptyFieldAsNull
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVSkipByteOrderMark(CSVSkipByteOrderMark *bool) *FileFormatTypeOptionsRequest {
	s.CSVSkipByteOrderMark = CSVSkipByteOrderMark
	return s
}

func (s *FileFormatTypeOptionsRequest) WithCSVEncoding(CSVEncoding *CSVEncoding) *FileFormatTypeOptionsRequest {
	s.CSVEncoding = CSVEncoding
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONCompression(JSONCompression *JSONCompression) *FileFormatTypeOptionsRequest {
	s.JSONCompression = JSONCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONDateFormat(JSONDateFormat *string) *FileFormatTypeOptionsRequest {
	s.JSONDateFormat = JSONDateFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONTimeFormat(JSONTimeFormat *string) *FileFormatTypeOptionsRequest {
	s.JSONTimeFormat = JSONTimeFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONTimestampFormat(JSONTimestampFormat *string) *FileFormatTypeOptionsRequest {
	s.JSONTimestampFormat = JSONTimestampFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONBinaryFormat(JSONBinaryFormat *BinaryFormat) *FileFormatTypeOptionsRequest {
	s.JSONBinaryFormat = JSONBinaryFormat
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONTrimSpace(JSONTrimSpace *bool) *FileFormatTypeOptionsRequest {
	s.JSONTrimSpace = JSONTrimSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONNullIf(JSONNullIf []NullString) *FileFormatTypeOptionsRequest {
	s.JSONNullIf = JSONNullIf
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONFileExtension(JSONFileExtension *string) *FileFormatTypeOptionsRequest {
	s.JSONFileExtension = JSONFileExtension
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONEnableOctal(JSONEnableOctal *bool) *FileFormatTypeOptionsRequest {
	s.JSONEnableOctal = JSONEnableOctal
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONAllowDuplicate(JSONAllowDuplicate *bool) *FileFormatTypeOptionsRequest {
	s.JSONAllowDuplicate = JSONAllowDuplicate
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONStripOuterArray(JSONStripOuterArray *bool) *FileFormatTypeOptionsRequest {
	s.JSONStripOuterArray = JSONStripOuterArray
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONStripNullValues(JSONStripNullValues *bool) *FileFormatTypeOptionsRequest {
	s.JSONStripNullValues = JSONStripNullValues
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONReplaceInvalidCharacters(JSONReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.JSONReplaceInvalidCharacters = JSONReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONIgnoreUTF8Errors(JSONIgnoreUTF8Errors *bool) *FileFormatTypeOptionsRequest {
	s.JSONIgnoreUTF8Errors = JSONIgnoreUTF8Errors
	return s
}

func (s *FileFormatTypeOptionsRequest) WithJSONSkipByteOrderMark(JSONSkipByteOrderMark *bool) *FileFormatTypeOptionsRequest {
	s.JSONSkipByteOrderMark = JSONSkipByteOrderMark
	return s
}

func (s *FileFormatTypeOptionsRequest) WithAvroCompression(AvroCompression *AvroCompression) *FileFormatTypeOptionsRequest {
	s.AvroCompression = AvroCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithAvroTrimSpace(AvroTrimSpace *bool) *FileFormatTypeOptionsRequest {
	s.AvroTrimSpace = AvroTrimSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithAvroReplaceInvalidCharacters(AvroReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.AvroReplaceInvalidCharacters = AvroReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithAvroNullIf(AvroNullIf *[]NullString) *FileFormatTypeOptionsRequest {
	s.AvroNullIf = AvroNullIf
	return s
}

func (s *FileFormatTypeOptionsRequest) WithORCTrimSpace(ORCTrimSpace *bool) *FileFormatTypeOptionsRequest {
	s.ORCTrimSpace = ORCTrimSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithORCReplaceInvalidCharacters(ORCReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.ORCReplaceInvalidCharacters = ORCReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithORCNullIf(ORCNullIf *[]NullString) *FileFormatTypeOptionsRequest {
	s.ORCNullIf = ORCNullIf
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetCompression(ParquetCompression *ParquetCompression) *FileFormatTypeOptionsRequest {
	s.ParquetCompression = ParquetCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetSnappyCompression(ParquetSnappyCompression *bool) *FileFormatTypeOptionsRequest {
	s.ParquetSnappyCompression = ParquetSnappyCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetBinaryAsText(ParquetBinaryAsText *bool) *FileFormatTypeOptionsRequest {
	s.ParquetBinaryAsText = ParquetBinaryAsText
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetTrimSpace(ParquetTrimSpace *bool) *FileFormatTypeOptionsRequest {
	s.ParquetTrimSpace = ParquetTrimSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetReplaceInvalidCharacters(ParquetReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.ParquetReplaceInvalidCharacters = ParquetReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithParquetNullIf(ParquetNullIf *[]NullString) *FileFormatTypeOptionsRequest {
	s.ParquetNullIf = ParquetNullIf
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLCompression(XMLCompression *XMLCompression) *FileFormatTypeOptionsRequest {
	s.XMLCompression = XMLCompression
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLIgnoreUTF8Errors(XMLIgnoreUTF8Errors *bool) *FileFormatTypeOptionsRequest {
	s.XMLIgnoreUTF8Errors = XMLIgnoreUTF8Errors
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLPreserveSpace(XMLPreserveSpace *bool) *FileFormatTypeOptionsRequest {
	s.XMLPreserveSpace = XMLPreserveSpace
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLStripOuterElement(XMLStripOuterElement *bool) *FileFormatTypeOptionsRequest {
	s.XMLStripOuterElement = XMLStripOuterElement
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLDisableSnowflakeData(XMLDisableSnowflakeData *bool) *FileFormatTypeOptionsRequest {
	s.XMLDisableSnowflakeData = XMLDisableSnowflakeData
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLDisableAutoConvert(XMLDisableAutoConvert *bool) *FileFormatTypeOptionsRequest {
	s.XMLDisableAutoConvert = XMLDisableAutoConvert
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLReplaceInvalidCharacters(XMLReplaceInvalidCharacters *bool) *FileFormatTypeOptionsRequest {
	s.XMLReplaceInvalidCharacters = XMLReplaceInvalidCharacters
	return s
}

func (s *FileFormatTypeOptionsRequest) WithXMLSkipByteOrderMark(XMLSkipByteOrderMark *bool) *FileFormatTypeOptionsRequest {
	s.XMLSkipByteOrderMark = XMLSkipByteOrderMark
	return s
}

func (s *FileFormatTypeOptionsRequest) WithComment(Comment *string) *FileFormatTypeOptionsRequest {
	s.Comment = Comment
	return s
}

func NewTableClusteringActionRequest() *TableClusteringActionRequest {
	return &TableClusteringActionRequest{}
}

func (s *TableClusteringActionRequest) WithClusterBy(ClusterBy []string) *TableClusteringActionRequest {
	s.ClusterBy = ClusterBy
	return s
}

func (s *TableClusteringActionRequest) WithRecluster(Recluster *TableReclusterActionRequest) *TableClusteringActionRequest {
	s.Recluster = Recluster
	return s
}

func (s *TableClusteringActionRequest) WithChangeReclusterState(ChangeReclusterState *ReclusterState) *TableClusteringActionRequest {
	s.ChangeReclusterState = ChangeReclusterState
	return s
}

func (s *TableClusteringActionRequest) WithDropClusteringKey(DropClusteringKey *bool) *TableClusteringActionRequest {
	s.DropClusteringKey = DropClusteringKey
	return s
}

func NewTableReclusterActionRequest() *TableReclusterActionRequest {
	return &TableReclusterActionRequest{}
}

func (s *TableReclusterActionRequest) WithMaxSize(MaxSize *int) *TableReclusterActionRequest {
	s.MaxSize = MaxSize
	return s
}

func (s *TableReclusterActionRequest) WithCondition(Condition *string) *TableReclusterActionRequest {
	s.Condition = Condition
	return s
}

func NewTableReclusterChangeStateRequest() *TableReclusterChangeStateRequest {
	return &TableReclusterChangeStateRequest{}
}

func (s *TableReclusterChangeStateRequest) WithState(State ReclusterState) *TableReclusterChangeStateRequest {
	s.State = State
	return s
}

func NewTableColumnActionRequest() *TableColumnActionRequest {
	return &TableColumnActionRequest{}
}

func (s *TableColumnActionRequest) WithAdd(Add *TableColumnAddActionRequest) *TableColumnActionRequest {
	s.Add = Add
	return s
}

func (s *TableColumnActionRequest) WithRename(Rename *TableColumnRenameActionRequest) *TableColumnActionRequest {
	s.Rename = Rename
	return s
}

func (s *TableColumnActionRequest) WithAlter(Alter []TableColumnAlterActionRequest) *TableColumnActionRequest {
	s.Alter = Alter
	return s
}

func (s *TableColumnActionRequest) WithSetMaskingPolicy(SetMaskingPolicy *TableColumnAlterSetMaskingPolicyActionRequest) *TableColumnActionRequest {
	s.SetMaskingPolicy = SetMaskingPolicy
	return s
}

func (s *TableColumnActionRequest) WithUnsetMaskingPolicy(UnsetMaskingPolicy *TableColumnAlterUnsetMaskingPolicyActionRequest) *TableColumnActionRequest {
	s.UnsetMaskingPolicy = UnsetMaskingPolicy
	return s
}

func (s *TableColumnActionRequest) WithSetTags(SetTags *TableColumnAlterSetTagsActionRequest) *TableColumnActionRequest {
	s.SetTags = SetTags
	return s
}

func (s *TableColumnActionRequest) WithUnsetTags(UnsetTags *TableColumnAlterUnsetTagsActionRequest) *TableColumnActionRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *TableColumnActionRequest) WithDropColumnsIfExists() *TableColumnActionRequest {
	s.DropColumnsIfExists = Bool(true)
	return s
}

func (s *TableColumnActionRequest) WithDropColumns(DropColumns []string) *TableColumnActionRequest {
	s.DropColumns = DropColumns
	return s
}

func NewTableColumnAddActionRequest(
	Name string,
	Type DataType,
) *TableColumnAddActionRequest {
	s := TableColumnAddActionRequest{}
	s.Name = Name
	s.Type = Type
	return &s
}

func (s *TableColumnAddActionRequest) WithIfNotExists() *TableColumnAddActionRequest {
	s.IfNotExists = Bool(true)
	return s
}

func (s *TableColumnAddActionRequest) WithDefaultValue(DefaultValue *ColumnDefaultValueRequest) *TableColumnAddActionRequest {
	s.DefaultValue = DefaultValue
	return s
}

func (s *TableColumnAddActionRequest) WithInlineConstraint(InlineConstraint *TableColumnAddInlineConstraintRequest) *TableColumnAddActionRequest {
	s.InlineConstraint = InlineConstraint
	return s
}

func (s *TableColumnAddActionRequest) WithMaskingPolicy(MaskingPolicy *ColumnMaskingPolicyRequest) *TableColumnAddActionRequest {
	s.MaskingPolicy = MaskingPolicy
	return s
}

func (s *TableColumnAddActionRequest) WithWith(With *bool) *TableColumnAddActionRequest {
	s.With = With
	return s
}

func (s *TableColumnAddActionRequest) WithTags(Tags []TagAssociation) *TableColumnAddActionRequest {
	s.Tags = Tags
	return s
}

func NewTableColumnAddInlineConstraintRequest() *TableColumnAddInlineConstraintRequest {
	return &TableColumnAddInlineConstraintRequest{}
}

func (s *TableColumnAddInlineConstraintRequest) WithNotNull(NotNull *bool) *TableColumnAddInlineConstraintRequest {
	s.NotNull = NotNull
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithName(Name string) *TableColumnAddInlineConstraintRequest {
	s.Name = Name
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithType(Type ColumnConstraintType) *TableColumnAddInlineConstraintRequest {
	s.Type = Type
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithForeignKey(ForeignKey *ColumnAddForeignKey) *TableColumnAddInlineConstraintRequest {
	s.ForeignKey = ForeignKey
	return s
}

func NewColumnAddForeignKeyRequest() *ColumnAddForeignKeyRequest {
	return &ColumnAddForeignKeyRequest{}
}

func (s *ColumnAddForeignKeyRequest) WithTableName(TableName string) *ColumnAddForeignKeyRequest {
	s.TableName = TableName
	return s
}

func (s *ColumnAddForeignKeyRequest) WithColumnName(ColumnName string) *ColumnAddForeignKeyRequest {
	s.ColumnName = ColumnName
	return s
}

func NewTableColumnRenameActionRequest(
	OldName string,
	NewName string,
) *TableColumnRenameActionRequest {
	s := TableColumnRenameActionRequest{}
	s.OldName = OldName
	s.NewName = NewName
	return &s
}

func NewTableColumnAlterActionRequest(
	Column bool,
	Name string,
) *TableColumnAlterActionRequest {
	s := TableColumnAlterActionRequest{}
	s.Column = Column
	s.Name = Name
	return &s
}

func (s *TableColumnAlterActionRequest) WithDropDefault(DropDefault *bool) *TableColumnAlterActionRequest {
	s.DropDefault = DropDefault
	return s
}

func (s *TableColumnAlterActionRequest) WithSetDefault(SetDefault *SequenceName) *TableColumnAlterActionRequest {
	s.SetDefault = SetDefault
	return s
}

func (s *TableColumnAlterActionRequest) WithNotNullConstraint(NotNullConstraint *TableColumnNotNullConstraintRequest) *TableColumnAlterActionRequest {
	s.NotNullConstraint = NotNullConstraint
	return s
}

func (s *TableColumnAlterActionRequest) WithType(Type *DataType) *TableColumnAlterActionRequest {
	s.Type = Type
	return s
}

func (s *TableColumnAlterActionRequest) WithComment(Comment *string) *TableColumnAlterActionRequest {
	s.Comment = Comment
	return s
}

func (s *TableColumnAlterActionRequest) WithUnsetComment(UnsetComment *bool) *TableColumnAlterActionRequest {
	s.UnsetComment = UnsetComment
	return s
}

func NewTableColumnAlterSetMaskingPolicyActionRequest(
	ColumnName string,
	MaskingPolicyName SchemaObjectIdentifier,
	Using []string,
) *TableColumnAlterSetMaskingPolicyActionRequest {
	s := TableColumnAlterSetMaskingPolicyActionRequest{}
	s.ColumnName = ColumnName
	s.MaskingPolicyName = MaskingPolicyName
	s.Using = Using
	return &s
}

func (s *TableColumnAlterSetMaskingPolicyActionRequest) WithForce(Force *bool) *TableColumnAlterSetMaskingPolicyActionRequest {
	s.Force = Force
	return s
}

func NewTableColumnAlterUnsetMaskingPolicyActionRequest(
	ColumnName string,
) *TableColumnAlterUnsetMaskingPolicyActionRequest {
	s := TableColumnAlterUnsetMaskingPolicyActionRequest{}
	s.ColumnName = ColumnName
	return &s
}

func NewTableColumnAlterSetTagsActionRequest(
	ColumnName string,
	Tags []TagAssociation,
) *TableColumnAlterSetTagsActionRequest {
	s := TableColumnAlterSetTagsActionRequest{}
	s.ColumnName = ColumnName
	s.Tags = Tags
	return &s
}

func NewTableColumnAlterUnsetTagsActionRequest(
	ColumnName string,
	Tags []ObjectIdentifier,
) *TableColumnAlterUnsetTagsActionRequest {
	s := TableColumnAlterUnsetTagsActionRequest{}
	s.ColumnName = ColumnName
	s.Tags = Tags
	return &s
}

func NewTableColumnNotNullConstraintRequest() *TableColumnNotNullConstraintRequest {
	return &TableColumnNotNullConstraintRequest{}
}

func (s *TableColumnNotNullConstraintRequest) WithSet(Set *bool) *TableColumnNotNullConstraintRequest {
	s.Set = Set
	return s
}

func (s *TableColumnNotNullConstraintRequest) WithDrop(Drop *bool) *TableColumnNotNullConstraintRequest {
	s.Drop = Drop
	return s
}

func NewTableConstraintActionRequest() *TableConstraintActionRequest {
	return &TableConstraintActionRequest{}
}

func (s *TableConstraintActionRequest) WithAdd(Add *OutOfLineConstraintRequest) *TableConstraintActionRequest {
	s.Add = Add
	return s
}

func (s *TableConstraintActionRequest) WithRename(Rename *TableConstraintRenameActionRequest) *TableConstraintActionRequest {
	s.Rename = Rename
	return s
}

func (s *TableConstraintActionRequest) WithAlter(Alter *TableConstraintAlterActionRequest) *TableConstraintActionRequest {
	s.Alter = Alter
	return s
}

func (s *TableConstraintActionRequest) WithDrop(Drop *TableConstraintDropActionRequest) *TableConstraintActionRequest {
	s.Drop = Drop
	return s
}

func NewTableConstraintRenameActionRequest() *TableConstraintRenameActionRequest {
	return &TableConstraintRenameActionRequest{}
}

func (s *TableConstraintRenameActionRequest) WithOldName(OldName string) *TableConstraintRenameActionRequest {
	s.OldName = OldName
	return s
}

func (s *TableConstraintRenameActionRequest) WithNewName(NewName string) *TableConstraintRenameActionRequest {
	s.NewName = NewName
	return s
}

func NewTableConstraintAlterActionRequest() *TableConstraintAlterActionRequest {
	return &TableConstraintAlterActionRequest{}
}

func (s *TableConstraintAlterActionRequest) WithConstraintName(ConstraintName *string) *TableConstraintAlterActionRequest {
	s.ConstraintName = ConstraintName
	return s
}

func (s *TableConstraintAlterActionRequest) WithPrimaryKey(PrimaryKey *bool) *TableConstraintAlterActionRequest {
	s.PrimaryKey = PrimaryKey
	return s
}

func (s *TableConstraintAlterActionRequest) WithUnique(Unique *bool) *TableConstraintAlterActionRequest {
	s.Unique = Unique
	return s
}

func (s *TableConstraintAlterActionRequest) WithForeignKey(ForeignKey *bool) *TableConstraintAlterActionRequest {
	s.ForeignKey = ForeignKey
	return s
}

func (s *TableConstraintAlterActionRequest) WithColumns(Columns []string) *TableConstraintAlterActionRequest {
	s.Columns = Columns
	return s
}

func (s *TableConstraintAlterActionRequest) WithEnforced(Enforced *bool) *TableConstraintAlterActionRequest {
	s.Enforced = Enforced
	return s
}

func (s *TableConstraintAlterActionRequest) WithNotEnforced(NotEnforced *bool) *TableConstraintAlterActionRequest {
	s.NotEnforced = NotEnforced
	return s
}

func (s *TableConstraintAlterActionRequest) WithValiate(Valiate *bool) *TableConstraintAlterActionRequest {
	s.Valiate = Valiate
	return s
}

func (s *TableConstraintAlterActionRequest) WithNoValidate(NoValidate *bool) *TableConstraintAlterActionRequest {
	s.NoValidate = NoValidate
	return s
}

func (s *TableConstraintAlterActionRequest) WithRely(Rely *bool) *TableConstraintAlterActionRequest {
	s.Rely = Rely
	return s
}

func (s *TableConstraintAlterActionRequest) WithNoRely(NoRely *bool) *TableConstraintAlterActionRequest {
	s.NoRely = NoRely
	return s
}

func NewTableConstraintDropActionRequest() *TableConstraintDropActionRequest {
	return &TableConstraintDropActionRequest{}
}

func (s *TableConstraintDropActionRequest) WithConstraintName(ConstraintName *string) *TableConstraintDropActionRequest {
	s.ConstraintName = ConstraintName
	return s
}

func (s *TableConstraintDropActionRequest) WithPrimaryKey(PrimaryKey *bool) *TableConstraintDropActionRequest {
	s.PrimaryKey = PrimaryKey
	return s
}

func (s *TableConstraintDropActionRequest) WithUnique(Unique *bool) *TableConstraintDropActionRequest {
	s.Unique = Unique
	return s
}

func (s *TableConstraintDropActionRequest) WithForeignKey(ForeignKey *bool) *TableConstraintDropActionRequest {
	s.ForeignKey = ForeignKey
	return s
}

func (s *TableConstraintDropActionRequest) WithColumns(Columns []string) *TableConstraintDropActionRequest {
	s.Columns = Columns
	return s
}

func (s *TableConstraintDropActionRequest) WithCascade(Cascade *bool) *TableConstraintDropActionRequest {
	s.Cascade = Cascade
	return s
}

func (s *TableConstraintDropActionRequest) WithRestrict(Restrict *bool) *TableConstraintDropActionRequest {
	s.Restrict = Restrict
	return s
}

func NewTableExternalTableActionRequest() *TableExternalTableActionRequest {
	return &TableExternalTableActionRequest{}
}

func (s *TableExternalTableActionRequest) WithAdd(Add *TableExternalTableColumnAddActionRequest) *TableExternalTableActionRequest {
	s.Add = Add
	return s
}

func (s *TableExternalTableActionRequest) WithRename(Rename *TableExternalTableColumnRenameActionRequest) *TableExternalTableActionRequest {
	s.Rename = Rename
	return s
}

func (s *TableExternalTableActionRequest) WithDrop(Drop *TableExternalTableColumnDropActionRequest) *TableExternalTableActionRequest {
	s.Drop = Drop
	return s
}

func NewTableSearchOptimizationActionRequest() *TableSearchOptimizationActionRequest {
	return &TableSearchOptimizationActionRequest{}
}

func (s *TableSearchOptimizationActionRequest) WithAddSearchOptimizationOn(AddSearchOptimizationOn []string) *TableSearchOptimizationActionRequest {
	s.AddSearchOptimizationOn = AddSearchOptimizationOn
	return s
}

func (s *TableSearchOptimizationActionRequest) WithDropSearchOptimizationOn(DropSearchOptimizationOn []string) *TableSearchOptimizationActionRequest {
	s.DropSearchOptimizationOn = DropSearchOptimizationOn
	return s
}

func NewTableSetRequest() *TableSetRequest {
	return &TableSetRequest{}
}

func (s *TableSetRequest) WithEnableSchemaEvolution(EnableSchemaEvolution *bool) *TableSetRequest {
	s.EnableSchemaEvolution = EnableSchemaEvolution
	return s
}

func (s *TableSetRequest) WithStageFileFormat(StageFileFormat StageFileFormatRequest) *TableSetRequest {
	s.StageFileFormat = &StageFileFormat
	return s
}

func (s *TableSetRequest) WithStageCopyOptions(StageCopyOptions StageCopyOptionsRequest) *TableSetRequest {
	s.StageCopyOptions = &StageCopyOptions
	return s
}

func (s *TableSetRequest) WithDataRetentionTimeInDays(DataRetentionTimeInDays *int) *TableSetRequest {
	s.DataRetentionTimeInDays = DataRetentionTimeInDays
	return s
}

func (s *TableSetRequest) WithMaxDataExtensionTimeInDays(MaxDataExtensionTimeInDays *int) *TableSetRequest {
	s.MaxDataExtensionTimeInDays = MaxDataExtensionTimeInDays
	return s
}

func (s *TableSetRequest) WithChangeTracking(ChangeTracking *bool) *TableSetRequest {
	s.ChangeTracking = ChangeTracking
	return s
}

func (s *TableSetRequest) WithDefaultDDLCollation(DefaultDDLCollation *string) *TableSetRequest {
	s.DefaultDDLCollation = DefaultDDLCollation
	return s
}

func (s *TableSetRequest) WithComment(Comment *string) *TableSetRequest {
	s.Comment = Comment
	return s
}

func NewTableExternalTableColumnAddActionRequest() *TableExternalTableColumnAddActionRequest {
	return &TableExternalTableColumnAddActionRequest{}
}

func (s *TableExternalTableColumnAddActionRequest) WithIfNotExists() *TableExternalTableColumnAddActionRequest {
	s.IfNotExists = Bool(true)
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithName(Name string) *TableExternalTableColumnAddActionRequest {
	s.Name = Name
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithType(Type DataType) *TableExternalTableColumnAddActionRequest {
	s.Type = Type
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithExpression(Expression string) *TableExternalTableColumnAddActionRequest {
	s.Expression = Expression
	return s
}

func NewTableExternalTableColumnRenameActionRequest() *TableExternalTableColumnRenameActionRequest {
	return &TableExternalTableColumnRenameActionRequest{}
}

func (s *TableExternalTableColumnRenameActionRequest) WithOldName(OldName string) *TableExternalTableColumnRenameActionRequest {
	s.OldName = OldName
	return s
}

func (s *TableExternalTableColumnRenameActionRequest) WithNewName(NewName string) *TableExternalTableColumnRenameActionRequest {
	s.NewName = NewName
	return s
}

func NewTableExternalTableColumnDropActionRequest() *TableExternalTableColumnDropActionRequest {
	return &TableExternalTableColumnDropActionRequest{}
}

func (s *TableExternalTableColumnDropActionRequest) WithColumns(Columns []string) *TableExternalTableColumnDropActionRequest {
	s.Columns = Columns
	return s
}

func (s *TableExternalTableColumnDropActionRequest) WithIfExists() *TableExternalTableColumnDropActionRequest {
	s.IfExists = Bool(true)
	return s
}

func NewShowTableRequest() *ShowTableRequest {
	return &ShowTableRequest{}
}

func (s *ShowTableRequest) WithTerse(Terse *bool) *ShowTableRequest {
	s.terse = Terse
	return s
}

func (s *ShowTableRequest) WithHistory(History *bool) *ShowTableRequest {
	s.history = History
	return s
}

func (s *ShowTableRequest) WithLikePattern(LikePattern string) *ShowTableRequest {
	s.likePattern = LikePattern
	return s
}

func (s *ShowTableRequest) WithIn(in *In) *ShowTableRequest {
	s.in = in
	return s
}

func (s *ShowTableRequest) WithStartsWith(StartsWith *string) *ShowTableRequest {
	s.startsWith = StartsWith
	return s
}

func (s *ShowTableRequest) WithLimitFrom(LimitFrom *LimitFrom) *ShowTableRequest {
	s.limitFrom = LimitFrom
	return s
}

func NewLimitFromRequest() *LimitFromRequest {
	return &LimitFromRequest{}
}

func (s *LimitFromRequest) WithRows(rows *int) *LimitFromRequest {
	s.rows = rows
	return s
}

func (s *LimitFromRequest) WithFrom(from *string) *LimitFromRequest {
	s.from = from
	return s
}

func NewDescribeTableColumnsRequest(
	id SchemaObjectIdentifier,
) *DescribeTableColumnsRequest {
	s := DescribeTableColumnsRequest{}
	s.id = id
	return &s
}

func NewDescribeTableStageRequest(
	id SchemaObjectIdentifier,
) *DescribeTableStageRequest {
	s := DescribeTableStageRequest{}
	s.id = id
	return &s
}
