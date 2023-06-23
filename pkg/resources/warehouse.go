package resources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	snowflakevalidation "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: snowflakevalidation.ValidateWarehouseSize,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			oldSize, err := sdk.ToWarehouseSize(old)
			if err != nil {
				return false
			}
			newSize, err := sdk.ToWarehouseSize(new)
			if err != nil {
				return false
			}
			return oldSize == newSize
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
		Type:        schema.TypeString,
		Description: "Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.",
		Optional:    true,
		Computed:    true,
		ValidateFunc: validation.StringInSlice([]string{
			string(sdk.ScalingPolicyStandard),
			string(sdk.ScalingPolicyEconomy),
		}, true),
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
		Deprecated:  "This field is deprecated and will be removed in the next major version of the provider. It doesn't do anything and should be removed from your configuration.",
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
		Type:     schema.TypeString,
		Optional: true,
		Default:  string(sdk.WarehouseTypeStandard),
		ValidateFunc: validation.StringInSlice([]string{
			string(sdk.WarehouseTypeStandard),
			string(sdk.WarehouseTypeSnowparkOptimized),
		}, true),
		Description: "Specifies a STANDARD or SNOWPARK-OPTIMIZED warehouse",
	},
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
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	name := d.Get("name").(string)
	objectIdentifier := sdk.NewAccountObjectIdentifier(name)
	whType := sdk.WarehouseType(d.Get("warehouse_type").(string))
	createOptions := &sdk.CreateWarehouseOptions{
		Comment:                         sdk.String(d.Get("comment").(string)),
		StatementTimeoutInSeconds:       sdk.Int(d.Get("statement_timeout_in_seconds").(int)),
		StatementQueuedTimeoutInSeconds: sdk.Int(d.Get("statement_queued_timeout_in_seconds").(int)),
		MaxConcurrencyLevel:             sdk.Int(d.Get("max_concurrency_level").(int)),
		EnableQueryAcceleration:         sdk.Bool(d.Get("enable_query_acceleration").(bool)),
		QueryAccelerationMaxScaleFactor: sdk.Int(d.Get("query_acceleration_max_scale_factor").(int)),
		WarehouseType:                   &whType,
	}

	if v, ok := d.GetOk("warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return err
		}
		createOptions.WarehouseSize = &size
	}
	if v, ok := d.GetOk("max_cluster_count"); ok {
		createOptions.MaxClusterCount = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("min_cluster_count"); ok {
		createOptions.MinClusterCount = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("scaling_policy"); ok {
		scalingPolicy := sdk.ScalingPolicy(v.(string))
		createOptions.ScalingPolicy = &scalingPolicy
	}
	if v, ok := d.GetOk("auto_suspend"); ok {
		createOptions.AutoSuspend = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("auto_resume"); ok {
		createOptions.AutoResume = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("initially_suspended"); ok {
		createOptions.InitiallySuspended = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("resource_monitor"); ok {
		createOptions.ResourceMonitor = sdk.String(v.(string))
	}

	err := client.Warehouses.Create(ctx, objectIdentifier, createOptions)
	if err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))

	return ReadWarehouse(d, meta)
}

// ReadWarehouse implements schema.ReadFunc.
func ReadWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return err
	}

	if err = d.Set("name", w.Name); err != nil {
		return err
	}
	if err = d.Set("comment", w.Comment); err != nil {
		return err
	}
	if err = d.Set("warehouse_type", w.Type); err != nil {
		return err
	}
	if err = d.Set("warehouse_size", w.Size); err != nil {
		return err
	}
	if err = d.Set("max_cluster_count", w.MaxClusterCount); err != nil {
		return err
	}
	if err = d.Set("min_cluster_count", w.MinClusterCount); err != nil {
		return err
	}
	if err = d.Set("scaling_policy", w.ScalingPolicy); err != nil {
		return err
	}
	if err = d.Set("auto_suspend", w.AutoSuspend); err != nil {
		return err
	}
	if err = d.Set("auto_resume", w.AutoResume); err != nil {
		return err
	}
	if err = d.Set("resource_monitor", w.ResourceMonitor); err != nil {
		return err
	}
	if err = d.Set("enable_query_acceleration", w.EnableQueryAcceleration); err != nil {
		return err
	}
	if w.EnableQueryAcceleration {
		if err = d.Set("query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor); err != nil {
			return err
		}
	}

	return nil
}

// UpdateWarehouse implements schema.UpdateFunc.
func UpdateWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	// Change name separately
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			newName := sdk.NewAccountObjectIdentifier(v.(string))
			err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
				NewName: newName,
			})
			if err != nil {
				return err
			}
			d.SetId(helpers.EncodeSnowflakeID(newName))
		} else {
			panic("name has to be set")
		}
	}

	// Batch SET operations and UNSET operations
	var runSet bool
	var runUnset bool
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}
	if d.HasChange("comment") {
		runSet = true
		set.Comment = sdk.String(d.Get("comment").(string))
	}
	if d.HasChange("warehouse_size") {
		runSet = true
		v := d.Get("warehouse_size")
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return err
		}
		set.WarehouseSize = &size
	}
	if d.HasChange("max_cluster_count") {
		if v, ok := d.GetOk("max_cluster_count"); ok {
			runSet = true
			set.MaxClusterCount = sdk.Int(v.(int))
		} else {
			runUnset = true
			unset.MaxClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("min_cluster_count") {
		if v, ok := d.GetOk("min_cluster_count"); ok {
			runSet = true
			set.MinClusterCount = sdk.Int(v.(int))
		} else {
			runUnset = true
			unset.MinClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("scaling_policy") {
		if v, ok := d.GetOk("scaling_policy"); ok {
			runSet = true
			scalingPolicy := sdk.ScalingPolicy(v.(string))
			set.ScalingPolicy = &scalingPolicy
		} else {
			runUnset = true
			unset.ScalingPolicy = sdk.Bool(true)
		}
	}
	if d.HasChange("auto_suspend") {
		if v, ok := d.GetOk("auto_suspend"); ok {
			runSet = true
			set.AutoSuspend = sdk.Int(v.(int))
		} else {
			runUnset = true
			unset.AutoSuspend = sdk.Bool(true)
		}
	}
	if d.HasChange("auto_resume") {
		if v, ok := d.GetOk("auto_resume"); ok {
			runSet = true
			set.AutoResume = sdk.Bool(v.(bool))
		} else {
			runUnset = true
			unset.AutoResume = sdk.Bool(true)
		}
	}
	if d.HasChange("resource_monitor") {
		if v, ok := d.GetOk("resource_monitor"); ok {
			runSet = true
			set.ResourceMonitor = sdk.NewAccountObjectIdentifier(v.(string))
		} else {
			runUnset = true
			unset.ResourceMonitor = sdk.Bool(true)
		}
	}
	if d.HasChange("statement_timeout_in_seconds") {
		runSet = true
		set.StatementTimeoutInSeconds = sdk.Int(d.Get("statement_timeout_in_seconds").(int))
	}
	if d.HasChange("statement_queued_timeout_in_seconds") {
		runSet = true
		set.StatementQueuedTimeoutInSeconds = sdk.Int(d.Get("statement_queued_timeout_in_seconds").(int))
	}
	if d.HasChange("max_concurrency_level") {
		runSet = true
		set.MaxConcurrencyLevel = sdk.Int(d.Get("max_concurrency_level").(int))
	}
	if d.HasChange("enable_query_acceleration") {
		runSet = true
		set.EnableQueryAcceleration = sdk.Bool(d.Get("enable_query_acceleration").(bool))
	}
	if d.HasChange("query_acceleration_max_scale_factor") {
		runSet = true
		set.QueryAccelerationMaxScaleFactor = sdk.Int(d.Get("query_acceleration_max_scale_factor").(int))
	}
	if d.HasChange("warehouse_type") {
		if v, ok := d.GetOk("warehouse_type"); ok {
			runSet = true
			whType := sdk.WarehouseType(v.(string))
			set.WarehouseType = &whType
		} else {
			runUnset = true
			unset.WarehouseType = sdk.Bool(true)
		}
	}

	// Apply SET and UNSET changes
	if runSet {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Set: &set,
		})
		if err != nil {
			return err
		}
	}
	if runUnset {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Unset: &unset,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteWarehouse implements schema.DeleteFunc.
func DeleteWarehouse(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Warehouses.Drop(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}
