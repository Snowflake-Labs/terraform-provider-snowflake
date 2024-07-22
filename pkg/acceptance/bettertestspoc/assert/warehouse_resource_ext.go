package assert

import r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

func (w *WarehouseResourceAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceAssert {
	w.AddAssertion(ValueSet("max_concurrency_level", "8"))
	return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseResourceAssert {
	w.AddAssertion(ValueSet("statement_queued_timeout_in_seconds", "0"))
	return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseResourceAssert {
	w.AddAssertion(ValueSet("statement_timeout_in_seconds", "172800"))
	return w
}

func (w *WarehouseResourceAssert) HasAllDefault() *WarehouseResourceAssert {
	return w.HasDefaultMaxConcurrencyLevel().
		HasNoType().
		HasNoSize().
		HasNoMaxClusterCount().
		HasNoMinClusterCount().
		HasNoScalingPolicy().
		HasAutoSuspend(r.IntDefaultString).
		HasAutoResume(r.BooleanDefault).
		HasNoInitiallySuspended().
		HasNoResourceMonitor().
		HasEnableQueryAcceleration(r.BooleanDefault).
		HasQueryAccelerationMaxScaleFactor(r.IntDefaultString).
		HasDefaultMaxConcurrencyLevel().
		HasDefaultStatementQueuedTimeoutInSeconds().
		HasDefaultStatementTimeoutInSeconds()
}
