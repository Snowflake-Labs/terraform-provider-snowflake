package main

import (
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"
)

func main() {
	for _, o := range generator.DatabaseRoleInterface.Operations {
		o.ObjectInterface = &generator.DatabaseRoleInterface
	}

	err := generator.InterfaceTemplate.Execute(os.Stdout, &generator.DatabaseRoleInterface)
	if err != nil {
		panic(err)
	}

	for _, o := range generator.DatabaseRoleInterface.Operations {
		generateOptionsStruct(o)
	}

	err = generator.ImplementationTemplate.Execute(os.Stdout, &generator.DatabaseRoleInterface)
	if err != nil {
		panic(err)
	}

	for _, o := range generator.DatabaseRoleInterface.Operations {
		generateTestFunc(o)
	}
}

func generateOptionsStruct(operation *generator.Operation) {
	err := generator.OptionsTemplate.Execute(os.Stdout, operation)
	if err != nil {
		panic(err)
	}
	for _, f := range operation.OptsStructFields {
		if len(f.Fields) > 0 {
			generateStruct(f)
		}
	}
}

func generateStruct(field *generator.Field) {
	err := generator.StructTemplate.Execute(os.Stdout, field)
	if err != nil {
		panic(err)
	}
	for _, f := range field.Fields {
		if len(f.Fields) > 0 {
			generateStruct(f)
		}
	}
}

func generateTestFunc(operation *generator.Operation) {
	err := generator.TestFuncTemplate.Execute(os.Stdout, operation)
	if err != nil {
		panic(err)
	}
}
