package model

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

// WithExternalOauthTokenUserMappingClaim was added to satisfy the default builders. The method itself is not generated because its type is not yet supported.
// This method will conflict the generated one when the type is supported.
func (e *ExternalOauthSecurityIntegrationModel) WithExternalOauthTokenUserMappingClaim(externalOauthTokenUserMappingClaim string) *ExternalOauthSecurityIntegrationModel {
	e.ExternalOauthTokenUserMappingClaim = tfconfig.StringVariable(externalOauthTokenUserMappingClaim)
	return e
}
