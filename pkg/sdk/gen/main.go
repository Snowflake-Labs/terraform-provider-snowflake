package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"
)

func main() {
	for _, o := range generator.DatabaseRoleInterface.Operations {
		o.ObjectInterface = &generator.DatabaseRoleInterface
	}

	runAllTemplates(os.Stdout)

	prefix := "database_role"
	runTemplateAndSave(generateInterface, filenameFor(prefix, ""))
	runTemplateAndSave(generateImplementation, filenameFor(prefix, "_impl"))
	runTemplateAndSave(generateUnitTests, filenameFor(prefix, "_test"))
	runTemplateAndSave(generateValidations, filenameFor(prefix, "_validations"))
}

func filenameFor(prefix string, part string) string {
	suffix := "_gen.go"
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
	printTo(writer, generator.InterfaceTemplate, &generator.DatabaseRoleInterface)
	for _, o := range generator.DatabaseRoleInterface.Operations {
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
	printTo(writer, generator.ImplementationTemplate, &generator.DatabaseRoleInterface)
}

func generateUnitTests(writer io.Writer) {
	printTo(writer, generator.TestFuncTemplate, &generator.DatabaseRoleInterface)
}

func generateValidations(writer io.Writer) {
	printTo(writer, generator.ValidationsImplTemplate, &generator.DatabaseRoleInterface)
}

// TODO: get rid of any
func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
