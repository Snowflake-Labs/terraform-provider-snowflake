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
		Parameter("Start", b.WithType(b.TypeInt), b.WithNoQuotes(), b.WithNoEquals(), b.WithSQL("START")).
		Parameter("Increment", b.WithType(b.TypeInt), b.WithNoQuotes(), b.WithNoEquals(), b.WithSQL("INCREMENT"))

	columnDefaultValue := b.QueryStruct("ColumnDefaultValue").
		OneOf(
			b.Keyword("Expression", b.WithPointerType(b.TypeString), b.WithSQL("DEFAULT")),
			b.Keyword("Identify", b.WithPointerType(columnIdentity), b.WithSQL("IDENTITY")),
		)

	maskingPolicy := b.QueryStruct("ColumnMaskingPolicy").
		Static("maskingPolicy", b.WithType(b.TypeBool), b.WithSQL("MASKING POLICY")).
		Identifier("Name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		Keyword("Using", b.WithSliceType(b.TypeString))

	tableColumn := b.QueryStruct("TableColumn").
		Keyword("Name", b.WithType(b.TypeString)).
		Keyword("Type", b.WithTypeT[sdk.DataType]()).
		Keyword("").Equals().String().End().
		Parameter("Collate", b.WithPointerType(b.TypeString), b.WithEquals()).
		Parameter("Comment", b.WithPointerType(b.TypeString)).
		// Parameter("Comment", b.Optional(), b.Text())
		// Parameter("UnsetTag", b.OptionalList(), b.Text())
		Keyword("DefailtValue", b.WithPointerType(columnDefaultValue)).
		Keyword("NotNull", b.WithTypeBoolPtr(), b.WithSQL("NOT NULL")).
		Keyword("MaskingPolicy", b.WithPointerType(maskingPolicy)).
		List("Tag", b.WithType(tagAssociation), b.WithDDL("parentheses"))

	ts := b.EnumType("TableScope").
		With("GlobalTableScope", "GLOBAL").
		With("LocalTableScope", "LOCAL")

	create := b.QueryStruct("CreateTableOptions").
		Create().
		OrReplace().
		OneOf(
			// first approach
			b.OneOf(
				b.Keyword("Local", b.WithType(b.TypeBoolPtr), b.WithSQL("LOCAL")),
				b.Keyword("Glocal", b.WithType(b.TypeBoolPtr), b.WithSQL("GLOBAL")),
			),
			// second approach
			b.OneOf(ts),
			b.OneOf(
				b.Keyword("Temp", b.WithType(b.TypeBoolPtr), b.WithSQL("TEMP")),
				b.Keyword("Temporary", b.WithType(b.TypeBoolPtr), b.WithSQL("TEMPORARY")),
				b.Keyword("Volatile", b.WithType(b.TypeBoolPtr), b.WithSQL("VOLATILE")),
				b.Keyword("Transient", b.WithType(b.TypeBoolPtr), b.WithSQL("TRANSIENT")),
			),
		).
		// should we explicitly write types or can we assume that the default is bool and if want something else then we should specify the type ?
		Static("TABLE", b.WithType(b.TypeBool)).
		IfNotExists().
		Identifier("name", b.WithTypeT[sdk.SchemaObjectIdentifier]()).
		Keyword("Columns", b.WithSliceType(tableColumn))

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
