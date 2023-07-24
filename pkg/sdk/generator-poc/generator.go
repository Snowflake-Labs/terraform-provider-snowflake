package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	blueprintPath = flag.String("blueprint", "", ".json file containing blueprint for generation; must be set")
)

func main() {
	fmt.Printf("Running generator on %s with args %#v\n", os.Getenv("GOFILE"), os.Args[1:])
	blueprint := loadBlueprint()
	fmt.Printf("Loaded blueprint %#v.\n", blueprint)
}

func loadBlueprint() Blueprint {
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
	if errOpenFile != nil {
		log.Panicln(errReadBytes)
	}

	var blueprint Blueprint
	if err := json.Unmarshal(byteValue, &blueprint); err != nil {
		log.Panicln(err)
	}
	return blueprint
}

type Blueprint struct {
	Name string
}
