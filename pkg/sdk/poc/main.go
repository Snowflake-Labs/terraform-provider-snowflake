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

var definitionMapping = map[string]*generator.Interface{
	"database_role": example.DatabaseRole,
}

func main() {
	fmt.Printf("Running generator on %s with args %#v\n", os.Getenv("GOFILE"), os.Args[1:])
	fileWithoutSuffix, _ := strings.CutSuffix(os.Getenv("GOFILE"), "_def.go")
	definition := getDefinition(fileWithoutSuffix)

	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
		o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
		setParent(o.OptsField)
	}

	runAllTemplates(os.Stdout, definition)

	runTemplateAndSave(definition, generateInterface, filenameFor(fileWithoutSuffix, ""))
	runTemplateAndSave(definition, generateDtos, filenameFor(fileWithoutSuffix, "_dto"))
	runTemplateAndSave(definition, generateImplementation, filenameFor(fileWithoutSuffix, "_impl"))
	runTemplateAndSave(definition, generateUnitTests, filename(fileWithoutSuffix, "_gen", "_test.go"))
	runTemplateAndSave(definition, generateValidations, filenameFor(fileWithoutSuffix, "_validations"))
	runTemplateAndSave(definition, generateIntegrationTests, filename(fileWithoutSuffix, "_gen_integration", "_test.go"))
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func getDefinition(fileWithoutSuffix string) *generator.Interface {
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

func runTemplateAndSave(def *generator.Interface, genFunc func(io.Writer, *generator.Interface), fileName string) {
	buffer := bytes.Buffer{}
	genFunc(&buffer, def)
	generator.WriteCodeToFile(&buffer, fileName)
}

func runAllTemplates(writer io.Writer, def *generator.Interface) {
	generateInterface(writer, def)
	generateImplementation(writer, def)
	generateUnitTests(writer, def)
	generateValidations(writer, def)
}

func generateInterface(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.InterfaceTemplate, def)
	for _, o := range def.Operations {
		generateOptionsStruct(writer, o)
	}
}

func generateOptionsStruct(writer io.Writer, operation *generator.Operation) {
	printTo(writer, generator.OptionsTemplate, operation)

	for _, f := range operation.OptsField.Fields {
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

func generateDtos(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.DtoTemplate, def)
}

func generateImplementation(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.ImplementationTemplate, def)
}

func generateUnitTests(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.TestFuncTemplate, def)
}

func generateValidations(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.ValidationsImplTemplate, def)
}

func generateIntegrationTests(writer io.Writer, def *generator.Interface) {
	generatePackageDirective(writer)
	printTo(writer, generator.IntegrationTestsTemplate, def)
}

func generatePackageDirective(writer io.Writer) {
	printTo(writer, generator.PackageTemplate, os.Getenv("GOPACKAGE"))
}

func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
