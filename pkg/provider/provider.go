package provider

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/db"
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
			"username": {
				Type:        schema.TypeString,
				Description: "Username for username+password authentication. Can come from the `SNOWFLAKE_USER` environment variable. Required unless using profile.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Password for username+password auth. Cannot be used with `browser_auth` or `private_key_path`. Can be sourced from `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Generating the token is left to other tools. Cannot be used with `browser_auth`, `private_key_path`, `oauth_refresh_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Description:   "Token for use with OAuth. Setup and generation of the token is left to other tools. Should be used in conjunction with `oauth_client_id`, `oauth_client_secret`, `oauth_endpoint`, `oauth_redirect_url`. Cannot be used with `browser_auth`, `private_key_path`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_OAUTH_REFRESH_TOKEN` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_ENDPOINT` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_OAUTH_REDIRECT_URL` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Description:   "Required when `oauth_refresh_token` is used. Can be sourced from `SNOWFLAKE_USE_BROWSER_AUTH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Description:   "Path to a private key for using keypair authentication. Cannot be used with `browser_auth`, `oauth_access_token` or `password`. Can be sourced from `SNOWFLAKE_PRIVATE_KEY_PATH` environment variable.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key"},
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
			"protocol": {
				Type:        schema.TypeString,
				Description: "Support custom protocols to snowflake go driver. Can be sourced from `SNOWFLAKE_PROTOCOL` environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROTOCOL", "https"),
			},
			"insecure_mode": {
				Type:        schema.TypeBool,
				Description: "If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_INSECURE_MODE", false),
			},
			"warehouse": {
				Type:        schema.TypeString,
				Description: "Sets the default warehouse. Optional. Can be sourced from SNOWFLAKE_WAREHOUSE environment variable.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_WAREHOUSE", nil),
			},
			"profile": {
				Type:        schema.TypeString,
				Description: "Sets the profile to read from ~/.snowflake/config file.",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_PROFILE", "default"),
			},
		},
		ResourcesMap:   getResources(),
		DataSourcesMap: getDataSources(),
		ConfigureFunc:  ConfigureProvider,
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
		"snowflake_account_parameter":                       resources.AccountParameter(),
		"snowflake_alert":                                   resources.Alert(),
		"snowflake_api_integration":                         resources.APIIntegration(),
		"snowflake_database":                                resources.Database(),
		"snowflake_database_role":                           resources.DatabaseRole(),
		"snowflake_email_notification_integration":          resources.EmailNotificationIntegration(),
		"snowflake_external_function":                       resources.ExternalFunction(),
		"snowflake_external_oauth_integration":              resources.ExternalOauthIntegration(),
		"snowflake_external_table":                          resources.ExternalTable(),
		"snowflake_failover_group":                          resources.FailoverGroup(),
		"snowflake_file_format":                             resources.FileFormat(),
		"snowflake_function":                                resources.Function(),
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

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	account := s.Get("account").(string)
	user := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	privateKeyPath := s.Get("private_key_path").(string)
	privateKey := s.Get("private_key").(string)
	privateKeyPassphrase := s.Get("private_key_passphrase").(string)
	oauthAccessToken := s.Get("oauth_access_token").(string)
	region := s.Get("region").(string)
	role := s.Get("role").(string)
	oauthRefreshToken := s.Get("oauth_refresh_token").(string)
	oauthClientID := s.Get("oauth_client_id").(string)
	oauthClientSecret := s.Get("oauth_client_secret").(string)
	oauthEndpoint := s.Get("oauth_endpoint").(string)
	oauthRedirectURL := s.Get("oauth_redirect_url").(string)
	host := s.Get("host").(string)
	protocol := s.Get("protocol").(string)
	port := s.Get("port").(int)
	warehouse := s.Get("warehouse").(string)
	insecureMode := s.Get("insecure_mode").(bool)
	profile := s.Get("profile").(string)

	if oauthRefreshToken != "" {
		accessToken, err := GetOauthAccessToken(oauthEndpoint, oauthClientID, oauthClientSecret, GetOauthData(oauthRefreshToken, oauthRedirectURL))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve access token from refresh token")
		}
		oauthAccessToken = accessToken
	}

	dsn, err := DSN(
		account,
		user,
		password,
		browserAuth,
		privateKeyPath,
		privateKey,
		privateKeyPassphrase,
		oauthAccessToken,
		region,
		role,
		host,
		protocol,
		port,
		warehouse,
		insecureMode,
		profile,
	)
	if err != nil {
		return nil, fmt.Errorf("could not build dsn for snowflake connection err = %w", err)
	}

	db, err := db.Open(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open snowflake database err = %w", err)
	}
	client := sdk.NewClientFromDB(db)
	sessionID, err := client.ContextFunctions.CurrentSession(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve session id err = %w", err)
	}
	log.Printf("[INFO] Snowflake DB connection opened, session ID : %s\n", sessionID)
	if err != nil {
		return nil, fmt.Errorf("could not open snowflake database err = %w", err)
	}

	return db, nil
}

func DSN(
	account string,
	user string,
	password string,
	browserAuth bool,
	privateKeyPath string,
	privateKey string,
	privateKeyPassphrase string,
	oauthAccessToken string,
	region string,
	role string,
	host string,
	protocol string,
	port int,
	warehouse string,
	insecureMode bool,
	profile string,
) (string, error) {
	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if region == "us-west-2" {
		region = ""
	}

	config := &gosnowflake.Config{
		Account:      account,
		User:         user,
		Region:       region,
		Role:         role,
		Port:         port,
		Protocol:     protocol,
		InsecureMode: insecureMode,
	}

	// If host is set trust it and do not use the region value
	if host != "" {
		config.Region = ""
		config.Host = host
	}

	// If warehouse is set
	if warehouse != "" {
		config.Warehouse = warehouse
	}

	if privateKeyPath != "" { //nolint:gocritic // todo: please fix this to pass gocritic
		privateKeyBytes, err := ReadPrivateKeyFile(privateKeyPath)
		if err != nil {
			return "", fmt.Errorf("private Key file could not be read err = %w", err)
		}
		rsaPrivateKey, err := ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
		if err != nil {
			return "", fmt.Errorf("private Key could not be parsed err = %w", err)
		}
		config.PrivateKey = rsaPrivateKey
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else if privateKey != "" {
		rsaPrivateKey, err := ParsePrivateKey([]byte(privateKey), []byte(privateKeyPassphrase))
		if err != nil {
			return "", fmt.Errorf("private Key could not be parsed err = %w", err)
		}
		config.PrivateKey = rsaPrivateKey
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else if browserAuth {
		config.Authenticator = gosnowflake.AuthTypeExternalBrowser
	} else if oauthAccessToken != "" {
		config.Authenticator = gosnowflake.AuthTypeOAuth
		config.Token = oauthAccessToken
	} else if password != "" {
		config.Password = password
	} else if account == "" && user == "" {
		// If account and user are empty then we need to fall back on using profile config
		log.Printf("[DEBUG] No account or user provided, falling back to profile %s\n", profile)
		profileConfig, err := sdk.ProfileConfig(profile)
		if err != nil {
			return "", errors.New("no authentication method provided")
		}
		config = sdk.MergeConfig(config, profileConfig)
	}
	config.Application = "terraform-provider-snowflake"
	return gosnowflake.DSN(config)
}

func ReadPrivateKeyFile(privateKeyPath string) ([]byte, error) {
	expandedPrivateKeyPath, err := homedir.Expand(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("invalid Path to private key err = %w", err)
	}

	privateKeyBytes, err := os.ReadFile(expandedPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key err = %w", err)
	}

	if len(privateKeyBytes) == 0 {
		return nil, errors.New("private key is empty")
	}

	return privateKeyBytes, nil
}

func ParsePrivateKey(privateKeyBytes []byte, passhrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passhrase) == 0 {
			return nil, fmt.Errorf("private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passhrase)
		if err != nil {
			return nil, fmt.Errorf("could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported err = %w", err)
		}
		return privateKey, nil
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key err = %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey not of type RSA")
	}
	return rsaPrivateKey, nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetOauthData(refreshToken, redirectURL string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", redirectURL)
	return data
}

func GetOauthRequest(dataContent io.Reader, endPoint, clientID, clientSecret string) (*http.Request, error) {
	request, err := http.NewRequest("POST", endPoint, dataContent)
	if err != nil {
		return nil, fmt.Errorf("request to the endpoint could not be completed %w", err)
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	return request, nil
}

func GetOauthAccessToken(
	endPoint,
	clientID,
	clientSecret string,
	data url.Values,
) (string, error) {
	client := &http.Client{}
	request, err := GetOauthRequest(strings.NewReader(data.Encode()), endPoint, clientID, clientSecret)
	if err != nil {
		return "", fmt.Errorf("oauth request returned an error")
	}

	var result Result

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("response status returned an err = %w", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("response status code: %s: %s err = %w", strconv.Itoa(response.StatusCode), http.StatusText(response.StatusCode), err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("response body was not able to be parsed err = %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON from Snowflake err = %w", err)
	}
	return result.AccessToken, nil
}

func GetDatabaseHandleFromEnv() (db *sql.DB, err error) {
	account := os.Getenv("SNOWFLAKE_ACCOUNT")
	user := os.Getenv("SNOWFLAKE_USER")
	password := os.Getenv("SNOWFLAKE_PASSWORD")
	browserAuth := os.Getenv("SNOWFLAKE_BROWSER_AUTH") == "true"
	privateKeyPath := os.Getenv("SNOWFLAKE_PRIVATE_KEY_PATH")
	privateKey := os.Getenv("SNOWFLAKE_PRIVATE_KEY")
	privateKeyPassphrase := os.Getenv("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE")
	oauthAccessToken := os.Getenv("SNOWFLAKE_OAUTH_ACCESS_TOKEN")
	region := os.Getenv("SNOWFLAKE_REGION")
	role := os.Getenv("SNOWFLAKE_ROLE")
	host := os.Getenv("SNOWFLAKE_HOST")
	warehouse := os.Getenv("SNOWFLAKE_WAREHOUSE")
	protocol := os.Getenv("SNOWFLAKE_PROTOCOL")
	profile := os.Getenv("SNOWFLAKE_PROFILE")
	if profile == "" {
		profile = "default"
	}
	port, err := strconv.Atoi(os.Getenv("SNOWFLAKE_PORT"))
	if err != nil {
		port = 443
	}
	dsn, err := DSN(
		account,
		user,
		password,
		browserAuth,
		privateKeyPath,
		privateKey,
		privateKeyPassphrase,
		oauthAccessToken,
		region,
		role,
		host,
		protocol,
		port,
		warehouse,
		false,
		profile,
	)
	if err != nil {
		return nil, err
	}
	db, err = sql.Open("snowflake", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
