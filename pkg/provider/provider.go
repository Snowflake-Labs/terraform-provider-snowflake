package provider

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// Provider returns a Terraform Provider using configuration from https://pkg.go.dev/github.com/snowflakedb/gosnowflake#Config
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Description: "Specifies your Snowflake account identifier assigned, by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html). Can also be sourced from the `SNOWFLAKE_ACCOUNT` environment variable. Required unless using `profile`.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"user": {
				Type:        schema.TypeString,
				Description: "Username. Can also be sourced from the `SNOWFLAKE_USER` environment variable. Required unless using `profile`.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for username+password authentication. Can also be sourced from the `SNOWFLAKE_USERNAME` environment variable. Required unless using `profile`.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USERNAME", nil),
				Deprecated:  "Use `user` instead of `username`",
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`. Can also be sourced from the `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			// todo: add database and schema once unqualified identifiers are supported
			"warehouse": {
				Type:        schema.TypeString,
				Description: "Specifies the virtual warehouse to use by default for queries, loading, etc. in the client session. Can also be sourced from the `SNOWFLAKE_WAREHOUSE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_WAREHOUSE", nil),
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Specifies the role to use by default for accessing Snowflake objects in the client session. Can also be sourced from the `SNOWFLAKE_ROLE` environment variable. .",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"validate_default_parameters": {
				Type:        schema.TypeBool,
				Description: "True by default. If false, disables the validation checks for Database, Schema, Warehouse and Role at the time a connection is established. Can also be sourced from the `SNOWFLAKE_VALIDATE_DEFAULT_PARAMETERS` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_VALIDATE_DEFAULT_PARAMETERS", nil),
			},
			"params": {
				Type:        schema.TypeMap,
				Description: "Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
			},
			"client_ip": {
				Type:        schema.TypeString,
				Description: "IP address for network checks. Can also be sourced from the `SNOWFLAKE_CLIENT_IP` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_CLIENT_IP", nil),
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Either http or https, defaults to https. Can also be sourced from the `SNOWFLAKE_PROTOCOL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROTOCOL", nil),
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					switch val.(string) {
					case "http", "https":
						return nil, nil
					default:
						errs := append(errs, fmt.Errorf("%q must be one of http or https", key))
						return warns, errs
					}
				},
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Supports passing in a custom host value to the snowflake go driver for use with privatelink. Can also be sourced from the `SNOWFLAKE_HOST` environment variable. ",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Support custom port values to snowflake go driver for use with privatelink. Can also be sourced from the `SNOWFLAKE_PORT` environment variable. ",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PORT", nil),
			},
			"authenticator": {
				Type:        schema.TypeString,
				Description: "Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid values include: Snowflake, OAuth, ExternalBrowser, Okta, JWT, TokenAccessor, UsernamePasswordMFA. Can also be sourced from the `SNOWFLAKE_AUTHENTICATOR` environment variable. It has to be set explicitly to JWT for private key authentication.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_AUTHENTICATOR", nil),
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					switch val.(string) {
					case "Snowflake", "OAuth", "ExternalBrowser", "Okta", "JWT", "TokenAccessor", "UsernamePasswordMFA":
						return nil, nil
					default:
						errs := append(errs, fmt.Errorf("%q must be one of Snowflake, OAuth, ExternalBrowser, Okta, JWT, TokenAccessor or UsernamePasswordMFA", key))
						return warns, errs
					}
				},
			},
			"passcode": {
				Type:          schema.TypeString,
				Description:   "Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login. Can also be sourced from the `SNOWFLAKE_PASSCODE` environment variable. ",
				Optional:      true,
				ConflictsWith: []string{"passcode_in_password"},
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSCODE", nil),
			},
			"passcode_in_password": {
				Type:          schema.TypeBool,
				Description:   "False by default. Set to true if the MFA passcode is embedded in the login password. Appends the MFA passcode to the end of the password. Can also be sourced from the `SNOWFLAKE_PASSCODE_IN_PASSWORD` environment variable. ",
				Optional:      true,
				ConflictsWith: []string{"passcode"},
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSCODE_IN_PASSWORD", nil),
			},
			"okta_url": {
				Type:        schema.TypeString,
				Description: "The URL of the Okta server. e.g. https://example.okta.com. Can also be sourced from the `SNOWFLAKE_OKTA_URL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OKTA_URL", nil),
			},
			"login_timeout": {
				Type:        schema.TypeInt,
				Description: "Login retry timeout EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_LOGIN_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_LOGIN_TIMEOUT", nil),
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Description: "request retry timeout EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_REQUEST_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REQUEST_TIMEOUT", nil),
			},
			"jwt_expire_timeout": {
				Type:        schema.TypeInt,
				Description: "JWT expire after timeout in seconds. Can also be sourced from the `SNOWFLAKE_JWT_EXPIRE_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_JWT_EXPIRE_TIMEOUT", nil),
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the client to complete the authentication. Default is 900 seconds. Can also be sourced from the `SNOWFLAKE_CLIENT_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_CLIENT_TIMEOUT", nil),
			},
			"jwt_client_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the JWT client to complete the authentication. Default is 10 seconds. Can also be sourced from the `SNOWFLAKE_JWT_CLIENT_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_JWT_CLIENT_TIMEOUT", nil),
			},
			"external_browser_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the external browser to complete the authentication. Default is 120 seconds. Can also be sourced from the `SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT", nil),
			},
			"insecure_mode": {
				Type:        schema.TypeBool,
				Description: "If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only. Can also be sourced from the `SNOWFLAKE_INSECURE_MODE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_INSECURE_MODE", nil),
			},
			"ocsp_fail_open": {
				Type:        schema.TypeBool,
				Description: "True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default. Can also be sourced from the `SNOWFLAKE_OCSP_FAIL_OPEN` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OCSP_FAIL_OPEN", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Token to use for OAuth and other forms of token based auth. Can also be sourced from the `SNOWFLAKE_TOKEN` environment variable.",
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN", nil),
			},
			"token_accessor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_endpoint": {
							Type:        schema.TypeString,
							Description: "The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT` environment variable.",
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT", nil),
						},
						"refresh_token": {
							Type:        schema.TypeString,
							Description: "The refresh token for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN` environment variable.",
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN", nil),
						},
						"client_id": {
							Type:        schema.TypeString,
							Description: "The client ID for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID` environment variable.",
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID", nil),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Description: "The client secret for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET` environment variable.",
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET", nil),
						},
						"redirect_uri": {
							Type:        schema.TypeString,
							Description: "The redirect URI for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI` environment variable.",
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI", nil),
						},
					},
				},
			},
			"keep_session_alive": {
				Type:        schema.TypeBool,
				Description: "Enables the session to persist even after the connection is closed. Can also be sourced from the `SNOWFLAKE_KEEP_SESSION_ALIVE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_KEEP_SESSION_ALIVE", nil),
			},
			"private_key": {
				Type:          schema.TypeString,
				Description:   "Private Key for username+private-key auth. Cannot be used with `browser_auth` or `password`. Can also be sourced from `SNOWFLAKE_PRIVATE_KEY` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY", nil),
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key_path", "oauth_refresh_token"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   "Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc. Can also be sourced from `SNOWFLAKE_PRIVATE_KEY_PASSPHRASE` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE", nil),
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "oauth_refresh_token"},
			},
			"disable_telemetry": {
				Type:        schema.TypeBool,
				Description: "Indicates whether to disable telemetry. Can also be sourced from the `SNOWFLAKE_DISABLE_TELEMETRY` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_DISABLE_TELEMETRY", nil),
			},
			"client_request_mfa_token": {
				Type:        schema.TypeBool,
				Description: "When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_REQUEST_MFA_TOKEN` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_CLIENT_REQUEST_MFA_TOKEN", nil),
			},
			"client_store_temporary_credential": {
				Type:        schema.TypeBool,
				Description: "When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_STORE_TEMPORARY_CREDENTIAL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_CLIENT_STORE_TEMPORARY_CREDENTIAL", nil),
			},
			"disable_query_context_cache": {
				Type:        schema.TypeBool,
				Description: "Should HTAP query context cache be disabled. Can also be sourced from the `SNOWFLAKE_DISABLE_QUERY_CONTEXT_CACHE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_DISABLE_QUERY_CONTEXT_CACHE", nil),
			},
			/*
				Feature not yet released as of latest gosnowflake release
				https://github.com/snowflakedb/gosnowflake/blob/master/dsn.go#L103
				"include_retry_reason": {
					Type:        schema.TypeBool,
					Description: "Should retried request contain retry reason. Can also be sourced from the `SNOWFLAKE_INCLUDE_RETRY_REASON` environment variable.",
					Optional:    true,
				},
			*/
			"profile": {
				Type:        schema.TypeString,
				Description: "Sets the profile to read from ~/.snowflake/config file. Can also be sourced from the `SNOWFLAKE_PROFILE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROFILE", "default"),
			},
			// Deprecated attributes
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
		ResourcesMap:       getResources(),
		DataSourcesMap:     getDataSources(),
		ConfigureFunc:      ConfigureProvider,
		ProviderMetaSchema: map[string]*schema.Schema{},
	}
}

func GetGrantResources() resources.TerraformGrantResources {
	grants := resources.TerraformGrantResources{
		"snowflake_account_grant":           resources.AccountGrant(),
		"snowflake_database_grant":          resources.DatabaseGrant(),
		"snowflake_external_table_grant":    resources.ExternalTableGrant(),
		"snowflake_failover_group_grant":    resources.FailoverGroupGrant(),
		"snowflake_file_format_grant":       resources.FileFormatGrant(),
		"snowflake_function_grant":          resources.FunctionGrant(),
		"snowflake_integration_grant":       resources.IntegrationGrant(),
		"snowflake_masking_policy_grant":    resources.MaskingPolicyGrant(),
		"snowflake_materialized_view_grant": resources.MaterializedViewGrant(),
		"snowflake_pipe_grant":              resources.PipeGrant(),
		"snowflake_procedure_grant":         resources.ProcedureGrant(),
		"snowflake_resource_monitor_grant":  resources.ResourceMonitorGrant(),
		"snowflake_row_access_policy_grant": resources.RowAccessPolicyGrant(),
		"snowflake_schema_grant":            resources.SchemaGrant(),
		"snowflake_sequence_grant":          resources.SequenceGrant(),
		"snowflake_stage_grant":             resources.StageGrant(),
		"snowflake_stream_grant":            resources.StreamGrant(),
		"snowflake_table_grant":             resources.TableGrant(),
		"snowflake_tag_grant":               resources.TagGrant(),
		"snowflake_task_grant":              resources.TaskGrant(),
		"snowflake_user_grant":              resources.UserGrant(),
		"snowflake_view_grant":              resources.ViewGrant(),
		"snowflake_warehouse_grant":         resources.WarehouseGrant(),
	}
	return grants
}

func getResources() map[string]*schema.Resource {
	// NOTE(): do not add grant resources here
	others := map[string]*schema.Resource{
		"snowflake_account":                                 resources.Account(),
		"snowflake_account_password_policy_attachment":      resources.AccountPasswordPolicyAttachment(),
		"snowflake_account_parameter":                       resources.AccountParameter(),
		"snowflake_alert":                                   resources.Alert(),
		"snowflake_api_integration":                         resources.APIIntegration(),
		"snowflake_database":                                resources.Database(),
		"snowflake_database_role":                           resources.DatabaseRole(),
		"snowflake_dynamic_table":                           resources.DynamicTable(),
		"snowflake_email_notification_integration":          resources.EmailNotificationIntegration(),
		"snowflake_external_function":                       resources.ExternalFunction(),
		"snowflake_external_oauth_integration":              resources.ExternalOauthIntegration(),
		"snowflake_external_table":                          resources.ExternalTable(),
		"snowflake_failover_group":                          resources.FailoverGroup(),
		"snowflake_file_format":                             resources.FileFormat(),
		"snowflake_function":                                resources.Function(),
		"snowflake_grant_account_role":                      resources.GrantAccountRole(),
		"snowflake_grant_application_role":                  resources.GrantApplicationRole(),
		"snowflake_grant_database_role":                     resources.GrantDatabaseRole(),
		"snowflake_grant_ownership":                         resources.GrantOwnership(),
		"snowflake_grant_privileges_to_role":                resources.GrantPrivilegesToRole(),
		"snowflake_grant_privileges_to_account_role":        resources.GrantPrivilegesToAccountRole(),
		"snowflake_grant_privileges_to_database_role":       resources.GrantPrivilegesToDatabaseRole(),
		"snowflake_grant_privileges_to_share":               resources.GrantPrivilegesToShare(),
		"snowflake_managed_account":                         resources.ManagedAccount(),
		"snowflake_masking_policy":                          resources.MaskingPolicy(),
		"snowflake_materialized_view":                       resources.MaterializedView(),
		"snowflake_network_policy":                          resources.NetworkPolicy(),
		"snowflake_network_policy_attachment":               resources.NetworkPolicyAttachment(),
		"snowflake_notification_integration":                resources.NotificationIntegration(),
		"snowflake_oauth_integration":                       resources.OAuthIntegration(),
		"snowflake_object_parameter":                        resources.ObjectParameter(),
		"snowflake_password_policy":                         resources.PasswordPolicy(),
		"snowflake_pipe":                                    resources.Pipe(),
		"snowflake_procedure":                               resources.Procedure(),
		"snowflake_resource_monitor":                        resources.ResourceMonitor(),
		"snowflake_role":                                    resources.Role(),
		"snowflake_role_grants":                             resources.RoleGrants(),
		"snowflake_role_ownership_grant":                    resources.RoleOwnershipGrant(),
		"snowflake_row_access_policy":                       resources.RowAccessPolicy(),
		"snowflake_saml_integration":                        resources.SAMLIntegration(),
		"snowflake_schema":                                  resources.Schema(),
		"snowflake_scim_integration":                        resources.SCIMIntegration(),
		"snowflake_sequence":                                resources.Sequence(),
		"snowflake_session_parameter":                       resources.SessionParameter(),
		"snowflake_share":                                   resources.Share(),
		"snowflake_stage":                                   resources.Stage(),
		"snowflake_storage_integration":                     resources.StorageIntegration(),
		"snowflake_stream":                                  resources.Stream(),
		"snowflake_table":                                   resources.Table(),
		"snowflake_table_column_masking_policy_application": resources.TableColumnMaskingPolicyApplication(),
		"snowflake_table_constraint":                        resources.TableConstraint(),
		"snowflake_tag":                                     resources.Tag(),
		"snowflake_tag_association":                         resources.TagAssociation(),
		"snowflake_tag_masking_policy_association":          resources.TagMaskingPolicyAssociation(),
		"snowflake_task":                                    resources.Task(),
		"snowflake_unsafe_execute":                          resources.UnsafeExecute(),
		"snowflake_user":                                    resources.User(),
		"snowflake_user_ownership_grant":                    resources.UserOwnershipGrant(),
		"snowflake_user_password_policy_attachment":         resources.UserPasswordPolicyAttachment(),
		"snowflake_user_public_keys":                        resources.UserPublicKeys(),
		"snowflake_view":                                    resources.View(),
		"snowflake_warehouse":                               resources.Warehouse(),
	}

	return mergeSchemas(
		others,
		GetGrantResources().GetTfSchemas(),
	)
}

func getDataSources() map[string]*schema.Resource {
	dataSources := map[string]*schema.Resource{
		"snowflake_accounts":                           datasources.Accounts(),
		"snowflake_alerts":                             datasources.Alerts(),
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
		"snowflake_parameters":                         datasources.Parameters(),
		"snowflake_pipes":                              datasources.Pipes(),
		"snowflake_procedures":                         datasources.Procedures(),
		"snowflake_resource_monitors":                  datasources.ResourceMonitors(),
		"snowflake_role":                               datasources.Role(),
		"snowflake_roles":                              datasources.Roles(),
		"snowflake_row_access_policies":                datasources.RowAccessPolicies(),
		"snowflake_schemas":                            datasources.Schemas(),
		"snowflake_sequences":                          datasources.Sequences(),
		"snowflake_shares":                             datasources.Shares(),
		"snowflake_stages":                             datasources.Stages(),
		"snowflake_storage_integrations":               datasources.StorageIntegrations(),
		"snowflake_streams":                            datasources.Streams(),
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

	return dataSources
}

var (
	configuredClient     *sdk.Client
	configureClientError error //nolint:errname
)

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	// hacky way to speed up our acceptance tests
	if os.Getenv("TF_ACC") != "" && os.Getenv("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE") == "true" {
		if configuredClient != nil {
			return &provider.Context{Client: configuredClient}, nil
		}
		if configureClientError != nil {
			return nil, configureClientError
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

	if v, ok := s.GetOk("validate_default_parameters"); ok && v.(bool) {
		config.ValidateDefaultParameters = gosnowflake.ConfigBoolTrue
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
		config.Authenticator = toAuthenticatorType(v.(string))
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
			return nil, fmt.Errorf("could not parse okta_url err = %w", err)
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

	if v, ok := s.GetOk("ocsp_fail_open"); ok && v.(bool) {
		config.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
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
				return nil, fmt.Errorf("could not retrieve access token from refresh token, err = %w", err)
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
		return nil, fmt.Errorf("could not retrieve private key: %w", err)
	}
	if v != nil {
		config.PrivateKey = v
	}

	if v, ok := s.GetOk("disable_telemetry"); ok && v.(bool) {
		config.DisableTelemetry = v.(bool)
	}

	if v, ok := s.GetOk("client_request_mfa_token"); ok && v.(bool) {
		config.ClientRequestMfaToken = gosnowflake.ConfigBoolTrue
	}

	if v, ok := s.GetOk("client_store_temporary_credential"); ok && v.(bool) {
		config.ClientStoreTemporaryCredential = gosnowflake.ConfigBoolTrue
	}

	if v, ok := s.GetOk("disable_query_context_cache"); ok && v.(bool) {
		config.DisableQueryContextCache = v.(bool)
	}

	/*
		Feature not yet released as of latest gosnowflake release
		https://github.com/snowflakedb/gosnowflake/blob/master/dsn.go#L103
		if v, ok := s.GetOk("include_retry_reason"); ok && v.(bool) {
			config.IncludeRetryParameters = v.(bool)
		}
	*/
	if v, ok := s.GetOk("profile"); ok && v.(string) != "" {
		profile := v.(string)
		if profile == "default" {
			defaultConfig := sdk.DefaultConfig()
			config = sdk.MergeConfig(config, defaultConfig)
		} else {
			profileConfig, err := sdk.ProfileConfig(profile)
			if err != nil {
				return "", errors.New("could not retrieve profile config: " + err.Error())
			}
			if profileConfig == nil {
				return "", errors.New("profile with name: " + profile + " not found in config file")
			}
			// merge any credentials found in profile with config
			config = sdk.MergeConfig(config, profileConfig)
		}
	}

	cl, clErr := sdk.NewClient(config)

	// needed for tests verifying different provider setups
	if os.Getenv("TF_ACC") != "" && os.Getenv("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE") == "true" {
		configuredClient = cl
		configureClientError = clErr
	} else {
		configuredClient = nil
		configureClientError = nil
	}

	if clErr != nil {
		return nil, clErr
	}

	return &provider.Context{Client: cl}, nil
}
