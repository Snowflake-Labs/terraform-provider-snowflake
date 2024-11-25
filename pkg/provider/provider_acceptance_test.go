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

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [this PR]: create this user and role and use in these tests?
func setUpLegacyServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, pass string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, userCleanup := acc.TestClient().User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Password: sdk.String(pass),
		Type:     sdk.Pointer(sdk.UserTypeLegacyService),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	acc.TestClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, acc.TestClient().Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	acc.TestClient().Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, acc.TestClient().Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	acc.TestClient().Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}

// TODO [this PR]: merge with the above func
func setUpServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, publicKey string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, userCleanup := acc.TestClient().User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Type:         sdk.Pointer(sdk.UserTypeService),
		RSAPublicKey: sdk.String(publicKey),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	acc.TestClient().Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, acc.TestClient().Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	acc.TestClient().Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, acc.TestClient().Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	acc.TestClient().Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}

func TestAcc_Provider_configHierarchy(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	accountDetailsConfig, err := sdk.ProfileConfig(testprofiles.OnlyAccountDetails)
	require.NoError(t, err)

	accountParts := strings.SplitN(accountDetailsConfig.Account, "-", 2)
	accountId := sdk.NewAccountIdentifier(accountParts[0], accountParts[1])
	warehouseId := acc.TestClient().Ids.SnowflakeWarehouseId()

	privateKey, publicKey, _ := random.GenerateRSAKeyPair(t)
	tmpUserId, tmpRoleId := setUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	profile := random.AlphaN(6)
	toml := helpers.TomlConfigForServiceUser(t, profile, tmpUserId, tmpRoleId, warehouseId, accountId, privateKey)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))

	tomlWithIncorrectCredentials := helpers.TomlIncorrectConfigForServiceUser(t, testprofiles.IncorrectUserAndPassword, accountId)
	configPathWithIncorrect := testhelpers.TestFile(t, random.AlphaN(10), []byte(tomlWithIncorrectCredentials))

	nonExistingUser := "non-existing-user"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// make sure that we fail for incorrect profile
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, configPathWithIncorrect)
				},
				Config:      providerConfig(testprofiles.IncorrectUserAndPassword),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// make sure that we succeed for the correct profile
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, configPath)
				},
				Config: providerConfig(profile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      providerConfigWithUserAndProfile(nonExistingUser, profile),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and key in provider's config should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, configPathWithIncorrect)
				},
				Config: providerConfigWithUserPrivateKeyAndProfile(tmpUserId.Name(), privateKey, tmpRoleId.Name(), testprofiles.IncorrectUserAndPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error (profile authenticator is set to JWT and that's why the error is about incorrect token)
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, nonExistingUser)
					t.Setenv(snowflakeenvs.ConfigPath, configPath)
				},
				Config:      providerConfig(profile),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and private key in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, tmpUserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, privateKey)
					t.Setenv(snowflakeenvs.Role, tmpRoleId.Name())
					t.Setenv(snowflakeenvs.ConfigPath, configPathWithIncorrect)
				},
				Config: providerConfigWithProfileAndAuthenticator(testprofiles.IncorrectUserAndPassword, sdk.AuthenticationTypeJwt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					t.Setenv(snowflakeenvs.ConfigPath, configPath)
				},
				Config:      providerConfigWithUserPrivateKeyAndProfile(nonExistingUser, privateKey, tmpRoleId.Name(), profile),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// there is no config (by setting the dir to something different from .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config:      providerConfigWithUserPrivateKeyAndProfile(tmpUserId.Name(), privateKey, tmpRoleId.Name(), testprofiles.Default),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					t.Setenv(snowflakeenvs.User, tmpUserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, privateKey)
					t.Setenv(snowflakeenvs.AccountName, accountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, accountId.OrganizationName())
					t.Setenv(snowflakeenvs.Role, tmpRoleId.Name())
				},
				Config:      providerConfigWithUser(nonExistingUser, testprofiles.Default),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// make sure the teardown is fine by using a correct env config at the end
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					testenvs.AssertEnvSet(t, snowflakeenvs.PrivateKey)
					testenvs.AssertEnvSet(t, snowflakeenvs.AccountName)
					testenvs.AssertEnvSet(t, snowflakeenvs.OrganizationName)
					testenvs.AssertEnvSet(t, snowflakeenvs.Role)
				},
				Config: providerConfigWithAuthenticator(sdk.AuthenticationTypeJwt),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	accountDetailsConfig, err := sdk.ProfileConfig(testprofiles.OnlyAccountDetails)
	require.NoError(t, err)

	accountParts := strings.SplitN(accountDetailsConfig.Account, "-", 2)
	accountId := sdk.NewAccountIdentifier(accountParts[0], accountParts[1])
	warehouseId := acc.TestClient().Ids.SnowflakeWarehouseId()

	privateKey, publicKey, _ := random.GenerateRSAKeyPair(t)
	tmpUserId, tmpRoleId := setUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	toml := helpers.FullTomlConfigForServiceUser(t, testprofiles.CompleteFields, tmpUserId, tmpRoleId, warehouseId, accountId, privateKey)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))

	oktaUrl, err := url.Parse("https://example.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, configPath)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig(testprofiles.CompleteFields),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					//assert.Equal(t, account, config.Account)
					assert.Equal(t, tmpUserId.Name(), config.User)
					assert.Equal(t, warehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpRoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("1.2.3.4"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					//assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", account), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.Equal(t, false, config.PasscodeInPassword)
					assert.Equal(t, oktaUrl, config.OktaURL)
					assert.Equal(t, 30*time.Second, config.LoginTimeout)
					assert.Equal(t, 40*time.Second, config.RequestTimeout)
					assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 10*time.Second, config.ClientTimeout)
					assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 1, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.Equal(t, true, config.InsecureMode)
					assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.Equal(t, true, config.KeepSessionAlive)
					assert.Equal(t, true, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, ".", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
					assert.Equal(t, true, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("bar"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_envConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	accountDetailsConfig, err := sdk.ProfileConfig(testprofiles.OnlyAccountDetails)
	require.NoError(t, err)

	accountParts := strings.SplitN(accountDetailsConfig.Account, "-", 2)
	accountId := sdk.NewAccountIdentifier(accountParts[0], accountParts[1])
	warehouseId := acc.TestClient().Ids.SnowflakeWarehouseId()

	privateKey, publicKey, _ := random.GenerateRSAKeyPair(t)
	tmpUserId, tmpRoleId := setUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	toml := helpers.FullInvalidTomlConfigForServiceUser(t, testprofiles.CompleteFieldsInvalid)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))

	oktaUrlFromEnv, err := url.Parse("https://example-env.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Account)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, configPath)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.AccountName, accountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, accountId.OrganizationName())
					t.Setenv(snowflakeenvs.User, tmpUserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, privateKey)
					t.Setenv(snowflakeenvs.Warehouse, warehouseId.Name())
					t.Setenv(snowflakeenvs.Protocol, "https")
					t.Setenv(snowflakeenvs.Port, "443")
					// do not set token - it should be propagated from TOML
					t.Setenv(snowflakeenvs.Role, tmpRoleId.Name())
					t.Setenv(snowflakeenvs.Authenticator, "SNOWFLAKE_JWT")
					t.Setenv(snowflakeenvs.ValidateDefaultParameters, "true")
					t.Setenv(snowflakeenvs.ClientIp, "2.2.2.2")
					t.Setenv(snowflakeenvs.Host, "")
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
					t.Setenv(snowflakeenvs.DriverTracing, string(sdk.DriverLogLevelWarning))
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfig(testprofiles.CompleteFieldsInvalid),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					//assert.Equal(t, account, config.Account)
					assert.Equal(t, tmpUserId.Name(), config.User)
					assert.Equal(t, warehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpRoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("2.2.2.2"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					//assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", account), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.Equal(t, false, config.PasscodeInPassword)
					assert.Equal(t, oktaUrlFromEnv, config.OktaURL)
					assert.Equal(t, 100*time.Second, config.LoginTimeout)
					assert.Equal(t, 200*time.Second, config.RequestTimeout)
					assert.Equal(t, 300*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 400*time.Second, config.ClientTimeout)
					assert.Equal(t, 500*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 600*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 2, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.Equal(t, true, config.InsecureMode)
					assert.Equal(t, gosnowflake.OCSPFailOpenFalse, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.Equal(t, true, config.KeepSessionAlive)
					assert.Equal(t, true, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, "../", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.ClientStoreTemporaryCredential)
					assert.Equal(t, true, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("bar"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_tfConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	accountDetailsConfig, err := sdk.ProfileConfig(testprofiles.OnlyAccountDetails)
	require.NoError(t, err)

	accountParts := strings.SplitN(accountDetailsConfig.Account, "-", 2)
	accountId := sdk.NewAccountIdentifier(accountParts[0], accountParts[1])
	warehouseId := acc.TestClient().Ids.SnowflakeWarehouseId()

	privateKey, publicKey, _ := random.GenerateRSAKeyPair(t)
	tmpUserId, tmpRoleId := setUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	toml := helpers.FullInvalidTomlConfigForServiceUser(t, testprofiles.CompleteFieldsInvalid)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))

	oktaUrlFromTf, err := url.Parse("https://example-tf.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, configPath)
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
					t.Setenv(snowflakeenvs.PrivateKey, "invalid")
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
					t.Setenv(snowflakeenvs.DriverTracing, "invalid")
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
				},
				Config: providerConfigAllFields(testprofiles.CompleteFieldsInvalid, tmpUserId, tmpRoleId, warehouseId, accountId, privateKey),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					//assert.Equal(t, account, config.Account)
					assert.Equal(t, tmpUserId.Name(), config.User)
					assert.Equal(t, warehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpRoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("3.3.3.3"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					//assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", account), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.Equal(t, false, config.PasscodeInPassword)
					assert.Equal(t, oktaUrlFromTf, config.OktaURL)
					assert.Equal(t, 101*time.Second, config.LoginTimeout)
					assert.Equal(t, 201*time.Second, config.RequestTimeout)
					assert.Equal(t, 301*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 401*time.Second, config.ClientTimeout)
					assert.Equal(t, 501*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 601*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 3, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.Equal(t, true, config.InsecureMode)
					assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.Equal(t, true, config.KeepSessionAlive)
					assert.Equal(t, true, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, "../../", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
					assert.Equal(t, true, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("piyo"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_useNonExistentDefaultParams(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	pass := random.Password()
	tmpUserId, tmpRoleId := setUpLegacyServiceUserWithAccessToTestDatabaseAndWarehouse(t, pass)

	nonExisting := "NON-EXISTENT"
	warehouse := acc.TestClient().Ids.SnowflakeWarehouseId().Name()

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
				Config:      providerConfigWithExplicitValidation(tmpUserId.Name(), pass, testprofiles.OnlyAccountDetails, nonExisting, warehouse, true),
				ExpectError: regexp.MustCompile("Role 'NON-EXISTENT' specified in the connect string does not exist or not authorized."),
			},
			{
				Config:      providerConfigWithExplicitValidation(tmpUserId.Name(), pass, testprofiles.OnlyAccountDetails, tmpRoleId.Name(), nonExisting, true),
				ExpectError: regexp.MustCompile("The requested warehouse does not exist or not authorized."),
			},
			// check that using a non-existing warehouse with disabled verification succeeds
			{
				Config: providerConfigWithExplicitValidation(tmpUserId.Name(), pass, testprofiles.OnlyAccountDetails, tmpRoleId.Name(), warehouse, false),
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
				Config: providerConfigWithProfileAndAuthenticator(testprofiles.JwtAuth, sdk.AuthenticationTypeJwt),
			},
			// authenticate with unencrypted private key with a legacy authenticator value
			// solves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2983
			{
				Config: providerConfigWithProfileAndAuthenticator(testprofiles.JwtAuth, sdk.AuthenticationTypeJwtLegacy),
			},
			// authenticate with encrypted private key
			{
				Config: providerConfigWithProfileAndAuthenticator(testprofiles.EncryptedJwtAuth, sdk.AuthenticationTypeJwt),
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
				Config: providerConfigWithProfileAndAuthenticator(testprofiles.Default, sdk.AuthenticationTypeSnowflake),
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

func providerConfigWithProfileAndAuthenticator(profile string, authenticator sdk.AuthenticationType) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	authenticator    = "%[2]s"
}
`, profile, authenticator) + datasourceConfig()
}

func emptyProviderConfig() string {
	return `
provider "snowflake" {
}` + datasourceConfig()
}

func providerConfigWithAuthenticator(authenticator sdk.AuthenticationType) string {
	return fmt.Sprintf(`
provider "snowflake" {
	authenticator    = "%[1]s"
}`, authenticator) + datasourceConfig()
}

func providerConfig(profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
}
`, profile) + datasourceConfig()
}

func providerConfigWithExplicitValidation(user, pass, profile, role, warehouse string, validate bool) string {
	return fmt.Sprintf(`
provider "snowflake" {
	user      = "%[1]s"
	password  = "%[2]s"
	profile   = "%[3]s"
	role      = "%[4]s"
	warehouse = "%[5]s"

	validate_default_parameters = "%[6]t"
}
`, user, pass, profile, role, warehouse, validate) + datasourceConfig()
}

func providerConfigWithClientStoreTemporaryCredential(profile, clientStoreTemporaryCredential string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	client_store_temporary_credential    = %[2]s
}
`, profile, clientStoreTemporaryCredential) + datasourceConfig()
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
	authenticator = "SNOWFLAKE_JWT"
	profile = "%[2]s"
}
`, user, profile) + datasourceConfig()
}

func providerConfigWithUserAndProfile(user string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	authenticator = "SNOWFLAKE_JWT"
	user = "%[1]s"
	profile = "%[2]s"
}
`, user, profile) + datasourceConfig()
}

func providerConfigWithUserPrivateKeyAndProfile(user string, key string, role string, profile string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	authenticator = "SNOWFLAKE_JWT"
	user = "%[1]s"
	private_key = <<EOT
%[2]sEOT
	role = "%[3]s"
	profile = "%[4]s"
}
`, user, key, role, profile) + datasourceConfig()
}

func datasourceConfig() string {
	return fmt.Sprintf(`
data snowflake_database "t" {
	name = "%s"
}`, acc.TestDatabaseName)
}

func providerConfigAllFields(profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountId sdk.AccountIdentifier, privateKey string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	organization_name = "%[2]s"
	account_name = "%[3]s"
	user = "%[4]s"
	private_key = <<EOT
%[7]sEOT
	warehouse = "%[5]s"
	protocol = "https"
	port = "443"
	role = "%[6]s"
	validate_default_parameters = true
	client_ip = "3.3.3.3"
	authenticator = "SNOWFLAKE_JWT"
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
	driver_tracing = "warning"
	tmp_directory_path = "../../"
	disable_console_login = true
	params = {
		foo = "piyo"
	}
}
`, profile, accountId.OrganizationName(), accountId.AccountName(), userId.Name(), warehouseId.Name(), roleId.Name(), privateKey) + datasourceConfig()
}

// TODO(SNOW-1348325): Use parameter data source with `IN SESSION` filtering.
func providerWithParamsConfig(profile string, statementTimeoutInSeconds int) string {
	return fmt.Sprintf(`
provider "snowflake" {
    profile = "%[1]s"
    params = {
        statement_timeout_in_seconds = %[2]d
    }
}
`, profile, statementTimeoutInSeconds) + unsafeExecuteShowSessionParameter()
}

func unsafeExecuteShowSessionParameter() string {
	return `
resource snowflake_unsafe_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
