package sdk

var (
	_ validatable = new(CreateSCIMSecurityIntegrationOptions)
	_ validatable = new(AlterSCIMIntegrationSecurityIntegrationOptions)
	_ validatable = new(DropSecurityIntegrationOptions)
	_ validatable = new(DescribeSecurityIntegrationOptions)
	_ validatable = new(ShowSecurityIntegrationOptions)
)

func (opts *CreateSCIMSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSCIMSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterSCIMIntegrationSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.NetworkPolicy, opts.Set.SyncPassword, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.NetworkPolicy, opts.Unset.SyncPassword, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Unset", "NetworkPolicy", "SyncPassword", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *DescribeSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}
