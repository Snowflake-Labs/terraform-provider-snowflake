package provider_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	incorrectConfig := acc.TestClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

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
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// make sure that we succeed for the correct profile
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt).WithProfile(tmpServiceUserConfig.Profile).WithUserId(ids.NonExistingAccountObjectIdentifier), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and key in provider's config should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(incorrectConfig.Profile).
					WithUserId(tmpServiceUser.UserId).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error (profile authenticator is set to JWT and that's why the error is about incorrect token)
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, ids.NonExistingAccountObjectIdentifier.Name())
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and private key in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(tmpServiceUserConfig.Profile).
					WithUserId(ids.NonExistingAccountObjectIdentifier).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// there is no config (by setting the dir to something different from .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(testprofiles.Default).
					WithUserId(tmpServiceUser.UserId).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.AccountName, tmpServiceUser.AccountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, tmpServiceUser.AccountId.OrganizationName())
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt).WithProfile(testprofiles.Default).WithUserId(ids.NonExistingAccountObjectIdentifier), datasourceModel()),
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
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_configureClientOnceSwitching(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	incorrectConfig := acc.TestClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

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
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// in this step we simulate the situation when we want to use client configured once, but it was faulty last time
			{
				PreConfig: func() {
					t.Setenv(string(testenvs.ConfigureClientOnce), "true")
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
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

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
	})

	oktaUrl, err := url.Parse("https://example.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("1.2.3.4"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
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

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullInvalidTomlConfigForServiceUser(t, profile)
	})

	oktaUrlFromEnv, err := url.Parse("https://example-env.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Account)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.AccountName, tmpServiceUser.AccountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, tmpServiceUser.AccountId.OrganizationName())
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.Warehouse, tmpServiceUser.WarehouseId.Name())
					t.Setenv(snowflakeenvs.Protocol, "https")
					t.Setenv(snowflakeenvs.Port, "443")
					// do not set token - it should be propagated from TOML
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
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
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("2.2.2.2"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
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

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullInvalidTomlConfigForServiceUser(t, profile)
	})

	oktaUrlFromTf, err := url.Parse("https://example-tf.com")
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
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
				Config: config.FromModels(t, providermodel.SnowflakeProvider().AllFields(tmpServiceUserConfig, tmpServiceUser), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := acc.TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("3.3.3.3"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
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

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	nonExisting := "NON-EXISTENT"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithRole(nonExisting).WithValidateDefaultParameters("true"), datasourceModel()),
				ExpectError: regexp.MustCompile("Role 'NON-EXISTENT' specified in the connect string does not exist or not authorized."),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithWarehouse(nonExisting).WithValidateDefaultParameters("true"), datasourceModel()),
				ExpectError: regexp.MustCompile("The requested warehouse does not exist or not authorized."),
			},
			// check that using a non-existing warehouse with disabled verification succeeds
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithWarehouse(nonExisting).WithValidateDefaultParameters("false"), datasourceModel()),
			},
		},
	})
}

// prove we can use tri-value booleans, similarly to the ones in resources
func TestAcc_Provider_triValueBoolean(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

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
				Config:            config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientStoreTemporaryCredentialBool(true), datasourceModel()),
			},
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientStoreTemporaryCredentialBool(true), datasourceModel()),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientStoreTemporaryCredential("true"), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_sessionParameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

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
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				// TODO [SNOW-1348325]: Use parameter data source with `IN SESSION` filtering.
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithParamsValue(
					tfconfig.ObjectVariable(
						map[string]tfconfig.Variable{
							"statement_timeout_in_seconds": tfconfig.IntegerVariable(31337),
						},
					),
				)) + unsafeExecuteShowSessionParameter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_unsafe_execute.t", "query_results.#", "1"),
					resource.TestCheckResourceAttr("snowflake_unsafe_execute.t", "query_results.0.value", "31337"),
				),
			},
		},
	})
}

func TestAcc_Provider_JwtAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	tmpIncorrectServiceUserConfig := acc.TestClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)
	tmpServiceUserWithEncryptedKeyConfig := acc.TestClient().TempTomlConfigForServiceUserWithEncryptedKey(t, tmpServiceUser)
	tmpIncorrectServiceUserWithEncryptedKeyConfig := acc.TestClient().TempIncorrectTomlConfigForServiceUserWithEncryptedKey(t, tmpServiceUser)

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
			// authenticate with incorrect private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpIncorrectServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpIncorrectServiceUserConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// authenticate with unencrypted private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
			},
			// authenticate with unencrypted private key with a legacy authenticator value
			// solves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2983
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwtLegacy), datasourceModel()),
			},
			// check encrypted private key with incorrect password
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpIncorrectServiceUserWithEncryptedKeyConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpIncorrectServiceUserWithEncryptedKeyConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				ExpectError: regexp.MustCompile("pkcs8: incorrect password"),
			},
			// authenticate with encrypted private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserWithEncryptedKeyConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserWithEncryptedKeyConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_SnowflakeAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpLegacyServiceUser := acc.TestClient().SetUpTemporaryLegacyServiceUser(t)
	tmpLegacyServiceUserConfig := acc.TestClient().TempTomlConfigForLegacyServiceUser(t, tmpLegacyServiceUser)
	incorrectLegacyServiceUserConfig := acc.TestClient().TempIncorrectTomlConfigForLegacyServiceUser(t, tmpLegacyServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectLegacyServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectLegacyServiceUserConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpLegacyServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpLegacyServiceUserConfig.Profile), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_invalidConfigurations(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := acc.TestClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := acc.TestClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() {
			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientIp("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("expected client_ip to contain a valid IP"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithProtocol("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid protocol: invalid"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithPort(123456789), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "port" to be a valid port number or 0, got: 123456789`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithAuthenticator("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid authenticator type: invalid"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithOktaUrl("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "okta_url" to have a host, got invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithLoginTimeout(-1), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected login_timeout to be at least \(0\), got -1`),
			},
			{
				Config: config.FromModels(
					t,
					providermodel.SnowflakeProvider().
						WithProfile(tmpServiceUserConfig.Profile).
						WithTokenAccessorValue(
							tfconfig.ObjectVariable(
								map[string]tfconfig.Variable{
									"token_endpoint": tfconfig.StringVariable("invalid"),
									"refresh_token":  tfconfig.StringVariable("refresh_token"),
									"client_id":      tfconfig.StringVariable("client_id"),
									"client_secret":  tfconfig.StringVariable("client_secret"),
									"redirect_uri":   tfconfig.StringVariable("redirect_uri"),
								},
							),
						),
					datasourceModel(),
				),
				ExpectError: regexp.MustCompile(`expected "token_endpoint" to have a host, got invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithDriverTracing("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`invalid driver log level: invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile("non-existing"), datasourceModel()),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`profile "non-existing" not found in file %s`, tmpServiceUserConfig.Path)),
			},
		},
	})
}

func datasourceModel() config.DatasourceModel {
	return datasourcemodel.Database("t", acc.TestDatabaseName)
}

func unsafeExecuteShowSessionParameter() string {
	return `
resource snowflake_unsafe_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
