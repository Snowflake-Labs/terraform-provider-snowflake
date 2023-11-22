package archtests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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

func filterFiles(files []File, filter func(*File) bool) []File {
	filteredFiles := make([]File, 0)
	for _, f := range files {
		if filter(&f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	return filteredFiles
}

func fileNameFilterProvider(regex *regexp.Regexp) func(f *File) bool {
	return func(f *File) bool {
		return regex.Match([]byte(f.fileName))
	}
}

func fileNameFilterWithExclusionsProvider(regex *regexp.Regexp, exclusionRegex ...*regexp.Regexp) func(f *File) bool {
	return func(f *File) bool {
		matches := regex.MatchString(f.fileName)
		for _, e := range exclusionRegex {
			matches = matches && !e.MatchString(f.fileName)
		}
		return matches
	}
}

func packageFilterProvider(packageName string) func(f *File) bool {
	return func(f *File) bool {
		return f.packageName == packageName
	}
}

func allExportedMethodsInFile(src *ast.File) []string {
	allExportedMethods := make([]string, 0)
	for _, d := range src.Decls {
		switch d.(type) {
		case *ast.FuncDecl:
			name := d.(*ast.FuncDecl).Name.Name
			if ast.IsExported(name) {
				allExportedMethods = append(allExportedMethods, name)
			}
		}
	}
	return allExportedMethods
}
