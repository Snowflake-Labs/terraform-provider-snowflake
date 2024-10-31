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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAcc_Provider_tomlConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user := acc.DefaultConfig(t).User
	pass := acc.DefaultConfig(t).Password
	account := acc.DefaultConfig(t).Account

	accountParts := strings.SplitN(account, "-", 2)
	filename := testhelpers.TestFile(t, "config.toml", getTomlConfig(accountParts[1], accountParts[0], user, pass))

	oktaUrl, err := url.Parse("https://example.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			t.Setenv("SNOWFLAKE_CONFIG_PATH", filename)
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
						Tracing:                        "info",
						TmpDirPath:                     ".",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
						DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
					}, config)

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
	filename := testhelpers.TestFile(t, "config.toml", getTomlConfig("incorrect", "incorrect", "incorrect", "incorrect"))

	oktaUrlFromEnv, err := url.Parse("https://example-env.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			t.Setenv("SNOWFLAKE_CONFIG_PATH", filename)
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
					// do not set warehouse, protocol, port, token and role - they should be propagated from TOML
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
					t.Setenv(snowflakeenvs.DriverTracing, "debug")
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfig(testprofiles.CompleteFields),
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
						Tracing:                        "DEBUG",
						TmpDirPath:                     "../",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolFalse,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolFalse,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolFalse,
						DisableConsoleLogin:            gosnowflake.ConfigBoolFalse,
					}, config)

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
	filename := testhelpers.TestFile(t, "config.toml", getTomlConfig("incorrect", "incorrect", "incorrect", "incorrect"))

	oktaUrlFromTf, err := url.Parse("https://example-tf.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			t.Setenv("SNOWFLAKE_CONFIG_PATH", filename)
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
					t.Setenv(snowflakeenvs.DriverTracing, "debug")
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfigAllFields(testprofiles.CompleteFields, orgName, accountName, user, pass),
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
						Protocol:                  "HTTPS",
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
						Tracing:                        "INFO",
						TmpDirPath:                     "../../",
						ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
						ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
						DisableQueryContextCache:       true,
						IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
						DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
					}, config)

					return nil
				},
			},
		},
	})
}

func getTomlConfig(accountName, orgName, user, password string) []byte {
	return []byte(fmt.Sprintf(`
[complete_fields]
accountname='%s'
organizationname='%s'
user='%s'
password='%s'
warehouse='SNOWFLAKE'
role='ACCOUNTADMIN'
clientip='1.2.3.4'
protocol='https'
port=443
oktaurl='https://example.com'
clienttimeout=10
jwtclienttimeout=20
logintimeout=30
requesttimeout=40
jwtexpiretimeout=50
externalbrowsertimeout=60
maxretrycount=1
authenticator='snowflake'
insecuremode=true
ocspfailopen=true
token='token'
keepsessionalive=true
disabletelemetry=true
validatedefaultparameters=true
clientrequestmfatoken=true
clientstoretemporarycredential=true
tracing='info'
tmpdirpath='.'
disablequerycontextcache=true
includeretryreason=true
disableconsolelogin=true

[complete_fields.params]
foo = 'bar'

`, accountName, orgName, user, password))
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
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.97.0"),
				Config:            providerConfigWithClientStoreTemporaryCredential(testprofiles.Default, `true`),
			},
			{
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
			{
				Config:      providerConfigWithLogLevel(testprofiles.Default, "invalid"),
				ExpectError: regexp.MustCompile(`invalid driver log level: INVALID`),
			},
			{
				Config: providerConfig("non-existing"),
				// .* is used to match the error message regarding of the home user location
				ExpectError: regexp.MustCompile(`could not retrieve profile config: profile "non-existing" not found in file .*.snowflake/config`),
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

func providerConfigWithClientStoreTemporaryCredential(profile, clientStoreTemporaryCredential string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	client_store_temporary_credential    = %[2]s
}
`, profile, clientStoreTemporaryCredential) + datasourceConfig()
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

func providerConfigWithLogLevel(profile, logLevel string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	driver_tracing    = "%[2]s"
}
`, profile, logLevel) + datasourceConfig()
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

func providerConfigWithNewAccountId(profile, orgName, accountName string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	account_name    = "%[2]s"
	organization_name    = "%[3]s"
}
`, profile, accountName, orgName) + datasourceConfig()
}

func providerConfigComplete(profile, user, password, orgName, accountName string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	user = "%[2]s"
	password = "%[3]s"
	organization_name = "%[4]s"
	account_name = "%[5]s"
	warehouse    = "SNOWFLAKE"
}
`, profile, user, password, orgName, accountName) + datasourceConfig()
}

func datasourceConfig() string {
	return fmt.Sprintf(`
data snowflake_database "t" {
	name = "%s"
}`, acc.TestDatabaseName)
}

func providerConfigAllFields(profile, orgName, accountName, user, password string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	organization_name = "%[2]s"
	account_name = "%[3]s"
	user = "%[4]s"
	password = "%[5]s"
	warehouse = "SNOWFLAKE"
	protocol = "https"
	port = "443"
	role = "ACCOUNTADMIN"
	validate_default_parameters = true
	client_ip = "3.3.3.3"
	authenticator = "snowflake"
	okta_url = "https://example-tf.com"
	login_timeout = 101
	request_timeout = 201
	jwt_expire_timeout = 301
	client_timeout = 401
	jwt_client_timeout = 501
	external_browser_timeout = 601
	insecure_mode = true
	ocsp_fail_open = true
	keep_session_alive = true
	disable_telemetry = true
	client_request_mfa_token = true
	client_store_temporary_credential = true
	disable_query_context_cache = true
	include_retry_reason = true
	max_retry_count = 3
	driver_tracing = "info"
	tmp_directory_path = "../../"
	disable_console_login = true
	params = {
		foo = "piyo"
	}
}
`, profile, orgName, accountName, user, password) + datasourceConfig()
}
