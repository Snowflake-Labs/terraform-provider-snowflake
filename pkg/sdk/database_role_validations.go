package sdk

import "errors"

var (
	_ validatableOpts = new(CreateDatabaseRoleOptions)
	_ validatableOpts = new(AlterDatabaseRoleOptions)
	_ validatableOpts = new(DropDatabaseRoleOptions)
	_ validatableOpts = new(ShowDatabaseRoleOptions)
)

var errDifferentDatabase = errors.New("database must be the same")

func (opts *CreateDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
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

func (opts *AlterDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(
		opts.Rename,
		opts.Set,
		opts.Unset,
	); !ok {
		errs = append(errs, errAlterNeedsExactlyOneAction)
	}
	if rename := opts.Rename; valueSet(rename) {
		if !validObjectidentifier(rename.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if opts.name.DatabaseName() != rename.Name.DatabaseName() {
			errs = append(errs, errDifferentDatabase)
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !unset.Comment {
			errs = append(errs, errAlterNeedsAtLeastOneProperty)
		}
	}
	return errors.Join(errs...)
}

func (opts *DropDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *ShowDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	if !validObjectidentifier(opts.Database) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, errPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}
