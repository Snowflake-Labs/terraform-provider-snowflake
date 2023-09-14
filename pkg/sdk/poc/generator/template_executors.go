package generator

import (
	"io"
	"log"
	"os"
	"text/template"
)

func GenerateInterface(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, InterfaceTemplate, def)
	for _, o := range def.Operations {
		generateOptionsStruct(writer, o)
	}
}

func generateOptionsStruct(writer io.Writer, operation *Operation) {
	printTo(writer, OptionsTemplate, operation)

	for _, f := range operation.HelperStructs {
		// TODO Merge with OptionsTemplate, but abstract documentation (move doc to field, every field (struct) will have it's own doc)
		// _, _ = writer.Write([]byte(fmt.Sprintf(
		//	"// %s is used to decode the result of a %s %s query.",
		//	f.Name,
		//	operation.Name,
		//	operation.ObjectInterface.Name,
		//)))
		generateStruct(writer, f)
	}

	for _, f := range operation.OptsField.Fields {
		if len(f.Fields) > 0 {
			generateStruct(writer, f)
		}
	}
}

func generateStruct(writer io.Writer, field *Field) {
	printTo(writer, StructTemplate, field)

	for _, f := range field.Fields {
		if len(f.Fields) > 0 {
			generateStruct(writer, f)
		}
	}
}

func GenerateDtos(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, DtoTemplate, def)
}

func GenerateImplementation(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ImplementationTemplate, def)
}

func GenerateUnitTests(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, TestFuncTemplate, def)
}

func GenerateValidations(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ValidationsImplTemplate, def)
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
