package sdk

import "fmt"

var (
	_ validatable = new(CreateExternalVolumeOptions)
	_ validatable = new(AlterExternalVolumeOptions)
	_ validatable = new(DropExternalVolumeOptions)
	_ validatable = new(DescribeExternalVolumeOptions)
	_ validatable = new(ShowExternalVolumeOptions)
)

func (opts *CreateExternalVolumeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateExternalVolumeOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}

	// Custom (not code generated) validations

	// Apply errExactlyOneOf to each element in storage locations list
	for i, storageLocation := range opts.StorageLocations {
		if !exactlyOneValueSet(storageLocation.S3StorageLocationParams, storageLocation.GCSStorageLocationParams, storageLocation.AzureStorageLocationParams) {
			errs = append(errs, errExactlyOneOf(fmt.Sprintf("CreateExternalVolumeOptions.StorageLocation[%d]", i), "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
		}
	}

	// Check the storage location list is not empty, as at least 1 storage location is required for an external volume
	if len(opts.StorageLocations) == 0 {
		errs = append(errs, errNotSet("CreateExternalVolumeOptions", "StorageLocations"))
	}

	return JoinErrors(errs...)
}

func (opts *AlterExternalVolumeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !exactlyOneValueSet(opts.RemoveStorageLocation, opts.Set, opts.AddStorageLocation) {
		errs = append(errs, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.AddStorageLocation) {
		if !exactlyOneValueSet(opts.AddStorageLocation.S3StorageLocationParams, opts.AddStorageLocation.GCSStorageLocationParams, opts.AddStorageLocation.AzureStorageLocationParams) {
			errs = append(errs, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropExternalVolumeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeExternalVolumeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowExternalVolumeOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
