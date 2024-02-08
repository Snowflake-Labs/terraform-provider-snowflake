package datasources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	account, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	extractedResourceMonitors, err := client.ResourceMonitors.Show(ctx, &sdk.ShowResourceMonitorOptions{})
	if err != nil {
		log.Printf("[DEBUG] unable to parse resource monitors in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	resourceMonitors := make([]map[string]any, len(extractedResourceMonitors))

	for i, resourceMonitor := range extractedResourceMonitors {
		resourceMonitors[i] = map[string]any{
			"name":         resourceMonitor.Name,
			"frequency":    resourceMonitor.Frequency,
			"credit_quota": fmt.Sprintf("%f", resourceMonitor.CreditQuota),
			"comment":      resourceMonitor.Comment,
		}
	}

	return d.Set("resource_monitors", resourceMonitors)
}
