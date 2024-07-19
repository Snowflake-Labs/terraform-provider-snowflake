//go:build exclude

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	allObjects := append(gen.SdkShowResultStructs, gen.AdditionalStructs...)
	allStructsDetails := make([]gencommons.Struct, len(allObjects))
	for idx, s := range allObjects {
		allStructsDetails[idx] = gencommons.ExtractStructDetails(s)
	}

	printAllStructsFields(allStructsDetails)
	printUniqueTypes(allStructsDetails)
	generateAllStructsToStdOut(allStructsDetails)
	saveAllGeneratedSchemas(allStructsDetails)
}

func printAllStructsFields(allStructs []gencommons.Struct) {
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

func printUniqueTypes(allStructs []gencommons.Struct) {
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

func generateAllStructsToStdOut(allStructs []gencommons.Struct) {
	for _, s := range allStructs {
		fmt.Println("===========================")
		fmt.Printf("Generated for %s\n", s.Name)
		fmt.Println("===========================")
		model := gen.ModelFromStructDetails(s)
		gen.Generate(model, os.Stdout)
	}
}

func saveAllGeneratedSchemas(allStructs []gencommons.Struct) {
	for _, s := range allStructs {
		buffer := bytes.Buffer{}
		model := gen.ModelFromStructDetails(s)
		gen.Generate(model, &buffer)
		filename := gencommons.ToSnakeCase(model.Name) + "_gen.go"
		writeCodeToFile(&buffer, filename)
	}
}

// TODO [SNOW-1501905]: this is copied, extract some generator helpers
func writeCodeToFile(buffer *bytes.Buffer, fileName string) {
	wd, errWd := os.Getwd()
	if errWd != nil {
		log.Panicln(errWd)
	}
	outputPath := filepath.Join(wd, fileName)
	src, errSrcFormat := format.Source(buffer.Bytes())
	if errSrcFormat != nil {
		log.Panicln(errSrcFormat)
	}
	if err := os.WriteFile(outputPath, src, 0o600); err != nil {
		log.Panicln(err)
	}
}
