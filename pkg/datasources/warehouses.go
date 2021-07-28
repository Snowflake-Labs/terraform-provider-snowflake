package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehousesSchema = map[string]*schema.Schema{
	"warehouses": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The warehouses in the database",
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
				"size": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"scaling_policy": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"state": {
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

func Warehouses() *schema.Resource {
	return &schema.Resource{
		Read:   ReadWarehouses,
		Schema: warehousesSchema,
	}
}

func ReadWarehouses(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	account, err := snowflake.ReadCurrentAccount(db)
	if err != nil {
		log.Print("[DEBUG] unable to retrieve current account")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	currentWarehouses, err := snowflake.ListWarehouses(db)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] no warehouses found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Printf("[DEBUG] unable to parse warehouses in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	warehouses := []map[string]interface{}{}

	for _, warehouse := range currentWarehouses {
		warehouseMap := map[string]interface{}{}

		warehouseMap["name"] = warehouse.Name
		warehouseMap["type"] = warehouse.Type
		warehouseMap["size"] = warehouse.Size
		warehouseMap["scaling_policy"] = warehouse.ScalingPolicy
		warehouseMap["state"] = warehouse.State
		warehouseMap["comment"] = warehouse.Comment

		warehouses = append(warehouses, warehouseMap)
	}

	return d.Set("warehouses", warehouses)
}
