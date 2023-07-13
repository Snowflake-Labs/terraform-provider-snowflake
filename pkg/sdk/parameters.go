package sdk

import "fmt"

type AccountParameter string

// There is a hierarchical relationship between the different parameter types. Account parameters can set any of account, user, session or object parameters
// https://docs.snowflake.com/en/sql-reference/parameters#parameter-hierarchy-and-types
// Account Parameters include Session Parameters, Object Parameters and User Parameters
const (
	// Account Parameters
	AccountParameterAllowClientMFACaching                        AccountParameter = "ALLOW_CLIENT_MFA_CACHING"
	AccountParameterAllowIDToken                                 AccountParameter = "ALLOW_ID_TOKEN" // #nosec G101
	AccountParameterClientEncryptionKeySize                      AccountParameter = "CLIENT_ENCRYPTION_KEY_SIZE"
	AccountParameterEnableInternalStagesPrivatelink              AccountParameter = "ENABLE_INTERNAL_STAGES_PRIVATELINK"
	AccountParameterEventTable                                   AccountParameter = "EVENT_TABLE"
	AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList AccountParameter = "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"
	AccountParameterInitialReplicationSizeLimitInTB              AccountParameter = "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"
	AccountParameterMinDataRetentionTimeInDays                   AccountParameter = "MIN_DATA_RETENTION_TIME_IN_DAYS"
	AccountParameterNetworkPolicy                                AccountParameter = "NETWORK_POLICY"
	AccountParameterPeriodicDataRekeying                         AccountParameter = "PERIODIC_DATA_REKEYING"
	AccountParameterPreventUnloadToInlineURL                     AccountParameter = "PREVENT_UNLOAD_TO_INLINE_URL"
	AccountParameterPreventUnloadToInternalStages                AccountParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES"
	AccountParameterRequireStorageIntegrationForStageCreation    AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"
	AccountParameterRequireStorageIntegrationForStageOperation   AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"
	AccountParameterSSOLoginPage                                 AccountParameter = "SSO_LOGIN_PAGE"

	// Session Parameters (inherited)
	AccountParameterAbortDetachedQuery               AccountParameter = "ABORT_DETACHED_QUERY"
	AccountParameterAutocommit                       AccountParameter = "AUTOCOMMIT"
	AccountParameterBinaryInputFormat                AccountParameter = "BINARY_INPUT_FORMAT"
	AccountParameterBinaryOutputFormat               AccountParameter = "BINARY_OUTPUT_FORMAT"
	AccountParameterDateInputFormat                  AccountParameter = "DATE_INPUT_FORMAT"
	AccountParameterDateOutputFormat                 AccountParameter = "DATE_OUTPUT_FORMAT"
	AccountParameterErrorOnNondeterministicMerge     AccountParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	AccountParameterErrorOnNondeterministicUpdate    AccountParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	AccountParameterGeographyOutputFormat            AccountParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	AccountParameterJSONIndent                       AccountParameter = "JSON_INDENT"
	AccountParameterLockTimeout                      AccountParameter = "LOCK_TIMEOUT"
	AccountParameterQueryTag                         AccountParameter = "QUERY_TAG"
	AccountParameterRowsPerResultset                 AccountParameter = "ROWS_PER_RESULTSET"
	AccountParameterSimulatedDataSharingConsumer     AccountParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	AccountParameterStatementTimeoutInSeconds        AccountParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	AccountParameterStrictJSONOutput                 AccountParameter = "STRICT_JSON_OUTPUT"
	AccountParameterTimestampDayIsAlways24h          AccountParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	AccountParameterTimestampInputFormat             AccountParameter = "TIMESTAMP_INPUT_FORMAT"
	AccountParameterTimestampLTZOutputFormat         AccountParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	AccountParameterTimestampNTZOutputFormat         AccountParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	AccountParameterTimestampOutputFormat            AccountParameter = "TIMESTAMP_OUTPUT_FORMAT"
	AccountParameterTimestampTypeMapping             AccountParameter = "TIMESTAMP_TYPE_MAPPING"
	AccountParameterTimestampTZOutputFormat          AccountParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	AccountParameterTimezone                         AccountParameter = "TIMEZONE"
	AccountParameterTimeInputFormat                  AccountParameter = "TIME_INPUT_FORMAT"
	AccountParameterTimeOutputFormat                 AccountParameter = "TIME_OUTPUT_FORMAT"
	AccountParameterTransactionDefaultIsolationLevel AccountParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	AccountParameterTwoDigitCenturyStart             AccountParameter = "TWO_DIGIT_CENTURY_START"
	AccountParameterUnsupportedDDLAction             AccountParameter = "UNSUPPORTED_DDL_ACTION"
	AccountParameterUseCachedResult                  AccountParameter = "USE_CACHED_RESULT"
	AccountParameterWeekOfYearPolicy                 AccountParameter = "WEEK_OF_YEAR_POLICY"
	AccountParameterWeekStart                        AccountParameter = "WEEK_START"

	// Object Parameters (inherited)
	AccountParameterDataRetentionTimeInDays             AccountParameter = "DATA_RETENTION_TIME_IN_DAYS"
	AccountParameterDefaultDDLCollation                 AccountParameter = "DEFAULT_DDL_COLLATION"
	AccountParameterLogLevel                            AccountParameter = "LOG_LEVEL"
	AccountParameterMaxConcurrencyLevel                 AccountParameter = "MAX_CONCURRENCY_LEVEL"
	AccountParameterMaxDataExtensionTimeInDays          AccountParameter = "MAX_DATA_EXTENSION_TIME_IN_DAYS"
	AccountParameterPipeExecutionPaused                 AccountParameter = "PIPE_EXECUTION_PAUSED"
	AccountParameterStatementQueuedTimeoutInSeconds     AccountParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	AccountParameterShareRestrictions                   AccountParameter = "SHARE_RESTRICTIONS"
	AccountParameterSuspendTaskAfterNumFailures         AccountParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	AccountParameterTraceLevel                          AccountParameter = "TRACE_LEVEL"
	AccountParameterUserTaskManagedInitialWarehouseSize AccountParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	AccountParameterUserTaskTimeoutMs                   AccountParameter = "USER_TASK_TIMEOUT_MS"

	// User Parameters (inherited)
	AccountParameterEnableUnredactedQuerySyntaxError AccountParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
)

type SessionParameter string

const (
	SessionParameterAbortDetachedQuery               SessionParameter = "ABORT_DETACHED_QUERY"
	SessionParameterAutocommit                       SessionParameter = "AUTOCOMMIT"
	SessionParameterBinaryInputFormat                SessionParameter = "BINARY_INPUT_FORMAT"
	SessionParameterBinaryOutputFormat               SessionParameter = "BINARY_OUTPUT_FORMAT"
	SessionParameterDateInputFormat                  SessionParameter = "DATE_INPUT_FORMAT"
	SessionParameterDateOutputFormat                 SessionParameter = "DATE_OUTPUT_FORMAT"
	SessionParameterErrorOnNondeterministicMerge     SessionParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	SessionParameterErrorOnNondeterministicUpdate    SessionParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	SessionParameterGeographyOutputFormat            SessionParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	SessionParameterJSONIndent                       SessionParameter = "JSON_INDENT"
	SessionParameterLockTimeout                      SessionParameter = "LOCK_TIMEOUT"
	SessionParameterQueryTag                         SessionParameter = "QUERY_TAG"
	SessionParameterRowsPerResultset                 SessionParameter = "ROWS_PER_RESULTSET"
	SessionParameterSimulatedDataSharingConsumer     SessionParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	SessionParameterStatementTimeoutInSeconds        SessionParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	SessionParameterStrictJSONOutput                 SessionParameter = "STRICT_JSON_OUTPUT"
	SessionParameterTimestampDayIsAlways24h          SessionParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	SessionParameterTimestampInputFormat             SessionParameter = "TIMESTAMP_INPUT_FORMAT"
	SessionParameterTimestampLTZOutputFormat         SessionParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	SessionParameterTimestampNTZOutputFormat         SessionParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	SessionParameterTimestampOutputFormat            SessionParameter = "TIMESTAMP_OUTPUT_FORMAT"
	SessionParameterTimestampTypeMapping             SessionParameter = "TIMESTAMP_TYPE_MAPPING"
	SessionParameterTimestampTZOutputFormat          SessionParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	SessionParameterTimezone                         SessionParameter = "TIMEZONE"
	SessionParameterTimeInputFormat                  SessionParameter = "TIME_INPUT_FORMAT"
	SessionParameterTimeOutputFormat                 SessionParameter = "TIME_OUTPUT_FORMAT"
	SessionParameterTransactionDefaultIsolationLevel SessionParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	SessionParameterTwoDigitCenturyStart             SessionParameter = "TWO_DIGIT_CENTURY_START"
	SessionParameterUnsupportedDDLAction             SessionParameter = "UNSUPPORTED_DDL_ACTION"
	SessionParameterUseCachedResult                  SessionParameter = "USE_CACHED_RESULT"
	SessionParameterWeekOfYearPolicy                 SessionParameter = "WEEK_OF_YEAR_POLICY"
	SessionParameterWeekStart                        SessionParameter = "WEEK_START"
)

type ObjectParameter string

const (
	// Object Parameters
	ObjectParameterDataRetentionTimeInDays             ObjectParameter = "DATA_RETENTION_TIME_IN_DAYS"
	ObjectParameterDefaultDDLCollation                 ObjectParameter = "DEFAULT_DDL_COLLATION"
	ObjectParameterLogLevel                            ObjectParameter = "LOG_LEVEL"
	ObjectParameterMaxConcurrencyLevel                 ObjectParameter = "MAX_CONCURRENCY_LEVEL"
	ObjectParameterMaxDataExtensionTimeInDays          ObjectParameter = "MAX_DATA_EXTENSION_TIME_IN_DAYS"
	ObjectParameterPipeExecutionPaused                 ObjectParameter = "PIPE_EXECUTION_PAUSED"
	ObjectParameterPreventUnloadToInternalStages       ObjectParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES" // also an account param
	ObjectParameterStatementQueuedTimeoutInSeconds     ObjectParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	ObjectParameterNetworkPolicy                       ObjectParameter = "NETWORK_POLICY" // also an account param
	ObjectParameterShareRestrictions                   ObjectParameter = "SHARE_RESTRICTIONS"
	ObjectParameterSuspendTaskAfterNumFailures         ObjectParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	ObjectParameterTraceLevel                          ObjectParameter = "TRACE_LEVEL"
	ObjectParameterUserTaskManagedInitialWarehouseSize ObjectParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	ObjectParameterUserTaskTimeoutMs                   ObjectParameter = "USER_TASK_TIMEOUT_MS"

	// User Parameters
	ObjectParameterEnableUnredactedQuerySyntaxError ObjectParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
)

type UserParameter string

const (
	// User Parameters
	UserParameterEnableUnredactedQuerySyntaxError UserParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"

	// Session Parameters (inherited)
	UserParameterAbortDetachedQuery               UserParameter = "ABORT_DETACHED_QUERY"
	UserParameterAutocommit                       UserParameter = "AUTOCOMMIT"
	UserParameterBinaryInputFormat                UserParameter = "BINARY_INPUT_FORMAT"
	UserParameterBinaryOutputFormat               UserParameter = "BINARY_OUTPUT_FORMAT"
	UserParameterDateInputFormat                  UserParameter = "DATE_INPUT_FORMAT"
	UserParameterDateOutputFormat                 UserParameter = "DATE_OUTPUT_FORMAT"
	UserParameterErrorOnNondeterministicMerge     UserParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	UserParameterErrorOnNondeterministicUpdate    UserParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	UserParameterGeographyOutputFormat            UserParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	UserParameterJSONIndent                       UserParameter = "JSON_INDENT"
	UserParameterLockTimeout                      UserParameter = "LOCK_TIMEOUT"
	UserParameterQueryTag                         UserParameter = "QUERY_TAG"
	UserParameterRowsPerResultset                 UserParameter = "ROWS_PER_RESULTSET"
	UserParameterSimulatedDataSharingConsumer     UserParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	UserParameterStatementTimeoutInSeconds        UserParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	UserParameterStrictJSONOutput                 UserParameter = "STRICT_JSON_OUTPUT"
	UserParameterTimestampDayIsAlways24h          UserParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	UserParameterTimestampInputFormat             UserParameter = "TIMESTAMP_INPUT_FORMAT"
	UserParameterTimestampLTZOutputFormat         UserParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	UserParameterTimestampNTZOutputFormat         UserParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	UserParameterTimestampOutputFormat            UserParameter = "TIMESTAMP_OUTPUT_FORMAT"
	UserParameterTimestampTypeMapping             UserParameter = "TIMESTAMP_TYPE_MAPPING"
	UserParameterTimestampTZOutputFormat          UserParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	UserParameterTimezone                         UserParameter = "TIMEZONE"
	UserParameterTimeInputFormat                  UserParameter = "TIME_INPUT_FORMAT"
	UserParameterTimeOutputFormat                 UserParameter = "TIME_OUTPUT_FORMAT"
	UserParameterTransactionDefaultIsolationLevel UserParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	UserParameterTwoDigitCenturyStart             UserParameter = "TWO_DIGIT_CENTURY_START"
	UserParameterUnsupportedDDLAction             UserParameter = "UNSUPPORTED_DDL_ACTION"
	UserParameterUseCachedResult                  UserParameter = "USE_CACHED_RESULT"
	UserParameterWeekOfYearPolicy                 UserParameter = "WEEK_OF_YEAR_POLICY"
	UserParameterWeekStart                        UserParameter = "WEEK_START"
)

type AccountParameters struct {
	// Account Parameters
	AllowClientMFACaching                        *bool    `ddl:"parameter" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                 *bool    `ddl:"parameter" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                      *int     `ddl:"parameter" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableInternalStagesPrivatelink              *bool    `ddl:"parameter" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EventTable                                   *string  `ddl:"parameter,single_quotes" sql:"EVENT_TABLE"`
	ExternalOAuthAddPrivilegedRolesToBlockedList *bool    `ddl:"parameter" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTB              *float64 `ddl:"parameter" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                   *int     `ddl:"parameter" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                                *string  `ddl:"parameter,single_quotes" sql:"NETWORK_POLICY"`
	PeriodicDataRekeying                         *bool    `ddl:"parameter" sql:"PERIODIC_DATA_REKEYING"`
	PreventUnloadToInlineURL                     *bool    `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                *bool    `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation    *bool    `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation   *bool    `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                                 *bool    `ddl:"parameter" sql:"SSO_LOGIN_PAGE"`
}

func (v *AccountParameters) validate() error {
	if valueSet(v.ClientEncryptionKeySize) {
		if !(*v.ClientEncryptionKeySize == 128 || *v.ClientEncryptionKeySize == 256) {
			return fmt.Errorf("CLIENT_ENCRYPTION_KEY_SIZE must be either 128 or 256")
		}
	}
	if valueSet(v.InitialReplicationSizeLimitInTB) {
		l := *v.InitialReplicationSizeLimitInTB
		if l < 0.0 || (l < 0.0 && l < 1.0) {
			return fmt.Errorf("%v must be 0.0 and above with a scale of at least 1 (e.g. 20.5, 32.25, 33.333, etc.)", l)
		}
		return nil
	}
	if valueSet(v.MinDataRetentionTimeInDays) {
		if ok := validateIntInRange(*v.MinDataRetentionTimeInDays, 0, 90); !ok {
			return fmt.Errorf("MIN_DATA_RETENTION_TIME_IN_DAYS must be between 0 and 90")
		}
	}
	return nil
}

type AccountParametersUnset struct {
	AllowClientMFACaching                        *bool `ddl:"keyword" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                 *bool `ddl:"keyword" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                      *bool `ddl:"keyword" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableInternalStagesPrivatelink              *bool `ddl:"keyword" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EventTable                                   *bool `ddl:"keyword" sql:"EVENT_TABLE"`
	ExternalOAuthAddPrivilegedRolesToBlockedList *bool `ddl:"keyword" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTB              *bool `ddl:"keyword" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                   *bool `ddl:"keyword" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                                *bool `ddl:"keyword,single_quotes" sql:"NETWORK_POLICY"`
	PeriodicDataRekeying                         *bool `ddl:"keyword" sql:"PERIODIC_DATA_REKEYING"`
	PreventUnloadToInlineURL                     *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation    *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation   *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                                 *bool `ddl:"keyword" sql:"SSO_LOGIN_PAGE"`
}

type GeographyOutputFormat string

const (
	GeographyOutputFormatGeoJSON GeographyOutputFormat = "GeoJSON"
	GeographyOutputFormatWKT     GeographyOutputFormat = "WKT"
	GeographyOutputFormatWKB     GeographyOutputFormat = "WKB"
	GeographyOutputFormatEWKT    GeographyOutputFormat = "EWKT"
)

type BinaryInputFormat string

const (
	BinaryInputFormatHex    BinaryInputFormat = "HEX"
	BinaryInputFormatBase64 BinaryInputFormat = "BASE64"
	BinaryInputFormatUTF8   BinaryInputFormat = "UTF8"
)

type TransactionDefaultIsolationLevel string

const (
	TransactionDefaultIsolationLevelReadCommitted TransactionDefaultIsolationLevel = "READ COMMITTED"
)

type UnsupportedDDLAction string

const (
	UnsupportedDDLActionIgnore UnsupportedDDLAction = "IGNORE"
	UnsupportedDDLActionFail   UnsupportedDDLAction = "FAIL"
)

type SessionParameters struct {
	AbortDetachedQuery               *bool                             `ddl:"parameter" sql:"ABORT_DETACHED_QUERY"`
	Autocommit                       *bool                             `ddl:"parameter" sql:"AUTOCOMMIT"`
	BinaryInputFormat                *string                           `ddl:"parameter,single_quotes" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat               *BinaryInputFormat                `ddl:"parameter,single_quotes" sql:"BINARY_OUTPUT_FORMAT"`
	DateInputFormat                  *string                           `ddl:"parameter,single_quotes" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                 *string                           `ddl:"parameter,single_quotes" sql:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge     *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate    *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat            *GeographyOutputFormat            `ddl:"parameter,single_quotes" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	JSONIndent                       *int                              `ddl:"parameter" sql:"JSON_INDENT"`
	LockTimeout                      *int                              `ddl:"parameter" sql:"LOCK_TIMEOUT"`
	QueryTag                         *string                           `ddl:"parameter,single_quotes" sql:"QUERY_TAG"`
	RowsPerResultset                 *int                              `ddl:"parameter" sql:"ROWS_PER_RESULTSET"`
	SimulatedDataSharingConsumer     *string                           `ddl:"parameter,single_quotes" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementTimeoutInSeconds        *int                              `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJSONOutput                 *bool                             `ddl:"parameter" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h          *bool                             `ddl:"parameter" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat             *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat         *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat         *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat            *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping             *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat          *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                         *string                           `ddl:"parameter,single_quotes" sql:"TIMEZONE"`
	TimeInputFormat                  *string                           `ddl:"parameter,single_quotes" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                 *string                           `ddl:"parameter,single_quotes" sql:"TIME_OUTPUT_FORMAT"`
	TransactionDefaultIsolationLevel *TransactionDefaultIsolationLevel `ddl:"parameter,single_quotes" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart             *int                              `ddl:"parameter" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction             *UnsupportedDDLAction             `ddl:"parameter,single_quotes" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                  *bool                             `ddl:"parameter" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                 *int                              `ddl:"parameter" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                        *int                              `ddl:"parameter" sql:"WEEK_START"`
}

func (v *SessionParameters) validate() error {
	if valueSet(v.JSONIndent) {
		if ok := validateIntInRange(*v.JSONIndent, 0, 16); !ok {
			return fmt.Errorf("JSON_INDENT must be between 0 and 16")
		}
	}
	if valueSet(v.LockTimeout) {
		if ok := validateIntGreaterThanOrEqual(*v.LockTimeout, 0); !ok {
			return fmt.Errorf("LOCK_TIMEOUT must be greater than or equal to 0")
		}
	}
	if valueSet(v.QueryTag) {
		if len(*v.QueryTag) > 2000 {
			return fmt.Errorf("QUERY_TAG must be less than 2000 characters")
		}
	}
	if valueSet(v.RowsPerResultset) {
		if ok := validateIntGreaterThanOrEqual(*v.RowsPerResultset, 0); !ok {
			return fmt.Errorf("ROWS_PER_RESULTSET must be greater than or equal to 0")
		}
	}
	if valueSet(v.TwoDigitCenturyStart) {
		if ok := validateIntInRange(*v.TwoDigitCenturyStart, 1900, 2100); !ok {
			return fmt.Errorf("TWO_DIGIT_CENTURY_START must be between 1900 and 2100")
		}
	}
	if valueSet(v.WeekOfYearPolicy) {
		if ok := validateIntInRange(*v.WeekOfYearPolicy, 0, 1); !ok {
			return fmt.Errorf("WEEK_OF_YEAR_POLICY must be either 0 or 1")
		}
	}
	if valueSet(v.WeekStart) {
		if ok := validateIntInRange(*v.WeekStart, 0, 1); !ok {
			return fmt.Errorf("WEEK_START must be either 0 or 1")
		}
	}
	return nil
}

type SessionParametersUnset struct {
	AbortDetachedQuery               *bool `ddl:"keyword" sql:"ABORT_DETACHED_QUERY"`
	Autocommit                       *bool `ddl:"keyword" sql:"AUTOCOMMIT"`
	BinaryInputFormat                *bool `ddl:"keyword" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat               *bool `ddl:"keyword" sql:"BINARY_OUTPUT_FORMAT"`
	DateInputFormat                  *bool `ddl:"keyword" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                 *bool `ddl:"keyword" sql:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge     *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate    *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat            *bool `ddl:"keyword" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	JSONIndent                       *bool `ddl:"keyword" sql:"JSON_INDENT"`
	LockTimeout                      *bool `ddl:"keyword" sql:"LOCK_TIMEOUT"`
	QueryTag                         *bool `ddl:"keyword" sql:"QUERY_TAG"`
	RowsPerResultset                 *bool `ddl:"keyword" sql:"ROWS_PER_RESULTSET"`
	SimulatedDataSharingConsumer     *bool `ddl:"keyword" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementTimeoutInSeconds        *bool `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJSONOutput                 *bool `ddl:"keyword" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h          *bool `ddl:"keyword" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat             *bool `ddl:"keyword" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat         *bool `ddl:"keyword" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat         *bool `ddl:"keyword" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat            *bool `ddl:"keyword" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping             *bool `ddl:"keyword" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat          *bool `ddl:"keyword" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                         *bool `ddl:"keyword" sql:"TIMEZONE"`
	TimeInputFormat                  *bool `ddl:"keyword" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                 *bool `ddl:"keyword" sql:"TIME_OUTPUT_FORMAT"`
	TransactionDefaultIsolationLevel *bool `ddl:"keyword" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart             *bool `ddl:"keyword" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction             *bool `ddl:"keyword" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                  *bool `ddl:"keyword" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                 *bool `ddl:"keyword" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                        *bool `ddl:"keyword" sql:"WEEK_START"`
}

func (v *SessionParametersUnset) validate() error {
	if ok := anyValueSet(v.AbortDetachedQuery, v.Autocommit, v.BinaryInputFormat, v.BinaryOutputFormat, v.DateInputFormat, v.DateOutputFormat, v.ErrorOnNondeterministicMerge, v.ErrorOnNondeterministicUpdate, v.GeographyOutputFormat, v.JSONIndent, v.LockTimeout, v.QueryTag, v.RowsPerResultset, v.SimulatedDataSharingConsumer, v.StatementTimeoutInSeconds, v.StrictJSONOutput, v.TimestampDayIsAlways24h, v.TimestampInputFormat, v.TimestampLTZOutputFormat, v.TimestampNTZOutputFormat, v.TimestampOutputFormat, v.TimestampTypeMapping, v.TimestampTZOutputFormat, v.Timezone, v.TimeInputFormat, v.TimeOutputFormat, v.TransactionDefaultIsolationLevel, v.TwoDigitCenturyStart, v.UnsupportedDDLAction, v.UseCachedResult, v.WeekOfYearPolicy, v.WeekStart); !ok {
		return fmt.Errorf("at least one session parameter must be set")
	}
	return nil
}

type LogLevel string

const (
	LogLevelTrace LogLevel = "TRACE"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"
	LogLevelOff   LogLevel = "OFF"
)

type TraceLevel string

const (
	TraceLevelAlways  TraceLevel = "ALWAYS"
	TraceLevelOnEvent TraceLevel = "ON_EVENT"
	TraceLevelOff     TraceLevel = "OFF"
)

type ObjectParameters struct {
	DataRetentionTimeInDays             *int           `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	DefaultDDLCollation                 *string        `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	LogLevel                            *LogLevel      `ddl:"parameter" sql:"LOG_LEVEL"`
	MaxConcurrencyLevel                 *int           `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	MaxDataExtensionTimeInDays          *int           `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	PipeExecutionPaused                 *bool          `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	PreventUnloadToInternalStages       *bool          `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	StatementQueuedTimeoutInSeconds     *int           `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	NetworkPolicy                       *string        `ddl:"parameter,single_quotes" sql:"NETWORK_POLICY"`
	ShareRestrictions                   *bool          `ddl:"parameter" sql:"SHARE_RESTRICTIONS"`
	SuspendTaskAfterNumFailures         *int           `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TraceLevel                          *TraceLevel    `ddl:"parameter" sql:"TRACE_LEVEL"`
	UserTaskManagedInitialWarehouseSize *WarehouseSize `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                   *int           `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
}

func (v *ObjectParameters) validate() error {
	if valueSet(v.DataRetentionTimeInDays) {
		if ok := validateIntInRange(*v.DataRetentionTimeInDays, 0, 90); !ok {
			return fmt.Errorf("DATA_RETENTION_TIME_IN_DAYS must be between 0 and 90")
		}
	}
	if valueSet(v.MaxConcurrencyLevel) {
		if ok := validateIntGreaterThanOrEqual(*v.MaxConcurrencyLevel, 1); !ok {
			return fmt.Errorf("MAX_CONCURRENCY_LEVEL must be greater than or equal to 1")
		}
	}

	if valueSet(v.MaxDataExtensionTimeInDays) {
		if ok := validateIntInRange(*v.MaxDataExtensionTimeInDays, 0, 90); !ok {
			return fmt.Errorf("MAX_DATA_EXTENSION_TIME_IN_DAYS must be between 0 and 90")
		}
	}

	if valueSet(v.StatementQueuedTimeoutInSeconds) {
		if ok := validateIntGreaterThanOrEqual(*v.StatementQueuedTimeoutInSeconds, 0); !ok {
			return fmt.Errorf("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS must be greater than or equal to 0")
		}
	}

	if valueSet(v.SuspendTaskAfterNumFailures) {
		if ok := validateIntGreaterThanOrEqual(*v.SuspendTaskAfterNumFailures, 0); !ok {
			return fmt.Errorf("SUSPEND_TASK_AFTER_NUM_FAILURES must be greater than or equal to 0")
		}
	}

	if valueSet(v.UserTaskTimeoutMs) {
		if ok := validateIntInRange(*v.UserTaskTimeoutMs, 0, 86400000); !ok {
			return fmt.Errorf("USER_TASK_TIMEOUT_MS must be between 0 and 86400000")
		}
	}
	return nil
}

type ObjectParametersUnset struct {
	DataRetentionTimeInDays             *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	DefaultDDLCollation                 *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	LogLevel                            *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	MaxConcurrencyLevel                 *bool `ddl:"keyword" sql:"MAX_CONCURRENCY_LEVEL"`
	MaxDataExtensionTimeInDays          *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	PipeExecutionPaused                 *bool `ddl:"keyword" sql:"PIPE_EXECUTION_PAUSED"`
	PreventUnloadToInternalStages       *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	StatementQueuedTimeoutInSeconds     *bool `ddl:"keyword" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	NetworkPolicy                       *bool `ddl:"keyword,single_quotes" sql:"NETWORK_POLICY"`
	ShareRestrictions                   *bool `ddl:"keyword" sql:"SHARE_RESTRICTIONS"`
	SuspendTaskAfterNumFailures         *bool `ddl:"keyword" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TraceLevel                          *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
	UserTaskManagedInitialWarehouseSize *bool `ddl:"keyword" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                   *bool `ddl:"keyword" sql:"USER_TASK_TIMEOUT_MS"`
}

type UserParameters struct {
	EnableUnredactedQuerySyntaxError *bool `ddl:"parameter" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
}

func (v *UserParameters) validate() error {
	return nil
}

type UserParametersUnset struct {
	EnableUnredactedQuerySyntaxError *bool `ddl:"keyword" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
}
