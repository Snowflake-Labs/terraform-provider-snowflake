//go:build exclude

package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	uniqueTypes := make(map[string]bool)
	allStructsDetails := make([]gen.Struct, len(gen.SdkShowResultStructs))
	for idx, s := range gen.SdkShowResultStructs {
		details := gen.ExtractStructDetails(s)
		allStructsDetails[idx] = details
		printFields(details)

		for _, f := range details.Fields {
			uniqueTypes[f.ConcreteType] = true
		}
	}
	fmt.Println("===========================")
	fmt.Println("Unique types")
	fmt.Println("===========================")
	keys := maps.Keys(uniqueTypes)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Println(k)
	}

	fmt.Println("===========================")
	fmt.Println("Generated")
	fmt.Println("===========================")
	for _, details := range allStructsDetails {
		model := gen.ModelFromStructDetails(details)
		err := gen.SchemaTemplate.Execute(os.Stdout, model)
		if err != nil {
			log.Panicln(err)
		}
		err = gen.ToSchemaMapperTemplate.Execute(os.Stdout, model)
		if err != nil {
			log.Panicln(err)
		}
	}
}

func printFields(s gen.Struct) {
	fmt.Println("===========================")
	fmt.Printf("%s\n", s.Name)
	fmt.Println("===========================")

	for _, field := range s.Fields {
		fmt.Println(gen.ColumnOutput(40, field.Name, field.ConcreteType, field.UnderlyingType))
	}
	fmt.Println()
}
