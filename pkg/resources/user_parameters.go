package resources

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	UserParametersSchema = make(map[string]*schema.Schema)
)

func init() {
	userParameterFields := []struct {
		Name                   sdk.UserParameter
		Type                   schema.ValueType
		Description            string
		SnowflakeDocsReference string
		//DiffSuppress schema.SchemaDiffSuppressFunc
		//ValidateDiag schema.SchemaValidateDiagFunc
	}{
		{Name: sdk.UserParameterAbortDetachedQuery, Type: schema.TypeBool, Description: "Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption).", SnowflakeDocsReference: "abort-detached-query"},
		{Name: sdk.UserParameterAutocommit, Type: schema.TypeBool, Description: "Specifies whether autocommit is enabled for the session. Autocommit determines whether a DML statement, when executed without an active transaction, is automatically committed after the statement successfully completes. For more information, see [Transactions](https://docs.snowflake.com/en/sql-reference/transactions).", SnowflakeDocsReference: "autocommit"},
		{Name: sdk.UserParameterBinaryInputFormat, Type: schema.TypeString, Description: "The format of VARCHAR values passed as input to VARCHAR-to-BINARY conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output).", SnowflakeDocsReference: "binary-input-format"},
		{Name: sdk.UserParameterBinaryOutputFormat, Type: schema.TypeString, Description: "The format for VARCHAR values returned as output by BINARY-to-VARCHAR conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output).", SnowflakeDocsReference: "binary-output-format"},
		{Name: sdk.UserParameterClientMemoryLimit, Type: schema.TypeInt, Description: "Parameter that specifies the maximum amount of memory the JDBC driver or ODBC driver should use for the result set from queries (in MB).", SnowflakeDocsReference: "client-memory-limit"},
		{Name: sdk.UserParameterClientMetadataRequestUseConnectionCtx, Type: schema.TypeBool, Description: "For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases/schemas to the current database/schema. The narrower search typically returns fewer rows and executes more quickly.", SnowflakeDocsReference: "client-metadata-request-use-connection-ctx"},
		{Name: sdk.UserParameterClientPrefetchThreads, Type: schema.TypeInt, Description: "Parameter that specifies the number of threads used by the client to pre-fetch large result sets. The driver will attempt to honor the parameter value, but defines the minimum and maximum values (depending on your system’s resources) to improve performance.", SnowflakeDocsReference: "client-prefetch-threads"},
		{Name: sdk.UserParameterClientResultChunkSize, Type: schema.TypeInt, Description: "Parameter that specifies the maximum size of each set (or chunk) of query results to download (in MB). The JDBC driver downloads query results in chunks.", SnowflakeDocsReference: "client-result-chunk-size"},
		{Name: sdk.UserParameterClientResultColumnCaseInsensitive, Type: schema.TypeBool, Description: "Parameter that indicates whether to match column name case-insensitively in ResultSet.get* methods in JDBC.", SnowflakeDocsReference: "client-result-column-case-insensitive"},
		{Name: sdk.UserParameterClientSessionKeepAlive, Type: schema.TypeBool, Description: "Parameter that indicates whether to force a user to log in again after a period of inactivity in the session.", SnowflakeDocsReference: "client-session-keep-alive"},
		{Name: sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, Type: schema.TypeInt, Description: "Number of seconds in-between client attempts to update the token for the session.", SnowflakeDocsReference: "client-session-keep-alive-heartbeat-frequency"},
		{Name: sdk.UserParameterClientTimestampTypeMapping, Type: schema.TypeString, Description: "Specifies the [TIMESTAMP_* variation](https://docs.snowflake.com/en/sql-reference/data-types-datetime.html#label-datatypes-timestamp-variations) to use when binding timestamp variables for JDBC or ODBC applications that use the bind API to load data.", SnowflakeDocsReference: "client-timestamp-type-mapping"},
		{Name: sdk.UserParameterDateInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "date-input-format"},
		{Name: sdk.UserParameterDateOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "date-output-format"},
		{Name: sdk.UserParameterEnableUnloadPhysicalTypeOptimization, Type: schema.TypeBool, Description: "Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table).", SnowflakeDocsReference: "enable-unload-physical-type-optimization"},
		{Name: sdk.UserParameterErrorOnNondeterministicMerge, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row.", SnowflakeDocsReference: "error-on-nondeterministic-merge"},
		{Name: sdk.UserParameterErrorOnNondeterministicUpdate, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row.", SnowflakeDocsReference: "error-on-nondeterministic-update"},
		{Name: sdk.UserParameterGeographyOutputFormat, Type: schema.TypeString, Description: "Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography).", SnowflakeDocsReference: "geography-output-format"},
		{Name: sdk.UserParameterGeometryOutputFormat, Type: schema.TypeString, Description: "Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry).", SnowflakeDocsReference: "geometry-output-format"},
		{Name: sdk.UserParameterJdbcTreatDecimalAsInt, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterJdbcTreatTimestampNtzAsUtc, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterJdbcUseSessionTimezone, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterJsonIndent, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterLockTimeout, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterLogLevel, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterMultiStatementCount, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterNoorderSequenceAsDefault, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterOdbcTreatDecimalAsInt, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterQueryTag, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterQuotedIdentifiersIgnoreCase, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterRowsPerResultset, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterS3StageVpceDnsName, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterSearchPath, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterSimulatedDataSharingConsumer, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterStatementQueuedTimeoutInSeconds, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterStatementTimeoutInSeconds, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterStrictJsonOutput, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampDayIsAlways24h, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampInputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampLtzOutputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampNtzOutputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampOutputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampTypeMapping, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimestampTzOutputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimezone, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimeInputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTimeOutputFormat, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTraceLevel, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTransactionAbortOnError, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTransactionDefaultIsolationLevel, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterTwoDigitCenturyStart, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterUnsupportedDdlAction, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterUseCachedResult, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterWeekOfYearPolicy, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterWeekStart, Type: schema., Description: "", SnowflakeDocsReference: ""},

		{Name: sdk.UserParameterEnableUnredactedQuerySyntaxError, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterNetworkPolicy, Type: schema., Description: "", SnowflakeDocsReference: ""},
		{Name: sdk.UserParameterPreventUnloadToInternalStages, Type: schema., Description: "", SnowflakeDocsReference: ""},
	}

	for _, field := range userParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		UserParametersSchema[fieldName] = &schema.Schema{
			Type:        field.Type,
			Description: field.Description,
			Computed:    true,
			Optional:    true,
			//ValidateDiagFunc: field.ValidateDiag,
			//DiffSuppressFunc: field.DiffSuppress,
		}
	}
}
