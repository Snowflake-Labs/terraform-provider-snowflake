package resourceparametersassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (f *ProcedureResourceParametersAssert) HasAllDefaults() *ProcedureResourceParametersAssert {
	return f.
		HasEnableConsoleOutput(false).
		HasLogLevel(sdk.LogLevelOff).
		HasMetricLevel(sdk.MetricLevelNone).
		HasTraceLevel(sdk.TraceLevelOff)
}
