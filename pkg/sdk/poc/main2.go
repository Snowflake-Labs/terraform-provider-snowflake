//go:build exclude

package main

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/example2"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator2"
	"os"
)

var definitionMapping = map[string]*generator2.Interface{
	"database_role_def.go": example2.DatabaseRole,
}

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])
	if def, ok := definitionMapping[file]; ok {
		runAllTemplatesToStdOut(def)
		//runAllTemplatesAndSave(def, file)
	} else {
		panic("missing definition in definitionMapping")
	}
}

func runAllTemplatesToStdOut(definition *generator2.Interface) {
	writer := os.Stdout
	generator2.GenerateInterface(writer, definition)
	generator2.GenerateDtos(writer, definition)
	generator2.GenerateImplementation(writer, definition)
	//generator2.GenerateUnitTests(writer, definition)
	generator2.GenerateValidations(writer, definition)
	//generator2.GenerateIntegrationTests(writer, definition)
}

//func runAllTemplatesAndSave(definition *generator.Interface, file string) {
//	fileWithoutSuffix, _ := strings.CutSuffix(file, "_def.go")
//	runTemplateAndSave(definition, generator.GenerateInterface, filenameFor(fileWithoutSuffix, ""))
//	runTemplateAndSave(definition, generator.GenerateDtos, filenameFor(fileWithoutSuffix, "_dto"))
//	runTemplateAndSave(definition, generator.GenerateImplementation, filenameFor(fileWithoutSuffix, "_impl"))
//	runTemplateAndSave(definition, generator.GenerateUnitTests, filename(fileWithoutSuffix, "_gen", "_test.go"))
//	runTemplateAndSave(definition, generator.GenerateValidations, filenameFor(fileWithoutSuffix, "_validations"))
//	runTemplateAndSave(definition, generator.GenerateIntegrationTests, filename(fileWithoutSuffix, "_gen_integration", "_test.go"))
//}
//
//func runTemplateAndSave(def *generator.Interface, genFunc func(io.Writer, *generator.Interface), fileName string) {
//	buffer := bytes.Buffer{}
//	genFunc(&buffer, def)
//	generator.WriteCodeToFile(&buffer, fileName)
//}
//
//func filenameFor(prefix string, part string) string {
//	return filename(prefix, part, "_gen.go")
//}
//
//func filename(prefix string, part string, suffix string) string {
//	return fmt.Sprintf("%s%s%s", prefix, part, suffix)
//}
