package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
