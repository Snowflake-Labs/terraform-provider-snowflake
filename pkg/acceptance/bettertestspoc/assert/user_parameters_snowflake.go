package assert

import (
	"strings"
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

// TODO: try to move this section to SnowflakeParametersAssert to not copy it for every object; persist the type-safe assertions
//////////////////////////////
// Generic parameter checks //
//////////////////////////////

func (w *UserParametersAssert) HasBoolParameterValue(parameterName sdk.UserParameter, expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(parameterName, expected))
	return w
}

func (w *UserParametersAssert) HasIntParameterValue(parameterName sdk.UserParameter, expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(parameterName, expected))
	return w
}

func (w *UserParametersAssert) HasStringParameterValue(parameterName sdk.UserParameter, expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(parameterName, expected))
	return w
}

func (w *UserParametersAssert) HasDefaultParameterValue(parameterName sdk.UserParameter) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterDefaultValueSet(parameterName))
	return w
}

func (w *UserParametersAssert) HasDefaultParameterValueOnLevel(parameterName sdk.UserParameter, parameterType sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterDefaultValueOnLevelSet(parameterName, parameterType))
	return w
}

///////////////////////////////
// Aggregated generic checks //
///////////////////////////////

// HasAllDefaults checks if all the parameters:
// - have a default value by comparing current value of the sdk.Parameter with its default
// - have an expected level
func (w *UserParametersAssert) HasAllDefaults() *UserParametersAssert {
	return w.
		HasDefaultParameterValueOnLevel(sdk.UserParameterEnableUnredactedQuerySyntaxError, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterNetworkPolicy, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterPreventUnloadToInternalStages, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterAbortDetachedQuery, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterAutocommit, sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnLevel(sdk.UserParameterBinaryInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterBinaryOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientMemoryLimit, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientMetadataRequestUseConnectionCtx, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientPrefetchThreads, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientResultChunkSize, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientResultColumnCaseInsensitive, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientSessionKeepAlive, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientTimestampTypeMapping, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterDateInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterDateOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterEnableUnloadPhysicalTypeOptimization, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterErrorOnNondeterministicMerge, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterErrorOnNondeterministicUpdate, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterGeographyOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterGeometryOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcTreatDecimalAsInt, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcTreatTimestampNtzAsUtc, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcUseSessionTimezone, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJsonIndent, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterLockTimeout, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterLogLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterMultiStatementCount, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterNoorderSequenceAsDefault, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterOdbcTreatDecimalAsInt, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterQueryTag, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterQuotedIdentifiersIgnoreCase, sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnLevel(sdk.UserParameterRowsPerResultset, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterS3StageVpceDnsName, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterSearchPath, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterSimulatedDataSharingConsumer, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStatementQueuedTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStatementTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStrictJsonOutput, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampDayIsAlways24h, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampLtzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampNtzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampTypeMapping, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampTzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimezone, sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimeInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimeOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTraceLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTransactionAbortOnError, sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTransactionDefaultIsolationLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTwoDigitCenturyStart, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterUnsupportedDdlAction, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterUseCachedResult, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterWeekOfYearPolicy, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterWeekStart, sdk.ParameterTypeSnowflakeDefault)
}

func (w *UserParametersAssert) HasAllDefaultsExplicit() *UserParametersAssert {
	return w.
		HasDefaultEnableUnredactedQuerySyntaxErrorValueExplicit().
		HasDefaultNetworkPolicyValueExplicit().
		HasDefaultPreventUnloadToInternalStagesValueExplicit().
		HasDefaultAbortDetachedQueryValueExplicit().
		HasDefaultAutocommitValueExplicit().
		HasDefaultBinaryInputFormatValueExplicit().
		HasDefaultBinaryOutputFormatValueExplicit().
		HasDefaultClientMemoryLimitValueExplicit().
		HasDefaultClientMetadataRequestUseConnectionCtxValueExplicit().
		HasDefaultClientPrefetchThreadsValueExplicit().
		HasDefaultClientResultChunkSizeValueExplicit().
		HasDefaultClientResultColumnCaseInsensitiveValueExplicit().
		HasDefaultClientSessionKeepAliveValueExplicit().
		HasDefaultClientSessionKeepAliveHeartbeatFrequencyValueExplicit().
		HasDefaultClientTimestampTypeMappingValueExplicit().
		HasDefaultDateInputFormatValueExplicit().
		HasDefaultDateOutputFormatValueExplicit().
		HasDefaultEnableUnloadPhysicalTypeOptimizationValueExplicit().
		HasDefaultErrorOnNondeterministicMergeValueExplicit().
		HasDefaultErrorOnNondeterministicUpdateValueExplicit().
		HasDefaultGeographyOutputFormatValueExplicit().
		HasDefaultGeometryOutputFormatValueExplicit().
		HasDefaultJdbcTreatDecimalAsIntValueExplicit().
		HasDefaultJdbcTreatTimestampNtzAsUtcValueExplicit().
		HasDefaultJdbcUseSessionTimezoneValueExplicit().
		HasDefaultJsonIndentValueExplicit().
		HasDefaultLockTimeoutValueExplicit().
		HasDefaultLogLevelValueExplicit().
		HasDefaultMultiStatementCountValueExplicit().
		HasDefaultNoorderSequenceAsDefaultValueExplicit().
		HasDefaultOdbcTreatDecimalAsIntValueExplicit().
		HasDefaultQueryTagValueExplicit().
		HasDefaultQuotedIdentifiersIgnoreCaseValueExplicit().
		HasDefaultRowsPerResultsetValueExplicit().
		HasDefaultS3StageVpceDnsNameValueExplicit().
		HasDefaultSearchPathValueExplicit().
		HasDefaultSimulatedDataSharingConsumerValueExplicit().
		HasDefaultStatementQueuedTimeoutInSecondsValueExplicit().
		HasDefaultStatementTimeoutInSecondsValueExplicit().
		HasDefaultStrictJsonOutputValueExplicit().
		HasDefaultTimestampDayIsAlways24hValueExplicit().
		HasDefaultTimestampInputFormatValueExplicit().
		HasDefaultTimestampLtzOutputFormatValueExplicit().
		HasDefaultTimestampNtzOutputFormatValueExplicit().
		HasDefaultTimestampOutputFormatValueExplicit().
		HasDefaultTimestampTypeMappingValueExplicit().
		HasDefaultTimestampTzOutputFormatValueExplicit().
		HasDefaultTimezoneValueExplicit().
		HasDefaultTimeInputFormatValueExplicit().
		HasDefaultTimeOutputFormatValueExplicit().
		HasDefaultTraceLevelValueExplicit().
		HasDefaultTransactionAbortOnErrorValueExplicit().
		HasDefaultTransactionDefaultIsolationLevelValueExplicit().
		HasDefaultTwoDigitCenturyStartValueExplicit().
		HasDefaultUnsupportedDdlActionValueExplicit().
		HasDefaultUseCachedResultValueExplicit().
		HasDefaultWeekOfYearPolicyValueExplicit().
		HasDefaultWeekStartValueExplicit()
}

///////////////////////////////
// Specific parameter checks //
///////////////////////////////

// value checks

func (w *UserParametersAssert) HasEnableUnredactedQuerySyntaxError(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, expected))
	return w
}

func (w *UserParametersAssert) HasNetworkPolicy(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterNetworkPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasPreventUnloadToInternalStages(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterPreventUnloadToInternalStages, expected))
	return w
}

func (w *UserParametersAssert) HasAbortDetachedQuery(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterAbortDetachedQuery, expected))
	return w
}

func (w *UserParametersAssert) HasAutocommit(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterAutocommit, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryInputFormat(expected sdk.BinaryInputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterBinaryInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryOutputFormat(expected sdk.BinaryOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterBinaryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasClientMemoryLimit(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterClientMemoryLimit, expected))
	return w
}

func (w *UserParametersAssert) HasClientMetadataRequestUseConnectionCtx(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterClientMetadataRequestUseConnectionCtx, expected))
	return w
}

func (w *UserParametersAssert) HasClientPrefetchThreads(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterClientPrefetchThreads, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultChunkSize(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterClientResultChunkSize, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultColumnCaseInsensitive(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterClientResultColumnCaseInsensitive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAlive(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterClientSessionKeepAlive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAliveHeartbeatFrequency(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, expected))
	return w
}

func (w *UserParametersAssert) HasClientTimestampTypeMapping(expected sdk.ClientTimestampTypeMapping) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterClientTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasDateInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterDateInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasDateOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterDateOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasEnableUnloadPhysicalTypeOptimization(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterEnableUnloadPhysicalTypeOptimization, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicMerge(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterErrorOnNondeterministicMerge, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicUpdate(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterErrorOnNondeterministicUpdate, expected))
	return w
}

func (w *UserParametersAssert) HasGeographyOutputFormat(expected sdk.GeographyOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterGeographyOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasGeometryOutputFormat(expected sdk.GeometryOutputFormat) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterGeometryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatDecimalAsInt(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterJdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatTimestampNtzAsUtc(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterJdbcTreatTimestampNtzAsUtc, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcUseSessionTimezone(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterJdbcUseSessionTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasJsonIndent(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterJsonIndent, expected))
	return w
}

func (w *UserParametersAssert) HasLockTimeout(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterLockTimeout, expected))
	return w
}

func (w *UserParametersAssert) HasLogLevel(expected sdk.LogLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterLogLevel, expected))
	return w
}

func (w *UserParametersAssert) HasMultiStatementCount(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterMultiStatementCount, expected))
	return w
}

func (w *UserParametersAssert) HasNoorderSequenceAsDefault(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterNoorderSequenceAsDefault, expected))
	return w
}

func (w *UserParametersAssert) HasOdbcTreatDecimalAsInt(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterOdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasQueryTag(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterQueryTag, expected))
	return w
}

func (w *UserParametersAssert) HasQuotedIdentifiersIgnoreCase(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterQuotedIdentifiersIgnoreCase, expected))
	return w
}

func (w *UserParametersAssert) HasRowsPerResultset(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterRowsPerResultset, expected))
	return w
}

func (w *UserParametersAssert) HasS3StageVpceDnsName(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterS3StageVpceDnsName, expected))
	return w
}

func (w *UserParametersAssert) HasSearchPath(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterSearchPath, expected))
	return w
}

func (w *UserParametersAssert) HasSimulatedDataSharingConsumer(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterSimulatedDataSharingConsumer, expected))
	return w
}

func (w *UserParametersAssert) HasStatementQueuedTimeoutInSeconds(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStatementTimeoutInSeconds(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterStatementTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStrictJsonOutput(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterStrictJsonOutput, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampDayIsAlways24h(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterTimestampDayIsAlways24h, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimestampInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampLtzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimestampLtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampNtzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimestampNtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimestampOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTypeMapping(expected sdk.TimestampTypeMapping) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTzOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimestampTzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimezone(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasTimeInputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimeInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimeOutputFormat(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterTimeOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTraceLevel(expected sdk.TraceLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTraceLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionAbortOnError(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterTransactionAbortOnError, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionDefaultIsolationLevel(expected sdk.TransactionDefaultIsolationLevel) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterStringUnderlyingValueSet(sdk.UserParameterTransactionDefaultIsolationLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTwoDigitCenturyStart(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterTwoDigitCenturyStart, expected))
	return w
}

// lowercase for ignore in snowflake by default but uppercase for FAIL
func (w *UserParametersAssert) HasUnsupportedDdlAction(expected string) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterValueSet(sdk.UserParameterUnsupportedDdlAction, expected))
	return w
}

func (w *UserParametersAssert) HasUseCachedResult(expected bool) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterBoolValueSet(sdk.UserParameterUseCachedResult, expected))
	return w
}

func (w *UserParametersAssert) HasWeekOfYearPolicy(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterWeekOfYearPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasWeekStart(expected int) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterIntValueSet(sdk.UserParameterWeekStart, expected))
	return w
}

// level checks

func (w *UserParametersAssert) HasEnableUnredactedQuerySyntaxErrorLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, expected))
	return w
}

func (w *UserParametersAssert) HasNetworkPolicyLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterNetworkPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasPreventUnloadToInternalStagesLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterPreventUnloadToInternalStages, expected))
	return w
}

func (w *UserParametersAssert) HasAbortDetachedQueryLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterAbortDetachedQuery, expected))
	return w
}

func (w *UserParametersAssert) HasAutocommitLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterAutocommit, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryInputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterBinaryInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasBinaryOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterBinaryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasClientMemoryLimitLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientMemoryLimit, expected))
	return w
}

func (w *UserParametersAssert) HasClientMetadataRequestUseConnectionCtxLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientMetadataRequestUseConnectionCtx, expected))
	return w
}

func (w *UserParametersAssert) HasClientPrefetchThreadsLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientPrefetchThreads, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultChunkSizeLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientResultChunkSize, expected))
	return w
}

func (w *UserParametersAssert) HasClientResultColumnCaseInsensitiveLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientResultColumnCaseInsensitive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAliveLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientSessionKeepAlive, expected))
	return w
}

func (w *UserParametersAssert) HasClientSessionKeepAliveHeartbeatFrequencyLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, expected))
	return w
}

func (w *UserParametersAssert) HasClientTimestampTypeMappingLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterClientTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasDateInputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterDateInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasDateOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterDateOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasEnableUnloadPhysicalTypeOptimizationLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterEnableUnloadPhysicalTypeOptimization, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicMergeLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterErrorOnNondeterministicMerge, expected))
	return w
}

func (w *UserParametersAssert) HasErrorOnNondeterministicUpdateLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterErrorOnNondeterministicUpdate, expected))
	return w
}

func (w *UserParametersAssert) HasGeographyOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterGeographyOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasGeometryOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterGeometryOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatDecimalAsIntLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterJdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcTreatTimestampNtzAsUtcLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterJdbcTreatTimestampNtzAsUtc, expected))
	return w
}

func (w *UserParametersAssert) HasJdbcUseSessionTimezoneLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterJdbcUseSessionTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasJsonIndentLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterJsonIndent, expected))
	return w
}

func (w *UserParametersAssert) HasLockTimeoutLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterLockTimeout, expected))
	return w
}

func (w *UserParametersAssert) HasLogLevelLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterLogLevel, expected))
	return w
}

func (w *UserParametersAssert) HasMultiStatementCountLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterMultiStatementCount, expected))
	return w
}

func (w *UserParametersAssert) HasNoorderSequenceAsDefaultLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterNoorderSequenceAsDefault, expected))
	return w
}

func (w *UserParametersAssert) HasOdbcTreatDecimalAsIntLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterOdbcTreatDecimalAsInt, expected))
	return w
}

func (w *UserParametersAssert) HasQueryTagLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterQueryTag, expected))
	return w
}

func (w *UserParametersAssert) HasQuotedIdentifiersIgnoreCaseLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterQuotedIdentifiersIgnoreCase, expected))
	return w
}

func (w *UserParametersAssert) HasRowsPerResultsetLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterRowsPerResultset, expected))
	return w
}

func (w *UserParametersAssert) HasS3StageVpceDnsNameLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterS3StageVpceDnsName, expected))
	return w
}

func (w *UserParametersAssert) HasSearchPathLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterSearchPath, expected))
	return w
}

func (w *UserParametersAssert) HasSimulatedDataSharingConsumerLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterSimulatedDataSharingConsumer, expected))
	return w
}

func (w *UserParametersAssert) HasStatementQueuedTimeoutInSecondsLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterStatementQueuedTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStatementTimeoutInSecondsLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterStatementTimeoutInSeconds, expected))
	return w
}

func (w *UserParametersAssert) HasStrictJsonOutputLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterStrictJsonOutput, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampDayIsAlways24hLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampDayIsAlways24h, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampInputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampLtzOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampLtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampNtzOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampNtzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTypeMappingLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampTypeMapping, expected))
	return w
}

func (w *UserParametersAssert) HasTimestampTzOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimestampTzOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimezoneLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimezone, expected))
	return w
}

func (w *UserParametersAssert) HasTimeInputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimeInputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTimeOutputFormatLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTimeOutputFormat, expected))
	return w
}

func (w *UserParametersAssert) HasTraceLevelLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTraceLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionAbortOnErrorLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTransactionAbortOnError, expected))
	return w
}

func (w *UserParametersAssert) HasTransactionDefaultIsolationLevelLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTransactionDefaultIsolationLevel, expected))
	return w
}

func (w *UserParametersAssert) HasTwoDigitCenturyStartLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterTwoDigitCenturyStart, expected))
	return w
}

func (w *UserParametersAssert) HasUnsupportedDdlActionLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterUnsupportedDdlAction, expected))
	return w
}

func (w *UserParametersAssert) HasUseCachedResultLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterUseCachedResult, expected))
	return w
}

func (w *UserParametersAssert) HasWeekOfYearPolicyLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterWeekOfYearPolicy, expected))
	return w
}

func (w *UserParametersAssert) HasWeekStartLevel(expected sdk.ParameterType) *UserParametersAssert {
	w.assertions = append(w.assertions, SnowflakeParameterLevelSet(sdk.UserParameterWeekStart, expected))
	return w
}

// default checks

func (w *UserParametersAssert) HasDefaultEnableUnredactedQuerySyntaxErrorValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterEnableUnredactedQuerySyntaxError)
}

func (w *UserParametersAssert) HasDefaultNetworkPolicyValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterNetworkPolicy)
}

func (w *UserParametersAssert) HasDefaultPreventUnloadToInternalStagesValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterPreventUnloadToInternalStages)
}

func (w *UserParametersAssert) HasDefaultAbortDetachedQueryValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterAbortDetachedQuery)
}

func (w *UserParametersAssert) HasDefaultAutocommitValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterAutocommit)
}

func (w *UserParametersAssert) HasDefaultBinaryInputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterBinaryInputFormat)
}

func (w *UserParametersAssert) HasDefaultBinaryOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterBinaryOutputFormat)
}

func (w *UserParametersAssert) HasDefaultClientMemoryLimitValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientMemoryLimit)
}

func (w *UserParametersAssert) HasDefaultClientMetadataRequestUseConnectionCtxValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientMetadataRequestUseConnectionCtx)
}

func (w *UserParametersAssert) HasDefaultClientPrefetchThreadsValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientPrefetchThreads)
}

func (w *UserParametersAssert) HasDefaultClientResultChunkSizeValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientResultChunkSize)
}

func (w *UserParametersAssert) HasDefaultClientResultColumnCaseInsensitiveValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientResultColumnCaseInsensitive)
}

func (w *UserParametersAssert) HasDefaultClientSessionKeepAliveValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientSessionKeepAlive)
}

func (w *UserParametersAssert) HasDefaultClientSessionKeepAliveHeartbeatFrequencyValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency)
}

func (w *UserParametersAssert) HasDefaultClientTimestampTypeMappingValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterClientTimestampTypeMapping)
}

func (w *UserParametersAssert) HasDefaultDateInputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterDateInputFormat)
}

func (w *UserParametersAssert) HasDefaultDateOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterDateOutputFormat)
}

func (w *UserParametersAssert) HasDefaultEnableUnloadPhysicalTypeOptimizationValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterEnableUnloadPhysicalTypeOptimization)
}

func (w *UserParametersAssert) HasDefaultErrorOnNondeterministicMergeValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterErrorOnNondeterministicMerge)
}

func (w *UserParametersAssert) HasDefaultErrorOnNondeterministicUpdateValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterErrorOnNondeterministicUpdate)
}

func (w *UserParametersAssert) HasDefaultGeographyOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterGeographyOutputFormat)
}

func (w *UserParametersAssert) HasDefaultGeometryOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterGeometryOutputFormat)
}

func (w *UserParametersAssert) HasDefaultJdbcTreatDecimalAsIntValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterJdbcTreatDecimalAsInt)
}

func (w *UserParametersAssert) HasDefaultJdbcTreatTimestampNtzAsUtcValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterJdbcTreatTimestampNtzAsUtc)
}

func (w *UserParametersAssert) HasDefaultJdbcUseSessionTimezoneValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterJdbcUseSessionTimezone)
}

func (w *UserParametersAssert) HasDefaultJsonIndentValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterJsonIndent)
}

func (w *UserParametersAssert) HasDefaultLockTimeoutValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterLockTimeout)
}

func (w *UserParametersAssert) HasDefaultLogLevelValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterLogLevel)
}

func (w *UserParametersAssert) HasDefaultMultiStatementCountValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterMultiStatementCount)
}

func (w *UserParametersAssert) HasDefaultNoorderSequenceAsDefaultValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterNoorderSequenceAsDefault)
}

func (w *UserParametersAssert) HasDefaultOdbcTreatDecimalAsIntValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterOdbcTreatDecimalAsInt)
}

func (w *UserParametersAssert) HasDefaultQueryTagValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterQueryTag)
}

func (w *UserParametersAssert) HasDefaultQuotedIdentifiersIgnoreCaseValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterQuotedIdentifiersIgnoreCase)
}

func (w *UserParametersAssert) HasDefaultRowsPerResultsetValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterRowsPerResultset)
}

func (w *UserParametersAssert) HasDefaultS3StageVpceDnsNameValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterS3StageVpceDnsName)
}

func (w *UserParametersAssert) HasDefaultSearchPathValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterSearchPath)
}

func (w *UserParametersAssert) HasDefaultSimulatedDataSharingConsumerValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterSimulatedDataSharingConsumer)
}

func (w *UserParametersAssert) HasDefaultStatementQueuedTimeoutInSecondsValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterStatementQueuedTimeoutInSeconds)
}

func (w *UserParametersAssert) HasDefaultStatementTimeoutInSecondsValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterStatementTimeoutInSeconds)
}

func (w *UserParametersAssert) HasDefaultStrictJsonOutputValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterStrictJsonOutput)
}

func (w *UserParametersAssert) HasDefaultTimestampDayIsAlways24hValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampDayIsAlways24h)
}

func (w *UserParametersAssert) HasDefaultTimestampInputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampInputFormat)
}

func (w *UserParametersAssert) HasDefaultTimestampLtzOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampLtzOutputFormat)
}

func (w *UserParametersAssert) HasDefaultTimestampNtzOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampNtzOutputFormat)
}

func (w *UserParametersAssert) HasDefaultTimestampOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampOutputFormat)
}

func (w *UserParametersAssert) HasDefaultTimestampTypeMappingValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampTypeMapping)
}

func (w *UserParametersAssert) HasDefaultTimestampTzOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimestampTzOutputFormat)
}

func (w *UserParametersAssert) HasDefaultTimezoneValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimezone)
}

func (w *UserParametersAssert) HasDefaultTimeInputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimeInputFormat)
}

func (w *UserParametersAssert) HasDefaultTimeOutputFormatValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTimeOutputFormat)
}

func (w *UserParametersAssert) HasDefaultTraceLevelValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTraceLevel)
}

func (w *UserParametersAssert) HasDefaultTransactionAbortOnErrorValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTransactionAbortOnError)
}

func (w *UserParametersAssert) HasDefaultTransactionDefaultIsolationLevelValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTransactionDefaultIsolationLevel)
}

func (w *UserParametersAssert) HasDefaultTwoDigitCenturyStartValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterTwoDigitCenturyStart)
}

func (w *UserParametersAssert) HasDefaultUnsupportedDdlActionValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterUnsupportedDdlAction)
}

func (w *UserParametersAssert) HasDefaultUseCachedResultValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterUseCachedResult)
}

func (w *UserParametersAssert) HasDefaultWeekOfYearPolicyValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterWeekOfYearPolicy)
}

func (w *UserParametersAssert) HasDefaultWeekStartValue() *UserParametersAssert {
	return w.HasDefaultParameterValue(sdk.UserParameterWeekStart)
}

// default checks explicit

func (w *UserParametersAssert) HasDefaultEnableUnredactedQuerySyntaxErrorValueExplicit() *UserParametersAssert {
	return w.HasEnableUnredactedQuerySyntaxError(false)
}

func (w *UserParametersAssert) HasDefaultNetworkPolicyValueExplicit() *UserParametersAssert {
	return w.HasNetworkPolicy("")
}

func (w *UserParametersAssert) HasDefaultPreventUnloadToInternalStagesValueExplicit() *UserParametersAssert {
	return w.HasPreventUnloadToInternalStages(false)
}

func (w *UserParametersAssert) HasDefaultAbortDetachedQueryValueExplicit() *UserParametersAssert {
	return w.HasAbortDetachedQuery(false)
}

func (w *UserParametersAssert) HasDefaultAutocommitValueExplicit() *UserParametersAssert {
	return w.HasAutocommit(true)
}

func (w *UserParametersAssert) HasDefaultBinaryInputFormatValueExplicit() *UserParametersAssert {
	return w.HasBinaryInputFormat(sdk.BinaryInputFormatHex)
}

func (w *UserParametersAssert) HasDefaultBinaryOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasBinaryOutputFormat(sdk.BinaryOutputFormatHex)
}

func (w *UserParametersAssert) HasDefaultClientMemoryLimitValueExplicit() *UserParametersAssert {
	return w.HasClientMemoryLimit(1536)
}

func (w *UserParametersAssert) HasDefaultClientMetadataRequestUseConnectionCtxValueExplicit() *UserParametersAssert {
	return w.HasClientMetadataRequestUseConnectionCtx(false)
}

func (w *UserParametersAssert) HasDefaultClientPrefetchThreadsValueExplicit() *UserParametersAssert {
	return w.HasClientPrefetchThreads(4)
}

func (w *UserParametersAssert) HasDefaultClientResultChunkSizeValueExplicit() *UserParametersAssert {
	return w.HasClientResultChunkSize(160)
}

func (w *UserParametersAssert) HasDefaultClientResultColumnCaseInsensitiveValueExplicit() *UserParametersAssert {
	return w.HasClientResultColumnCaseInsensitive(false)
}

func (w *UserParametersAssert) HasDefaultClientSessionKeepAliveValueExplicit() *UserParametersAssert {
	return w.HasClientSessionKeepAlive(false)
}

func (w *UserParametersAssert) HasDefaultClientSessionKeepAliveHeartbeatFrequencyValueExplicit() *UserParametersAssert {
	return w.HasClientSessionKeepAliveHeartbeatFrequency(3600)
}

func (w *UserParametersAssert) HasDefaultClientTimestampTypeMappingValueExplicit() *UserParametersAssert {
	return w.HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingLtz)
}

func (w *UserParametersAssert) HasDefaultDateInputFormatValueExplicit() *UserParametersAssert {
	return w.HasDateInputFormat("AUTO")
}

func (w *UserParametersAssert) HasDefaultDateOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasDateOutputFormat("YYYY-MM-DD")
}

func (w *UserParametersAssert) HasDefaultEnableUnloadPhysicalTypeOptimizationValueExplicit() *UserParametersAssert {
	return w.HasEnableUnloadPhysicalTypeOptimization(true)
}

func (w *UserParametersAssert) HasDefaultErrorOnNondeterministicMergeValueExplicit() *UserParametersAssert {
	return w.HasErrorOnNondeterministicMerge(true)
}

func (w *UserParametersAssert) HasDefaultErrorOnNondeterministicUpdateValueExplicit() *UserParametersAssert {
	return w.HasErrorOnNondeterministicUpdate(false)
}

func (w *UserParametersAssert) HasDefaultGeographyOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasGeographyOutputFormat(sdk.GeographyOutputFormatGeoJSON)
}

func (w *UserParametersAssert) HasDefaultGeometryOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasGeometryOutputFormat(sdk.GeometryOutputFormatGeoJSON)
}

func (w *UserParametersAssert) HasDefaultJdbcTreatDecimalAsIntValueExplicit() *UserParametersAssert {
	return w.HasJdbcTreatDecimalAsInt(true)
}

func (w *UserParametersAssert) HasDefaultJdbcTreatTimestampNtzAsUtcValueExplicit() *UserParametersAssert {
	return w.HasJdbcTreatTimestampNtzAsUtc(false)
}

func (w *UserParametersAssert) HasDefaultJdbcUseSessionTimezoneValueExplicit() *UserParametersAssert {
	return w.HasJdbcUseSessionTimezone(true)
}

func (w *UserParametersAssert) HasDefaultJsonIndentValueExplicit() *UserParametersAssert {
	return w.HasJsonIndent(2)
}

func (w *UserParametersAssert) HasDefaultLockTimeoutValueExplicit() *UserParametersAssert {
	return w.HasLockTimeout(43200)
}

func (w *UserParametersAssert) HasDefaultLogLevelValueExplicit() *UserParametersAssert {
	return w.HasLogLevel(sdk.LogLevelOff)
}

func (w *UserParametersAssert) HasDefaultMultiStatementCountValueExplicit() *UserParametersAssert {
	return w.HasMultiStatementCount(1)
}

func (w *UserParametersAssert) HasDefaultNoorderSequenceAsDefaultValueExplicit() *UserParametersAssert {
	return w.HasNoorderSequenceAsDefault(true)
}

func (w *UserParametersAssert) HasDefaultOdbcTreatDecimalAsIntValueExplicit() *UserParametersAssert {
	return w.HasOdbcTreatDecimalAsInt(false)
}

func (w *UserParametersAssert) HasDefaultQueryTagValueExplicit() *UserParametersAssert {
	return w.HasQueryTag("")
}

func (w *UserParametersAssert) HasDefaultQuotedIdentifiersIgnoreCaseValueExplicit() *UserParametersAssert {
	return w.HasQuotedIdentifiersIgnoreCase(false)
}

func (w *UserParametersAssert) HasDefaultRowsPerResultsetValueExplicit() *UserParametersAssert {
	return w.HasRowsPerResultset(0)
}

func (w *UserParametersAssert) HasDefaultS3StageVpceDnsNameValueExplicit() *UserParametersAssert {
	return w.HasS3StageVpceDnsName("")
}

func (w *UserParametersAssert) HasDefaultSearchPathValueExplicit() *UserParametersAssert {
	return w.HasSearchPath("$current, $public")
}

func (w *UserParametersAssert) HasDefaultSimulatedDataSharingConsumerValueExplicit() *UserParametersAssert {
	return w.HasSimulatedDataSharingConsumer("")
}

func (w *UserParametersAssert) HasDefaultStatementQueuedTimeoutInSecondsValueExplicit() *UserParametersAssert {
	return w.HasStatementQueuedTimeoutInSeconds(0)
}

func (w *UserParametersAssert) HasDefaultStatementTimeoutInSecondsValueExplicit() *UserParametersAssert {
	return w.HasStatementTimeoutInSeconds(172800)
}

func (w *UserParametersAssert) HasDefaultStrictJsonOutputValueExplicit() *UserParametersAssert {
	return w.HasStrictJsonOutput(false)
}

func (w *UserParametersAssert) HasDefaultTimestampDayIsAlways24hValueExplicit() *UserParametersAssert {
	return w.HasTimestampDayIsAlways24h(false)
}

func (w *UserParametersAssert) HasDefaultTimestampInputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimestampInputFormat("AUTO")
}

func (w *UserParametersAssert) HasDefaultTimestampLtzOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimestampLtzOutputFormat("")
}

func (w *UserParametersAssert) HasDefaultTimestampNtzOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS.FF3")
}

func (w *UserParametersAssert) HasDefaultTimestampOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM")
}

func (w *UserParametersAssert) HasDefaultTimestampTypeMappingValueExplicit() *UserParametersAssert {
	return w.HasTimestampTypeMapping(sdk.TimestampTypeMappingNtz)
}

func (w *UserParametersAssert) HasDefaultTimestampTzOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimestampTzOutputFormat("")
}

func (w *UserParametersAssert) HasDefaultTimezoneValueExplicit() *UserParametersAssert {
	return w.HasTimezone("America/Los_Angeles")
}

func (w *UserParametersAssert) HasDefaultTimeInputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimeInputFormat("AUTO")
}

func (w *UserParametersAssert) HasDefaultTimeOutputFormatValueExplicit() *UserParametersAssert {
	return w.HasTimeOutputFormat("HH24:MI:SS")
}

func (w *UserParametersAssert) HasDefaultTraceLevelValueExplicit() *UserParametersAssert {
	return w.HasTraceLevel(sdk.TraceLevelOff)
}

func (w *UserParametersAssert) HasDefaultTransactionAbortOnErrorValueExplicit() *UserParametersAssert {
	return w.HasTransactionAbortOnError(false)
}

func (w *UserParametersAssert) HasDefaultTransactionDefaultIsolationLevelValueExplicit() *UserParametersAssert {
	return w.HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted)
}

func (w *UserParametersAssert) HasDefaultTwoDigitCenturyStartValueExplicit() *UserParametersAssert {
	return w.HasTwoDigitCenturyStart(1970)
}

// lowercase for ignore in snowflake by default but uppercase for FAIL
func (w *UserParametersAssert) HasDefaultUnsupportedDdlActionValueExplicit() *UserParametersAssert {
	return w.HasUnsupportedDdlAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))
}

func (w *UserParametersAssert) HasDefaultUseCachedResultValueExplicit() *UserParametersAssert {
	return w.HasUseCachedResult(true)
}

func (w *UserParametersAssert) HasDefaultWeekOfYearPolicyValueExplicit() *UserParametersAssert {
	return w.HasWeekOfYearPolicy(0)
}

func (w *UserParametersAssert) HasDefaultWeekStartValueExplicit() *UserParametersAssert {
	return w.HasWeekStart(0)
}
