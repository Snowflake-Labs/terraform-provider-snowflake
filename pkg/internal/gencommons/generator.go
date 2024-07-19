package gencommons

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

// TODO: describe
type ObjectNameProvider interface {
	ObjectName() string
}

// TODO: describe
// TODO: use
// TODO: better func
type GenerationModel interface {
	SomeFunc()
}

// TODO: add type for objects provider?
// TODO: add erorrs to any of these functions?
type Generator[T ObjectNameProvider, M GenerationModel] struct {
	objectsProvider func() []T
	modelProvider   func(T) M
	// TODO: add filename to model?
	filenameProvider func(T, M) string
	templates        []*template.Template

	additionalObjectDebugLogProviders []func([]T)
}

func NewGenerator[T ObjectNameProvider, M GenerationModel](objectsProvider func() []T, modelProvider func(T) M, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	return &Generator[T, M]{
		objectsProvider:  objectsProvider,
		modelProvider:    modelProvider,
		filenameProvider: filenameProvider,
		templates:        templates,

		additionalObjectDebugLogProviders: make([]func([]T), 0),
	}
}

func (g *Generator[T, M]) WithAdditionalObjectsDebugLogs(objectLogsProvider func([]T)) *Generator[T, M] {
	g.additionalObjectDebugLogProviders = append(g.additionalObjectDebugLogProviders, objectLogsProvider)
	return g
}

func (g *Generator[_, _]) Run() error {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	objects := g.objectsProvider()

	// TODO: print conditionally from invocation flag
	for _, p := range g.additionalObjectDebugLogProviders {
		p(objects)
	}

	// TODO: do not generate twice?
	// TODO: print conditionally from invocation flag
	if err := GenerateAndPrintForAllObjects(objects, g.modelProvider, g.templates...); err != nil {
		return err
	}
	if err := GenerateAndSaveForAllObjects(
		objects,
		g.modelProvider,
		g.filenameProvider,
		g.templates...,
	); err != nil {
		return err
	}

	return nil
}

func (g *Generator[_, _]) RunAndHandleOsReturn() {
	err := g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
