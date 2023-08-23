//go:build exclude

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/example"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

var definitionMapping = map[string]generator.Interface{
	"database_role": example.DatabaseRole,
}

func main() {
	fmt.Printf("Running generator on %s with args %#v\n", os.Getenv("GOFILE"), os.Args[1:])
	fileWithoutSuffix, _ := strings.CutSuffix(os.Getenv("GOFILE"), "_def.go")
	definition := getDefinition(fileWithoutSuffix)

	for _, o := range definition.Operations {
		o.ObjectInterface = &definition
	}

	runAllTemplates(os.Stdout)

	runTemplateAndSave(generateInterface, filenameFor(fileWithoutSuffix, ""))
	runTemplateAndSave(generateImplementation, filenameFor(fileWithoutSuffix, "_impl"))
	runTemplateAndSave(generateUnitTests, filename(fileWithoutSuffix, "_gen", "_test.go"))
	runTemplateAndSave(generateValidations, filenameFor(fileWithoutSuffix, "_validations"))
}

func getDefinition(fileWithoutSuffix string) generator.Interface {
	def, ok := definitionMapping[fileWithoutSuffix]
	if !ok {
		log.Panicf("Definition for key %s not found", os.Getenv("GOFILE"))
	}
	return def
}

func filenameFor(prefix string, part string) string {
	return filename(prefix, part, "_gen.go")
}

func filename(prefix string, part string, suffix string) string {
	return fmt.Sprintf("%s%s%s", prefix, part, suffix)
}

func runTemplateAndSave(genFunc func(io.Writer), fileName string) {
	buffer := bytes.Buffer{}
	genFunc(&buffer)
	generator.WriteCodeToFile(&buffer, fileName)
}

func runAllTemplates(writer io.Writer) {
	generateInterface(writer)
	generateImplementation(writer)
	generateUnitTests(writer)
	generateValidations(writer)
}

func generateInterface(writer io.Writer) {
	generatePackageDirective(writer)
	printTo(writer, generator.InterfaceTemplate, &example.DatabaseRole)
	for _, o := range example.DatabaseRole.Operations {
		generateOptionsStruct(writer, o)
	}
}

func generateOptionsStruct(writer io.Writer, operation *generator.Operation) {
	printTo(writer, generator.OptionsTemplate, operation)

	for _, f := range operation.OptsStructFields {
		if len(f.Fields) > 0 {
			generateStruct(writer, f)
		}
	}
}

func generateStruct(writer io.Writer, field *generator.Field) {
	printTo(writer, generator.StructTemplate, field)

	for _, f := range field.Fields {
		if len(f.Fields) > 0 {
			generateStruct(writer, f)
		}
	}
}

func generateImplementation(writer io.Writer) {
	generatePackageDirective(writer)
	printTo(writer, generator.ImplementationTemplate, &example.DatabaseRole)
}

func generateUnitTests(writer io.Writer) {
	generatePackageDirective(writer)
	printTo(writer, generator.TestFuncTemplate, &example.DatabaseRole)
}

func generateValidations(writer io.Writer) {
	generatePackageDirective(writer)
	printTo(writer, generator.ValidationsImplTemplate, &example.DatabaseRole)
}

func generatePackageDirective(writer io.Writer) {
	printTo(writer, generator.PackageTemplate, os.Getenv("GOPACKAGE"))
}

// TODO: get rid of any
func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
