package sdk

var (
	_ validatable = new(CreateSaml2SecurityIntegrationOptions)
	_ validatable = new(CreateScimSecurityIntegrationOptions)
	_ validatable = new(AlterSaml2SecurityIntegrationOptions)
	_ validatable = new(AlterScimSecurityIntegrationOptions)
	_ validatable = new(DropSecurityIntegrationOptions)
	_ validatable = new(DescribeSecurityIntegrationOptions)
	_ validatable = new(ShowSecurityIntegrationOptions)
)

func (opts *CreateSaml2SecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSaml2SecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateScimSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateScimSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterSaml2SecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.RefreshSaml2SnowflakePrivateKey, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterSaml2SecurityIntegrationOptions", "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.Saml2Issuer, opts.Set.Saml2SsoUrl, opts.Set.Saml2Provider, opts.Set.Saml2X509Cert, opts.Set.AllowedUserDomains, opts.Set.AllowedEmailPatterns, opts.Set.Saml2SpInitiatedLoginPageLabel, opts.Set.Saml2EnableSpInitiated, opts.Set.Saml2SnowflakeX509Cert, opts.Set.Saml2SignRequest, opts.Set.Saml2RequestedNameidFormat, opts.Set.Saml2PostLogoutRedirectUrl, opts.Set.Saml2ForceAuthn, opts.Set.Saml2SnowflakeIssuerUrl, opts.Set.Saml2SnowflakeAcsUrl, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Set", "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider", "Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns", "Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Saml2ForceAuthn, opts.Unset.Saml2RequestedNameidFormat, opts.Unset.Saml2PostLogoutRedirectUrl, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Unset", "Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterScimSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterScimSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.NetworkPolicy, opts.Set.SyncPassword, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.NetworkPolicy, opts.Unset.SyncPassword, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
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
