package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Running generator on %s with args %#v\n", os.Getenv("GOFILE"), os.Args[1:])
}
