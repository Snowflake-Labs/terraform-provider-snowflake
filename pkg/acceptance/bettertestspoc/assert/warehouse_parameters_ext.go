package assert

func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}

func (w *WarehouseParametersAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseParametersAssert {
	return w.
		HasStatementQueuedTimeoutInSeconds(0).
		HasStatementQueuedTimeoutInSecondsLevel("")
}

func (w *WarehouseParametersAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseParametersAssert {
	return w.
		HasStatementTimeoutInSeconds(172800).
		HasStatementTimeoutInSecondsLevel("")
}
