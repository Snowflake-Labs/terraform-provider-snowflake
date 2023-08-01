//go:build exclude

package main

import (
	b "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/gen/builder"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/gen/generator"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

/*
   CREATE [ OR REPLACE ]
     [ { [ { LOCAL | GLOBAL } ] TEMP | TEMPORARY | VOLATILE | TRANSIENT } ]
     TABLE [ IF NOT EXISTS ] <table_name>
     ( <col_name> <col_type>
       [ COLLATE '<collation_specification>' ]
       [ COMMENT '<string_literal>' ]
       [ { DEFAULT <expr>
         | { AUTOINCREMENT | IDENTITY } [ { ( <start_num> , <step_num> ) | START <num> INCREMENT <num> } ] } ]
       [ NOT NULL ]
*/

// Common types
var (
	tagAssociation = b.QueryStruct("TagAssociation").
		Identifier("Name", b.TypeT[sdk.AccountObjectIdentifier]()).
		AssignText("Value", b.ParameterOptions().SingleQuotes())

	// ...
)

func main() {
	columnIdentity := b.QueryStruct("ColumnIdentity").
		AssignNumber("Start", b.ParameterOptions().NoQuotes().NoEquals()).
		AssignNumber("Increment", b.ParameterOptions().NoQuotes().NoEquals())

	columnDefaultValue := b.QueryStruct("ColumnDefaultValue").
		OneOf(
			b.OptionalText("Expression", b.KeywordOptions().SQLPrefix("DEFAULT")),
			b.OptionalValue("Identity", b.TypeOfQueryStruct(columnIdentity), b.KeywordOptions().SQLPrefix("IDENTITY")),
		)

	maskingPolicy := b.QueryStruct("ColumnMaskingPolicy").
		SQL("MASKING POLICY").
		Identifier("Name", b.TypeT[sdk.SchemaObjectIdentifier]()).
		List("Using", b.TypeString, nil)

	tableColumn := b.QueryStruct("TableColumn").
		Text("Name", nil).
		Value("Typer", b.TypeT[sdk.DataType](), nil).
		OptionalText("Collate", b.KeywordOptions().SQLPrefix("COLLATE").SingleQuotes()).
		OptionalText("Comment", b.KeywordOptions().SQLPrefix("COMMENT").SingleQuotes()).
		OptionalValue("DefaultValue", b.TypeOfQueryStruct(columnDefaultValue), nil).
		OptionalSQL("NOT NULL", nil).
		OptionalValue("MaskingPolicy", b.TypeOfQueryStruct(maskingPolicy), nil).
		List("Tag", b.TypeOfQueryStruct(tagAssociation), b.ListOptions().Parentheses())

	ts := b.EnumType[string]("TableScope").
		With("GlobalTableScope", "GLOBAL").
		With("LocalTableScope", "LOCAL")

	kind := b.EnumType[string]("TableKind").
		With("TableKindTemp", "TEMP").
		With("LocalTableTemporary", "TEMPORARY").
		With("LocalTableVolatile", "VOLATILE").
		With("LocalTableTransient", "TRANSIENT")

	create := b.QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf(
			b.OneOf("Scope", ts),
			b.OneOf("Kind", kind),
		).
		SQL("TABLE").
		IfNotExists().
		Identifier("name", b.TypeT[*sdk.SchemaObjectIdentifier]()).
		List("Columns", b.TypeOfQueryStruct(tableColumn), nil)

	generator.GenerateAll(
		ts,
		kind,
		tagAssociation,
		columnIdentity,
		columnDefaultValue,
		maskingPolicy,
		tableColumn,
		create,
	)
}
