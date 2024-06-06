package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowWarehouseSchema should be generated later based on the sdk.Warehouse
var ShowWarehouseSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"size": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"min_cluster_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"max_cluster_count": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"started_clusters": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"running": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"queued": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"auto_suspend": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"auto_resume": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"available": {
		Type:     schema.TypeFloat,
		Computed: true,
	},
	"provisioning": {
		Type:     schema.TypeFloat,
		Computed: true,
	},
	"quiescing": {
		Type:     schema.TypeFloat,
		Computed: true,
	},
	"other": {
		Type:     schema.TypeFloat,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"resumed_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"updated_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enable_query_acceleration": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"query_acceleration_max_scale_factor": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"resource_monitor": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"scaling_policy": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

// TODO: better name?
// TODO: interface (e.g. asMap)? in SDK?
func WarehouseToSchema(warehouse *sdk.Warehouse) map[string]any {
	warehouseSchema := make(map[string]any)
	warehouseSchema["name"] = warehouse.Name
	warehouseSchema["state"] = warehouse.State
	warehouseSchema["type"] = warehouse.Type
	warehouseSchema["size"] = warehouse.Size
	warehouseSchema["min_cluster_count"] = warehouse.MinClusterCount
	warehouseSchema["max_cluster_count"] = warehouse.MaxClusterCount
	warehouseSchema["started_clusters"] = warehouse.StartedClusters
	warehouseSchema["running"] = warehouse.Running
	warehouseSchema["queued"] = warehouse.Queued
	warehouseSchema["is_default"] = warehouse.IsDefault
	warehouseSchema["is_current"] = warehouse.IsCurrent
	warehouseSchema["auto_suspend"] = warehouse.AutoSuspend
	warehouseSchema["auto_resume"] = warehouse.AutoResume
	warehouseSchema["available"] = warehouse.Available
	warehouseSchema["provisioning"] = warehouse.Provisioning
	warehouseSchema["quiescing"] = warehouse.Quiescing
	warehouseSchema["other"] = warehouse.Other
	warehouseSchema["created_on"] = warehouse.CreatedOn.String()
	warehouseSchema["resumed_on"] = warehouse.ResumedOn.String()
	warehouseSchema["updated_on"] = warehouse.UpdatedOn.String()
	warehouseSchema["owner"] = warehouse.Owner
	warehouseSchema["comment"] = warehouse.Comment
	warehouseSchema["enable_query_acceleration"] = warehouse.EnableQueryAcceleration
	warehouseSchema["query_acceleration_max_scale_factor"] = warehouse.QueryAccelerationMaxScaleFactor
	warehouseSchema["resource_monitor"] = warehouse.ResourceMonitor
	warehouseSchema["scaling_policy"] = warehouse.ScalingPolicy
	warehouseSchema["owner_role_type"] = warehouse.OwnerRoleType
	return warehouseSchema
}
