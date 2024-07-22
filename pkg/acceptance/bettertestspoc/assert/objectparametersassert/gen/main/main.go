//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func main() {
	gencommons.NewGenerator(
		getAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getAllSnowflakeObjectParameters() []gen.SnowflakeObjectParameters {
	return allObjectsParameters
}

func getFilename(_ gen.SnowflakeObjectParameters, model gen.SnowflakeObjectParametersAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_parameters_snowflake" + "_gen.go"
}

// TODO: use SDK definition after parameters rework (+ preprocessing here)
var allObjectsParameters = []gen.SnowflakeObjectParameters{
	{
		Name:   "User",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeUser,
		Parameters: []gen.SnowflakeParameter{
			{string(sdk.UserParameterEnableUnredactedQuerySyntaxError), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterNetworkPolicy), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterPreventUnloadToInternalStages), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterAbortDetachedQuery), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterAutocommit), "bool", "true", "sdk.ParameterTypeAccount"},
			{string(sdk.UserParameterBinaryInputFormat), "sdk.BinaryInputFormat", "sdk.BinaryInputFormatHex", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterBinaryOutputFormat), "sdk.BinaryOutputFormat", "sdk.BinaryOutputFormatHex", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientMemoryLimit), "int", "1536", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientMetadataRequestUseConnectionCtx), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientPrefetchThreads), "int", "4", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientResultChunkSize), "int", "160", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientResultColumnCaseInsensitive), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientSessionKeepAlive), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency), "int", "3600", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterClientTimestampTypeMapping), "sdk.ClientTimestampTypeMapping", "sdk.ClientTimestampTypeMappingLtz", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterDateInputFormat), "string", "AUTO", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterDateOutputFormat), "string", "YYYY-MM-DD", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterEnableUnloadPhysicalTypeOptimization), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterErrorOnNondeterministicMerge), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterErrorOnNondeterministicUpdate), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterGeographyOutputFormat), "sdk.GeographyOutputFormat", "sdk.GeographyOutputFormatGeoJSON", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterGeometryOutputFormat), "sdk.GeometryOutputFormat", "sdk.GeometryOutputFormatGeoJSON", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterJdbcTreatDecimalAsInt), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterJdbcTreatTimestampNtzAsUtc), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterJdbcUseSessionTimezone), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterJsonIndent), "int", "2", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterLockTimeout), "int", "43200", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterLogLevel), "sdk.LogLevel", "sdk.LogLevelOff", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterMultiStatementCount), "int", "1", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterNoorderSequenceAsDefault), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterOdbcTreatDecimalAsInt), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterQueryTag), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterQuotedIdentifiersIgnoreCase), "bool", "false", "sdk.ParameterTypeAccount"},
			{string(sdk.UserParameterRowsPerResultset), "int", "0", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterS3StageVpceDnsName), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterSearchPath), "string", "$current, $public", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterSimulatedDataSharingConsumer), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterStatementQueuedTimeoutInSeconds), "int", "0", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterStatementTimeoutInSeconds), "int", "172800", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterStrictJsonOutput), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampDayIsAlways24h), "bool", "false", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampInputFormat), "string", "AUTO", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampLtzOutputFormat), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampNtzOutputFormat), "string", "YYYY-MM-DD HH24:MI:SS.FF3", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampOutputFormat), "string", "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampTypeMapping), "sdk.TimestampTypeMapping", "sdk.TimestampTypeMappingNtz", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimestampTzOutputFormat), "string", "", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimezone), "string", "America/Los_Angeles", "sdk.ParameterTypeAccount"},
			{string(sdk.UserParameterTimeInputFormat), "string", "AUTO", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTimeOutputFormat), "string", "HH24:MI:SS", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTraceLevel), "sdk.TraceLevel", "sdk.TraceLevelOff", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTransactionAbortOnError), "bool", "false", "sdk.ParameterTypeAccount"},
			{string(sdk.UserParameterTransactionDefaultIsolationLevel), "sdk.TransactionDefaultIsolationLevel", "sdk.TransactionDefaultIsolationLevelReadCommitted", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterTwoDigitCenturyStart), "int", "1970", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterUnsupportedDdlAction), "sdk.UnsupportedDDLAction", "sdk.UnsupportedDDLActionIgnore", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterUseCachedResult), "bool", "true", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterWeekOfYearPolicy), "int", "0", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.UserParameterWeekStart), "int", "0", "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Warehouse",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeWarehouse,
		Parameters: []gen.SnowflakeParameter{
			{string(sdk.WarehouseParameterMaxConcurrencyLevel), "int", "8", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds), "int", "0", "sdk.ParameterTypeSnowflakeDefault"},
			{string(sdk.WarehouseParameterStatementTimeoutInSeconds), "int", "172800", "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
}
