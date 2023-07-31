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
		Identifier("Name", b.WithTypeT[sdk.AccountObjectIdentifier]()).
		Parameter("Value", b.WithType(b.TypeString), b.WithSQL("single_quotes"))

	// ...
)

func main() {
	//columnIdentity := b.QueryStruct("ColumnIdentity").
	//	Number("Start", b.SQLPrefix("START")).
	//	Number("Increment", b.SQLPrefix("INCREMENT")).
	//
	//columnDefaultValue := b.QueryStruct("ColumnDefaultValue").
	//	OneOf(
	//		b.OptionalText("Expression", b.SQLPrefix("DEFAULT")),
	//		b.OptionalValue("Identity", columnIdentity, b.SQLPrefix("IDENTITY")),
	//	)

	optsExample := b.QueryStruct("🥸").
		Static2("", b.StaticOpts().Quotes()).
		Keyword2(b.KeywordOpts().SQL("qafsd")).
		Parameter2(b.ParameterOpts().Paren(true).Equals(false))

	_ = optsExample

	maskingPolicy := b.QueryStruct("ColumnMaskingPolicy").
		SQL("MASKING POLICY").
		Identifier("Name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		List("Using", b.TypeString)

	tableColumn := b.QueryStruct("TableColumn").
		Text("Name").
		Value("Type", b.WithTypeT[sdk.DataType]()).
		OptionalText("Collate", b.SQLPrefix("COLLATE"), b.SingleQuotedText()).
		OptionalText("Comment", b.SQLPrefix("COMMENT"), b.SingleQuotedText()).
		OptionalValue("DefaultValue", columnDefaultValue).
		OptionalSQL("NOT NULL").
		OptionalValue("MaskingPolicy", b.Link(maskingPolicy)).
		Tag()

	ts := b.EnumType[string]("TableScope").
		With("GlobalTableScope", "GLOBAL").
		With("LocalTableScope", "LOCAL")

	create := b.QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf("field",
			b.OneOf("scope",
				b.OptionalSQL("LOCAL"),
				b.OptionalSQL("GLOBAL"),
			),
			b.OneOf("scopeV2", ts),
			b.OneOf("fieldType",
				b.OptionalSQL("TEMP"),
				b.OptionalSQL("TEMPORARY"),
				b.OptionalSQL("VOLATILE"),
				b.OptionalSQL("TRANSIENT"),
			),
		).
		SQL("TABLE"). // can we derive field name from sql or vice versa ?
		IfNotExists().
		Identifier("name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		List("Columns", tableColumn).
		Assignment()

	generator.GenerateAll(
		b.API{},
		ts,
		tagAssociation,
		columnIdentity,
		columnDefaultValue,
		maskingPolicy,
		tableColumn,
		create,
	)
}
