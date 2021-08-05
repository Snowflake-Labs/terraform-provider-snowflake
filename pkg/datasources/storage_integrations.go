package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var storageIntegrationsSchema = map[string]*schema.Schema{
	"storage_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The storage integrations in the database",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func StorageIntegrations() *schema.Resource {
	return &schema.Resource{
		Read:   ReadStorageIntegrations,
		Schema: storageIntegrationsSchema,
	}
}

func ReadStorageIntegrations(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	account, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	currentStorageIntegrations, err := snowflake.ListStorageIntegrations(db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] no storage integrations found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse storage integrations in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	storageIntegrations := []map[string]interface{}{}

	for _, storageIntegration := range currentStorageIntegrations {
		storageIntegrationMap := map[string]interface{}{}

		storageIntegrationMap["name"] = storageIntegration.Name.String
		storageIntegrationMap["type"] = storageIntegration.IntegrationType.String
		storageIntegrationMap["comment"] = storageIntegration.Comment.String
		storageIntegrationMap["enabled"] = storageIntegration.Enabled.Bool

		storageIntegrations = append(storageIntegrations, storageIntegrationMap)
	}

	return d.Set("storage_integrations", storageIntegrations)
}
