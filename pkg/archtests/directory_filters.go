package archtests

import (
	"go/parser"
	"go/token"
	"io/fs"
)

func allFilesInDirectory(path string) ([]File, error) {
	return filesInDirectory(path, nil)
}

func filesInDirectory(path string, filter func(fi fs.FileInfo) bool) ([]File, error) {
	packagesDict, err := parser.ParseDir(token.NewFileSet(), path, filter, 0)
	if err != nil {
		return nil, err
	}
	files := make([]File, 0)
	for packageName, astPackage := range packagesDict {
		for fileName, fileSrc := range astPackage.Files {
			files = append(files, File{packageName, fileName, fileSrc})
		}
	}
	return files, nil
}
