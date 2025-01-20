package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
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
	UnsetAccountParameter(ctx context.Context, parameter AccountParameter) error
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
	case AccountParameterCortexEnabledCrossRegion:
		opts.Set.Parameters.AccountParameters.CortexEnabledCrossRegion = &value
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
	case AccountParameterEnablePersonalDatabase:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnablePersonalDatabase = b
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
	case AccountParameterEnableUnhandledExceptionsReporting:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnableUnhandledExceptionsReporting = b
	case AccountParameterEnforceNetworkRulesForInternalStages:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.EnforceNetworkRulesForInternalStages = b
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
		opts.Set.Parameters.AccountParameters.InitialReplicationSizeLimitInTB = &value
	case AccountParameterMinDataRetentionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MIN_DATA_RETENTION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.AccountParameters.MinDataRetentionTimeInDays = Pointer(v)
	case AccountParameterNetworkPolicy:
		opts.Set.Parameters.AccountParameters.NetworkPolicy = &value
	case AccountParameterOAuthAddPrivilegedRolesToBlockedList:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.AccountParameters.OAuthAddPrivilegedRolesToBlockedList = b
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

// TODO(SNOW-1866453): add integration tests
func (parameters *parameters) UnsetAccountParameter(ctx context.Context, parameter AccountParameter) error {
	opts := AlterAccountOptions{
		Unset: &AccountUnset{
			Parameters: &AccountLevelParametersUnset{
				AccountParameters: &AccountParametersUnset{},
			},
		},
	}
	switch parameter {
	case AccountParameterAllowClientMFACaching:
		opts.Unset.Parameters.AccountParameters.AllowClientMFACaching = Pointer(true)
	case AccountParameterAllowIDToken:
		opts.Unset.Parameters.AccountParameters.AllowIDToken = Pointer(true)
	case AccountParameterClientEncryptionKeySize:
		opts.Unset.Parameters.AccountParameters.ClientEncryptionKeySize = Pointer(true)
	case AccountParameterCortexEnabledCrossRegion:
		opts.Unset.Parameters.AccountParameters.CortexEnabledCrossRegion = Pointer(true)
	case AccountParameterEnableIdentifierFirstLogin:
		opts.Unset.Parameters.AccountParameters.EnableIdentifierFirstLogin = Pointer(true)
	case AccountParameterEnableInternalStagesPrivatelink:
		opts.Unset.Parameters.AccountParameters.EnableInternalStagesPrivatelink = Pointer(true)
	case AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository:
		opts.Unset.Parameters.AccountParameters.EnableTriSecretAndRekeyOptOutForImageRepository = Pointer(true)
	case AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage:
		opts.Unset.Parameters.AccountParameters.EnableTriSecretAndRekeyOptOutForSpcsBlockStorage = Pointer(true)
	case AccountParameterEnablePersonalDatabase:
		opts.Unset.Parameters.AccountParameters.EnablePersonalDatabase = Pointer(true)
	case AccountParameterEnableUnhandledExceptionsReporting:
		opts.Unset.Parameters.AccountParameters.EnableUnhandledExceptionsReporting = Pointer(true)
	case AccountParameterEnableUnredactedQuerySyntaxError:
		opts.Unset.Parameters.AccountParameters.EnableUnredactedQuerySyntaxError = Pointer(true)
	case AccountParameterEnforceNetworkRulesForInternalStages:
		opts.Unset.Parameters.AccountParameters.EnforceNetworkRulesForInternalStages = Pointer(true)
	case AccountParameterEventTable:
		opts.Unset.Parameters.AccountParameters.EventTable = Pointer(true)
	case AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList:
		opts.Unset.Parameters.AccountParameters.ExternalOAuthAddPrivilegedRolesToBlockedList = Pointer(true)
	case AccountParameterInitialReplicationSizeLimitInTB:
		opts.Unset.Parameters.AccountParameters.InitialReplicationSizeLimitInTB = Pointer(true)
	case AccountParameterMinDataRetentionTimeInDays:
		opts.Unset.Parameters.AccountParameters.MinDataRetentionTimeInDays = Pointer(true)
	case AccountParameterNetworkPolicy:
		opts.Unset.Parameters.AccountParameters.NetworkPolicy = Pointer(true)
	case AccountParameterOAuthAddPrivilegedRolesToBlockedList:
		opts.Unset.Parameters.AccountParameters.OAuthAddPrivilegedRolesToBlockedList = Pointer(true)
	case AccountParameterPeriodicDataRekeying:
		opts.Unset.Parameters.AccountParameters.PeriodicDataRekeying = Pointer(true)
	case AccountParameterPreventLoadFromInlineURL:
		opts.Unset.Parameters.AccountParameters.PreventLoadFromInlineURL = Pointer(true)
	case AccountParameterPreventUnloadToInlineURL:
		opts.Unset.Parameters.AccountParameters.PreventUnloadToInlineURL = Pointer(true)
	case AccountParameterPreventUnloadToInternalStages:
		opts.Unset.Parameters.AccountParameters.PreventUnloadToInternalStages = Pointer(true)
	case AccountParameterRequireStorageIntegrationForStageCreation:
		opts.Unset.Parameters.AccountParameters.RequireStorageIntegrationForStageCreation = Pointer(true)
	case AccountParameterRequireStorageIntegrationForStageOperation:
		opts.Unset.Parameters.AccountParameters.RequireStorageIntegrationForStageOperation = Pointer(true)
	case AccountParameterSSOLoginPage:
		opts.Unset.Parameters.AccountParameters.SSOLoginPage = Pointer(true)
	default:
		return parameters.UnsetSessionParameterOnAccount(ctx, SessionParameter(parameter))
	}
	return parameters.client.Accounts.Alter(ctx, &opts)
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

func (parameters *parameters) UnsetSessionParameterOnAccount(ctx context.Context, parameter SessionParameter) error {
	sp := &SessionParametersUnset{}
	err := sp.setParam(parameter)
	if err == nil {
		opts := AlterAccountOptions{Unset: &AccountUnset{Parameters: &AccountLevelParametersUnset{SessionParameters: sp}}}
		err = parameters.client.Accounts.Alter(ctx, &opts)
		if err != nil {
			return err
		}
		return nil
	} else {
		if strings.Contains(err.Error(), "session parameter is not supported") {
			return parameters.UnsetObjectParameterOnAccount(ctx, ObjectParameter(parameter))
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
			return fmt.Errorf("DATA_RETENTION_TIME_IN_DAYS object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.DataRetentionTimeInDays = Pointer(v)
	case ObjectParameterDefaultDDLCollation:
		opts.Set.Parameters.ObjectParameters.DefaultDDLCollation = &value
	case ObjectParameterLogLevel:
		opts.Set.Parameters.ObjectParameters.LogLevel = Pointer(LogLevel(value))
	case ObjectParameterMaxConcurrencyLevel:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MAX_CONCURRENCY_LEVEL object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.MaxConcurrencyLevel = Pointer(v)
	case ObjectParameterMaxDataExtensionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("MAX_DATA_EXTENSION_TIME_IN_DAYS object parameter is an integer, got %v", value)
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
			return fmt.Errorf("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.StatementQueuedTimeoutInSeconds = Pointer(v)
	case ObjectParameterStatementTimeoutInSeconds:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("STATEMENT_TIMEOUT_IN_SECONDS object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.StatementTimeoutInSeconds = Pointer(v)
	case ObjectParameterNetworkPolicy:
		opts.Set.Parameters.ObjectParameters.NetworkPolicy = &value
	case ObjectParameterShareRestrictions:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.ShareRestrictions = b
	case ObjectParameterStorageSerializationPolicy:
		opts.Set.Parameters.ObjectParameters.StorageSerializationPolicy = &value
	case ObjectParameterSuspendTaskAfterNumFailures:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("SUSPEND_TASK_AFTER_NUM_FAILURES object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.SuspendTaskAfterNumFailures = Pointer(v)
	case ObjectParameterTaskAutoRetryAttempts:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("TASK_AUTO_RETRY_ATTEMPTS object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.TaskAutoRetryAttempts = Pointer(v)
	case ObjectParameterTraceLevel:
		opts.Set.Parameters.ObjectParameters.TraceLevel = Pointer(TraceLevel(value))
	case ObjectParameterUserTaskManagedInitialWarehouseSize:
		opts.Set.Parameters.ObjectParameters.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSize(value))
	case ObjectParameterUserTaskTimeoutMs:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("USER_TASK_TIMEOUT_MS object parameter is an integer, got %v", value)
		}
		opts.Set.Parameters.ObjectParameters.UserTaskTimeoutMs = Pointer(v)
	case ObjectParameterEnableUnredactedQuerySyntaxError:
		b, err := parseBooleanParameter(string(parameter), value)
		if err != nil {
			return err
		}
		opts.Set.Parameters.ObjectParameters.EnableUnredactedQuerySyntaxError = b
	case ObjectParameterCatalog:
		opts.Set.Parameters.ObjectParameters.Catalog = &value
	default:
		return fmt.Errorf("Invalid object parameter: %v", string(parameter))
	}
	err := parameters.client.Accounts.Alter(ctx, &opts)
	if err != nil {
		return err
	}
	return nil
}

func (parameters *parameters) UnsetObjectParameterOnAccount(ctx context.Context, parameter ObjectParameter) error {
	opts := AlterAccountOptions{Unset: &AccountUnset{Parameters: &AccountLevelParametersUnset{ObjectParameters: &ObjectParametersUnset{}}}}
	switch parameter {
	case ObjectParameterCatalog:
		opts.Unset.Parameters.ObjectParameters.Catalog = Pointer(true)
	case ObjectParameterDataRetentionTimeInDays:
		opts.Unset.Parameters.ObjectParameters.DataRetentionTimeInDays = Pointer(true)
	case ObjectParameterDefaultDDLCollation:
		opts.Unset.Parameters.ObjectParameters.DefaultDDLCollation = Pointer(true)
	case ObjectParameterLogLevel:
		opts.Unset.Parameters.ObjectParameters.LogLevel = Pointer(true)
	case ObjectParameterMaxConcurrencyLevel:
		opts.Unset.Parameters.ObjectParameters.MaxConcurrencyLevel = Pointer(true)
	case ObjectParameterMaxDataExtensionTimeInDays:
		opts.Unset.Parameters.ObjectParameters.MaxDataExtensionTimeInDays = Pointer(true)
	case ObjectParameterPipeExecutionPaused:
		opts.Unset.Parameters.ObjectParameters.PipeExecutionPaused = Pointer(true)
	case ObjectParameterPreventUnloadToInternalStages:
		opts.Unset.Parameters.ObjectParameters.PreventUnloadToInternalStages = Pointer(true)
	case ObjectParameterStatementQueuedTimeoutInSeconds:
		opts.Unset.Parameters.ObjectParameters.StatementQueuedTimeoutInSeconds = Pointer(true)
	case ObjectParameterStatementTimeoutInSeconds:
		opts.Unset.Parameters.ObjectParameters.StatementTimeoutInSeconds = Pointer(true)
	case ObjectParameterNetworkPolicy:
		opts.Unset.Parameters.ObjectParameters.NetworkPolicy = Pointer(true)
	case ObjectParameterShareRestrictions:
		opts.Unset.Parameters.ObjectParameters.ShareRestrictions = Pointer(true)
	case ObjectParameterStorageSerializationPolicy:
		opts.Unset.Parameters.ObjectParameters.StorageSerializationPolicy = Pointer(true)
	case ObjectParameterSuspendTaskAfterNumFailures:
		opts.Unset.Parameters.ObjectParameters.SuspendTaskAfterNumFailures = Pointer(true)
	case ObjectParameterTaskAutoRetryAttempts:
		opts.Unset.Parameters.ObjectParameters.TaskAutoRetryAttempts = Pointer(true)
	case ObjectParameterTraceLevel:
		opts.Unset.Parameters.ObjectParameters.TraceLevel = Pointer(true)
	case ObjectParameterUserTaskManagedInitialWarehouseSize:
		opts.Unset.Parameters.ObjectParameters.UserTaskManagedInitialWarehouseSize = Pointer(true)
	case ObjectParameterUserTaskTimeoutMs:
		opts.Unset.Parameters.ObjectParameters.UserTaskTimeoutMs = Pointer(true)
	case ObjectParameterEnableUnredactedQuerySyntaxError:
		opts.Unset.Parameters.ObjectParameters.EnableUnredactedQuerySyntaxError = Pointer(true)
	default:
		return fmt.Errorf("invalid object parameter: %v", string(parameter))
	}
	return parameters.client.Accounts.Alter(ctx, &opts)
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
	AccountParameterCortexEnabledCrossRegion                         AccountParameter = "CORTEX_ENABLED_CROSS_REGION"
	AccountParameterEnableIdentifierFirstLogin                       AccountParameter = "ENABLE_IDENTIFIER_FIRST_LOGIN"
	AccountParameterEnableInternalStagesPrivatelink                  AccountParameter = "ENABLE_INTERNAL_STAGES_PRIVATELINK"
	AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository  AccountParameter = "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"   // #nosec G101
	AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage AccountParameter = "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE" // #nosec G101
	AccountParameterEnableUnhandledExceptionsReporting               AccountParameter = "ENABLE_UNHANDLED_EXCEPTIONS_REPORTING"
	AccountParameterEnforceNetworkRulesForInternalStages             AccountParameter = "ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES"
	AccountParameterEventTable                                       AccountParameter = "EVENT_TABLE"
	AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList     AccountParameter = "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"
	AccountParameterInitialReplicationSizeLimitInTB                  AccountParameter = "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"
	AccountParameterMinDataRetentionTimeInDays                       AccountParameter = "MIN_DATA_RETENTION_TIME_IN_DAYS"
	AccountParameterNetworkPolicy                                    AccountParameter = "NETWORK_POLICY"
	AccountParameterOAuthAddPrivilegedRolesToBlockedList             AccountParameter = "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"
	AccountParameterPeriodicDataRekeying                             AccountParameter = "PERIODIC_DATA_REKEYING"
	AccountParameterPreventLoadFromInlineURL                         AccountParameter = "PREVENT_LOAD_FROM_INLINE_URL"
	AccountParameterPreventUnloadToInlineURL                         AccountParameter = "PREVENT_UNLOAD_TO_INLINE_URL"
	AccountParameterRequireStorageIntegrationForStageCreation        AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"
	AccountParameterRequireStorageIntegrationForStageOperation       AccountParameter = "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"
	AccountParameterSSOLoginPage                                     AccountParameter = "SSO_LOGIN_PAGE"

	// Session Parameters (inherited)
	AccountParameterAbortDetachedQuery                       AccountParameter = "ABORT_DETACHED_QUERY"
	AccountParameterActivePythonProfiler                     AccountParameter = "ACTIVE_PYTHON_PROFILER"
	AccountParameterAutocommit                               AccountParameter = "AUTOCOMMIT"
	AccountParameterBinaryInputFormat                        AccountParameter = "BINARY_INPUT_FORMAT"
	AccountParameterBinaryOutputFormat                       AccountParameter = "BINARY_OUTPUT_FORMAT"
	AccountParameterClientEnableLogInfoStatementParameters   AccountParameter = "CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS"
	AccountParameterClientMemoryLimit                        AccountParameter = "CLIENT_MEMORY_LIMIT"
	AccountParameterClientMetadataRequestUseConnectionCtx    AccountParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	AccountParameterClientMetadataUseSessionDatabase         AccountParameter = "CLIENT_METADATA_USE_SESSION_DATABASE"
	AccountParameterClientPrefetchThreads                    AccountParameter = "CLIENT_PREFETCH_THREADS"
	AccountParameterClientResultChunkSize                    AccountParameter = "CLIENT_RESULT_CHUNK_SIZE"
	AccountParameterClientResultColumnCaseInsensitive        AccountParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	AccountParameterClientSessionKeepAlive                   AccountParameter = "CLIENT_SESSION_KEEP_ALIVE"
	AccountParameterClientSessionKeepAliveHeartbeatFrequency AccountParameter = "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"
	AccountParameterClientTimestampTypeMapping               AccountParameter = "CLIENT_TIMESTAMP_TYPE_MAPPING"
	AccountParameterCsvTimestampFormat                       AccountParameter = "CSV_TIMESTAMP_FORMAT"
	AccountParameterDateInputFormat                          AccountParameter = "DATE_INPUT_FORMAT"
	AccountParameterDateOutputFormat                         AccountParameter = "DATE_OUTPUT_FORMAT"
	AccountParameterEnableUnloadPhysicalTypeOptimization     AccountParameter = "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"
	AccountParameterErrorOnNondeterministicMerge             AccountParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	AccountParameterErrorOnNondeterministicUpdate            AccountParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	AccountParameterGeographyOutputFormat                    AccountParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	AccountParameterGeometryOutputFormat                     AccountParameter = "GEOMETRY_OUTPUT_FORMAT"
	AccountParameterHybridTableLockTimeout                   AccountParameter = "HYBRID_TABLE_LOCK_TIMEOUT"
	AccountParameterJdbcTreatDecimalAsInt                    AccountParameter = "JDBC_TREAT_DECIMAL_AS_INT"
	AccountParameterJdbcTreatTimestampNtzAsUtc               AccountParameter = "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"
	AccountParameterJdbcUseSessionTimezone                   AccountParameter = "JDBC_USE_SESSION_TIMEZONE"
	AccountParameterJsonIndent                               AccountParameter = "JSON_INDENT"
	AccountParameterJsTreatIntegerAsBigInt                   AccountParameter = "JS_TREAT_INTEGER_AS_BIGINT"
	AccountParameterLockTimeout                              AccountParameter = "LOCK_TIMEOUT"
	AccountParameterMultiStatementCount                      AccountParameter = "MULTI_STATEMENT_COUNT"
	AccountParameterNoorderSequenceAsDefault                 AccountParameter = "NOORDER_SEQUENCE_AS_DEFAULT"
	AccountParameterOdbcTreatDecimalAsInt                    AccountParameter = "ODBC_TREAT_DECIMAL_AS_INT"
	AccountParameterPythonProfilerModules                    AccountParameter = "PYTHON_PROFILER_MODULES"
	AccountParameterPythonProfilerTargetStage                AccountParameter = "PYTHON_PROFILER_TARGET_STAGE"
	AccountParameterQueryTag                                 AccountParameter = "QUERY_TAG"
	AccountParameterQuotedIdentifiersIgnoreCase              AccountParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	AccountParameterRowsPerResultset                         AccountParameter = "ROWS_PER_RESULTSET"
	AccountParameterS3StageVpceDnsName                       AccountParameter = "S3_STAGE_VPCE_DNS_NAME"
	AccountParameterSearchPath                               AccountParameter = "SEARCH_PATH"
	AccountParameterSimulatedDataSharingConsumer             AccountParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	AccountParameterStrictJsonOutput                         AccountParameter = "STRICT_JSON_OUTPUT"
	AccountParameterTimeInputFormat                          AccountParameter = "TIME_INPUT_FORMAT"
	AccountParameterTimeOutputFormat                         AccountParameter = "TIME_OUTPUT_FORMAT"
	AccountParameterTimestampDayIsAlways24h                  AccountParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	AccountParameterTimestampInputFormat                     AccountParameter = "TIMESTAMP_INPUT_FORMAT"
	AccountParameterTimestampLtzOutputFormat                 AccountParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	AccountParameterTimestampNtzOutputFormat                 AccountParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	AccountParameterTimestampOutputFormat                    AccountParameter = "TIMESTAMP_OUTPUT_FORMAT"
	AccountParameterTimestampTypeMapping                     AccountParameter = "TIMESTAMP_TYPE_MAPPING"
	AccountParameterTimestampTzOutputFormat                  AccountParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	AccountParameterTimezone                                 AccountParameter = "TIMEZONE"
	AccountParameterTransactionAbortOnError                  AccountParameter = "TRANSACTION_ABORT_ON_ERROR"
	AccountParameterTransactionDefaultIsolationLevel         AccountParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	AccountParameterTwoDigitCenturyStart                     AccountParameter = "TWO_DIGIT_CENTURY_START"
	AccountParameterUnsupportedDdlAction                     AccountParameter = "UNSUPPORTED_DDL_ACTION"
	AccountParameterUseCachedResult                          AccountParameter = "USE_CACHED_RESULT"
	AccountParameterWeekOfYearPolicy                         AccountParameter = "WEEK_OF_YEAR_POLICY"
	AccountParameterWeekStart                                AccountParameter = "WEEK_START"

	// Object Parameters (inherited)
	AccountParameterCatalog                                 AccountParameter = "CATALOG"
	AccountParameterDataRetentionTimeInDays                 AccountParameter = "DATA_RETENTION_TIME_IN_DAYS"
	AccountParameterDefaultDDLCollation                     AccountParameter = "DEFAULT_DDL_COLLATION"
	AccountParameterExternalVolume                          AccountParameter = "EXTERNAL_VOLUME"
	AccountParameterLogLevel                                AccountParameter = "LOG_LEVEL"
	AccountParameterMaxConcurrencyLevel                     AccountParameter = "MAX_CONCURRENCY_LEVEL"
	AccountParameterMaxDataExtensionTimeInDays              AccountParameter = "MAX_DATA_EXTENSION_TIME_IN_DAYS"
	AccountParameterPipeExecutionPaused                     AccountParameter = "PIPE_EXECUTION_PAUSED"
	AccountParameterReplaceInvalidCharacters                AccountParameter = "REPLACE_INVALID_CHARACTERS"
	AccountParameterStatementQueuedTimeoutInSeconds         AccountParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	AccountParameterStatementTimeoutInSeconds               AccountParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	AccountParameterStorageSerializationPolicy              AccountParameter = "STORAGE_SERIALIZATION_POLICY"
	AccountParameterShareRestrictions                       AccountParameter = "SHARE_RESTRICTIONS"
	AccountParameterSuspendTaskAfterNumFailures             AccountParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	AccountParameterTraceLevel                              AccountParameter = "TRACE_LEVEL"
	AccountParameterUserTaskManagedInitialWarehouseSize     AccountParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	AccountParameterUserTaskTimeoutMs                       AccountParameter = "USER_TASK_TIMEOUT_MS"
	AccountParameterTaskAutoRetryAttempts                   AccountParameter = "TASK_AUTO_RETRY_ATTEMPTS"
	AccountParameterUserTaskMinimumTriggerIntervalInSeconds AccountParameter = "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"
	AccountParameterMetricLevel                             AccountParameter = "METRIC_LEVEL"
	AccountParameterEnableConsoleOutput                     AccountParameter = "ENABLE_CONSOLE_OUTPUT"

	// User Parameters (inherited)
	AccountParameterEnablePersonalDatabase           AccountParameter = "ENABLE_PERSONAL_DATABASE"
	AccountParameterEnableUnredactedQuerySyntaxError AccountParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
	AccountParameterPreventUnloadToInternalStages    AccountParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES"
)

var AllAccountParameters = []AccountParameter{
	AccountParameterAllowClientMFACaching,
	AccountParameterAllowIDToken,
	AccountParameterClientEncryptionKeySize,
	AccountParameterCortexEnabledCrossRegion,
	AccountParameterEnableIdentifierFirstLogin,
	AccountParameterEnableInternalStagesPrivatelink,
	AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository,
	AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage,
	AccountParameterEnableUnhandledExceptionsReporting,
	AccountParameterEnforceNetworkRulesForInternalStages,
	AccountParameterEventTable,
	AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList,
	AccountParameterInitialReplicationSizeLimitInTB,
	AccountParameterMinDataRetentionTimeInDays,
	AccountParameterNetworkPolicy,
	AccountParameterOAuthAddPrivilegedRolesToBlockedList,
	AccountParameterPeriodicDataRekeying,
	AccountParameterPreventLoadFromInlineURL,
	AccountParameterPreventUnloadToInlineURL,
	AccountParameterRequireStorageIntegrationForStageCreation,
	AccountParameterRequireStorageIntegrationForStageOperation,
	AccountParameterSSOLoginPage,

	AccountParameterAbortDetachedQuery,
	AccountParameterActivePythonProfiler,
	AccountParameterAutocommit,
	AccountParameterBinaryInputFormat,
	AccountParameterBinaryOutputFormat,
	AccountParameterClientEnableLogInfoStatementParameters,
	AccountParameterClientMemoryLimit,
	AccountParameterClientMetadataRequestUseConnectionCtx,
	AccountParameterClientMetadataUseSessionDatabase,
	AccountParameterClientPrefetchThreads,
	AccountParameterClientResultChunkSize,
	AccountParameterClientSessionKeepAlive,
	AccountParameterClientSessionKeepAliveHeartbeatFrequency,
	AccountParameterClientTimestampTypeMapping,
	AccountParameterEnableUnloadPhysicalTypeOptimization,
	AccountParameterClientResultColumnCaseInsensitive,
	AccountParameterCsvTimestampFormat,
	AccountParameterDateInputFormat,
	AccountParameterDateOutputFormat,
	AccountParameterErrorOnNondeterministicMerge,
	AccountParameterErrorOnNondeterministicUpdate,
	AccountParameterGeographyOutputFormat,
	AccountParameterGeometryOutputFormat,
	AccountParameterHybridTableLockTimeout,
	AccountParameterJdbcTreatDecimalAsInt,
	AccountParameterJdbcTreatTimestampNtzAsUtc,
	AccountParameterJdbcUseSessionTimezone,
	AccountParameterJsonIndent,
	AccountParameterJsTreatIntegerAsBigInt,
	AccountParameterLockTimeout,
	AccountParameterMultiStatementCount,
	AccountParameterNoorderSequenceAsDefault,
	AccountParameterOdbcTreatDecimalAsInt,
	AccountParameterPythonProfilerModules,
	AccountParameterPythonProfilerTargetStage,
	AccountParameterQueryTag,
	AccountParameterQuotedIdentifiersIgnoreCase,
	AccountParameterRowsPerResultset,
	AccountParameterS3StageVpceDnsName,
	AccountParameterSearchPath,
	AccountParameterSimulatedDataSharingConsumer,
	AccountParameterStatementTimeoutInSeconds,
	AccountParameterStrictJsonOutput,
	AccountParameterTimeInputFormat,
	AccountParameterTimeOutputFormat,
	AccountParameterTimestampDayIsAlways24h,
	AccountParameterTimestampInputFormat,
	AccountParameterTimestampLtzOutputFormat,
	AccountParameterTimestampNtzOutputFormat,
	AccountParameterTimestampOutputFormat,
	AccountParameterTimestampTypeMapping,
	AccountParameterTimestampTzOutputFormat,
	AccountParameterTimezone,
	AccountParameterTransactionAbortOnError,
	AccountParameterTransactionDefaultIsolationLevel,
	AccountParameterTwoDigitCenturyStart,
	AccountParameterUnsupportedDdlAction,
	AccountParameterUseCachedResult,
	AccountParameterWeekOfYearPolicy,
	AccountParameterWeekStart,

	AccountParameterCatalog,
	AccountParameterDataRetentionTimeInDays,
	AccountParameterDefaultDDLCollation,
	AccountParameterExternalVolume,
	AccountParameterLogLevel,
	AccountParameterMaxConcurrencyLevel,
	AccountParameterMaxDataExtensionTimeInDays,
	AccountParameterPipeExecutionPaused,
	AccountParameterPreventUnloadToInternalStages,
	AccountParameterReplaceInvalidCharacters,
	AccountParameterStatementQueuedTimeoutInSeconds,
	AccountParameterStorageSerializationPolicy,
	AccountParameterShareRestrictions,
	AccountParameterSuspendTaskAfterNumFailures,
	AccountParameterTraceLevel,
	AccountParameterUserTaskManagedInitialWarehouseSize,
	AccountParameterUserTaskTimeoutMs,
	AccountParameterTaskAutoRetryAttempts,
	AccountParameterUserTaskMinimumTriggerIntervalInSeconds,
	AccountParameterMetricLevel,
	AccountParameterEnableConsoleOutput,
	AccountParameterEnableUnredactedQuerySyntaxError,
	AccountParameterEnablePersonalDatabase,
}

func ToAccountParameter(s string) (AccountParameter, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllAccountParameters, AccountParameter(s)) {
		return "", fmt.Errorf("invalid account parameter: %s", s)
	}
	return AccountParameter(s), nil
}

type SessionParameter string

const (
	SessionParameterAbortDetachedQuery                       SessionParameter = "ABORT_DETACHED_QUERY"
	SessionParameterActivePythonProfiler                     SessionParameter = "ACTIVE_PYTHON_PROFILER"
	SessionParameterAutocommit                               SessionParameter = "AUTOCOMMIT"
	SessionParameterBinaryInputFormat                        SessionParameter = "BINARY_INPUT_FORMAT"
	SessionParameterBinaryOutputFormat                       SessionParameter = "BINARY_OUTPUT_FORMAT"
	SessionParameterClientEnableLogInfoStatementParameters   SessionParameter = "CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS"
	SessionParameterClientMemoryLimit                        SessionParameter = "CLIENT_MEMORY_LIMIT"
	SessionParameterClientMetadataRequestUseConnectionCtx    SessionParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	SessionParameterClientPrefetchThreads                    SessionParameter = "CLIENT_PREFETCH_THREADS"
	SessionParameterClientResultChunkSize                    SessionParameter = "CLIENT_RESULT_CHUNK_SIZE"
	SessionParameterClientResultColumnCaseInsensitive        SessionParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	SessionParameterClientMetadataUseSessionDatabase         SessionParameter = "CLIENT_METADATA_USE_SESSION_DATABASE"
	SessionParameterClientSessionKeepAlive                   SessionParameter = "CLIENT_SESSION_KEEP_ALIVE"
	SessionParameterClientSessionKeepAliveHeartbeatFrequency SessionParameter = "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"
	SessionParameterClientTimestampTypeMapping               SessionParameter = "CLIENT_TIMESTAMP_TYPE_MAPPING"
	SessionParameterCsvTimestampFormat                       SessionParameter = "CSV_TIMESTAMP_FORMAT"
	SessionParameterDateInputFormat                          SessionParameter = "DATE_INPUT_FORMAT"
	SessionParameterDateOutputFormat                         SessionParameter = "DATE_OUTPUT_FORMAT"
	SessionParameterEnableUnloadPhysicalTypeOptimization     SessionParameter = "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"
	SessionParameterErrorOnNondeterministicMerge             SessionParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	SessionParameterErrorOnNondeterministicUpdate            SessionParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	SessionParameterGeographyOutputFormat                    SessionParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	SessionParameterGeometryOutputFormat                     SessionParameter = "GEOMETRY_OUTPUT_FORMAT"
	SessionParameterHybridTableLockTimeout                   SessionParameter = "HYBRID_TABLE_LOCK_TIMEOUT"
	SessionParameterJdbcTreatDecimalAsInt                    SessionParameter = "JDBC_TREAT_DECIMAL_AS_INT"
	SessionParameterJdbcTreatTimestampNtzAsUtc               SessionParameter = "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"
	SessionParameterJdbcUseSessionTimezone                   SessionParameter = "JDBC_USE_SESSION_TIMEZONE"
	SessionParameterJsonIndent                               SessionParameter = "JSON_INDENT"
	SessionParameterJsTreatIntegerAsBigInt                   SessionParameter = "JS_TREAT_INTEGER_AS_BIGINT"
	SessionParameterLockTimeout                              SessionParameter = "LOCK_TIMEOUT"
	SessionParameterLogLevel                                 SessionParameter = "LOG_LEVEL"
	SessionParameterMultiStatementCount                      SessionParameter = "MULTI_STATEMENT_COUNT"
	SessionParameterNoorderSequenceAsDefault                 SessionParameter = "NOORDER_SEQUENCE_AS_DEFAULT"
	SessionParameterOdbcTreatDecimalAsInt                    SessionParameter = "ODBC_TREAT_DECIMAL_AS_INT"
	SessionParameterPythonProfilerModules                    SessionParameter = "PYTHON_PROFILER_MODULES"
	SessionParameterPythonProfilerTargetStage                SessionParameter = "PYTHON_PROFILER_TARGET_STAGE"
	SessionParameterQueryTag                                 SessionParameter = "QUERY_TAG"
	SessionParameterQuotedIdentifiersIgnoreCase              SessionParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	SessionParameterRowsPerResultset                         SessionParameter = "ROWS_PER_RESULTSET"
	SessionParameterS3StageVpceDnsName                       SessionParameter = "S3_STAGE_VPCE_DNS_NAME"
	SessionParameterSearchPath                               SessionParameter = "SEARCH_PATH"
	SessionParameterSimulatedDataSharingConsumer             SessionParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	SessionParameterStatementQueuedTimeoutInSeconds          SessionParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	SessionParameterStatementTimeoutInSeconds                SessionParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	SessionParameterStrictJsonOutput                         SessionParameter = "STRICT_JSON_OUTPUT"
	SessionParameterTimestampDayIsAlways24h                  SessionParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	SessionParameterTimestampInputFormat                     SessionParameter = "TIMESTAMP_INPUT_FORMAT"
	SessionParameterTimestampLTZOutputFormat                 SessionParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	SessionParameterTimestampNTZOutputFormat                 SessionParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	SessionParameterTimestampOutputFormat                    SessionParameter = "TIMESTAMP_OUTPUT_FORMAT"
	SessionParameterTimestampTypeMapping                     SessionParameter = "TIMESTAMP_TYPE_MAPPING"
	SessionParameterTimestampTZOutputFormat                  SessionParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	SessionParameterTimezone                                 SessionParameter = "TIMEZONE"
	SessionParameterTimeInputFormat                          SessionParameter = "TIME_INPUT_FORMAT"
	SessionParameterTimeOutputFormat                         SessionParameter = "TIME_OUTPUT_FORMAT"
	SessionParameterTraceLevel                               SessionParameter = "TRACE_LEVEL"
	SessionParameterTransactionAbortOnError                  SessionParameter = "TRANSACTION_ABORT_ON_ERROR"
	SessionParameterTransactionDefaultIsolationLevel         SessionParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	SessionParameterTwoDigitCenturyStart                     SessionParameter = "TWO_DIGIT_CENTURY_START"
	SessionParameterUnsupportedDDLAction                     SessionParameter = "UNSUPPORTED_DDL_ACTION"
	SessionParameterUseCachedResult                          SessionParameter = "USE_CACHED_RESULT"
	SessionParameterWeekOfYearPolicy                         SessionParameter = "WEEK_OF_YEAR_POLICY"
	SessionParameterWeekStart                                SessionParameter = "WEEK_START"
)

type ObjectParameter string

const (
	// Object Parameters
	ObjectParameterDataRetentionTimeInDays                 ObjectParameter = "DATA_RETENTION_TIME_IN_DAYS"
	ObjectParameterDefaultDDLCollation                     ObjectParameter = "DEFAULT_DDL_COLLATION"
	ObjectParameterLogLevel                                ObjectParameter = "LOG_LEVEL"
	ObjectParameterMaxConcurrencyLevel                     ObjectParameter = "MAX_CONCURRENCY_LEVEL"
	ObjectParameterMaxDataExtensionTimeInDays              ObjectParameter = "MAX_DATA_EXTENSION_TIME_IN_DAYS"
	ObjectParameterPipeExecutionPaused                     ObjectParameter = "PIPE_EXECUTION_PAUSED"
	ObjectParameterPreventUnloadToInternalStages           ObjectParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES" // also an account param
	ObjectParameterStatementQueuedTimeoutInSeconds         ObjectParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	ObjectParameterStatementTimeoutInSeconds               ObjectParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	ObjectParameterNetworkPolicy                           ObjectParameter = "NETWORK_POLICY" // also an account param
	ObjectParameterShareRestrictions                       ObjectParameter = "SHARE_RESTRICTIONS"
	ObjectParameterSuspendTaskAfterNumFailures             ObjectParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	ObjectParameterTraceLevel                              ObjectParameter = "TRACE_LEVEL"
	ObjectParameterUserTaskManagedInitialWarehouseSize     ObjectParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	ObjectParameterUserTaskTimeoutMs                       ObjectParameter = "USER_TASK_TIMEOUT_MS"
	ObjectParameterCatalog                                 ObjectParameter = "CATALOG"
	ObjectParameterExternalVolume                          ObjectParameter = "EXTERNAL_VOLUME"
	ObjectParameterReplaceInvalidCharacters                ObjectParameter = "REPLACE_INVALID_CHARACTERS"
	ObjectParameterStorageSerializationPolicy              ObjectParameter = "STORAGE_SERIALIZATION_POLICY"
	ObjectParameterTaskAutoRetryAttempts                   ObjectParameter = "TASK_AUTO_RETRY_ATTEMPTS"
	ObjectParameterUserTaskMinimumTriggerIntervalInSeconds ObjectParameter = "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"
	ObjectParameterQuotedIdentifiersIgnoreCase             ObjectParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	ObjectParameterMetricLevel                             ObjectParameter = "METRIC_LEVEL"
	ObjectParameterEnableConsoleOutput                     ObjectParameter = "ENABLE_CONSOLE_OUTPUT"

	// User Parameters
	ObjectParameterEnableUnredactedQuerySyntaxError ObjectParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
)

type UserParameter string

const (
	// User Parameters
	UserParameterEnableUnredactedQuerySyntaxError UserParameter = "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"
	UserParameterNetworkPolicy                    UserParameter = "NETWORK_POLICY"
	UserParameterPreventUnloadToInternalStages    UserParameter = "PREVENT_UNLOAD_TO_INTERNAL_STAGES"

	// Session Parameters (inherited)
	UserParameterAbortDetachedQuery                       UserParameter = "ABORT_DETACHED_QUERY"
	UserParameterAutocommit                               UserParameter = "AUTOCOMMIT"
	UserParameterBinaryInputFormat                        UserParameter = "BINARY_INPUT_FORMAT"
	UserParameterBinaryOutputFormat                       UserParameter = "BINARY_OUTPUT_FORMAT"
	UserParameterClientMemoryLimit                        UserParameter = "CLIENT_MEMORY_LIMIT"
	UserParameterClientMetadataRequestUseConnectionCtx    UserParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	UserParameterClientPrefetchThreads                    UserParameter = "CLIENT_PREFETCH_THREADS"
	UserParameterClientResultChunkSize                    UserParameter = "CLIENT_RESULT_CHUNK_SIZE"
	UserParameterClientResultColumnCaseInsensitive        UserParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	UserParameterClientSessionKeepAlive                   UserParameter = "CLIENT_SESSION_KEEP_ALIVE"
	UserParameterClientSessionKeepAliveHeartbeatFrequency UserParameter = "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"
	UserParameterClientTimestampTypeMapping               UserParameter = "CLIENT_TIMESTAMP_TYPE_MAPPING"
	UserParameterDateInputFormat                          UserParameter = "DATE_INPUT_FORMAT"
	UserParameterDateOutputFormat                         UserParameter = "DATE_OUTPUT_FORMAT"
	UserParameterEnableUnloadPhysicalTypeOptimization     UserParameter = "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"
	UserParameterErrorOnNondeterministicMerge             UserParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	UserParameterErrorOnNondeterministicUpdate            UserParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	UserParameterGeographyOutputFormat                    UserParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	UserParameterGeometryOutputFormat                     UserParameter = "GEOMETRY_OUTPUT_FORMAT"
	UserParameterJdbcTreatDecimalAsInt                    UserParameter = "JDBC_TREAT_DECIMAL_AS_INT"
	UserParameterJdbcTreatTimestampNtzAsUtc               UserParameter = "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"
	UserParameterJdbcUseSessionTimezone                   UserParameter = "JDBC_USE_SESSION_TIMEZONE"
	UserParameterJsonIndent                               UserParameter = "JSON_INDENT"
	UserParameterLockTimeout                              UserParameter = "LOCK_TIMEOUT"
	UserParameterLogLevel                                 UserParameter = "LOG_LEVEL"
	UserParameterMultiStatementCount                      UserParameter = "MULTI_STATEMENT_COUNT"
	UserParameterNoorderSequenceAsDefault                 UserParameter = "NOORDER_SEQUENCE_AS_DEFAULT"
	UserParameterOdbcTreatDecimalAsInt                    UserParameter = "ODBC_TREAT_DECIMAL_AS_INT"
	UserParameterQueryTag                                 UserParameter = "QUERY_TAG"
	UserParameterQuotedIdentifiersIgnoreCase              UserParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	UserParameterRowsPerResultset                         UserParameter = "ROWS_PER_RESULTSET"
	UserParameterS3StageVpceDnsName                       UserParameter = "S3_STAGE_VPCE_DNS_NAME"
	UserParameterSearchPath                               UserParameter = "SEARCH_PATH"
	UserParameterSimulatedDataSharingConsumer             UserParameter = "SIMULATED_DATA_SHARING_CONSUMER"
	UserParameterStatementQueuedTimeoutInSeconds          UserParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	UserParameterStatementTimeoutInSeconds                UserParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	UserParameterStrictJsonOutput                         UserParameter = "STRICT_JSON_OUTPUT"
	UserParameterTimestampDayIsAlways24h                  UserParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	UserParameterTimestampInputFormat                     UserParameter = "TIMESTAMP_INPUT_FORMAT"
	UserParameterTimestampLtzOutputFormat                 UserParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	UserParameterTimestampNtzOutputFormat                 UserParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	UserParameterTimestampOutputFormat                    UserParameter = "TIMESTAMP_OUTPUT_FORMAT"
	UserParameterTimestampTypeMapping                     UserParameter = "TIMESTAMP_TYPE_MAPPING"
	UserParameterTimestampTzOutputFormat                  UserParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	UserParameterTimezone                                 UserParameter = "TIMEZONE"
	UserParameterTimeInputFormat                          UserParameter = "TIME_INPUT_FORMAT"
	UserParameterTimeOutputFormat                         UserParameter = "TIME_OUTPUT_FORMAT"
	UserParameterTraceLevel                               UserParameter = "TRACE_LEVEL"
	UserParameterTransactionAbortOnError                  UserParameter = "TRANSACTION_ABORT_ON_ERROR"
	UserParameterTransactionDefaultIsolationLevel         UserParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	UserParameterTwoDigitCenturyStart                     UserParameter = "TWO_DIGIT_CENTURY_START"
	UserParameterUnsupportedDdlAction                     UserParameter = "UNSUPPORTED_DDL_ACTION"
	UserParameterUseCachedResult                          UserParameter = "USE_CACHED_RESULT"
	UserParameterWeekOfYearPolicy                         UserParameter = "WEEK_OF_YEAR_POLICY"
	UserParameterWeekStart                                UserParameter = "WEEK_START"
)

var AllUserParameters = []UserParameter{
	UserParameterAbortDetachedQuery,
	UserParameterAutocommit,
	UserParameterBinaryInputFormat,
	UserParameterBinaryOutputFormat,
	UserParameterClientMemoryLimit,
	UserParameterClientMetadataRequestUseConnectionCtx,
	UserParameterClientPrefetchThreads,
	UserParameterClientResultChunkSize,
	UserParameterClientResultColumnCaseInsensitive,
	UserParameterClientSessionKeepAlive,
	UserParameterClientSessionKeepAliveHeartbeatFrequency,
	UserParameterClientTimestampTypeMapping,
	UserParameterDateInputFormat,
	UserParameterDateOutputFormat,
	UserParameterEnableUnloadPhysicalTypeOptimization,
	UserParameterErrorOnNondeterministicMerge,
	UserParameterErrorOnNondeterministicUpdate,
	UserParameterGeographyOutputFormat,
	UserParameterGeometryOutputFormat,
	UserParameterJdbcTreatDecimalAsInt,
	UserParameterJdbcTreatTimestampNtzAsUtc,
	UserParameterJdbcUseSessionTimezone,
	UserParameterJsonIndent,
	UserParameterLockTimeout,
	UserParameterLogLevel,
	UserParameterMultiStatementCount,
	UserParameterNoorderSequenceAsDefault,
	UserParameterOdbcTreatDecimalAsInt,
	UserParameterQueryTag,
	UserParameterQuotedIdentifiersIgnoreCase,
	UserParameterRowsPerResultset,
	UserParameterS3StageVpceDnsName,
	UserParameterSearchPath,
	UserParameterSimulatedDataSharingConsumer,
	UserParameterStatementQueuedTimeoutInSeconds,
	UserParameterStatementTimeoutInSeconds,
	UserParameterStrictJsonOutput,
	UserParameterTimestampDayIsAlways24h,
	UserParameterTimestampInputFormat,
	UserParameterTimestampLtzOutputFormat,
	UserParameterTimestampNtzOutputFormat,
	UserParameterTimestampOutputFormat,
	UserParameterTimestampTypeMapping,
	UserParameterTimestampTzOutputFormat,
	UserParameterTimezone,
	UserParameterTimeInputFormat,
	UserParameterTimeOutputFormat,
	UserParameterTraceLevel,
	UserParameterTransactionAbortOnError,
	UserParameterTransactionDefaultIsolationLevel,
	UserParameterTwoDigitCenturyStart,
	UserParameterUnsupportedDdlAction,
	UserParameterUseCachedResult,
	UserParameterWeekOfYearPolicy,
	UserParameterWeekStart,
	UserParameterEnableUnredactedQuerySyntaxError,
	UserParameterNetworkPolicy,
	UserParameterPreventUnloadToInternalStages,
}

type TaskParameter string

const (
	// Task Parameters
	TaskParameterSuspendTaskAfterNumFailures             TaskParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	TaskParameterTaskAutoRetryAttempts                   TaskParameter = "TASK_AUTO_RETRY_ATTEMPTS"
	TaskParameterUserTaskManagedInitialWarehouseSize     TaskParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	TaskParameterUserTaskMinimumTriggerIntervalInSeconds TaskParameter = "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"
	TaskParameterUserTaskTimeoutMs                       TaskParameter = "USER_TASK_TIMEOUT_MS"

	// Session Parameters (inherited)
	TaskParameterAbortDetachedQuery                       TaskParameter = "ABORT_DETACHED_QUERY"
	TaskParameterAutocommit                               TaskParameter = "AUTOCOMMIT"
	TaskParameterBinaryInputFormat                        TaskParameter = "BINARY_INPUT_FORMAT"
	TaskParameterBinaryOutputFormat                       TaskParameter = "BINARY_OUTPUT_FORMAT"
	TaskParameterClientMemoryLimit                        TaskParameter = "CLIENT_MEMORY_LIMIT"
	TaskParameterClientMetadataRequestUseConnectionCtx    TaskParameter = "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"
	TaskParameterClientPrefetchThreads                    TaskParameter = "CLIENT_PREFETCH_THREADS"
	TaskParameterClientResultChunkSize                    TaskParameter = "CLIENT_RESULT_CHUNK_SIZE"
	TaskParameterClientResultColumnCaseInsensitive        TaskParameter = "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"
	TaskParameterClientSessionKeepAlive                   TaskParameter = "CLIENT_SESSION_KEEP_ALIVE"
	TaskParameterClientSessionKeepAliveHeartbeatFrequency TaskParameter = "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"
	TaskParameterClientTimestampTypeMapping               TaskParameter = "CLIENT_TIMESTAMP_TYPE_MAPPING"
	TaskParameterDateInputFormat                          TaskParameter = "DATE_INPUT_FORMAT"
	TaskParameterDateOutputFormat                         TaskParameter = "DATE_OUTPUT_FORMAT"
	TaskParameterEnableUnloadPhysicalTypeOptimization     TaskParameter = "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"
	TaskParameterErrorOnNondeterministicMerge             TaskParameter = "ERROR_ON_NONDETERMINISTIC_MERGE"
	TaskParameterErrorOnNondeterministicUpdate            TaskParameter = "ERROR_ON_NONDETERMINISTIC_UPDATE"
	TaskParameterGeographyOutputFormat                    TaskParameter = "GEOGRAPHY_OUTPUT_FORMAT"
	TaskParameterGeometryOutputFormat                     TaskParameter = "GEOMETRY_OUTPUT_FORMAT"
	TaskParameterJdbcTreatTimestampNtzAsUtc               TaskParameter = "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"
	TaskParameterJdbcUseSessionTimezone                   TaskParameter = "JDBC_USE_SESSION_TIMEZONE"
	TaskParameterJsonIndent                               TaskParameter = "JSON_INDENT"
	TaskParameterLockTimeout                              TaskParameter = "LOCK_TIMEOUT"
	TaskParameterLogLevel                                 TaskParameter = "LOG_LEVEL"
	TaskParameterMultiStatementCount                      TaskParameter = "MULTI_STATEMENT_COUNT"
	TaskParameterNoorderSequenceAsDefault                 TaskParameter = "NOORDER_SEQUENCE_AS_DEFAULT"
	TaskParameterOdbcTreatDecimalAsInt                    TaskParameter = "ODBC_TREAT_DECIMAL_AS_INT"
	TaskParameterQueryTag                                 TaskParameter = "QUERY_TAG"
	TaskParameterQuotedIdentifiersIgnoreCase              TaskParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	TaskParameterRowsPerResultset                         TaskParameter = "ROWS_PER_RESULTSET"
	TaskParameterS3StageVpceDnsName                       TaskParameter = "S3_STAGE_VPCE_DNS_NAME"
	TaskParameterSearchPath                               TaskParameter = "SEARCH_PATH"
	TaskParameterStatementQueuedTimeoutInSeconds          TaskParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	TaskParameterStatementTimeoutInSeconds                TaskParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
	TaskParameterStrictJsonOutput                         TaskParameter = "STRICT_JSON_OUTPUT"
	TaskParameterTimestampDayIsAlways24h                  TaskParameter = "TIMESTAMP_DAY_IS_ALWAYS_24H"
	TaskParameterTimestampInputFormat                     TaskParameter = "TIMESTAMP_INPUT_FORMAT"
	TaskParameterTimestampLtzOutputFormat                 TaskParameter = "TIMESTAMP_LTZ_OUTPUT_FORMAT"
	TaskParameterTimestampNtzOutputFormat                 TaskParameter = "TIMESTAMP_NTZ_OUTPUT_FORMAT"
	TaskParameterTimestampOutputFormat                    TaskParameter = "TIMESTAMP_OUTPUT_FORMAT"
	TaskParameterTimestampTypeMapping                     TaskParameter = "TIMESTAMP_TYPE_MAPPING"
	TaskParameterTimestampTzOutputFormat                  TaskParameter = "TIMESTAMP_TZ_OUTPUT_FORMAT"
	TaskParameterTimezone                                 TaskParameter = "TIMEZONE"
	TaskParameterTimeInputFormat                          TaskParameter = "TIME_INPUT_FORMAT"
	TaskParameterTimeOutputFormat                         TaskParameter = "TIME_OUTPUT_FORMAT"
	TaskParameterTraceLevel                               TaskParameter = "TRACE_LEVEL"
	TaskParameterTransactionAbortOnError                  TaskParameter = "TRANSACTION_ABORT_ON_ERROR"
	TaskParameterTransactionDefaultIsolationLevel         TaskParameter = "TRANSACTION_DEFAULT_ISOLATION_LEVEL"
	TaskParameterTwoDigitCenturyStart                     TaskParameter = "TWO_DIGIT_CENTURY_START"
	TaskParameterUnsupportedDdlAction                     TaskParameter = "UNSUPPORTED_DDL_ACTION"
	TaskParameterUseCachedResult                          TaskParameter = "USE_CACHED_RESULT"
	TaskParameterWeekOfYearPolicy                         TaskParameter = "WEEK_OF_YEAR_POLICY"
	TaskParameterWeekStart                                TaskParameter = "WEEK_START"
)

var AllTaskParameters = []TaskParameter{
	// Task Parameters
	TaskParameterSuspendTaskAfterNumFailures,
	TaskParameterTaskAutoRetryAttempts,
	TaskParameterUserTaskManagedInitialWarehouseSize,
	TaskParameterUserTaskMinimumTriggerIntervalInSeconds,
	TaskParameterUserTaskTimeoutMs,

	// Session Parameters (inherited)
	TaskParameterAbortDetachedQuery,
	TaskParameterAutocommit,
	TaskParameterBinaryInputFormat,
	TaskParameterBinaryOutputFormat,
	TaskParameterClientMemoryLimit,
	TaskParameterClientMetadataRequestUseConnectionCtx,
	TaskParameterClientPrefetchThreads,
	TaskParameterClientResultChunkSize,
	TaskParameterClientResultColumnCaseInsensitive,
	TaskParameterClientSessionKeepAlive,
	TaskParameterClientSessionKeepAliveHeartbeatFrequency,
	TaskParameterClientTimestampTypeMapping,
	TaskParameterDateInputFormat,
	TaskParameterDateOutputFormat,
	TaskParameterEnableUnloadPhysicalTypeOptimization,
	TaskParameterErrorOnNondeterministicMerge,
	TaskParameterErrorOnNondeterministicUpdate,
	TaskParameterGeographyOutputFormat,
	TaskParameterGeometryOutputFormat,
	TaskParameterJdbcTreatTimestampNtzAsUtc,
	TaskParameterJdbcUseSessionTimezone,
	TaskParameterJsonIndent,
	TaskParameterLockTimeout,
	TaskParameterLogLevel,
	TaskParameterMultiStatementCount,
	TaskParameterNoorderSequenceAsDefault,
	TaskParameterOdbcTreatDecimalAsInt,
	TaskParameterQueryTag,
	TaskParameterQuotedIdentifiersIgnoreCase,
	TaskParameterRowsPerResultset,
	TaskParameterS3StageVpceDnsName,
	TaskParameterSearchPath,
	TaskParameterStatementQueuedTimeoutInSeconds,
	TaskParameterStatementTimeoutInSeconds,
	TaskParameterStrictJsonOutput,
	TaskParameterTimestampDayIsAlways24h,
	TaskParameterTimestampInputFormat,
	TaskParameterTimestampLtzOutputFormat,
	TaskParameterTimestampNtzOutputFormat,
	TaskParameterTimestampOutputFormat,
	TaskParameterTimestampTypeMapping,
	TaskParameterTimestampTzOutputFormat,
	TaskParameterTimezone,
	TaskParameterTimeInputFormat,
	TaskParameterTimeOutputFormat,
	TaskParameterTraceLevel,
	TaskParameterTransactionAbortOnError,
	TaskParameterTransactionDefaultIsolationLevel,
	TaskParameterTwoDigitCenturyStart,
	TaskParameterUnsupportedDdlAction,
	TaskParameterUseCachedResult,
	TaskParameterWeekOfYearPolicy,
	TaskParameterWeekStart,
}

type WarehouseParameter string

const (
	WarehouseParameterMaxConcurrencyLevel             WarehouseParameter = "MAX_CONCURRENCY_LEVEL"
	WarehouseParameterStatementQueuedTimeoutInSeconds WarehouseParameter = "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"
	WarehouseParameterStatementTimeoutInSeconds       WarehouseParameter = "STATEMENT_TIMEOUT_IN_SECONDS"
)

var AllSchemaParameters = []ObjectParameter{
	ObjectParameterDataRetentionTimeInDays,
	ObjectParameterMaxDataExtensionTimeInDays,
	ObjectParameterExternalVolume,
	ObjectParameterCatalog,
	ObjectParameterReplaceInvalidCharacters,
	ObjectParameterDefaultDDLCollation,
	ObjectParameterStorageSerializationPolicy,
	ObjectParameterLogLevel,
	ObjectParameterTraceLevel,
	ObjectParameterSuspendTaskAfterNumFailures,
	ObjectParameterTaskAutoRetryAttempts,
	ObjectParameterUserTaskManagedInitialWarehouseSize,
	ObjectParameterUserTaskTimeoutMs,
	ObjectParameterUserTaskMinimumTriggerIntervalInSeconds,
	ObjectParameterQuotedIdentifiersIgnoreCase,
	ObjectParameterEnableConsoleOutput,
	ObjectParameterPipeExecutionPaused,
}

type DatabaseParameter string

const (
	DatabaseParameterDataRetentionTimeInDays                 DatabaseParameter = "DATA_RETENTION_TIME_IN_DAYS"
	DatabaseParameterMaxDataExtensionTimeInDays              DatabaseParameter = "MAX_DATA_EXTENSION_TIME_IN_DAYS"
	DatabaseParameterExternalVolume                          DatabaseParameter = "EXTERNAL_VOLUME"
	DatabaseParameterCatalog                                 DatabaseParameter = "CATALOG"
	DatabaseParameterReplaceInvalidCharacters                DatabaseParameter = "REPLACE_INVALID_CHARACTERS"
	DatabaseParameterDefaultDdlCollation                     DatabaseParameter = "DEFAULT_DDL_COLLATION"
	DatabaseParameterStorageSerializationPolicy              DatabaseParameter = "STORAGE_SERIALIZATION_POLICY"
	DatabaseParameterLogLevel                                DatabaseParameter = "LOG_LEVEL"
	DatabaseParameterTraceLevel                              DatabaseParameter = "TRACE_LEVEL"
	DatabaseParameterSuspendTaskAfterNumFailures             DatabaseParameter = "SUSPEND_TASK_AFTER_NUM_FAILURES"
	DatabaseParameterTaskAutoRetryAttempts                   DatabaseParameter = "TASK_AUTO_RETRY_ATTEMPTS"
	DatabaseParameterUserTaskManagedInitialWarehouseSize     DatabaseParameter = "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"
	DatabaseParameterUserTaskTimeoutMs                       DatabaseParameter = "USER_TASK_TIMEOUT_MS"
	DatabaseParameterUserTaskMinimumTriggerIntervalInSeconds DatabaseParameter = "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"
	DatabaseParameterQuotedIdentifiersIgnoreCase             DatabaseParameter = "QUOTED_IDENTIFIERS_IGNORE_CASE"
	DatabaseParameterEnableConsoleOutput                     DatabaseParameter = "ENABLE_CONSOLE_OUTPUT"
)

type FunctionParameter string

const (
	FunctionParameterEnableConsoleOutput FunctionParameter = "ENABLE_CONSOLE_OUTPUT"
	FunctionParameterLogLevel            FunctionParameter = "LOG_LEVEL"
	FunctionParameterMetricLevel         FunctionParameter = "METRIC_LEVEL"
	FunctionParameterTraceLevel          FunctionParameter = "TRACE_LEVEL"
)

var AllFunctionParameters = []FunctionParameter{
	FunctionParameterEnableConsoleOutput,
	FunctionParameterLogLevel,
	FunctionParameterMetricLevel,
	FunctionParameterTraceLevel,
}

type ProcedureParameter string

const (
	ProcedureParameterAutoEventLogging    ProcedureParameter = "AUTO_EVENT_LOGGING"
	ProcedureParameterEnableConsoleOutput ProcedureParameter = "ENABLE_CONSOLE_OUTPUT"
	ProcedureParameterLogLevel            ProcedureParameter = "LOG_LEVEL"
	ProcedureParameterMetricLevel         ProcedureParameter = "METRIC_LEVEL"
	ProcedureParameterTraceLevel          ProcedureParameter = "TRACE_LEVEL"
)

var AllProcedureParameters = []ProcedureParameter{
	ProcedureParameterEnableConsoleOutput,
	ProcedureParameterLogLevel,
	ProcedureParameterMetricLevel,
	ProcedureParameterTraceLevel,
}

// AccountParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#account-parameters.
type AccountParameters struct {
	// Account Parameters
	AllowClientMFACaching                            *bool   `ddl:"parameter" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                     *bool   `ddl:"parameter" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                          *int    `ddl:"parameter" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	CortexEnabledCrossRegion                         *string `ddl:"parameter,single_quotes" sql:"CORTEX_ENABLED_CROSS_REGION"`
	EnableIdentifierFirstLogin                       *bool   `ddl:"parameter" sql:"ENABLE_IDENTIFIER_FIRST_LOGIN"`
	EnableInternalStagesPrivatelink                  *bool   `ddl:"parameter" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EnablePersonalDatabase                           *bool   `ddl:"parameter" sql:"ENABLE_PERSONAL_DATABASE"`
	EnableUnredactedQuerySyntaxError                 *bool   `ddl:"parameter" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	EnableTriSecretAndRekeyOptOutForImageRepository  *bool   `ddl:"parameter" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"`
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage *bool   `ddl:"parameter" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE"`
	EnableUnhandledExceptionsReporting               *bool   `ddl:"parameter" sql:"ENABLE_UNHANDLED_EXCEPTIONS_REPORTING"`
	EnforceNetworkRulesForInternalStages             *bool   `ddl:"keyword" sql:"ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES"`
	EventTable                                       *string `ddl:"parameter,single_quotes" sql:"EVENT_TABLE"`
	ExternalOAuthAddPrivilegedRolesToBlockedList     *bool   `ddl:"parameter" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	// InitialReplicationSizeLimitInTB is a string because values like 3.0 get rounded to 3, resulting in an error in Snowflake.
	// This is still validated below.
	InitialReplicationSizeLimitInTB            *string `ddl:"parameter" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                 *int    `ddl:"parameter" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                              *string `ddl:"parameter,single_quotes" sql:"NETWORK_POLICY"`
	OAuthAddPrivilegedRolesToBlockedList       *bool   `ddl:"parameter" sql:"OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	PeriodicDataRekeying                       *bool   `ddl:"parameter" sql:"PERIODIC_DATA_REKEYING"`
	PreventLoadFromInlineURL                   *bool   `ddl:"parameter" sql:"PREVENT_LOAD_FROM_INLINE_URL"`
	PreventUnloadToInlineURL                   *bool   `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages              *bool   `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation  *bool   `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation *bool   `ddl:"parameter" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                               *bool   `ddl:"parameter" sql:"SSO_LOGIN_PAGE"`
}

func (v *AccountParameters) validate() error {
	var errs []error
	if valueSet(v.ClientEncryptionKeySize) {
		if !slices.Contains([]int{128, 256}, *v.ClientEncryptionKeySize) {
			errs = append(errs, fmt.Errorf("ClientEncryptionKeySize must be either 128 or 256, got %d", *v.ClientEncryptionKeySize))
		}
	}
	if valueSet(v.InitialReplicationSizeLimitInTB) {
		value, err := strconv.ParseFloat(*v.InitialReplicationSizeLimitInTB, 64)
		if err != nil || value < 0 {
			return fmt.Errorf("InitialReplicationSizeLimitInTB must be a non-negative float, got %v", *v.InitialReplicationSizeLimitInTB)
		}

	}
	if valueSet(v.MinDataRetentionTimeInDays) {
		if !validateIntInRangeInclusive(*v.MinDataRetentionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("AccountParameters", "MinDataRetentionTimeInDays", 0, 90))
		}
	}
	return errors.Join(errs...)
}

type AccountParametersUnset struct {
	AllowClientMFACaching                            *bool `ddl:"keyword" sql:"ALLOW_CLIENT_MFA_CACHING"`
	AllowIDToken                                     *bool `ddl:"keyword" sql:"ALLOW_ID_TOKEN"`
	ClientEncryptionKeySize                          *bool `ddl:"keyword" sql:"CLIENT_ENCRYPTION_KEY_SIZE"`
	CortexEnabledCrossRegion                         *bool `ddl:"keyword" sql:"CORTEX_ENABLED_CROSS_REGION"`
	EnableIdentifierFirstLogin                       *bool `ddl:"keyword" sql:"ENABLE_IDENTIFIER_FIRST_LOGIN"`
	EnableInternalStagesPrivatelink                  *bool `ddl:"keyword" sql:"ENABLE_INTERNAL_STAGES_PRIVATELINK"`
	EnablePersonalDatabase                           *bool `ddl:"keyword" sql:"ENABLE_PERSONAL_DATABASE"`
	EnableTriSecretAndRekeyOptOutForImageRepository  *bool `ddl:"keyword" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY"`
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage *bool `ddl:"keyword" sql:"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE"`
	EnableUnhandledExceptionsReporting               *bool `ddl:"keyword" sql:"ENABLE_UNHANDLED_EXCEPTIONS_REPORTING"`
	EventTable                                       *bool `ddl:"keyword" sql:"EVENT_TABLE"`
	EnableUnredactedQuerySyntaxError                 *bool `ddl:"keyword" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	EnforceNetworkRulesForInternalStages             *bool `ddl:"keyword" sql:"ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES"`
	ExternalOAuthAddPrivilegedRolesToBlockedList     *bool `ddl:"keyword" sql:"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	InitialReplicationSizeLimitInTB                  *bool `ddl:"keyword" sql:"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB"`
	MinDataRetentionTimeInDays                       *bool `ddl:"keyword" sql:"MIN_DATA_RETENTION_TIME_IN_DAYS"`
	NetworkPolicy                                    *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	OAuthAddPrivilegedRolesToBlockedList             *bool `ddl:"keyword" sql:"OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"`
	PeriodicDataRekeying                             *bool `ddl:"keyword" sql:"PERIODIC_DATA_REKEYING"`
	PreventLoadFromInlineURL                         *bool `ddl:"keyword" sql:"PREVENT_LOAD_FROM_INLINE_URL"`
	PreventUnloadToInlineURL                         *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INLINE_URL"`
	PreventUnloadToInternalStages                    *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	RequireStorageIntegrationForStageCreation        *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"`
	RequireStorageIntegrationForStageOperation       *bool `ddl:"keyword" sql:"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION"`
	SSOLoginPage                                     *bool `ddl:"keyword" sql:"SSO_LOGIN_PAGE"`
}

type ActivePythonProfiler string

const (
	ActivePythonProfilerLine   ActivePythonProfiler = "LINE"
	ActivePythonProfilerMemory ActivePythonProfiler = "MEMORY"
)

func ToActivePythonProfiler(s string) (ActivePythonProfiler, error) {
	switch strings.ToUpper(s) {
	case string(ActivePythonProfilerLine):
		return ActivePythonProfilerLine, nil
	case string(ActivePythonProfilerMemory):
		return ActivePythonProfilerMemory, nil
	default:
		return "", fmt.Errorf("invalid active python profiler: %s", s)
	}
}

type GeographyOutputFormat string

const (
	GeographyOutputFormatGeoJSON GeographyOutputFormat = "GeoJSON"
	GeographyOutputFormatWKT     GeographyOutputFormat = "WKT"
	GeographyOutputFormatWKB     GeographyOutputFormat = "WKB"
	GeographyOutputFormatEWKT    GeographyOutputFormat = "EWKT"
	GeographyOutputFormatEWKB    GeographyOutputFormat = "EWKB"
)

func ToGeographyOutputFormat(s string) (GeographyOutputFormat, error) {
	switch strings.ToUpper(s) {
	case strings.ToUpper(string(GeographyOutputFormatGeoJSON)):
		return GeographyOutputFormatGeoJSON, nil
	case string(GeographyOutputFormatWKT):
		return GeographyOutputFormatWKT, nil
	case string(GeographyOutputFormatWKB):
		return GeographyOutputFormatWKB, nil
	case string(GeographyOutputFormatEWKT):
		return GeographyOutputFormatEWKT, nil
	case string(GeographyOutputFormatEWKB):
		return GeographyOutputFormatEWKB, nil
	default:
		return "", fmt.Errorf("invalid geography output format: %s", s)
	}
}

type GeometryOutputFormat string

const (
	GeometryOutputFormatGeoJSON GeometryOutputFormat = "GeoJSON"
	GeometryOutputFormatWKT     GeometryOutputFormat = "WKT"
	GeometryOutputFormatWKB     GeometryOutputFormat = "WKB"
	GeometryOutputFormatEWKT    GeometryOutputFormat = "EWKT"
	GeometryOutputFormatEWKB    GeometryOutputFormat = "EWKB"
)

func ToGeometryOutputFormat(s string) (GeometryOutputFormat, error) {
	switch strings.ToUpper(s) {
	case strings.ToUpper(string(GeometryOutputFormatGeoJSON)):
		return GeometryOutputFormatGeoJSON, nil
	case string(GeometryOutputFormatWKT):
		return GeometryOutputFormatWKT, nil
	case string(GeometryOutputFormatWKB):
		return GeometryOutputFormatWKB, nil
	case string(GeometryOutputFormatEWKT):
		return GeometryOutputFormatEWKT, nil
	case string(GeometryOutputFormatEWKB):
		return GeometryOutputFormatEWKB, nil
	default:
		return "", fmt.Errorf("invalid geometry output format: %s", s)
	}
}

type BinaryInputFormat string

const (
	BinaryInputFormatHex    BinaryInputFormat = "HEX"
	BinaryInputFormatBase64 BinaryInputFormat = "BASE64"
	BinaryInputFormatUTF8   BinaryInputFormat = "UTF8"
)

func ToBinaryInputFormat(s string) (BinaryInputFormat, error) {
	switch strings.ToUpper(s) {
	case string(BinaryInputFormatHex):
		return BinaryInputFormatHex, nil
	case string(BinaryInputFormatBase64):
		return BinaryInputFormatBase64, nil
	case string(BinaryInputFormatUTF8), "UTF-8":
		return BinaryInputFormatUTF8, nil
	default:
		return "", fmt.Errorf("invalid binary input format: %s", s)
	}
}

type BinaryOutputFormat string

const (
	BinaryOutputFormatHex    BinaryOutputFormat = "HEX"
	BinaryOutputFormatBase64 BinaryOutputFormat = "BASE64"
)

func ToBinaryOutputFormat(s string) (BinaryOutputFormat, error) {
	switch strings.ToUpper(s) {
	case string(BinaryOutputFormatHex):
		return BinaryOutputFormatHex, nil
	case string(BinaryOutputFormatBase64):
		return BinaryOutputFormatBase64, nil
	default:
		return "", fmt.Errorf("invalid binary output format: %s", s)
	}
}

type ClientTimestampTypeMapping string

const (
	ClientTimestampTypeMappingLtz ClientTimestampTypeMapping = "TIMESTAMP_LTZ"
	ClientTimestampTypeMappingNtz ClientTimestampTypeMapping = "TIMESTAMP_NTZ"
)

func ToClientTimestampTypeMapping(s string) (ClientTimestampTypeMapping, error) {
	switch strings.ToUpper(s) {
	case string(ClientTimestampTypeMappingLtz):
		return ClientTimestampTypeMappingLtz, nil
	case string(ClientTimestampTypeMappingNtz):
		return ClientTimestampTypeMappingNtz, nil
	default:
		return "", fmt.Errorf("invalid client timestamp type mapping: %s", s)
	}
}

type TimestampTypeMapping string

const (
	TimestampTypeMappingLtz TimestampTypeMapping = "TIMESTAMP_LTZ"
	TimestampTypeMappingNtz TimestampTypeMapping = "TIMESTAMP_NTZ"
	TimestampTypeMappingTz  TimestampTypeMapping = "TIMESTAMP_TZ"
)

func ToTimestampTypeMapping(s string) (TimestampTypeMapping, error) {
	switch strings.ToUpper(s) {
	case string(TimestampTypeMappingLtz):
		return TimestampTypeMappingLtz, nil
	case string(TimestampTypeMappingNtz):
		return TimestampTypeMappingNtz, nil
	case string(TimestampTypeMappingTz):
		return TimestampTypeMappingTz, nil
	default:
		return "", fmt.Errorf("invalid timestamp type mapping: %s", s)
	}
}

type TransactionDefaultIsolationLevel string

const (
	TransactionDefaultIsolationLevelReadCommitted TransactionDefaultIsolationLevel = "READ COMMITTED"
)

func ToTransactionDefaultIsolationLevel(s string) (TransactionDefaultIsolationLevel, error) {
	switch strings.ToUpper(s) {
	case string(TransactionDefaultIsolationLevelReadCommitted):
		return TransactionDefaultIsolationLevelReadCommitted, nil
	default:
		return "", fmt.Errorf("invalid transaction default isolation level: %s", s)
	}
}

type UnsupportedDDLAction string

const (
	UnsupportedDDLActionIgnore UnsupportedDDLAction = "IGNORE"
	UnsupportedDDLActionFail   UnsupportedDDLAction = "FAIL"
)

func ToUnsupportedDDLAction(s string) (UnsupportedDDLAction, error) {
	switch strings.ToUpper(s) {
	case string(UnsupportedDDLActionIgnore):
		return UnsupportedDDLActionIgnore, nil
	case string(UnsupportedDDLActionFail):
		return UnsupportedDDLActionFail, nil
	default:
		return "", fmt.Errorf("invalid ddl action: %s", s)
	}
}

// SessionParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#session-parameters.
type SessionParameters struct {
	AbortDetachedQuery                       *bool                             `ddl:"parameter" sql:"ABORT_DETACHED_QUERY"`
	ActivePythonProfiler                     *ActivePythonProfiler             `ddl:"parameter,single_quotes" sql:"ACTIVE_PYTHON_PROFILER"`
	Autocommit                               *bool                             `ddl:"parameter" sql:"AUTOCOMMIT"`
	BinaryInputFormat                        *BinaryInputFormat                `ddl:"parameter,single_quotes" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat                       *BinaryOutputFormat               `ddl:"parameter,single_quotes" sql:"BINARY_OUTPUT_FORMAT"`
	ClientEnableLogInfoStatementParameters   *bool                             `ddl:"parameter" sql:"CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS"`
	ClientMemoryLimit                        *int                              `ddl:"parameter" sql:"CLIENT_MEMORY_LIMIT"`
	ClientMetadataRequestUseConnectionCtx    *bool                             `ddl:"parameter" sql:"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"`
	ClientPrefetchThreads                    *int                              `ddl:"parameter" sql:"CLIENT_PREFETCH_THREADS"`
	ClientResultChunkSize                    *int                              `ddl:"parameter" sql:"CLIENT_RESULT_CHUNK_SIZE"`
	ClientResultColumnCaseInsensitive        *bool                             `ddl:"parameter" sql:"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"`
	ClientMetadataUseSessionDatabase         *bool                             `ddl:"parameter" sql:"CLIENT_METADATA_USE_SESSION_DATABASE"`
	ClientSessionKeepAlive                   *bool                             `ddl:"parameter" sql:"CLIENT_SESSION_KEEP_ALIVE"`
	ClientSessionKeepAliveHeartbeatFrequency *int                              `ddl:"parameter" sql:"CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"`
	ClientTimestampTypeMapping               *ClientTimestampTypeMapping       `ddl:"parameter,single_quotes" sql:"CLIENT_TIMESTAMP_TYPE_MAPPING"`
	CsvTimestampFormat                       *string                           `ddl:"parameter" sql:"CSV_TIMESTAMP_FORMAT"`
	DateInputFormat                          *string                           `ddl:"parameter,single_quotes" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                         *string                           `ddl:"parameter,single_quotes" sql:"DATE_OUTPUT_FORMAT"`
	EnableUnloadPhysicalTypeOptimization     *bool                             `ddl:"parameter" sql:"ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"`
	ErrorOnNondeterministicMerge             *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate            *bool                             `ddl:"parameter" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat                    *GeographyOutputFormat            `ddl:"parameter,single_quotes" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	GeometryOutputFormat                     *GeometryOutputFormat             `ddl:"parameter,single_quotes" sql:"GEOMETRY_OUTPUT_FORMAT"`
	HybridTableLockTimeout                   *int                              `ddl:"parameter" sql:"HYBRID_TABLE_LOCK_TIMEOUT"`
	JdbcTreatDecimalAsInt                    *bool                             `ddl:"parameter" sql:"JDBC_TREAT_DECIMAL_AS_INT"`
	JdbcTreatTimestampNtzAsUtc               *bool                             `ddl:"parameter" sql:"JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"`
	JdbcUseSessionTimezone                   *bool                             `ddl:"parameter" sql:"JDBC_USE_SESSION_TIMEZONE"`
	JsonIndent                               *int                              `ddl:"parameter" sql:"JSON_INDENT"`
	JsTreatIntegerAsBigInt                   *bool                             `ddl:"parameter" sql:"JS_TREAT_INTEGER_AS_BIGINT"`
	LockTimeout                              *int                              `ddl:"parameter" sql:"LOCK_TIMEOUT"`
	LogLevel                                 *LogLevel                         `ddl:"parameter" sql:"LOG_LEVEL"`
	MultiStatementCount                      *int                              `ddl:"parameter" sql:"MULTI_STATEMENT_COUNT"`
	NoorderSequenceAsDefault                 *bool                             `ddl:"parameter" sql:"NOORDER_SEQUENCE_AS_DEFAULT"`
	OdbcTreatDecimalAsInt                    *bool                             `ddl:"parameter" sql:"ODBC_TREAT_DECIMAL_AS_INT"`
	PythonProfilerModules                    *string                           `ddl:"parameter" sql:"PYTHON_PROFILER_MODULES"`
	PythonProfilerTargetStage                *string                           `ddl:"parameter" sql:"PYTHON_PROFILER_TARGET_STAGE"`
	QueryTag                                 *string                           `ddl:"parameter,single_quotes" sql:"QUERY_TAG"`
	QuotedIdentifiersIgnoreCase              *bool                             `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	RowsPerResultset                         *int                              `ddl:"parameter" sql:"ROWS_PER_RESULTSET"`
	S3StageVpceDnsName                       *string                           `ddl:"parameter,single_quotes" sql:"S3_STAGE_VPCE_DNS_NAME"`
	SearchPath                               *string                           `ddl:"parameter,single_quotes" sql:"SEARCH_PATH"`
	SimulatedDataSharingConsumer             *string                           `ddl:"parameter,single_quotes" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementQueuedTimeoutInSeconds          *int                              `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds                *int                              `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJsonOutput                         *bool                             `ddl:"parameter" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h                  *bool                             `ddl:"parameter" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat                     *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat                 *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat                 *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat                    *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping                     *TimestampTypeMapping             `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat                  *string                           `ddl:"parameter,single_quotes" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                                 *string                           `ddl:"parameter,single_quotes" sql:"TIMEZONE"`
	TimeInputFormat                          *string                           `ddl:"parameter,single_quotes" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                         *string                           `ddl:"parameter,single_quotes" sql:"TIME_OUTPUT_FORMAT"`
	TraceLevel                               *TraceLevel                       `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	TransactionAbortOnError                  *bool                             `ddl:"parameter" sql:"TRANSACTION_ABORT_ON_ERROR"`
	TransactionDefaultIsolationLevel         *TransactionDefaultIsolationLevel `ddl:"parameter,single_quotes" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart                     *int                              `ddl:"parameter" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction                     *UnsupportedDDLAction             `ddl:"parameter,single_quotes" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                          *bool                             `ddl:"parameter" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                         *int                              `ddl:"parameter" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                                *int                              `ddl:"parameter" sql:"WEEK_START"`
}

func (v *SessionParameters) validate() error {
	var errs []error
	// Do not validate input and output formats, because there are a lot of them, and may be not supported in Go itself.
	// See https://docs.snowflake.com/en/sql-reference/date-time-input-output#supported-formats-for-auto-detection.
	if valueSet(v.ClientPrefetchThreads) {
		if !validateIntGreaterThanOrEqual(*v.ClientPrefetchThreads, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientPrefetchThreads", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.ClientResultChunkSize) {
		if !validateIntGreaterThanOrEqual(*v.ClientResultChunkSize, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientResultChunkSize", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.ClientSessionKeepAliveHeartbeatFrequency) {
		if !validateIntGreaterThanOrEqual(*v.ClientSessionKeepAliveHeartbeatFrequency, 0) {
			errs = append(errs, errIntValue("SessionParameters", "ClientSessionKeepAliveHeartbeatFrequency", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.HybridTableLockTimeout) {
		if !validateIntGreaterThanOrEqual(*v.HybridTableLockTimeout, 0) {
			errs = append(errs, errIntValue("SessionParameters", "HybridTableLockTimeout", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.JsonIndent) {
		if !validateIntGreaterThanOrEqual(*v.JsonIndent, 0) {
			errs = append(errs, errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.LockTimeout) {
		if !validateIntGreaterThanOrEqual(*v.LockTimeout, 0) {
			errs = append(errs, errIntValue("SessionParameters", "LockTimeout", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.MultiStatementCount) {
		if !validateIntGreaterThanOrEqual(*v.MultiStatementCount, 0) {
			errs = append(errs, errIntValue("SessionParameters", "MultiStatementCount", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.RowsPerResultset) {
		if !validateIntGreaterThanOrEqual(*v.RowsPerResultset, 0) {
			errs = append(errs, errIntValue("SessionParameters", "RowsPerResultset", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.StatementQueuedTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementQueuedTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("SessionParameters", "StatementQueuedTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.StatementTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("SessionParameters", "StatementTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.TwoDigitCenturyStart) {
		if !validateIntGreaterThanOrEqual(*v.TwoDigitCenturyStart, 1900) {
			errs = append(errs, errIntValue("SessionParameters", "TwoDigitCenturyStart", IntErrGreaterOrEqual, 1900))
		}
	}
	if valueSet(v.WeekOfYearPolicy) {
		if !validateIntGreaterThanOrEqual(*v.WeekOfYearPolicy, 0) {
			errs = append(errs, errIntValue("SessionParameters", "WeekOfYearPolicy", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.WeekStart) {
		if !validateIntGreaterThanOrEqual(*v.WeekStart, 0) {
			errs = append(errs, errIntValue("SessionParameters", "WeekStart", IntErrGreaterOrEqual, 0))
		}
	}
	return errors.Join(errs...)
}

type SessionParametersUnset struct {
	AbortDetachedQuery                       *bool `ddl:"keyword" sql:"ABORT_DETACHED_QUERY"`
	ActivePythonProfiler                     *bool `ddl:"keyword" sql:"ACTIVE_PYTHON_PROFILER"`
	Autocommit                               *bool `ddl:"keyword" sql:"AUTOCOMMIT"`
	BinaryInputFormat                        *bool `ddl:"keyword" sql:"BINARY_INPUT_FORMAT"`
	BinaryOutputFormat                       *bool `ddl:"keyword" sql:"BINARY_OUTPUT_FORMAT"`
	ClientEnableLogInfoStatementParameters   *bool `ddl:"keyword" sql:"CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS"`
	ClientMemoryLimit                        *bool `ddl:"keyword" sql:"CLIENT_MEMORY_LIMIT"`
	ClientMetadataRequestUseConnectionCtx    *bool `ddl:"keyword" sql:"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX"`
	ClientPrefetchThreads                    *bool `ddl:"keyword" sql:"CLIENT_PREFETCH_THREADS"`
	ClientResultChunkSize                    *bool `ddl:"keyword" sql:"CLIENT_RESULT_CHUNK_SIZE"`
	ClientResultColumnCaseInsensitive        *bool `ddl:"keyword" sql:"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE"`
	ClientMetadataUseSessionDatabase         *bool `ddl:"keyword" sql:"CLIENT_METADATA_USE_SESSION_DATABASE"`
	ClientSessionKeepAlive                   *bool `ddl:"keyword" sql:"CLIENT_SESSION_KEEP_ALIVE"`
	ClientSessionKeepAliveHeartbeatFrequency *bool `ddl:"keyword" sql:"CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY"`
	ClientTimestampTypeMapping               *bool `ddl:"keyword" sql:"CLIENT_TIMESTAMP_TYPE_MAPPING"`
	CsvTimestampFormat                       *bool `ddl:"keyword" sql:"CSV_TIMESTAMP_FORMAT"`
	DateInputFormat                          *bool `ddl:"keyword" sql:"DATE_INPUT_FORMAT"`
	DateOutputFormat                         *bool `ddl:"keyword" sql:"DATE_OUTPUT_FORMAT"`
	EnableUnloadPhysicalTypeOptimization     *bool `ddl:"keyword" sql:"ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION"`
	ErrorOnNondeterministicMerge             *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_MERGE"`
	ErrorOnNondeterministicUpdate            *bool `ddl:"keyword" sql:"ERROR_ON_NONDETERMINISTIC_UPDATE"`
	GeographyOutputFormat                    *bool `ddl:"keyword" sql:"GEOGRAPHY_OUTPUT_FORMAT"`
	GeometryOutputFormat                     *bool `ddl:"keyword" sql:"GEOMETRY_OUTPUT_FORMAT"`
	HybridTableLockTimeout                   *bool `ddl:"keyword" sql:"HYBRID_TABLE_LOCK_TIMEOUT"`
	JdbcTreatDecimalAsInt                    *bool `ddl:"keyword" sql:"JDBC_TREAT_DECIMAL_AS_INT"`
	JdbcTreatTimestampNtzAsUtc               *bool `ddl:"keyword" sql:"JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC"`
	JdbcUseSessionTimezone                   *bool `ddl:"keyword" sql:"JDBC_USE_SESSION_TIMEZONE"`
	JsonIndent                               *bool `ddl:"keyword" sql:"JSON_INDENT"`
	JsTreatIntegerAsBigInt                   *bool `ddl:"keyword" sql:"JS_TREAT_INTEGER_AS_BIGINT"`
	LockTimeout                              *bool `ddl:"keyword" sql:"LOCK_TIMEOUT"`
	LogLevel                                 *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	MultiStatementCount                      *bool `ddl:"keyword" sql:"MULTI_STATEMENT_COUNT"`
	NoorderSequenceAsDefault                 *bool `ddl:"keyword" sql:"NOORDER_SEQUENCE_AS_DEFAULT"`
	OdbcTreatDecimalAsInt                    *bool `ddl:"keyword" sql:"ODBC_TREAT_DECIMAL_AS_INT"`
	PythonProfilerModules                    *bool `ddl:"keyword" sql:"PYTHON_PROFILER_MODULES"`
	PythonProfilerTargetStage                *bool `ddl:"keyword" sql:"PYTHON_PROFILER_TARGET_STAGE"`
	QueryTag                                 *bool `ddl:"keyword" sql:"QUERY_TAG"`
	QuotedIdentifiersIgnoreCase              *bool `ddl:"keyword" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	RowsPerResultset                         *bool `ddl:"keyword" sql:"ROWS_PER_RESULTSET"`
	S3StageVpceDnsName                       *bool `ddl:"keyword" sql:"S3_STAGE_VPCE_DNS_NAME"`
	SearchPath                               *bool `ddl:"keyword" sql:"SEARCH_PATH"`
	SimulatedDataSharingConsumer             *bool `ddl:"keyword" sql:"SIMULATED_DATA_SHARING_CONSUMER"`
	StatementQueuedTimeoutInSeconds          *bool `ddl:"keyword" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds                *bool `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	StrictJsonOutput                         *bool `ddl:"keyword" sql:"STRICT_JSON_OUTPUT"`
	TimestampDayIsAlways24h                  *bool `ddl:"keyword" sql:"TIMESTAMP_DAY_IS_ALWAYS_24H"`
	TimestampInputFormat                     *bool `ddl:"keyword" sql:"TIMESTAMP_INPUT_FORMAT"`
	TimestampLTZOutputFormat                 *bool `ddl:"keyword" sql:"TIMESTAMP_LTZ_OUTPUT_FORMAT"`
	TimestampNTZOutputFormat                 *bool `ddl:"keyword" sql:"TIMESTAMP_NTZ_OUTPUT_FORMAT"`
	TimestampOutputFormat                    *bool `ddl:"keyword" sql:"TIMESTAMP_OUTPUT_FORMAT"`
	TimestampTypeMapping                     *bool `ddl:"keyword" sql:"TIMESTAMP_TYPE_MAPPING"`
	TimestampTZOutputFormat                  *bool `ddl:"keyword" sql:"TIMESTAMP_TZ_OUTPUT_FORMAT"`
	Timezone                                 *bool `ddl:"keyword" sql:"TIMEZONE"`
	TimeInputFormat                          *bool `ddl:"keyword" sql:"TIME_INPUT_FORMAT"`
	TimeOutputFormat                         *bool `ddl:"keyword" sql:"TIME_OUTPUT_FORMAT"`
	TraceLevel                               *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
	TransactionAbortOnError                  *bool `ddl:"keyword" sql:"TRANSACTION_ABORT_ON_ERROR"`
	TransactionDefaultIsolationLevel         *bool `ddl:"keyword" sql:"TRANSACTION_DEFAULT_ISOLATION_LEVEL"`
	TwoDigitCenturyStart                     *bool `ddl:"keyword" sql:"TWO_DIGIT_CENTURY_START"`
	UnsupportedDDLAction                     *bool `ddl:"keyword" sql:"UNSUPPORTED_DDL_ACTION"`
	UseCachedResult                          *bool `ddl:"keyword" sql:"USE_CACHED_RESULT"`
	WeekOfYearPolicy                         *bool `ddl:"keyword" sql:"WEEK_OF_YEAR_POLICY"`
	WeekStart                                *bool `ddl:"keyword" sql:"WEEK_START"`
}

func (v *SessionParametersUnset) validate() error {
	if !anyValueSet(v.AbortDetachedQuery, v.ActivePythonProfiler, v.Autocommit, v.BinaryInputFormat, v.BinaryOutputFormat, v.ClientEnableLogInfoStatementParameters, v.ClientMemoryLimit, v.ClientMetadataRequestUseConnectionCtx, v.ClientPrefetchThreads, v.ClientResultChunkSize, v.ClientResultColumnCaseInsensitive, v.ClientMetadataUseSessionDatabase, v.ClientSessionKeepAlive, v.ClientSessionKeepAliveHeartbeatFrequency, v.ClientTimestampTypeMapping, v.CsvTimestampFormat, v.DateInputFormat, v.DateOutputFormat, v.EnableUnloadPhysicalTypeOptimization, v.ErrorOnNondeterministicMerge, v.ErrorOnNondeterministicUpdate, v.GeographyOutputFormat, v.GeometryOutputFormat, v.HybridTableLockTimeout, v.JdbcTreatDecimalAsInt, v.JdbcTreatTimestampNtzAsUtc, v.JdbcUseSessionTimezone, v.JsonIndent, v.JsTreatIntegerAsBigInt, v.LockTimeout, v.LogLevel, v.MultiStatementCount, v.NoorderSequenceAsDefault, v.OdbcTreatDecimalAsInt, v.PythonProfilerModules, v.PythonProfilerTargetStage, v.QueryTag, v.QuotedIdentifiersIgnoreCase, v.RowsPerResultset, v.S3StageVpceDnsName, v.SearchPath, v.SimulatedDataSharingConsumer, v.StatementQueuedTimeoutInSeconds, v.StatementTimeoutInSeconds, v.StrictJsonOutput, v.TimestampDayIsAlways24h, v.TimestampInputFormat, v.TimestampLTZOutputFormat, v.TimestampNTZOutputFormat, v.TimestampOutputFormat, v.TimestampTypeMapping, v.TimestampTZOutputFormat, v.Timezone, v.TimeInputFormat, v.TimeOutputFormat, v.TraceLevel, v.TransactionAbortOnError, v.TransactionDefaultIsolationLevel, v.TwoDigitCenturyStart, v.UnsupportedDDLAction, v.UseCachedResult, v.WeekOfYearPolicy, v.WeekStart) {
		return errors.Join(errAtLeastOneOf("SessionParametersUnset", "AbortDetachedQuery", "ActivePythonProfiler", "Autocommit", "BinaryInputFormat", "BinaryOutputFormat", "ClientEnableLogInfoStatementParameters", "ClientMemoryLimit", "ClientMetadataRequestUseConnectionCtx", "ClientPrefetchThreads", "ClientResultChunkSize", "ClientResultColumnCaseInsensitive", "ClientMetadataUseSessionDatabase", "ClientSessionKeepAlive", "ClientSessionKeepAliveHeartbeatFrequency", "ClientTimestampTypeMapping", "CsvTimestampFormat", "DateInputFormat", "DateOutputFormat", "EnableUnloadPhysicalTypeOptimization", "ErrorOnNondeterministicMerge", "ErrorOnNondeterministicUpdate", "GeographyOutputFormat", "GeometryOutputFormat", "HybridTableLockTimeout", "JdbcTreatDecimalAsInt", "JdbcTreatTimestampNtzAsUtc", "JdbcUseSessionTimezone", "JsonIndent", "JsTreatIntegerAsBigInt", "LockTimeout", "LogLevel", "MultiStatementCount", "NoorderSequenceAsDefault", "OdbcTreatDecimalAsInt", "PythonProfilerModules", "PythonProfilerTargetStage", "QueryTag", "QuotedIdentifiersIgnoreCase", "RowsPerResultset", "S3StageVpceDnsName", "SearchPath", "SimulatedDataSharingConsumer", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "StrictJsonOutput", "TimestampDayIsAlways24h", "TimestampInputFormat", "TimestampLTZOutputFormat", "TimestampNTZOutputFormat", "TimestampOutputFormat", "TimestampTypeMapping", "TimestampTZOutputFormat", "Timezone", "TimeInputFormat", "TimeOutputFormat", "TraceLevel", "TransactionAbortOnError", "TransactionDefaultIsolationLevel", "TwoDigitCenturyStart", "UnsupportedDDLAction", "UseCachedResult", "WeekOfYearPolicy", "WeekStart"))
	}
	return nil
}

// ObjectParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#object-parameters.
type ObjectParameters struct {
	Catalog                                 *string        `ddl:"parameter" sql:"CATALOG"`
	DataRetentionTimeInDays                 *int           `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	DefaultDDLCollation                     *string        `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	EnableUnredactedQuerySyntaxError        *bool          `ddl:"parameter" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	LogLevel                                *LogLevel      `ddl:"parameter" sql:"LOG_LEVEL"`
	MaxConcurrencyLevel                     *int           `ddl:"parameter" sql:"MAX_CONCURRENCY_LEVEL"`
	MaxDataExtensionTimeInDays              *int           `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	PipeExecutionPaused                     *bool          `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	PreventUnloadToInternalStages           *bool          `ddl:"parameter" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	StatementQueuedTimeoutInSeconds         *int           `ddl:"parameter" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds               *int           `ddl:"parameter" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	NetworkPolicy                           *string        `ddl:"parameter,single_quotes" sql:"NETWORK_POLICY"`
	ShareRestrictions                       *bool          `ddl:"parameter" sql:"SHARE_RESTRICTIONS"`
	SuspendTaskAfterNumFailures             *int           `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	StorageSerializationPolicy              *string        `ddl:"parameter" sql:"STORAGE_SERIALIZATION_POLICY"`
	TraceLevel                              *TraceLevel    `ddl:"parameter" sql:"TRACE_LEVEL"`
	TaskAutoRetryAttempts                   *int           `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskMinimumTriggerIntervalInSeconds *int           `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	UserTaskTimeoutMs                       *int           `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
}

func (v *ObjectParameters) validate() error {
	var errs []error
	if valueSet(v.DataRetentionTimeInDays) {
		if !validateIntInRangeInclusive(*v.DataRetentionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("ObjectParameters", "DataRetentionTimeInDays", 0, 90))
		}
	}
	if valueSet(v.MaxConcurrencyLevel) {
		if !validateIntGreaterThanOrEqual(*v.MaxConcurrencyLevel, 1) {
			errs = append(errs, errIntValue("ObjectParameters", "MaxConcurrencyLevel", IntErrGreaterOrEqual, 1))
		}
	}
	if valueSet(v.MaxDataExtensionTimeInDays) {
		if !validateIntInRangeInclusive(*v.MaxDataExtensionTimeInDays, 0, 90) {
			errs = append(errs, errIntBetween("ObjectParameters", "MaxDataExtensionTimeInDays", 0, 90))
		}
	}
	if valueSet(v.StatementQueuedTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementQueuedTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "StatementQueuedTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.StatementTimeoutInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.StatementTimeoutInSeconds, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "StatementTimeoutInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.SuspendTaskAfterNumFailures) {
		if !validateIntGreaterThanOrEqual(*v.SuspendTaskAfterNumFailures, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "SuspendTaskAfterNumFailures", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.TaskAutoRetryAttempts) {
		if !validateIntGreaterThanOrEqual(*v.TaskAutoRetryAttempts, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "SuspendTaskAfterNumFailures", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.UserTaskMinimumTriggerIntervalInSeconds) {
		if !validateIntGreaterThanOrEqual(*v.UserTaskMinimumTriggerIntervalInSeconds, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "UserTaskMinimumTriggerIntervalInSeconds", IntErrGreaterOrEqual, 0))
		}
	}
	if valueSet(v.UserTaskTimeoutMs) {
		if !validateIntGreaterThanOrEqual(*v.UserTaskTimeoutMs, 0) {
			errs = append(errs, errIntValue("ObjectParameters", "UserTaskTimeoutMs", IntErrGreaterOrEqual, 0))
		}
	}
	return errors.Join(errs...)
}

type ObjectParametersUnset struct {
	Catalog                             *bool `ddl:"keyword" sql:"CATALOG"`
	DataRetentionTimeInDays             *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	DefaultDDLCollation                 *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	EnableUnredactedQuerySyntaxError    *bool `ddl:"keyword" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	LogLevel                            *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	MaxConcurrencyLevel                 *bool `ddl:"keyword" sql:"MAX_CONCURRENCY_LEVEL"`
	MaxDataExtensionTimeInDays          *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	PipeExecutionPaused                 *bool `ddl:"keyword" sql:"PIPE_EXECUTION_PAUSED"`
	PreventUnloadToInternalStages       *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
	StatementQueuedTimeoutInSeconds     *bool `ddl:"keyword" sql:"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS"`
	StatementTimeoutInSeconds           *bool `ddl:"keyword" sql:"STATEMENT_TIMEOUT_IN_SECONDS"`
	NetworkPolicy                       *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	ShareRestrictions                   *bool `ddl:"keyword" sql:"SHARE_RESTRICTIONS"`
	SuspendTaskAfterNumFailures         *bool `ddl:"keyword" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	StorageSerializationPolicy          *bool `ddl:"keyword" sql:"STORAGE_SERIALIZATION_POLICY"`
	TaskAutoRetryAttempts               *bool `ddl:"keyword" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
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
	Session   *bool                               `ddl:"keyword" sql:"SESSION"`
	Account   *bool                               `ddl:"keyword" sql:"ACCOUNT"`
	User      AccountObjectIdentifier             `ddl:"identifier" sql:"USER"`
	Warehouse AccountObjectIdentifier             `ddl:"identifier" sql:"WAREHOUSE"`
	Database  AccountObjectIdentifier             `ddl:"identifier" sql:"DATABASE"`
	Schema    DatabaseObjectIdentifier            `ddl:"identifier" sql:"SCHEMA"`
	Task      SchemaObjectIdentifier              `ddl:"identifier" sql:"TASK"`
	Table     SchemaObjectIdentifier              `ddl:"identifier" sql:"TABLE"`
	Function  SchemaObjectIdentifierWithArguments `ddl:"identifier" sql:"FUNCTION"`
	Procedure SchemaObjectIdentifierWithArguments `ddl:"identifier" sql:"PROCEDURE"`
}

func (v *ParametersIn) validate() error {
	if !anyValueSet(v.Session, v.Account, v.User, v.Warehouse, v.Database, v.Schema, v.Task, v.Table, v.Function, v.Procedure) {
		return errors.Join(errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table", "Function", "Procedure"))
	}
	return nil
}

type ParameterType string

const (
	ParameterTypeSnowflakeDefault ParameterType = ""
	ParameterTypeAccount          ParameterType = "ACCOUNT"
	ParameterTypeUser             ParameterType = "USER"
	ParameterTypeSession          ParameterType = "SESSION"
	ParameterTypeObject           ParameterType = "OBJECT"
	ParameterTypeWarehouse        ParameterType = "WAREHOUSE"
	ParameterTypeDatabase         ParameterType = "DATABASE"
	ParameterTypeSchema           ParameterType = "SCHEMA"
	ParameterTypeTask             ParameterType = "TASK"
	ParameterTypeFunction         ParameterType = "FUNCTION"
	ParameterTypeProcedure        ParameterType = "PROCEDURE"
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
	case ObjectTypeFunction:
		opts.In.Function = object.Name.(SchemaObjectIdentifierWithArguments)
	case ObjectTypeProcedure:
		opts.In.Procedure = object.Name.(SchemaObjectIdentifierWithArguments)
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
