package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	blueprintPath = flag.String("blueprint", "", ".json file containing blueprint for generation; must be set")
)

func main() {
	fmt.Printf("Running generator on %s with args %#v\n", os.Getenv("GOFILE"), os.Args[1:])

	blueprint := loadBlueprint()
	fmt.Printf("Loaded blueprint %#v.\n", blueprint)

	gen := setUpGenerator(blueprint)
	fmt.Printf("Generator set up for package \"%s\" with output name \"%s\".\n", gen.outputPackage, gen.outputName)
}

func loadBlueprint() *Blueprint {
	flag.Parse()

	if len(*blueprintPath) == 0 {
		flag.Usage()
		log.Panicln("Blueprint .json file was not specified.")
	}

	blueprintFile, errOpenFile := os.Open(*blueprintPath)
	if errOpenFile != nil {
		log.Panicln(errOpenFile)
	}
	defer blueprintFile.Close()
	fmt.Printf("Opened blueprint file %s\n", blueprintFile.Name())

	byteValue, errReadBytes := io.ReadAll(blueprintFile)
	if errReadBytes != nil {
		log.Panicln(errReadBytes)
	}

	var blueprint Blueprint
	if err := json.Unmarshal(byteValue, &blueprint); err != nil {
		log.Panicln(err)
	}
	return &blueprint
}

type Blueprint struct {
	Name string
}

func setUpGenerator(blueprint *Blueprint) *Generator {
	wd, errWd := os.Getwd()
	if errWd != nil {
		log.Panicln(errWd)
	}

	file := os.Getenv("GOFILE")
	fileWithoutSuffix, _ := strings.CutSuffix(file, ".go")
	baseName := fmt.Sprintf("%s_generated.go", fileWithoutSuffix)
	outputName := filepath.Join(wd, baseName)

	return &Generator{
		buffer:        bytes.Buffer{},
		blueprint:     blueprint,
		outputPackage: os.Getenv("GOPACKAGE"),
		outputName:    outputName,
	}
}

type Generator struct {
	buffer        bytes.Buffer
	blueprint     *Blueprint
	outputPackage string
	outputName    string
}
