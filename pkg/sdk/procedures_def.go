package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var procedureArgument = g.NewQueryStruct("ProcedureArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes()).
	Text("ArgDataType", g.KeywordOptions().NoQuotes())

var procedureArgumentType = g.NewQueryStruct("ProcedureArgumentType").
	Text("ArgDataType", g.KeywordOptions().NoQuotes())

var procedureColumn = g.NewQueryStruct("ProcedureColumn").
	Text("ColumnName", g.KeywordOptions().NoQuotes()).
	Text("ColumnDataType", g.KeywordOptions().NoQuotes())

var procedureSecret = g.NewQueryStruct("ProcedureSecret").
	Text("SecretVariableName", g.KeywordOptions().SingleQuotes()).
	Text("SecretName", g.KeywordOptions().NoQuotes())

var procedureReturns = g.NewQueryStruct("ProcedureReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			Text("ResultDataType", g.KeywordOptions()).
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

var procedureReturns2 = g.NewQueryStruct("ProcedureReturns2").
	Text("ResultDataType", g.KeywordOptions()).OptionalSQL("NOT NULL")

var procedureReturns3 = g.NewQueryStruct("ProcedureReturns3").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			Text("ResultDataType", g.KeywordOptions()),
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

var procedureSet = g.NewQueryStruct("ProcedureSet").
	OptionalTextAssignment("LOG_LEVEL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("TRACE_LEVEL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var procedureUnset = g.NewQueryStruct("ProcedureUnset").OptionalSQL("COMMENT")

var (
	procedureStrictOrNot   = g.NewQueryStruct("ProcedureStrictOrNot").OptionalSQL("STRICT").OptionalSQL("CALLED ON NULL INPUT")
	procedureVolatileOrNot = g.NewQueryStruct("ProcedureVolatileOrNot").OptionalSQL("VOLATILE").OptionalSQL("IMMUTABLE")
	procedureExecuteAs     = g.NewQueryStruct("ProcedureExecuteAs").OptionalSQL("CALLER").OptionalSQL("OWNER")
	procedureImport        = g.NewQueryStruct("ProcedureImport").Text("Import", g.KeywordOptions().SingleQuotes())
	procedurePackage       = g.NewQueryStruct("ProcedurePackage").Text("Package", g.KeywordOptions().SingleQuotes())
)

var ProceduresDef = g.NewInterface(
	"Procedures",
	"Procedure",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateProcedureForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure",
	g.NewQueryStruct("CreateProcedureForJava").
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
		OptionalQueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		SQL("LANGUAGE JAVA").
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListQueryStructField(
			"Secrets",
			procedureSecret,
			g.ParameterOptions().Parentheses().SQL("SECRETS"),
		).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"StrictOrNot",
			procedureStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			procedureVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
).CustomOperation(
	"CreateProcedureForJavaScript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure",
	g.NewQueryStruct("CreateProcedureForJavaScript").
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
		OptionalQueryStructField(
			"Returns",
			procedureReturns2,
			g.KeywordOptions().SQL("RETURNS"),
		).
		SQL("LANGUAGE JAVASCRIPT").
		OptionalQueryStructField(
			"StrictOrNot",
			procedureStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			procedureVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
).CustomOperation(
	"CreateProcedureForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure",
	g.NewQueryStruct("CreateProcedureForPython").
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
		OptionalQueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		SQL("LANGUAGE PYTHON").
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListQueryStructField(
			"Secrets",
			procedureSecret,
			g.ParameterOptions().Parentheses().SQL("SECRETS"),
		).
		OptionalQueryStructField(
			"StrictOrNot",
			procedureStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			procedureVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
).CustomOperation(
	"CreateProcedureForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure",
	g.NewQueryStruct("CreateProcedureForScala").
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
		OptionalQueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS"),
		).
		SQL("LANGUAGE SCALA").
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"StrictOrNot",
			procedureStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			procedureVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
).CustomOperation(
	"CreateProcedureForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure",
	g.NewQueryStruct("CreateProcedureForSQL").
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
		OptionalQueryStructField(
			"Returns",
			procedureReturns3,
			g.KeywordOptions().SQL("RETURNS"),
		).
		SQL("LANGUAGE SQL").
		OptionalQueryStructField(
			"StrictOrNot",
			procedureStrictOrNot,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"VolatileOrNot",
			procedureVolatileOrNot,
			g.KeywordOptions(),
		).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		TextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes()),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-procedure",
	g.NewQueryStruct("AlterProcedure").
		Alter().
		SQL("PROCEDURE").
		IfExists().
		Name().
		ListQueryStructField(
			"ArgumentTypes",
			procedureArgumentType,
			g.ParameterOptions().Parentheses().NoEquals()).
		OptionalQueryStructField(
			"Set",
			procedureSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			procedureUnset,
			g.KeywordOptions().SQL("UNSET"),
		).
		OptionalQueryStructField(
			"ExecuteAs",
			procedureExecuteAs,
			g.KeywordOptions().SQL("EXECUTE AS"),
		).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		SetTags().UnsetTags(),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-procedure",
	g.NewQueryStruct("DropProcedure").
		Drop().
		SQL("PROCEDURE").
		IfExists().
		Name().
		ListQueryStructField(
			"ArgumentTypes",
			procedureArgumentType,
			g.ParameterOptions().Parentheses().NoEquals(),
		).WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-procedures",
	g.DbStruct("procedureRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("min_num_arguments", "int").
		Field("max_num_arguments", "int").
		Field("arguments", "string").
		Field("is_table_function", "string"),
	g.PlainStruct("Procedure").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("MinNumArguments", "int").
		Field("MaxNumArguments", "int").
		Field("Arguments", "string").
		Field("IsTableFunction", "string"),
	g.NewQueryStruct("ShowProcedures").
		Show().
		SQL("PROCEDURES").
		OptionalLike().OptionalIn(),
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/describe-procedure",
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
		ListQueryStructField(
			"ArgumentTypes",
			procedureArgumentType,
			g.ParameterOptions().Parentheses().NoEquals(),
		).WithValidation(g.ValidIdentifier, "name"),
)
