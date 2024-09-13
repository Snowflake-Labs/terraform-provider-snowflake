package sdk

var (
	_ validatable = new(CreateWithOAuthClientCredentialsFlowSecretOptions)
	_ validatable = new(CreateWithOAuthAuthorizationCodeFlowSecretOptions)
	_ validatable = new(CreateWithBasicAuthenticationSecretOptions)
	_ validatable = new(CreateWithGenericStringSecretOptions)
	_ validatable = new(AlterSecretOptions)
)

func (opts *CreateWithOAuthClientCredentialsFlowSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithOAuthClientCredentialsFlowSecretOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateWithOAuthAuthorizationCodeFlowSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithOAuthAuthorizationCodeFlowSecretOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateWithBasicAuthenticationSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithBasicAuthenticationSecretOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateWithGenericStringSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateWithGenericStringSecretOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !exactlyOneValueSet(opts.Set, opts.Unset) {
		errs = append(errs, errExactlyOneOf("AlterSecretOptions", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if !exactlyOneValueSet(opts.Set.SetForOAuthClientCredentialsFlow, opts.Set.SetForOAuthAuthorizationFlow, opts.Set.SetForBasicAuthentication, opts.Set.SetForGenericString) {
			errs = append(errs, errExactlyOneOf("AlterSecretOptions.Set", "SetForOAuthClientCredentialsFlow", "SetForOAuthAuthorizationFlow", "SetForBasicAuthentication", "SetForGenericString"))
		}
	}
	return JoinErrors(errs...)
}
