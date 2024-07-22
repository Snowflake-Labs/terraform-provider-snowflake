package assert

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
	w.AddAssertion(ValueSet("show_output.#", "1"))
	return &w
}

func ImportedWarehouseShowOutput(t *testing.T, id string) *WarehouseShowOutputAssert {
	t.Helper()
	w := WarehouseShowOutputAssert{
		NewImportedResourceAssert(id, "show_output"),
	}
	w.AddAssertion(ValueSet("show_output.#", "1"))
	return &w
}

func (w *WarehouseShowOutputAssert) HasType(expected sdk.WarehouseType) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("type", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasSize(expected sdk.WarehouseSize) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("size", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasMaxClusterCount(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("max_cluster_count", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasMinClusterCount(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("min_cluster_count", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasScalingPolicy(expected sdk.ScalingPolicy) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("scaling_policy", string(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoSuspend(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("auto_suspend", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoResume(expected bool) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("auto_resume", strconv.FormatBool(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasResourceMonitor(expected string) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("resource_monitor", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasComment(expected string) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("comment", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasEnableQueryAcceleration(expected bool) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("enable_query_acceleration", strconv.FormatBool(expected)))
	return w
}

func (w *WarehouseShowOutputAssert) HasQueryAccelerationMaxScaleFactor(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(showOutputValueSet("query_acceleration_max_scale_factor", strconv.Itoa(expected)))
	return w
}
