package gencommons

import (
	"fmt"
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Run() error {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	// generating objects
	allObjects := append(gen.SdkShowResultStructs, gen.AdditionalStructs...)
	allStructsDetails := make([]StructDetails, len(allObjects))
	for idx, s := range allObjects {
		allStructsDetails[idx] = ExtractStructDetails(s)
	}

	// additional debug logs
	// printAllStructsFields(allStructsDetails)
	// printUniqueTypes(allStructsDetails)

	// printing objects to sdt out
	_ = GenerateAndPrintForAllObjects(allStructsDetails, gen.ModelFromStructDetails, gen.AllTemplates...)

	// saving objects to file
	_ = GenerateAndSaveForAllObjects(
		allStructsDetails,
		gen.ModelFromStructDetails,
		func(_ StructDetails, model gen.ShowResultSchemaModel) string {
			return ToSnakeCase(model.Name) + "_gen.go"
		},
		gen.AllTemplates...,
	)

	return nil
}
