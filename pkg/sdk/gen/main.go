package main

import (
	"io"
	"os"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"
)

func main() {
	for _, o := range generator.DatabaseRoleInterface.Operations {
		o.ObjectInterface = &generator.DatabaseRoleInterface
	}

	runAllTemplates(os.Stdout)
}

func runAllTemplates(writer io.Writer) {
	printTo(writer, generator.InterfaceTemplate, &generator.DatabaseRoleInterface)

	for _, o := range generator.DatabaseRoleInterface.Operations {
		generateOptionsStruct(writer, o)
	}

	printTo(writer, generator.ImplementationTemplate, &generator.DatabaseRoleInterface)

	printTo(writer, generator.TestFuncTemplate, &generator.DatabaseRoleInterface)

	printTo(writer, generator.ValidationsImplTemplate, &generator.DatabaseRoleInterface)
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

// TODO: get rid of any
func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		panic(err)
	}
}
