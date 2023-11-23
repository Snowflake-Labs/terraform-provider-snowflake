package architest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
)

var (
	AcceptanceTestFileRegex *regexp.Regexp
	AcceptanceTestNameRegex *regexp.Regexp
	TestFileRegex           *regexp.Regexp
	TestNameRegex           *regexp.Regexp
)

func init() {
	var err error
	AcceptanceTestFileRegex, err = regexp.Compile("^.*_acceptance_test.go$")
	if err != nil {
		panic(err)
	}
	AcceptanceTestNameRegex, err = regexp.Compile("^TestAcc_.*$")
	if err != nil {
		panic(err)
	}
	TestFileRegex, err = regexp.Compile("^.*_test.go$")
	if err != nil {
		panic(err)
	}
	TestNameRegex, err = regexp.Compile("^Test.*$")
	if err != nil {
		panic(err)
	}
}

type File struct {
	packageName string
	fileName    string
	fileSrc     *ast.File
}

func NewFile(packageName string, fileName string, fileSrc *ast.File) *File {
	return &File{
		packageName: packageName,
		fileName:    fileName,
		fileSrc:     fileSrc,
	}
}

func NewFileFromPath(path string) *File {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		panic(err)
	}
	return NewFile(file.Name.Name, filepath.Base(path), file)
}

func (f *File) PackageName() string {
	return f.packageName
}

func (f *File) FileName() string {
	return f.fileName
}

func (f *File) ExportedMethods() Methods {
	allExportedMethods := make(Methods, 0)
	for _, d := range f.fileSrc.Decls {
		switch d.(type) {
		case *ast.FuncDecl:
			name := d.(*ast.FuncDecl).Name.Name
			if ast.IsExported(name) {
				allExportedMethods = append(allExportedMethods, *NewMethod(name, f))
			}
		}
	}
	return allExportedMethods
}
