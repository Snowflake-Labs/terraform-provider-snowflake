package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var functionArgument = g.NewQueryStruct("FunctionArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes()).
	Text("ArgDataType", g.KeywordOptions().NoQuotes())

var functionArgumentType = g.NewQueryStruct("FunctionArgumentType").
	Text("ArgDataType", g.KeywordOptions().NoQuotes())

var functionColumn = g.NewQueryStruct("FunctionColumn").
	Text("ColumnName", g.KeywordOptions().NoQuotes()).
	Text("ColumnDataType", g.KeywordOptions().NoQuotes())

var functionSecret = g.NewQueryStruct("FunctionSecret").
	Text("SecretVariableName", g.KeywordOptions().SingleQuotes()).
	Text("SecretName", g.KeywordOptions().NoQuotes())

var functionReturns = g.NewQueryStruct("FunctionReturns").
	OptionalText("ResultDataType", g.KeywordOptions()).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("FunctionReturnsTable").
			ListQueryStructField(
				"Columns",
				functionColumn,
				g.ParameterOptions().Parentheses().NoEquals(),
			),
		g.KeywordOptions().SQL("TABLE"),
	)

var functionSet = g.NewQueryStruct("FunctionSet").
	OptionalTextAssignment("LOG_LEVEL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("TRACE_LEVEL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalSQL("SECURE")

var functionUnset = g.NewQueryStruct("FunctionUnset").
	OptionalSQL("SECURE").
	OptionalSQL("COMMENT").
	OptionalSQL("LOG_LEVEL").
	OptionalSQL("TRACE_LEVEL")

var (
	functionNullOrNot     = g.NewQueryStruct("FunctionNullOrNot").OptionalSQL("NULL").OptionalSQL("NOT NULL")
	functionStrictOrNot   = g.NewQueryStruct("FunctionStrictOrNot").OptionalSQL("STRICT").OptionalSQL("CALLED ON NULL INPUT")
	functionVolatileOrNot = g.NewQueryStruct("FunctionVolatileOrNot").OptionalSQL("VOLATILE").OptionalSQL("IMMUTABLE")
	functionImports       = g.NewQueryStruct("FunctionImports").Text("Import", g.KeywordOptions().SingleQuotes())
	functionPackages      = g.NewQueryStruct("FunctionPackages").Text("Package", g.KeywordOptions().SingleQuotes())
	functionDefinition    = g.NewQueryStruct("FunctionDefinition").Text("Definition", g.KeywordOptions())
)

var FunctionsDef = g.NewInterface(
	"Functions",
	"Function",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateFunctionForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunctionForJava").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSQL("COPY GRANTS").
		OptionalQueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		SQL("LANGUAGE JAVA").
		OptionalQueryStructField(
			"StrictOrNot",
			functionStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			functionVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		OptionalTextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListQueryStructField(
			"Secrets",
			functionSecret,
			g.ParameterOptions().Parentheses().SQL("SECRETS"),
		).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"FunctionDefinition",
			functionDefinition,
			g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateFunctionForJavascript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunctionForJavascript").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSQL("COPY GRANTS").
		OptionalQueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		SQL("LANGUAGE JAVASCRIPT").
		OptionalQueryStructField(
			"StrictOrNot",
			functionStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			functionVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"FunctionDefinition",
			functionDefinition,
			g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateFunctionForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunctionForPython").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSQL("COPY GRANTS").
		OptionalQueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		SQL("LANGUAGE PYTHON").
		OptionalQueryStructField(
			"StrictOrNot",
			functionStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			functionVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		OptionalTextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListQueryStructField(
			"Secrets",
			functionSecret,
			g.ParameterOptions().Parentheses().SQL("SECRETS"),
		).
		OptionalQueryStructField(
			"FunctionDefinition",
			functionDefinition,
			g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateFunctionForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunctionForScala").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSQL("COPY GRANTS").
		OptionalQueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		SQL("LANGUAGE SCALA").
		OptionalQueryStructField(
			"StrictOrNot",
			functionStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			functionVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		OptionalTextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"FunctionDefinition",
			functionDefinition,
			g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateFunctionForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunctionForSQL").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSQL("COPY GRANTS").
		OptionalQueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			functionVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalSQL("MEMOIZABLE").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"FunctionDefinition",
			functionDefinition,
			g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		ListQueryStructField(
			"ArgumentTypes",
			functionArgumentType,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalQueryStructField(
			"Set",
			functionSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			functionUnset,
			g.KeywordOptions().SQL("UNSET"),
		).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		SetTags().UnsetTags().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-function",
	g.NewQueryStruct("DropFunction").
		Drop().
		SQL("FUNCTION").
		IfExists().
		Name().
		ListQueryStructField(
			"ArgumentTypes",
			functionArgumentType,
			g.ParameterOptions().Parentheses().NoEquals()).
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-functions",
	g.DbStruct("functionRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("min_num_arguments", "int").
		Field("max_num_arguments", "int").
		Field("arguments", "string").
		Field("is_table_function", "string").
		Field("is_secure", "string").
		Field("is_external_function", "string").
		Field("language", "string").
		Field("is_memoizable", "string"),
	g.PlainStruct("Function").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("MinNumArguments", "int").
		Field("MaxNumArguments", "int").
		Field("Arguments", "string").
		Field("IsTableFunction", "bool").
		Field("IsSecure", "bool").
		Field("IsExternalFunction", "bool").
		Field("Language", "string").
		Field("IsMemoizable", "bool"),
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("USER FUNCTIONS").
		OptionalLike().
		OptionalIn(),
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/describe-function",
	g.DbStruct("functionDetailRow").
		Field("property", "string").
		Field("value", "string"),
	g.PlainStruct("FunctionDetail").
		Field("Property", "string").
		Field("Value", "string"),
	g.NewQueryStruct("DescribeFunction").
		Describe().
		SQL("FUNCTION").
		Name().
		ListQueryStructField(
			"ArgumentTypes",
			functionArgumentType,
			g.ParameterOptions().Parentheses().NoEquals()).
		WithValidation(g.ValidIdentifier, "name"),
)
