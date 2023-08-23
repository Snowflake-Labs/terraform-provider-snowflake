package main

import (
	"os"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"
)

func main() {
	for _, o := range generator.DatabaseRoleInterface.Operations {
		o.ObjectInterface = &generator.DatabaseRoleInterface
	}

	printToStdOut(generator.InterfaceTemplate, &generator.DatabaseRoleInterface)

	for _, o := range generator.DatabaseRoleInterface.Operations {
		generateOptionsStruct(o)
	}

	printToStdOut(generator.ImplementationTemplate, &generator.DatabaseRoleInterface)

	printToStdOut(generator.TestFuncTemplate, &generator.DatabaseRoleInterface)

	printToStdOut(generator.ValidationsImplTemplate, &generator.DatabaseRoleInterface)
}

func generateOptionsStruct(operation *generator.Operation) {
	printToStdOut(generator.OptionsTemplate, operation)

	for _, f := range operation.OptsStructFields {
		if len(f.Fields) > 0 {
			generateStruct(f)
		}
	}
}

func generateStruct(field *generator.Field) {
	printToStdOut(generator.StructTemplate, field)

	for _, f := range field.Fields {
		if len(f.Fields) > 0 {
			generateStruct(f)
		}
	}
}

// TODO: get rid of any
func printToStdOut(template *template.Template, model any) {
	err := template.Execute(os.Stdout, model)
	if err != nil {
		panic(err)
	}
}
