package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var functionArgument = g.NewQueryStruct("FunctionArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes()).
	DataType("ArgDataType")

var functionArgumentType = g.NewQueryStruct("FunctionArgumentType").
	DataType("ArgDataType")

var functionColumn = g.NewQueryStruct("FunctionColumn").
	Text("ColumnName", g.KeywordOptions().NoQuotes()).
	DataType("ColumnDataType")

var functionReturns = g.NewQueryStruct("FunctionReturns").
	DataType("ResultDataType").
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
/*
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
*/
var (
	//functionNullOrNot     = g.NewQueryStruct("FunctionNullOrNot").OptionalSQL("NULL").OptionalSQL("NOT NULL")
	//functionStrictOrNot   = g.NewQueryStruct("FunctionStrictOrNot").OptionalSQL("STRICT").OptionalSQL("CALLED ON NULL INPUT")
	//functionVolatileOrNot = g.NewQueryStruct("FunctionVolatileOrNot").OptionalSQL("VOLATILE").OptionalSQL("IMMUTABLE")
	functionImports       = g.NewQueryStruct("FunctionImports").Text("Import", g.KeywordOptions().SingleQuotes())
	functionPackages      = g.NewQueryStruct("FunctionPackages").Text("Package", g.KeywordOptions().SingleQuotes())
	functionSecret 	  = g.NewQueryStruct("FunctionSecret"). .Text("SecretVariableName", g.KeywordOptions().SingleQuotes()).Text("SecretName", g.KeywordOptions().SingleQuotes())
)

var FunctionsDef = g.NewInterface(
	"Functions",
	"Function",
	g.KindOfT[SchemaObjectIdentifier](),
).WithEnums(
	g.NewEnum("FunctionsNullInputBehaviorType", []string{"CALLED ON NULL INPUT", "RETURNS NULL ON NULL INPUT", "STRICT"}),
	g.NewEnum("FunctionsResultsBehaviorType", []string{"Volatile", "Immutable"}),
).
CustomOperation("CreateJavaFunction",
"https://docs.snowflake.com/en/sql-reference/sql/create-function",
g.NewQueryStruct("CreateJavaFunction").
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
		g.ParameterOptions().Parentheses().NoEquals().Required(),
	).
	OptionalSQL("COPY GRANTS").
	QueryStructField(
		"Returns",
		functionReturns,
		g.KeywordOptions().SQL("RETURNS").Required(),
	).
	OptionalSQL("NOT NULL").
	Assignment("LANGUAGE", "LanguageType", g.ParameterOptions().NoEquals().NoQuotes()).
	Enum("FunctionsNullInputBehavior", "FunctionsNullInputBehaviorType").
	Enum("FunctionsResultsBehavior", "*FunctionsResultsBehaviorType").
	OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	ListAssignment("IMPORTS", "string", g.ParameterOptions().Parentheses().SingleQuotes()).
	ListAssignment("PACKAGES", "string", g.ParameterOptions().Parentheses().SingleQuotes()).
	TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
	ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", g.KindOfT[AccountObjectIdentifier](), g.ParameterOptions().Parentheses()).
	ListQueryStructField("SECRETS", g.KindOfT[UDFSecret](), g.ParameterOptions().Parentheses()).
	OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
	TextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes().Required()),
)




/*
.CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-function",
	g.NewQueryStruct("CreateFunction").
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
			g.ParameterOptions().Parentheses().NoEquals().Required()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		OptionalQueryStructField(
			"NullOrNot",
			functionNullOrNot,
			g.KeywordOptions(),
		).
		OptionalAssignment("LANGUAGE","LanguageType", g.ParameterOptions().NoEquals().NoQuotes()).
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
		OptionalSQL("MEMOIZABLE").
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
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
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
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")),
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
	g.DbStruct("functionDBRow").
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
).
ShowByIdOperation().
DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/describe-function",
	g.DbStruct("functionDetailDBRow").
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

type LanguageType string

const (
	LanguageTypeJava       LanguageType = "JAVA"
	LanguageTypeJavascript LanguageType = "JAVASCRIPT"
	LanguageTypePython     LanguageType = "PYTHON"
	LanguageTypeScala      LanguageType = "SCALA"
	LanguageTypeSQL        LanguageType = "SQL"
)

*/
