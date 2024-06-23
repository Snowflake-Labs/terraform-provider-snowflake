package gen

import (
	"io"
	"log"
)

// TODO: handle panics better
// TODO: test and describe
func Generate(model ShowResultSchemaModel, writer io.Writer) {
	err := SchemaTemplate.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
	err = ToSchemaMapperTemplate.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}
