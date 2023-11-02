package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

// TODO: define
var viewDbRow = g.DbStruct("viewDBRow").
	Field("name", "string")

// TODO: define
var view = g.PlainStruct("View").
	Field("Name", "string")

var viewColumn = g.NewQueryStruct("ViewColumn").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().NoEquals())

var viewColumnMaskingPolicy = g.NewQueryStruct("ViewColumnMaskingPolicy").
	Text("Name", g.KeywordOptions().Required()).
	Identifier("MaskingPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	List("USING", g.KindOfT[string](), nil). // TODO: double quotes here?
	WithTags()

var viewRowAccessPolicy = g.NewQueryStruct("ViewRowAccessPolicy").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	List("ON", g.KindOfT[string](), nil). // TODO: double quotes here?
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var ViewsDef = g.NewInterface(
	"Views",
	"View",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-view",
		g.NewQueryStruct("CreateView").
			Create().
			OrReplace().
			OptionalSQL("SECURE").
			// There are multiple variants in the docs: { [ { LOCAL | GLOBAL } ] TEMP | TEMPORARY | VOLATILE }
			// but from description they are all the same. For the sake of simplicity only one option is used here.
			OptionalSQL("TEMPORARY").
			OptionalSQL("RECURSIVE").
			SQL("VIEW").
			IfNotExists().
			Name().
			ListQueryStructField("Columns", viewColumn, g.ListOptions().Parentheses()).
			ListQueryStructField("ColumnsMaskingPolicies", viewColumnMaskingPolicy, g.ListOptions().NoParentheses().NoEquals()).
			OptionalSQL("COPY GRANTS").
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			// In the current docs ROW ACCESS POLICY and TAG are specified twice.
			// It is a mistake probably so here they are present only once.
			OptionalQueryStructField("RowAccessPolicy", viewRowAccessPolicy, g.KeywordOptions()).
			WithTags().
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-view",
		g.NewQueryStruct("AlterView").
			Alter().
			SQL("VIEW").
			IfExists().
			Name().
			OptionalSQL("RESUME").
			SetTags().
			UnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Resume", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-view",
		g.NewQueryStruct("DropView").
			Drop().
			SQL("VIEW").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-views",
		viewDbRow,
		view,
		g.NewQueryStruct("ShowViews").
			Show().
			Terse().
			SQL("VIEWS").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimit(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-view",
		viewDbRow,
		view,
		g.NewQueryStruct("DescribeView").
			Describe().
			SQL("VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
