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

func (sessionParameters *SessionParameters) setParam(parameter SessionParameter, value string) error {
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.AbortDetachedQuery = b
	case SessionParameterAutocommit:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.Autocommit = b
	case SessionParameterBinaryInputFormat:
		sessionParameters.BinaryInputFormat = Pointer(BinaryInputFormat(value))
	case SessionParameterBinaryOutputFormat:
		sessionParameters.BinaryOutputFormat = Pointer(BinaryOutputFormat(value))
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.ClientMetadataRequestUseConnectionCtx = b
	case SessionParameterClientMetadataUseSessionDatabase:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.ClientMetadataUseSessionDatabase = b
	case SessionParameterClientResultColumnCaseInsensitive:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.ClientResultColumnCaseInsensitive = b
	case SessionParameterDateInputFormat:
		sessionParameters.DateInputFormat = &value
	case SessionParameterDateOutputFormat:
		sessionParameters.DateOutputFormat = &value
	case SessionParameterErrorOnNondeterministicMerge:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.ErrorOnNondeterministicMerge = b
	case SessionParameterErrorOnNondeterministicUpdate:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.ErrorOnNondeterministicUpdate = b
	case SessionParameterGeographyOutputFormat:
		sessionParameters.GeographyOutputFormat = Pointer(GeographyOutputFormat(value))
	case SessionParameterJSONIndent:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("JSON_INDENT session parameter is an integer, got %v", value)
		}
		sessionParameters.JSONIndent = Pointer(v)
	case SessionParameterLockTimeout:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("LOCK_TIMEOUT session parameter is an integer, got %v", value)
		}
		sessionParameters.LockTimeout = Pointer(v)
	case SessionParameterMultiStatementCount:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MULTI_STATEMENT_COUNT session parameter is an integer, got %v", value)
		}
		sessionParameters.MultiStatementCount = Pointer(v)

	case SessionParameterQueryTag:
		sessionParameters.QueryTag = &value
	case SessionParameterQuotedIdentifiersIgnoreCase:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.QuotedIdentifiersIgnoreCase = b
	case SessionParameterRowsPerResultset:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("ROWS_PER_RESULTSET session parameter is an integer, got %v", value)
		}
		sessionParameters.RowsPerResultset = Pointer(v)
	case SessionParameterSimulatedDataSharingConsumer:
		sessionParameters.SimulatedDataSharingConsumer = &value
	case SessionParameterStatementTimeoutInSeconds:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("STATEMENT_TIMEOUT_IN_SECONDS session parameter is an integer, got %v", value)
		}
		sessionParameters.StatementTimeoutInSeconds = Pointer(v)
	case SessionParameterStrictJSONOutput:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.StrictJSONOutput = b
	case SessionParameterTimestampDayIsAlways24h:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.TimestampDayIsAlways24h = b
	case SessionParameterTimestampInputFormat:
		sessionParameters.TimestampInputFormat = &value
	case SessionParameterTimestampLTZOutputFormat:
		sessionParameters.TimestampLTZOutputFormat = &value
	case SessionParameterTimestampNTZOutputFormat:
		sessionParameters.TimestampNTZOutputFormat = &value
	case SessionParameterTimestampOutputFormat:
		sessionParameters.TimestampOutputFormat = &value
	case SessionParameterTimestampTypeMapping:
		sessionParameters.TimestampTypeMapping = &value
	case SessionParameterTimestampTZOutputFormat:
		sessionParameters.TimestampTZOutputFormat = &value
	case SessionParameterTimezone:
		sessionParameters.Timezone = &value
	case SessionParameterTimeInputFormat:
		sessionParameters.TimeInputFormat = &value
	case SessionParameterTimeOutputFormat:
		sessionParameters.TimeOutputFormat = &value
	case SessionParameterTransactionDefaultIsolationLevel:
		sessionParameters.TransactionDefaultIsolationLevel = Pointer(TransactionDefaultIsolationLevel(value))
	case SessionParameterTwoDigitCenturyStart:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("TWO_DIGIT_CENTURY_START session parameter is an integer, got %v", value)
		}
		sessionParameters.TwoDigitCenturyStart = Pointer(v)
	case SessionParameterUnsupportedDDLAction:
		sessionParameters.UnsupportedDDLAction = Pointer(UnsupportedDDLAction(value))
	case SessionParameterUseCachedResult:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		sessionParameters.UseCachedResult = b
	case SessionParameterWeekOfYearPolicy:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("WEEK_OF_YEAR_POLICY session parameter is an integer, got %v", value)
		}
		sessionParameters.WeekOfYearPolicy = Pointer(v)
	case SessionParameterWeekStart:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("WEEK_START session parameter is an integer, got %v", value)
		}
		sessionParameters.WeekStart = Pointer(v)
	default:
		return fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
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
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		sessionParametersUnset.AbortDetachedQuery = Bool(true)
	case SessionParameterAutocommit:
		sessionParametersUnset.Autocommit = Bool(true)
	case SessionParameterBinaryInputFormat:
		sessionParametersUnset.BinaryInputFormat = Bool(true)
	case SessionParameterBinaryOutputFormat:
		sessionParametersUnset.BinaryOutputFormat = Bool(true)
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		sessionParametersUnset.ClientMetadataRequestUseConnectionCtx = Bool(true)
	case SessionParameterClientMetadataUseSessionDatabase:
		sessionParametersUnset.ClientMetadataUseSessionDatabase = Bool(true)
	case SessionParameterClientResultColumnCaseInsensitive:
		sessionParametersUnset.ClientResultColumnCaseInsensitive = Bool(true)
	case SessionParameterDateInputFormat:
		sessionParametersUnset.DateInputFormat = Bool(true)
	case SessionParameterDateOutputFormat:
		sessionParametersUnset.DateOutputFormat = Bool(true)
	case SessionParameterErrorOnNondeterministicMerge:
		sessionParametersUnset.ErrorOnNondeterministicMerge = Bool(true)
	case SessionParameterErrorOnNondeterministicUpdate:
		sessionParametersUnset.ErrorOnNondeterministicUpdate = Bool(true)
	case SessionParameterGeographyOutputFormat:
		sessionParametersUnset.GeographyOutputFormat = Bool(true)
	case SessionParameterJSONIndent:
		sessionParametersUnset.JSONIndent = Bool(true)
	case SessionParameterLockTimeout:
		sessionParametersUnset.LockTimeout = Bool(true)
	case SessionParameterMultiStatementCount:
		sessionParametersUnset.MultiStatementCount = Bool(true)
	case SessionParameterQueryTag:
		sessionParametersUnset.QueryTag = Bool(true)
	case SessionParameterQuotedIdentifiersIgnoreCase:
		sessionParametersUnset.QuotedIdentifiersIgnoreCase = Bool(true)
	case SessionParameterRowsPerResultset:
		sessionParametersUnset.RowsPerResultset = Bool(true)
	case SessionParameterSimulatedDataSharingConsumer:
		sessionParametersUnset.SimulatedDataSharingConsumer = Bool(true)
	case SessionParameterStatementTimeoutInSeconds:
		sessionParametersUnset.StatementTimeoutInSeconds = Bool(true)
	case SessionParameterStrictJSONOutput:
		sessionParametersUnset.StrictJSONOutput = Bool(true)
	case SessionParameterTimestampDayIsAlways24h:
		sessionParametersUnset.TimestampDayIsAlways24h = Bool(true)
	case SessionParameterTimestampInputFormat:
		sessionParametersUnset.TimestampInputFormat = Bool(true)
	case SessionParameterTimestampLTZOutputFormat:
		sessionParametersUnset.TimestampLTZOutputFormat = Bool(true)
	case SessionParameterTimestampNTZOutputFormat:
		sessionParametersUnset.TimestampNTZOutputFormat = Bool(true)
	case SessionParameterTimestampOutputFormat:
		sessionParametersUnset.TimestampOutputFormat = Bool(true)
	case SessionParameterTimestampTypeMapping:
		sessionParametersUnset.TimestampTypeMapping = Bool(true)
	case SessionParameterTimestampTZOutputFormat:
		sessionParametersUnset.TimestampTZOutputFormat = Bool(true)
	case SessionParameterTimezone:
		sessionParametersUnset.Timezone = Bool(true)
	case SessionParameterTimeInputFormat:
		sessionParametersUnset.TimeInputFormat = Bool(true)
	case SessionParameterTimeOutputFormat:
		sessionParametersUnset.TimeOutputFormat = Bool(true)
	case SessionParameterTransactionDefaultIsolationLevel:
		sessionParametersUnset.TransactionDefaultIsolationLevel = Bool(true)
	case SessionParameterTwoDigitCenturyStart:
		sessionParametersUnset.TwoDigitCenturyStart = Bool(true)
	case SessionParameterUnsupportedDDLAction:
		sessionParametersUnset.UnsupportedDDLAction = Bool(true)
	case SessionParameterUseCachedResult:
		sessionParametersUnset.UseCachedResult = Bool(true)
	case SessionParameterWeekOfYearPolicy:
		sessionParametersUnset.WeekOfYearPolicy = Bool(true)
	case SessionParameterWeekStart:
		sessionParametersUnset.WeekStart = Bool(true)
	default:
		return fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
	return nil
}
