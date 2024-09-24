package sdk

var (
	_ validatable = new(CreateWithOAuthClientCredentialsFlowSecretOptions)
	_ validatable = new(CreateWithOAuthAuthorizationCodeFlowSecretOptions)
	_ validatable = new(CreateWithBasicAuthenticationSecretOptions)
	_ validatable = new(CreateWithGenericStringSecretOptions)
	_ validatable = new(AlterSecretOptions)
	_ validatable = new(DropSecretOptions)
	_ validatable = new(ShowSecretOptions)
	_ validatable = new(DescribeSecretOptions)
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
		if moreThanOneValueSet(opts.Set.SetForOAuthClientCredentialsFlow, opts.Set.SetForOAuthAuthorizationFlow, opts.Set.SetForBasicAuthentication, opts.Set.SetForGenericString) {
			errs = append(errs, errMoreThanOneOf("AlterSecretOptions.Set", "SetForOAuthClientCredentialsFlow", "SetForOAuthAuthorizationFlow", "SetForBasicAuthentication", "SetForGenericString"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeSecretOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
