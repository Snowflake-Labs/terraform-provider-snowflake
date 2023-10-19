package sdk

import (
	"fmt"
	"strconv"
)

func GetSessionParametersFrom(params map[string]string) (*SessionParameters, error) {
	sessionParameters := &SessionParameters{}
	for k, v := range params {
		err := sessionParameters.setParam(SessionParameter(k), v)
		if err != nil {
			return nil, err
		}
	}
	return sessionParameters, nil
}

// TODO [SNOW-884987]: use this method in SetSessionParameterOnAccount and in SetSessionParameterOnUser
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
	}
	return fmt.Errorf("%s session parameter is not supported", string(parameter))
}
