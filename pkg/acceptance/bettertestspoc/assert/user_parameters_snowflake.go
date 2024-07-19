package assert

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO: make assertions naming consistent (resource paramaters vs snowflake parameters)
type UserParametersAssert struct {
	*SnowflakeParametersAssert[sdk.AccountObjectIdentifier]
}

func UserParameters(t *testing.T, id sdk.AccountObjectIdentifier) *UserParametersAssert {
	t.Helper()
	return &UserParametersAssert{
		NewSnowflakeParametersAssertWithProvider(id, sdk.ObjectTypeUser, acc.TestClient().Parameter.ShowUserParameters),
	}
}

func UserParametersPrefetched(t *testing.T, id sdk.AccountObjectIdentifier, parameters []*sdk.Parameter) *UserParametersAssert {
	t.Helper()
	return &UserParametersAssert{
		NewSnowflakeParametersAssertWithParameters(id, sdk.ObjectTypeUser, parameters),
	}
}

func (w *UserParametersAssert) HasEnableUnredactedQuerySyntaxError(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, expected))
	return w
}

func (w *UserParametersAssert) HasNetworkPolicy(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterNetworkPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasPreventUnloadToInternalStages(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterPreventUnloadToInternalStages, expected))
	return w
}

func (w *UserParametersAssert) HasAbortDetachedQuery(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterAbortDetachedQuery, expected))
	return w
}

func (w *UserParametersAssert) HasAutocommit(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterAutocommit, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryInputFormat(expected sdk.BinaryInputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterBinaryInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryOutputFormat(expected sdk.BinaryOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterBinaryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasClientMemoryLimit(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterClientMemoryLimit, expected))
	return w
}

func (w *UserParametersAssert) HasClientMetadataRequestUseConnectionCtx(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterClientMetadataRequestUseConnectionCtx, expected))
	return w
}

func (w *UserParametersAssert) HasClientPrefetchThreads(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterClientPrefetchThreads, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultChunkSize(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterClientResultChunkSize, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultColumnCaseInsensitive(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterClientResultColumnCaseInsensitive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAlive(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterClientSessionKeepAlive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAliveHeartbeatFrequency(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, expected))
	return w
}

func (w *UserParametersAssert) HasClientTimestampTypeMapping(expected sdk.ClientTimestampTypeMapping) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterClientTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasDateInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterDateInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasDateOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterDateOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasEnableUnloadPhysicalTypeOptimization(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterEnableUnloadPhysicalTypeOptimization, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicMerge(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterErrorOnNondeterministicMerge, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicUpdate(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterErrorOnNondeterministicUpdate, expected))
	return w
}

func (w *UserParametersAssert) HasGeographyOutputFormat(expected sdk.GeographyOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterGeographyOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasGeometryOutputFormat(expected sdk.GeometryOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterGeometryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatDecimalAsInt(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterJdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatTimestampNtzAsUtc(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterJdbcTreatTimestampNtzAsUtc, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcUseSessionTimezone(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterJdbcUseSessionTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasJsonIndent(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterJsonIndent, expected))
	return w
}

func (w *UserParametersAssert) HasLockTimeout(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterLockTimeout, expected))
	return w
}

func (w *UserParametersAssert) HasLogLevel(expected sdk.LogLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterLogLevel, expected))
	return w
}

func (w *UserParametersAssert) HasMultiStatementCount(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterMultiStatementCount, expected))
	return w
}

func (w *UserParametersAssert) HasNoorderSequenceAsDefault(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterNoorderSequenceAsDefault, expected))
	return w
}

func (w *UserParametersAssert) HasOdbcTreatDecimalAsInt(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterOdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasQueryTag(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterQueryTag, expected))
	return w
}

func (w *UserParametersAssert) HasQuotedIdentifiersIgnoreCase(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterQuotedIdentifiersIgnoreCase, expected))
	return w
}

func (w *UserParametersAssert) HasRowsPerResultset(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterRowsPerResultset, expected))
	return w
}

func (w *UserParametersAssert) HasS3StageVpceDnsName(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterS3StageVpceDnsName, expected))
	return w
}

func (w *UserParametersAssert) HasSearchPath(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterSearchPath, expected))
	return w
}

func (w *UserParametersAssert) HasSimulatedDataSharingConsumer(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterSimulatedDataSharingConsumer, expected))
	return w
}

func (w *UserParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStatementTimeoutInSeconds(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterStatementTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStrictJsonOutput(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterStrictJsonOutput, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampDayIsAlways24h(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterTimestampDayIsAlways24h, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimestampInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampLtzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimestampLtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampNtzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimestampNtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimestampOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTypeMapping(expected sdk.TimestampTypeMapping) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimestampTzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimezone(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasTimeInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimeInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimeOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterTimeOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTraceLevel(expected sdk.TraceLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTraceLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionAbortOnError(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterTransactionAbortOnError, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionDefaultIsolationLevel(expected sdk.TransactionDefaultIsolationLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTransactionDefaultIsolationLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTwoDigitCenturyStart(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterTwoDigitCenturyStart, expected))
	return w
}

// lowercase for ignore in snowflake by default but uppercase for FAIL
func (w *UserParametersAssert) HasUnsupportedDdlAction(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterValueSet(sdk.UserParameterUnsupportedDdlAction, expected))
	return w
}

func (w *UserParametersAssert) HasUseCachedResult(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterBoolValueSet(sdk.UserParameterUseCachedResult, expected))
	return w
}

func (w *UserParametersAssert) HasWeekOfYearPolicy(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterWeekOfYearPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasWeekStart(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, snowflakeParameterIntValueSet(sdk.UserParameterWeekStart, expected))
	return w
}
