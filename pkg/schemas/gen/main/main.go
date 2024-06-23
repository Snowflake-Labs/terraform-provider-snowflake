//go:build exclude

package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	allStructsDetails := make([]gen.Struct, len(gen.SdkShowResultStructs))
	for idx, s := range gen.SdkShowResultStructs {
		allStructsDetails[idx] = gen.ExtractStructDetails(s)
	}

	printAllStructsFields(allStructsDetails)
	printUniqueTypes(allStructsDetails)
	generateAllStructsToStdOut(allStructsDetails)
}

func printAllStructsFields(allStructs []gen.Struct) {
	for _, s := range allStructs {
		fmt.Println("===========================")
		fmt.Printf("%s\n", s.Name)
		fmt.Println("===========================")
		for _, field := range s.Fields {
			fmt.Println(gen.ColumnOutput(40, field.Name, field.ConcreteType, field.UnderlyingType))
		}
		fmt.Println()
	}
}

func printUniqueTypes(allStructs []gen.Struct) {
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

func generateAllStructsToStdOut(allStructs []gen.Struct) {
	for _, s := range allStructs {
		fmt.Println("===========================")
		fmt.Printf("Generated for %s\n", s.Name)
		fmt.Println("===========================")
		gen.Generate(s, os.Stdout)
	}
}
