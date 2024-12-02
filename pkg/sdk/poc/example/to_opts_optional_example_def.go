package example

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ../main.go

var ToOptsOptionalExample = g.NewInterface(
	"ToOptsOptionalExamples",
	"ToOptsOptionalExample",
	g.KindOfT[DatabaseObjectIdentifier](),
).AlterOperation("https://example.com",
	g.NewQueryStruct("Alter").
		Alter().
		IfExists().
		Name().
		OptionalQueryStructField(
			"OptionalField",
			g.NewQueryStruct("OptionalField").
				List("SomeList", "DatabaseObjectIdentifier", g.ListOptions()),
			g.KeywordOptions(),
		).
		QueryStructField(
			"RequiredField",
			g.NewQueryStruct("RequiredField").
				List("SomeRequiredList", "DatabaseObjectIdentifier", g.ListOptions()),
			g.KeywordOptions(),
		),
)
