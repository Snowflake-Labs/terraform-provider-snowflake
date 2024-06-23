package gen

import (
	"io"
	"log"
)

// TODO: handle panics better
// TODO: test and describe
func Generate(structDetails Struct, writer io.Writer) {
	model := ModelFromStructDetails(structDetails)
	err := SchemaTemplate.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
	err = ToSchemaMapperTemplate.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
