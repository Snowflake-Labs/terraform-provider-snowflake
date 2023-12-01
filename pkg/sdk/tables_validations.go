package sdk

import "errors"

var (
	_ validatable = new(createTableOptions)
	_ validatable = new(createTableAsSelectOptions)
	_ validatable = new(createTableLikeOptions)
	_ validatable = new(createTableCloneOptions)
	_ validatable = new(createTableUsingTemplateOptions)
	_ validatable = new(alterTableOptions)
	_ validatable = new(dropTableOptions)
	_ validatable = new(showTableOptions)
	_ validatable = new(describeTableColumnsOptions)
	_ validatable = new(describeTableStageOptions)
)

func (opts *createTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if len(opts.Columns) == 0 {
		errs = append(errs, errNotSet("createTableOptions", "Columns"))
	}
	for _, column := range opts.Columns {
		if column.DefaultValue != nil {
			if ok := exactlyOneValueSet(
				column.DefaultValue.Expression,
				column.DefaultValue.Identity,
			); !ok {
				errs = append(errs, errExactlyOneOf("DefaultValue", "Expression", "Identity"))
			}
		}
		if column.MaskingPolicy != nil {
			if !ValidObjectIdentifier(column.MaskingPolicy.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
		for _, tag := range column.Tags {
			if !ValidObjectIdentifier(tag.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	if outOfLineConstraint := opts.OutOfLineConstraint; valueSet(outOfLineConstraint) {
		if foreignKey := outOfLineConstraint.ForeignKey; valueSet(foreignKey) {
			if !ValidObjectIdentifier(foreignKey.TableName) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	for _, stageFileFormat := range opts.StageFileFormat {
		if ok := exactlyOneValueSet(
			stageFileFormat.FormatName,
			stageFileFormat.Type,
		); !ok {
			errs = append(errs, errExactlyOneOf("StageFileFormat", "FormatName", "FormatType"))
		}
	}

	if opts.RowAccessPolicy != nil {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}

	return errors.Join(errs...)
}

func (opts *createTableAsSelectOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if len(opts.Columns) == 0 {
		errs = append(errs, errNotSet("createTableAsSelectOptions", "Columns"))
	}
	return errors.Join(errs...)
}

func (opts *createTableLikeOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.SourceTable) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *createTableUsingTemplateOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *createTableCloneOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.SourceTable) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *alterTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
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
		errs = append(errs, errExactlyOneOf("alterTableOptions", "NewName", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllAccessRowPolicies"))
	}
	if opts.NewName != nil {
		if !ValidObjectIdentifier(*opts.NewName) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if opts.SwapWith != nil {
		if !ValidObjectIdentifier(*opts.SwapWith) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	if clusteringAction := opts.ClusteringAction; valueSet(clusteringAction) {
		if ok := exactlyOneValueSet(
			clusteringAction.ClusterBy,
			clusteringAction.Recluster,
			clusteringAction.ChangeReclusterState,
			clusteringAction.DropClusteringKey,
		); !ok {
			errs = append(errs, errExactlyOneOf("ClusteringAction", "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey"))
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
			errs = append(errs, errExactlyOneOf("ColumnAction", "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetTags", "UnsetTags", "DropColumns"))
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
				errs = append(errs, errExactlyOneOf("TableColumnAlterAction", "DropDefault", "SetDefault", "NotNullConstraint", "Type", "Comment", "UnsetComment"))
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
			); !ok {
				errs = append(errs, errExactlyOneOf("TableConstraintAlterAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
			}
		}
		if dropAction := constraintAction.Drop; valueSet(dropAction) {
			if ok := exactlyOneValueSet(
				dropAction.ConstraintName,
				dropAction.PrimaryKey,
				dropAction.Unique,
				dropAction.ForeignKey,
			); !ok {
				errs = append(errs, errExactlyOneOf("TableConstraintDropAction", "ConstraintName", "PrimaryKey", "Unique", "ForeignKey", "Columns"))
			}
		}
	}
	if externalAction := opts.ExternalTableAction; valueSet(externalAction) {
		if ok := exactlyOneValueSet(
			externalAction.Add,
			externalAction.Rename,
			externalAction.Drop,
		); !ok {
			errs = append(errs, errExactlyOneOf("TableExternalTableAction", "Add", "Rename", "Drop"))
		}
	}
	if searchOptimizationAction := opts.SearchOptimizationAction; valueSet(searchOptimizationAction) {
		if ok := exactlyOneValueSet(
			searchOptimizationAction.Add,
			searchOptimizationAction.Drop,
		); !ok {
			errs = append(errs, errExactlyOneOf("TableSearchOptimizationAction", "Add", "Drop"))
		}
	}
	return errors.Join(errs...)
}

func (opts *dropTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *showTableOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}

func (opts *describeTableColumnsOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (opts *describeTableStageOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}
