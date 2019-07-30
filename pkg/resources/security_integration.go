package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

	securityIntegrationPropertiesSet := map[string]struct{}{}
	for _, prop := range securityIntegrationProperties {
		securityIntegrationPropertiesSet[prop] = struct{}{}
	}

	var property, propertyType, propertyValue, propertyDefault *string
	rows, err := db.Query(fmt.Sprintf("DESCRIBE SECURITY INTEGRATION '%s'", id))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(property, propertyType, propertyValue, propertyDefault)
		if err != nil {
			return err
		}

		normalizedPropertyName := strings.ToLower(*property)
		_, ok := securityIntegrationPropertiesSet[normalizedPropertyName]
		if !ok {
			log.Printf("[DEBUG] unrecognized SECURITY INTEGRATION property %s, skipping.", normalizedPropertyName)
			continue
		}
		data.Set(normalizedPropertyName, *propertyValue)
	}

	return rows.Err()
}
func UpdateSecurityIntegration(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource(
		"this does not seem to be used",
		securityIntegrationProperties,
		securityIntegrationSchema,
		snowflake.SecurityIntegration,
		ReadSecurityIntegration)(data, meta)
}
func DeleteSecurityIntegration(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource(
		"this does not seem to be used",
		snowflake.SecurityIntegration)(data, meta)
}

func SecurityIntegrationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.SecurityIntegration(id).Show()
	rows, err := db.Query(stmt)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), rows.Err()

}
