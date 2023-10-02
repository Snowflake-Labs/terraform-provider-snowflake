package generator

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
)

func WriteCodeToFile(buffer *bytes.Buffer, fileName string) {
	wd, errWd := os.Getwd()
	if errWd != nil {
		log.Panicln(errWd)
	}
	outputPath := filepath.Join(wd, fileName)
	src, errSrcFormat := format.Source(buffer.Bytes())
	if errSrcFormat != nil {
		log.Panicln(errSrcFormat)
	}
	if err := os.WriteFile(outputPath, src, 0o600); err != nil {
		log.Panicln(err)
	}
}
