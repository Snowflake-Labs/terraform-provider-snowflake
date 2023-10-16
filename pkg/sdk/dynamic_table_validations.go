package sdk

import (
	"errors"
)

var (
	_ validatable = new(createDynamicTableOptions)
	_ validatable = new(alterDynamicTableOptions)
	_ validatable = new(dropDynamicTableOptions)
	_ validatable = new(showDynamicTableOptions)
	_ validatable = new(describeDynamicTableOptions)
	_ validatable = new(DynamicTableSet)
)

func (tl *TargetLag) validate() error {
	if tl == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if everyValueSet(tl.Lagtime, tl.Downstream) {
		errs = append(errs, errOneOf("TargetLag", "Lagtime", "Downstream"))
	}
	return errors.Join(errs...)
}

func (opts *createDynamicTableOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if !validObjectidentifier(opts.warehouse) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (dts *DynamicTableSet) validate() error {
	var errs []error
	if valueSet(dts.TargetLag) {
		errs = append(errs, dts.TargetLag.validate())
	}

	if valueSet(dts.Warehouse) {
		if !validObjectidentifier(*dts.Warehouse) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (opts *alterDynamicTableOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.Suspend, opts.Resume, opts.Refresh, opts.Set); !ok {
		errs = append(errs, errExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set"))
	}
	if !anyValueSet(opts.Suspend, opts.Resume, opts.Refresh, opts.Set) {
		errs = append(errs, errAtLeastOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set"))
	}
	if valueSet(opts.Set) && valueSet(opts.Set.TargetLag) {
		errs = append(errs, opts.Set.TargetLag.validate())
	}
	return errors.Join(errs...)
}

func (opts *showDynamicTableOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, errPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("showDynamicTableOptions.In", "Account", "Database", "Schema"))
	}
	return errors.Join(errs...)
}

func (opts *dropDynamicTableOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error

	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *describeDynamicTableOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
