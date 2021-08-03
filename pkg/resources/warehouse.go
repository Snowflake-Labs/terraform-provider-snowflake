package resources

import (
	"database/sql"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// warehouseCreateProperties are only available via the CREATE statement
var warehouseCreateProperties = []string{"initially_suspended", "wait_for_provisioning"}

var warehouseProperties = []string{
	"comment", "warehouse_size", "max_cluster_count", "min_cluster_count",
	"scaling_policy", "auto_suspend", "auto_resume",
	"resource_monitor", "max_concurrency_level", "statement_queued_timeout_in_seconds",
	"statement_timeout_in_seconds",
}

var warehouseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the virtual warehouse; must be unique for your account.",
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"warehouse_size": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateFunc: validation.StringInSlice([]string{
			"XSMALL", "X-SMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE",
			"X-LARGE", "XXLARGE", "X2LARGE", "2X-LARGE", "XXXLARGE", "X3LARGE",
			"3X-LARGE", "X4LARGE", "4X-LARGE", "X5LARGE", "5X-LARGE", "X6LARGE",
			"6X-LARGE",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
		Description: "Specifies the size of the virtual warehouse. Larger warehouse sizes 5X-Large and 6X-Large are currently in preview and only available on Amazon Web Services (AWS).",
	},
	"max_cluster_count": {
		Type:         schema.TypeInt,
		Description:  "Specifies the maximum number of server clusters for the warehouse.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"min_cluster_count": {
		Type:         schema.TypeInt,
		Description:  "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"scaling_policy": {
		Type:         schema.TypeString,
		Description:  "Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"STANDARD", "ECONOMY"}, true),
	},
	"auto_suspend": {
		Type:         schema.TypeInt,
		Description:  "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	// @TODO add a disable_auto_suspend property that sets the value of auto_suspend to NULL
	"auto_resume": {
		Type:        schema.TypeBool,
		Description: "Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.",
		Optional:    true,
		Computed:    true,
	},
	"initially_suspended": {
		Type:        schema.TypeBool,
		Description: "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
		Optional:    true,
		ForceNew:    true,
	},
	"resource_monitor": {
		Type:        schema.TypeString,
		Description: "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		Optional:    true,
		Computed:    true,
	},
	"wait_for_provisioning": {
		Type:        schema.TypeBool,
		Description: "Specifies whether the warehouse, after being resized, waits for all the servers to provision before executing any queued or new queries.",
		Optional:    true,
		ForceNew:    true,
	},
	"statement_timeout_in_seconds": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system",
	},
	"statement_queued_timeout_in_seconds": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Object parameter that specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on a warehouse before it is canceled by the system.",
	},
	"max_concurrency_level": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Object parameter that specifies the concurrency level for SQL statements (i.e. queries and DML) executed by a warehouse.",
	},
}

// Warehouse returns a pointer to the resource representing a warehouse
func Warehouse() *schema.Resource {
	return &schema.Resource{
		Create: CreateWarehouse,
		Read:   ReadWarehouse,
		Delete: DeleteWarehouse,
		Update: UpdateWarehouse,

		Schema: warehouseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateWarehouse implements schema.CreateFunc
func CreateWarehouse(d *schema.ResourceData, meta interface{}) error {
	props := append(warehouseProperties, warehouseCreateProperties...)
	return CreateResource("warehouse", props, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(d, meta)
}

// ReadWarehouse implements schema.ReadFunc
func ReadWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	stmt := snowflake.Warehouse(d.Id()).Show()

	row := snowflake.QueryRow(db, stmt)
	w, err := snowflake.ScanWarehouse(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Printf("[DEBUG] warehouse (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", w.Name)
	if err != nil {
		return err
	}
	err = d.Set("comment", w.Comment)
	if err != nil {
		return err
	}
	err = d.Set("warehouse_size", w.Size)
	if err != nil {
		return err
	}
	err = d.Set("max_cluster_count", w.MaxClusterCount)
	if err != nil {
		return err
	}
	err = d.Set("min_cluster_count", w.MinClusterCount)
	if err != nil {
		return err
	}
	err = d.Set("scaling_policy", w.ScalingPolicy)
	if err != nil {
		return err
	}
	err = d.Set("auto_suspend", w.AutoSuspend)
	if err != nil {
		return err
	}
	err = d.Set("auto_resume", w.AutoResume)
	if err != nil {
		return err
	}
	err = d.Set("statement_timeout_in_seconds", w.StatementTimeoutInSeconds)
	if err != nil {
		return err
	}
	err = d.Set("statement_queued_timeout_in_seconds", w.StatementQueuedTimeoutInSeconds)
	if err != nil {
		return err
	}
	err = d.Set("max_concurrency_level", w.MaxConcurrencyLevel)
	if err != nil {
		return err
	}
	err = d.Set("resource_monitor", w.ResourceMonitor)

	return err
}

// UpdateWarehouse implements schema.UpdateFunc
func UpdateWarehouse(d *schema.ResourceData, meta interface{}) error {
	return UpdateResource("warehouse", warehouseProperties, warehouseSchema, snowflake.Warehouse, ReadWarehouse)(d, meta)
}

// DeleteWarehouse implements schema.DeleteFunc
func DeleteWarehouse(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("warehouse", snowflake.Warehouse)(d, meta)
}
