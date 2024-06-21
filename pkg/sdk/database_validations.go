package sdk

type ParameterWithType struct {
	Key  ObjectParameter
	Type ParameterType // TODO: Remove ?
}

// DatabaseParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#object-parameters
var DatabaseParameters = []ParameterWithType{
	{ObjectParameterDataRetentionTimeInDays, ParameterTypeNumber},
	{ObjectParameterMaxDataExtensionTimeInDays, ParameterTypeNumber},
	{ObjectParameterSuspendTaskAfterNumFailures, ParameterTypeNumber},
	{ObjectParameterTaskAutoRetryAttempts, ParameterTypeNumber},
	{ObjectParameterUserTaskTimeoutMs, ParameterTypeNumber},
	{ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, ParameterTypeNumber},

	{ObjectParameterReplaceInvalidCharacters, ParameterTypeBoolean},
	{ObjectParameterQuotedIdentifiersIgnoreCase, ParameterTypeBoolean},
	{ObjectParameterEnableConsoleOutput, ParameterTypeBoolean},

	{ObjectParameterExternalVolume, ParameterTypeString},
	{ObjectParameterCatalog, ParameterTypeString},
	{ObjectParameterDefaultDDLCollation, ParameterTypeString},
	{ObjectParameterStorageSerializationPolicy, ParameterTypeString},
	{ObjectParameterLogLevel, ParameterTypeString},
	{ObjectParameterTraceLevel, ParameterTypeString},
	{ObjectParameterUserTaskManagedInitialWarehouseSize, ParameterTypeString},
}
