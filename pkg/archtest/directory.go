package archtest

import (
	"go/parser"
	"go/token"
	"io/fs"
)

type Directory struct {
	path string
}

func NewDirectory(path string) *Directory {
	return &Directory{path: path}
}

func (d *Directory) AllFiles() ([]File, error) {
	return d.files(nil)
}

func (d *Directory) files(filter func(fi fs.FileInfo) bool) ([]File, error) {
	packagesDict, err := parser.ParseDir(token.NewFileSet(), d.path, filter, 0)
	if err != nil {
		return nil, err
	}
	files := make([]File, 0)
	for packageName, astPackage := range packagesDict {
		for fileName, fileSrc := range astPackage.Files {
			files = append(files, *NewFile(packageName, fileName, fileSrc))
		}
	}
	return files, nil
}
