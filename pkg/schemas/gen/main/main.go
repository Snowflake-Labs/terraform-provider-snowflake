//go:build exclude

package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

func main() {
	gencommons.NewGenerator(
		getStructDetails,
		gen.ModelFromStructDetails,
		getFilename,
		gen.AllTemplates,
	).
		WithAdditionalObjectsDebugLogs(printAllStructsFields).
		WithAdditionalObjectsDebugLogs(printUniqueTypes).
		WithObjectFilter(filterByNameFromEnv).
		RunAndHandleOsReturn()
}

func getStructDetails() []gencommons.StructDetails {
	allObjects := append(gen.SdkShowResultStructs, gen.AdditionalStructs...)
	allStructsDetails := make([]gencommons.StructDetails, len(allObjects))
	for idx, s := range allObjects {
		allStructsDetails[idx] = gencommons.ExtractStructDetails(s)
	}
	return allStructsDetails
}

func getFilename(_ gencommons.StructDetails, model gen.ShowResultSchemaModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_gen.go"
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

// TODO: describe filtering: make generate-show-output-schemas SF_TF_GENERATOR_EXT_ALLOWED_OBJECT_NAMES="sdk.Warehouse,sdk.User"
// TODO: move this filter to commons and consider extracting this as a command line param
func filterByNameFromEnv(o gencommons.StructDetails) bool {
	allowedObjectNamesString := os.Getenv("SF_TF_GENERATOR_EXT_ALLOWED_OBJECT_NAMES")
	if allowedObjectNamesString == "" {
		return true
	}
	allowedObjectNames := strings.Split(allowedObjectNamesString, ",")
	return slices.Contains(allowedObjectNames, o.ObjectName())
}
