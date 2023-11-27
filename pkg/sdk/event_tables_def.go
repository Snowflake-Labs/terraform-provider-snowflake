package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var eventTableSet = g.NewQueryStruct("EventTableSet").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var eventTableUnset = g.NewQueryStruct("EventTableUnset").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("CHANGE_TRACKING").
	OptionalSQL("COMMENT")

var eventTableDropRowAccessPolicy = g.NewQueryStruct("EventTableDropRowAccessPolicy").
	SQL("ROW ACCESS POLICY").
	Identifier("Name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions())

var eventTableClusteringAction = g.NewQueryStruct("EventTableClusteringAction").
	PredefinedQueryStructField("ClusterBy", "*[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
	OptionalSQL("SUSPEND RECLUSTER").
	OptionalSQL("RESUME RECLUSTER").
	OptionalSQL("DROP CLUSTERING KEY")

var searchOptimization = g.NewQueryStruct("SearchOptimization").
	SQL("SEARCH OPTIMIZATION").
	PredefinedQueryStructField("On", "[]string", g.KeywordOptions().SQL("ON"))

var eventTableSearchOptimizationAction = g.NewQueryStruct("EventTableSearchOptimizationAction").
	OptionalQueryStructField(
		"Add",
		searchOptimization,
		g.KeywordOptions().SQL("ADD"),
	).
	OptionalQueryStructField(
		"Drop",
		searchOptimization,
		g.KeywordOptions().SQL("DROP"),
	)

var EventTablesDef = g.NewInterface(
	"EventTables",
	"EventTable",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-event-table",
	g.NewQueryStruct("CreateEventTable").
		Create().
		OrReplace().
		SQL("EVENT TABLE").
		IfNotExists().
		Name().
		PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("COPY GRANTS").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("RowAccessPolicy", "*RowAccessPolicy", g.KeywordOptions()).
		OptionalTags().WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-event-tables",
	g.DbStruct("eventTableRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		Field("owner", "string").
		Field("comment", "string").
		Field("owner_role_type", "string").
		Field("change_tracking", "string"),
	g.PlainStruct("EventTable").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("OwnerRoleType", "string").
		Field("ChangeTracking", "bool"),
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("EVENT TABLES").
		OptionalLike().
		OptionalIn().
		OptionalTextAssignment("STARTS WITH", g.ParameterOptions().SingleQuotes().NoEquals()).
		OptionalNumberAssignment("LIMIT", g.ParameterOptions()).
		OptionalTextAssignment("FROM", g.ParameterOptions().SingleQuotes().NoEquals()),
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/describe-event-table",
	g.DbStruct("eventTableDetailsRow").
		Field("name", "string").
		Field("kind", "string").
		Field("comment", "string"),
	g.PlainStruct("EventTableDetails").
		Field("Name", "string").
		Field("Kind", "string").
		Field("Comment", "string"),
	g.NewQueryStruct("DescribeEventTable").
		Describe().
		SQL("EVENT TABLE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-event-table",
	g.NewQueryStruct("AlterEventTable").
		Alter().
		SQL("TABLE").
		IfNotExists().
		Name().
		OptionalQueryStructField(
			"Set",
			eventTableSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			eventTableUnset,
			g.KeywordOptions().SQL("UNSET"),
		).
		PredefinedQueryStructField("AddRowAccessPolicy", "*RowAccessPolicy", g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField(
			"DropRowAccessPolicy",
			eventTableDropRowAccessPolicy,
			g.KeywordOptions().SQL("DROP"),
		).
		OptionalSQL("DROP ALL ROW ACCESS POLICIES").
		OptionalQueryStructField(
			"ClusteringAction",
			eventTableClusteringAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"SearchOptimizationAction",
			eventTableSearchOptimizationAction,
			g.KeywordOptions(),
		).
		PredefinedQueryStructField("SetTags", "[]TagAssociation", g.KeywordOptions().SQL("SET TAG")).
		PredefinedQueryStructField("UnsetTags", "[]ObjectIdentifier", g.KeywordOptions().SQL("UNSET TAG")).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		WithValidation(g.ValidIdentifier, "name"),
)
