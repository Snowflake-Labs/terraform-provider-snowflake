package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var functionArgument = g.NewQueryStruct("FunctionArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ArgDataType", "DataType", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("DefaultValue", "*string", g.ParameterOptions().NoEquals().SQL("DEFAULT"))

var functionColumn = g.NewQueryStruct("FunctionColumn").
	Text("ColumnName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ColumnDataType", "DataType", g.KeywordOptions().NoQuotes().Required())

var functionReturns = g.NewQueryStruct("FunctionReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("FunctionReturnsResultDataType").
			PredefinedQueryStructField("ResultDataType", "DataType", g.KeywordOptions().NoQuotes().Required()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("FunctionReturnsTable").
			ListQueryStructField(
				"Columns",
				functionColumn,
				g.ParameterOptions().Parentheses().NoEquals(),
			),
		g.KeywordOptions().SQL("TABLE"),
	).WithValidation(g.ExactlyOneValueSet, "ResultDataType", "Table")

var (
	functionImports  = g.NewQueryStruct("FunctionImport").Text("Import", g.KeywordOptions().SingleQuotes())
	functionPackages = g.NewQueryStruct("FunctionPackage").Text("Package", g.KeywordOptions().SingleQuotes())
)

var FunctionsDef = g.NewInterface(
	"Functions",
	"Function",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#java-handler",
	g.NewQueryStruct("CreateForJava").
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
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE JAVA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
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
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "Secret", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateForJavascript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#javascript-handler",
	g.NewQueryStruct("CreateForJavascript").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#python-handler",
	g.NewQueryStruct("CreateForPython").
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
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE PYTHON").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
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
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "Secret", g.ParameterOptions().Parentheses()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#scala-handler",
	g.NewQueryStruct("CreateForScala").
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
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		PredefinedQueryStructField("ResultDataType", "DataType", g.ParameterOptions().NoEquals().SQL("RETURNS").Required()).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE SCALA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
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
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#sql-handler",
	g.NewQueryStruct("CreateForSQL").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Name().
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalSQL("MEMOIZABLE").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().MustParentheses().Required()).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SET LOG_LEVEL", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SET TRACE_LEVEL", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("SET SECURE").
		OptionalSQL("UNSET SECURE").
		OptionalSQL("UNSET LOG_LEVEL").
		OptionalSQL("UNSET TRACE_LEVEL").
		OptionalSQL("UNSET COMMENT").
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "SetSecure", "UnsetLogLevel", "UnsetTraceLevel", "UnsetSecure", "UnsetComment", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-function",
	g.NewQueryStruct("DropFunction").
		Drop().
		SQL("FUNCTION").
		IfExists().
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().MustParentheses().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-user-functions",
	g.DbStruct("functionRow").
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
		Field("is_memoizable", "sql.NullString"),
	g.PlainStruct("Function").
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
		Field("IsMemoizable", "bool"),
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("USER FUNCTIONS").
		OptionalLike().
		OptionalIn(),
).ShowByIdOperation().DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-function",
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
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().MustParentheses().Required()).
		WithValidation(g.ValidIdentifier, "name"),
)
