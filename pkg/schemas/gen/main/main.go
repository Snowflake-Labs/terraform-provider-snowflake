//go:build exclude

package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	allObjects := append(gen.SdkShowResultStructs, gen.AdditionalStructs...)
	allStructsDetails := make([]gencommons.StructDetails, len(allObjects))
	for idx, s := range allObjects {
		allStructsDetails[idx] = gencommons.ExtractStructDetails(s)
	}

	printAllStructsFields(allStructsDetails)
	printUniqueTypes(allStructsDetails)
	// TODO: handle objects
	gencommons.GenerateAndPrintForAllObjects(allStructsDetails, gen.ModelFromStructDetails, gen.AllTemplates)
	gencommons.GenerateAndSaveForAllObjects(
		allStructsDetails,
		gen.ModelFromStructDetails,
		func(_, model gen.ShowResultSchemaModel) { return gencommons.ToSnakeCase(model.Name) + "_gen.go" },
		gen.AllTemplates,
	)
}

func printAllStructsFields(allStructs []gencommons.StructDetails) {
	for _, s := range allStructs {
		fmt.Println("===========================")
		fmt.Printf("%s\n", s.Name)
		fmt.Println("===========================")
		for _, field := range s.Fields {
			fmt.Println(gencommons.ColumnOutput(40, field.Name, field.ConcreteType, field.UnderlyingType))
		}
		fmt.Println()
	}
}

func printUniqueTypes(allStructs []gencommons.StructDetails) {
	uniqueTypes := make(map[string]bool)
	for _, s := range allStructs {
		for _, f := range s.Fields {
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
}
