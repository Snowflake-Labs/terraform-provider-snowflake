package sdk

import "errors"

var (
	_ validatable = (*CreateRoleOptions)(nil)
	_ validatable = (*AlterRoleOptions)(nil)
	_ validatable = (*DropRoleOptions)(nil)
	_ validatable = (*ShowRoleOptions)(nil)
	_ validatable = (*GrantRoleOptions)(nil)
	_ validatable = (*RevokeRoleOptions)(nil)
)

func (opts *CreateRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *AlterRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
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

func (opts *DropRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return errInvalidObjectIdentifier
	}
	return nil
}

func (opts *ShowRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, errPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.InClass) && !validObjectidentifier(opts.InClass.Class) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *GrantRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if (opts.Grant.Role != nil && opts.Grant.User != nil) || (opts.Grant.Role == nil && opts.Grant.User == nil) {
		errs = append(errs, errors.New("only one grant option can be set [TO ROLE or TO USER]"))
	}
	if opts.Grant.Role != nil && !validObjectidentifier(opts.Grant.Role) {
		errs = append(errs, errors.New("invalid object identifier for granted role"))
	}
	if opts.Grant.User != nil && !validObjectidentifier(opts.Grant.User) {
		errs = append(errs, errors.New("invalid object identifier for granted user"))
	}
	return errors.Join(errs...)
}

func (opts *RevokeRoleOptions) validate() error {
	if opts == nil {
		return errNilOptions
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, errInvalidObjectIdentifier)
	}
	if (opts.Revoke.Role != nil && opts.Revoke.User != nil) || (opts.Revoke.Role == nil && opts.Revoke.User == nil) {
		errs = append(errs, errors.New("only one revoke option can be set [FROM ROLE or FROM USER]"))
	}
	return errors.Join(errs...)
}
