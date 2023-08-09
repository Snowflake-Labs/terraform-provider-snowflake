package sdk

import "errors"

var (
	_ validatableOpts = (*CreateExternalTableOpts)(nil)
	_ validatableOpts = (*CreateWithManualPartitioningExternalTableOpts)(nil)
	_ validatableOpts = (*CreateDeltaLakeExternalTableOpts)(nil)
	_ validatableOpts = (*AlterExternalTableOptions)(nil)
	_ validatableOpts = (*AlterExternalTablePartitionOptions)(nil)
	_ validatableOpts = (*DropExternalTableOptions)(nil)
	_ validatableOpts = (*ShowExternalTableOptions)(nil)
	_ validatableOpts = (*describeExternalTableColumns)(nil)
	_ validatableOpts = (*describeExternalTableStage)(nil)
)

// +-OK
func (opts *CreateExternalTableOpts) validateProp() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDeltaLakeExternalTableOpts", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "FileFormat"))
	}
	// TODO call validate() underlying props
	return errors.Join(errs...)
}

// +-OK
func (opts *CreateWithManualPartitioningExternalTableOpts) validateProp() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDeltaLakeExternalTableOpts", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "FileFormat"))
	}
	// TODO call validate() underlying props
	return errors.Join(errs...)
}

// +-OK
func (opts *CreateDeltaLakeExternalTableOpts) validateProp() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDeltaLakeExternalTableOpts", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOpts", "FileFormat"))
	}
	// TODO call validate() underlying props
	return errors.Join(errs...)
}

// +-OK
func (opts *AlterExternalTableOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if anyValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.Set, opts.Unset) &&
		!exactlyOneValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.Set, opts.Unset) {
		errs = append(errs, errOneOf("AlterExternalTableOptions", "Refresh", "AddFiles", "RemoveFiles", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		// TODO Check if Tags and AUTO_REFRESH is set ?
	}
	return errors.Join(errs...)
}

func (opts *AlterExternalTablePartitionOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.AddPartitions, opts.DropPartition) {
		errs = append(errs, errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"))
	}
	return errors.Join(errs...)
}

func (opts *DropExternalTableOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.DropOption) {
		if err := opts.DropOption.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *ShowExternalTableOptions) validateProp() error {
	return nil
}

func (v *describeExternalTableColumns) validateProp() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *describeExternalTableStage) validateProp() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (cpp *CloudProviderParams) validate() error {
	if anyValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) && exactlyOneValueSet(cpp.GoogleCloudStorage, cpp.MicrosoftAzure) {
		return errOneOf("CloudProviderParams", "GoogleCloudStorage", "MicrosoftAzure")
	}
	return nil
}

func (opts *ExternalTableFileFormat) validate() error {
	if everyValueSet(opts.Name, opts.Type) {
		return errOneOf("ExternalTableFileFormat", "Name", "Type")
	}
	return nil
}

func (opts *ExternalTableDropOption) validate() error {
	if anyValueSet(opts.Restrict, opts.Cascade) && !exactlyOneValueSet(opts.Restrict, opts.Cascade) {
		return errOneOf("ExternalTableDropOption", "Restrict", "Cascade")
	}
	return nil
}
