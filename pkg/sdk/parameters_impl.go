package sdk

import (
	"fmt"
	"strconv"
)

func GetSessionParametersFrom(params map[string]any) (*SessionParameters, error) {
	sessionParameters := &SessionParameters{}
	for k, v := range params {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expecting string value for parameter %s (current value: %v)", k, v)
		}
		err := sessionParameters.setParam(SessionParameter(k), s)
		if err != nil {
			return nil, err
		}
	}
	return sessionParameters, nil
}

// TODO [SNOW-1348330]: get type based on the tag in SessionParameters struct and handle in a generic way
// TODO [SNOW-1348330]: use sdk.ToX for the enums
func (sessionParameters *SessionParameters) setParam(parameter SessionParameter, value string) error {
	var err error
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		err = setBooleanValue(parameter, value, &sessionParameters.AbortDetachedQuery)
	case SessionParameterAutocommit:
		err = setBooleanValue(parameter, value, &sessionParameters.Autocommit)
	case SessionParameterBinaryInputFormat:
		sessionParameters.BinaryInputFormat = Pointer(BinaryInputFormat(value))
	case SessionParameterBinaryOutputFormat:
		sessionParameters.BinaryOutputFormat = Pointer(BinaryOutputFormat(value))
	case SessionParameterClientMemoryLimit:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientMemoryLimit)
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientMetadataRequestUseConnectionCtx)
	case SessionParameterClientPrefetchThreads:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientPrefetchThreads)
	case SessionParameterClientResultChunkSize:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientResultChunkSize)
	case SessionParameterClientResultColumnCaseInsensitive:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientResultColumnCaseInsensitive)
	case SessionParameterClientMetadataUseSessionDatabase:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientMetadataUseSessionDatabase)
	case SessionParameterClientSessionKeepAlive:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientSessionKeepAlive)
	case SessionParameterClientSessionKeepAliveHeartbeatFrequency:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientSessionKeepAliveHeartbeatFrequency)
	case SessionParameterClientTimestampTypeMapping:
		sessionParameters.ClientTimestampTypeMapping = Pointer(ClientTimestampTypeMapping(value))
	case SessionParameterDateInputFormat:
		sessionParameters.DateInputFormat = &value
	case SessionParameterDateOutputFormat:
		sessionParameters.DateOutputFormat = &value
	case SessionParameterEnableUnloadPhysicalTypeOptimization:
		err = setBooleanValue(parameter, value, &sessionParameters.EnableUnloadPhysicalTypeOptimization)
	case SessionParameterErrorOnNondeterministicMerge:
		err = setBooleanValue(parameter, value, &sessionParameters.ErrorOnNondeterministicMerge)
	case SessionParameterErrorOnNondeterministicUpdate:
		err = setBooleanValue(parameter, value, &sessionParameters.ErrorOnNondeterministicUpdate)
	case SessionParameterGeographyOutputFormat:
		sessionParameters.GeographyOutputFormat = Pointer(GeographyOutputFormat(value))
	case SessionParameterGeometryOutputFormat:
		sessionParameters.GeometryOutputFormat = Pointer(GeometryOutputFormat(value))
	case SessionParameterJdbcTreatDecimalAsInt:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcTreatDecimalAsInt)
	case SessionParameterJdbcTreatTimestampNtzAsUtc:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcTreatTimestampNtzAsUtc)
	case SessionParameterJdbcUseSessionTimezone:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcUseSessionTimezone)
	case SessionParameterJSONIndent:
		err = setIntegerValue(parameter, value, &sessionParameters.JSONIndent)
	case SessionParameterLockTimeout:
		err = setIntegerValue(parameter, value, &sessionParameters.LockTimeout)
	case SessionParameterLogLevel:
		sessionParameters.LogLevel = Pointer(LogLevel(value))
	case SessionParameterMultiStatementCount:
		err = setIntegerValue(parameter, value, &sessionParameters.MultiStatementCount)
	case SessionParameterNoorderSequenceAsDefault:
		err = setBooleanValue(parameter, value, &sessionParameters.NoorderSequenceAsDefault)
	case SessionParameterOdbcTreatDecimalAsInt:
		err = setBooleanValue(parameter, value, &sessionParameters.OdbcTreatDecimalAsInt)
	case SessionParameterQueryTag:
		sessionParameters.QueryTag = &value
	case SessionParameterQuotedIdentifiersIgnoreCase:
		err = setBooleanValue(parameter, value, &sessionParameters.QuotedIdentifiersIgnoreCase)
	case SessionParameterRowsPerResultset:
		err = setIntegerValue(parameter, value, &sessionParameters.RowsPerResultset)
	case SessionParameterS3StageVpceDnsName:
		sessionParameters.S3StageVpceDnsName = &value
	case SessionParameterSearchPath:
		sessionParameters.SearchPath = &value
	case SessionParameterSimulatedDataSharingConsumer:
		sessionParameters.SimulatedDataSharingConsumer = &value
	case SessionParameterStatementQueuedTimeoutInSeconds:
		err = setIntegerValue(parameter, value, &sessionParameters.StatementQueuedTimeoutInSeconds)
	case SessionParameterStatementTimeoutInSeconds:
		err = setIntegerValue(parameter, value, &sessionParameters.StatementTimeoutInSeconds)
	case SessionParameterStrictJSONOutput:
		err = setBooleanValue(parameter, value, &sessionParameters.StrictJSONOutput)
	case SessionParameterTimestampDayIsAlways24h:
		err = setBooleanValue(parameter, value, &sessionParameters.TimestampDayIsAlways24h)
	case SessionParameterTimestampInputFormat:
		sessionParameters.TimestampInputFormat = &value
	case SessionParameterTimestampLTZOutputFormat:
		sessionParameters.TimestampLTZOutputFormat = &value
	case SessionParameterTimestampNTZOutputFormat:
		sessionParameters.TimestampNTZOutputFormat = &value
	case SessionParameterTimestampOutputFormat:
		sessionParameters.TimestampOutputFormat = &value
	case SessionParameterTimestampTypeMapping:
		sessionParameters.TimestampTypeMapping = Pointer(TimestampTypeMapping(value))
	case SessionParameterTimestampTZOutputFormat:
		sessionParameters.TimestampTZOutputFormat = &value
	case SessionParameterTimezone:
		sessionParameters.Timezone = &value
	case SessionParameterTimeInputFormat:
		sessionParameters.TimeInputFormat = &value
	case SessionParameterTimeOutputFormat:
		sessionParameters.TimeOutputFormat = &value
	case SessionParameterTraceLevel:
		sessionParameters.TraceLevel = Pointer(TraceLevel(value))
	case SessionParameterTransactionAbortOnError:
		err = setBooleanValue(parameter, value, &sessionParameters.TransactionAbortOnError)
	case SessionParameterTransactionDefaultIsolationLevel:
		sessionParameters.TransactionDefaultIsolationLevel = Pointer(TransactionDefaultIsolationLevel(value))
	case SessionParameterTwoDigitCenturyStart:
		err = setIntegerValue(parameter, value, &sessionParameters.TwoDigitCenturyStart)
	case SessionParameterUnsupportedDDLAction:
		sessionParameters.UnsupportedDDLAction = Pointer(UnsupportedDDLAction(value))
	case SessionParameterUseCachedResult:
		err = setBooleanValue(parameter, value, &sessionParameters.UseCachedResult)
	case SessionParameterWeekOfYearPolicy:
		err = setIntegerValue(parameter, value, &sessionParameters.WeekOfYearPolicy)
	case SessionParameterWeekStart:
		err = setIntegerValue(parameter, value, &sessionParameters.WeekStart)
	default:
		err = fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
	return err
}

func setBooleanValue(parameter SessionParameter, value string, setField **bool) error {
	b, err := parseBooleanParameter(string(parameter), value)
	if err != nil {
		return err
	}
	*setField = b
	return nil
}

func setIntegerValue(parameter SessionParameter, value string, setField **int) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s session parameter is an integer, got %v", parameter, value)
	}
	*setField = Pointer(v)
	return nil
}

func GetSessionParametersUnsetFrom(params map[string]any) (*SessionParametersUnset, error) {
	sessionParametersUnset := &SessionParametersUnset{}
	for k := range params {
		err := sessionParametersUnset.setParam(SessionParameter(k))
		if err != nil {
			return nil, err
		}
	}
	return sessionParametersUnset, nil
}

func (sessionParametersUnset *SessionParametersUnset) setParam(parameter SessionParameter) error {
	var unsetField **bool
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		unsetField = &sessionParametersUnset.AbortDetachedQuery
	case SessionParameterAutocommit:
		unsetField = &sessionParametersUnset.Autocommit
	case SessionParameterBinaryInputFormat:
		unsetField = &sessionParametersUnset.BinaryInputFormat
	case SessionParameterBinaryOutputFormat:
		unsetField = &sessionParametersUnset.BinaryOutputFormat
	case SessionParameterClientMemoryLimit:
		unsetField = &sessionParametersUnset.ClientMemoryLimit
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		unsetField = &sessionParametersUnset.ClientMetadataRequestUseConnectionCtx
	case SessionParameterClientPrefetchThreads:
		unsetField = &sessionParametersUnset.ClientPrefetchThreads
	case SessionParameterClientResultChunkSize:
		unsetField = &sessionParametersUnset.ClientResultChunkSize
	case SessionParameterClientResultColumnCaseInsensitive:
		unsetField = &sessionParametersUnset.ClientResultColumnCaseInsensitive
	case SessionParameterClientMetadataUseSessionDatabase:
		unsetField = &sessionParametersUnset.ClientMetadataUseSessionDatabase
	case SessionParameterClientSessionKeepAlive:
		unsetField = &sessionParametersUnset.ClientSessionKeepAlive
	case SessionParameterClientSessionKeepAliveHeartbeatFrequency:
		unsetField = &sessionParametersUnset.ClientSessionKeepAliveHeartbeatFrequency
	case SessionParameterClientTimestampTypeMapping:
		unsetField = &sessionParametersUnset.ClientTimestampTypeMapping
	case SessionParameterDateInputFormat:
		unsetField = &sessionParametersUnset.DateInputFormat
	case SessionParameterDateOutputFormat:
		unsetField = &sessionParametersUnset.DateOutputFormat
	case SessionParameterEnableUnloadPhysicalTypeOptimization:
		unsetField = &sessionParametersUnset.EnableUnloadPhysicalTypeOptimization
	case SessionParameterErrorOnNondeterministicMerge:
		unsetField = &sessionParametersUnset.ErrorOnNondeterministicMerge
	case SessionParameterErrorOnNondeterministicUpdate:
		unsetField = &sessionParametersUnset.ErrorOnNondeterministicUpdate
	case SessionParameterGeographyOutputFormat:
		unsetField = &sessionParametersUnset.GeographyOutputFormat
	case SessionParameterGeometryOutputFormat:
		unsetField = &sessionParametersUnset.GeometryOutputFormat
	case SessionParameterJdbcTreatDecimalAsInt:
		unsetField = &sessionParametersUnset.JdbcTreatDecimalAsInt
	case SessionParameterJdbcTreatTimestampNtzAsUtc:
		unsetField = &sessionParametersUnset.JdbcTreatTimestampNtzAsUtc
	case SessionParameterJdbcUseSessionTimezone:
		unsetField = &sessionParametersUnset.JdbcUseSessionTimezone
	case SessionParameterJSONIndent:
		unsetField = &sessionParametersUnset.JSONIndent
	case SessionParameterLockTimeout:
		unsetField = &sessionParametersUnset.LockTimeout
	case SessionParameterLogLevel:
		unsetField = &sessionParametersUnset.LogLevel
	case SessionParameterMultiStatementCount:
		unsetField = &sessionParametersUnset.MultiStatementCount
	case SessionParameterNoorderSequenceAsDefault:
		unsetField = &sessionParametersUnset.NoorderSequenceAsDefault
	case SessionParameterOdbcTreatDecimalAsInt:
		unsetField = &sessionParametersUnset.OdbcTreatDecimalAsInt
	case SessionParameterQueryTag:
		unsetField = &sessionParametersUnset.QueryTag
	case SessionParameterQuotedIdentifiersIgnoreCase:
		unsetField = &sessionParametersUnset.QuotedIdentifiersIgnoreCase
	case SessionParameterRowsPerResultset:
		unsetField = &sessionParametersUnset.RowsPerResultset
	case SessionParameterS3StageVpceDnsName:
		unsetField = &sessionParametersUnset.S3StageVpceDnsName
	case SessionParameterSearchPath:
		unsetField = &sessionParametersUnset.SearchPath
	case SessionParameterSimulatedDataSharingConsumer:
		unsetField = &sessionParametersUnset.SimulatedDataSharingConsumer
	case SessionParameterStatementQueuedTimeoutInSeconds:
		unsetField = &sessionParametersUnset.StatementQueuedTimeoutInSeconds
	case SessionParameterStatementTimeoutInSeconds:
		unsetField = &sessionParametersUnset.StatementTimeoutInSeconds
	case SessionParameterStrictJSONOutput:
		unsetField = &sessionParametersUnset.StrictJSONOutput
	case SessionParameterTimestampDayIsAlways24h:
		unsetField = &sessionParametersUnset.TimestampDayIsAlways24h
	case SessionParameterTimestampInputFormat:
		unsetField = &sessionParametersUnset.TimestampInputFormat
	case SessionParameterTimestampLTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampLTZOutputFormat
	case SessionParameterTimestampNTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampNTZOutputFormat
	case SessionParameterTimestampOutputFormat:
		unsetField = &sessionParametersUnset.TimestampOutputFormat
	case SessionParameterTimestampTypeMapping:
		unsetField = &sessionParametersUnset.TimestampTypeMapping
	case SessionParameterTimestampTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampTZOutputFormat
	case SessionParameterTimezone:
		unsetField = &sessionParametersUnset.Timezone
	case SessionParameterTimeInputFormat:
		unsetField = &sessionParametersUnset.TimeInputFormat
	case SessionParameterTimeOutputFormat:
		unsetField = &sessionParametersUnset.TimeOutputFormat
	case SessionParameterTraceLevel:
		unsetField = &sessionParametersUnset.TraceLevel
	case SessionParameterTransactionAbortOnError:
		unsetField = &sessionParametersUnset.TransactionAbortOnError
	case SessionParameterTransactionDefaultIsolationLevel:
		unsetField = &sessionParametersUnset.TransactionDefaultIsolationLevel
	case SessionParameterTwoDigitCenturyStart:
		unsetField = &sessionParametersUnset.TwoDigitCenturyStart
	case SessionParameterUnsupportedDDLAction:
		unsetField = &sessionParametersUnset.UnsupportedDDLAction
	case SessionParameterUseCachedResult:
		unsetField = &sessionParametersUnset.UseCachedResult
	case SessionParameterWeekOfYearPolicy:
		unsetField = &sessionParametersUnset.WeekOfYearPolicy
	case SessionParameterWeekStart:
		unsetField = &sessionParametersUnset.WeekStart
	default:
		return fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
	*unsetField = Bool(true)
	return nil
}
