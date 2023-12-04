package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var procedureArgument = g.NewQueryStruct("ProcedureArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ArgDataType", "DataType", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("DefaultValue", "*string", g.ParameterOptions().NoEquals().SQL("DEFAULT"))

var procedureColumn = g.NewQueryStruct("ProcedureColumn").
	Text("ColumnName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ColumnDataType", "DataType", g.KeywordOptions().NoQuotes().Required())

var procedureReturns = g.NewQueryStruct("ProcedureReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			PredefinedQueryStructField("ResultDataType", "DataType", g.KeywordOptions().NoQuotes().Required()).
			OptionalSQL("NULL").OptionalSQL("NOT NULL"),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("ProcedureReturnsTable").
			ListQueryStructField(
				"Columns",
				procedureColumn,
				g.ParameterOptions().Parentheses().NoEquals(),
			),
		g.KeywordOptions().SQL("TABLE"),
	)

var procedureSQLReturns = g.NewQueryStruct("ProcedureSQLReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			PredefinedQueryStructField("ResultDataType", "DataType", g.KeywordOptions().NoQuotes().Required()),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("ProcedureReturnsTable").
			ListQueryStructField(
				"Columns",
				procedureColumn,
				g.ParameterOptions().Parentheses().NoEquals(),
			),
		g.KeywordOptions().SQL("TABLE"),
	).
	OptionalSQL("NOT NULL")

var (
	procedureImport  = g.NewQueryStruct("ProcedureImport").Text("Import", g.KeywordOptions().SingleQuotes().Required())
	procedurePackage = g.NewQueryStruct("ProcedurePackage").Text("Package", g.KeywordOptions().SingleQuotes().Required())
)

var ProceduresDef = g.NewInterface(
	"Procedures",
	"Procedure",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#java-handler",
	g.NewQueryStruct("CreateForJava").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Name().
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ParameterOptions().Parentheses().NoEquals(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE JAVA").
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "Secret", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForJavaScript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#javascript-handler",
	g.NewQueryStruct("CreateForJavaScript").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Name().
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ParameterOptions().Parentheses().NoEquals(),
		).
		OptionalSQL("COPY GRANTS").
		PredefinedQueryStructField("ResultDataType", "DataType", g.ParameterOptions().NoEquals().SQL("RETURNS").Required()).
		OptionalSQL("NOT NULL").
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#python-handler",
	g.NewQueryStruct("CreateForPython").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Name().
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ParameterOptions().Parentheses().NoEquals(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE PYTHON").
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "Secret", g.ParameterOptions().Parentheses()).
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#scala-handler",
	g.NewQueryStruct("CreateForScala").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Name().
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ParameterOptions().Parentheses().NoEquals(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SCALA").
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#snowflake-scripting-handler",
	g.NewQueryStruct("CreateForSQL").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Name().
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ParameterOptions().Parentheses().NoEquals(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureSQLReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SQL").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-procedure",
	g.NewQueryStruct("AlterProcedure").
		Alter().
		SQL("PROCEDURE").
		IfExists().
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().Parentheses().Required()).
		OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SET LOG_LEVEL", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("SET TRACE_LEVEL", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("UNSET COMMENT").
		OptionalSetTags().
		OptionalUnsetTags().
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.KeywordOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "UnsetComment", "SetTags", "UnsetTags", "ExecuteAs"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-procedure",
	g.NewQueryStruct("DropProcedure").
		Drop().
		SQL("PROCEDURE").
		IfExists().
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().Parentheses().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-procedures",
	g.DbStruct("procedureRow").
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
		Field("is_secure", "sql.NullString"),
	g.PlainStruct("Procedure").
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
		Field("IsSecure", "bool"),
	g.NewQueryStruct("ShowProcedures").
		Show().
		SQL("PROCEDURES").
		OptionalLike().
		OptionalIn(), // TODO: 'In' struct for procedures not support keyword "CLASS" now
).ShowByIdOperation().DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-procedure",
	g.DbStruct("procedureDetailRow").
		Field("property", "string").
		Field("value", "string"),
	g.PlainStruct("ProcedureDetail").
		Field("Property", "string").
		Field("Value", "string"),
	g.NewQueryStruct("DescribeProcedure").
		Describe().
		SQL("PROCEDURE").
		Name().
		PredefinedQueryStructField("ArgumentDataTypes", "[]DataType", g.KeywordOptions().Parentheses().Required()).
		WithValidation(g.ValidIdentifier, "name"),
)
