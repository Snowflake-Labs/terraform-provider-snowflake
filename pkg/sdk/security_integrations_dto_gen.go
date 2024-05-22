package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateOauthPartnerSecurityIntegrationOptions] = new(CreateOauthPartnerSecurityIntegrationRequest)
	_ optionsProvider[CreateOauthCustomSecurityIntegrationOptions]  = new(CreateOauthCustomSecurityIntegrationRequest)
	_ optionsProvider[CreateSaml2SecurityIntegrationOptions]                 = new(CreateSaml2SecurityIntegrationRequest)
	_ optionsProvider[CreateScimSecurityIntegrationOptions]                  = new(CreateScimSecurityIntegrationRequest)
	_ optionsProvider[AlterOauthPartnerSecurityIntegrationOptions]  = new(AlterOauthPartnerSecurityIntegrationRequest)
	_ optionsProvider[AlterOauthCustomSecurityIntegrationOptions]   = new(AlterOauthCustomSecurityIntegrationRequest)
	_ optionsProvider[AlterSaml2SecurityIntegrationOptions]                  = new(AlterSaml2SecurityIntegrationRequest)
	_ optionsProvider[AlterScimSecurityIntegrationOptions]                   = new(AlterScimSecurityIntegrationRequest)
	_ optionsProvider[DropSecurityIntegrationOptions]                        = new(DropSecurityIntegrationRequest)
	_ optionsProvider[DescribeSecurityIntegrationOptions]                    = new(DescribeSecurityIntegrationRequest)
	_ optionsProvider[ShowSecurityIntegrationOptions]                        = new(ShowSecurityIntegrationRequest)
)

type CreateOauthPartnerSecurityIntegrationRequest struct {
	OrReplace                 *bool
	IfNotExists               *bool
	name                      AccountObjectIdentifier              // required
	OauthClient               OauthSecurityIntegrationClientOption // required
	OauthRedirectUri          *string
	Enabled                   *bool
	OauthIssueRefreshTokens   *bool
	OauthRefreshTokenValidity *int
	OauthUseSecondaryRoles    *OauthSecurityIntegrationUseSecondaryRolesOption
	BlockedRolesList          *BlockedRolesListRequest
	Comment                   *string
}

func (r *CreateOauthPartnerSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type BlockedRolesListRequest struct {
	BlockedRolesList []AccountObjectIdentifier
}

type CreateOauthCustomSecurityIntegrationRequest struct {
	OrReplace                   *bool
	IfNotExists                 *bool
	name                        AccountObjectIdentifier                  // required
	OauthClientType             OauthSecurityIntegrationClientTypeOption // required
	OauthRedirectUri            string                                   // required
	Enabled                     *bool
	OauthAllowNonTlsRedirectUri *bool
	OauthEnforcePkce            *bool
	OauthUseSecondaryRoles      *OauthSecurityIntegrationUseSecondaryRolesOption
	PreAuthorizedRolesList      *PreAuthorizedRolesListRequest
	BlockedRolesList            *BlockedRolesListRequest
	OauthIssueRefreshTokens     *bool
	OauthRefreshTokenValidity   *int
	NetworkPolicy               *AccountObjectIdentifier
	OauthClientRsaPublicKey     *string
	OauthClientRsaPublicKey2    *string
	Comment                     *string
}

func (r *CreateOauthCustomSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type PreAuthorizedRolesListRequest struct {
	PreAuthorizedRolesList []AccountObjectIdentifier
}

type CreateSaml2SecurityIntegrationRequest struct {
	OrReplace                      *bool
	IfNotExists                    *bool
	name                           AccountObjectIdentifier // required
	Enabled                        bool                    // required
	Saml2Issuer                    string                  // required
	Saml2SsoUrl                    string                  // required
	Saml2Provider                  string                  // required
	Saml2X509Cert                  string                  // required
	AllowedUserDomains             []UserDomain
	AllowedEmailPatterns           []EmailPattern
	Saml2SpInitiatedLoginPageLabel *string
	Saml2EnableSpInitiated         *bool
	Saml2SnowflakeX509Cert         *string
	Saml2SignRequest               *bool
	Saml2RequestedNameidFormat     *string
	Saml2PostLogoutRedirectUrl     *string
	Saml2ForceAuthn                *bool
	Saml2SnowflakeIssuerUrl        *string
	Saml2SnowflakeAcsUrl           *string
	Comment                        *string
}

func (r *CreateSaml2SecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type CreateScimSecurityIntegrationRequest struct {
	OrReplace     *bool
	IfNotExists   *bool
	name          AccountObjectIdentifier                 // required
	Enabled       bool                                    // required
	ScimClient    ScimSecurityIntegrationScimClientOption // required
	RunAsRole     ScimSecurityIntegrationRunAsRoleOption  // required
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

func (r *CreateScimSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

type AlterOauthPartnerSecurityIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Set       *OauthPartnerIntegrationSetRequest
	Unset     *OauthPartnerIntegrationUnsetRequest
}

type OauthPartnerIntegrationSetRequest struct {
	Enabled                   *bool
	OauthRedirectUri          *string
	OauthIssueRefreshTokens   *bool
	OauthRefreshTokenValidity *int
	OauthUseSecondaryRoles    *OauthSecurityIntegrationUseSecondaryRolesOption
	BlockedRolesList          *BlockedRolesListRequest
	Comment                   *string
}

type OauthPartnerIntegrationUnsetRequest struct {
	Enabled                *bool
	OauthUseSecondaryRoles *bool
}

type AlterOauthCustomSecurityIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Set       *OauthCustomIntegrationSetRequest
	Unset     *OauthCustomIntegrationUnsetRequest
}

type OauthCustomIntegrationSetRequest struct {
	Enabled                     *bool
	OauthRedirectUri            *string
	OauthAllowNonTlsRedirectUri *bool
	OauthEnforcePkce            *bool
	OauthUseSecondaryRoles      *OauthSecurityIntegrationUseSecondaryRolesOption
	PreAuthorizedRolesList      *PreAuthorizedRolesListRequest
	BlockedRolesList            *BlockedRolesListRequest
	OauthIssueRefreshTokens     *bool
	OauthRefreshTokenValidity   *int
	NetworkPolicy               *AccountObjectIdentifier
	OauthClientRsaPublicKey     *string
	OauthClientRsaPublicKey2    *string
	Comment                     *string
}

type OauthCustomIntegrationUnsetRequest struct {
	Enabled                  *bool
	OauthUseSecondaryRoles   *bool
	NetworkPolicy            *bool
	OauthClientRsaPublicKey  *bool
	OauthClientRsaPublicKey2 *bool
}

type AlterSaml2SecurityIntegrationRequest struct {
	IfExists                        *bool
	name                            AccountObjectIdentifier // required
	SetTags                         []TagAssociation
	UnsetTags                       []ObjectIdentifier
	Set                             *Saml2IntegrationSetRequest
	Unset                           *Saml2IntegrationUnsetRequest
	RefreshSaml2SnowflakePrivateKey *bool
}

type Saml2IntegrationSetRequest struct {
	Enabled                        *bool
	Saml2Issuer                    *string
	Saml2SsoUrl                    *string
	Saml2Provider                  *string
	Saml2X509Cert                  *string
	AllowedUserDomains             []UserDomain
	AllowedEmailPatterns           []EmailPattern
	Saml2SpInitiatedLoginPageLabel *string
	Saml2EnableSpInitiated         *bool
	Saml2SnowflakeX509Cert         *string
	Saml2SignRequest               *bool
	Saml2RequestedNameidFormat     *string
	Saml2PostLogoutRedirectUrl     *string
	Saml2ForceAuthn                *bool
	Saml2SnowflakeIssuerUrl        *string
	Saml2SnowflakeAcsUrl           *string
	Comment                        *string
}

type Saml2IntegrationUnsetRequest struct {
	Saml2ForceAuthn            *bool
	Saml2RequestedNameidFormat *bool
	Saml2PostLogoutRedirectUrl *bool
	Comment                    *bool
}

type AlterScimSecurityIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Set       *ScimIntegrationSetRequest
	Unset     *ScimIntegrationUnsetRequest
}

type ScimIntegrationSetRequest struct {
	Enabled       *bool
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

type ScimIntegrationUnsetRequest struct {
	Enabled       *bool
	NetworkPolicy *bool
	SyncPassword  *bool
	Comment       *bool
}

type DropSecurityIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type DescribeSecurityIntegrationRequest struct {
	name AccountObjectIdentifier // required
}

type ShowSecurityIntegrationRequest struct {
	Like *Like
}
