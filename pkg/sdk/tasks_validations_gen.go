package sdk

import "errors"

var (
	_ validatable = new(CreateTaskOptions)
)

func (opts *CreateTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateTaskOptions", "OrReplace", "IfNotExists"))
	}
	if valueSet(opts.Warehouse) {
		if !validObjectidentifier(opts.Warehouse.Warehouse) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if ok := exactlyOneValueSet(opts.Warehouse.Warehouse, opts.Warehouse.UserTaskManagedInitialWarehouseSize); !ok {
			errs = append(errs, errExactlyOneOf("Warehouse", "UserTaskManagedInitialWarehouseSize"))
		}
	}
	return errors.Join(errs...)
}
