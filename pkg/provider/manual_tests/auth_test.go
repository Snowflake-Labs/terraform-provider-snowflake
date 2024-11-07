package manual

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// This is a manual test for authenticating with Okta.
func TestAcc_Provider_OktaAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfigWithAuthenticator("okta", sdk.AuthenticationTypeOkta),
			},
		},
	})
}

// This test requires manual action due to MFA. Make sure the user does not have a positive `mins_to_bypass_mfa` in `SHOW USERS`.
func TestAcc_Provider_UsernamePasswordMfaAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// ensure MFA is checked here - accept login on your MFA device
			{
				Config: providerConfigWithAuthenticator(testprofiles.Default, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
			// check that MFA login is cached - this step should not require manual action
			{
				Config: providerConfigWithAuthenticator(testprofiles.Default, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
		},
	})
}

// This test requires manual action due to MFA. Make sure the user does not have a positive `mins_to_bypass_mfa` in `SHOW USERS`.
func TestAcc_Provider_UsernamePasswordMfaAuthWithPasscode(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// ensure MFA is checked here - accept access to keychain on your device
			{
				Config: providerConfigWithAuthenticator(testprofiles.DefaultWithPasscode, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
			// check that MFA login is cached - this step should not require manual action
			{
				Config: providerConfigWithAuthenticator(testprofiles.DefaultWithPasscode, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
		},
	})
}

func providerConfigWithAuthenticator(profile string, authenticator sdk.AuthenticationType) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	authenticator    = "%[2]s"
}
`, profile, authenticator) + datasourceConfig()
}

func datasourceConfig() string {
	return `
data snowflake_database "t" {
	name = "SNOWFLAKE"
}`
}
