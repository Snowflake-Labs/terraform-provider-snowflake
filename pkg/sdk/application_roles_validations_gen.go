package sdk

import "errors"

var (
	_ validatable = new(CreateApplicationRoleOptions)
	_ validatable = new(AlterApplicationRoleOptions)
	_ validatable = new(DropApplicationRoleOptions)
	_ validatable = new(ShowApplicationRoleOptions)
	_ validatable = new(GrantApplicationRoleOptions)
	_ validatable = new(RevokeApplicationRoleOptions)
)

func (opts *CreateApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateApplicationRoleOptions", "OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment); !ok {
		errs = append(errs, errExactlyOneOf("RenameTo", "SetComment", "UnsetComment"))
	}
	if valueSet(opts.RenameTo) && !validObjectidentifier(opts.RenameTo) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *DropApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.ApplicationName) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *GrantApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if valueSet(opts.GrantTo) {
		if ok := exactlyOneValueSet(opts.GrantTo.ParentRole, opts.GrantTo.ApplicationRole, opts.GrantTo.Application); !ok {
			errs = append(errs, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
		}
		if valueSet(opts.GrantTo.ParentRole) && !validObjectidentifier(opts.GrantTo.ParentRole) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if valueSet(opts.GrantTo.ApplicationRole) && !validObjectidentifier(opts.GrantTo.ApplicationRole) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if valueSet(opts.GrantTo.Application) && !validObjectidentifier(opts.GrantTo.Application) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (opts *RevokeApplicationRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if valueSet(opts.RevokeFrom) {
		if ok := exactlyOneValueSet(opts.RevokeFrom.ParentRole, opts.RevokeFrom.ApplicationRole, opts.RevokeFrom.Application); !ok {
			errs = append(errs, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
		}
		if valueSet(opts.RevokeFrom.ParentRole) && !validObjectidentifier(opts.RevokeFrom.ParentRole) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if valueSet(opts.RevokeFrom.ApplicationRole) && !validObjectidentifier(opts.RevokeFrom.ApplicationRole) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
		if valueSet(opts.RevokeFrom.Application) && !validObjectidentifier(opts.RevokeFrom.Application) {
			errs = append(errs, errInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}
