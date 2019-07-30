package resources

import (
	"database/sql"
	"fmt"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
)

var securityIntegrationProperties = []string{"type", "enabled", "oauth_client", "oauth_client_type", "oauth_redirect_uri"}
var securityIntegrationSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"type": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "OAUTH",
		ValidateFunc: SetValidationFunc(map[string]struct{}{
			"OAUTH": struct{}{},
		}),
	},
	"enabled": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"oauth_client": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "CUSTOM",
		ValidateFunc: SetValidationFunc(map[string]struct{}{
			"CUSTOM":          struct{}{},
			"TABLEAU_DESKTOP": struct{}{},
			"TABLEAU_SERVER":  struct{}{},
		}),
	},
	"oauth_client_type": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "CONFIDENTIAL",
		ValidateFunc: SetValidationFunc(map[string]struct{}{
			"CONFIDENTIAL": struct{}{},
			"PUBLIC":       struct{}{},
		}),
	},
	"oauth_redirect_uri": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
}

// SecurityIntegration https://docs.snowflake.net/manuals/user-guide/oauth.html
func SecurityIntegration() *schema.Resource {
	return &schema.Resource{
		Create: nil,
		Read:   nil,
		Delete: nil,
		Update: nil,

		Schema: securityIntegrationSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateSecurityIntegration(data *schema.ResourceData, meta interface{}) error {
	return CreateResource(
		"security integration",
		securityIntegrationProperties,
		securityIntegrationSchema,
		snowflake.SecurityIntegration,
		ReadSecurityIntegration)(data, meta)
}
func ReadSecurityIntegration(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	row := db.QueryRow(fmt.Sprintf("DESCRIBE SECURITY INTEGRATION '%s'", id))

}
func UpdateSecurityIntegration(data *schema.ResourceData, meta interface{}) error {}
func DeleteSecurityIntegration(data *schema.ResourceData, meta interface{}) error {}
