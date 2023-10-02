package sdk

import (
	"errors"
	"fmt"
)

var (
	_ validatable = (*CreateExternalTableOptions)(nil)
	_ validatable = (*CreateWithManualPartitioningExternalTableOptions)(nil)
	_ validatable = (*CreateDeltaLakeExternalTableOptions)(nil)
	_ validatable = (*CreateExternalTableUsingTemplateOptions)(nil)
	_ validatable = (*AlterExternalTableOptions)(nil)
	_ validatable = (*AlterExternalTablePartitionOptions)(nil)
	_ validatable = (*DropExternalTableOptions)(nil)
	_ validatable = (*ShowExternalTableOptions)(nil)
	_ validatable = (*describeExternalTableColumnsOptions)(nil)
	_ validatable = (*describeExternalTableStageOptions)(nil)
)

func (opts *CreateExternalTableOptions) validate() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateExternalTableOptions", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateExternalTableOptions", "FileFormat"))
	} else {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateWithManualPartitioningExternalTableOptions) validate() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithManualPartitioningExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateWithManualPartitioningExternalTableOptions", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateWithManualPartitioningExternalTableOptions", "FileFormat"))
	} else {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateWithManualPartitioningExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateWithManualPartitioningExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateDeltaLakeExternalTableOptions) validate() error {
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDeltaLakeExternalTableOptions", "OrReplace", "IfNotExists"))
	}
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOptions", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateDeltaLakeExternalTableOptions", "FileFormat"))
	} else {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateDeltaLakeExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateDeltaLakeExternalTableOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *CreateExternalTableUsingTemplateOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !valueSet(opts.Query) {
		errs = append(errs, errNotSet("CreateExternalTableUsingTemplateOptions", "Query"))
	}
	if !valueSet(opts.Location) {
		errs = append(errs, errNotSet("CreateExternalTableUsingTemplateOptions", "Location"))
	}
	if !valueSet(opts.FileFormat) {
		errs = append(errs, errNotSet("CreateExternalTableUsingTemplateOptions", "FileFormat"))
	} else {
		for i, ff := range opts.FileFormat {
			if !valueSet(ff.Name) && !valueSet(ff.Type) {
				errs = append(errs, errNotSet(fmt.Sprintf("CreateExternalTableUsingTemplateOptions.FileFormat[%d]", i), "Name or Type"))
			}
			if valueSet(ff.Name) && valueSet(ff.Type) {
				errs = append(errs, errOneOf(fmt.Sprintf("CreateExternalTableUsingTemplateOptions.FileFormat[%d]", i), "Name or Type"))
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *AlterExternalTableOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
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
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.AddPartitions, opts.DropPartition) {
		errs = append(errs, errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"))
	}
	return errors.Join(errs...)
}

func (opts *DropExternalTableOptions) validate() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
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

func (v *describeExternalTableColumnsOptions) validate() error {
	if !validObjectidentifier(v.name) {
		return errInvalidObjectIdentifier
	}
	return nil
}

func (v *describeExternalTableStageOptions) validate() error {
	if !validObjectidentifier(v.name) {
		return errInvalidObjectIdentifier
	}
	return nil
}

func (cpp *CloudProviderParams) validate() error {
	if anyValueSet(cpp.GoogleCloudStorageIntegration, cpp.MicrosoftAzureIntegration) && exactlyOneValueSet(cpp.GoogleCloudStorageIntegration, cpp.MicrosoftAzureIntegration) {
		return errOneOf("CloudProviderParams", "GoogleCloudStorageIntegration", "MicrosoftAzureIntegration")
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
