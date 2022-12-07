package resources

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// warehouseCreateProperties are only available via the CREATE statement.
var warehouseCreateProperties = []string{"initially_suspended", "wait_for_provisioning"}

var warehouseProperties = []string{
	"comment", "warehouse_size", "max_cluster_count", "min_cluster_count",
	"scaling_policy", "auto_suspend", "auto_resume",
	"resource_monitor", "max_concurrency_level", "statement_queued_timeout_in_seconds",
	"statement_timeout_in_seconds", "enable_query_acceleration", "query_acceleration_max_scale_factor",
	"warehouse_type",
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
				return strings.ToUpper(strings.ReplaceAll(s, "-", ""))
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
		ValidateFunc: validation.IntAtLeast(1),
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
		Default:     172800,
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
		Default:     8,
		Description: "Object parameter that specifies the concurrency level for SQL statements (i.e. queries and DML) executed by a warehouse.",
	},
	"enable_query_acceleration": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources.",
	},
	"query_acceleration_max_scale_factor": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      8,
		ValidateFunc: validation.IntBetween(0, 100),
		Description:  "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
	},
	"warehouse_type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "STANDARD",
		ValidateFunc: validation.StringInSlice([]string{"STANDARD", "SNOWPARK-OPTIMIZED"}, true),
		Description:  "Specifies a STANDARD or SNOWPARK-OPTIMIZED warehouse",
	},
	"tag": tagReferenceSchema,
}

// Warehouse returns a pointer to the resource representing a warehouse.
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

// CreateWarehouse implements schema.CreateFunc.
func CreateWarehouse(d *schema.ResourceData, meta interface{}) error {
	props := append(warehouseProperties, warehouseCreateProperties...) //nolint:gocritic // todo: please fix this to pass gocritic
	return CreateResource(
		"warehouse",
		props,
		warehouseSchema,
		func(name string) *snowflake.Builder {
			return snowflake.Warehouse(name).Builder
		},
		ReadWarehouse,
	)(d, meta)
}

// ReadWarehouse implements schema.ReadFunc.
func ReadWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	warehouseBuilder := snowflake.Warehouse(d.Id())
	stmt := warehouseBuilder.Show()

	row := snowflake.QueryRow(db, stmt)
	w, err := snowflake.ScanWarehouse(row)
	if errors.Is(err, sql.ErrNoRows) {
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
	err = d.Set("auto_suspend", w.AutoSuspend.Int64)
	if err != nil {
		return err
	}
	err = d.Set("auto_resume", w.AutoResume)
	if err != nil {
		return err
	}
	err = d.Set("resource_monitor", w.ResourceMonitor)
	if err != nil {
		return err
	}
	err = d.Set("enable_query_acceleration", w.EnableQueryAcceleration)
	if err != nil {
		return err
	}
	err = d.Set("query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor)
	if err != nil {
		return err
	}
	err = d.Set("warehouse_type", w.WarehouseType)
	if err != nil {
		return err
	}

	stmt = warehouseBuilder.ShowParameters()
	paramRows, err := snowflake.Query(db, stmt)
	if err != nil {
		return err
	}

	warehouseParams, err := snowflake.ScanWarehouseParameters(paramRows)
	if err != nil {
		return err
	}

	for _, param := range warehouseParams {
		log.Printf("[TRACE] %+v\n", param)

		var value interface{}
		if strings.EqualFold(param.Type, "number") {
			i, err := strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				return err
			}
			value = i
		} else {
			value = param.Value
		}

		key := strings.ToLower(param.Key)
		// lintignore:R001
		err = d.Set(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateWarehouse implements schema.UpdateFunc.
func UpdateWarehouse(d *schema.ResourceData, meta interface{}) error {
	return UpdateResource(
		"warehouse",
		warehouseProperties,
		warehouseSchema,
		func(name string) *snowflake.Builder {
			return snowflake.Warehouse(name).Builder
		},
		ReadWarehouse,
	)(d, meta)
}

// DeleteWarehouse implements schema.DeleteFunc.
func DeleteWarehouse(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource(
		"warehouse", func(name string) *snowflake.Builder {
			return snowflake.Warehouse(name).Builder
		},
	)(d, meta)
}
