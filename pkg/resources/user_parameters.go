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
		{Name: sdk.UserParameterJdbcTreatDecimalAsInt, Type: schema.TypeBool, Description: "Specifies how JDBC processes columns that have a scale of zero (0).", SnowflakeDocsReference: "jdbc-treat-decimal-as-int"},
		{Name: sdk.UserParameterJdbcTreatTimestampNtzAsUtc, Type: schema.TypeBool, Description: "Specifies how JDBC processes TIMESTAMP_NTZ values.", SnowflakeDocsReference: "jdbc-treat-timestamp-ntz-as-utc"},
		{Name: sdk.UserParameterJdbcUseSessionTimezone, Type: schema.TypeBool, Description: "Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class.", SnowflakeDocsReference: "jdbc-use-session-timezone"},
		{Name: sdk.UserParameterJsonIndent, Type: schema.TypeInt, Description: "Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element.", SnowflakeDocsReference: "json-indent"},
		{Name: sdk.UserParameterLockTimeout, Type: schema.TypeInt, Description: "Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement.", SnowflakeDocsReference: "lock-timeout"},
		{Name: sdk.UserParameterLogLevel, Type: schema.TypeString, Description: "Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested. For more information about log levels, see [Setting log level](https://docs.snowflake.com/en/developer-guide/logging-tracing/logging-log-level).", SnowflakeDocsReference: "log-level"},
		{Name: sdk.UserParameterMultiStatementCount, Type: schema.TypeInt, Description: "Number of statements to execute when using the multi-statement capability.", SnowflakeDocsReference: "multi-statement-count"},
		{Name: sdk.UserParameterNoorderSequenceAsDefault, Type: schema.TypeBool, Description: "Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values).", SnowflakeDocsReference: "noorder-sequence-as-default"},
		{Name: sdk.UserParameterOdbcTreatDecimalAsInt, Type: schema.TypeBool, Description: "Specifies how ODBC processes columns that have a scale of zero (0).", SnowflakeDocsReference: "odbc-treat-decimal-as-int"},
		{Name: sdk.UserParameterQueryTag, Type: schema.TypeString, Description: "Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions.", SnowflakeDocsReference: "query-tag"},
		{Name: sdk.UserParameterQuotedIdentifiersIgnoreCase, Type: schema.TypeBool, Description: "Specifies whether letters in double-quoted object identifiers are stored and resolved as uppercase letters. By default, Snowflake preserves the case of alphabetic characters when storing and resolving double-quoted identifiers (see [Identifier resolution](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing)). You can use this parameter in situations in which [third-party applications always use double quotes around identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing-parameter).", SnowflakeDocsReference: "quoted-identifiers-ignore-case"},
		{Name: sdk.UserParameterRowsPerResultset, Type: schema.TypeInt, Description: "Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum.", SnowflakeDocsReference: "rows-per-resultset"},
		{Name: sdk.UserParameterS3StageVpceDnsName, Type: schema.TypeString, Description: "Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation).", SnowflakeDocsReference: "s3-stage-vpce-dns-name"},
		{Name: sdk.UserParameterSearchPath, Type: schema.TypeString, Description: "Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name.", SnowflakeDocsReference: "search-path"},
		{Name: sdk.UserParameterSimulatedDataSharingConsumer, Type: schema.TypeString, Description: "Specifies the name of a consumer account to simulate for testing/validating shared data, particularly shared secure views. When this parameter is set in a session, shared views return rows as if executed in the specified consumer account rather than the provider account. For more information, see [Introduction to Secure Data Sharing](https://docs.snowflake.com/en/user-guide/data-sharing-intro) and [Working with shares](https://docs.snowflake.com/en/user-guide/data-sharing-provider).", SnowflakeDocsReference: "simulated-data-sharing-consumer"},
		{Name: sdk.UserParameterStatementQueuedTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged.", SnowflakeDocsReference: "statement-queued-timeout-in-seconds"},
		{Name: sdk.UserParameterStatementTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system.", SnowflakeDocsReference: "statement-timeout-in-seconds"},
		{Name: sdk.UserParameterStrictJsonOutput, Type: schema.TypeBool, Description: "This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON.", SnowflakeDocsReference: "strict-json-output"},
		{Name: sdk.UserParameterTimestampDayIsAlways24h, Type: schema.TypeBool, Description: "Specifies whether the [DATEADD](https://docs.snowflake.com/en/sql-reference/functions/dateadd) function (and its aliases) always consider a day to be exactly 24 hours for expressions that span multiple days.", SnowflakeDocsReference: "timestamp-day-is-always-24h"},
		{Name: sdk.UserParameterTimestampInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported timestamp format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of timestamps stored in the system during the session).", SnowflakeDocsReference: "timestamp-input-format"},
		{Name: sdk.UserParameterTimestampLtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_LTZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "timestamp-ltz-output-format"},
		{Name: sdk.UserParameterTimestampNtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_NTZ data type.", SnowflakeDocsReference: "timestamp-ntz-output-format"},
		{Name: sdk.UserParameterTimestampOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "timestamp-output-format"},
		{Name: sdk.UserParameterTimestampTypeMapping, Type: schema.TypeString, Description: "Specifies the TIMESTAMP_* variation that the TIMESTAMP data type alias maps to.", SnowflakeDocsReference: "timestamp-type-mapping"},
		{Name: sdk.UserParameterTimestampTzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_TZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "timestamp-tz-output-format"},
		{Name: sdk.UserParameterTimezone, Type: schema.TypeString, Description: "Specifies the time zone for the session. You can specify a [time zone name](https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab) or a [link name](https://data.iana.org/time-zones/tzdb-2021a/backward) from release 2021a of the [IANA Time Zone Database](https://www.iana.org/time-zones) (e.g. America/Los_Angeles, Europe/London, UTC, Etc/GMT, etc.).", SnowflakeDocsReference: "timezone"},
		{Name: sdk.UserParameterTimeInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session).", SnowflakeDocsReference: "time-input-format"},
		{Name: sdk.UserParameterTimeOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).", SnowflakeDocsReference: "time-output-format"},
		{Name: sdk.UserParameterTraceLevel, Type: schema.TypeString, Description: "Controls how trace events are ingested into the event table. For more information about trace levels, see [Setting trace level](https://docs.snowflake.com/en/developer-guide/logging-tracing/tracing-trace-level).", SnowflakeDocsReference: "trace-level"},
		{Name: sdk.UserParameterTransactionAbortOnError, Type: schema.TypeBool, Description: "Specifies the action to perform when a statement issued within a non-autocommit transaction returns with an error.", SnowflakeDocsReference: "transaction-abort-on-error"},
		{Name: sdk.UserParameterTransactionDefaultIsolationLevel, Type: schema.TypeString, Description: "Specifies the isolation level for transactions in the user session.", SnowflakeDocsReference: "transaction-default-isolation-level"},
		{Name: sdk.UserParameterTwoDigitCenturyStart, Type: schema.TypeInt, Description: "Specifies the “century start” year for 2-digit years (i.e. the earliest year such dates can represent). This parameter prevents ambiguous dates when importing or converting data with the `YY` date format component (i.e. years represented as 2 digits).", SnowflakeDocsReference: "two-digit-century-start"},
		{Name: sdk.UserParameterUnsupportedDdlAction, Type: schema.TypeString, Description: "Determines if an unsupported (i.e. non-default) value specified for a constraint property returns an error.", SnowflakeDocsReference: "unsupported-ddl-action"},
		{Name: sdk.UserParameterUseCachedResult, Type: schema.TypeBool, Description: "Specifies whether to reuse persisted query results, if available, when a matching query is submitted.", SnowflakeDocsReference: "use-cached-result"},
		{Name: sdk.UserParameterWeekOfYearPolicy, Type: schema.TypeInt, Description: "Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year.", SnowflakeDocsReference: "week-of-year-policy"},
		{Name: sdk.UserParameterWeekStart, Type: schema.TypeInt, Description: "Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week.", SnowflakeDocsReference: "week-start"},

		{Name: sdk.UserParameterEnableUnredactedQuerySyntaxError, Type: schema.TypeBool, Description: "Controls whether query text is redacted if a SQL query fails due to a syntax or parsing error. If `FALSE`, the content of a failed query is redacted in the views, pages, and functions that provide a query history. Only users with a role that is granted or inherits the AUDIT privilege can set the ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR parameter. When using the ALTER USER command to set the parameter to `TRUE` for a particular user, modify the user that you want to see the query text, not the user who executed the query (if those are different users).", SnowflakeDocsReference: "enable-unredacted-query-syntax-error"},
		{Name: sdk.UserParameterNetworkPolicy, Type: schema.TypeString, Description: "Specifies the network policy to enforce for your account. Network policies enable restricting access to your account based on users’ IP address. For more details, see [Controlling network traffic with network policies](https://docs.snowflake.com/en/user-guide/network-policies). Any existing network policy (created using [CREATE NETWORK POLICY](https://docs.snowflake.com/en/sql-reference/sql/create-network-policy)).", SnowflakeDocsReference: "network-policy"},
		{Name: sdk.UserParameterPreventUnloadToInternalStages, Type: schema.TypeBool, Description: "Specifies whether to prevent data unload operations to internal (Snowflake) stages using [COPY INTO <location>](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location) statements.", SnowflakeDocsReference: "prevent-unload-to-internal-stages"},
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
