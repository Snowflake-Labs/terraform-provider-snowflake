package assert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehouseParametersAssert struct {
	*ResourceAssert
}

func WarehouseParameters(t *testing.T, name string) *WarehouseParametersAssert {
	t.Helper()
	w := WarehouseParametersAssert{
		NewResourceAssert(name, "parameters"),
	}
	w.AddAssertion(ValueSet("parameters.#", "1"))
	return &w
}

func ImportedWarehouseParameters(t *testing.T, id string) *WarehouseParametersAssert {
	t.Helper()
	w := WarehouseParametersAssert{
		NewImportedResourceAssert(id, "imported parameters"),
	}
	w.AddAssertion(ValueSet("parameters.#", "1"))
	return &w
}

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevel(expected int) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterIntValueSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterIntValueSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterIntValueSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevelLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterLevelSet(sdk.WarehouseParameterMaxConcurrencyLevel, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterLevelSet(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.AddAssertion(ResourceParameterLevelSet(sdk.WarehouseParameterStatementTimeoutInSeconds, expected))
	return w
}
