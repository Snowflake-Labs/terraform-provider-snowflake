package provider

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Provider is a provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Description: "The name of the Snowflake account. Can also come from the `SNOWFLAKE_ACCOUNT` environment variable. Required unless using profile.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"user": {
				Type:        schema.TypeString,
				Description: "Username. Can come from the `SNOWFLAKE_USER` environment variable. Required unless using profile.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for username+password authentication. Can come from the `SNOWFLAKE_USER` environment variable. Required unless using profile.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
				Deprecated:  "Use `user` instead",
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`. Can be sourced from `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			// todo: add database and schema once unqualified identifiers are supported
			"warehouse": {
				Type:        schema.TypeString,
				Description: "Sets the default warehouse. Optional. Can be sourced from SNOWFLAKE_WAREHOUSE environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_WAREHOUSE", nil),
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Snowflake role to use for operations. If left unset, default role for user will be used. Can be sourced from the `SNOWFLAKE_ROLE` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": {
				Type:        schema.TypeString,
				Description: "[Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use.  Required if using the [legacy format for the `account` identifier](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#format-2-legacy-account-locator-in-a-region) in the form of `<cloud_region_id>.<cloud>`. Can be sourced from the `SNOWFLAKE_REGION` environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
			},
			"validate_default_parameters": {
				Type:        schema.TypeBool,
				Description: "ValidateDefaultParameters disable the validation checks for Database, Schema, Warehouse and Role at the time a connection is established",
				Optional:    true,
			},
			"params": {
				Type:        schema.TypeMap,
				Description: "Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
			},
			"session_params": {
				Type:        schema.TypeMap,
				Description: "Sets session parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
				Deprecated:  "Use `params` instead",
			},
			"client_ip": {
				Type:        schema.TypeString,
				Description: "IP address for network check",
				Optional:    true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Description: "Either http or https. Can be sourced from `SNOWFLAKE_PROTOCOL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROTOCOL", "https"),
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Supports passing in a custom host value to the snowflake go driver for use with privatelink.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Support custom port values to snowflake go driver for use with privatelink. Can be sourced from `SNOWFLAKE_PORT` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PORT", 443),
			},
			"authenticator": {
				Type:        schema.TypeString,
				Description: "Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid values include: Snowflake, OAuth, ExternalBrowser, Okta, JWT, TokenAccessor, UsernamePasswordMFA",
				Optional:    true,
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
				Description:   "Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login.",
				Optional:      true,
				ConflictsWith: []string{"passcode_in_password"},
			},
			"passcode_in_password": {
				Type:          schema.TypeBool,
				Description:   "False by default. Set to true if the MFA passcode is embedded in the login password. Appends the MFA passcode to the end of the password.",
				Optional:      true,
				ConflictsWith: []string{"passcode"},
			},
			"okta_url": {
				Type:        schema.TypeString,
				Description: "The URL of the Okta server. e.g. https://example.okta.com",
				Optional:    true,
			},
			"login_timeout": {
				Type:        schema.TypeInt,
				Description: "Login retry timeout EXCLUDING network roundtrip and read out http response.",
				Optional:    true,
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Description: "request retry timeout EXCLUDING network roundtrip and read out http response",
				Optional:    true,
			},
			"jwt_expire_timeout": {
				Type:        schema.TypeInt,
				Description: "JWT expire after timeout in seconds.",
				Optional:    true,
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the client to complete the authentication. Default is 900 seconds.",
				Optional:    true,
			},
			"jwt_client_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the JWT client to complete the authentication. Default is 10 seconds.",
				Optional:    true,
			},
			"external_browser_timeout": {
				Type:        schema.TypeInt,
				Description: "The timeout in seconds for the external browser to complete the authentication. Default is 120 seconds.",
				Optional:    true,
			},
			"insecure_mode": {
				Type:        schema.TypeBool,
				Description: "If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_INSECURE_MODE", false),
			},
			"oscp_fail_open": {
				Type:        schema.TypeBool,
				Description: "True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default.",
				Optional:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Token to use for OAuth other forms of token based auth",
				Sensitive:   true,
			},
			"token_accessor": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Provider{
					Schema: map[string]*schema.Schema{
						"token_endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
							Sensitive:   true,
							Description: "The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token",
						},
						"refresh_token": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
							Sensitive:   true,
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
							Sensitive:   true,
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
							Sensitive:   true,
						},
						"redirect_uri": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
							Sensitive:   true,
						},
					},
				},
			},
			"keep_session_alive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the session to persist even after the connection is closed",
			},
			"private_key": {
				Type:          schema.TypeString,
				Description:   "Private Key for username+private-key auth. Cannot be used with `browser_auth` or `password`. Can be sourced from `SNOWFLAKE_PRIVATE_KEY` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key_path", "oauth_refresh_token"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   "Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "oauth_refresh_token"},
			},
			"disable_telemetry": {
				Type:        schema.TypeBool,
				Description: "ndicates whether to disable telemetry",
				Optional:    true,
			},
			"client_request_mfa_token": {
				Type:        schema.TypeBool,
				Description: "When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux.",
				Optional:    true,
			},
			"client_store_temporary_credential": {
				Type:        schema.TypeBool,
				Description: "When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux.",
				Optional:    true,
			},
			"disable_query_context_cache": {
				Type:        schema.TypeBool,
				Description: "Should HTAP query context cache be disabled",
				Optional:    true,
			},
			"profile": {
				Type:        schema.TypeString,
				Description: "Sets the profile to read from ~/.snowflake/config file.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROFILE", "default"),
			},
			// Deprecated attributes
			"oauth_access_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Generating the token is left to other tools. Cannot be used with `browser_auth`, `private_key_path`, `oauth_refresh_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
				Deprecated:    "Use `token` instead",
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Setup and generation of the token is left to other tools. Should be used in conjunction with `oauth_client_id`, `oauth_client_secret`, `oauth_endpoint`, `oauth_redirect_url`. Cannot be used with `browser_auth`, `private_key_path`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_REFRESH_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.refresh_token` instead",
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.client_id` instead",
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.client_secret` instead",
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_ENDPOINT` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
				Deprecated:    "Use `token_accessor.0.token_endpoint` instead",
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_REDIRECT_URL` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
				Deprecated:    "Use `token_accessor.0.redirect_uri` instead",
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_USE_BROWSER_AUTH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				Deprecated:    "Use `authenticator` instead",
				ConflictsWith: []string{"password", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Description:   "Path to a private key for using keypair authentication. Cannot be used with `browser_auth`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
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
		"snowflake_grant_privileges_to_role":                resources.GrantPrivilegesToRole(),
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
		"snowflake_user":                                    resources.User(),
		"snowflake_user_ownership_grant":                    resources.UserOwnershipGrant(),
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

type ProviderMeta struct {
	client *sdk.Client
	db     *sql.DB
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	config := &gosnowflake.Config{
		Application: "terraform-provider-snowflake",
	}

	if v, ok := s.GetOk("account"); ok {
		config.Account = v.(string)
	}
	if v, ok := s.GetOk("user"); ok {
		config.User = v.(string)
	}
	// backwards compatibility until we can remove this
	if v, ok := s.GetOk("username"); ok {
		config.User = v.(string)
	}
	if v, ok := s.GetOk("password"); ok {
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

	if v, ok := s.GetOk("protocol"); ok {
		config.Protocol = v.(string)
	}

	if v, ok := s.GetOk("host"); ok && v.(string) != "" {
		config.Host = v.(string)
	}

	if v, ok := s.GetOk("port"); ok {
		config.Port = v.(int)
	}

	// backwards compatibility until we can remove this
	if v, ok := s.GetOk("browser_auth"); ok && v.(bool) {
		config.Authenticator = gosnowflake.AuthTypeExternalBrowser
	}

	if v, ok := s.GetOk("authenticator"); ok && v.(string) != "" {
		authenticator := v.(string)
		switch authenticator {
		case "Snowflake":
			config.Authenticator = gosnowflake.AuthTypeSnowflake
		case "OAuth":
			config.Authenticator = gosnowflake.AuthTypeOAuth
		case "ExternalBrowser":
			config.Authenticator = gosnowflake.AuthTypeExternalBrowser
		case "Okta":
			config.Authenticator = gosnowflake.AuthTypeOkta
		case "JWT":
			config.Authenticator = gosnowflake.AuthTypeJwt
		case "TokenAccessor":
			config.Authenticator = gosnowflake.AuthTypeTokenAccessor
		case "UsernamePasswordMFA":
			config.Authenticator = gosnowflake.AuthTypeUsernamePasswordMFA
		default:
			return nil, fmt.Errorf("invalid authenticator %s", authenticator)
		}
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

	if v, ok := s.GetOk("login_timeout"); ok {
		config.LoginTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("request_timeout"); ok {
		config.RequestTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("jwt_expire_timeout"); ok {
		config.JWTExpireTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("client_timeout"); ok {
		config.ClientTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("jwt_client_timeout"); ok {
		config.JWTClientTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("external_browser_timeout"); ok {
		config.ExternalBrowserTimeout = time.Duration(int64(v.(int)))
	}

	if v, ok := s.GetOk("insecure_mode"); ok && v.(bool) {
		config.InsecureMode = v.(bool)
	}

	if v, ok := s.GetOk("oscp_fail_open"); ok && v.(bool) {
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
				return nil, fmt.Errorf("could not retrieve access token from refresh token")
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
	if v, err := getPrivateKey(privateKeyPath, privateKey, privateKeyPassphrase); err != nil && v != nil {
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

	if v, ok := s.GetOk("profile"); ok && v.(string) != "" {
		profile := v.(string)
		if profile == "default" {
			defaultConfig := sdk.DefaultConfig()
			if defaultConfig.Account == "" || defaultConfig.User == "" {
				return "", errors.New("account and User must be set in provider config, ~/.snowflake/config, or as an environment variable")
			}
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
	client, err := sdk.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ProviderMeta{
		client: client,
		db:     client.GetConn().DB,
	}, nil
}
