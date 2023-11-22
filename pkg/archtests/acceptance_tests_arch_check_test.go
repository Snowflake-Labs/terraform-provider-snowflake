package archtests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesPath := "../resources/"
	acceptanceTestFileRegex, err := regexp.Compile("^.*_acceptance_test.go$")
	require.NoError(t, err)
	acceptanceTestFileNameRegex, err := regexp.Compile("^TestAcc_.*$")

	t.Run("acceptance tests for resources have the right package", func(t *testing.T) {
		packagesDict, err := parser.ParseDir(token.NewFileSet(), resourcesPath, fileNameFilterProvider(acceptanceTestFileRegex), 0)
		require.NoError(t, err)

		assert.Len(t, packagesDict, 1)
		assert.Contains(t, packagesDict, "resources_test")
		assert.NotContains(t, packagesDict, "resources")
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		packagesDict, err := parser.ParseDir(token.NewFileSet(), resourcesPath, fileNameFilterProvider(acceptanceTestFileRegex), 0)
		require.NoError(t, err)

		allMatchingFiles := packagesDict["resources_test"].Files
		for name, src := range allMatchingFiles {
			exportedMethods := allExportedMethodsInFile(src)
			for _, method := range exportedMethods {
				assert.Truef(t, acceptanceTestFileNameRegex.Match([]byte(method)), "filename %s contains exported method %s which does not match %s", name, method, acceptanceTestFileNameRegex.String())
			}
		}
	})
}

func fileNameFilterProvider(regex *regexp.Regexp) func(fi fs.FileInfo) bool {
	return func(fi fs.FileInfo) bool {
		return regex.Match([]byte(fi.Name()))
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
