package sdk

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

var alterServiceSet = g.NewQueryStruct("CortexSearchServiceSet").
	// Fields
	OptionalTextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "TargetLag", "Warehouse", "Comment")

var CortexSearchServiceDef = g.NewInterface(
	"CortexSearchServices",
	"CortexSearchService",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/create-cortex-search",
	g.NewQueryStruct("CreateCortexSearchService").
		// Fields
		Create().
		OrReplace().
		SQL("CORTEX SEARCH SERVICE").
		IfNotExists().
		Name().
		TextAssignment("ON", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("ON")).
		OptionalQueryStructField(
			"Attributes",
			g.NewQueryStruct("Attributes").
				SQL("ATTRIBUTES").
				List("Columns", "string", g.ListOptions().NoEquals().NoParentheses()),
			g.KeywordOptions(),
		).
		Identifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().Required().SQL("WAREHOUSE")).
		TextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes().Required()).
		OptionalComment().
		PredefinedQueryStructField("QueryDefinition", "string", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("AS")).
		// Validations
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "TargetLag").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/alter-cortex-search",
	g.NewQueryStruct("AlterCortexSearchService").
		// Fields
		Alter().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		OptionalQueryStructField("Set", alterServiceSet, g.KeywordOptions().SQL("SET")).
		// Validations
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set"),
).ShowOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/show-cortex-search",
	// Fields
	g.DbStruct("cortexSearchServiceRow").
		Field("created_on", "time.Time").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		Field("comment", "string"),
	g.PlainStruct("CortexSearchService").
		Field("CreatedOn", "time.Time").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Comment", "string"),
	g.NewQueryStruct("ShowCortexSearchService").
		Show().
		SQL("CORTEX SEARCH SERVICES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
).ShowByIdOperation().DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/desc-cortex-search",
	g.DbStruct("cortexSearchServiceDetailsRow").
		Field("name", "string").
		Field("schema", "string").
		Field("database", "string").
		Field("warehouse", "string").
		Field("target_lag", "string").
		Field("search_column", "string").
		OptionalText("included_columns").
		Field("service_url", "string").
		OptionalText("refreshed_on").
		OptionalNumber("num_rows_indexed").
		OptionalText("comment"),
	g.PlainStruct("CortexSearchServiceDetails").
		Field("Name", "string").
		Field("Schema", "string").
		Field("Database", "string").
		Field("Warehouse", "string").
		Field("TargetLag", "string").
		Field("On", "string").
		Field("Attributes", "[]string").
		Field("ServiceUrl", "string").
		Field("RefreshedOn", "string").
		Field("NumRowsIndexed", "int").
		Field("Comment", "string"),
	g.NewQueryStruct("DescribeCortexSearchService").
		Describe().
		SQL("CORTEX SEARCH SERVICE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/drop-cortex-search",
	g.NewQueryStruct("DropCortexSearchService").
		// Fields
		Drop().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
