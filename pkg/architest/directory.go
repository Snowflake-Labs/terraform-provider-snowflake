package architest

import (
	"go/parser"
	"go/token"
)

type Directory struct {
	path  string
	files Files
}

func NewDirectory(path string) *Directory {
	packagesDict, err := parser.ParseDir(token.NewFileSet(), path, nil, 0)
	if err != nil {
		panic(err)
	}
	files := make(Files, 0)
	for packageName, astPackage := range packagesDict {
		for fileName, fileSrc := range astPackage.Files {
			files = append(files, *NewFile(packageName, fileName, fileSrc))
		}
	}
	return &Directory{
		path:  path,
		files: files,
	}
}

func (d *Directory) AllFiles() Files {
	return d.files
}

func (d *Directory) Files(filter FileFilter) Files {
	return d.AllFiles().Filter(filter)
}
