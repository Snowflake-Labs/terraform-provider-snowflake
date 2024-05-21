package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	_ validatable = new(ShowParametersOptions)
	_ validatable = new(AccountParameters)
	_ validatable = new(SessionParameters)
	_ validatable = new(ObjectParameters)
	_ validatable = new(UserParameters)
	_ validatable = new(setParameterOnObject)
)

var _ Parameters = (*parameters)(nil)

type Parameters interface {
	SetAccountParameter(ctx context.Context, parameter AccountParameter, value string) error
	SetSessionParameterOnAccount(ctx context.Context, parameter SessionParameter, value string) error
	SetSessionParameterOnUser(ctx context.Context, userID AccountObjectIdentifier, parameter SessionParameter, value string) error
	SetObjectParameterOnAccount(ctx context.Context, parameter ObjectParameter, value string) error
	SetObjectParameterOnObject(ctx context.Context, object Object, parameter ObjectParameter, value string) error
	ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error)
	ShowAccountParameter(ctx context.Context, parameter AccountParameter) (*Parameter, error)
	ShowSessionParameter(ctx context.Context, parameter SessionParameter) (*Parameter, error)
	ShowUserParameter(ctx context.Context, parameter UserParameter, user AccountObjectIdentifier) (*Parameter, error)
	ShowObjectParameter(ctx context.Context, parameter ObjectParameter, object Object) (*Parameter, error)
}

type parameters struct {
	client *Client
}

func (parameters *parameters) SetAccountParameter(ctx context.Context, parameter AccountParameter, value string) error {
	opts := AlterAccountOptions{
		Set: &AccountSet{
			Parameters: &AccountLevelParameters{
				AccountParameters: &AccountParameters{},
			},
		},
	}
	switch parameter {
	case AccountParameterAllowClientMFACaching:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.AllowClientMFACaching = b
	case AccountParameterAllowIDToken:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.AllowIDToken = b
	case AccountParameterClientEncryptionKeySize:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("CLIENT_ENCRYPTION_KEY_SIZE session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.ClientEncryptionKeySize = Pointer(v)
	case AccountParameterEnableIdentifierFirstLogin:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnableIdentifierFirstLogin = b
	case AccountParameterEnableInternalStagesPrivatelink:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.AllowIDToken = b
	case AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnableTriSecretAndRekeyOptOutForImageRepository = b
	case AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnableTriSecretAndRekeyOptOutForSpcsBlockStorage = b
	case AccountParameterEventTable:
		opts.Set.Parameters.AccountParameters.EventTable = &value
	case AccountParameterEnableUnredactedQuerySyntaxError:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnableUnredactedQuerySyntaxError = b
	case AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.ExternalOAuthAddPrivilegedRolesToBlockedList = b
	case AccountParameterInitialReplicationSizeLimitInTB:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("INITIAL_REPLICATION_SIZE_LIMIT_IN_TB session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.InitialReplicationSizeLimitInTB = Pointer(v)

	case AccountParameterMinDataRetentionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MIN_DATA_RETENTION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.MinDataRetentionTimeInDays = Pointer(v)
	case AccountParameterNetworkPolicy:
		opts.Set.Parameters.AccountParameters.NetworkPolicy = &value
	case AccountParameterPeriodicDataRekeying:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.PeriodicDataRekeying = b
	case AccountParameterPreventLoadFromInlineURL:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.PreventLoadFromInlineURL = b
	case AccountParameterPreventUnloadToInlineURL:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.PreventUnloadToInlineURL = b
	case AccountParameterPreventUnloadToInternalStages:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.PreventUnloadToInternalStages = b
	case AccountParameterRequireStorageIntegrationForStageCreation:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageCreation = b
	case AccountParameterRequireStorageIntegrationForStageOperation:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.RequireStorageIntegrationForStageOperation = b
	case AccountParameterSSOLoginPage:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.SSOLoginPage = b
	default:
		return parameters.SetSessionParameterOnAccount(ctx, SessionParameter(parameter), value)
	}
	if err := parameters.client.Accounts.Alter(ctx, &opts); err != nil {
		return err
	}
	return nil
}

func (parameters *parameters) SetSessionParameterOnAccount(ctx context.Context, parameter SessionParameter, value string) error {
	sp := &SessionParameters{}
	err := sp.setParam(parameter, value)
	if err == nil {
		opts := AlterAccountOptions{Set: &AccountSet{Parameters: &AccountLevelParameters{SessionParameters: sp}}}
		err = parameters.client.Accounts.Alter(ctx, &opts)
		if err != nil {
			return err
		}
		return nil
	} else {
		if strings.Contains(err.Error(), "session parameter is not supported") {
			return parameters.SetObjectParameterOnAccount(ctx, ObjectParameter(parameter), value)
		}
		return err
	}
}

func (parameters *parameters) SetSessionParameterOnUser(ctx context.Context, userId AccountObjectIdentifier, parameter SessionParameter, value string) error {
	sp := &SessionParameters{}
	err := sp.setParam(parameter, value)
	if err != nil {
		return err
	}
	opts := AlterUserOptions{Set: &UserSet{SessionParameters: sp}}
	err = parameters.client.Users.Alter(ctx, userId, &opts)
	if err != nil {
		return err
	}
	return nil
}

func (parameters *parameters) SetObjectParameterOnAccount(ctx context.Context, parameter ObjectParameter, value string) error {
	opts := AlterAccountOptions{Set: &AccountSet{Parameters: &AccountLevelParameters{ObjectParameters: &ObjectParameters{}}}}
	switch parameter {
	case ObjectParameterDataRetentionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("DATA_RETENTION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.DataRetentionTimeInDays = Pointer(v)
	case ObjectParameterDefaultDDLCollation:
		opts.Set.Parameters.ObjectParameters.DefaultDDLCollation = &value
	case ObjectParameterLogLevel:
		opts.Set.Parameters.ObjectParameters.LogLevel = Pointer(LogLevel(value))
	case ObjectParameterMaxConcurrencyLevel:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MAX_CONCURRENCY_LEVEL session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.MaxConcurrencyLevel = Pointer(v)
	case ObjectParameterMaxDataExtensionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MAX_DATA_EXTENSION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.MaxDataExtensionTimeInDays = Pointer(v)
	case ObjectParameterPipeExecutionPaused:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.PipeExecutionPaused = b
	case ObjectParameterPreventUnloadToInternalStages:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.PreventUnloadToInternalStages = b
	case ObjectParameterStatementQueuedTimeoutInSeconds:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.StatementQueuedTimeoutInSeconds = Pointer(v)
	case ObjectParameterNetworkPolicy:
		opts.Set.Parameters.ObjectParameters.NetworkPolicy = &value
	case ObjectParameterShareRestrictions:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.ShareRestrictions = b
	case ObjectParameterSuspendTaskAfterNumFailures:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("SUSPEND_TASK_AFTER_NUM_FAILURES session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.SuspendTaskAfterNumFailures = Pointer(v)
	case ObjectParameterTraceLevel:
		opts.Set.Parameters.ObjectParameters.TraceLevel = Pointer(TraceLevel(value))
	case ObjectParameterUserTaskManagedInitialWarehouseSize:
		opts.Set.Parameters.ObjectParameters.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSize(value))
	case ObjectParameterUserTaskTimeoutMs:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("USER_TASK_TIMEOUT_MS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.UserTaskTimeoutMs = Pointer(v)
	case ObjectParameterEnableUnredactedQuerySyntaxError:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.EnableUnredactedQuerySyntaxError = b
	default:
		return fmt.Errorf("Invalid object parameter: %v", string(parameter))
	}
	err := parameters.client.Accounts.Alter(ctx, &opts)
	if err != nil {
		return err
	}
	return nil
}

type setParameterOnObject struct {
	alter            bool             `ddl:"static" sql:"ALTER"`
	objectType       ObjectType       `ddl:"keyword"`
	objectIdentifier ObjectIdentifier `ddl:"identifier"`
	set              bool             `ddl:"static" sql:"SET"`
	parameterKey     ObjectParameter  `ddl:"keyword"`
	equals           bool             `ddl:"static" sql:"="`
	parameterValue   string           `ddl:"keyword"`
}

// TODO: add validations
func (v *setParameterOnObject) validate() error {
	return nil
}

func (parameters *parameters) SetObjectParameterOnObject(ctx context.Context, object Object, parameter ObjectParameter, value string) error {
	opts := &setParameterOnObject{
		objectType:       object.ObjectType,
		objectIdentifier: object.Name,
		parameterKey:     parameter,
		parameterValue:   value,
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = parameters.client.exec(ctx, sql)
	return err
}

func parseBooleanParameter(parameter, value string) (_ *bool, err error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("Boolean value (\"true\"/\"false\") expected for %v parameter, got %v instead", parameter, value)
	}
	return &b, nil
}

type AccountParameter string

// There is a hierarchical relationship between the different parameter types. Account parameters can set any of account, user, session or object parameters
// https://docs.snowflake.com/en/sql-reference/parameters#parameter-hierarchy-and-types
// Account Parameters include Session Parameters, Object Parameters and User Parameters
const (
	// Account Parameters
	AccountParameterAllowClientMFACaching                            AccountParameter = "ALLOW_CLIENT_MFA_CACHING"
	AccountParameterAllowIDToken                                     AccountParameter = "ALLOW_ID_TOKEN" // #nosec G101
	AccountParameterClientEncryptionKeySize                          AccountParameter = "CLIENT_ENCRYPTION_KEY_SIZE"
	AccountParameterEnableIdentifierFirstLogin                       AccountParameter = "ENABLE_IDENTIFIER_FIRST_LOGIN"
	AccountParameterEnableInternalStagesPrivatelink                  AccountParameter = "ENABLE_INTERNAL_STAGES_PRIVATELINK"
	AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository  AccountParameter = "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"   // #nosec G101
	AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage AccountParameter = "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE" // #nosec G101
	AccountParameterEventTable                                       AccountParameter = "EVENT_TABLE"
	AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList     AccountParameter = "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"
	AccountParameterInitialReplicationSizeLimitInTB                  AccountParameter = "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"
	AccountParameterMinDataRetentionTimeInDays                       AccountParameter = "MIN_DATA_RETENTION_TIME_IN_DAYS"
	AccountParameterNetworkPolicy                                    AccountParameter = "NETWORK_POLICY"
	AccountParameterPeriodicDataRekeying                             AccountParameter = "PERIODIC_DATA_REKEYING"
	AccountParameterPreventLoadFromInlineURL                         AccountParameter = "PREVENT_LOAD_FROM_INLINE_URL"
	AccountParameterPreventUnloadToInlineURL                         AccountParameter = "PREVENT_UNLOAD_TO_INLINE_URL"
	AccountParameterPreventUnloadToInternalStages                    AccountParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES"
	AccountParameterRequireStorageIntegrationForStageCreation        AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"
	AccountParameterRequireStorageIntegrationForStageOperation       AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"
	AccountParameterSSOLoginPage                                     AccountParameter = "SSO_LOGIN_PAGE"

	// Session Parameters (inherited)
	AccountParameterAbortDetachedQuery                    AccountParameter = "ABORT_DETACHED_QUERY"
	AccountParameterAutocommit                            AccountParameter = "AUTOCOMMIT"
	AccountParameterBinaryInputFormat                     AccountParameter = "BINARY_INPUT_FORMAT"
	AccountParameterBinaryOutputFormat                    AccountParameter = "BINARY_OUTPUT_FORMAT"
	AccountParameterClientMetadataRequestUseConnectionCtx AccountParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	AccountParameterClientMetadataUseSessionDatabase      AccountParameter = "CLIENT_METADATA_USE_SESSION_DATABASE"
	AccountParameterClientResultColumnCaseInsensitive     AccountParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	AccountParameterDateInputFormat                       AccountParameter = "DATE_INPUT_FORMAT"
	AccountParameterGeographyOutputFormat                 AccountParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	AccountParameterDateOutputFormat                      AccountParameter = "DATE_OUTPUT_FORMAT"
	AccountParameterErrorOnNondeterministicMerge          AccountParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	AccountParameterErrorOnNondeterministicUpdate         AccountParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	AccountParameterJSONIndent                            AccountParameter = "JSON_INDENT"
	AccountParameterLockTimeout                           AccountParameter = "LOCK_TIMEOUT"
	AccountParameterMultiStatementCount                   AccountParameter = "MULTI_STATEMENT_COUNT"
	AccountParameterQueryTag                              AccountParameter = "QUERY_TAG"
	AccountParameterQuotedIdentifiersIgnoreCase           AccountParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	AccountParameterRowsPerResultset                      AccountParameter = "ROWS_PER_RESULTSET"
	AccountParameterS3StageVpceDnsName                    AccountParameter = "S3_STAGE_VPCE_DNS_NAME"
	AccountParameterSimulatedDataSharingConsumer          AccountParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	AccountParameterStatementTimeoutInSeconds             AccountParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	AccountParameterStrictJSONOutput                      AccountParameter = "STRICT_JSON_OUTPUT"
	AccountParameterTimeInputFormat                       AccountParameter = "TIME_INPUT_FORMAT"
	AccountParameterTimeOutputFormat                      AccountParameter = "TIME_OUTPUT_FORMAT"
	AccountParameterTimestampDayIsAlways24h               AccountParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	AccountParameterTimestampInputFormat                  AccountParameter = "TIMESTAMP_INPUT_FORMAT"
	AccountParameterTimestampLtzOutputFormat              AccountParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	AccountParameterTimestampNtzOutputFormat              AccountParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	AccountParameterTimestampOutputFormat                 AccountParameter = "TIMESTAMP_OUTPUT_FORMAT"
	AccountParameterTimestampTypeMapping                  AccountParameter = "TIMESTAMP_TYPE_MAPPING"
	AccountParameterTimestampTzOutputFormat               AccountParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	AccountParameterTimezone                              AccountParameter = "TIMEZONE"
	AccountParameterTransactionAbortOnError               AccountParameter = "TRANSACTION_ABORT_ON_ERROR"
	AccountParameterTransactionDefaultIsolationLevel      AccountParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	AccountParameterTwoDigitCenturyStart                  AccountParameter = "TWO_DIGIT_CENTURY_START"
	AccountParameterUnsupportedDdlAction                  AccountParameter = "UNSUPPORTED_DDL_ACTION"
	AccountParameterUseCachedResult                       AccountParameter = "USE_CACHED_RESULT"
	AccountParameterWeekOfYearPolicy                      AccountParameter = "WEEK_OF_YEAR_POLICY"
	AccountParameterWeekStart                             AccountParameter = "WEEK_START"

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
	SessionParameterAbortDetachedQuery                    SessionParameter = "ABORT_DETACHED_QUERY"
	SessionParameterAutocommit                            SessionParameter = "AUTOCOMMIT"
	SessionParameterBinaryInputFormat                     SessionParameter = "BINARY_INPUT_FORMAT"
	SessionParameterBinaryOutputFormat                    SessionParameter = "BINARY_OUTPUT_FORMAT"
	SessionParameterClientMetadataRequestUseConnectionCtx SessionParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	SessionParameterClientMetadataUseSessionDatabase      SessionParameter = "CLIENT_METADATA_USE_SESSION_DATABASE"
	SessionParameterClientResultColumnCaseInsensitive     SessionParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	SessionParameterDateInputFormat                       SessionParameter = "DATE_INPUT_FORMAT"
	SessionParameterGeographyOutputFormat                 SessionParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	SessionParameterDateOutputFormat                      SessionParameter = "DATE_OUTPUT_FORMAT"
	SessionParameterErrorOnNondeterministicMerge          SessionParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	SessionParameterErrorOnNondeterministicUpdate         SessionParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	SessionParameterJSONIndent                            SessionParameter = "JSON_INDENT"
	SessionParameterLockTimeout                           SessionParameter = "LOCK_TIMEOUT"
	SessionParameterMultiStatementCount                   SessionParameter = "MULTI_STATEMENT_COUNT"
	SessionParameterQueryTag                              SessionParameter = "QUERY_TAG"
	SessionParameterQuotedIdentifiersIgnoreCase           SessionParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	SessionParameterRowsPerResultset                      SessionParameter = "ROWS_PER_RESULTSET"
	SessionParameterS3StageVpceDnsName                    SessionParameter = "S3_STAGE_VPCE_DNS_NAME"
	SessionParameterSimulatedDataSharingConsumer          SessionParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	SessionParameterStatementTimeoutInSeconds             SessionParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	SessionParameterStrictJSONOutput                      SessionParameter = "STRICT_JSON_OUTPUT"
	SessionParameterTimeInputFormat                       SessionParameter = "TIME_INPUT_FORMAT"
	SessionParameterTimeOutputFormat                      SessionParameter = "TIME_OUTPUT_FORMAT"
	SessionParameterTimestampDayIsAlways24h               SessionParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	SessionParameterTimestampInputFormat                  SessionParameter = "TIMESTAMP_INPUT_FORMAT"
	SessionParameterTimestampLTZOutputFormat              SessionParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	SessionParameterTimestampNTZOutputFormat              SessionParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	SessionParameterTimestampOutputFormat                 SessionParameter = "TIMESTAMP_OUTPUT_FORMAT"
	SessionParameterTimestampTypeMapping                  SessionParameter = "TIMESTAMP_TYPE_MAPPING"
	SessionParameterTimestampTZOutputFormat               SessionParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	SessionParameterTimezone                              SessionParameter = "TIMEZONE"
	SessionParameterTransactionAbortOnError               SessionParameter = "TRANSACTION_ABORT_ON_ERROR"
	SessionParameterTransactionDefaultIsolationLevel      SessionParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	SessionParameterTwoDigitCenturyStart                  SessionParameter = "TWO_DIGIT_CENTURY_START"
	SessionParameterUnsupportedDDLAction                  SessionParameter = "UNSUPPORTED_DDL_ACTION"
	SessionParameterUseCachedResult                       SessionParameter = "USE_CACHED_RESULT"
	SessionParameterWeekOfYearPolicy                      SessionParameter = "WEEK_OF_YEAR_POLICY"
	SessionParameterWeekStart                             SessionParameter = "WEEK_START"
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
	ObjectParameterCatalog                             ObjectParameter = "CATALOG"
	ObjectParameterExternalVolume                      ObjectParameter = "EXTERNAL_VOLUME"

	// User Parameters
	ObjectParameterEnableUnredactedQuerySyntaxError ObjectParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
)

type UserParameter string

const (
	// User Parameters
	UserParameterEnableUnredactedQuerySyntaxError UserParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"

	// Session Parameters (inherited)
	UserParameterAbortDetachedQuery                    UserParameter = "ABORT_DETACHED_QUERY"
	UserParameterAutocommit                            UserParameter = "AUTOCOMMIT"
	UserParameterBinaryInputFormat                     UserParameter = "BINARY_INPUT_FORMAT"
	UserParameterBinaryOutputFormat                    UserParameter = "BINARY_OUTPUT_FORMAT"
	UserParameterClientMetadataRequestUseConnectionCtx UserParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	UserParameterClientMetadataUseSessionDatabase      UserParameter = "CLIENT_METADATA_USE_SESSION_DATABASE"
	UserParameterClientResultColumnCaseInsensitive     UserParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	UserParameterDateInputFormat                       UserParameter = "DATE_INPUT_FORMAT"
	UserParameterDateOutputFormat                      UserParameter = "DATE_OUTPUT_FORMAT"
	UserParameterErrorOnNondeterministicMerge          UserParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	UserParameterErrorOnNondeterministicUpdate         UserParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	UserParameterGeographyOutputFormat                 UserParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	UserParameterJsonIndent                            UserParameter = "JSON_INDENT"
	UserParameterLockTimeout                           UserParameter = "LOCK_TIMEOUT"
	UserParameterMultiStatementCount                   UserParameter = "MULTI_STATEMENT_COUNT"
	UserParameterNetworkPolicy                         UserParameter = "NETWORK_POLICY"
	UserParameterQueryTag                              UserParameter = "QUERY_TAG"
	UserParameterQuotedIdentifiersIgnoreCase           UserParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	UserParameterRowsPerResultset                      UserParameter = "ROWS_PER_RESULTSET"
	UserParameterS3StageVpceDnsName                    UserParameter = "S3_STAGE_VPCE_DNS_NAME"
	UserParameterSimulatedDataSharingConsumer          UserParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	UserParameterStatementTimeoutInSeconds             UserParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	UserParameterStrictJsonOutput                      UserParameter = "STRICT_JSON_OUTPUT"
	UserParameterTimeInputFormat                       UserParameter = "TIME_INPUT_FORMAT"
	UserParameterTimeOutputFormat                      UserParameter = "TIME_OUTPUT_FORMAT"
	UserParameterTimestampDayIsAlways24h               UserParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	UserParameterTimestampInputFormat                  UserParameter = "TIMESTAMP_INPUT_FORMAT"
	UserParameterTimestampLtzOutputFormat              UserParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	UserParameterTimestampNtzOutputFormat              UserParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	UserParameterTimestampOutputFormat                 UserParameter = "TIMESTAMP_OUTPUT_FORMAT"
	UserParameterTimestampTypeMapping                  UserParameter = "TIMESTAMP_TYPE_MAPPING"
	UserParameterTimestampTzOutputFormat               UserParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	UserParameterTimezone                              UserParameter = "TIMEZONE"
	UserParameterTransactionDefaultIsolationLevel      UserParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	UserParameterTwoDigitCenturyStart                  UserParameter = "TWO_DIGIT_CENTURY_START"
	UserParameterUnsupportedDdlAction                  UserParameter = "UNSUPPORTED_DDL_ACTION"
	UserParameterUseCachedResult                       UserParameter = "USE_CACHED_RESULT"
	UserParameterWeekOfYearPolicy                      UserParameter = "WEEK_OF_YEAR_POLICY"
	UserParameterWeekStart                             UserParameter = "WEEK_START"
)

// AccountParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#account-parameters.
type AccountParameters struct {
	// Account Parameters
	AllowClientMFACaching                            *bool    `ddl:"parameter" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                     *bool    `ddl:"parameter" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                          *int     `ddl:"parameter" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableIdentifierFirstLogin                       *bool    `ddl:"parameter" sql:"ENABLE_IDENTIFIER_FIRST_LOGIN"`
	EnableInternalStagesPrivatelink                  *bool    `ddl:"parameter" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EnableUnredactedQuerySyntaxError                 *bool    `ddl:"parameter" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	EnableTriSecretAndRekeyOptOutForImageRepository  *bool    `ddl:"parameter" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"`
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage *bool    `ddl:"parameter" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE"`
	EventTable                                       *string  `ddl:"parameter,single_quotes" sql:"EVENT_TABLE"`
	ExternalOAuthAddPrivilegedRolesToBlockedList     *bool    `ddl:"parameter" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTB                  *float64 `ddl:"parameter" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                       *int     `ddl:"parameter" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                                    *string  `ddl:"parameter,single_quotes" sql:"NETWORK_POLICY"`
	PeriodicDataRekeying                             *bool    `ddl:"parameter" sql:"PERIODIC_DATA_REKEYING"`
	PreventLoadFromInlineURL                         *bool    `ddl:"parameter" sql:"PREVENT_LOAD_FROM_INLINE_URL"`
	PreventUnloadToInlineURL                         *bool    `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                    *bool    `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation        *bool    `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation       *bool    `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                                     *bool    `ddl:"parameter" sql:"SSO_LOGIN_PAGE"`
}

func (v *AccountParameters) validate() error {
	var errs []error
	if valueSet(v.ClientEncryptionKeySize) {
		if !(*v.ClientEncryptionKeySize == 128 || *v.ClientEncryptionKeySize == 256) {
			errs = append(errs, fmt.Errorf("CLIENT_ENCRYPTION_KEY_SIZE must be either 128 or 256"))
		}
	}
	if valueSet(v.InitialReplicationSizeLimitInTB) {
		l := *v.InitialReplicationSizeLimitInTB
		if l < 0.0 || (l < 0.0 && l < 1.0) {
			errs = append(errs, fmt.Errorf("%v must be 0.0 and above with a scale of at least 1 (e.g. 20.5, 32.25, 33.333, etc.)", l))
		}
	}
	if valueSet(v.MinDataRetentionTimeInDays) {
		if !validateIntInRange(*v.MinDataRetentionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("AccountParameters", "MinDataRetentionTimeInDays", 0, 90))
		}
	}
	return errors.Join(errs...)
}

type AccountParametersUnset struct {
	AllowClientMFACaching                            *bool `ddl:"keyword" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                     *bool `ddl:"keyword" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                          *bool `ddl:"keyword" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	EnableIdentifierFirstLogin                       *bool `ddl:"keyword" sql:"ENABLE_IDENTIFIER_FIRST_LOGIN"`
	EnableInternalStagesPrivatelink                  *bool `ddl:"keyword" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EnableTriSecretAndRekeyOptOutForImageRepository  *bool `ddl:"keyword" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"`
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage *bool `ddl:"keyword" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE"`
	EventTable                                       *bool `ddl:"keyword" sql:"EVENT_TABLE"`
	ExternalOAuthAddPrivilegedRolesToBlockedList     *bool `ddl:"keyword" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTB                  *bool `ddl:"keyword" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                       *bool `ddl:"keyword" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                                    *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	PeriodicDataRekeying                             *bool `ddl:"keyword" sql:"PERIODIC_DATA_REKEYING"`
	PreventUnloadToInlineURL                         *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                    *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation        *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation       *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                                     *bool `ddl:"keyword" sql:"SSO_LOGIN_PAGE"`
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

type BinaryOutputFormat string

const (
	BinaryOutputFormatHex    BinaryOutputFormat = "HEX"
	BinaryOutputFormatBase64 BinaryOutputFormat = "BASE64"
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

// SessionParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#session-parameters.
type SessionParameters struct {
	AbortDetachedQuery                    *bool                             `ddl:"parameter" sql:"ABORT_DETACHED_QUERY"`
	Autocommit                            *bool                             `ddl:"parameter" sql:"AUTOCOMMIT"`
	BinaryInputFormat                     *BinaryInputFormat                `ddl:"parameter,single_quotes" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat                    *BinaryOutputFormat               `ddl:"parameter,single_quotes" sql:"BINARY_OUTPUT_FORMAT"`
	ClientMetadataRequestUseConnectionCtx *bool                             `ddl:"parameter" sql:"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"`
	ClientMetadataUseSessionDatabase      *bool                             `ddl:"parameter" sql:"CLIENT_METADATA_USE_SESSION_DATABASE"`
	ClientResultColumnCaseInsensitive     *bool                             `ddl:"parameter" sql:"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"`
	DateInputFormat                       *string                           `ddl:"parameter,single_quotes" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                      *string                           `ddl:"parameter,single_quotes" sql:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge          *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate         *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat                 *GeographyOutputFormat            `ddl:"parameter,single_quotes" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	JSONIndent                            *int                              `ddl:"parameter" sql:"JSON_INDENT"`
	LockTimeout                           *int                              `ddl:"parameter" sql:"LOCK_TIMEOUT"`
	MultiStatementCount                   *int                              `ddl:"parameter" sql:"MULTI_STATEMENT_COUNT"`
	QueryTag                              *string                           `ddl:"parameter,single_quotes" sql:"QUERY_TAG"`
	QuotedIdentifiersIgnoreCase           *bool                             `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	RowsPerResultset                      *int                              `ddl:"parameter" sql:"ROWS_PER_RESULTSET"`
	S3StageVpceDnsName                    *string                           `ddl:"parameter,single_quotes" sql:"S3_STAGE_VPCE_DNS_NAME"`
	SimulatedDataSharingConsumer          *string                           `ddl:"parameter,single_quotes" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementTimeoutInSeconds             *int                              `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJSONOutput                      *bool                             `ddl:"parameter" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h               *bool                             `ddl:"parameter" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat                  *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat              *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat              *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat                 *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping                  *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat               *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                              *string                           `ddl:"parameter,single_quotes" sql:"TIMEZONE"`
	TimeInputFormat                       *string                           `ddl:"parameter,single_quotes" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                      *string                           `ddl:"parameter,single_quotes" sql:"TIME_OUTPUT_FORMAT"`
	TransactionAbortOnError               *bool                             `ddl:"parameter" sql:"TRANSACTION_ABORT_ON_ERROR"`
	TransactionDefaultIsolationLevel      *TransactionDefaultIsolationLevel `ddl:"parameter,single_quotes" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart                  *int                              `ddl:"parameter" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction                  *UnsupportedDDLAction             `ddl:"parameter,single_quotes" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                       *bool                             `ddl:"parameter" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                      *int                              `ddl:"parameter" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                             *int                              `ddl:"parameter" sql:"WEEK_START"`
}

func (v *SessionParameters) validate() error {
	var errs []error
	if valueSet(v.JSONIndent) {
		if !validateIntInRange(*v.JSONIndent, 0, 16) {
			errs = append(errs, errIntBetween("SessionParameters", "JSONIndent", 0, 16))
		}
	}
	if valueSet(v.LockTimeout) {
		if !validateIntGreaterThanOrEqual(*v.LockTimeout, 0) {
			errs = append(errs, errIntValue("SessionParameters", "LockTimeout", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.QueryTag) {
		if len(*v.QueryTag) > 2000 {
			errs = append(errs, errIntValue("SessionParameters", "QueryTag", IntErrLess, 2000))
		}
	}
	if valueSet(v.RowsPerResultset) {
		if !validateIntGreaterThanOrEqual(*v.RowsPerResultset, 0) {
			errs = append(errs, errIntValue("SessionParameters", "RowsPerResultset", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.TwoDigitCenturyStart) {
		if !validateIntInRange(*v.TwoDigitCenturyStart, 1900, 2100) {
			errs = append(errs, errIntBetween("SessionParameters", "TwoDigitCenturyStart", 1900, 2100))
		}
	}
	if valueSet(v.WeekOfYearPolicy) {
		if !validateIntInRange(*v.WeekOfYearPolicy, 0, 1) {
			errs = append(errs, fmt.Errorf("WEEK_OF_YEAR_POLICY must be either 0 or 1"))
		}
	}
	if valueSet(v.WeekStart) {
		if !validateIntInRange(*v.WeekStart, 0, 1) {
			errs = append(errs, fmt.Errorf("WEEK_START must be either 0 or 1"))
		}
	}
	return errors.Join(errs...)
}

type SessionParametersUnset struct {
	AbortDetachedQuery                    *bool `ddl:"keyword" sql:"ABORT_DETACHED_QUERY"`
	Autocommit                            *bool `ddl:"keyword" sql:"AUTOCOMMIT"`
	BinaryInputFormat                     *bool `ddl:"keyword" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat                    *bool `ddl:"keyword" sql:"BINARY_OUTPUT_FORMAT"`
	ClientMetadataRequestUseConnectionCtx *bool `ddl:"keyword" sql:"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"`
	ClientMetadataUseSessionDatabase      *bool `ddl:"keyword" sql:"CLIENT_METADATA_USE_SESSION_DATABASE"`
	ClientResultColumnCaseInsensitive     *bool `ddl:"keyword" sql:"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"`
	DateInputFormat                       *bool `ddl:"keyword" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                      *bool `ddl:"keyword" sql:"DATE_OUTPUT_FORMAT"`
	ErrorOnNondeterministicMerge          *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate         *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat                 *bool `ddl:"keyword" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	JSONIndent                            *bool `ddl:"keyword" sql:"JSON_INDENT"`
	LockTimeout                           *bool `ddl:"keyword" sql:"LOCK_TIMEOUT"`
	MultiStatementCount                   *bool `ddl:"keyword" sql:"MULTI_STATEMENT_COUNT"`
	QueryTag                              *bool `ddl:"keyword" sql:"QUERY_TAG"`
	QuotedIdentifiersIgnoreCase           *bool `ddl:"keyword" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	RowsPerResultset                      *bool `ddl:"keyword" sql:"ROWS_PER_RESULTSET"`
	S3StageVpceDnsName                    *bool `ddl:"keyword" sql:"S3_STAGE_VPCE_DNS_NAME"`
	SimulatedDataSharingConsumer          *bool `ddl:"keyword" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementTimeoutInSeconds             *bool `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJSONOutput                      *bool `ddl:"keyword" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h               *bool `ddl:"keyword" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat                  *bool `ddl:"keyword" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat              *bool `ddl:"keyword" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat              *bool `ddl:"keyword" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat                 *bool `ddl:"keyword" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping                  *bool `ddl:"keyword" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat               *bool `ddl:"keyword" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                              *bool `ddl:"keyword" sql:"TIMEZONE"`
	TimeInputFormat                       *bool `ddl:"keyword" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                      *bool `ddl:"keyword" sql:"TIME_OUTPUT_FORMAT"`
	TransactionDefaultIsolationLevel      *bool `ddl:"keyword" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart                  *bool `ddl:"keyword" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction                  *bool `ddl:"keyword" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                       *bool `ddl:"keyword" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                      *bool `ddl:"keyword" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                             *bool `ddl:"keyword" sql:"WEEK_START"`
}

func (v *SessionParametersUnset) validate() error {
	if !anyValueSet(v.AbortDetachedQuery, v.Autocommit, v.BinaryInputFormat, v.BinaryOutputFormat, v.ClientMetadataRequestUseConnectionCtx, v.ClientMetadataUseSessionDatabase, v.ClientResultColumnCaseInsensitive, v.DateInputFormat, v.DateOutputFormat, v.ErrorOnNondeterministicMerge, v.ErrorOnNondeterministicUpdate, v.GeographyOutputFormat, v.JSONIndent, v.LockTimeout, v.MultiStatementCount, v.QueryTag, v.QuotedIdentifiersIgnoreCase, v.RowsPerResultset, v.SimulatedDataSharingConsumer, v.StatementTimeoutInSeconds, v.StrictJSONOutput, v.TimestampDayIsAlways24h, v.TimestampInputFormat, v.TimestampLTZOutputFormat, v.TimestampNTZOutputFormat, v.TimestampOutputFormat, v.TimestampTypeMapping, v.TimestampTZOutputFormat, v.Timezone, v.TimeInputFormat, v.TimeOutputFormat, v.TransactionDefaultIsolationLevel, v.TwoDigitCenturyStart, v.UnsupportedDDLAction, v.UseCachedResult, v.WeekOfYearPolicy, v.WeekStart) {
		return errors.Join(errAtLeastOneOf("SessionParametersUnset", "AbortDetachedQuery", "Autocommit", "BinaryInputFormat", "BinaryOutputFormat", "DateInputFormat", "DateOutputFormat", "ErrorOnNondeterministicMerge", "ErrorOnNondeterministicUpdate", "GeographyOutputFormat", "JSONIndent", "LockTimeout", "QueryTag", "RowsPerResultset", "SimulatedDataSharingConsumer", "StatementTimeoutInSeconds", "StrictJSONOutput", "TimestampDayIsAlways24h", "TimestampInputFormat", "TimestampLTZOutputFormat", "TimestampNTZOutputFormat", "TimestampOutputFormat", "TimestampTypeMapping", "TimestampTZOutputFormat", "Timezone", "TimeInputFormat", "TimeOutputFormat", "TransactionDefaultIsolationLevel", "TwoDigitCenturyStart", "UnsupportedDDLAction", "UseCachedResult", "WeekOfYearPolicy", "WeekStart"))
	}
	return nil
}

// ObjectParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#object-parameters.
type ObjectParameters struct {
	DataRetentionTimeInDays             *int           `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	DefaultDDLCollation                 *string        `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	EnableUnredactedQuerySyntaxError    *bool          `ddl:"parameter" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
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
	var errs []error
	if valueSet(v.DataRetentionTimeInDays) {
		if !validateIntInRange(*v.DataRetentionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("ObjectParameters", "DataRetentionTimeInDays", 0, 90))
		}
	}
	if valueSet(v.MaxConcurrencyLevel) {
		if !validateIntGreaterThanOrEqual(*v.MaxConcurrencyLevel, 1) {
			errs = append(errs, errIntValue("ObjectParameters", "MaxConcurrencyLevel", IntErrGreaterOrEqual, 1))
		}
	}
	if valueSet(v.MaxDataExtensionTimeInDays) {
		if !validateIntInRange(*v.MaxDataExtensionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("ObjectParameters", "MaxDataExtensionTimeInDays", 0, 90))
		}
	}
	if valueSet(v.StatementQueuedTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementQueuedTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "StatementQueuedTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.SuspendTaskAfterNumFailures) {
		if !validateIntGreaterThanOrEqual(*v.SuspendTaskAfterNumFailures, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "SuspendTaskAfterNumFailures", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.UserTaskTimeoutMs) {
		if !validateIntInRange(*v.UserTaskTimeoutMs, 0, 86400000) {
			errs = append(errs, errIntBetween("ObjectParameters", "UserTaskTimeoutMs", 0, 86400000))
		}
	}
	return errors.Join(errs...)
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
	NetworkPolicy                       *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
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

// ShowParametersOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-parameters.
type ShowParametersOptions struct {
	show       bool          `ddl:"static" sql:"SHOW"`
	parameters bool          `ddl:"static" sql:"PARAMETERS"`
	Like       *Like         `ddl:"keyword" sql:"LIKE"`
	In         *ParametersIn `ddl:"keyword" sql:"IN"`
}

func (opts *ShowParametersOptions) validate() error {
	if valueSet(opts.In) {
		if err := opts.In.validate(); err != nil {
			return err
		}
	}
	return nil
}

type ParametersIn struct {
	Session   *bool                    `ddl:"keyword" sql:"SESSION"`
	Account   *bool                    `ddl:"keyword" sql:"ACCOUNT"`
	User      AccountObjectIdentifier  `ddl:"identifier" sql:"USER"`
	Warehouse AccountObjectIdentifier  `ddl:"identifier" sql:"WAREHOUSE"`
	Database  AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
	Schema    DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	Task      SchemaObjectIdentifier   `ddl:"identifier" sql:"TASK"`
	Table     SchemaObjectIdentifier   `ddl:"identifier" sql:"TABLE"`
}

func (v *ParametersIn) validate() error {
	if !anyValueSet(v.Session, v.Account, v.User, v.Warehouse, v.Database, v.Schema, v.Task, v.Table) {
		return errors.Join(errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table"))
	}
	return nil
}

type ParameterType string

const (
	ParameterTypeAccount ParameterType = "ACCOUNT"
	ParameterTypeUser    ParameterType = "USER"
	ParameterTypeSession ParameterType = "SESSION"
	ParameterTypeObject  ParameterType = "OBJECT"
)

type Parameter struct {
	Key         string
	Value       string
	Default     string
	Level       ParameterType
	Description string
}

type parameterRow struct {
	Key         sql.NullString `db:"key"`
	Value       sql.NullString `db:"value"`
	Default     sql.NullString `db:"default"`
	Level       sql.NullString `db:"level"`
	Description sql.NullString `db:"description"`
}

func (row *parameterRow) toParameter() *Parameter {
	return &Parameter{
		Key:         row.Key.String,
		Value:       row.Value.String,
		Default:     row.Default.String,
		Level:       ParameterType(row.Level.String),
		Description: row.Description.String,
	}
}

func (v *parameters) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	if opts == nil {
		opts = &ShowParametersOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	rows := []parameterRow{}
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	parameters := make([]*Parameter, len(rows))
	for i, row := range rows {
		parameters[i] = row.toParameter()
	}
	return parameters, nil
}

func (v *parameters) ShowAccountParameter(ctx context.Context, parameter AccountParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Account: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *parameters) ShowSessionParameter(ctx context.Context, parameter SessionParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Session: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *parameters) ShowUserParameter(ctx context.Context, parameter UserParameter, userId AccountObjectIdentifier) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			User: userId,
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *parameters) ShowObjectParameter(ctx context.Context, parameter ObjectParameter, object Object) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{},
	}
	switch object.ObjectType {
	case ObjectTypeWarehouse:
		opts.In.Warehouse = object.Name.(AccountObjectIdentifier)
	case ObjectTypeDatabase:
		opts.In.Database = object.Name.(AccountObjectIdentifier)
	case ObjectTypeSchema:
		opts.In.Schema = object.Name.(DatabaseObjectIdentifier)
	case ObjectTypeTask:
		opts.In.Task = object.Name.(SchemaObjectIdentifier)
	case ObjectTypeTable:
		opts.In.Table = object.Name.(SchemaObjectIdentifier)
	case ObjectTypeUser:
		opts.In.User = object.Name.(AccountObjectIdentifier)
	default:
		return nil, fmt.Errorf("unsupported object type %s", object.Name)
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}
