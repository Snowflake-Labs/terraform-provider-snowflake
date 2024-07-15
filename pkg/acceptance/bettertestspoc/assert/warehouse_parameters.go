package assert

import (
	"strconv"
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
	w.assertions = append(w.assertions, valueSet("parameters.#", "1"))
	return &w
}

func ImportedWarehouseParameters(t *testing.T, id string) *WarehouseParametersAssert {
	t.Helper()
	w := WarehouseParametersAssert{
		NewImportedResourceAssert(id, "imported parameters"),
	}
	w.assertions = append(w.assertions, valueSet("parameters.#", "1"))
	return &w
}

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevel(expected int) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterValueSet("max_concurrency_level", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterValueSet("statement_queued_timeout_in_seconds", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSeconds(expected int) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterValueSet("statement_timeout_in_seconds", strconv.Itoa(expected)))
	return w
}

func (w *WarehouseParametersAssert) HasMaxConcurrencyLevelLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterLevelSet("max_concurrency_level", string(expected)))
	return w
}

func (w *WarehouseParametersAssert) HasStatementQueuedTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterLevelSet("statement_queued_timeout_in_seconds", string(expected)))
	return w
}

func (w *WarehouseParametersAssert) HasStatementTimeoutInSecondsLevel(expected sdk.ParameterType) *WarehouseParametersAssert {
	w.assertions = append(w.assertions, parameterLevelSet("statement_timeout_in_seconds", string(expected)))
	return w
}
