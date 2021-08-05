package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceMonitorsSchema = map[string]*schema.Schema{
	"resource_monitors": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The resource monitors in the database",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"frequency": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"credit_quota": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func ResourceMonitors() *schema.Resource {
	return &schema.Resource{
		Read:   ReadResourceMonitors,
		Schema: resourceMonitorsSchema,
	}
}

func ReadResourceMonitors(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	account, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	currentResourceMonitors, err := snowflake.ListResourceMonitors(db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] no resource monitors found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse resource monitors in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	resourceMonitors := []map[string]interface{}{}

	for _, resourceMonitor := range currentResourceMonitors {
		resourceMonitorMap := map[string]interface{}{}

		resourceMonitorMap["name"] = resourceMonitor.Name.String
		resourceMonitorMap["frequency"] = resourceMonitor.Frequency.String
		resourceMonitorMap["credit_quota"] = resourceMonitor.CreditQuota.String
		resourceMonitorMap["comment"] = resourceMonitor.Comment.String

		resourceMonitors = append(resourceMonitors, resourceMonitorMap)
	}

	return d.Set("resource_monitors", resourceMonitors)
}
