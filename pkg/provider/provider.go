package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ACCOUNT", nil),
				// TODO Validation
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_USER", nil),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", "us-west-2"),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
