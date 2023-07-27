//go:build exclude

package main

import (
	b "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/gen/builder"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/gen/generator"
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

type DataType string

type SchemaObjectIdentifier string

// Common types
var (
	tagAssociation = b.QueryStruct("TagAssociation").
		Identifier("Name").
		Parameter("Value", b.WithSQL("single_quotes"))

	// ...
)

func main() {
	columnIdentity := b.QueryStruct("ColumnIdentity").
		Parameter("Start", b.WithType(b.TypeInt), b.WithNoQuotes(), b.WithNoEquals(), b.WithSQL("START")).
		Parameter("Increment", b.WithType(b.TypeInt), b.WithNoQuotes(), b.WithNoEquals(), b.WithSQL("INCREMENT"))

	columnDefaultValue := b.QueryStruct("ColumnDefaultValue").
		OneOf(
			b.Keyword("Expression", b.WithPointerType(b.TypeString), b.WithSQL("DEFAULT")),
			b.Keyword("Identify", b.WithPointerType(columnIdentity), b.WithSQL("IDENTITY")),
		)

	maskingPolicy := b.QueryStruct("ColumnMaskingPolicy").
		Static("maskingPolicy", b.WithType(b.TypeBool), b.WithSQL("MASKING POLICY")).
		Identifier("Name", b.WithTypeT[SchemaObjectIdentifier]()).
		Keyword("Using", b.WithSliceType(b.TypeString))

	tableColumn := b.QueryStruct("TableColumn").
		Keyword("Name", b.WithType(b.TypeString)).
		Keyword("Type", b.WithTypeT[DataType]()).
		Parameter("Collate", b.WithPointerType(b.TypeString)).
		Parameter("Comment", b.WithPointerType(b.TypeString)).
		Keyword("DefailtValue", b.WithPointerType(columnDefaultValue)).
		Keyword("NotNull", b.WithTypeBoolPtr(), b.WithSQL("NOT NULL")).
		Keyword("MaskingPolicy", b.WithPointerType(maskingPolicy)).
		List("Tag", b.WithType(tagAssociation), b.WithDDL("parentheses"))

	create := b.QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf(
			b.OneOf(
				b.Keyword("Local", b.WithSQL("LOCAL")),
				b.Keyword("Glocal", b.WithSQL("GLOBAL")),
			),
			b.OneOf(
				b.Keyword("Temp", b.WithSQL("TEMP")),
				b.Keyword("Temporary", b.WithSQL("TEMPORARY")),
				b.Keyword("Volatile", b.WithSQL("VOLATILE")),
				b.Keyword("Transient", b.WithSQL("TRANSIENT")),
			),
		).
		Static("TABLE").
		IfNotExists().
		Identifier("name").
		Keyword("Columns", b.WithSliceType(tableColumn))

	generator.GenerateAll(
		b.API{},
		tagAssociation,
		columnIdentity,
		columnDefaultValue,
		maskingPolicy,
		tableColumn,
		create,
	)
}
