package sdk

var (
	_ validatable = new(CreateAuthenticationPolicyOptions)
	_ validatable = new(AlterAuthenticationPolicyOptions)
	_ validatable = new(DropAuthenticationPolicyOptions)
	_ validatable = new(ShowAuthenticationPolicyOptions)
	_ validatable = new(DescribeAuthenticationPolicyOptions)
)

func (opts *CreateAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateAuthenticationPolicyOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.RenameTo) {
		errs = append(errs, errExactlyOneOf("AlterAuthenticationPolicyOptions", "Set", "Unset", "RenameTo"))
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.AuthenticationMethods, opts.Set.MfaAuthenticationMethods, opts.Set.MfaEnrollment, opts.Set.ClientTypes, opts.Set.SecurityIntegrations, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set", "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.ClientTypes, opts.Unset.AuthenticationMethods, opts.Unset.Comment, opts.Unset.SecurityIntegrations, opts.Unset.MfaAuthenticationMethods, opts.Unset.MfaEnrollment) {
			errs = append(errs, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Unset", "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeAuthenticationPolicyOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
