package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SecretWithAuthorizationCodeResourceAssert) HasOauthRefreshTokenNoEmpty() *SecretWithAuthorizationCodeResourceAssert {
	s.AddAssertion(assert.ValuePresent("oauth_refresh_token_expiry_time"))
	return s
}
