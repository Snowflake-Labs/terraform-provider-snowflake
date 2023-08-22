package main

import (
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"
)

func main() {
	for _, o := range generator.DatabaseRoleInterface.Operations {
		o.ObjectInterface = &generator.DatabaseRoleInterface
	}

	err := generator.InterfaceTemplate.Execute(os.Stdout, generator.DatabaseRoleInterface)
	if err != nil {
		panic(err)
	}

	for _, o := range generator.DatabaseRoleInterface.Operations {
		err := generator.OptionsTemplate.Execute(os.Stdout, o)
		if err != nil {
			panic(err)
		}
	}
}
