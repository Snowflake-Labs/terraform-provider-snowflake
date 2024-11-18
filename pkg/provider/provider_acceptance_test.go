package provider_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
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

func TestAcc_Provider_tomlConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user := acc.DefaultConfig(t).User
	pass := acc.DefaultConfig(t).Password
	account := acc.DefaultConfig(t).Account

	oktaUrl, err := url.Parse("https://example.com")
	require.NoError(t, err)

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
				Config: providerConfig(testprofiles.CompleteFields),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, &gosnowflake.Config{
						Account:                   account,
						User:                      user,
						Password:                  pass,
						Warehouse:                 "SNOWFLAKE",
						Role:                      "ACCOUNTADMIN",
						ValidateDefaultParameters: gosnowflake.ConfigBoolTrue,
						ClientIP:                  net.ParseIP("1.2.3.4"),
						Protocol:                  "https",
						Host:                      fmt.Sprintf("%s.snowflakecomputing.com", account),
						Params: map[string]*string{
							"foo": sdk.Pointer("bar"),
						},
						Port:                           443,
						Authenticator:                  gosnowflake.AuthTypeSnowflake,
						PasscodeInPassword:             false,
						OktaURL:                        oktaUrl,
						LoginTimeout:                   30 * time.Second,
						RequestTimeout:                 40 * time.Second,
						JWTExpireTimeout:               50 * time.Second,
						ClientTimeout:                  10 * time.Second,
						JWTClientTimeout:               20 * time.Second,
						ExternalBrowserTimeout:         60 * time.Second,
						MaxRetryCount:                  1,
						Application:                    "terraform-provider-snowflake",
						InsecureMode:                   true,
						OCSPFailOpen:                   gosnowflake.OCSPFailOpenTrue,
						Token:                          "token",
						KeepSessionAlive:               true,
						DisableTelemetry:               true,
						Tracing:                        string(sdk.DriverLogLevelInfo),
						TmpDirPath:                     ".",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
						DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
					}, config)
					assert.Equal(t, string(sdk.DriverLogLevelInfo), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_envConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user := acc.DefaultConfig(t).User
	pass := acc.DefaultConfig(t).Password
	account := acc.DefaultConfig(t).Account

	accountParts := strings.SplitN(account, "-", 2)

	oktaUrlFromEnv, err := url.Parse("https://example-env.com")
	require.NoError(t, err)

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
				PreConfig: func() {
					t.Setenv(snowflakeenvs.AccountName, accountParts[1])
					t.Setenv(snowflakeenvs.OrganizationName, accountParts[0])
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
					t.Setenv(snowflakeenvs.Warehouse, "SNOWFLAKE")
					t.Setenv(snowflakeenvs.Protocol, "https")
					t.Setenv(snowflakeenvs.Port, "443")
					// do not set token - it should be propagated from TOML
					t.Setenv(snowflakeenvs.Role, "ACCOUNTADMIN")
					t.Setenv(snowflakeenvs.Authenticator, "snowflake")
					t.Setenv(snowflakeenvs.ValidateDefaultParameters, "false")
					t.Setenv(snowflakeenvs.ClientIp, "2.2.2.2")
					t.Setenv(snowflakeenvs.Host, "")
					t.Setenv(snowflakeenvs.Authenticator, "")
					t.Setenv(snowflakeenvs.Passcode, "")
					t.Setenv(snowflakeenvs.PasscodeInPassword, "false")
					t.Setenv(snowflakeenvs.OktaUrl, "https://example-env.com")
					t.Setenv(snowflakeenvs.LoginTimeout, "100")
					t.Setenv(snowflakeenvs.RequestTimeout, "200")
					t.Setenv(snowflakeenvs.JwtExpireTimeout, "300")
					t.Setenv(snowflakeenvs.ClientTimeout, "400")
					t.Setenv(snowflakeenvs.JwtClientTimeout, "500")
					t.Setenv(snowflakeenvs.ExternalBrowserTimeout, "600")
					t.Setenv(snowflakeenvs.InsecureMode, "false")
					t.Setenv(snowflakeenvs.OcspFailOpen, "false")
					t.Setenv(snowflakeenvs.KeepSessionAlive, "false")
					t.Setenv(snowflakeenvs.DisableTelemetry, "false")
					t.Setenv(snowflakeenvs.ClientRequestMfaToken, "false")
					t.Setenv(snowflakeenvs.ClientStoreTemporaryCredential, "false")
					t.Setenv(snowflakeenvs.DisableQueryContextCache, "false")
					t.Setenv(snowflakeenvs.IncludeRetryReason, "false")
					t.Setenv(snowflakeenvs.MaxRetryCount, "2")
					t.Setenv(snowflakeenvs.DriverTracing, string(sdk.DriverLogLevelDebug))
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfig(testprofiles.CompleteFieldsInvalid),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, &gosnowflake.Config{
						Account:                   account,
						User:                      user,
						Password:                  pass,
						Warehouse:                 "SNOWFLAKE",
						Role:                      "ACCOUNTADMIN",
						ValidateDefaultParameters: gosnowflake.ConfigBoolFalse,
						ClientIP:                  net.ParseIP("2.2.2.2"),
						Protocol:                  "https",
						Params: map[string]*string{
							"foo": sdk.Pointer("bar"),
						},
						Host:                           fmt.Sprintf("%s.snowflakecomputing.com", account),
						Port:                           443,
						Authenticator:                  gosnowflake.AuthTypeSnowflake,
						PasscodeInPassword:             false,
						OktaURL:                        oktaUrlFromEnv,
						LoginTimeout:                   100 * time.Second,
						RequestTimeout:                 200 * time.Second,
						JWTExpireTimeout:               300 * time.Second,
						ClientTimeout:                  400 * time.Second,
						JWTClientTimeout:               500 * time.Second,
						ExternalBrowserTimeout:         600 * time.Second,
						MaxRetryCount:                  2,
						Application:                    "terraform-provider-snowflake",
						InsecureMode:                   true,
						OCSPFailOpen:                   gosnowflake.OCSPFailOpenFalse,
						Token:                          "token",
						KeepSessionAlive:               true,
						DisableTelemetry:               true,
						Tracing:                        string(sdk.DriverLogLevelDebug),
						TmpDirPath:                     "../",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolFalse,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolFalse,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolFalse,
						DisableConsoleLogin:            gosnowflake.ConfigBoolFalse,
					}, config)
					assert.Equal(t, string(sdk.DriverLogLevelDebug), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_tfConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user := acc.DefaultConfig(t).User
	pass := acc.DefaultConfig(t).Password
	account := acc.DefaultConfig(t).Account

	accountParts := strings.SplitN(account, "-", 2)
	orgName, accountName := accountParts[0], accountParts[1]

	oktaUrlFromTf, err := url.Parse("https://example-tf.com")
	require.NoError(t, err)

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
				PreConfig: func() {
					t.Setenv(snowflakeenvs.OrganizationName, "invalid")
					t.Setenv(snowflakeenvs.AccountName, "invalid")
					t.Setenv(snowflakeenvs.User, "invalid")
					t.Setenv(snowflakeenvs.Password, "invalid")
					t.Setenv(snowflakeenvs.Warehouse, "invalid")
					t.Setenv(snowflakeenvs.Protocol, "invalid")
					t.Setenv(snowflakeenvs.Port, "-1")
					t.Setenv(snowflakeenvs.Token, "")
					t.Setenv(snowflakeenvs.Role, "invalid")
					t.Setenv(snowflakeenvs.ValidateDefaultParameters, "false")
					t.Setenv(snowflakeenvs.ClientIp, "2.2.2.2")
					t.Setenv(snowflakeenvs.Host, "")
					t.Setenv(snowflakeenvs.Authenticator, "invalid")
					t.Setenv(snowflakeenvs.Passcode, "")
					t.Setenv(snowflakeenvs.PasscodeInPassword, "false")
					t.Setenv(snowflakeenvs.OktaUrl, "https://example-env.com")
					t.Setenv(snowflakeenvs.LoginTimeout, "100")
					t.Setenv(snowflakeenvs.RequestTimeout, "200")
					t.Setenv(snowflakeenvs.JwtExpireTimeout, "300")
					t.Setenv(snowflakeenvs.ClientTimeout, "400")
					t.Setenv(snowflakeenvs.JwtClientTimeout, "500")
					t.Setenv(snowflakeenvs.ExternalBrowserTimeout, "600")
					t.Setenv(snowflakeenvs.InsecureMode, "false")
					t.Setenv(snowflakeenvs.OcspFailOpen, "false")
					t.Setenv(snowflakeenvs.KeepSessionAlive, "false")
					t.Setenv(snowflakeenvs.DisableTelemetry, "false")
					t.Setenv(snowflakeenvs.ClientRequestMfaToken, "false")
					t.Setenv(snowflakeenvs.ClientStoreTemporaryCredential, "false")
					t.Setenv(snowflakeenvs.DisableQueryContextCache, "false")
					t.Setenv(snowflakeenvs.IncludeRetryReason, "false")
					t.Setenv(snowflakeenvs.MaxRetryCount, "2")
					t.Setenv(snowflakeenvs.DriverTracing, string(sdk.DriverLogLevelDebug))
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfigAllFields(testprofiles.CompleteFieldsInvalid, orgName, accountName, user, pass),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, &gosnowflake.Config{
						Account:                   account,
						User:                      user,
						Password:                  pass,
						Warehouse:                 "SNOWFLAKE",
						Role:                      "ACCOUNTADMIN",
						ValidateDefaultParameters: gosnowflake.ConfigBoolTrue,
						ClientIP:                  net.ParseIP("3.3.3.3"),
						Protocol:                  "https",
						Params: map[string]*string{
							"foo": sdk.Pointer("piyo"),
						},
						Host:                           fmt.Sprintf("%s.snowflakecomputing.com", account),
						Port:                           443,
						Authenticator:                  gosnowflake.AuthTypeSnowflake,
						PasscodeInPassword:             false,
						OktaURL:                        oktaUrlFromTf,
						LoginTimeout:                   101 * time.Second,
						RequestTimeout:                 201 * time.Second,
						JWTExpireTimeout:               301 * time.Second,
						ClientTimeout:                  401 * time.Second,
						JWTClientTimeout:               501 * time.Second,
						ExternalBrowserTimeout:         601 * time.Second,
						MaxRetryCount:                  3,
						Application:                    "terraform-provider-snowflake",
						InsecureMode:                   true,
						OCSPFailOpen:                   gosnowflake.OCSPFailOpenTrue,
						Token:                          "token",
						KeepSessionAlive:               true,
						DisableTelemetry:               true,
						Tracing:                        string(sdk.DriverLogLevelInfo),
						TmpDirPath:                     "../../",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
						DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
					}, config)
					assert.Equal(t, string(sdk.DriverLogLevelInfo), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
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
			// check that using a non-existing warehouse with disabled verification succeeds
			{
				Config: providerConfigWithWarehouseAndDisabledValidation(testprofiles.Default, nonExisting),
			},
		},
	})
}

// prove we can use tri-value booleans, similarly to the ones in resources
func TestAcc_Provider_triValueBoolean(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	resource.Test(t, resource.TestCase{
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
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.97.0"),
				Config:            providerConfigWithClientStoreTemporaryCredential(testprofiles.Default, `true`),
			},
			{
				// Use the default TOML config again.
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, "")
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   providerConfigWithClientStoreTemporaryCredential(testprofiles.Default, `true`),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   providerConfigWithClientStoreTemporaryCredential(testprofiles.Default, `"true"`),
			},
		},
	})
}

func TestAcc_Provider_sessionParameters(t *testing.T) {
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
				Config: providerWithParamsConfig(testprofiles.Default, 31337),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_unsafe_execute.t", "query_results.#", "1"),
					resource.TestCheckResourceAttr("snowflake_unsafe_execute.t", "query_results.0.value", "31337"),
				),
			},
		},
	})
}

func TestAcc_Provider_JwtAuth(t *testing.T) {
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
			// authenticate with unencrypted private key
			{
				Config: providerConfigWithAuthenticator(testprofiles.JwtAuth, sdk.AuthenticationTypeJwt),
			},
			// authenticate with unencrypted private key with a legacy authenticator value
			// solves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2983
			{
				Config: providerConfigWithAuthenticator(testprofiles.JwtAuth, sdk.AuthenticationTypeJwtLegacy),
			},
			// authenticate with encrypted private key
			{
				Config: providerConfigWithAuthenticator(testprofiles.EncryptedJwtAuth, sdk.AuthenticationTypeJwt),
			},
		},
	})
}

func TestAcc_Provider_SnowflakeAuth(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfigWithAuthenticator(testprofiles.Default, sdk.AuthenticationTypeSnowflake),
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
				ExpectError: regexp.MustCompile("invalid protocol: invalid"),
			},
			{
				Config:      providerConfigWithPort(testprofiles.Default, 123456789),
				ExpectError: regexp.MustCompile(`expected "port" to be a valid port number or 0, got: 123456789`),
			},
			{
				Config:      providerConfigWithAuthType(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile("invalid authenticator type: invalid"),
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
			{
				Config:      providerConfigWithLogLevel(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile(`invalid driver log level: invalid`),
			},
			{
				Config: providerConfig("non-existing"),
				// .* is used to match the error message regarding of the home user location
				ExpectError: regexp.MustCompile(`profile "non-existing" not found in file .*.snowflake/config`),
			},
		},
	})
}

func providerConfigWithAuthenticator(profile string, authenticator sdk.AuthenticationType) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithAuthenticator(string(authenticator))) + datasourceConfig(t)
}

func emptyProviderConfig() string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider()) + datasourceConfig(t)
}

func providerConfig(profile string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile)) + datasourceConfig(t)
}

func providerConfigWithRole(profile, role string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithRole(role)) + datasourceConfig(t)
}

func providerConfigWithWarehouse(profile, warehouse string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithWarehouse(warehouse)) + datasourceConfig(t)
}

func providerConfigWithClientStoreTemporaryCredential(profile, clientStoreTemporaryCredential string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithClientStoreTemporaryCredential(clientStoreTemporaryCredential)) + datasourceConfig(t)
}

func providerConfigWithWarehouseAndDisabledValidation(profile, warehouse string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithWarehouse(warehouse).WithValidateDefaultParameters("false")) + datasourceConfig(t)
}

func providerConfigWithProtocol(profile, protocol string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithProtocol(protocol)) + datasourceConfig(t)
}

func providerConfigWithPort(profile string, port int) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithPort(port)) + datasourceConfig(t)
}

func providerConfigWithAuthType(profile, authType string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithAuthenticator(authType)) + datasourceConfig(t)
}

func providerConfigWithOktaUrl(profile, oktaUrl string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithOktaUrl(oktaUrl)) + datasourceConfig(t)
}

func providerConfigWithTimeout(profile, timeoutName string, timeoutSeconds int) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithLoginTimeout(timeoutSeconds)) + datasourceConfig(t)
}

func providerConfigWithTokenEndpoint(profile, tokenEndpoint string) string {
	return config.ProviderFromModelPoc(
		t,
		providermodel.SnowflakeProvider().
			WithProfile(profile).
			WithTokenAccessorValue(
				tfconfig.ObjectVariable(
					map[string]tfconfig.Variable{
						"token_endpoint": tfconfig.StringVariable(tokenEndpoint),
						"refresh_token":  tfconfig.StringVariable("refresh_token"),
						"client_id":      tfconfig.StringVariable("client_id"),
						"client_secret":  tfconfig.StringVariable("client_secret"),
						"redirect_uri":   tfconfig.StringVariable("redirect_uri"),
					},
				),
			),
	) + datasourceConfig(t)
}

func providerConfigWithLogLevel(profile, logLevel string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithDriverTracing(logLevel)) + datasourceConfig(t)
}

func providerConfigWithClientIp(profile, clientIp string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithClientIp(clientIp)) + datasourceConfig(t)
}

func providerConfigWithUser(user string, profile string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithUser(user)) + datasourceConfig(t)
}

func providerConfigWithUserAndPassword(user string, pass string, profile string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithUser(user).WithPassword(pass)) + datasourceConfig(t)
}

func providerConfigWithNewAccountId(profile, orgName, accountName string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithAccountName(accountName).WithOrganizationName(orgName)) + datasourceConfig(t)
}

func providerConfigComplete(profile, user, password, orgName, accountName string) string {
	return config.ProviderFromModelPoc(
		t,
		providermodel.SnowflakeProvider().
			WithProfile(profile).
			WithUser(user).
			WithPassword(password).
			WithOrganizationName(orgName).
			WithAccountName(accountName).
			WithWarehouse("SNOWFLAKE"),
	) + datasourceConfig(t)
}

func datasourceConfig(t *testing.T) string {
	return config.ResourceFromModelPoc(t, model.Database("t", acc.TestDatabaseName))
}

func providerConfigAllFields(profile, orgName, accountName, user, password string) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().AllFields(profile, orgName, accountName, user, password)) + datasourceConfig(t)
}

// TODO(SNOW-1348325): Use parameter data source with `IN SESSION` filtering.
func providerWithParamsConfig(profile string, statementTimeoutInSeconds int) string {
	return config.ProviderFromModelPoc(t, providermodel.SnowflakeProvider().WithProfile(profile).WithParamsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"statement_timeout_in_seconds": tfconfig.IntegerVariable(statementTimeoutInSeconds),
			},
		),
	)) + unsafeExecuteShowSessionParameter()
}

func unsafeExecuteShowSessionParameter() string {
	return `
resource snowflake_unsafe_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
