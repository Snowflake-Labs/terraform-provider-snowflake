package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = (*CreateExternalTableOpts)(nil)
	_ validatable = (*CreateWithManualPartitioningExternalTableOpts)(nil)
	_ validatable = (*CreateDeltaLakeExternalTableOpts)(nil)
	_ validatable = (*CreateExternalTableUsingTemplateOpts)(nil)
	_ validatable = (*AlterExternalTableOptions)(nil)
	_ validatable = (*AlterExternalTablePartitionOptions)(nil)
	_ validatable = (*DropExternalTableOptions)(nil)
	_ validatable = (*ShowExternalTableOptions)(nil)
	_ validatable = (*describeExternalTableColumns)(nil)
	_ validatable = (*describeExternalTableStage)(nil)
)

func (opts *CreateExternalTableOpts) validate() error {
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

func (opts *CreateWithManualPartitioningExternalTableOpts) validate() error {
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

func (opts *CreateDeltaLakeExternalTableOpts) validate() error {
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

func (opts *CreateExternalTableUsingTemplateOpts) validate() error {
	var errs []error
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

func (opts *AlterExternalTableOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if anyValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.AutoRefresh, opts.SetTag, opts.UnsetTag) &&
		!exactlyOneValueSet(opts.Refresh, opts.AddFiles, opts.RemoveFiles, opts.AutoRefresh, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errOneOf("AlterExternalTableOptions", "Refresh", "AddFiles", "RemoveFiles", "AutoRefresh", "SetTag", "UnsetTag"))
	}
	return errors.Join(errs...)
}

func (opts *AlterExternalTablePartitionOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.AddPartitions, opts.DropPartition) {
		errs = append(errs, errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"))
	}
	return errors.Join(errs...)
}

func (opts *DropExternalTableOptions) validate() error {
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

func (opts *ShowExternalTableOptions) validate() error {
	return nil
}

func (v *describeExternalTableColumns) validate() error {
	if !validObjectidentifier(v.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *describeExternalTableStage) validate() error {
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
	var errs []error
	if everyValueSet(opts.Name, opts.Type) {
		errs = append(errs, errOneOf("ExternalTableFileFormat", "Name", "Type"))
	}
	fields := externalTableFileFormatTypeOptionsFieldsByType(opts.Options)
	for formatType := range fields {
		if *opts.Type == formatType {
			continue
		}
		if anyValueSet(fields[formatType]...) {
			errs = append(errs, fmt.Errorf("cannot set %s fields when TYPE = %s", formatType, *opts.Type))
		}
	}
	return errors.Join(errs...)
}

func (opts *ExternalTableDropOption) validate() error {
	if anyValueSet(opts.Restrict, opts.Cascade) && !exactlyOneValueSet(opts.Restrict, opts.Cascade) {
		return errOneOf("ExternalTableDropOption", "Restrict", "Cascade")
	}
	return nil
}

func externalTableFileFormatTypeOptionsFieldsByType(opts *ExternalTableFileFormatTypeOptions) map[ExternalTableFileFormatType][]any {
	return map[ExternalTableFileFormatType][]any{
		ExternalTableFileFormatTypeCSV: {
			opts.CSVCompression,
			opts.CSVRecordDelimiter,
			opts.CSVFieldDelimiter,
			opts.CSVSkipHeader,
			opts.CSVSkipBlankLines,
			opts.CSVEscapeUnenclosedField,
			opts.CSVTrimSpace,
			opts.CSVFieldOptionallyEnclosedBy,
			opts.CSVNullIf,
			opts.CSVEmptyFieldAsNull,
			opts.CSVEncoding,
		},
		ExternalTableFileFormatTypeJSON: {
			opts.JSONCompression,
			opts.JSONAllowDuplicate,
			opts.JSONStripOuterArray,
			opts.JSONStripNullValues,
			opts.JSONReplaceInvalidCharacters,
		},
		ExternalTableFileFormatTypeAvro: {
			opts.AvroCompression,
			opts.AvroReplaceInvalidCharacters,
		},
		ExternalTableFileFormatTypeORC: {
			opts.ORCTrimSpace,
			opts.ORCReplaceInvalidCharacters,
			opts.ORCNullIf,
		},
		ExternalTableFileFormatTypeParquet: {
			opts.ParquetCompression,
			opts.ParquetBinaryAsText,
			opts.ParquetReplaceInvalidCharacters,
		},
	}
}
