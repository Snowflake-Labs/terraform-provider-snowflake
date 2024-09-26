package resourceassert

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopes(expected []sdk.ApiIntegrationScope) *SecretWithClientCredentialsResourceAssert {
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("oauth_scopes.%d", i), v.Scope))
	}
	return s
}
