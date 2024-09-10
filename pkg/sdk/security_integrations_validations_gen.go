package sdk

var (
	_ validatable = new(CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions)
	_ validatable = new(CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions)
	_ validatable = new(CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions)
	_ validatable = new(CreateExternalOauthSecurityIntegrationOptions)
	_ validatable = new(CreateOauthForPartnerApplicationsSecurityIntegrationOptions)
	_ validatable = new(CreateOauthForCustomClientsSecurityIntegrationOptions)
	_ validatable = new(CreateSaml2SecurityIntegrationOptions)
	_ validatable = new(CreateScimSecurityIntegrationOptions)
	_ validatable = new(AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions)
	_ validatable = new(AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions)
	_ validatable = new(AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions)
	_ validatable = new(AlterExternalOauthSecurityIntegrationOptions)
	_ validatable = new(AlterOauthForPartnerApplicationsSecurityIntegrationOptions)
	_ validatable = new(AlterOauthForCustomClientsSecurityIntegrationOptions)
	_ validatable = new(AlterSaml2SecurityIntegrationOptions)
	_ validatable = new(AlterScimSecurityIntegrationOptions)
	_ validatable = new(DropSecurityIntegrationOptions)
	_ validatable = new(DescribeSecurityIntegrationOptions)
	_ validatable = new(ShowSecurityIntegrationOptions)
)

func (opts *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateExternalOauthSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if everyValueSet(opts.ExternalOauthBlockedRolesList, opts.ExternalOauthAllowedRolesList) {
		errs = append(errs, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList"))
	}
	if !exactlyOneValueSet(opts.ExternalOauthJwsKeysUrl, opts.ExternalOauthRsaPublicKey) {
		errs = append(errs, errExactlyOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey"))
	}
	if everyValueSet(opts.ExternalOauthJwsKeysUrl, opts.ExternalOauthRsaPublicKey2) {
		errs = append(errs, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateExternalOauthSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOauthForPartnerApplicationsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	if opts.OauthClient == OauthSecurityIntegrationClientLooker && opts.OauthRedirectUri == nil {
		errs = append(errs, NewError("OauthRedirectUri is required when OauthClient is LOOKER"))
	}
	return JoinErrors(errs...)
}

func (opts *CreateOauthForCustomClientsSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateOauthForCustomClientsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

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
	if opts.ScimClient == ScimSecurityIntegrationScimClientAzure && opts.SyncPassword != nil {
		errs = append(errs, NewError("SyncPassword is not supported for Azure scim client"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.OauthTokenEndpoint, opts.Set.OauthClientAuthMethod, opts.Set.OauthClientId, opts.Set.OauthClientSecret, opts.Set.OauthGrantClientCredentials, opts.Set.OauthAccessTokenValidity, opts.Set.OauthRefreshTokenValidity, opts.Set.OauthAllowedScopes, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions.Set", "Enabled", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantClientCredentials", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions.Unset", "Enabled", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.OauthAuthorizationEndpoint, opts.Set.OauthTokenEndpoint, opts.Set.OauthClientAuthMethod, opts.Set.OauthClientId, opts.Set.OauthClientSecret, opts.Set.OauthGrantAuthorizationCode, opts.Set.OauthAccessTokenValidity, opts.Set.OauthRefreshTokenValidity, opts.Set.OauthAllowedScopes, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions.Set", "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantAuthorizationCode", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "OauthAllowedScopes", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions.Unset", "Enabled", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.OauthAuthorizationEndpoint, opts.Set.OauthTokenEndpoint, opts.Set.OauthClientAuthMethod, opts.Set.OauthClientId, opts.Set.OauthClientSecret, opts.Set.OauthGrantJwtBearer, opts.Set.OauthAccessTokenValidity, opts.Set.OauthRefreshTokenValidity, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions.Set", "Enabled", "OauthAuthorizationEndpoint", "OauthTokenEndpoint", "OauthClientAuthMethod", "OauthClientId", "OauthClientSecret", "OauthGrantJwtBearer", "OauthAccessTokenValidity", "OauthRefreshTokenValidity", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions.Unset", "Enabled", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterExternalOauthSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterExternalOauthSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if everyValueSet(opts.Set.ExternalOauthBlockedRolesList, opts.Set.ExternalOauthAllowedRolesList) {
			errs = append(errs, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList"))
		}
		if everyValueSet(opts.Set.ExternalOauthJwsKeysUrl, opts.Set.ExternalOauthRsaPublicKey) {
			errs = append(errs, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey"))
		}
		if everyValueSet(opts.Set.ExternalOauthJwsKeysUrl, opts.Set.ExternalOauthRsaPublicKey2) {
			errs = append(errs, errOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "ExternalOauthJwsKeysUrl", "ExternalOauthRsaPublicKey2"))
		}
		if !anyValueSet(opts.Set.Enabled, opts.Set.ExternalOauthType, opts.Set.ExternalOauthIssuer, opts.Set.ExternalOauthTokenUserMappingClaim, opts.Set.ExternalOauthSnowflakeUserMappingAttribute, opts.Set.ExternalOauthJwsKeysUrl, opts.Set.ExternalOauthBlockedRolesList, opts.Set.ExternalOauthAllowedRolesList, opts.Set.ExternalOauthRsaPublicKey, opts.Set.ExternalOauthRsaPublicKey2, opts.Set.ExternalOauthAudienceList, opts.Set.ExternalOauthAnyRoleMode, opts.Set.ExternalOauthScopeDelimiter, opts.Set.ExternalOauthScopeMappingAttribute, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterExternalOauthSecurityIntegrationOptions.Set", "Enabled", "ExternalOauthType", "ExternalOauthIssuer", "ExternalOauthTokenUserMappingClaim", "ExternalOauthSnowflakeUserMappingAttribute", "ExternalOauthJwsKeysUrl", "ExternalOauthBlockedRolesList", "ExternalOauthAllowedRolesList", "ExternalOauthRsaPublicKey", "ExternalOauthRsaPublicKey2", "ExternalOauthAudienceList", "ExternalOauthAnyRoleMode", "ExternalOauthScopeDelimiter", "ExternalOauthScopeMappingAttribute", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.ExternalOauthAudienceList) {
			errs = append(errs, errAtLeastOneOf("AlterExternalOauthSecurityIntegrationOptions.Unset", "Enabled", "ExternalOauthAudienceList"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.OauthIssueRefreshTokens, opts.Set.OauthRedirectUri, opts.Set.OauthRefreshTokenValidity, opts.Set.OauthUseSecondaryRoles, opts.Set.BlockedRolesList, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Set", "Enabled", "OauthIssueRefreshTokens", "OauthRedirectUri", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "BlockedRolesList", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.OauthUseSecondaryRoles) {
			errs = append(errs, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Unset", "Enabled", "OauthUseSecondaryRoles"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterOauthForCustomClientsSecurityIntegrationOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.SetTags, opts.UnsetTags) {
		errs = append(errs, errExactlyOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Enabled, opts.Set.OauthRedirectUri, opts.Set.OauthAllowNonTlsRedirectUri, opts.Set.OauthEnforcePkce, opts.Set.PreAuthorizedRolesList, opts.Set.BlockedRolesList, opts.Set.OauthIssueRefreshTokens, opts.Set.OauthRefreshTokenValidity, opts.Set.OauthUseSecondaryRoles, opts.Set.NetworkPolicy, opts.Set.OauthClientRsaPublicKey, opts.Set.OauthClientRsaPublicKey2, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Set", "Enabled", "OauthRedirectUri", "OauthAllowNonTlsRedirectUri", "OauthEnforcePkce", "PreAuthorizedRolesList", "BlockedRolesList", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "NetworkPolicy", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.NetworkPolicy, opts.Unset.OauthUseSecondaryRoles, opts.Unset.OauthClientRsaPublicKey, opts.Unset.OauthClientRsaPublicKey2) {
			errs = append(errs, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "OauthUseSecondaryRoles", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2"))
		}
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
		if !anyValueSet(opts.Unset.Enabled, opts.Unset.NetworkPolicy, opts.Unset.SyncPassword) {
			errs = append(errs, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "SyncPassword"))
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
