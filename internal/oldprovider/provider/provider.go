package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/oldprovider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/oldprovider/db"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/oldprovider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/sdk"
)

// Provider is a provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Description:  "",
				Optional:    true,
				//DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
			},
			"user": {
				Type:        schema.TypeString,
				Description: "",
				Optional:    true,
				//DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "",
				Optional:      true,
				//DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				//ConflictsWith: []string{"browser_auth", "private_key_path", "private_key", "private_key_passphrase", "oauth_access_token", "oauth_refresh_token"},
			},
			"role": {
				Type:        schema.TypeString,
				Description: "",
				Optional:    true,
				//DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
		},
		ResourcesMap:   getResources(),
		DataSourcesMap: getDataSources(),
		ConfigureFunc:  ConfigureProvider,
	}
}

func getResources() map[string]*schema.Resource {
	// NOTE(): do not add grant resources here
	return map[string]*schema.Resource{
		"snowflake_warehouse": resources.Warehouse(),
	}
}

func getDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"snowflake_warehouses": datasources.Warehouses(),
	}
}

func ConfigureProvider(s *schema.ResourceData) (interface{}, error) {
	account := s.Get("account").(string)
	user := s.Get("user").(string)
	password := s.Get("password").(string)
	role := s.Get("role").(string)
	config := &gosnowflake.Config{
		Account:  account,
		User:     user,
		Role:     role,
		Password: password,
	}
	dsn, err := gosnowflake.DSN(config)
	if err != nil {
		return nil, fmt.Errorf("could not build dsn for snowflake connection err = %w", err)
	}

	db, err := db.Open(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open snowflake database err = %w", err)
	}
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	sessionID, err := client.ContextFunctions.CurrentSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve session id err = %w", err)
	}
	log.Printf("[INFO] Snowflake DB connection opened, session ID : %s\n", sessionID)
	if err != nil {
		return nil, fmt.Errorf("could not open snowflake database err = %w", err)
	}

	return db, nil
}
