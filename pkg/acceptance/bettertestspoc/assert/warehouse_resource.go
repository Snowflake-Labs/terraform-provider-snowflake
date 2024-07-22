package assert

import (
	"testing"
)

type WarehouseResourceAssert struct {
	*ResourceAssert
}

func WarehouseResource(t *testing.T, name string) *WarehouseResourceAssert {
	t.Helper()

	return &WarehouseResourceAssert{
		ResourceAssert: NewResourceAssert(name, "resource"),
	}
}

func ImportedWarehouseResource(t *testing.T, id string) *WarehouseResourceAssert {
	t.Helper()

	return &WarehouseResourceAssert{
		ResourceAssert: NewImportedResourceAssert(id, "imported resource"),
	}
}

func (w *WarehouseResourceAssert) HasName(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("name", expected))
	return w
}

func (w *WarehouseResourceAssert) HasType(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("warehouse_type", expected))
	return w
}

func (w *WarehouseResourceAssert) HasSize(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("warehouse_size", expected))
	return w
}

func (w *WarehouseResourceAssert) HasMaxClusterCount(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("max_cluster_count", expected))
	return w
}

func (w *WarehouseResourceAssert) HasMinClusterCount(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("min_cluster_count", expected))
	return w
}

func (w *WarehouseResourceAssert) HasScalingPolicy(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("scaling_policy", expected))
	return w
}

func (w *WarehouseResourceAssert) HasAutoSuspend(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("auto_suspend", expected))
	return w
}

func (w *WarehouseResourceAssert) HasAutoResume(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("auto_resume", expected))
	return w
}

func (w *WarehouseResourceAssert) HasInitiallySuspended(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("initially_suspended", expected))
	return w
}

func (w *WarehouseResourceAssert) HasResourceMonitor(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("resource_monitor", expected))
	return w
}

func (w *WarehouseResourceAssert) HasComment(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("comment", expected))
	return w
}

func (w *WarehouseResourceAssert) HasEnableQueryAcceleration(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("enable_query_acceleration", expected))
	return w
}

func (w *WarehouseResourceAssert) HasQueryAccelerationMaxScaleFactor(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("query_acceleration_max_scale_factor", expected))
	return w
}

func (w *WarehouseResourceAssert) HasMaxConcurrencyLevel(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("max_concurrency_level", expected))
	return w
}

func (w *WarehouseResourceAssert) HasStatementQueuedTimeoutInSeconds(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("statement_queued_timeout_in_seconds", expected))
	return w
}

func (w *WarehouseResourceAssert) HasStatementTimeoutInSeconds(expected string) *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueSet("statement_timeout_in_seconds", expected))
	return w
}

func (w *WarehouseResourceAssert) HasNoName() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("name"))
	return w
}

func (w *WarehouseResourceAssert) HasNoType() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("warehouse_type"))
	return w
}

func (w *WarehouseResourceAssert) HasNoSize() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("warehouse_size"))
	return w
}

func (w *WarehouseResourceAssert) HasNoMaxClusterCount() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("max_cluster_count"))
	return w
}

func (w *WarehouseResourceAssert) HasNoMinClusterCount() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("min_cluster_count"))
	return w
}

func (w *WarehouseResourceAssert) HasNoScalingPolicy() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("scaling_policy"))
	return w
}

func (w *WarehouseResourceAssert) HasNoAutoSuspend() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("auto_suspend"))
	return w
}

func (w *WarehouseResourceAssert) HasNoAutoResume() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("auto_resume"))
	return w
}

func (w *WarehouseResourceAssert) HasNoInitiallySuspended() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("initially_suspended"))
	return w
}

func (w *WarehouseResourceAssert) HasNoResourceMonitor() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("resource_monitor"))
	return w
}

func (w *WarehouseResourceAssert) HasNoComment() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("comment"))
	return w
}

func (w *WarehouseResourceAssert) HasNoEnableQueryAcceleration() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("enable_query_acceleration"))
	return w
}

func (w *WarehouseResourceAssert) HasNoQueryAccelerationMaxScaleFactor() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("query_acceleration_max_scale_factor"))
	return w
}

func (w *WarehouseResourceAssert) HasNoMaxConcurrencyLevel() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("max_concurrency_level"))
	return w
}

func (w *WarehouseResourceAssert) HasNoStatementQueuedTimeoutInSeconds() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("statement_queued_timeout_in_seconds"))
	return w
}

func (w *WarehouseResourceAssert) HasNoStatementTimeoutInSeconds() *WarehouseResourceAssert {
	w.assertions = append(w.assertions, ValueNotSet("statement_timeout_in_seconds"))
	return w
}
