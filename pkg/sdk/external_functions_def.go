package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var externalFunctionArgument = g.NewQueryStruct("ExternalFunctionArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ArgDataType", "DataType", g.KeywordOptions().NoQuotes().Required())

var externalFunctionHeader = g.NewQueryStruct("ExternalFunctionHeader").
	Text("Name", g.KeywordOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("Value", "string", g.ParameterOptions().SingleQuotes().Required())

var externalFunctionContextHeader = g.NewQueryStruct("ExternalFunctionContextHeader").Text("ContextFunction", g.KeywordOptions().NoQuotes())

var externalFunctionSet = g.NewQueryStruct("ExternalFunctionSet").
	OptionalIdentifier("ApiIntegration", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION =")).
	ListQueryStructField(
		"Headers",
		externalFunctionHeader,
		g.ParameterOptions().Parentheses().SQL("HEADERS"),
	).
	ListQueryStructField(
		"ContextHeaders",
		externalFunctionContextHeader,
		g.ParameterOptions().Parentheses().SQL("CONTEXT_HEADERS"),
	).
	OptionalNumberAssignment("MAX_BATCH_ROWS", g.ParameterOptions()).
	OptionalTextAssignment("COMPRESSION", g.ParameterOptions()).
	OptionalIdentifier("RequestTranslator", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REQUEST_TRANSLATOR =")).
	OptionalIdentifier("ResponseTranslator", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RESPONSE_TRANSLATOR =")).
	WithValidation(g.ExactlyOneValueSet, "ApiIntegration", "Headers", "ContextHeaders", "MaxBatchRows", "Compression", "RequestTranslator", "ResponseTranslator")

var externalFunctionUnset = g.NewQueryStruct("ExternalFunctionUnset").
	OptionalSQL("COMMENT").
	OptionalSQL("HEADERS").
	OptionalSQL("CONTEXT_HEADERS").
	OptionalSQL("MAX_BATCH_ROWS").
	OptionalSQL("COMPRESSION").
	OptionalSQL("SECURE").
	OptionalSQL("REQUEST_TRANSLATOR").
	OptionalSQL("RESPONSE_TRANSLATOR")

var ExternalFunctionsDef = g.NewInterface(
	"ExternalFunctions",
	"ExternalFunction",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-function",
	g.NewQueryStruct("CreateExternalFunction").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("EXTERNAL FUNCTION").
		Name().
		ListQueryStructField(
			"Arguments",
			externalFunctionArgument,
			g.ListOptions().MustParentheses()).
		PredefinedQueryStructField("ResultDataType", "DataType", g.ParameterOptions().NoEquals().SQL("RETURNS").Required()).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		Identifier("ApiIntegration", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION =").Required()).
		ListQueryStructField(
			"Headers",
			externalFunctionHeader,
			g.ParameterOptions().Parentheses().SQL("HEADERS"),
		).
		ListQueryStructField(
			"ContextHeaders",
			externalFunctionContextHeader,
			g.ParameterOptions().Parentheses().SQL("CONTEXT_HEADERS"),
		).
		OptionalNumberAssignment("MAX_BATCH_ROWS", g.ParameterOptions()).
		OptionalTextAssignment("COMPRESSION", g.ParameterOptions()).
		OptionalIdentifier("RequestTranslator", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REQUEST_TRANSLATOR =")).
		OptionalIdentifier("ResponseTranslator", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RESPONSE_TRANSLATOR =")).
		TextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "ApiIntegration").
		WithValidation(g.ValidIdentifierIfSet, "ApiIntegration").
		WithValidation(g.ValidIdentifierIfSet, "RequestTranslator").
		WithValidation(g.ValidateValueSet, "As").
		WithValidation(g.ValidIdentifierIfSet, "ResponseTranslator"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterExternalFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().MustParentheses().Required()).
		OptionalQueryStructField(
			"Set",
			externalFunctionSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			externalFunctionUnset,
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-external-functions",
	g.DbStruct("externalFunctionRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("is_builtin", "string").
		Field("is_aggregate", "string").
		Field("is_ansi", "string").
		Field("min_num_arguments", "int").
		Field("max_num_arguments", "int").
		Field("arguments", "string").
		Field("description", "string").
		Field("catalog_name", "string").
		Field("is_table_function", "string").
		Field("valid_for_clustering", "string").
		Field("is_secure", "sql.NullString").
		Field("is_external_function", "string").
		Field("language", "string").
		Field("is_memoizable", "sql.NullString").
		Field("is_data_metric", "sql.NullString"),
	g.PlainStruct("ExternalFunction").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("IsBuiltin", "bool").
		Field("IsAggregate", "bool").
		Field("IsAnsi", "bool").
		Field("MinNumArguments", "int").
		Field("MaxNumArguments", "int").
		Field("Arguments", "string").
		Field("Description", "string").
		Field("CatalogName", "string").
		Field("IsTableFunction", "bool").
		Field("ValidForClustering", "bool").
		Field("IsSecure", "bool").
		Field("IsExternalFunction", "bool").
		Field("Language", "string").
		Field("IsMemoizable", "bool").
		Field("IsDataMetric", "bool"),
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("EXTERNAL FUNCTIONS").
		OptionalLike(),
).ShowByIdOperation()
