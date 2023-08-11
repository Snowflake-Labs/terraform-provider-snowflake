package sdk

import "errors"

func (opts *CreateRoleOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterRoleOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueNil(opts.RenameTo, opts.Set, opts.Unset) {
		errs = append(errs, errors.New("no alter action specified"))
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set, opts.Unset) {
		errs = append(errs, errOneOf("RenameTo", "Set", "Unset"))
	}
	return errors.Join(errs...)
}

func (opts *DropRoleOptions) validateProp() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *GrantRoleOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Grant.Role, opts.Grant.User) {
		errs = append(errs, errors.New("only one grant option can be set [TO ROLE or TO USER]"))
	}
	return errors.Join(errs...)
}

func (opts *RevokeRoleOptions) validateProp() error {
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Revoke.Role, opts.Revoke.User) {
		errs = append(errs, errors.New("only one revoke option can be set [FROM ROLE or FROM USER]"))
	}
	return errors.Join(errs...)
}
