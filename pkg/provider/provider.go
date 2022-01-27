package provider

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"io"
	"io/ioutil"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/datasources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/db"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"

	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_refresh_token"},
			},
			"oauth_refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REFRESH_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_ID", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_refresh_token", "oauth_client_secret", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_CLIENT_SECRET", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_refresh_token", "oauth_endpoint", "oauth_redirect_url"},
			},
			"oauth_endpoint": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ENDPOINT", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_refresh_token", "oauth_redirect_url"},
			},
			"oauth_redirect_url": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_REDIRECT_URL", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "password", "oauth_access_token"},
				RequiredWith:  []string{"oauth_client_id", "oauth_client_secret", "oauth_endpoint", "oauth_refresh_token"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token", "private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
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
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
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
		"snowflake_task_grant":              resources.TaskGrant(),
		"snowflake_view_grant":              resources.ViewGrant(),
		"snowflake_warehouse_grant":         resources.WarehouseGrant(),
	}
	return grants
}

func getResources() map[string]*schema.Resource {
	// NOTE(): do not add grant resources here
	others := map[string]*schema.Resource{
		"snowflake_api_integration":           resources.APIIntegration(),
		"snowflake_database":                  resources.Database(),
		"snowflake_external_function":         resources.ExternalFunction(),
		"snowflake_file_format":               resources.FileFormat(),
		"snowflake_function":                  resources.Function(),
		"snowflake_managed_account":           resources.ManagedAccount(),
		"snowflake_masking_policy":            resources.MaskingPolicy(),
		"snowflake_materialized_view":         resources.MaterializedView(),
		"snowflake_network_policy_attachment": resources.NetworkPolicyAttachment(),
		"snowflake_network_policy":            resources.NetworkPolicy(),
		"snowflake_oauth_integration":         resources.OAuthIntegration(),
		"snowflake_pipe":                      resources.Pipe(),
		"snowflake_procedure":                 resources.Procedure(),
		"snowflake_resource_monitor":          resources.ResourceMonitor(),
		"snowflake_role":                      resources.Role(),
		"snowflake_role_grants":               resources.RoleGrants(),
		"snowflake_row_access_policy":         resources.RowAccessPolicy(),
		"snowflake_saml_integration":          resources.SAMLIntegration(),
		"snowflake_schema":                    resources.Schema(),
		"snowflake_scim_integration":          resources.SCIMIntegration(),
		"snowflake_sequence":                  resources.Sequence(),
		"snowflake_share":                     resources.Share(),
		"snowflake_stage":                     resources.Stage(),
		"snowflake_storage_integration":       resources.StorageIntegration(),
		"snowflake_notification_integration":  resources.NotificationIntegration(),
		"snowflake_stream":                    resources.Stream(),
		"snowflake_table":                     resources.Table(),
		"snowflake_external_table":            resources.ExternalTable(),
		"snowflake_tag":                       resources.Tag(),
		"snowflake_task":                      resources.Task(),
		"snowflake_user":                      resources.User(),
		"snowflake_user_public_keys":          resources.UserPublicKeys(),
		"snowflake_view":                      resources.View(),
		"snowflake_warehouse":                 resources.Warehouse(),
	}

	return mergeSchemas(
		others,
		GetGrantResources().GetTfSchemas(),
	)
}

func getDataSources() map[string]*schema.Resource {
	dataSources := map[string]*schema.Resource{
		"snowflake_current_account":                    datasources.CurrentAccount(),
		"snowflake_system_generate_scim_access_token":  datasources.SystemGenerateSCIMAccessToken(),
		"snowflake_system_get_aws_sns_iam_policy":      datasources.SystemGetAWSSNSIAMPolicy(),
		"snowflake_system_get_privatelink_config":      datasources.SystemGetPrivateLinkConfig(),
		"snowflake_system_get_snowflake_platform_info": datasources.SystemGetSnowflakePlatformInfo(),
		"snowflake_schemas":                            datasources.Schemas(),
		"snowflake_tables":                             datasources.Tables(),
		"snowflake_views":                              datasources.Views(),
		"snowflake_materialized_views":                 datasources.MaterializedViews(),
		"snowflake_stages":                             datasources.Stages(),
		"snowflake_file_formats":                       datasources.FileFormats(),
		"snowflake_sequences":                          datasources.Sequences(),
		"snowflake_streams":                            datasources.Streams(),
		"snowflake_tasks":                              datasources.Tasks(),
		"snowflake_pipes":                              datasources.Pipes(),
		"snowflake_masking_policies":                   datasources.MaskingPolicies(),
		"snowflake_external_functions":                 datasources.ExternalFunctions(),
		"snowflake_external_tables":                    datasources.ExternalTables(),
		"snowflake_warehouses":                         datasources.Warehouses(),
		"snowflake_resource_monitors":                  datasources.ResourceMonitors(),
		"snowflake_storage_integrations":               datasources.StorageIntegrations(),
		"snowflake_row_access_policies":                datasources.RowAccessPolicies(),
		"snowflake_functions":                          datasources.Functions(),
		"snowflake_procedures":                         datasources.Procedures(),
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

	if oauthRefreshToken != "" {
		accessToken, err := GetOauthAccessToken(oauthEndpoint, oauthClientID, oauthClientSecret, GetOauthData(oauthRefreshToken, oauthRedirectURL))
		if err != nil {
			return nil, errors.Wrap(err, "could not retreive access token from refresh token")
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
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not build dsn for snowflake connection")
	}

	db, err := db.Open(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open snowflake database.")
	}

	return db, nil
}

func DSN(
	account,
	user,
	password string,
	browserAuth bool,
	privateKeyPath,
	privateKey,
	privateKeyPassphrase,
	oauthAccessToken,
	region,
	role string) (string, error) {

	// us-west-2 is their default region, but if you actually specify that it won't trigger their default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if region == "us-west-2" {
		region = ""
	}

	config := gosnowflake.Config{
		Account: account,
		User:    user,
		Region:  region,
		Role:    role,
	}

	if privateKeyPath != "" {
		privateKeyBytes, err := ReadPrivateKeyFile(privateKeyPath)
		if err != nil {
			return "", errors.Wrap(err, "Private Key file could not be read")
		}
		rsaPrivateKey, err := ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
		if err != nil {
			return "", errors.Wrap(err, "Private Key could not be parsed")
		}
		config.PrivateKey = rsaPrivateKey
		config.Authenticator = gosnowflake.AuthTypeJwt
	} else if privateKey != "" {
		rsaPrivateKey, err := ParsePrivateKey([]byte(privateKey), []byte(privateKeyPassphrase))
		if err != nil {
			return "", errors.Wrap(err, "Private Key could not be parsed")
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
	} else {
		return "", errors.New("no authentication method provided")
	}

	return gosnowflake.DSN(&config)
}

func ReadPrivateKeyFile(privateKeyPath string) ([]byte, error) {
	expandedPrivateKeyPath, err := homedir.Expand(privateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Path to private key")
	}

	privateKeyBytes, err := ioutil.ReadFile(expandedPrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read private key")
	}

	if len(privateKeyBytes) == 0 {
		return nil, errors.New("Private key is empty")
	}

	return privateKeyBytes, nil
}

func ParsePrivateKey(privateKeyBytes []byte, passhrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("Could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passhrase) == 0 {
			return nil, fmt.Errorf("Private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passhrase)
		if err != nil {
			return nil, errors.Wrap(
				err,
				"Could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported",
			)
		}
		return privateKey, nil
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse private key")
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

func GetOauthData(refreshToken, redirectUrl string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", redirectUrl)
	return data
}

func GetOauthRequest(dataContent io.Reader, endPoint, clientId, clientSecret string) (*http.Request, error) {
	request, err := http.NewRequest("POST", endPoint, dataContent)
	if err != nil {
		return nil, errors.Wrap(err, "Request to the endpoint could not be completed")
	}
	request.SetBasicAuth(clientId, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	return request, nil

}

func GetOauthAccessToken(
	endPoint,
	client_id,
	client_secret string,
	data url.Values) (string, error) {

	client := &http.Client{}
	request, err := GetOauthRequest(strings.NewReader(data.Encode()), endPoint, client_id, client_secret)
	if err != nil {
		return "", errors.Wrap(err, "Oauth request returned an error:")
	}

	var result Result

	response, err := client.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "Response status returned an error:")
	}
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Response status code: %s: %s", strconv.Itoa(response.StatusCode), http.StatusText(response.StatusCode)))
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "Response body was not able to be parsed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", errors.Wrap(err, "Error parsing JSON from Snowflake")
	}
	return result.AccessToken, nil
}
