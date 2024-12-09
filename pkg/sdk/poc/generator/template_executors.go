package generator

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"text/template"
)

var (
	generatedStructs []string
	generatedDtos    []string
)

func GenerateInterface(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, InterfaceTemplate, def)
	for _, o := range def.Operations {
		if o.OptsField != nil {
			generateOptionsStruct(writer, o)
		}
	}
	for _, m := range def.HelperMethods {
		if m != nil {
			generateHelperMethods(writer, m)
		}
	}
}

func generateHelperMethods(writer io.Writer, hm *HelperMethod) {
	printTo(writer, HelperMethodTemplate, hm)
}

func generateOptionsStruct(writer io.Writer, operation *Operation) {
	printTo(writer, OperationStructTemplate, operation)

	for _, f := range operation.HelperStructs {
		if !slices.Contains(generatedStructs, f.KindNoPtr()) {
			generateStruct(writer, f)
		}
	}

	for _, f := range operation.OptsField.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
			generateStruct(writer, f)
		}
	}
}

func generateStruct(writer io.Writer, field *Field) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		fmt.Println("Generating: " + field.KindNoPtr())
		printTo(writer, StructTemplate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	}

	for _, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			generateStruct(writer, f)
		}
	}
}

func GenerateDtos(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, DtoTemplate, def)
	for _, o := range def.Operations {
		if o.OptsField != nil {
			generateDtoDecls(writer, o.OptsField)
		}
	}
}

func generateDtoDecls(writer io.Writer, field *Field) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		printTo(writer, DtoDeclTemplate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for _, f := range field.Fields {
			if f.IsStruct() {
				generateDtoDecls(writer, f)
			}
		}
	}
}

func GenerateImplementation(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ImplementationTemplate, def)
}

func GenerateUnitTests(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, UnitTestsTemplate, def)
}

func GenerateValidations(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ValidationsTemplate, def)
}

func GenerateIntegrationTests(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, IntegrationTestsTemplate, def)
}

func generatePackageDirective(writer io.Writer) {
	printTo(writer, PackageTemplate, os.Getenv("GOPACKAGE"))
}

func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
