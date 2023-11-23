package architest_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/architest"
	"github.com/stretchr/testify/assert"
)

// TODO: create jira issue in tests stabilization (+one-pager?)
func Test_Directory(t *testing.T) {
	tests := []struct {
		directory            string
		expectedFileNames    []string
		expectedPackageNames []string
	}{
		{directory: "testdata/dir1", expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}, expectedPackageNames: []string{"dir1"}},
		{directory: "testdata/dir2", expectedFileNames: []string{"testdata/dir2/sample1.go", "testdata/dir2/sample1_test.go"}, expectedPackageNames: []string{"dir2", "dir2_test"}},
		{directory: "testdata/dir3", expectedFileNames: []string{"testdata/dir3/sample1.go", "testdata/dir3/sample1_acceptance_test.go"}, expectedPackageNames: []string{"dir3", "dir3_test"}},
	}
	for _, tt := range tests {
		t.Run("list all files in the given directory", func(t *testing.T) {
			dir := architest.NewDirectory(tt.directory)

			allFiles := dir.AllFiles()
			assert.Len(t, allFiles, len(tt.expectedFileNames))

			fileNames := make([]string, 0, len(allFiles))
			for _, f := range allFiles {
				fileNames = append(fileNames, f.FileName())
			}
			assert.ElementsMatch(t, fileNames, tt.expectedFileNames)

			packageNames := make(map[string]bool)
			for _, f := range allFiles {
				packageNames[f.PackageName()] = true
			}
			assert.Len(t, packageNames, len(tt.expectedPackageNames))
			for _, name := range tt.expectedPackageNames {
				assert.Contains(t, packageNames, name)
			}
		})
	}
}

func Test_Files(t *testing.T) {
	tests := []struct {
		filePath            string
		expectedMethodNames []string
	}{
		{filePath: "testdata/dir1/sample1.go", expectedMethodNames: []string{}},
		{filePath: "testdata/dir1/sample2.go", expectedMethodNames: []string{"A"}},
	}
	for _, tt := range tests {
		t.Run("list all methods in file", func(t *testing.T) {
			file := architest.NewFileFromPath(tt.filePath)

			exportedMethods := file.ExportedMethods()
			assert.Len(t, exportedMethods, len(tt.expectedMethodNames))

			methodNames := make([]string, 0, len(exportedMethods))
			for _, m := range exportedMethods {
				methodNames = append(methodNames, m.Name())
			}
			assert.ElementsMatch(t, methodNames, tt.expectedMethodNames)
		})
	}
}
