package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var materializedViewColumn = g.NewQueryStruct("MaterializedViewColumn").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().NoEquals())

var materializedViewColumnMaskingPolicy = g.NewQueryStruct("MaterializedViewColumnMaskingPolicy").
	Text("Name", g.KeywordOptions().Required()).
	Identifier("MaskingPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	NamedListWithParens("USING", g.KindOfT[string](), nil). // TODO: double quotes here?
	OptionalTags()

var materializedViewRowAccessPolicy = g.NewQueryStruct("MaterializedViewRowAccessPolicy").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	NamedListWithParens("ON", g.KindOfT[string](), g.KeywordOptions().Required()). // TODO: double quotes here?
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var materializedViewSet = g.NewQueryStruct("MaterializedViewSet").
	OptionalSQL("SECURE").
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Secure", "Comment")

var materializedViewUnset = g.NewQueryStruct("MaterializedViewUnset").
	OptionalSQL("SECURE").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Secure", "Comment")

var MaterializedViewsDef = g.NewInterface(
	"MaterializedViews",
	"MaterializedView",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-materialized-view",
		g.NewQueryStruct("CreateMaterializedView").
			Create().
			OrReplace().
			OptionalSQL("SECURE").
			SQL("MATERIALIZED VIEW").
			IfNotExists().
			Name().
			OptionalCopyGrants().
			ListQueryStructField("Columns", materializedViewColumn, g.ListOptions().Parentheses()).
			ListQueryStructField("ColumnsMaskingPolicies", materializedViewColumnMaskingPolicy, g.ListOptions().NoParentheses().NoEquals()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			// In the current docs ROW ACCESS POLICY and TAG are specified twice.
			// It is a mistake probably so here they are present only once.
			OptionalQueryStructField("RowAccessPolicy", materializedViewRowAccessPolicy, g.KeywordOptions()).
			OptionalTags().
			NamedListWithParens("CLUSTER BY", g.KindOfT[string](), g.KeywordOptions()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view",
		g.NewQueryStruct("AlterMaterializedView").
			Alter().
			SQL("MATERIALIZED VIEW").
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			NamedListWithParens("CLUSTER BY", g.KindOfT[string](), g.KeywordOptions()).
			OptionalSQL("DROP CLUSTERING KEY").
			OptionalSQL("SUSPEND RECLUSTER").
			OptionalSQL("RESUME RECLUSTER").
			OptionalSQL("SUSPEND").
			OptionalSQL("RESUME").
			OptionalQueryStructField("Set", materializedViewSet, g.KeywordOptions().SQL("SET")).
			OptionalQueryStructField("Unset", materializedViewUnset, g.KeywordOptions().SQL("UNSET")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"),
	)
