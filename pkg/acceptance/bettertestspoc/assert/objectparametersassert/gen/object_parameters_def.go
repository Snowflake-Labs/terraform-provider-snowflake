package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type SnowflakeObjectParameters struct {
	Name              string
	IdType            string
	Level             sdk.ParameterType
	AdditionalImports []string
	Parameters        []SnowflakeParameter
}

func (p SnowflakeObjectParameters) ObjectName() string {
	return p.Name
}

type SnowflakeParameter struct {
	ParameterName string
	ParameterType string
	DefaultValue  string
	DefaultLevel  string
}

func GetAllSnowflakeObjectParameters() []SnowflakeObjectParameters {
	return allObjectsParameters
}

// TODO [SNOW-1501905]: use SDK definition after parameters rework (+ preprocessing here)
var allObjectsParameters = []SnowflakeObjectParameters{
	{
		Name:   "User",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeUser,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.UserParameterEnableUnredactedQuerySyntaxError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterNetworkPolicy), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterPreventUnloadToInternalStages), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterAbortDetachedQuery), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterAutocommit), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterBinaryInputFormat), ParameterType: "sdk.BinaryInputFormat", DefaultValue: "sdk.BinaryInputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterBinaryOutputFormat), ParameterType: "sdk.BinaryOutputFormat", DefaultValue: "sdk.BinaryOutputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientMemoryLimit), ParameterType: "int", DefaultValue: "1536", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientMetadataRequestUseConnectionCtx), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientPrefetchThreads), ParameterType: "int", DefaultValue: "4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientResultChunkSize), ParameterType: "int", DefaultValue: "160", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientResultColumnCaseInsensitive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientSessionKeepAlive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientTimestampTypeMapping), ParameterType: "sdk.ClientTimestampTypeMapping", DefaultValue: "sdk.ClientTimestampTypeMappingLtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterDateInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterDateOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterEnableUnloadPhysicalTypeOptimization), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterErrorOnNondeterministicMerge), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterErrorOnNondeterministicUpdate), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterGeographyOutputFormat), ParameterType: "sdk.GeographyOutputFormat", DefaultValue: "sdk.GeographyOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterGeometryOutputFormat), ParameterType: "sdk.GeometryOutputFormat", DefaultValue: "sdk.GeometryOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcTreatTimestampNtzAsUtc), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcUseSessionTimezone), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJsonIndent), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterLockTimeout), ParameterType: "int", DefaultValue: "43200", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterMultiStatementCount), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterNoorderSequenceAsDefault), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterOdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterQueryTag), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterRowsPerResultset), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterS3StageVpceDnsName), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterSimulatedDataSharingConsumer), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStrictJsonOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampDayIsAlways24h), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampLtzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampNtzOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampTypeMapping), ParameterType: "sdk.TimestampTypeMapping", DefaultValue: "sdk.TimestampTypeMappingNtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampTzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterTransactionDefaultIsolationLevel), ParameterType: "sdk.TransactionDefaultIsolationLevel", DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO [SNOW-1501905]: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
			{ParameterName: string(sdk.UserParameterUnsupportedDdlAction), ParameterType: "sdk.UnsupportedDDLAction", DefaultValue: "sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterUseCachedResult), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterWeekOfYearPolicy), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterWeekStart), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
		// for the quickfix above
		AdditionalImports: []string{"strings"},
	},
	{
		Name:   "Warehouse",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeWarehouse,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.WarehouseParameterMaxConcurrencyLevel), ParameterType: "int", DefaultValue: "8", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Database",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeDatabase,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.DatabaseParameterDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterMaxDataExtensionTimeInDays), ParameterType: "int", DefaultValue: "14", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterExternalVolume), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterCatalog), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterReplaceInvalidCharacters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterDefaultDdlCollation), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterStorageSerializationPolicy), ParameterType: "sdk.StorageSerializationPolicy", DefaultValue: "sdk.StorageSerializationPolicyOptimized", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterSuspendTaskAfterNumFailures), ParameterType: "int", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterTaskAutoRetryAttempts), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: "sdk.WarehouseSizeMedium", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.DatabaseParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Task",
		IdType: "sdk.SchemaObjectIdentifier",
		Level:  sdk.ParameterTypeTask,
		Parameters: []SnowflakeParameter{
			// ABORT_DETACHED_QUERY
			// ACTIVE_PYTHON_PROFILER
			{ParameterName: string(sdk.TaskParameterAutocommit), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeAccount"},
			// AUTOCOMMIT_API_SUPPORTED
			{ParameterName: string(sdk.TaskParameterBinaryInputFormat), ParameterType: "sdk.BinaryInputFormat", DefaultValue: "sdk.BinaryInputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterBinaryOutputFormat), ParameterType: "sdk.BinaryOutputFormat", DefaultValue: "sdk.BinaryOutputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// CLIENT_ENABLE_CONSERVATIVE_MEMORY_USAGE
			// CLIENT_ENABLE_DEFAULT_OVERWRITE_IN_PUT
			// CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS
			{ParameterName: string(sdk.TaskParameterClientMemoryLimit), ParameterType: "int", DefaultValue: "1536", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientMetadataRequestUseConnectionCtx), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// CLIENT_METADATA_USE_SESSION_DATABASE
			{ParameterName: string(sdk.TaskParameterClientPrefetchThreads), ParameterType: "int", DefaultValue: "4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientResultChunkSize), ParameterType: "int", DefaultValue: "160", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientResultColumnCaseInsensitive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// CLIENT_SESSION_CLONE
			{ParameterName: string(sdk.TaskParameterClientSessionKeepAlive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientTimestampTypeMapping), ParameterType: "sdk.ClientTimestampTypeMapping", DefaultValue: "sdk.ClientTimestampTypeMappingLtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// CSV_TIMESTAMP_FORMAT
			// C_API_QUERY_RESULT_FORMAT
			{ParameterName: string(sdk.TaskParameterDateInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterDateOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// DOTNET_QUERY_RESULT_FORMAT
			// DYNAMIC_TABLES_VIEW_VERSION
			// DYNAMIC_TABLE_GRAPH_HISTORY_VIEW_VERSION
			// DYNAMIC_TABLE_REFRESH_HISTORY_VIEW_VERSION
			// ENABLE_CONSOLE_OUTPUT
			// ENABLE_PROVIDER_LISTING_PROGRAMMATIC_ACCESS_DESCRIBE_LISTING
			// ENABLE_SNOW_API_FOR_COMPUTE_POOL
			// ENABLE_SNOW_API_FOR_COPILOT
			// ENABLE_SNOW_API_FOR_DATABASE
			// ENABLE_SNOW_API_FOR_FUNCTION
			// ENABLE_SNOW_API_FOR_GRANT
			// ENABLE_SNOW_API_FOR_ICEBERG
			// ENABLE_SNOW_API_FOR_IMAGE_REPOSITORY
			// ENABLE_SNOW_API_FOR_RESULT
			// ENABLE_SNOW_API_FOR_ROLE
			// ENABLE_SNOW_API_FOR_SCHEMA
			// ENABLE_SNOW_API_FOR_SERVICE
			// ENABLE_SNOW_API_FOR_SESSION
			// ENABLE_SNOW_API_FOR_STAGE
			// ENABLE_SNOW_API_FOR_TABLE
			// ENABLE_SNOW_API_FOR_TASK
			// ENABLE_SNOW_API_FOR_USER
			// ENABLE_SNOW_API_FOR_WAREHOUSE
			{ParameterName: string(sdk.TaskParameterEnableUnloadPhysicalTypeOptimization), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterErrorOnNondeterministicMerge), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterErrorOnNondeterministicUpdate), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterGeographyOutputFormat), ParameterType: "sdk.GeographyOutputFormat", DefaultValue: "sdk.GeographyOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterGeometryOutputFormat), ParameterType: "sdk.GeometryOutputFormat", DefaultValue: "sdk.GeometryOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// GO_QUERY_RESULT_FORMAT
			// HYBRID_TABLE_LOCK_TIMEOUT
			// INCLUDE_DT_WITH_TABLE_KIND_IN_SHOW_OBJECTS
			// INCLUDE_DYNAMIC_TABLES_WITH_TABLE_KIND
			// JDBC_FORMAT_DATE_WITH_TIMEZONE
			// JDBC_QUERY_RESULT_FORMAT
			{ParameterName: string(sdk.TaskParameterJdbcTreatTimestampNtzAsUtc), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterJdbcUseSessionTimezone), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterJsonIndent), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// LANGUAGE
			{ParameterName: string(sdk.TaskParameterLockTimeout), ParameterType: "int", DefaultValue: "43200", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// METRIC_LEVEL
			{ParameterName: string(sdk.TaskParameterMultiStatementCount), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterNoorderSequenceAsDefault), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// ODBC_QUERY_RESULT_FORMAT
			{ParameterName: string(sdk.TaskParameterOdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// PYTHON_CONNECTOR_QUERY_RESULT_FORMAT
			// PYTHON_CONNECTOR_USE_NANOARROW
			// PYTHON_PROFILER_MODULES
			// PYTHON_PROFILER_TARGET_STAGE
			// PYTHON_SNOWPARK_AUTO_CLEAN_UP_TEMP_TABLE_ENABLED
			// PYTHON_SNOWPARK_COMPILATION_STAGE_ENABLED
			// PYTHON_SNOWPARK_ELIMINATE_NUMERIC_SQL_VALUE_CAST_ENABLED
			// PYTHON_SNOWPARK_USE_CTE_OPTIMIZATION
			// PYTHON_SNOWPARK_USE_LARGE_QUERY_BREAKDOWN_OPTIMIZATION
			// PYTHON_SNOWPARK_USE_LOGICAL_TYPE_FOR_CREATE_DATAFRAME
			// PYTHON_SNOWPARK_USE_SCOPED_TEMP_OBJECTS
			// PYTHON_SNOWPARK_USE_SQL_SIMPLIFIER
			// QA_TEST_NAME
			// QUERY_RESULT_FORMAT
			{ParameterName: string(sdk.TaskParameterQueryTag), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.TaskParameterRowsPerResultset), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterS3StageVpceDnsName), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// SHOW_EXTERNAL_TABLE_KIND_AS_TABLE
			// SNOWPARK_HIDE_INTERNAL_ALIAS
			// SNOWPARK_LAZY_ANALYSIS
			// SNOWPARK_REQUEST_TIMEOUT_IN_SECONDS
			// SNOWPARK_STORED_PROC_IS_FINAL_TABLE_QUERY
			// SNOWPARK_USE_SCOPED_TEMP_OBJECTS
			// SQL_API_NULLABLE_IN_RESULT_SET
			// SQL_API_QUERY_RESULT_FORMAT
			{ParameterName: string(sdk.TaskParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterStrictJsonOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO: SUSPEND_TASK_AFTER_NUM_FAILURES
			// TODO: TASK_AUTO_RETRY_ATTEMPTS
			{ParameterName: string(sdk.TaskParameterTimestampDayIsAlways24h), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampLtzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampNtzOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampTypeMapping), ParameterType: "sdk.TimestampTypeMapping", DefaultValue: "sdk.TimestampTypeMappingNtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampTzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.TaskParameterTransactionDefaultIsolationLevel), ParameterType: "sdk.TransactionDefaultIsolationLevel", DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO UI_QUERY_RESULT_FORMAT
			// TODO [SNOW-1501905]: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
			{ParameterName: string(sdk.TaskParameterUnsupportedDdlAction), ParameterType: "sdk.UnsupportedDDLAction", DefaultValue: "sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: "Medium", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// {ParameterName: string(sdk.UserParameterUnsupportedUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// {ParameterName: string(sdk.UserParameterUnsupportedUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUseCachedResult), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterWeekOfYearPolicy), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterWeekStart), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
		// for the quickfix above
		AdditionalImports: []string{"strings"},
	},
}
