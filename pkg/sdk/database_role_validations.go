package sdk

var (
	_ validatableOpts = &CreateDatabaseRoleOptions{}
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
