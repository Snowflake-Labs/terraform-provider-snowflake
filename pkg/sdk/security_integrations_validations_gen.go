package sdk

var (
	_ validatable = new(CreateSAML2SecurityIntegrationOptions)
	_ validatable = new(CreateSCIMSecurityIntegrationOptions)
	_ validatable = new(AlterSAML2IntegrationSecurityIntegrationOptions)
	_ validatable = new(AlterSCIMIntegrationSecurityIntegrationOptions)
	_ validatable = new(DropSecurityIntegrationOptions)
	_ validatable = new(DescribeSecurityIntegrationOptions)
	_ validatable = new(ShowSecurityIntegrationOptions)
)

func (opts *CreateSAML2SecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSAML2SecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

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

func (opts *AlterSAML2IntegrationSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.Saml2Issuer, opts.Set.Saml2SsoUrl, opts.Set.Saml2Provider, opts.Set.Saml2X509Cert, opts.Set.AllowedUserDomains, opts.Set.AllowedEmailPatterns, opts.Set.Saml2SpInitiatedLoginPageLabel, opts.Set.Saml2EnableSpInitiated, opts.Set.Saml2SnowflakeX509Cert, opts.Set.Saml2SignRequest, opts.Set.Saml2RequestedNameidFormat, opts.Set.Saml2PostLogoutRedirectUrl, opts.Set.Saml2ForceAuthn, opts.Set.Saml2SnowflakeIssuerUrl, opts.Set.Saml2SnowflakeAcsUrl, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterSAML2IntegrationSecurityIntegrationOptions.Set", "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider", "Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns", "Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.Saml2ForceAuthn) {
			errs = append(errs, errAtLeastOneOf("AlterSAML2IntegrationSecurityIntegrationOptions.Unset", "Enabled", "Saml2ForceAuthn"))
		}
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
