package provider

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/snowflakedb/gosnowflake"
)

func init() {
	// useful links:
	// - https://github.com/hashicorp/terraform-plugin-docs/issues/10#issuecomment-767682837
	// - https://github.com/hashicorp/terraform-plugin-docs/issues/156#issuecomment-1600427216
	schema.ResourceDescriptionBuilder = func(r *schema.Resource) string {
		desc := r.Description
		if r.DeprecationMessage != "" {
			deprecationMessage := r.DeprecationMessage
			replacement, path, ok := docs.GetDeprecatedResourceReplacement(deprecationMessage)
			if ok {
				deprecationMessage = strings.ReplaceAll(deprecationMessage, replacement, docs.RelativeLink(replacement, path))
			}
			// <deprecation> tag is a hack to split description into two parts (deprecation/real description) nicely. This tag won't be rendered.
			// Check resources.md.tmpl for usage example.
			desc = fmt.Sprintf("~> **Deprecation** %v <deprecation>\n\n%s", deprecationMessage, r.Description)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider returns a Terraform Provider using configuration. It is based on https://pkg.go.dev/github.com/snowflakedb/gosnowflake#Config.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Specifies your Snowflake account identifier assigned, by Snowflake. The [account locator](https://docs.snowflake.com/en/user-guide/admin-account-identifier#format-2-account-locator-in-a-region) format is not supported. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html). Required unless using `profile`.", snowflakeenvs.Account),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Account, nil),
			},
			"user": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Username. Required unless using `profile`.", snowflakeenvs.User),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.User, nil),
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"password": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`.", snowflakeenvs.Password),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.Password, nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"warehouse": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Specifies the virtual warehouse to use by default for queries, loading, etc. in the client session.", snowflakeenvs.Warehouse),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Warehouse, nil),
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"role": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Specifies the role to use by default for accessing Snowflake objects in the client session.", snowflakeenvs.Role),
				Optional:         true,
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Role, nil),
			},
			"validate_default_parameters": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("True by default. If false, disables the validation checks for Database, Schema, Warehouse and Role at the time a connection is established.", snowflakeenvs.ValidateDefaultParameters),
				Optional:         true,
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ValidateDefaultParameters, provider.BooleanDefault),
			},
			"params": {
				Type:        schema.TypeMap,
				Description: "Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
			},
			"client_ip": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("IP address for network checks.", snowflakeenvs.ClientIp),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientIp, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"protocol": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription(fmt.Sprintf("A protocol used in the connection. Valid options are: %v.", docs.PossibleValuesListed(allProtocols)), snowflakeenvs.Protocol),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Protocol, nil),
				ValidateDiagFunc: validators.NormalizeValidation(toProtocol),
			},
			"host": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Specifies a custom host value used by the driver for privatelink connections.", snowflakeenvs.Host),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Host, nil),
			},
			"port": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("Specifies a custom port value used by the driver for privatelink connections.", snowflakeenvs.Port),
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsPortNumberOrZero),
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Port, nil),
			},
			"authenticator": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid values include: Snowflake, OAuth, ExternalBrowser, Okta, JWT, TokenAccessor, UsernamePasswordMFA. It has to be set explicitly to JWT for private key authentication.", snowflakeenvs.Authenticator),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Authenticator, string(authenticationTypeSnowflake)),
				ValidateDiagFunc: validators.NormalizeValidation(toAuthenticatorType),
			},
			"passcode": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login.", snowflakeenvs.Passcode),
				Optional:      true,
				ConflictsWith: []string{"passcode_in_password"},
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.Passcode, nil),
			},
			"passcode_in_password": {
				Type:          schema.TypeBool,
				Description:   envNameFieldDescription("False by default. Set to true if the MFA passcode is embedded to the configured password.", snowflakeenvs.PasscodeInPassword),
				Optional:      true,
				ConflictsWith: []string{"passcode"},
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PasscodeInPassword, nil),
			},
			"okta_url": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("The URL of the Okta server. e.g. https://example.okta.com.", snowflakeenvs.OktaUrl),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.OktaUrl, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"login_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("Login retry timeout in seconds EXCLUDING network roundtrip and read out http response.", snowflakeenvs.LoginTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.LoginTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("request retry timeout in seconds EXCLUDING network roundtrip and read out http response.", snowflakeenvs.RequestTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.RequestTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"jwt_expire_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("JWT expire after timeout in seconds.", snowflakeenvs.JwtExpireTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.JwtExpireTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"client_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the client to complete the authentication.", snowflakeenvs.ClientTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"jwt_client_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the JWT client to complete the authentication.", snowflakeenvs.JwtClientTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.JwtClientTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"external_browser_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the external browser to complete the authentication.", snowflakeenvs.ExternalBrowserTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ExternalBrowserTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"insecure_mode": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only.", snowflakeenvs.InsecureMode),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.InsecureMode, nil),
			},
			"ocsp_fail_open": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default.", snowflakeenvs.OcspFailOpen),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.OcspFailOpen, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"token": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Token to use for OAuth and other forms of token based auth.", snowflakeenvs.Token),
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Token, nil),
			},
			"token_accessor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_endpoint": {
							Type:             schema.TypeString,
							Description:      envNameFieldDescription("The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorTokenEndpoint),
							Required:         true,
							Sensitive:        true,
							DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorTokenEndpoint, nil),
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"refresh_token": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The refresh token for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRefreshToken),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorRefreshToken, nil),
						},
						"client_id": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The client ID for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientId),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorClientId, nil),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The client secret for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientSecret),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorClientSecret, nil),
						},
						"redirect_uri": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The redirect URI for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRedirectUri),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorRedirectUri, nil),
						},
					},
				},
			},
			"keep_session_alive": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Enables the session to persist even after the connection is closed.", snowflakeenvs.KeepSessionAlive),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.KeepSessionAlive, nil),
			},
			"private_key": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Private Key for username+private-key auth. Cannot be used with `browser_auth` or `password`.", snowflakeenvs.PrivateKey),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PrivateKey, nil),
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key_path", "oauth_refresh_token"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc.", snowflakeenvs.PrivateKeyPassphrase),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PrivateKeyPassphrase, nil),
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "oauth_refresh_token"},
			},
			"disable_telemetry": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Disables telemetry in the driver.", snowflakeenvs.DisableTelemetry),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.DisableTelemetry, nil),
			},
			"client_request_mfa_token": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux.", snowflakeenvs.ClientRequestMfaToken),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientRequestMfaToken, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"client_store_temporary_credential": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux.", snowflakeenvs.ClientStoreTemporaryCredential),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientStoreTemporaryCredential, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"disable_query_context_cache": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Disables HTAP query context cache in the driver.", snowflakeenvs.DisableQueryContextCache),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.DisableQueryContextCache, nil),
			},
			"profile": {
				Type: schema.TypeString,
				// TODO(SNOW-1754364): Note that a default file path is already filled on sdk side.
				Description: envNameFieldDescription("Sets the profile to read from ~/.snowflake/config file.", snowflakeenvs.Profile),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Profile, "default"),
			},
			// Deprecated attributes
			"username": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Username for username+password authentication. Required unless using `profile`.", snowflakeenvs.Username),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Username, nil),
				Deprecated:  "Use `user` instead of `username`",
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Snowflake region, such as \"eu-central-1\", with this parameter. However, since this parameter is deprecated, it is best to specify the region as part of the account parameter. For details, see the description of the account parameter. [Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use.  Required if using the [legacy format for the `account` identifier](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#format-2-legacy-account-locator-in-a-region) in the form of `<cloud_region_id>.<cloud>`. Can also be sourced from the `SNOWFLAKE_REGION` environment variable. ",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", nil),
				Deprecated:  "Specify the region as part of the account parameter",
			},
			"session_params": {
				Type:        schema.TypeMap,
				Description: "Sets session parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
				Deprecated:  "Use `params` instead",
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Generating the token is left to other tools. Cannot be used with `browser_auth`, `private_key_path`, `oauth_refresh_token` or `password`. Can also be sourced from `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
				Deprecated:    "Use `token` instead",
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Setup and generation of the token is left to other tools. Should be used in conjunction with `oauth_client_id`, `oauth_client_secret`, `oauth_endpoint`, `oauth_redirect_url`. Cannot be used with `browser_auth`, `private_key_path`, `oauth_access_token` or `password`. Can also be sourced from `SNOWFLAKE_OAUTH_REFRESH_TOKEN` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.refresh_token` instead",
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.client_id` instead",
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.client_secret` instead",
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_ENDPOINT` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.token_endpoint` instead",
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_REDIRECT_URL` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
				Deprecated:    "Use `token_accessor.0.redirect_uri` instead",
			},
			"browser_auth": {
				Type:        schema.TypeBool,
				Description: "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_USE_BROWSER_AUTH` environment variable.",
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Deprecated:  "Use `authenticator` instead",
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Description:   "Path to a private key for using keypair authentication. Cannot be used with `browser_auth`, `oauth_access_token` or `password`. Can also be sourced from `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key"},
				Deprecated:    "use the [file Function](https://developer.hashicorp.com/terraform/language/functions/file) instead",
			},
		},
		ResourcesMap:         getResources(),
		DataSourcesMap:       getDataSources(),
		ConfigureContextFunc: ConfigureProvider,
		ProviderMetaSchema:   map[string]*schema.Schema{},
	}
}

func getResources() map[string]*schema.Resource {
	resourceList := map[string]*schema.Resource{
		"snowflake_account": resources.Account(),
		"snowflake_account_authentication_policy_attachment":                     resources.AccountAuthenticationPolicyAttachment(),
		"snowflake_account_role":                                                 resources.AccountRole(),
		"snowflake_account_password_policy_attachment":                           resources.AccountPasswordPolicyAttachment(),
		"snowflake_account_parameter":                                            resources.AccountParameter(),
		"snowflake_alert":                                                        resources.Alert(),
		"snowflake_api_authentication_integration_with_authorization_code_grant": resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant(),
		"snowflake_api_authentication_integration_with_client_credentials":       resources.ApiAuthenticationIntegrationWithClientCredentials(),
		"snowflake_api_authentication_integration_with_jwt_bearer":               resources.ApiAuthenticationIntegrationWithJwtBearer(),
		"snowflake_api_integration":                                              resources.APIIntegration(),
		"snowflake_authentication_policy":                                        resources.AuthenticationPolicy(),
		"snowflake_cortex_search_service":                                        resources.CortexSearchService(),
		"snowflake_database_old":                                                 resources.DatabaseOld(),
		"snowflake_database":                                                     resources.Database(),
		"snowflake_database_role":                                                resources.DatabaseRole(),
		"snowflake_dynamic_table":                                                resources.DynamicTable(),
		"snowflake_email_notification_integration":                               resources.EmailNotificationIntegration(),
		"snowflake_external_function":                                            resources.ExternalFunction(),
		"snowflake_external_oauth_integration":                                   resources.ExternalOauthIntegration(),
		"snowflake_external_table":                                               resources.ExternalTable(),
		"snowflake_external_volume":                                              resources.ExternalVolume(),
		"snowflake_failover_group":                                               resources.FailoverGroup(),
		"snowflake_file_format":                                                  resources.FileFormat(),
		"snowflake_function":                                                     resources.Function(),
		"snowflake_grant_account_role":                                           resources.GrantAccountRole(),
		"snowflake_grant_application_role":                                       resources.GrantApplicationRole(),
		"snowflake_grant_database_role":                                          resources.GrantDatabaseRole(),
		"snowflake_grant_ownership":                                              resources.GrantOwnership(),
		"snowflake_grant_privileges_to_account_role":                             resources.GrantPrivilegesToAccountRole(),
		"snowflake_grant_privileges_to_database_role":                            resources.GrantPrivilegesToDatabaseRole(),
		"snowflake_grant_privileges_to_share":                                    resources.GrantPrivilegesToShare(),
		"snowflake_legacy_service_user":                                          resources.LegacyServiceUser(),
		"snowflake_managed_account":                                              resources.ManagedAccount(),
		"snowflake_masking_policy":                                               resources.MaskingPolicy(),
		"snowflake_materialized_view":                                            resources.MaterializedView(),
		"snowflake_network_policy":                                               resources.NetworkPolicy(),
		"snowflake_network_policy_attachment":                                    resources.NetworkPolicyAttachment(),
		"snowflake_network_rule":                                                 resources.NetworkRule(),
		"snowflake_notification_integration":                                     resources.NotificationIntegration(),
		"snowflake_oauth_integration":                                            resources.OAuthIntegration(),
		"snowflake_oauth_integration_for_partner_applications":                   resources.OauthIntegrationForPartnerApplications(),
		"snowflake_oauth_integration_for_custom_clients":                         resources.OauthIntegrationForCustomClients(),
		"snowflake_object_parameter":                                             resources.ObjectParameter(),
		"snowflake_password_policy":                                              resources.PasswordPolicy(),
		"snowflake_pipe":                                                         resources.Pipe(),
		"snowflake_primary_connection":                                           resources.PrimaryConnection(),
		"snowflake_procedure":                                                    resources.Procedure(),
		"snowflake_resource_monitor":                                             resources.ResourceMonitor(),
		"snowflake_role":                                                         resources.Role(),
		"snowflake_row_access_policy":                                            resources.RowAccessPolicy(),
		"snowflake_saml_integration":                                             resources.SAMLIntegration(),
		"snowflake_saml2_integration":                                            resources.SAML2Integration(),
		"snowflake_schema":                                                       resources.Schema(),
		"snowflake_scim_integration":                                             resources.SCIMIntegration(),
		"snowflake_secondary_connection":                                         resources.SecondaryConnection(),
		"snowflake_secondary_database":                                           resources.SecondaryDatabase(),
		"snowflake_secret_with_authorization_code_grant":                         resources.SecretWithAuthorizationCodeGrant(),
		"snowflake_secret_with_basic_authentication":                             resources.SecretWithBasicAuthentication(),
		"snowflake_secret_with_client_credentials":                               resources.SecretWithClientCredentials(),
		"snowflake_secret_with_generic_string":                                   resources.SecretWithGenericString(),
		"snowflake_sequence":                                                     resources.Sequence(),
		"snowflake_service_user":                                                 resources.ServiceUser(),
		"snowflake_session_parameter":                                            resources.SessionParameter(),
		"snowflake_share":                                                        resources.Share(),
		"snowflake_shared_database":                                              resources.SharedDatabase(),
		"snowflake_stage":                                                        resources.Stage(),
		"snowflake_storage_integration":                                          resources.StorageIntegration(),
		"snowflake_stream":                                                       resources.Stream(),
		"snowflake_stream_on_directory_table":                                    resources.StreamOnDirectoryTable(),
		"snowflake_stream_on_external_table":                                     resources.StreamOnExternalTable(),
		"snowflake_stream_on_table":                                              resources.StreamOnTable(),
		"snowflake_stream_on_view":                                               resources.StreamOnView(),
		"snowflake_streamlit":                                                    resources.Streamlit(),
		"snowflake_table":                                                        resources.Table(),
		"snowflake_table_column_masking_policy_application":                      resources.TableColumnMaskingPolicyApplication(),
		"snowflake_table_constraint":                                             resources.TableConstraint(),
		"snowflake_tag":                                                          resources.Tag(),
		"snowflake_tag_association":                                              resources.TagAssociation(),
		"snowflake_tag_masking_policy_association":                               resources.TagMaskingPolicyAssociation(),
		"snowflake_task":                                                         resources.Task(),
		"snowflake_unsafe_execute":                                               resources.UnsafeExecute(),
		"snowflake_user":                                                         resources.User(),
		"snowflake_user_authentication_policy_attachment":                        resources.UserAuthenticationPolicyAttachment(),
		"snowflake_user_password_policy_attachment":                              resources.UserPasswordPolicyAttachment(),
		"snowflake_user_public_keys":                                             resources.UserPublicKeys(),
		"snowflake_view":                                                         resources.View(),
		"snowflake_warehouse":                                                    resources.Warehouse(),
	}

	if os.Getenv(string(testenvs.EnableObjectRenamingTest)) != "" {
		resourceList["snowflake_object_renaming"] = resources.ObjectRenamingListsAndSets()
	}

	return resourceList
}

func getDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"snowflake_accounts":                           datasources.Accounts(),
		"snowflake_alerts":                             datasources.Alerts(),
		"snowflake_connections":                        datasources.Connections(),
		"snowflake_cortex_search_services":             datasources.CortexSearchServices(),
		"snowflake_current_account":                    datasources.CurrentAccount(),
		"snowflake_current_role":                       datasources.CurrentRole(),
		"snowflake_database":                           datasources.Database(),
		"snowflake_database_role":                      datasources.DatabaseRole(),
		"snowflake_database_roles":                     datasources.DatabaseRoles(),
		"snowflake_databases":                          datasources.Databases(),
		"snowflake_dynamic_tables":                     datasources.DynamicTables(),
		"snowflake_external_functions":                 datasources.ExternalFunctions(),
		"snowflake_external_tables":                    datasources.ExternalTables(),
		"snowflake_failover_groups":                    datasources.FailoverGroups(),
		"snowflake_file_formats":                       datasources.FileFormats(),
		"snowflake_functions":                          datasources.Functions(),
		"snowflake_grants":                             datasources.Grants(),
		"snowflake_masking_policies":                   datasources.MaskingPolicies(),
		"snowflake_materialized_views":                 datasources.MaterializedViews(),
		"snowflake_network_policies":                   datasources.NetworkPolicies(),
		"snowflake_parameters":                         datasources.Parameters(),
		"snowflake_pipes":                              datasources.Pipes(),
		"snowflake_procedures":                         datasources.Procedures(),
		"snowflake_resource_monitors":                  datasources.ResourceMonitors(),
		"snowflake_role":                               datasources.Role(),
		"snowflake_roles":                              datasources.Roles(),
		"snowflake_row_access_policies":                datasources.RowAccessPolicies(),
		"snowflake_schemas":                            datasources.Schemas(),
		"snowflake_secrets":                            datasources.Secrets(),
		"snowflake_security_integrations":              datasources.SecurityIntegrations(),
		"snowflake_sequences":                          datasources.Sequences(),
		"snowflake_shares":                             datasources.Shares(),
		"snowflake_stages":                             datasources.Stages(),
		"snowflake_storage_integrations":               datasources.StorageIntegrations(),
		"snowflake_streams":                            datasources.Streams(),
		"snowflake_streamlits":                         datasources.Streamlits(),
		"snowflake_system_generate_scim_access_token":  datasources.SystemGenerateSCIMAccessToken(),
		"snowflake_system_get_aws_sns_iam_policy":      datasources.SystemGetAWSSNSIAMPolicy(),
		"snowflake_system_get_privatelink_config":      datasources.SystemGetPrivateLinkConfig(),
		"snowflake_system_get_snowflake_platform_info": datasources.SystemGetSnowflakePlatformInfo(),
		"snowflake_tables":                             datasources.Tables(),
		"snowflake_tasks":                              datasources.Tasks(),
		"snowflake_users":                              datasources.Users(),
		"snowflake_views":                              datasources.Views(),
		"snowflake_warehouses":                         datasources.Warehouses(),
	}
}

var (
	configuredClient     *sdk.Client
	configureClientError error //nolint:errname
)

func ConfigureProvider(ctx context.Context, s *schema.ResourceData) (any, diag.Diagnostics) {
	// hacky way to speed up our acceptance tests
	if os.Getenv("TF_ACC") != "" && os.Getenv("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE") == "true" {
		if configuredClient != nil {
			return &provider.Context{Client: configuredClient}, nil
		}
		if configureClientError != nil {
			return nil, diag.FromErr(configureClientError)
		}
	}

	config := &gosnowflake.Config{
		Application: "terraform-provider-snowflake",
	}

	if v, ok := s.GetOk("account"); ok && v.(string) != "" {
		config.Account = v.(string)
	}

	// backwards compatibility until we can remove this
	if v, ok := s.GetOk("username"); ok && v.(string) != "" {
		config.User = v.(string)
	}

	if v, ok := s.GetOk("user"); ok && v.(string) != "" {
		config.User = v.(string)
	}

	if v, ok := s.GetOk("password"); ok && v.(string) != "" {
		config.Password = v.(string)
	}

	if v, ok := s.GetOk("warehouse"); ok && v.(string) != "" {
		config.Warehouse = v.(string)
	}

	if v, ok := s.GetOk("role"); ok && v.(string) != "" {
		config.Role = v.(string)
	}

	if v, ok := s.GetOk("region"); ok && v.(string) != "" {
		config.Region = v.(string)
	}

	if v := s.Get("validate_default_parameters").(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if parsed {
			config.ValidateDefaultParameters = gosnowflake.ConfigBoolTrue
		} else {
			config.ValidateDefaultParameters = gosnowflake.ConfigBoolFalse
		}
	}

	m := make(map[string]interface{})
	if v, ok := s.GetOk("params"); ok {
		m = v.(map[string]interface{})
	}

	// backwards compatibility until we can remove this
	if v, ok := s.GetOk("session_params"); ok {
		m = v.(map[string]interface{})
	}

	params := make(map[string]*string)
	for key, value := range m {
		strValue := value.(string)
		params[key] = &strValue
	}
	config.Params = params

	if v, ok := s.GetOk("client_ip"); ok && v.(string) != "" {
		config.ClientIP = net.ParseIP(v.(string))
	}

	if v, ok := s.GetOk("protocol"); ok && v.(string) != "" {
		config.Protocol = v.(string)
	}

	if v, ok := s.GetOk("host"); ok && v.(string) != "" {
		config.Host = v.(string)
	}

	if v, ok := s.GetOk("port"); ok && v.(int) > 0 {
		config.Port = v.(int)
	}

	// backwards compatibility until we can remove this
	if v, ok := s.GetOk("browser_auth"); ok && v.(bool) {
		config.Authenticator = gosnowflake.AuthTypeExternalBrowser
	}

	if v, ok := s.GetOk("authenticator"); ok && v.(string) != "" {
		authType, err := toAuthenticatorType(v.(string))
		if err != nil {
			return "", diag.FromErr(err)
		}
		config.Authenticator = authType
	}

	if v, ok := s.GetOk("passcode"); ok && v.(string) != "" {
		config.Passcode = v.(string)
	}

	if v, ok := s.GetOk("passcode_in_password"); ok && v.(bool) {
		config.PasscodeInPassword = v.(bool)
	}
	if v, ok := s.GetOk("okta_url"); ok && v.(string) != "" {
		oktaURL, err := url.Parse(v.(string))
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("could not parse okta_url err = %w", err))
		}
		config.OktaURL = oktaURL
	}

	if v, ok := s.GetOk("login_timeout"); ok && v.(int) > 0 {
		config.LoginTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("request_timeout"); ok && v.(int) > 0 {
		config.RequestTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("jwt_expire_timeout"); ok && v.(int) > 0 {
		config.JWTExpireTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("client_timeout"); ok && v.(int) > 0 {
		config.ClientTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("jwt_client_timeout"); ok && v.(int) > 0 {
		config.JWTClientTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("external_browser_timeout"); ok && v.(int) > 0 {
		config.ExternalBrowserTimeout = time.Second * time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("insecure_mode"); ok && v.(bool) {
		config.InsecureMode = v.(bool)
	}

	if v := s.Get("ocsp_fail_open").(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if parsed {
			config.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
		} else {
			config.OCSPFailOpen = gosnowflake.OCSPFailOpenFalse
		}
	}

	if v, ok := s.GetOk("token"); ok && v.(string) != "" {
		config.Token = v.(string)
		config.Authenticator = gosnowflake.AuthTypeOAuth
	}

	if v, ok := s.GetOk("token_accessor"); ok {
		if len(v.([]interface{})) > 0 {
			tokenAccessor := v.([]interface{})[0].(map[string]interface{})
			tokenEndpoint := tokenAccessor["token_endpoint"].(string)
			refreshToken := tokenAccessor["refresh_token"].(string)
			clientID := tokenAccessor["client_id"].(string)
			clientSecret := tokenAccessor["client_secret"].(string)
			redirectURI := tokenAccessor["redirect_uri"].(string)
			accessToken, err := GetAccessTokenWithRefreshToken(tokenEndpoint, clientID, clientSecret, refreshToken, redirectURI)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("could not retrieve access token from refresh token, err = %w", err))
			}
			config.Token = accessToken
			config.Authenticator = gosnowflake.AuthTypeOAuth
		}
	}

	if v, ok := s.GetOk("keep_session_alive"); ok && v.(bool) {
		config.KeepSessionAlive = v.(bool)
	}

	privateKeyPath := s.Get("private_key_path").(string)
	privateKey := s.Get("private_key").(string)
	privateKeyPassphrase := s.Get("private_key_passphrase").(string)
	v, err := getPrivateKey(privateKeyPath, privateKey, privateKeyPassphrase)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("could not retrieve private key: %w", err))
	}
	if v != nil {
		config.PrivateKey = v
	}

	if v, ok := s.GetOk("disable_telemetry"); ok && v.(bool) {
		config.DisableTelemetry = v.(bool)
	}

	if v := s.Get("client_request_mfa_token").(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if parsed {
			config.ClientRequestMfaToken = gosnowflake.ConfigBoolTrue
		} else {
			config.ClientRequestMfaToken = gosnowflake.ConfigBoolFalse
		}
	}

	if v := s.Get("client_store_temporary_credential").(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if parsed {
			config.ClientStoreTemporaryCredential = gosnowflake.ConfigBoolTrue
		} else {
			config.ClientStoreTemporaryCredential = gosnowflake.ConfigBoolFalse
		}
	}

	if v, ok := s.GetOk("disable_query_context_cache"); ok && v.(bool) {
		config.DisableQueryContextCache = v.(bool)
	}

	if v, ok := s.GetOk("profile"); ok && v.(string) != "" {
		profile := v.(string)
		if profile == "default" {
			defaultConfig := sdk.DefaultConfig()
			config = sdk.MergeConfig(config, defaultConfig)
		} else {
			profileConfig, err := sdk.ProfileConfig(profile)
			if err != nil {
				return "", diag.FromErr(errors.New("could not retrieve profile config: " + err.Error()))
			}
			if profileConfig == nil {
				return "", diag.FromErr(errors.New("profile with name: " + profile + " not found in config file"))
			}
			// merge any credentials found in profile with config
			config = sdk.MergeConfig(config, profileConfig)
		}
	}

	client, clientErr := sdk.NewClient(config)

	// needed for tests verifying different provider setups
	if os.Getenv("TF_ACC") != "" && os.Getenv("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE") == "true" {
		configuredClient = client
		configureClientError = clientErr
	} else {
		configuredClient = nil
		configureClientError = nil
	}

	if clientErr != nil {
		return nil, diag.FromErr(clientErr)
	}

	return &provider.Context{Client: client}, nil
}
