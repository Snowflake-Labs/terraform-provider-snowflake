package sdk

import "errors"

var (
	_ validatableOpts = new(createTableOptions)
	_ validatableOpts = new(createTableAsSelectOptions)
	_ validatableOpts = new(createTableLikeOptions)
	_ validatableOpts = new(createTableCloneOptions)
	_ validatableOpts = new(createTableUsingTemplateOptions)
	_ validatableOpts = new(alterTableOptions)
	_ validatableOpts = new(dropTableOptions)
	_ validatableOpts = new(showTableOptions)
)

func (opts *createTableOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	if len(opts.Columns) == 0 {
		errs = append(errs, errTableNeedsAtLeastOneColumn)
	}
	for _, column := range opts.Columns {
		if column.DefaultValue != nil {
			if ok := exactlyOneValueSet(
				column.DefaultValue.Expression,
				column.DefaultValue.Identity,
			); !ok {
				errs = append(errs, errColumnDefaultValueNeedsExactlyOneValue)
			}
		}
		if column.MaskingPolicy != nil {
			if !validObjectidentifier(column.MaskingPolicy.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier2())
			}
		}
		for _, tag := range column.Tags {
			if !validObjectidentifier(tag.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier2())
			}

		}
	}
	if outOfLineConstraint := opts.OutOfLineConstraint; valueSet(outOfLineConstraint) {
		if foreignKey := outOfLineConstraint.ForeignKey; valueSet(foreignKey) {
			if !validObjectidentifier(foreignKey.TableName) {
				errs = append(errs, ErrInvalidObjectIdentifier2())
			}
		}
	}
	for _, stageFileFormat := range opts.StageFileFormat {
		if ok := exactlyOneValueSet(
			stageFileFormat.InnerValue.FormatName,
			stageFileFormat.InnerValue.FormatType,
		); !ok {
			errs = append(errs, errStageFileFormatValueNeedsExactlyOneValue)
		}
	}

	if opts.RowAccessPolicy != nil {
		if !validObjectidentifier(opts.RowAccessPolicy.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier2())
		}
	}

	return errors.Join(errs...)
}

func (opts *createTableAsSelectOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	if len(opts.Columns) == 0 {
		errs = append(errs, errTableNeedsAtLeastOneColumn)
	}
	return errors.Join(errs...)
}
func (opts *createTableLikeOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	if !validObjectidentifier(opts.SourceTable) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	return errors.Join(errs...)
}
func (opts *createTableUsingTemplateOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	return errors.Join(errs...)
}
func (opts *createTableCloneOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	if !validObjectidentifier(opts.SourceTable) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	return errors.Join(errs...)
}

func (opts *alterTableOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	if ok := exactlyOneValueSet(
		opts.NewName,
		opts.SwapWith,
		opts.ClusteringAction,
		opts.ColumnAction,
		opts.ConstraintAction,
		opts.ExternalTableAction,
		opts.SearchOptimizationAction,
		opts.Set,
		opts.SetTags,
		opts.UnsetTags,
		opts.Unset,
		opts.AddRowAccessPolicy,
		opts.DropRowAccessPolicy,
		opts.DropAndAddRowAccessPolicy,
		opts.DropAllAccessRowPolicies,
	); !ok {
		errs = append(errs, errAlterTableNeedsExactlyOneAction)
	}
	if opts.NewName != nil {
		if !validObjectidentifier(*opts.NewName) {
			errs = append(errs, ErrInvalidObjectIdentifier2())
		}
	}
	if opts.SwapWith != nil {
		if !validObjectidentifier(*opts.SwapWith) {
			errs = append(errs, ErrInvalidObjectIdentifier2())
		}
	}
	if clusteringAction := opts.ClusteringAction; valueSet(clusteringAction) {
		if ok := exactlyOneValueSet(
			clusteringAction.ClusterBy,
			clusteringAction.Recluster,
			clusteringAction.ChangeReclusterState,
			clusteringAction.DropClusteringKey,
		); !ok {
			errs = append(errs, errTableClusteringActionNeedsExactlyOneAction)
		}
	}
	if columnAction := opts.ColumnAction; valueSet(columnAction) {
		if ok := exactlyOneValueSet(
			columnAction.Add,
			columnAction.Rename,
			columnAction.Alter,
			columnAction.SetMaskingPolicy,
			columnAction.UnsetMaskingPolicy,
			columnAction.SetTags,
			columnAction.UnsetTags,
			columnAction.DropColumns,
		); !ok {
			errs = append(errs, errTableColumnActionNeedsExactlyOneAction)
		}
		for _, alterAction := range columnAction.Alter {
			if ok := exactlyOneValueSet(
				alterAction.DropDefault,
				alterAction.SetDefault,
				alterAction.NotNullConstraint,
				alterAction.Type,
				alterAction.Comment,
				alterAction.UnsetComment,
			); !ok {
				errs = append(errs, errTableColumnAlterActionNeedsExactlyOneAction)
			}
		}
	}
	if constraintAction := opts.ConstraintAction; valueSet(constraintAction) {
		if alterAction := constraintAction.Alter; valueSet(alterAction) {
			if ok := exactlyOneValueSet(
				alterAction.ConstraintName,
				alterAction.PrimaryKey,
				alterAction.Unique,
				alterAction.ForeignKey,
				alterAction.Columns,
			); !ok {
				errs = append(errs, errTableConstraintAlterActionNeedsExactlyOneAction)
			}
		}
		if dropAction := constraintAction.Drop; valueSet(dropAction) {
			if ok := exactlyOneValueSet(
				dropAction.ConstraintName,
				dropAction.PrimaryKey,
				dropAction.Unique,
				dropAction.ForeignKey,
				dropAction.Columns,
			); !ok {
				errs = append(errs, errTableConstraintDropActionNeedsExactlyOneAction)
			}
		}
	}
	if externalAction := opts.ExternalTableAction; valueSet(externalAction) {
		if ok := exactlyOneValueSet(
			externalAction.Add,
			externalAction.Rename,
			externalAction.Drop,
		); !ok {
			errs = append(errs, errTableExternalActionNeedsExactlyOneAction)
		}
	}
	if searchOptimizationAction := opts.SearchOptimizationAction; valueSet(searchOptimizationAction) {
		if ok := exactlyOneValueSet(
			searchOptimizationAction.Add,
			searchOptimizationAction.Drop,
		); !ok {
			errs = append(errs, errTableSearchOptimizationActionNeedsExactlyOneAction)
		}
	}
	return errors.Join(errs...)
}

func (opts *dropTableOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier2())
	}
	return errors.Join(errs...)
}
func (opts *showTableOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, errPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}

var (
	errTableNeedsAtLeastOneColumn               = errors.New("table create statement needs at least one column")
	errColumnDefaultValueNeedsExactlyOneValue   = errors.New("column default value needs exactly one of {Expression, Identity}")
	errStageFileFormatValueNeedsExactlyOneValue = errors.New("stage file format value needs exactly one of {FormatName, FormatType}")
	errAlterTableNeedsExactlyOneAction          = errors.New(`
stage file format value needs exactly one of {
     NewName,
     SwapWith,
     ClusteringAction,
     ColumnAction,
     ConstraintAction,
     ExternalTableAction,
     SearchOptimizationAction,
     Set,
     SetTags,
     UnsetTags,
     Unset,
     AddRowAccessPolicy,
     DropRowAccessPolicy,
     DropAndAddRowAccessPolicy,
     DropAllAccessRowPolicies,
}`)
	errTableClusteringActionNeedsExactlyOneAction         = errors.New("alter table clustering action needs exactly one of {ClusterBy, Recluster, ChangeReclusterState,DropClusteringKey}")
	errTableColumnActionNeedsExactlyOneAction             = errors.New("alter table column action needs exactly one of {Add,Rename,Alter,SetMaskingPolicy,UnsetMaskingPolicy,SetTags,UnsetTags,DropColumns}")
	errTableColumnAlterActionNeedsExactlyOneAction        = errors.New("alter table column alter action needs exactly one of {DropDefault,SetDefault,NotNullConstraint,Type,Comment,UnsetComment}")
	errTableConstraintAlterActionNeedsExactlyOneAction    = errors.New("alter table constraint alter action needs exactly one of {ConstraintName,PrimaryKey,Unique,ForeignKey,Columns}")
	errTableConstraintDropActionNeedsExactlyOneAction     = errors.New("alter table constraint drop action needs exactly one of {ConstraintName,PrimaryKey,Unique,ForeignKey,Columns}")
	errTableExternalActionNeedsExactlyOneAction           = errors.New("alter table external action needs exactly one of {Add, Rename, Drop}")
	errTableSearchOptimizationActionNeedsExactlyOneAction = errors.New("alter table search optimization action needs exactly one of {Add, Drop}")
)
