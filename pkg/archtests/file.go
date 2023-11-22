package archtests

import (
	"go/ast"
	"regexp"
)

var (
	acceptanceTestFileRegex *regexp.Regexp
	acceptanceTestNameRegex *regexp.Regexp
	testFileRegex           *regexp.Regexp
	testNameRegex           *regexp.Regexp
)

func init() {
	var err error
	acceptanceTestFileRegex, err = regexp.Compile("^.*_acceptance_test.go$")
	if err != nil {
		panic(err)
	}
	acceptanceTestNameRegex, err = regexp.Compile("^TestAcc_.*$")
	if err != nil {
		panic(err)
	}
	testFileRegex, err = regexp.Compile("^.*_test.go$")
	if err != nil {
		panic(err)
	}
	testNameRegex, err = regexp.Compile("^Test.*$")
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

func (f *File) allExportedMethods() []Method {
	allExportedMethods := make([]Method, 0)
	for _, d := range f.fileSrc.Decls {
		switch d.(type) {
		case *ast.FuncDecl:
			name := d.(*ast.FuncDecl).Name.Name
			if ast.IsExported(name) {
				allExportedMethods = append(allExportedMethods, *NewMethod(name))
			}
		}
	}
	return allExportedMethods
}
