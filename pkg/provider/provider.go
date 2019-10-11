package provider

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/db"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ConflictsWith: []string{"browser_auth"},
			},
			"browser_auth": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password", "path_private_key"},
			},
			"path_private_key": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PRIVATE_KEY_PATH", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth"},
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
			"snowflake_database":         resources.Database(),
			"snowflake_database_grant":   resources.DatabaseGrant(),
			"snowflake_managed_account":  resources.ManagedAccount(),
			"snowflake_resource_monitor": resources.ResourceMonitor(),
			"snowflake_role":             resources.Role(),
			"snowflake_role_grants":      resources.RoleGrants(),
			"snowflake_schema":           resources.Schema(),
			"snowflake_schema_grant":     resources.SchemaGrant(),
			"snowflake_share":            resources.Share(),
			"snowflake_user":             resources.User(),
			"snowflake_view":             resources.View(),
			"snowflake_view_grant":       resources.ViewGrant(),
			"snowflake_table_grant":      resources.TableGrant(),
			"snowflake_warehouse":        resources.Warehouse(),
			"snowflake_warehouse_grant":  resources.WarehouseGrant(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  ConfigureProvider,
	}
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	dsn, err := DSN(s) // got an error

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
	username := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	pathPrivateKey := s.Get("path_private_key").(string)
	region := s.Get("region").(string)
	role := s.Get("role").(string)

	// us-west-2 is their default region, but if you actually specify that it won't trigger their default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if region == "us-west-2" {
		region = ""
	}

	var dsn string
	var err error

	if len(pathPrivateKey) != 0 {

		// reading the private key
		privateKeyBytes, err := ioutil.ReadFile(pathPrivateKey)
		if err != nil || len(privateKeyBytes) == 0 { // both conditionals had to be false in order to get past this
			return "", errors.Errorf("Could not read private key: %s", err)
		}

		// reads and unmarshals a private key
		privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
		if err != nil {
			return "", errors.Errorf("Could not parse private key: %s", err)
		}

		rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return "", errors.Errorf("privateKey not of type RSA")
		}

		dsn, err = gosnowflake.DSN(&gosnowflake.Config{
			Account:       account,
			User:          username,
			Region:        region,
			Role:          role,
			PrivateKey:    rsaPrivateKey,
			Authenticator: gosnowflake.AuthTypeJwt,
		})

		if err != nil {
			return "", errors.Errorf("JWT authentication not successful: %s", err)
		}

	} else if browserAuth {
		dsn, err = gosnowflake.DSN(&gosnowflake.Config{
			Account:       account,
			User:          username,
			Region:        region,
			Role:          role,
			Authenticator: gosnowflake.AuthTypeExternalBrowser,
		})
	} else {
		dsn, err = gosnowflake.DSN(&gosnowflake.Config{
			Account:  account,
			User:     username,
			Region:   region,
			Password: password,
			Role:     role,
		})
	}

	return dsn, err
}
