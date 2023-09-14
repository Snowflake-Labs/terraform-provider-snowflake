package generator

import (
	"golang.org/x/exp/slices"
	"io"
	"log"
	"os"
	"text/template"
)

// TODO Move all the logic from templates here, so we could have more control and e.g.
//	 do not generate structs that already have been generated (reference to the same object twice)

// TODO We could hold references to struct - generatedStructs []*Field which is better,
//	 but we're passing dbStruct or plainStruct
//	 we're then creating new *Field internally (because it's better from DSL point of view) and it's problematic
//	 but we can either:
//		1. build them in DSL file and pass *Field's
//		2. we can keep internal cache, so when someone passes dbStruct or plainStruct or any other *Field abstraction
//			it will create *Field if missing or return already existing instance pointing to the same memory address
// 			which we could filter here with generatedStructs []*Field

var generatedStructs []string

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
}

func GenerateImplementation(writer io.Writer, def *Interface) {
	generatePackageDirective(writer)
	printTo(writer, ImplementationTemplate, def)
	// TODO ToOpts template
	// TODO Convert
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
