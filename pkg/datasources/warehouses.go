package datasources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	account, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.SetId(fmt.Sprintf("%s.%s", account.Account, account.Region))

	result, err := client.Warehouses.Show(ctx, nil)
	if err != nil {
		return err
	}

	warehouses := []map[string]interface{}{}

	for _, warehouse := range result {
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
