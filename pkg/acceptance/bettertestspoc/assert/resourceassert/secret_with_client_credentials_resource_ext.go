package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopes(expected []sdk.ApiIntegrationScope) *SecretWithClientCredentialsResourceAssert {
	for _, v := range expected {
		s.AddAssertion(assert.ValueSet("oauth_scopes.*", v.Scope))
	}
	return s
}
