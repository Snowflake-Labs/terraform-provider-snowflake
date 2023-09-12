package generator2

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

	for _, f := range operation.Options.Fields {
		s := f.IntoStruct()
		if s != nil {
			printTo(writer, StructTemplate, s)
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

//	func GenerateUnitTests(writer io.Writer, def *Interface) {
//		generatePackageDirective(writer)
//		printTo(writer, TestFuncTemplate, def)
//	}

func GenerateValidations(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ValidationsImplTemplate, def)
}

//	func GenerateIntegrationTests(writer io.Writer, def *Interface) {
//		generatePackageDirective(writer)
//		printTo(writer, IntegrationTestsTemplate, def)
//	}

func generatePackageDirective(writer io.Writer) {
	printTo(writer, PackageTemplate, os.Getenv("GOPACKAGE"))
}

func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
