//go:build exclude

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

	runAllTemplatesToStdOut(definition)

	runAllTemplatesAndSave(definition, fileWithoutSuffix)
}

func getDefinition(fileWithoutSuffix string) *generator.Interface {
	def, ok := definitionMapping[fileWithoutSuffix]
	if !ok {
		log.Panicf("Definition for key %s not found", os.Getenv("GOFILE"))
	}
	preprocessDefinition(def)
	return def
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *generator.Interface) {
	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
		o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
		setParent(o.OptsField)
	}
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func runAllTemplatesToStdOut(definition *generator.Interface) {
	writer := os.Stdout
	generator.GenerateInterface(writer, definition)
	generator.GenerateDtos(writer, definition)
	generator.GenerateImplementation(writer, definition)
	generator.GenerateUnitTests(writer, definition)
	generator.GenerateValidations(writer, definition)
	generator.GenerateIntegrationTests(writer, definition)
}

func runAllTemplatesAndSave(definition *generator.Interface, fileWithoutSuffix string) {
	runTemplateAndSave(definition, generator.GenerateInterface, filenameFor(fileWithoutSuffix, ""))
	runTemplateAndSave(definition, generator.GenerateDtos, filenameFor(fileWithoutSuffix, "_dto"))
	runTemplateAndSave(definition, generator.GenerateImplementation, filenameFor(fileWithoutSuffix, "_impl"))
	runTemplateAndSave(definition, generator.GenerateUnitTests, filename(fileWithoutSuffix, "_gen", "_test.go"))
	runTemplateAndSave(definition, generator.GenerateValidations, filenameFor(fileWithoutSuffix, "_validations"))
	runTemplateAndSave(definition, generator.GenerateIntegrationTests, filename(fileWithoutSuffix, "_gen_integration", "_test.go"))
}

func runTemplateAndSave(def *generator.Interface, genFunc func(io.Writer, *generator.Interface), fileName string) {
	buffer := bytes.Buffer{}
	genFunc(&buffer, def)
	generator.WriteCodeToFile(&buffer, fileName)
}

func filenameFor(prefix string, part string) string {
	return filename(prefix, part, "_gen.go")
}

func filename(prefix string, part string, suffix string) string {
	return fmt.Sprintf("%s%s%s", prefix, part, suffix)
}
