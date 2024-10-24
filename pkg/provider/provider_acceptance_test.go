package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Provider_configHierarchy(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user := acc.DefaultConfig(t).User
	pass := acc.DefaultConfig(t).Password
	account := acc.DefaultConfig(t).Account
	role := acc.DefaultConfig(t).Role
	host := acc.DefaultConfig(t).Host

	nonExistingUser := "non-existing-user"

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
			// make sure that we fail for incorrect profile
			{
				Config:      providerConfig(testprofiles.IncorrectUserAndPassword),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in provider's config should not be rewritten by a faulty config
			{
				Config: providerConfigWithUserAndPassword(user, pass, testprofiles.IncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, nonExistingUser)
				},
				Config:      providerConfig(testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
				},
				Config: providerConfig(testprofiles.IncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// there is no config (by setting the dir to something different than .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config:      providerConfigWithUserAndPassword(user, pass, testprofiles.Default),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
					t.Setenv(snowflakeenvs.Account, account)
					t.Setenv(snowflakeenvs.Role, role)
					t.Setenv(snowflakeenvs.Host, host)
				},
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// make sure the teardown is fine by using a correct env config at the end
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					testenvs.AssertEnvSet(t, snowflakeenvs.Password)
					testenvs.AssertEnvSet(t, snowflakeenvs.Account)
					testenvs.AssertEnvSet(t, snowflakeenvs.Role)
					testenvs.AssertEnvSet(t, snowflakeenvs.Host)
				},
				Config: emptyProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_configureClientOnceSwitching(t *testing.T) {
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
			// client setup is incorrect
			{
				Config:      providerConfig(testprofiles.IncorrectUserAndPassword),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// in this step we simulate the situation when we want to use client configured once, but it was faulty last time
			{
				PreConfig: func() {
					t.Setenv(string(testenvs.ConfigureClientOnce), "true")
				},
				Config: emptyProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_useNonExistentDefaultParams(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	nonExisting := "NON-EXISTENT"

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
				Config:      providerConfigWithRole(testprofiles.Default, nonExisting),
				ExpectError: regexp.MustCompile("Role 'NON-EXISTENT' specified in the connect string does not exist or not authorized."),
			},
			{
				Config:      providerConfigWithWarehouse(testprofiles.Default, nonExisting),
				ExpectError: regexp.MustCompile("The requested warehouse does not exist or not authorized."),
			},
			{
				Config: providerConfigWithWarehouseAndDisabledValidation(testprofiles.Default, nonExisting),
			},
		},
	})
}

func TestAcc_Provider_invalidConfigurations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      providerConfigWithClientIp(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile("expected client_ip to contain a valid IP"),
			},
			{
				Config:      providerConfigWithProtocol(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile("invalid protocol: INVALID"),
			},
			{
				Config:      providerConfigWithPort(testprofiles.Default, 123456789),
				ExpectError: regexp.MustCompile(`expected "port" to be a valid port number or 0, got: 123456789`),
			},
			{
				Config:      providerConfigWithAuthType(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile("invalid authenticator type: INVALID"),
			},
			{
				Config:      providerConfigWithOktaUrl(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile(`expected "okta_url" to have a host, got invalid`),
			},
			{
				Config:      providerConfigWithTimeout(testprofiles.Default, "login_timeout", -1),
				ExpectError: regexp.MustCompile(`expected login_timeout to be at least \(0\), got -1`),
			},
			{
				Config:      providerConfigWithTokenEndpoint(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile(`expected "token_endpoint" to have a host, got invalid`),
			},
		},
	})
}

// TODO(SNOW-1754319): for JWT auth flow, check setting authenticator value as `SNOWFLAKE_JWT`.
// This will ensure https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2983 is solved.

func emptyProviderConfig() string {
	return `
provider "snowflake" {
}` + datasourceConfig()
}

func providerConfig(profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
}
`, profile) + datasourceConfig()
}

func providerConfigWithRole(profile, role string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	role    = "%[2]s"
}
`, profile, role) + datasourceConfig()
}

func providerConfigWithWarehouse(profile, warehouse string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	warehouse    = "%[2]s"
}
`, profile, warehouse) + datasourceConfig()
}

func providerConfigWithWarehouseAndDisabledValidation(profile, warehouse string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	warehouse    = "%[2]s"
	validate_default_parameters = "false"
}
`, profile, warehouse) + datasourceConfig()
}

func providerConfigWithProtocol(profile, protocol string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	protocol    = "%[2]s"
}
`, profile, protocol) + datasourceConfig()
}

func providerConfigWithPort(profile string, port int) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	port    = %[2]d
}
`, profile, port) + datasourceConfig()
}

func providerConfigWithAuthType(profile, authType string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	authenticator    = "%[2]s"
}
`, profile, authType) + datasourceConfig()
}

func providerConfigWithOktaUrl(profile, oktaUrl string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	okta_url    = "%[2]s"
}
`, profile, oktaUrl) + datasourceConfig()
}

func providerConfigWithTimeout(profile, timeoutName string, timeoutSeconds int) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	%[2]s    = %[3]d
}
`, profile, timeoutName, timeoutSeconds) + datasourceConfig()
}

func providerConfigWithTokenEndpoint(profile, tokenEndpoint string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	token_accessor {
		token_endpoint = "%[2]s"
		refresh_token = "refresh_token"
		client_id = "client_id"
		client_secret = "client_secret"
		redirect_uri = "redirect_uri"
	}
}
`, profile, tokenEndpoint) + datasourceConfig()
}

func providerConfigWithClientIp(profile, clientIp string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	client_ip    = "%[2]s"
}
`, profile, clientIp) + datasourceConfig()
}

func providerConfigWithUser(user string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	user = "%[1]s"
	profile = "%[2]s"
}
`, user, profile) + datasourceConfig()
}

func providerConfigWithUserAndPassword(user string, pass string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	user = "%[1]s"
	password = "%[2]s"
	profile = "%[3]s"
}
`, user, pass, profile) + datasourceConfig()
}

func datasourceConfig() string {
	return fmt.Sprintf(`
data snowflake_database "t" {
	name = "%s"
}`, acc.TestDatabaseName)
}
