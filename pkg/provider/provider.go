package provider

import (
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/db"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
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
				ConflictsWith: []string{"password"},
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
			"snowflake_database":        resources.Database(),
			"snowflake_managed_account": resources.ManagedAccount(),
			"snowflake_role":            resources.Role(),
			"snowflake_role_grants":     resources.RoleGrants(),
			"snowflake_user":            resources.User(),
			"snowflake_warehouse":       resources.Warehouse(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  ConfigureProvider,
	}
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	dsn, err := DSN(s)

	log.Printf("[DEBUG] connecting to %s", dsn)
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
	region := s.Get("region").(string)
	role := s.Get("role").(string)

	// us-west-2 is their default region, but if you actually specify that it won't trigger their default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if region == "us-west-2" {
		region = ""
	}

	dsn, err := gosnowflake.DSN(&gosnowflake.Config{
		Account:  account,
		User:     username,
		Region:   region,
		Password: password,
		Role:     role,
	})

	if browserAuth {
		dsn, err = gosnowflake.DSN(&gosnowflake.Config{
			Account:       account,
			User:          username,
			Region:        region,
			Role:          role,
			Authenticator: "externalbrowser",
		})
	}

	return dsn, err
}
