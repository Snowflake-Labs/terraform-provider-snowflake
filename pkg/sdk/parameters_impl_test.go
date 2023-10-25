package sdk

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSessionParameters_setParam(t *testing.T) {
	tests := []struct {
		parameter     SessionParameter
		value         string
		expectedValue any
		accessor      func(*SessionParameters) any
	}{
		{parameter: SessionParameterAbortDetachedQuery, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.AbortDetachedQuery }},
		{parameter: SessionParameterAutocommit, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.Autocommit }},
		{parameter: SessionParameterBinaryInputFormat, value: "some", expectedValue: BinaryInputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.BinaryInputFormat }},
		{parameter: SessionParameterBinaryOutputFormat, value: "some", expectedValue: BinaryOutputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.BinaryOutputFormat }},
		{parameter: SessionParameterClientMetadataRequestUseConnectionCtx, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientMetadataRequestUseConnectionCtx }},
		{parameter: SessionParameterClientMetadataUseSessionDatabase, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientMetadataUseSessionDatabase }},
		{parameter: SessionParameterClientResultColumnCaseInsensitive, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientResultColumnCaseInsensitive }},
		{parameter: SessionParameterDateInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.DateInputFormat }},
		{parameter: SessionParameterGeographyOutputFormat, value: "some", expectedValue: GeographyOutputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.GeographyOutputFormat }},
		{parameter: SessionParameterDateOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.DateOutputFormat }},
		{parameter: SessionParameterErrorOnNondeterministicMerge, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ErrorOnNondeterministicMerge }},
		{parameter: SessionParameterErrorOnNondeterministicUpdate, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ErrorOnNondeterministicUpdate }},
		{parameter: SessionParameterJSONIndent, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.JSONIndent }},
		{parameter: SessionParameterLockTimeout, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.LockTimeout }},
		{parameter: SessionParameterMultiStatementCount, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.MultiStatementCount }},
		{parameter: SessionParameterQueryTag, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.QueryTag }},
		{parameter: SessionParameterQuotedIdentifiersIgnoreCase, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.QuotedIdentifiersIgnoreCase }},
		{parameter: SessionParameterRowsPerResultset, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.RowsPerResultset }},
		{parameter: SessionParameterSimulatedDataSharingConsumer, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.SimulatedDataSharingConsumer }},
		{parameter: SessionParameterStatementTimeoutInSeconds, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.StatementTimeoutInSeconds }},
		{parameter: SessionParameterStrictJSONOutput, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.StrictJSONOutput }},
		{parameter: SessionParameterTimeInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimeInputFormat }},
		{parameter: SessionParameterTimeOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimeOutputFormat }},
		{parameter: SessionParameterTimestampDayIsAlways24h, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.TimestampDayIsAlways24h }},
		{parameter: SessionParameterTimestampInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampInputFormat }},
		{parameter: SessionParameterTimestampLTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampLTZOutputFormat }},
		{parameter: SessionParameterTimestampNTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampNTZOutputFormat }},
		{parameter: SessionParameterTimestampOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampOutputFormat }},
		{parameter: SessionParameterTimestampTypeMapping, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampTypeMapping }},
		{parameter: SessionParameterTimestampTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampTZOutputFormat }},
		{parameter: SessionParameterTimezone, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.Timezone }},
		{parameter: SessionParameterTransactionDefaultIsolationLevel, value: "some", expectedValue: TransactionDefaultIsolationLevel("some"), accessor: func(sp *SessionParameters) any { return *sp.TransactionDefaultIsolationLevel }},
		{parameter: SessionParameterTwoDigitCenturyStart, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.TwoDigitCenturyStart }},
		{parameter: SessionParameterUnsupportedDDLAction, value: "some", expectedValue: UnsupportedDDLAction("some"), accessor: func(sp *SessionParameters) any { return *sp.UnsupportedDDLAction }},
		{parameter: SessionParameterUseCachedResult, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.UseCachedResult }},
		{parameter: SessionParameterWeekOfYearPolicy, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.WeekOfYearPolicy }},
		{parameter: SessionParameterWeekStart, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.WeekStart }},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test valid value '%s' for parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			sessionParameters := &SessionParameters{}

			err := sessionParameters.setParam(tt.parameter, tt.value)

			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, tt.accessor(sessionParameters))
		})
	}

	invalidCases := []struct {
		parameter SessionParameter
		value     string
	}{
		{parameter: SessionParameterAbortDetachedQuery, value: "true123"},
		{parameter: SessionParameterAutocommit, value: "true123"},
		//{parameter: SessionParameterBinaryInputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterBinaryOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterClientMetadataRequestUseConnectionCtx, value: "true123"},
		{parameter: SessionParameterClientMetadataUseSessionDatabase, value: "true123"},
		{parameter: SessionParameterClientResultColumnCaseInsensitive, value: "true123"},
		//{parameter: SessionParameterDateInputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterGeographyOutputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterDateOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterErrorOnNondeterministicMerge, value: "true123"},
		{parameter: SessionParameterErrorOnNondeterministicUpdate, value: "true123"},
		{parameter: SessionParameterJSONIndent, value: "aaa"},
		{parameter: SessionParameterLockTimeout, value: "aaa"},
		{parameter: SessionParameterMultiStatementCount, value: "aaa"},
		//{parameter: SessionParameterQueryTag, value: "some"}, // add validation
		{parameter: SessionParameterQuotedIdentifiersIgnoreCase, value: "true123"},
		{parameter: SessionParameterRowsPerResultset, value: "aaa"},
		//{parameter: SessionParameterSimulatedDataSharingConsumer, value: "some"}, // add validation
		{parameter: SessionParameterStatementTimeoutInSeconds, value: "aaa"},
		{parameter: SessionParameterStrictJSONOutput, value: "true123"},
		//{parameter: SessionParameterTimeInputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimeOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterTimestampDayIsAlways24h, value: "true123"},
		//{parameter: SessionParameterTimestampInputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimestampLTZOutputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimestampNTZOutputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimestampOutputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimestampTypeMapping, value: "some"}, // add validation
		//{parameter: SessionParameterTimestampTZOutputFormat, value: "some"}, // add validation
		//{parameter: SessionParameterTimezone, value: "some"}, // add validation
		//{parameter: SessionParameterTransactionDefaultIsolationLevel, value: "some"}, // add validation
		{parameter: SessionParameterTwoDigitCenturyStart, value: "aaa"},
		//{parameter: SessionParameterUnsupportedDDLAction, value: "some"}, // add validation
		{parameter: SessionParameterUseCachedResult, value: "true123"},
		{parameter: SessionParameterWeekOfYearPolicy, value: "aaa"},
		{parameter: SessionParameterWeekStart, value: "aaa"},
	}
	for _, tt := range invalidCases {
		t.Run(fmt.Sprintf("test invalid value '%s' for parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			sessionParameters := &SessionParameters{}

			err := sessionParameters.setParam(tt.parameter, tt.value)

			require.Error(t, err)
		})
	}
}
