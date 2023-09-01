package sdk

import "errors"

var (
	_ validatable = new(createDatabaseRoleOptions)
	_ validatable = new(alterDatabaseRoleOptions)
	_ validatable = new(dropDatabaseRoleOptions)
	_ validatable = new(showDatabaseRoleOptions)
	_ validatable = new(grantDatabaseRoleOptions)
	_ validatable = new(revokeDatabaseRoleOptions)
	_ validatable = new(grantDatabaseRoleToShareOptions)
	_ validatable = new(revokeDatabaseRoleFromShareOptions)
)

func (opts *createDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, errOneOf("createDatabaseRoleOptions", "OrReplace", "IfNotExists"))
	}
	return errors.Join(errs...)
}

func (opts *createDatabaseRoleOptions) validate2() error {
	if opts == nil {
		return NewError("Nil options")
	}
	var errs []error
	if !validObjectidentifier(opts.name) {
		errs = append(errs, NewError("Invalid object identifier"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) && *opts.OrReplace && *opts.IfNotExists {
		errs = append(errs, NewError("One of IfNotExists / OrReplace"))
	}
	return errors.Join(errs...)
}

func (opts *createDatabaseRoleOptions) validate3() error {
	return validateAll(
		opts,
		validObjectIdentifier(opts.name),
		oneOf(opts.OrReplace, opts.IfNotExists),
	)
}

func validObjectIdentifier(objectIdentifier ObjectIdentifier) error {
	if !validObjectidentifier(objectIdentifier) {
		return NewError("Invalid object identifier")
	}
	return nil
}

func oneOf(fields ...any) error {
	return nil
}

func validateAll(opts any, validations ...error) error {
	if opts == nil {
		return NewError("Nil options")
	}
	return errors.Join(validations...)
}

func (opts *alterDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Rename, opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("alterDatabaseRoleOptions", "Rename", "Set", "Unset"))
	}
	if rename := opts.Rename; valueSet(rename) {
		if !ValidObjectIdentifier(rename.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
		if opts.name.DatabaseName() != rename.Name.DatabaseName() {
			errs = append(errs, ErrDifferentDatabase)
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !unset.Comment {
			errs = append(errs, errAtLeastOneOf("alterDatabaseRoleOptions.Unset", "Comment"))
		}
	}
	return errors.Join(errs...)
}

func (opts *dropDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *showDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.Database) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}

func (opts *grantDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.ParentRole.DatabaseRoleName, opts.ParentRole.AccountRoleName); !ok {
		errs = append(errs, errOneOf("DatabaseRoleName", "AccountRoleName"))
	}
	return errors.Join(errs...)
}

func (opts *revokeDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if ok := exactlyOneValueSet(opts.ParentRole.DatabaseRoleName, opts.ParentRole.AccountRoleName); !ok {
		errs = append(errs, errOneOf("DatabaseRoleName", "AccountRoleName"))
	}
	return errors.Join(errs...)
}

func (opts *grantDatabaseRoleToShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.Share) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (opts *revokeDatabaseRoleFromShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.Share) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}
