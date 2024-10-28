package model

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (s *SecretWithClientCredentialsModel) WithOauthScopes(oauthScopes []string) *SecretWithClientCredentialsModel {
	oauthScopesStringVariables := make([]tfconfig.Variable, len(oauthScopes))
	for i, v := range oauthScopes {
		oauthScopesStringVariables[i] = tfconfig.StringVariable(v)
	}

	s.OauthScopes = tfconfig.SetVariable(oauthScopesStringVariables...)
	return s
}
