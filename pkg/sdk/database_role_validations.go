package sdk

import "errors"

var (
	_ validatableOpts = new(CreateDatabaseRoleOptions)
	_ validatableOpts = new(AlterDatabaseRoleOptions)
	_ validatableOpts = new(DropDatabaseRoleOptions)
	_ validatableOpts = new(ShowDatabaseRoleOptions)
)

var (
	errDifferentDatabase = errors.New("database must be the same")
)

func (opts *CreateDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *AlterDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if ok := exactlyOneValueSet(
		opts.Rename,
		opts.Set,
		opts.Unset,
	); !ok {
		return errAlterNeedsExactlyOneAction
	}
	if rename := opts.Rename; valueSet(rename) {
		if !validObjectidentifier(rename.Name) {
			return ErrInvalidObjectIdentifier
		}
		if opts.name.DatabaseName() != rename.Name.DatabaseName() {
			return errDifferentDatabase
		}
	}
	if unset := opts.Unset; valueSet(unset) {
		if !unset.Comment {
			return errAlterNeedsAtLeastOneProperty
		}
	}
	return nil
}

func (opts *DropDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (opts *ShowDatabaseRoleOptions) validateProp() error {
	if opts == nil {
		return errNilOptions
	}
	if !validObjectidentifier(opts.Database) {
		return ErrInvalidObjectIdentifier
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return errPatternRequiredForLikeKeyword
	}
	return nil
}
