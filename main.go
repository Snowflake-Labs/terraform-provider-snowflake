package main

import (
	"flag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const ProviderAddr = "registry.terraform.io/Snowflake-Labs/snowflake"

func main() {
	debug := flag.Bool("debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        *debug,
		ProviderAddr: ProviderAddr,
		ProviderFunc: provider.Provider,
	})
}

// createTable := 
// 	b.QueryStruct("CreateTableOpts").
// 	create().
// 	Field("OrReplace").Sql("OR REPLACE").Optional().Type(T[string]).End().
// 	Field("ClusterBy").Sql("CLUSTER BY").Optional().ListWithParen().Type(T[string]).End().
// 	Field("ClusterBy").Sql("CLUSTER BY").Optional().ListWithParen()
// 	Field("Type").Optional.NoSql.Equals().Text.End
// 	Field("Type").Required().Type(T[sdk.DataType]).End().
// 	Field("table").SQL("TABLE").Required().End().
// 	Field("IfNotExists").Sql("IF NOT EXISTS").Optional().End().
// 	Field("TableName").Required().Type(T[string]).End().
// 	Field("LeftParen").Sql("(").Required().End().
// 	Field("Columns").List().Type(T[tableColumn.Type]).End().
// 	Field("OutOfLineConstraint").Optional().Type(T[outOfLineConstraint.Type]).End().
// 	Field("RightParen").Sql(")").Required().End().
// 	Field("ClusterBy").Sql("CLUSTER BY").Optional().ListWithParen().Type(T[string]).End().
// 	Field("EnableSchemaEvolution").Sql("ENABLE SCHEMA EVOLUTION").Optional().Equals().Type(T[bool]).End().
// 	Field("Something with quotes").Sql("FOO").Optional().Equals().Type(T[string]).SingleQuotes().End().
//
// 	Field("something").Optional().Sql("FOO").Value(Equals().T[string]().SingleQuotes)
//
//
// 	Field("ClusterBy").Optional().Sql("CLUSTER BY").Value(b.Equals().List(enum.Type).Parens().Commas()).
// 	Field("EnableSchemaEvolution").Optional().Sql("ENABLE...").Value(b.Var(T[bool]).Parens()]).End().
//
// 	Field("RowAccessPolicy").Optional().Sql("ROW ACCESS POLICY").Value(b.Identifier).Sql("ON").Value(b.List(T[string]).Parens().Commas()).End().
//
// 	Field("RowAccessPolicy").Optional().Sql("ROW ACCESS POLICY").Value(rowAccessPolicy.Type).End().
// 	OneOf(
// 		
// 	)
//
//
// 	Field("ClusterBy").Sql("CLUSTER BY").Optional().ListWithParen().Type(T[string]).End().
// 	Build()
