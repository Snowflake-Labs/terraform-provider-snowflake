package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSaml2SecurityIntegrationOptions] = new(CreateSaml2SecurityIntegrationRequest)
	_ optionsProvider[CreateScimSecurityIntegrationOptions]  = new(CreateScimSecurityIntegrationRequest)
	_ optionsProvider[AlterSaml2SecurityIntegrationOptions]  = new(AlterSaml2SecurityIntegrationRequest)
	_ optionsProvider[AlterScimSecurityIntegrationOptions]   = new(AlterScimSecurityIntegrationRequest)
	_ optionsProvider[DropSecurityIntegrationOptions]        = new(DropSecurityIntegrationRequest)
	_ optionsProvider[DescribeSecurityIntegrationOptions]    = new(DescribeSecurityIntegrationRequest)
	_ optionsProvider[ShowSecurityIntegrationOptions]        = new(ShowSecurityIntegrationRequest)
)

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
