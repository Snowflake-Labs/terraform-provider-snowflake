package sdk

var (
	_ validatable = new(GrantApplicationRoleOptions)
	_ validatable = new(RevokeApplicationRoleOptions)
	_ validatable = new(ShowApplicationRoleOptions)
)

func (opts *GrantApplicationRoleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.To) {
		if !exactlyOneValueSet(opts.To.RoleName, opts.To.ApplicationRoleName, opts.To.ApplicationName) {
			errs = append(errs, errExactlyOneOf("GrantApplicationRoleOptions.To", "RoleName", "ApplicationRoleName", "ApplicationName"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *RevokeApplicationRoleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.From) {
		if !exactlyOneValueSet(opts.From.RoleName, opts.From.ApplicationRoleName, opts.From.ApplicationName) {
			errs = append(errs, errExactlyOneOf("RevokeApplicationRoleOptions.From", "RoleName", "ApplicationRoleName", "ApplicationName"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowApplicationRoleOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.ApplicationName) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
