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
	columnIdentity := b.QueryStruct("ColumnIdentity").
		Assignment("Start", b.Number(), b.NoEquals()).
		Assignment("Increment", b.Number(), b.NoEquals())

	columnDefaultValue := b.QueryStruct("ColumnDefaultValue").
		OneOf(
			b.Value("Expression", b.WithPointerType(b.TypeString), b.WithSQLPrefix("DEFAULT")),
			b.Value("Identity", b.WithPointerType(columnIdentity), b.WithSQLPrefix("Identity")),
		)

	maskingPolicy := b.QueryStruct("ColumnMaskingPolicy").
		SQL("MASKING POLICY").
		Identifier("Name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		List("Using", b.WithSliceType(b.TypeString))

	tableColumn := b.QueryStruct("TableColumn").
		Value("Name", b.WithType(b.TypeString)).
		Value("Type", b.WithTypeT[sdk.DataType]()).
		OptionalAssignment("Collate", b.NoEquals(), b.SingleQuotedText()).
		OptionalAssignment("Comment", b.NoEquals(), b.SingleQuotedText()).
		OptionalValue("DefailtValue", columnDefaultValue).
		OptionalSQL("NOT NULL").
		OptionalValue("MaskingPolicy", maskingPolicy).
		Tag()

	ts := b.EnumType[string]("TableScope").
		With("GlobalTableScope", "GLOBAL").
		With("LocalTableScope", "LOCAL")

	create := b.QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf(
			b.OneOf(
				b.OptionalSQL("LOCAL"),
				b.OptionalSQL("GLOBAL"),
			),
			b.OneOf(ts),
			b.OneOf(
				b.OptionalSQL("TEMP"),
				b.OptionalSQL("TEMPORARY"),
				b.OptionalSQL("VOLATILE"),
				b.OptionalSQL("TRANSIENT"),
			),
		).
		SQL("TABLE"). // can we derive fieldname from sql or vice versa ?
		IfNotExists().
		Identifier("name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		List("Columns", b.WithSliceType(tableColumn))

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
