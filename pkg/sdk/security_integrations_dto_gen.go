package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSAML2SecurityIntegrationOptions]           = new(CreateSAML2SecurityIntegrationRequest)
	_ optionsProvider[CreateSCIMSecurityIntegrationOptions]            = new(CreateSCIMSecurityIntegrationRequest)
	_ optionsProvider[AlterSAML2IntegrationSecurityIntegrationOptions] = new(AlterSAML2IntegrationSecurityIntegrationRequest)
	_ optionsProvider[AlterSCIMIntegrationSecurityIntegrationOptions]  = new(AlterSCIMIntegrationSecurityIntegrationRequest)
	_ optionsProvider[DropSecurityIntegrationOptions]                  = new(DropSecurityIntegrationRequest)
	_ optionsProvider[DescribeSecurityIntegrationOptions]              = new(DescribeSecurityIntegrationRequest)
	_ optionsProvider[ShowSecurityIntegrationOptions]                  = new(ShowSecurityIntegrationRequest)
)

type CreateSAML2SecurityIntegrationRequest struct {
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

type CreateSCIMSecurityIntegrationRequest struct {
	OrReplace     *bool
	IfNotExists   *bool
	name          AccountObjectIdentifier // required
	Enabled       bool                    // required
	ScimClient    string                  // required
	RunAsRole     string                  // required
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

type AlterSAML2IntegrationSecurityIntegrationRequest struct {
	IfExists                        *bool
	name                            AccountObjectIdentifier // required
	SetTags                         []TagAssociation
	UnsetTags                       []ObjectIdentifier
	Set                             *SAML2IntegrationSetRequest
	Unset                           *SAML2IntegrationUnsetRequest
	RefreshSaml2SnowflakePrivateKey *bool
}

type SAML2IntegrationSetRequest struct {
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

type SAML2IntegrationUnsetRequest struct {
	Enabled                    *bool
	Saml2ForceAuthn            *bool
	Saml2RequestedNameidFormat *bool
	Saml2PostLogoutRedirectUrl *bool
	Comment                    *bool
}

type AlterSCIMIntegrationSecurityIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	Set       *SCIMIntegrationSetRequest
	Unset     *SCIMIntegrationUnsetRequest
}

type SCIMIntegrationSetRequest struct {
	Enabled       *bool
	NetworkPolicy *AccountObjectIdentifier
	SyncPassword  *bool
	Comment       *string
}

type SCIMIntegrationUnsetRequest struct {
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
