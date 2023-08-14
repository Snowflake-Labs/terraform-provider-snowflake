package sdk

import "errors"

func (opts *CreateRoleOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
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
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueNil(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errors.New("no alter action specified"))
	}
	if anyValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) &&
		!exactlyOneValueSet(opts.RenameTo, opts.SetComment, opts.UnsetComment, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errOneOf("RenameTo", "SetComment", "UnsetComment", "SetTags", "UnsetTags"))
	}
	return errors.Join(errs...)
}

func (opts *DropRoleOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *ShowRoleOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.InClass) && !validObjectidentifier(opts.InClass.Class) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *GrantRoleOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Grant.Role, opts.Grant.User) {
		errs = append(errs, errors.New("only one grant option can be set [TO ROLE or TO USER]"))
	}
	if valueSet(opts.Grant.Role) && !validObjectidentifier(opts.Grant.Role) {
		errs = append(errs, errors.New("invalid object identifier for granted role"))
	}
	if valueSet(opts.Grant.User) && !validObjectidentifier(opts.Grant.User) {
		errs = append(errs, errors.New("invalid object identifier for granted user"))
	}
	return errors.Join(errs...)
}

func (opts *RevokeRoleOptions) validateProp() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Revoke.Role, opts.Revoke.User) {
		errs = append(errs, errors.New("only one revoke option can be set [FROM ROLE or FROM USER]"))
	}
	return errors.Join(errs...)
}
