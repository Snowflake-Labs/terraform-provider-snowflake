package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ Tables = (*tables)(nil)

type tables struct {
	client *Client
}

func (v *tables) Create(ctx context.Context, request *CreateTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateAsSelect(ctx context.Context, request *CreateTableAsSelectRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateUsingTemplate(ctx context.Context, request *CreateTableUsingTemplateRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateLike(ctx context.Context, request *CreateTableLikeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateClone(ctx context.Context, request *CreateTableCloneRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) Alter(ctx context.Context, request *AlterTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) Drop(ctx context.Context, request *DropTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) Show(ctx context.Context, request *ShowTableRequest) ([]Table, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[tableDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList := convertRows[tableDBRow, Table](dbRows)

	return resultList, nil
}

func (v *tables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Table, error) {
	requestWithName := NewShowTableRequest().WithLikePattern(id.Name())
	tables, err := v.Show(ctx, requestWithName)
	if err != nil {
		return nil, err
	}

	return collections.FindOne(tables, func(r Table) bool { return r.Name == id.Name() })
}

func (s *AlterTableRequest) toOpts() *alterTableOptions {
	var clusteringAction *TableClusteringAction
	if s.ClusteringAction != nil {
		var reclusterAction *TableReclusterAction
		if s.ClusteringAction.Recluster != nil {
			reclusterAction = &TableReclusterAction{
				MaxSize:   s.ClusteringAction.Recluster.MaxSize,
				Condition: s.ClusteringAction.Recluster.Condition,
			}
		}
		var changeReclusterChange *TableReclusterChangeState
		if s.ClusteringAction.ChangeReclusterState != nil {
			changeReclusterChange = &TableReclusterChangeState{State: s.ClusteringAction.ChangeReclusterState}
		}
		clusteringAction = &TableClusteringAction{
			ClusterBy:            s.ClusteringAction.ClusterBy,
			Recluster:            reclusterAction,
			ChangeReclusterState: changeReclusterChange,
			DropClusteringKey:    s.ClusteringAction.DropClusteringKey,
		}
	}
	var columnAction *TableColumnAction
	if s.ColumnAction != nil {
		columnAction = convertTableColumnAction(*s.ColumnAction)
	}
	var constraintAction *TableConstraintAction
	if s.ConstraintAction != nil {
		constraintAction = convertTableConstraintAction(*s.ConstraintAction)
	}
	var externalTableAction *TableExternalTableAction
	if s.ExternalTableAction != nil {
		externalTableAction = convertTableExternalAction(*s.ExternalTableAction)
	}
	var searchOptimizationAction *TableSearchOptimizationAction
	if s.SearchOptimizationAction != nil {
		searchOptimizationAction = convertSearchOptimizationAction(*s.SearchOptimizationAction)
	}
	var tableSet *TableSet
	if s.Set != nil {
		tableSet = convertAlterTableSet(*s.Set)
	}

	tagAssociations := make([]TagAssociation, 0, len(s.SetTags))
	for _, tagRequest := range s.SetTags {
		tagAssociations = append(tagAssociations, TagAssociation(tagRequest))
	}
	var tableUnset *TableUnset
	if s.Unset != nil {
		tableUnset = &TableUnset{
			DataRetentionTimeInDays:    Bool(s.Unset.DataRetentionTimeInDays),
			MaxDataExtensionTimeInDays: Bool(s.Unset.MaxDataExtensionTimeInDays),
			ChangeTracking:             Bool(s.Unset.ChangeTracking),
			DefaultDDLCollation:        Bool(s.Unset.DefaultDDLCollation),
			EnableSchemaEvolution:      Bool(s.Unset.EnableSchemaEvolution),
			Comment:                    Bool(s.Unset.Comment),
		}
	}
	var addRowAccessPolicy *AddRowAccessPolicy
	if s.AddRowAccessPolicy != nil {
		addRowAccessPolicy = &AddRowAccessPolicy{
			PolicyName:  s.AddRowAccessPolicy.PolicyName,
			ColumnNames: s.AddRowAccessPolicy.ColumnName,
		}
	}
	var dropAndAddRowAccessPolicy *DropAndAddRowAccessPolicy
	if s.DropAndAddRowAccessPolicy != nil {
		addRowAccessPolicy := &AddRowAccessPolicy{
			PolicyName:  s.DropAndAddRowAccessPolicy.AddedPolicy.PolicyName,
			ColumnNames: s.DropAndAddRowAccessPolicy.AddedPolicy.ColumnName,
		}
		dropAndAddRowAccessPolicy = &DropAndAddRowAccessPolicy{
			DroppedPolicyName: s.DropAndAddRowAccessPolicy.DroppedPolicyName,
			AddedPolicy:       addRowAccessPolicy,
		}
	}

	return &alterTableOptions{
		IfExists:                  s.IfExists,
		name:                      s.name,
		NewName:                   s.NewName,
		SwapWith:                  s.SwapWith,
		ClusteringAction:          clusteringAction,
		ColumnAction:              columnAction,
		ConstraintAction:          constraintAction,
		ExternalTableAction:       externalTableAction,
		SearchOptimizationAction:  searchOptimizationAction,
		Set:                       tableSet,
		SetTags:                   tagAssociations,
		UnsetTags:                 s.UnsetTags,
		Unset:                     tableUnset,
		AddRowAccessPolicy:        addRowAccessPolicy,
		DropRowAccessPolicy:       s.DropRowAccessPolicy,
		DropAndAddRowAccessPolicy: dropAndAddRowAccessPolicy,
		DropAllAccessRowPolicies:  s.DropAllAccessRowPolicies,
	}
}

func convertAlterTableSet(request TableSetRequest) *TableSet {
	stageFileFormats := make([]StageFileFormat, 0, len(request.StageFileFormat))
	for _, stageFileFormat := range request.StageFileFormat {
		var options *FileFormatTypeOptions
		if stageFileFormat.Options != nil {
			options = stageFileFormat.Options.toOpts()
		}
		stageFileFormats = append(stageFileFormats, StageFileFormat{
			FormatName: stageFileFormat.FormatName,
			Type:       stageFileFormat.Type,
			Options:    options,
		})
	}
	stageCopyOptions := make([]StageCopyOption, 0, len(request.StageCopyOptions))
	for _, stageCopyOption := range request.StageCopyOptions {
		stageCopyOptions = append(stageCopyOptions, StageCopyOption{
			InnerValue: *stageCopyOption.toOpts(),
		})
	}

	return &TableSet{
		EnableSchemaEvolution:      request.EnableSchemaEvolution,
		StageFileFormat:            stageFileFormats,
		StageCopyOptions:           stageCopyOptions,
		DataRetentionTimeInDays:    request.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: request.MaxDataExtensionTimeInDays,
		ChangeTracking:             request.ChangeTracking,
		DefaultDDLCollation:        request.DefaultDDLCollation,
		Comment:                    request.Comment,
	}
}

func convertSearchOptimizationAction(request TableSearchOptimizationActionRequest) *TableSearchOptimizationAction {
	if len(request.AddSearchOptimizationOn) > 0 {
		return &TableSearchOptimizationAction{
			Add: &AddSearchOptimization{
				On: request.AddSearchOptimizationOn,
			},
		}
	}
	if len(request.DropSearchOptimizationOn) > 0 {
		return &TableSearchOptimizationAction{
			Drop: &DropSearchOptimization{
				On: request.DropSearchOptimizationOn,
			},
		}
	}
	return nil
}

func convertTableExternalAction(request TableExternalTableActionRequest) *TableExternalTableAction {
	if request.Add != nil {
		return &TableExternalTableAction{
			Add: &TableExternalTableColumnAddAction{
				Name:       request.Add.Name,
				Type:       request.Add.Type,
				Expression: request.Add.Expression,
			},
		}
	}
	if request.Rename != nil {
		return &TableExternalTableAction{
			Rename: &TableExternalTableColumnRenameAction{
				OldName: request.Rename.OldName,
				NewName: request.Rename.NewName,
			},
		}
	}
	if request.Drop != nil {
		return &TableExternalTableAction{
			Drop: &TableExternalTableColumnDropAction{
				Columns: request.Drop.Columns,
			},
		}
	}
	return nil
}

func convertTableConstraintAction(request TableConstraintActionRequest) *TableConstraintAction {
	if request.Add != nil {
		var foreignKey *OutOfLineForeignKey
		if request.Add.ForeignKey != nil {
			var foreignKeyOnAction *ForeignKeyOnAction
			if request.Add.ForeignKey.On != nil {
				foreignKeyOnAction = &ForeignKeyOnAction{
					OnUpdate: request.Add.ForeignKey.On.OnUpdate,
					OnDelete: request.Add.ForeignKey.On.OnDelete,
				}
			}
			foreignKey = &OutOfLineForeignKey{
				TableName:   request.Add.ForeignKey.TableName,
				ColumnNames: request.Add.ForeignKey.ColumnNames,
				Match:       request.Add.ForeignKey.Match,
				On:          foreignKeyOnAction,
			}
		}
		outOfLineConstrait := AlterOutOfLineConstraint{
			Name:               request.Add.Name,
			Type:               request.Add.Type,
			Columns:            request.Add.Columns,
			ForeignKey:         foreignKey,
			Enforced:           request.Add.Enforced,
			NotEnforced:        request.Add.NotEnforced,
			Deferrable:         request.Add.Deferrable,
			NotDeferrable:      request.Add.NotDeferrable,
			InitiallyDeferred:  request.Add.InitiallyDeferred,
			InitiallyImmediate: request.Add.InitiallyImmediate,
			Enable:             request.Add.Enable,
			Disable:            request.Add.Disable,
			Validate:           request.Add.Validate,
			NoValidate:         request.Add.NoValidate,
			Rely:               request.Add.Rely,
			NoRely:             request.Add.NoRely,
		}
		return &TableConstraintAction{
			Add: &outOfLineConstrait,
		}
	}
	if request.Rename != nil {
		return &TableConstraintAction{
			Rename: &TableConstraintRenameAction{
				OldName: request.Rename.OldName,
				NewName: request.Rename.NewName,
			},
		}
	}
	if request.Alter != nil {
		return &TableConstraintAction{
			Alter: &TableConstraintAlterAction{
				ConstraintName: request.Alter.ConstraintName,
				PrimaryKey:     request.Alter.PrimaryKey,
				Unique:         request.Alter.Unique,
				ForeignKey:     request.Alter.ForeignKey,
				Columns:        request.Alter.Columns,
				Enforced:       request.Alter.Enforced,
				NotEnforced:    request.Alter.NotEnforced,
				Validate:       request.Alter.Valiate,
				NoValidate:     request.Alter.NoValidate,
				Rely:           request.Alter.Rely,
				NoRely:         request.Alter.NoRely,
			},
		}
	}
	if request.Drop != nil {
		return &TableConstraintAction{
			Drop: &TableConstraintDropAction{
				ConstraintName: request.Drop.ConstraintName,
				PrimaryKey:     request.Drop.PrimaryKey,
				Unique:         request.Drop.Unique,
				ForeignKey:     request.Drop.ForeignKey,
				Columns:        request.Drop.Columns,
				Cascade:        request.Drop.Cascade,
				Restrict:       request.Drop.Restrict,
			},
		}
	}
	return nil
}

func convertTableColumnAction(request TableColumnActionRequest) *TableColumnAction {
	if request.Add != nil {
		var defaultValue *ColumnDefaultValue
		if request.Add.DefaultValue != nil {
			defaultValue = &ColumnDefaultValue{
				request.Add.DefaultValue.expression,
				&ColumnIdentity{
					Start:     request.Add.DefaultValue.identity.Start,
					Increment: request.Add.DefaultValue.identity.Increment,
				},
			}
		}
		var inlineConstraint *TableColumnAddInlineConstraint
		if request.Add.InlineConstraint != nil {
			var foreignKey *ColumnAddForeignKey
			if request.Add.InlineConstraint.ForeignKey != nil {
				foreignKey = &ColumnAddForeignKey{
					TableName:  request.Add.InlineConstraint.ForeignKey.TableName,
					ColumnName: request.Add.InlineConstraint.ForeignKey.ColumnName,
				}
			}
			inlineConstraint = &TableColumnAddInlineConstraint{
				NotNull:    request.Add.InlineConstraint.NotNull,
				Name:       request.Add.InlineConstraint.Name,
				Type:       request.Add.InlineConstraint.Type,
				ForeignKey: foreignKey,
			}
		}
		return &TableColumnAction{
			Add: &TableColumnAddAction{
				Column:           request.Add.Column,
				Name:             request.Add.Name,
				Type:             request.Add.Type,
				DefaultValue:     defaultValue,
				InlineConstraint: inlineConstraint,
			},
		}
	}
	if request.Rename != nil {
		return &TableColumnAction{
			Rename: &TableColumnRenameAction{
				OldName: request.Rename.OldName,
				NewName: request.Rename.NewName,
			},
		}
	}
	if len(request.Alter) > 0 {
		var alterActions []TableColumnAlterAction
		for _, alterAction := range request.Alter {
			var notNullConstraint *TableColumnNotNullConstraint
			if alterAction.NotNullConstraint != nil {
				notNullConstraint = &TableColumnNotNullConstraint{
					Set:  alterAction.NotNullConstraint.Set,
					Drop: alterAction.NotNullConstraint.Drop,
				}
			}
			alterActions = append(alterActions, TableColumnAlterAction{
				Column:            Bool(alterAction.Column),
				Name:              alterAction.Name,
				DropDefault:       alterAction.DropDefault,
				SetDefault:        alterAction.SetDefault,
				NotNullConstraint: notNullConstraint,
				Type:              alterAction.Type,
				Comment:           alterAction.Comment,
				UnsetComment:      alterAction.UnsetComment,
			})
		}
		return &TableColumnAction{
			Alter: alterActions,
		}
	}
	if request.SetMaskingPolicy != nil {
		return &TableColumnAction{
			SetMaskingPolicy: &TableColumnAlterSetMaskingPolicyAction{
				ColumnName:        request.SetMaskingPolicy.ColumnName,
				MaskingPolicyName: request.SetMaskingPolicy.MaskingPolicyName,
				Using:             request.SetMaskingPolicy.Using,
				Force:             request.SetMaskingPolicy.Force,
			},
		}
	}
	if request.UnsetMaskingPolicy != nil {
		return &TableColumnAction{
			UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{
				ColumnName: request.UnsetMaskingPolicy.ColumnName,
			},
		}
	}
	if request.SetTags != nil {
		return &TableColumnAction{
			SetTags: &TableColumnAlterSetTagsAction{
				ColumnName: request.SetTags.ColumnName,
				Tags:       request.SetTags.Tags,
			},
		}
	}
	if request.UnsetTags != nil {
		return &TableColumnAction{
			UnsetTags: &TableColumnAlterUnsetTagsAction{
				ColumnName: request.UnsetTags.ColumnName,
				Tags:       request.UnsetTags.Tags,
			},
		}
	}
	if len(request.DropColumns) > 0 {
		return &TableColumnAction{
			DropColumns: &TableColumnAlterDropColumns{
				Columns: request.DropColumns,
			},
		}
	}
	return nil
}

func (s *CreateTableRequest) toOpts() *createTableOptions {
	tagAssociations := make([]TagAssociation, 0, len(s.Tags))
	for _, tagRequest := range s.Tags {
		tagAssociations = append(tagAssociations, TagAssociation(tagRequest))
	}
	var rowAccessPolicy *RowAccessPolicy
	if s.RowAccessPolicy != nil {
		rowAccessPolicy = &RowAccessPolicy{
			Name: s.RowAccessPolicy.Name,
			On:   s.RowAccessPolicy.On,
		}
	}
	var outOfLineConstrait *CreateOutOfLineConstraint
	if s.OutOfLineConstraint != nil {
		var foreignKey *OutOfLineForeignKey
		if s.OutOfLineConstraint.ForeignKey != nil {
			var foreignKeyOnAction *ForeignKeyOnAction
			if s.OutOfLineConstraint.ForeignKey.On != nil {
				foreignKeyOnAction = &ForeignKeyOnAction{
					OnUpdate: s.OutOfLineConstraint.ForeignKey.On.OnUpdate,
					OnDelete: s.OutOfLineConstraint.ForeignKey.On.OnDelete,
				}
			}
			foreignKey = &OutOfLineForeignKey{
				TableName:   s.OutOfLineConstraint.ForeignKey.TableName,
				ColumnNames: s.OutOfLineConstraint.ForeignKey.ColumnNames,
				Match:       s.OutOfLineConstraint.ForeignKey.Match,
				On:          foreignKeyOnAction,
			}
		}
		outOfLineConstrait = &CreateOutOfLineConstraint{
			Name:               s.OutOfLineConstraint.Name,
			Type:               s.OutOfLineConstraint.Type,
			Columns:            s.OutOfLineConstraint.Columns,
			ForeignKey:         foreignKey,
			Enforced:           s.OutOfLineConstraint.Enforced,
			NotEnforced:        s.OutOfLineConstraint.NotEnforced,
			Deferrable:         s.OutOfLineConstraint.Deferrable,
			NotDeferrable:      s.OutOfLineConstraint.NotDeferrable,
			InitiallyDeferred:  s.OutOfLineConstraint.InitiallyDeferred,
			InitiallyImmediate: s.OutOfLineConstraint.InitiallyImmediate,
			Enable:             s.OutOfLineConstraint.Enable,
			Disable:            s.OutOfLineConstraint.Disable,
			Validate:           s.OutOfLineConstraint.Validate,
			NoValidate:         s.OutOfLineConstraint.NoValidate,
			Rely:               s.OutOfLineConstraint.Rely,
			NoRely:             s.OutOfLineConstraint.NoRely,
		}
	}

	return &createTableOptions{
		OrReplace:                  s.orReplace,
		IfNotExists:                s.ifNotExists,
		Scope:                      s.scope,
		Kind:                       s.kind,
		name:                       s.name,
		Columns:                    convertColumns(s.columns),
		ClusterBy:                  s.clusterBy,
		OutOfLineConstraint:        outOfLineConstrait,
		EnableSchemaEvolution:      s.enableSchemaEvolution,
		StageCopyOptions:           convertStageCopyOptions(s.stageCopyOptions),
		StageFileFormat:            convertStageFileFormatOptions(s.stageFileFormat),
		DataRetentionTimeInDays:    s.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: s.MaxDataExtensionTimeInDays,
		ChangeTracking:             s.ChangeTracking,
		DefaultDDLCollation:        s.DefaultDDLCollation,
		CopyGrants:                 s.CopyGrants,
		Tags:                       tagAssociations,
		Comment:                    s.Comment,
		RowAccessPolicy:            rowAccessPolicy,
	}
}

func (s *CreateTableAsSelectRequest) toOpts() *createTableAsSelectOptions {
	columns := make([]TableAsSelectColumn, 0, len(s.columns))
	for _, column := range s.columns {
		var maskingPolicy *TableAsSelectColumnMaskingPolicy
		if column.maskingPolicyName != nil {
			maskingPolicy = &TableAsSelectColumnMaskingPolicy{
				Name: *column.maskingPolicyName,
			}
		}
		columns = append(columns, TableAsSelectColumn{
			Name:          column.name,
			Type:          column.type_,
			MaskingPolicy: maskingPolicy,
		})
	}
	return &createTableAsSelectOptions{
		OrReplace: s.orReplace,
		name:      s.name,
		Columns:   columns,
	}
}

func (s *CreateTableUsingTemplateRequest) toOpts() *createTableUsingTemplateOptions {
	return &createTableUsingTemplateOptions{
		OrReplace:  s.orReplace,
		name:       s.name,
		CopyGrants: s.copyGrants,
		Query:      []string{s.Query},
	}
}

func (s *CreateTableLikeRequest) toOpts() *createTableLikeOptions {
	return &createTableLikeOptions{
		OrReplace:   s.orReplace,
		name:        s.name,
		CopyGrants:  s.copyGrants,
		SourceTable: s.sourceTable,
		ClusterBy:   s.clusterBy,
	}
}

func (s *CreateTableCloneRequest) toOpts() *createTableCloneOptions {
	var clonePoint *ClonePoint
	if s.ClonePoint != nil {
		clonePoint = &ClonePoint{
			Moment: s.ClonePoint.Moment,
			At: TimeTravel{
				Timestamp: s.ClonePoint.At.Timestamp,
				Offset:    s.ClonePoint.At.Offset,
				Statement: s.ClonePoint.At.Statement,
			},
		}
	}
	return &createTableCloneOptions{
		OrReplace:   s.orReplace,
		name:        s.name,
		CopyGrants:  s.copyGrants,
		SourceTable: s.sourceTable,
		ClonePoint:  clonePoint,
	}
}

func convertStageCopyOptions(copyOptionRequests []StageCopyOptionsRequest) []StageCopyOption {
	copyOptions := make([]StageCopyOption, 0, len(copyOptionRequests))
	for _, request := range copyOptionRequests {
		innerValue := request.toOpts()
		copyOptions = append(copyOptions, StageCopyOption{
			InnerValue: *innerValue,
		})
	}
	return copyOptions
}

func (v *StageCopyOptionsRequest) toOpts() *StageCopyOptionsInnerValue {
	return &StageCopyOptionsInnerValue{
		OnError:           v.OnError.toOpts(),
		SizeLimit:         v.SizeLimit,
		Purge:             v.Purge,
		ReturnFailedOnly:  v.ReturnFailedOnly,
		MatchByColumnName: v.MatchByColumnName,
		EnforceLength:     v.EnforceLength,
		TruncateColumns:   v.Truncatecolumns,
		Force:             v.Force,
	}
}

func (v *StageCopyOnErrorOptionsRequest) toOpts() *StageCopyOnErrorOptions {
	return &StageCopyOnErrorOptions{
		Continue:       v.Continue,
		SkipFile:       v.SkipFile,
		AbortStatement: v.AbortStatement,
	}
}

func convertStageFileFormatOptions(stageFileFormatRequests []StageFileFormatRequest) []StageFileFormat {
	fileFormats := make([]StageFileFormat, 0, len(stageFileFormatRequests))
	for _, request := range stageFileFormatRequests {
		var options *FileFormatTypeOptions
		if request.Options != nil {
			options = request.Options.toOpts()
		}
		format := StageFileFormat{
			FormatName: request.FormatName,
			Type:       request.Type,
			Options:    options,
		}
		fileFormats = append(fileFormats, format)
	}
	return fileFormats
}

func (v *FileFormatTypeOptionsRequest) toOpts() *FileFormatTypeOptions {
	if v == nil {
		return nil
	}
	return &FileFormatTypeOptions{
		CSVCompression:                  v.CSVCompression,
		CSVRecordDelimiter:              v.CSVRecordDelimiter,
		CSVFieldDelimiter:               v.CSVFieldDelimiter,
		CSVFileExtension:                v.CSVFileExtension,
		CSVParseHeader:                  v.CSVParseHeader,
		CSVSkipHeader:                   v.CSVSkipHeader,
		CSVSkipBlankLines:               v.CSVSkipBlankLines,
		CSVDateFormat:                   v.CSVDateFormat,
		CSVTimeFormat:                   v.CSVTimeFormat,
		CSVTimestampFormat:              v.CSVTimestampFormat,
		CSVBinaryFormat:                 v.CSVBinaryFormat,
		CSVEscape:                       v.CSVEscape,
		CSVEscapeUnenclosedField:        v.CSVEscapeUnenclosedField,
		CSVTrimSpace:                    v.CSVTrimSpace,
		CSVFieldOptionallyEnclosedBy:    v.CSVFieldOptionallyEnclosedBy,
		CSVNullIf:                       v.CSVNullIf,
		CSVErrorOnColumnCountMismatch:   v.CSVErrorOnColumnCountMismatch,
		CSVReplaceInvalidCharacters:     v.CSVReplaceInvalidCharacters,
		CSVEmptyFieldAsNull:             v.CSVEmptyFieldAsNull,
		CSVSkipByteOrderMark:            v.CSVSkipByteOrderMark,
		CSVEncoding:                     v.CSVEncoding,
		JSONCompression:                 v.JSONCompression,
		JSONDateFormat:                  v.JSONDateFormat,
		JSONTimeFormat:                  v.JSONTimeFormat,
		JSONTimestampFormat:             v.JSONTimestampFormat,
		JSONBinaryFormat:                v.JSONBinaryFormat,
		JSONTrimSpace:                   v.JSONTrimSpace,
		JSONNullIf:                      v.JSONNullIf,
		JSONFileExtension:               v.JSONFileExtension,
		JSONEnableOctal:                 v.JSONEnableOctal,
		JSONAllowDuplicate:              v.JSONAllowDuplicate,
		JSONStripOuterArray:             v.JSONStripOuterArray,
		JSONStripNullValues:             v.JSONStripNullValues,
		JSONReplaceInvalidCharacters:    v.JSONReplaceInvalidCharacters,
		JSONIgnoreUTF8Errors:            v.JSONIgnoreUTF8Errors,
		JSONSkipByteOrderMark:           v.JSONSkipByteOrderMark,
		AvroCompression:                 v.AvroCompression,
		AvroTrimSpace:                   v.AvroTrimSpace,
		AvroReplaceInvalidCharacters:    v.AvroReplaceInvalidCharacters,
		AvroNullIf:                      v.AvroNullIf,
		ORCTrimSpace:                    v.ORCTrimSpace,
		ORCReplaceInvalidCharacters:     v.ORCReplaceInvalidCharacters,
		ORCNullIf:                       v.ORCNullIf,
		ParquetCompression:              v.ParquetCompression,
		ParquetSnappyCompression:        v.ParquetSnappyCompression,
		ParquetBinaryAsText:             v.ParquetBinaryAsText,
		ParquetTrimSpace:                v.ParquetTrimSpace,
		ParquetReplaceInvalidCharacters: v.ParquetReplaceInvalidCharacters,
		ParquetNullIf:                   v.ParquetNullIf,
		XMLCompression:                  v.XMLCompression,
		XMLIgnoreUTF8Errors:             v.XMLIgnoreUTF8Errors,
		XMLPreserveSpace:                v.XMLPreserveSpace,
		XMLStripOuterElement:            v.XMLStripOuterElement,
		XMLDisableSnowflakeData:         v.XMLDisableSnowflakeData,
		XMLDisableAutoConvert:           v.XMLDisableAutoConvert,
		XMLReplaceInvalidCharacters:     v.XMLReplaceInvalidCharacters,
		XMLSkipByteOrderMark:            v.XMLSkipByteOrderMark,
	}
}

func convertColumns(columnRequests []TableColumnRequest) []TableColumn {
	columns := make([]TableColumn, 0, len(columnRequests))
	for _, columnRequest := range columnRequests {
		var defaultValue *ColumnDefaultValue
		if columnRequest.defaultValue != nil {
			var columnIdentity *ColumnIdentity
			if columnRequest.defaultValue.identity != nil {
				columnIdentity = &ColumnIdentity{
					Start:     columnRequest.defaultValue.identity.Start,
					Increment: columnRequest.defaultValue.identity.Increment,
				}
			}
			defaultValue = &ColumnDefaultValue{
				columnRequest.defaultValue.expression,
				columnIdentity,
			}
		}
		var inlineConstraint *ColumnInlineConstraint
		if columnRequest.inlineConstraint != nil {
			var foreignKey *InlineForeignKey
			if columnRequest.inlineConstraint.foreignKey != nil {
				var onActionRequest *ForeignKeyOnAction
				if columnRequest.inlineConstraint.foreignKey.On != nil {
					onActionRequest = &ForeignKeyOnAction{
						OnUpdate: columnRequest.inlineConstraint.foreignKey.On.OnUpdate,
						OnDelete: columnRequest.inlineConstraint.foreignKey.On.OnDelete,
					}
				}
				foreignKey = &InlineForeignKey{
					TableName:  columnRequest.inlineConstraint.foreignKey.TableName,
					ColumnName: columnRequest.inlineConstraint.foreignKey.ColumnName,
					Match:      columnRequest.inlineConstraint.foreignKey.Match,
					On:         onActionRequest,
				}
			}
			inlineConstraint = &ColumnInlineConstraint{
				NotNull:            columnRequest.notNull,
				Name:               &columnRequest.inlineConstraint.Name,
				Type:               &columnRequest.inlineConstraint.type_,
				ForeignKey:         foreignKey,
				Enforced:           columnRequest.inlineConstraint.enforced,
				NotEnforced:        columnRequest.inlineConstraint.notEnforced,
				Deferrable:         columnRequest.inlineConstraint.deferrable,
				NotDeferrable:      columnRequest.inlineConstraint.notDeferrable,
				InitiallyDeferred:  columnRequest.inlineConstraint.initiallyDeferred,
				InitiallyImmediate: columnRequest.inlineConstraint.initiallyImmediate,
				Enable:             columnRequest.inlineConstraint.enable,
				Disable:            columnRequest.inlineConstraint.disable,
				Validate:           columnRequest.inlineConstraint.validate,
				NoValidate:         columnRequest.inlineConstraint.noValidate,
				Rely:               columnRequest.inlineConstraint.rely,
				NoRely:             columnRequest.inlineConstraint.noRely,
			}
		}
		var maskingPolicy *ColumnMaskingPolicy
		if columnRequest.maskingPolicy != nil {
			maskingPolicy = &ColumnMaskingPolicy{
				With:  columnRequest.maskingPolicy.with,
				Name:  columnRequest.maskingPolicy.name,
				Using: columnRequest.maskingPolicy.using,
			}
		}
		columns = append(columns, TableColumn{
			Name:             columnRequest.name,
			Type:             columnRequest.type_,
			Collate:          columnRequest.collate,
			Comment:          columnRequest.comment,
			DefaultValue:     defaultValue,
			MaskingPolicy:    maskingPolicy,
			NotNull:          columnRequest.notNull,
			With:             columnRequest.with,
			Tags:             columnRequest.tags,
			InlineConstraint: inlineConstraint,
		})
	}
	return columns
}
