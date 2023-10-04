package sdk

import "errors"

var (
	_ validatable = new(CreateTaskOptions)
	_ validatable = new(AlterTaskOptions)
	_ validatable = new(DropTaskOptions)
	_ validatable = new(ShowTaskOptions)
	_ validatable = new(DescribeTaskOptions)
	_ validatable = new(ExecuteTaskOptions)
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
		if ok := exactlyOneValueSet(opts.Warehouse.Warehouse, opts.Warehouse.UserTaskManagedInitialWarehouseSize); !ok {
			errs = append(errs, errExactlyOneOf("Warehouse", "UserTaskManagedInitialWarehouseSize"))
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *AlterTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Resume, opts.Suspend, opts.RemoveAfter, opts.AddAfter, opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags, opts.ModifyAs, opts.ModifyWhen); !ok {
		errs = append(errs, errExactlyOneOf("Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "ModifyAs", "ModifyWhen"))
	}
	if valueSet(opts.Set) {
		if ok := anyValueSet(opts.Set.Warehouse, opts.Set.Schedule, opts.Set.Config, opts.Set.AllowOverlappingExecution, opts.Set.UserTaskTimeoutMs, opts.Set.SuspendTaskAfterNumFailures, opts.Set.Comment, opts.Set.SessionParameters); !ok {
			errs = append(errs, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParameters"))
		}
		if valueSet(opts.Set.Warehouse) && !validObjectidentifier(opts.Set.Warehouse) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if valueSet(opts.Set.SessionParameters) {
			if err := opts.Set.SessionParameters.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if valueSet(opts.Unset) {
		if ok := anyValueSet(opts.Unset.Warehouse, opts.Unset.Schedule, opts.Unset.Config, opts.Unset.AllowOverlappingExecution, opts.Unset.UserTaskTimeoutMs, opts.Unset.SuspendTaskAfterNumFailures, opts.Unset.Comment, opts.Unset.SessionParametersUnset); !ok {
			errs = append(errs, errAtLeastOneOf("Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "Comment", "SessionParametersUnset"))
		}
		if valueSet(opts.Unset.SessionParametersUnset) {
			if err := opts.Unset.SessionParametersUnset.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (opts *DropTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	return errors.Join(errs...)
}

func (opts *DescribeTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ExecuteTaskOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
