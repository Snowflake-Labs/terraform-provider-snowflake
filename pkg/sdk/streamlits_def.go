package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var streamlitSet = g.NewQueryStruct("StreamlitSet").
	OptionalTextAssignment("ROOT_LOCATION", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes().Required()).
	OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var StreamlitsDef = g.NewInterface(
	"Streamlits",
	"Streamlit",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-streamlit",
	g.NewQueryStruct("CreateStreamlit").
		Create().
		OrReplace().
		SQL("STREAMLIT").
		IfNotExists().
		Name().
		TextAssignment("ROOT_LOCATION", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes().Required()).
		OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-streamlit",
	g.NewQueryStruct("AlterStreamlit").
		Alter().
		SQL("STREAMLIT").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			streamlitSet,
			g.KeywordOptions().SQL("SET"),
		).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-streamlit",
	g.NewQueryStruct("DropStreamlit").
		Drop().
		SQL("STREAMLIT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-streamlits",
	g.DbStruct("streamlitsRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		Field("title", "sql.NullString").
		Field("owner", "string").
		Field("comment", "sql.NullString").
		Field("query_warehouse", "sql.NullString").
		Field("url_id", "string").
		Field("owner_role_type", "string"),
	g.PlainStruct("Streamlit").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Title", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("QueryWarehouse", "string").
		Field("UrlId", "string").
		Field("OwnerRoleType", "string"),
	g.NewQueryStruct("ShowStreamlits").
		Show().
		Terse().
		SQL("STREAMLITS").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
).ShowByIdOperation().DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-streamlit",
	g.DbStruct("streamlitsDetailRow").
		Field("name", "string").
		Field("title", "sql.NullString").
		Field("root_location", "string").
		Field("main_file", "string").
		Field("query_warehouse", "sql.NullString").
		Field("url_id", "string"),
	g.PlainStruct("StreamlitDetail").
		Field("Name", "string").
		Field("Title", "string").
		Field("RootLocation", "string").
		Field("MainFile", "string").
		Field("QueryWarehouse", "string").
		Field("UrlId", "string"),
	g.NewQueryStruct("DescribeStreamlit").
		Describe().
		SQL("STREAMLIT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
