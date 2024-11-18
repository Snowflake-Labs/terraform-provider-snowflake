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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.IncorrectUserAndPassword), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithUser(nonExistingUser), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in provider's config should not be rewritten by a faulty config
			{
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.IncorrectUserAndPassword).WithUser(user).WithPassword(pass), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, nonExistingUser)
				},
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// correct user and password in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, user)
					t.Setenv(snowflakeenvs.Password, pass)
				},
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.IncorrectUserAndPassword), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", acc.TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithUser(nonExistingUser), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// there is no config (by setting the dir to something different than .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithUser(user).WithPassword(pass), datasourceModel()),
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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithUser(nonExistingUser), datasourceModel()),
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider(), datasourceModel()),
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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.IncorrectUserAndPassword), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			// in this step we simulate the situation when we want to use client configured once, but it was faulty last time
			{
				PreConfig: func() {
					t.Setenv(string(testenvs.ConfigureClientOnce), "true")
				},
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider(), datasourceModel()),
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.CompleteFields), datasourceModel()),
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.CompleteFieldsInvalid), datasourceModel()),
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().AllFields(testprofiles.CompleteFieldsInvalid, orgName, accountName, user, pass), datasourceModel()),
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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithRole(nonExisting), datasourceModel()),
				ExpectError: regexp.MustCompile("Role 'NON-EXISTENT' specified in the connect string does not exist or not authorized."),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithWarehouse(nonExisting), datasourceModel()),
				ExpectError: regexp.MustCompile("The requested warehouse does not exist or not authorized."),
			},
			// check that using a non-existing warehouse with disabled verification succeeds
			{
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithWarehouse(nonExisting).WithValidateDefaultParameters("false"), datasourceModel()),
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
				Config:            config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientStoreTemporaryCredential(`true`), datasourceModel()),
			},
			{
				// Use the default TOML config again.
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, "")
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientStoreTemporaryCredential(`true`), datasourceModel()),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientStoreTemporaryCredential(`"true"`), datasourceModel()),
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
				// TODO(SNOW-1348325): Use parameter data source with `IN SESSION` filtering.
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithParamsValue(
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.JwtAuth,).WithAuthenticator(string(sdk.AuthenticationTypeJwt)), datasourceModel()),
			},
			// authenticate with unencrypted private key with a legacy authenticator value
			// solves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2983
			{
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.JwtAuth).WithAuthenticator(string(sdk.AuthenticationTypeJwtLegacy)), datasourceModel()),
			},
			// authenticate with encrypted private key
			{
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.EncryptedJwtAuth).WithAuthenticator(string(sdk.AuthenticationTypeJwt)), datasourceModel()),
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
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithAuthenticator(string(sdk.AuthenticationTypeSnowflake)), datasourceModel()),
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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientIp("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("expected client_ip to contain a valid IP"),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithProtocol("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid protocol: invalid"),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithPort(123456789), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "port" to be a valid port number or 0, got: 123456789`),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithAuthenticator("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid authenticator type: invalid"),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithOktaUrl("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "okta_url" to have a host, got invalid`),
			},
			{
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithLoginTimeout(-1), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected login_timeout to be at least \(0\), got -1`),
			},
			{
				Config: config.ConfigFromModelsPoc(
					t,
					providermodel.SnowflakeProvider().
						WithProfile(testprofiles.Default).
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
				Config:      config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithDriverTracing("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`invalid driver log level: invalid`),
			},
			{
				Config: config.ConfigFromModelsPoc(t, providermodel.SnowflakeProvider().WithProfile("non-existing"), datasourceModel()),
				// .* is used to match the error message regarding of the home user location
				ExpectError: regexp.MustCompile(`profile "non-existing" not found in file .*.snowflake/config`),
			},
		},
	})
}

func datasourceModel() config.ResourceModel {
	return model.Database("t", acc.TestDatabaseName)
}

func unsafeExecuteShowSessionParameter() string {
	return `
resource snowflake_unsafe_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
