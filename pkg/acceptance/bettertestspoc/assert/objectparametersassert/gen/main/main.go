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
			{ParameterName: string(sdk.UserParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterTransactionDefaultIsolationLevel), ParameterType: "sdk.TransactionDefaultIsolationLevel", DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
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
		Parameters: []gen.SnowflakeParameter{
			{ParameterName: string(sdk.WarehouseParameterMaxConcurrencyLevel), ParameterType: "int", DefaultValue: "8", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
}
