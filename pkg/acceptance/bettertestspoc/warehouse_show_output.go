package bettertestspoc

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehouseShowOutputAssert struct {
	*ResourceAssert
}

func WarehouseShowOutput(t *testing.T, name string) *WarehouseShowOutputAssert {
	t.Helper()
	w := WarehouseShowOutputAssert{
		NewResourceAssert(name, "show_output"),
	}
	w.assertions = append(w.assertions, valueSet("show_output.#", "1"))
	return &w
}

func ImportedWarehouseShowOutput(t *testing.T, id string) *WarehouseShowOutputAssert {
	t.Helper()
	w := WarehouseShowOutputAssert{
		NewImportedResourceAssert(id, "show_output"),
	}
	w.assertions = append(w.assertions, valueSet("show_output.#", "1"))
	return &w
}

func (w *WarehouseShowOutputAssert) HasType(expected sdk.WarehouseType) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("type", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasSize(expected sdk.WarehouseSize) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("size", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasMaxClusterCount(expected int) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("max_cluster_count", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasMinClusterCount(expected int) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("min_cluster_count", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasScalingPolicy(expected sdk.ScalingPolicy) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("scaling_policy", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoSuspend(expected int) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("auto_suspend", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoResume(expected bool) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("auto_resume", strconv.FormatBool(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasResourceMonitor(expected string) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("resource_monitor", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasComment(expected string) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("comment", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasEnableQueryAcceleration(expected bool) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("enable_query_acceleration", strconv.FormatBool(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasQueryAccelerationMaxScaleFactor(expected int) *WarehouseShowOutputAssert {
	w.assertions = append(w.assertions, showOutputValueSet("query_acceleration_max_scale_factor", strconv.Itoa(expected)))
	return w
}
