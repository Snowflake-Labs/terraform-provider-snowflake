package provider

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/db"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"golang.org/x/crypto/ssh"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "private_key_path"},
			},
			"browser_auth": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "private_key_path"},
			},
			"private_key_path": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth", "password"},
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snowflake_account_grant":          resources.AccountGrant(),
			"snowflake_database":               resources.Database(),
			"snowflake_database_grant":         resources.DatabaseGrant(),
			"snowflake_integration_grant":      resources.IntegrationGrant(),
			"snowflake_managed_account":        resources.ManagedAccount(),
			"snowflake_pipe":                   resources.Pipe(),
			"snowflake_resource_monitor":       resources.ResourceMonitor(),
			"snowflake_resource_monitor_grant": resources.ResourceMonitorGrant(),
			"snowflake_role":                   resources.Role(),
			"snowflake_role_grants":            resources.RoleGrants(),
			"snowflake_schema":                 resources.Schema(),
			"snowflake_schema_grant":           resources.SchemaGrant(),
			"snowflake_share":                  resources.Share(),
			"snowflake_stage":                  resources.Stage(),
			"snowflake_stage_grant":            resources.StageGrant(),
			"snowflake_storage_integration":    resources.StorageIntegration(),
			"snowflake_user":                   resources.User(),
			"snowflake_view":                   resources.View(),
			"snowflake_view_grant":             resources.ViewGrant(),
			"snowflake_table_grant":            resources.TableGrant(),
			"snowflake_warehouse":              resources.Warehouse(),
			"snowflake_warehouse_grant":        resources.WarehouseGrant(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  ConfigureProvider,
	}
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	dsn, err := DSN(s)

	if err != nil {
		return nil, errors.Wrap(err, "could not build dsn for snowflake connection")
	}

	db, err := db.Open(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open snowflake database.")
	}

	return db, nil
}

func DSN(s *schema.ResourceData) (string, error) {
	account := s.Get("account").(string)
	user := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	privateKeyPath := s.Get("private_key_path").(string)
	region := s.Get("region").(string)
	role := s.Get("role").(string)

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
	} else {
		config.Password = password
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
