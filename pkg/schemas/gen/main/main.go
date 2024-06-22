//go:build exclude

package main

import (
	"fmt"
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
)

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	for _, s := range gen.SdkShowResultStructs {
		printFields(s)
	}
}

func printFields(s any) {
	structDetails := gen.ExtractStructDetails(s)

	fmt.Println("===========================")
	fmt.Printf("%s\n", structDetails.Name)
	fmt.Println("===========================")

	for _, field := range structDetails.Fields {
		fmt.Println(gen.ColumnOutput(40, field.Name, field.ConcreteType, field.UnderlyingType))
	}
	fmt.Println()
}
