package assert

import (
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
	w.AddAssertion(ResourceShowOutputStringUnderlyingValueSet("type", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasSize(expected sdk.WarehouseSize) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputStringUnderlyingValueSet("size", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasMaxClusterCount(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputIntValueSet("max_cluster_count", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasMinClusterCount(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputIntValueSet("min_cluster_count", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasScalingPolicy(expected sdk.ScalingPolicy) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputStringUnderlyingValueSet("scaling_policy", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoSuspend(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputIntValueSet("auto_suspend", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasAutoResume(expected bool) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputBoolValueSet("auto_resume", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasResourceMonitor(expected string) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputValueSet("resource_monitor", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasComment(expected string) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputValueSet("comment", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasEnableQueryAcceleration(expected bool) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputBoolValueSet("enable_query_acceleration", expected))
	return w
}

func (w *WarehouseShowOutputAssert) HasQueryAccelerationMaxScaleFactor(expected int) *WarehouseShowOutputAssert {
	w.AddAssertion(ResourceShowOutputIntValueSet("query_acceleration_max_scale_factor", expected))
	return w
}
