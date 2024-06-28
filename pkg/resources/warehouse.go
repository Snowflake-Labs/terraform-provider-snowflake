package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the virtual warehouse; must be unique for your account.",
	},
	"warehouse_type": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseType),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseType), IgnoreChangeToCurrentSnowflakeValueInShow("type")),
		Description:      fmt.Sprintf("Specifies warehouse type. Valid values are (case-insensitive): %s. Warehouse needs to be suspended to change its type. Provider will handle automatic suspension and resumption if needed.", possibleValuesListed(sdk.ValidWarehouseTypesString)),
	},
	"warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size")),
		Description:      fmt.Sprintf("Specifies the size of the virtual warehouse. Valid values are (case-insensitive): %s. Consult [warehouse documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for the details. Note: removing the size from config will result in the resource recreation.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 10)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count"),
		Description:      "Specifies the maximum number of server clusters for the warehouse.",
	},
	"min_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 10)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count"),
		Description:      "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
	},
	"scaling_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToScalingPolicy),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseType), IgnoreChangeToCurrentSnowflakeValueInShow("scaling_policy")),
		Description:      fmt.Sprintf("Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseScalingPoliciesString)),
	},
	"auto_suspend": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend"),
		Description:      "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Default:          IntDefault,
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")),
		Description:      "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the warehouse.",
	},
	"enable_query_acceleration": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enable_query_acceleration"),
		Description:      booleanStringFieldDescription("Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources."),
		Default:          BooleanDefault,
	},
	"query_acceleration_max_scale_factor": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("query_acceleration_max_scale_factor"),
		Description:      "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
		Default:          IntDefault,
	},
	strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Object parameter that specifies the concurrency level for SQL statements (i.e. queries and DML) executed by a warehouse.",
	},
	strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Object parameter that specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on a warehouse before it is canceled by the system.",
	},
	strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 604800)),
		Description:      "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system",
	},
	showOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSE` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseSchema,
		},
	},
	parametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchema,
		},
	},
}

// TODO: merge with DatabaseParametersCustomDiff and extract common
var warehouseParametersCustomDiff = func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	if d.Id() == "" {
		return nil
	}

	client := meta.(*provider.Context).Client
	params, err := client.Parameters.ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Warehouse: helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier),
		},
	})
	if err != nil {
		return err
	}

	return customdiff.All(
		IntParameterValueComputedIf("max_concurrency_level", params, sdk.ParameterTypeWarehouse, sdk.AccountParameterMaxConcurrencyLevel),
		IntParameterValueComputedIf("statement_queued_timeout_in_seconds", params, sdk.ParameterTypeWarehouse, sdk.AccountParameterStatementQueuedTimeoutInSeconds),
		IntParameterValueComputedIf("statement_timeout_in_seconds", params, sdk.ParameterTypeWarehouse, sdk.AccountParameterStatementTimeoutInSeconds),
	)(ctx, d, meta)
}

// Warehouse returns a pointer to the resource representing a warehouse.
func Warehouse() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateWarehouse,
		UpdateContext: UpdateWarehouse,
		ReadContext:   GetReadWarehouseFunc(true),
		DeleteContext: DeleteWarehouse,
		Description:   "Resource used to manage warehouse objects. For more information, check [warehouse documentation](https://docs.snowflake.com/en/sql-reference/commands-warehouse).",

		Schema: warehouseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportWarehouse,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(showOutputAttributeName, "warehouse_type", "warehouse_size", "max_cluster_count", "min_cluster_count", "scaling_policy", "auto_suspend", "auto_resume", "resource_monitor", "comment", "enable_query_acceleration", "query_acceleration_max_scale_factor"),
			ComputedIfAnyAttributeChanged(parametersAttributeName, strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel)), strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)), strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds))),
			customdiff.ForceNewIfChange("warehouse_size", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != "" && new.(string) == ""
			}),
			warehouseParametersCustomDiff,
		),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092WarehouseSizeStateUpgrader,
			},
		},
	}
}

func ImportWarehouse(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting warehouse import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", w.Name); err != nil {
		return nil, err
	}
	if err = d.Set("warehouse_type", w.Type); err != nil {
		return nil, err
	}
	if err = d.Set("warehouse_size", w.Size); err != nil {
		return nil, err
	}
	if err = d.Set("max_cluster_count", w.MaxClusterCount); err != nil {
		return nil, err
	}
	if err = d.Set("min_cluster_count", w.MinClusterCount); err != nil {
		return nil, err
	}
	if err = d.Set("scaling_policy", w.ScalingPolicy); err != nil {
		return nil, err
	}
	if err = d.Set("auto_suspend", w.AutoSuspend); err != nil {
		return nil, err
	}
	if err = d.Set("auto_resume", booleanStringFromBool(w.AutoResume)); err != nil {
		return nil, err
	}
	if err = d.Set("resource_monitor", w.ResourceMonitor.Name()); err != nil {
		return nil, err
	}
	if err = d.Set("comment", w.Comment); err != nil {
		return nil, err
	}
	if err = d.Set("enable_query_acceleration", booleanStringFromBool(w.EnableQueryAcceleration)); err != nil {
		return nil, err
	}
	if err = d.Set("query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// CreateWarehouse implements schema.CreateFunc.
func CreateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	createOptions := &sdk.CreateWarehouseOptions{}

	if v, ok := d.GetOk("warehouse_type"); ok {
		warehouseType, err := sdk.ToWarehouseType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.WarehouseType = &warehouseType
	}
	if v, ok := d.GetOk("warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return diag.FromErr(err)
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
		scalingPolicy, err := sdk.ToScalingPolicy(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.ScalingPolicy = &scalingPolicy
	}
	if v := d.Get("auto_suspend").(int); v != IntDefault {
		createOptions.AutoSuspend = sdk.Int(v)
	}
	if v := d.Get("auto_resume").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.AutoResume = sdk.Bool(parsed)
	}
	if v, ok := d.GetOk("initially_suspended"); ok {
		createOptions.InitiallySuspended = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("resource_monitor"); ok {
		createOptions.ResourceMonitor = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	if v, ok := d.GetOk("comment"); ok {
		createOptions.Comment = sdk.String(v.(string))
	}
	if v := d.Get("enable_query_acceleration").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.EnableQueryAcceleration = sdk.Bool(parsed)
	}
	if v := d.Get("query_acceleration_max_scale_factor").(int); v != IntDefault {
		createOptions.QueryAccelerationMaxScaleFactor = sdk.Int(v)
	}
	if v := GetPropertyAsPointerWithPossibleZeroValues[int](d, "max_concurrency_level"); v != nil {
		createOptions.MaxConcurrencyLevel = v
	}
	if v := GetPropertyAsPointerWithPossibleZeroValues[int](d, "statement_queued_timeout_in_seconds"); v != nil {
		createOptions.StatementQueuedTimeoutInSeconds = v
	}
	if v := GetPropertyAsPointerWithPossibleZeroValues[int](d, "statement_timeout_in_seconds"); v != nil {
		createOptions.StatementTimeoutInSeconds = v
	}

	err := client.Warehouses.Create(ctx, id, createOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return GetReadWarehouseFunc(false)(ctx, d, meta)
}

// TODO: move
func GetPropertyAsPointerWithPossibleZeroValues[T any](d *schema.ResourceData, property string) *T {
	if d.GetRawConfig().AsValueMap()[property].IsNull() {
		return nil
	}
	value := d.Get(property)
	typedValue, ok := value.(T)
	if !ok {
		return nil
	}
	return &typedValue
}

func GetReadWarehouseFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

		w, err := client.Warehouses.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query warehouse. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Warehouse: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		warehouseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{
				Warehouse: id,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"type", "warehouse_type", string(w.Type), w.Type, nil},
				showMapping{"size", "warehouse_size", string(w.Size), w.Size, nil},
				showMapping{"max_cluster_count", "max_cluster_count", w.MaxClusterCount, w.MaxClusterCount, nil},
				showMapping{"min_cluster_count", "min_cluster_count", w.MinClusterCount, w.MinClusterCount, nil},
				showMapping{"scaling_policy", "scaling_policy", string(w.ScalingPolicy), w.ScalingPolicy, nil},
				showMapping{"auto_suspend", "auto_suspend", w.AutoSuspend, w.AutoSuspend, nil},
				showMapping{"auto_resume", "auto_resume", w.AutoResume, fmt.Sprintf("%t", w.AutoResume), nil},
				showMapping{"resource_monitor", "resource_monitor", w.ResourceMonitor.Name(), w.ResourceMonitor.Name(), nil},
				showMapping{"enable_query_acceleration", "enable_query_acceleration", w.EnableQueryAcceleration, fmt.Sprintf("%t", w.EnableQueryAcceleration), nil},
				showMapping{"query_acceleration_max_scale_factor", "query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor, w.QueryAccelerationMaxScaleFactor, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		// These are all identity sets, needed for the case where:
		// - previous config was empty (therefore Snowflake defaults had been used)
		// - new config have the same values that are already in SF
		if !d.GetRawConfig().IsNull() {
			if v := d.GetRawConfig().AsValueMap()["warehouse_type"]; !v.IsNull() {
				if err = d.Set("warehouse_type", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["warehouse_size"]; !v.IsNull() {
				if err = d.Set("warehouse_size", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["max_cluster_count"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("max_cluster_count", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["min_cluster_count"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("min_cluster_count", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["scaling_policy"]; !v.IsNull() {
				if err = d.Set("scaling_policy", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["auto_suspend"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("auto_suspend", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["auto_resume"]; !v.IsNull() {
				if err = d.Set("auto_resume", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["resource_monitor"]; !v.IsNull() {
				if err = d.Set("resource_monitor", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["enable_query_acceleration"]; !v.IsNull() {
				if err = d.Set("enable_query_acceleration", v.AsString()); err != nil {
					return diag.FromErr(err)
				}
			}
			if v := d.GetRawConfig().AsValueMap()["query_acceleration_max_scale_factor"]; !v.IsNull() {
				intVal, _ := v.AsBigFloat().Int64()
				if err = d.Set("query_acceleration_max_scale_factor", intVal); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = d.Set("name", w.Name); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("comment", w.Comment); err != nil {
			return diag.FromErr(err)
		}

		if diags := handleWarehouseParameterRead(d, warehouseParameters); diags != nil {
			return diags
		}

		if err = d.Set(showOutputAttributeName, []map[string]any{schemas.WarehouseToSchema(w)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(parametersAttributeName, []map[string]any{schemas.WarehouseParametersToSchema(warehouseParameters)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func handleWarehouseParameterRead(d *schema.ResourceData, warehouseParameters []*sdk.Parameter) diag.Diagnostics {
	for _, parameter := range warehouseParameters {
		switch parameter.Key {
		case
			string(sdk.ObjectParameterMaxConcurrencyLevel),
			string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds),
			string(sdk.ObjectParameterStatementTimeoutInSeconds):
			value, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

// UpdateWarehouse implements schema.UpdateFunc.
func UpdateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	// Change name separately
	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	// Batch SET operations and UNSET operations
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}
	if d.HasChange("warehouse_type") {
		if v, ok := d.GetOk("warehouse_type"); ok {
			warehouseType, err := sdk.ToWarehouseType(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WarehouseType = &warehouseType
		} else {
			// TODO [SNOW-1473453]: UNSET of type does not work
			// unset.WarehouseType = sdk.Bool(true)
			set.WarehouseType = &sdk.WarehouseTypeStandard
		}
	}
	if d.HasChange("warehouse_size") {
		n := d.Get("warehouse_size").(string)
		size, err := sdk.ToWarehouseSize(n)
		if err != nil {
			return diag.FromErr(err)
		}
		set.WarehouseSize = &size
		// For now, we always want to wait for the resize completion. In the future, we may parametrize it.
		set.WaitForCompletion = sdk.Bool(true)
	}
	if d.HasChange("max_cluster_count") {
		if v, ok := d.GetOk("max_cluster_count"); ok {
			set.MaxClusterCount = sdk.Int(v.(int))
		} else {
			unset.MaxClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("min_cluster_count") {
		if v, ok := d.GetOk("min_cluster_count"); ok {
			set.MinClusterCount = sdk.Int(v.(int))
		} else {
			unset.MinClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("scaling_policy") {
		if v, ok := d.GetOk("scaling_policy"); ok {
			scalingPolicy, err := sdk.ToScalingPolicy(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.ScalingPolicy = &scalingPolicy
		} else {
			// TODO [SNOW-1473453]: UNSET of scaling policy does not work
			// unset.ScalingPolicy = sdk.Bool(true)
			set.ScalingPolicy = &sdk.ScalingPolicyStandard
		}
	}
	if d.HasChange("auto_suspend") {
		if v := d.Get("auto_suspend").(int); v != IntDefault {
			set.AutoSuspend = sdk.Int(v)
		} else {
			// TODO [SNOW-1473453]: UNSET of auto suspend works incorrectly
			// unset.AutoSuspend = sdk.Bool(true)
			set.AutoSuspend = sdk.Int(600)
		}
	}
	if d.HasChange("auto_resume") {
		if v := d.Get("auto_resume").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.AutoResume = sdk.Bool(parsed)
		} else {
			// TODO [SNOW-1473453]: UNSET of auto resume works incorrectly
			// unset.AutoResume = sdk.Bool(true)
			set.AutoResume = sdk.Bool(true)
		}
	}
	if d.HasChange("resource_monitor") {
		if v, ok := d.GetOk("resource_monitor"); ok {
			set.ResourceMonitor = sdk.NewAccountObjectIdentifier(v.(string))
		} else {
			unset.ResourceMonitor = sdk.Bool(true)
		}
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.Comment = sdk.String(v.(string))
		} else {
			unset.Comment = sdk.Bool(true)
		}
	}
	if d.HasChange("enable_query_acceleration") {
		if v := d.Get("enable_query_acceleration").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.EnableQueryAcceleration = sdk.Bool(parsed)
		} else {
			unset.EnableQueryAcceleration = sdk.Bool(true)
		}
	}
	if d.HasChange("query_acceleration_max_scale_factor") {
		if v := d.Get("query_acceleration_max_scale_factor").(int); v != IntDefault {
			set.QueryAccelerationMaxScaleFactor = sdk.Int(v)
		} else {
			unset.QueryAccelerationMaxScaleFactor = sdk.Bool(true)
		}
	}
	if d.HasChange("max_concurrency_level") {
		if v := d.Get("max_concurrency_level").(int); v != IntDefault {
			set.MaxConcurrencyLevel = sdk.Int(v)
		} else {
			unset.MaxConcurrencyLevel = sdk.Bool(true)
		}
	}
	if d.HasChange("statement_queued_timeout_in_seconds") {
		if v := d.Get("statement_queued_timeout_in_seconds").(int); v != IntDefault {
			set.StatementQueuedTimeoutInSeconds = sdk.Int(v)
		} else {
			unset.StatementQueuedTimeoutInSeconds = sdk.Bool(true)
		}
	}

	if updateParamDiags := handleWarehouseParametersChanges(d, &set, &unset); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	// Apply SET and UNSET changes
	if (set != sdk.WarehouseSet{}) {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Set: &set,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if (unset != sdk.WarehouseUnset{}) {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Unset: &unset,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadWarehouseFunc(false)(ctx, d, meta)
}

func handleWarehouseParametersChanges(d *schema.ResourceData, set *sdk.WarehouseSet, unset *sdk.WarehouseUnset) diag.Diagnostics {
	return JoinDiags(
		handleValuePropertyChange[int](d, "max_concurrency_level", &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		handleValuePropertyChange[int](d, "statement_queued_timeout_in_seconds", &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleValuePropertyChange[int](d, "statement_timeout_in_seconds", &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	)
}

// DeleteWarehouse implements schema.DeleteFunc.
func DeleteWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Warehouses.Drop(ctx, id, &sdk.DropWarehouseOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
