package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
)

var _ SecurityIntegrations = (*securityIntegrations)(nil)

type securityIntegrations struct {
	client *Client
}

func (v *securityIntegrations) CreateSnowflakeOauthPartner(ctx context.Context, request *CreateSnowflakeOauthPartnerSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateSnowflakeOauthCustom(ctx context.Context, request *CreateSnowflakeOauthCustomSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateSaml2(ctx context.Context, request *CreateSaml2SecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateScim(ctx context.Context, request *CreateScimSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterSnowflakeOauthPartner(ctx context.Context, request *AlterSnowflakeOauthPartnerSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterSnowflakeOauthCustom(ctx context.Context, request *AlterSnowflakeOauthCustomSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterSaml2(ctx context.Context, request *AlterSaml2SecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterScim(ctx context.Context, request *AlterScimSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error) {
	opts := &DescribeSecurityIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[securityIntegrationDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[securityIntegrationDescRow, SecurityIntegrationProperty](rows), nil
}

func (v *securityIntegrations) Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[securityIntegrationShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[securityIntegrationShowRow, SecurityIntegration](dbRows)
	return resultList, nil
}

func (v *securityIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error) {
	securityIntegrations, err := v.Show(ctx, NewShowSecurityIntegrationRequest().WithLike(&Like{
		Pattern: String(id.Name()),
	}))
	if err != nil {
		return nil, err
	}
	return collections.FindOne(securityIntegrations, func(r SecurityIntegration) bool { return r.Name == id.Name() })
}

func (r *CreateSnowflakeOauthPartnerSecurityIntegrationRequest) toOpts() *CreateSnowflakeOauthPartnerSecurityIntegrationOptions {
	opts := &CreateSnowflakeOauthPartnerSecurityIntegrationOptions{
		OrReplace:                 r.OrReplace,
		IfNotExists:               r.IfNotExists,
		name:                      r.name,
		OauthClient:               r.OauthClient,
		OauthRedirectUri:          r.OauthRedirectUri,
		Enabled:                   r.Enabled,
		OauthIssueRefreshTokens:   r.OauthIssueRefreshTokens,
		OauthRefreshTokenValidity: r.OauthRefreshTokenValidity,
		OauthUseSecondaryRoles:    r.OauthUseSecondaryRoles,

		Comment: r.Comment,
	}
	if r.BlockedRolesList != nil {
		opts.BlockedRolesList = &BlockedRolesList{
			BlockedRolesList: r.BlockedRolesList.BlockedRolesList,
		}
	}
	return opts
}

func (r *CreateSnowflakeOauthCustomSecurityIntegrationRequest) toOpts() *CreateSnowflakeOauthCustomSecurityIntegrationOptions {
	opts := &CreateSnowflakeOauthCustomSecurityIntegrationOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		OauthClientType:             r.OauthClientType,
		OauthRedirectUri:            r.OauthRedirectUri,
		Enabled:                     r.Enabled,
		OauthAllowNonTlsRedirectUri: r.OauthAllowNonTlsRedirectUri,
		OauthEnforcePkce:            r.OauthEnforcePkce,
		OauthUseSecondaryRoles:      r.OauthUseSecondaryRoles,

		OauthIssueRefreshTokens:   r.OauthIssueRefreshTokens,
		OauthRefreshTokenValidity: r.OauthRefreshTokenValidity,
		NetworkPolicy:             r.NetworkPolicy,
		OauthClientRsaPublicKey:   r.OauthClientRsaPublicKey,
		OauthClientRsaPublicKey2:  r.OauthClientRsaPublicKey2,
		Comment:                   r.Comment,
	}
	if r.PreAuthorizedRolesList != nil {
		opts.PreAuthorizedRolesList = &PreAuthorizedRolesList{
			PreAuthorizedRolesList: r.PreAuthorizedRolesList.PreAuthorizedRolesList,
		}
	}
	if r.BlockedRolesList != nil {
		opts.BlockedRolesList = &BlockedRolesList{
			BlockedRolesList: r.BlockedRolesList.BlockedRolesList,
		}
	}
	return opts
}

func (r *CreateSaml2SecurityIntegrationRequest) toOpts() *CreateSaml2SecurityIntegrationOptions {
	opts := &CreateSaml2SecurityIntegrationOptions{
		OrReplace:                      r.OrReplace,
		IfNotExists:                    r.IfNotExists,
		name:                           r.name,
		Enabled:                        r.Enabled,
		Saml2Issuer:                    r.Saml2Issuer,
		Saml2SsoUrl:                    r.Saml2SsoUrl,
		Saml2Provider:                  r.Saml2Provider,
		Saml2X509Cert:                  r.Saml2X509Cert,
		AllowedUserDomains:             r.AllowedUserDomains,
		AllowedEmailPatterns:           r.AllowedEmailPatterns,
		Saml2SpInitiatedLoginPageLabel: r.Saml2SpInitiatedLoginPageLabel,
		Saml2EnableSpInitiated:         r.Saml2EnableSpInitiated,
		Saml2SnowflakeX509Cert:         r.Saml2SnowflakeX509Cert,
		Saml2SignRequest:               r.Saml2SignRequest,
		Saml2RequestedNameidFormat:     r.Saml2RequestedNameidFormat,
		Saml2PostLogoutRedirectUrl:     r.Saml2PostLogoutRedirectUrl,
		Saml2ForceAuthn:                r.Saml2ForceAuthn,
		Saml2SnowflakeIssuerUrl:        r.Saml2SnowflakeIssuerUrl,
		Saml2SnowflakeAcsUrl:           r.Saml2SnowflakeAcsUrl,
		Comment:                        r.Comment,
	}
	return opts
}

func (r *CreateScimSecurityIntegrationRequest) toOpts() *CreateScimSecurityIntegrationOptions {
	opts := &CreateScimSecurityIntegrationOptions{
		OrReplace:     r.OrReplace,
		IfNotExists:   r.IfNotExists,
		name:          r.name,
		Enabled:       r.Enabled,
		ScimClient:    r.ScimClient,
		RunAsRole:     r.RunAsRole,
		NetworkPolicy: r.NetworkPolicy,
		SyncPassword:  r.SyncPassword,
		Comment:       r.Comment,
	}
	return opts
}

func (r *AlterSnowflakeOauthPartnerSecurityIntegrationRequest) toOpts() *AlterSnowflakeOauthPartnerSecurityIntegrationOptions {
	opts := &AlterSnowflakeOauthPartnerSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &SnowflakeOauthPartnerIntegrationSet{
			Enabled:                   r.Set.Enabled,
			OauthRedirectUri:          r.Set.OauthRedirectUri,
			OauthIssueRefreshTokens:   r.Set.OauthIssueRefreshTokens,
			OauthRefreshTokenValidity: r.Set.OauthRefreshTokenValidity,
			OauthUseSecondaryRoles:    r.Set.OauthUseSecondaryRoles,

			Comment: r.Set.Comment,
		}
		if r.Set.BlockedRolesList != nil {
			opts.Set.BlockedRolesList = &BlockedRolesList{
				BlockedRolesList: r.Set.BlockedRolesList.BlockedRolesList,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &SnowflakeOauthPartnerIntegrationUnset{
			Enabled:                r.Unset.Enabled,
			OauthUseSecondaryRoles: r.Unset.OauthUseSecondaryRoles,
		}
	}
	return opts
}

func (r *AlterSnowflakeOauthCustomSecurityIntegrationRequest) toOpts() *AlterSnowflakeOauthCustomSecurityIntegrationOptions {
	opts := &AlterSnowflakeOauthCustomSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &SnowflakeOauthCustomIntegrationSet{
			Enabled:                     r.Set.Enabled,
			OauthRedirectUri:            r.Set.OauthRedirectUri,
			OauthAllowNonTlsRedirectUri: r.Set.OauthAllowNonTlsRedirectUri,
			OauthEnforcePkce:            r.Set.OauthEnforcePkce,
			OauthUseSecondaryRoles:      r.Set.OauthUseSecondaryRoles,

			OauthIssueRefreshTokens:   r.Set.OauthIssueRefreshTokens,
			OauthRefreshTokenValidity: r.Set.OauthRefreshTokenValidity,
			NetworkPolicy:             r.Set.NetworkPolicy,
			OauthClientRsaPublicKey:   r.Set.OauthClientRsaPublicKey,
			OauthClientRsaPublicKey2:  r.Set.OauthClientRsaPublicKey2,
			Comment:                   r.Set.Comment,
		}
		if r.Set.PreAuthorizedRolesList != nil {
			opts.Set.PreAuthorizedRolesList = &PreAuthorizedRolesList{
				PreAuthorizedRolesList: r.Set.PreAuthorizedRolesList.PreAuthorizedRolesList,
			}
		}
		if r.Set.BlockedRolesList != nil {
			opts.Set.BlockedRolesList = &BlockedRolesList{
				BlockedRolesList: r.Set.BlockedRolesList.BlockedRolesList,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &SnowflakeOauthCustomIntegrationUnset{
			Enabled:                  r.Unset.Enabled,
			OauthUseSecondaryRoles:   r.Unset.OauthUseSecondaryRoles,
			NetworkPolicy:            r.Unset.NetworkPolicy,
			OauthClientRsaPublicKey:  r.Unset.OauthClientRsaPublicKey,
			OauthClientRsaPublicKey2: r.Unset.OauthClientRsaPublicKey2,
		}
	}
	return opts
}

func (r *AlterSaml2SecurityIntegrationRequest) toOpts() *AlterSaml2SecurityIntegrationOptions {
	opts := &AlterSaml2SecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,

		RefreshSaml2SnowflakePrivateKey: r.RefreshSaml2SnowflakePrivateKey,
	}
	if r.Set != nil {
		opts.Set = &Saml2IntegrationSet{
			Enabled:                        r.Set.Enabled,
			Saml2Issuer:                    r.Set.Saml2Issuer,
			Saml2SsoUrl:                    r.Set.Saml2SsoUrl,
			Saml2Provider:                  r.Set.Saml2Provider,
			Saml2X509Cert:                  r.Set.Saml2X509Cert,
			AllowedUserDomains:             r.Set.AllowedUserDomains,
			AllowedEmailPatterns:           r.Set.AllowedEmailPatterns,
			Saml2SpInitiatedLoginPageLabel: r.Set.Saml2SpInitiatedLoginPageLabel,
			Saml2EnableSpInitiated:         r.Set.Saml2EnableSpInitiated,
			Saml2SnowflakeX509Cert:         r.Set.Saml2SnowflakeX509Cert,
			Saml2SignRequest:               r.Set.Saml2SignRequest,
			Saml2RequestedNameidFormat:     r.Set.Saml2RequestedNameidFormat,
			Saml2PostLogoutRedirectUrl:     r.Set.Saml2PostLogoutRedirectUrl,
			Saml2ForceAuthn:                r.Set.Saml2ForceAuthn,
			Saml2SnowflakeIssuerUrl:        r.Set.Saml2SnowflakeIssuerUrl,
			Saml2SnowflakeAcsUrl:           r.Set.Saml2SnowflakeAcsUrl,
			Comment:                        r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &Saml2IntegrationUnset{
			Saml2ForceAuthn:            r.Unset.Saml2ForceAuthn,
			Saml2RequestedNameidFormat: r.Unset.Saml2RequestedNameidFormat,
			Saml2PostLogoutRedirectUrl: r.Unset.Saml2PostLogoutRedirectUrl,
			Comment:                    r.Unset.Comment,
		}
	}
	return opts
}

func (r *AlterScimSecurityIntegrationRequest) toOpts() *AlterScimSecurityIntegrationOptions {
	opts := &AlterScimSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ScimIntegrationSet{
			Enabled:       r.Set.Enabled,
			NetworkPolicy: r.Set.NetworkPolicy,
			SyncPassword:  r.Set.SyncPassword,
			Comment:       r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ScimIntegrationUnset{
			Enabled:       r.Unset.Enabled,
			NetworkPolicy: r.Unset.NetworkPolicy,
			SyncPassword:  r.Unset.SyncPassword,
			Comment:       r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropSecurityIntegrationRequest) toOpts() *DropSecurityIntegrationOptions {
	opts := &DropSecurityIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeSecurityIntegrationRequest) toOpts() *DescribeSecurityIntegrationOptions {
	opts := &DescribeSecurityIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r securityIntegrationDescRow) convert() *SecurityIntegrationProperty {
	return &SecurityIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}

func (r *ShowSecurityIntegrationRequest) toOpts() *ShowSecurityIntegrationOptions {
	opts := &ShowSecurityIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r securityIntegrationShowRow) convert() *SecurityIntegration {
	s := &SecurityIntegration{
		Name:            r.Name,
		IntegrationType: r.Type,
		Enabled:         r.Enabled,
		CreatedOn:       r.CreatedOn,
		Category:        r.Category,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s
}
