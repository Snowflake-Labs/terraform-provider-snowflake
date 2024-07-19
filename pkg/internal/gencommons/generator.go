package gencommons

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	preprocessArgs()

	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	// TODO: describe running with build flags: make generate-show-output-schemas SF_TF_GENERATOR_ARGS='--dry-run additional-logs'
	var additionalLogs = flag.Bool("additional-logs", false, "print additional object debug logs")
	var dryRun = flag.Bool("dry-run", false, "generate to std out instead of saving")
	flag.Parse()

	objects := g.objectsProvider()

	if *additionalLogs {
		for _, p := range g.additionalObjectDebugLogProviders {
			p(objects)
		}
	}

	if *dryRun {
		if err := generateAndPrintForAllObjects(objects, g.modelProvider, g.templates...); err != nil {
			return err
		}
	} else {
		if err := generateAndSaveForAllObjects(
			objects,
			g.modelProvider,
			g.filenameProvider,
			g.templates...,
		); err != nil {
			return err
		}
	}

	return nil
}

// TODO: temporary hacky solution to allow easy passing multiple args from the make command
func preprocessArgs() {
	rest := os.Args[1:]
	newArgs := make([]string, 0)
	newArgs = append(newArgs, os.Args[0])
	for _, a := range rest {
		newArgs = append(newArgs, strings.Split(a, " ")...)
	}
	os.Args = newArgs
}

func (g *Generator[_, _]) RunAndHandleOsReturn() {
	err := g.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func generateAndSaveForAllObjects[T ObjectNameProvider, M GenerationModel](objects []T, modelProvider func(T) M, filenameProvider func(T, M) string, templates ...*template.Template) error {
	var errs []error
	for _, s := range objects {
		buffer := bytes.Buffer{}
		model := modelProvider(s)
		if err := executeAllTemplates(model, &buffer, templates...); err != nil {
			errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
			continue
		}
		filename := filenameProvider(s, model)
		if err := WriteCodeToFile(&buffer, filename); err != nil {
			errs = append(errs, fmt.Errorf("saving output for object %s to file %s failed with err: %w", s.ObjectName(), filename, err))
			continue
		}
	}
	return errors.Join(errs...)
}

func generateAndPrintForAllObjects[T ObjectNameProvider, M GenerationModel](objects []T, modelProvider func(T) M, templates ...*template.Template) error {
	var errs []error
	for _, s := range objects {
		fmt.Println("===========================")
		fmt.Printf("Generating for object %s\n", s.ObjectName())
		fmt.Println("===========================")
		if err := executeAllTemplates(modelProvider(s), os.Stdout, templates...); err != nil {
			errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
			continue
		}
	}
	return errors.Join(errs...)
}

func executeAllTemplates[M GenerationModel](model M, writer io.Writer, templates ...*template.Template) error {
	var errs []error
	for _, t := range templates {
		if err := t.Execute(writer, model); err != nil {
			errs = append(errs, fmt.Errorf("template execution for template %s failed with err: %w", t.Name(), err))
		}
	}
	return errors.Join(errs...)
}
