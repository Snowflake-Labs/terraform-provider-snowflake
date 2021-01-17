package provider

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/datasources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/db"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"golang.org/x/crypto/ssh"
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
				ConflictsWith: []string{"browser_auth", "private_key_path", "oauth_access_token"},
			},
			"oauth_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_OAUTH_ACCESS_TOKEN", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path", "password"},
			},
			"browser_auth": {
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path", "oauth_access_token"},
			},
			"private_key_path": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password", "oauth_access_token"},
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
		ResourcesMap: getResources(),
		DataSourcesMap: map[string]*schema.Resource{
			"snowflake_system_get_aws_sns_iam_policy": datasources.SystemGetAWSSNSIAMPolicy(),
		},
		ConfigureFunc: ConfigureProvider,
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
		"snowflake_materialized_view_grant": resources.MaterializedViewGrant(),
		"snowflake_procedure_grant":         resources.ProcedureGrant(),
		"snowflake_resource_monitor_grant":  resources.ResourceMonitorGrant(),
		"snowflake_schema_grant":            resources.SchemaGrant(),
		"snowflake_sequence_grant":          resources.SequenceGrant(),
		"snowflake_stage_grant":             resources.StageGrant(),
		"snowflake_stream_grant":            resources.StreamGrant(),
		"snowflake_table_grant":             resources.TableGrant(),
		"snowflake_view_grant":              resources.ViewGrant(),
		"snowflake_warehouse_grant":         resources.WarehouseGrant(),
	}
	return grants
}

func getResources() map[string]*schema.Resource {
	others := map[string]*schema.Resource{
		"snowflake_database":                  resources.Database(),
		"snowflake_managed_account":           resources.ManagedAccount(),
		"snowflake_masking_policy":            resources.MaskingPolicy(),
		"snowflake_network_policy_attachment": resources.NetworkPolicyAttachment(),
		"snowflake_network_policy":            resources.NetworkPolicy(),
		"snowflake_pipe":                      resources.Pipe(),
		"snowflake_resource_monitor":          resources.ResourceMonitor(),
		"snowflake_role_grants":               resources.RoleGrants(),
		"snowflake_role":                      resources.Role(),
		"snowflake_schema":                    resources.Schema(),
		"snowflake_share":                     resources.Share(),
		"snowflake_stage":                     resources.Stage(),
		"snowflake_storage_integration":       resources.StorageIntegration(),
		"snowflake_stream":                    resources.Stream(),
		"snowflake_table":                     resources.Table(),
		"snowflake_external_table":            resources.ExternalTable(),
		"snowflake_task":                      resources.Task(),
		"snowflake_user":                      resources.User(),
		"snowflake_view":                      resources.View(),
		"snowflake_warehouse":                 resources.Warehouse(),
	}

	return mergeSchemas(
		others,
		GetGrantResources().GetTfSchemas(),
	)
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	account := s.Get("account").(string)
	user := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	privateKeyPath := s.Get("private_key_path").(string)
	oauthAccessToken := s.Get("oauth_access_token").(string)
	region := s.Get("region").(string)
	role := s.Get("role").(string)

	dsn, err := DSN(account, user, password, browserAuth, privateKeyPath, oauthAccessToken, region, role)

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

		rsaPrivateKey, err := ParsePrivateKey(privateKeyPath)
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

func ParsePrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
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
