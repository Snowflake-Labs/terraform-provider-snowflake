package architest

import (
	"go/parser"
	"go/token"
)

type directory struct {
	path  string
	files Files
}

type FilesProvider interface {
	AllFiles() Files
	Files(filter FileFilter) Files
}

func Directory(path string) FilesProvider {
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
	return &directory{
		path:  path,
		files: files,
	}
}

func (d *directory) AllFiles() Files {
	return d.files
}

func (d *directory) Files(filter FileFilter) Files {
	return d.AllFiles().Filter(filter)
}
