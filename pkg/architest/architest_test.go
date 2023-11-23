package architest_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/architest"
	"github.com/stretchr/testify/assert"
)

func Test_Directory(t *testing.T) {
	tests1 := []struct {
		directory            string
		expectedFileNames    []string
		expectedPackageNames []string
	}{
		{directory: "testdata/dir1", expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go", "testdata/dir1/different1.go"}, expectedPackageNames: []string{"dir1"}},
		{directory: "testdata/dir2", expectedFileNames: []string{"testdata/dir2/sample1.go", "testdata/dir2/sample1_test.go"}, expectedPackageNames: []string{"dir2", "dir2_test"}},
		{directory: "testdata/dir3", expectedFileNames: []string{"testdata/dir3/sample1.go", "testdata/dir3/sample1_acceptance_test.go"}, expectedPackageNames: []string{"dir3", "dir3_test"}},
	}
	for _, tt := range tests1 {
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

	tests2 := []struct {
		directory         string
		filter            architest.FileFilter
		expectedFileNames []string
	}{
		{directory: "testdata/dir1", filter: architest.FileNameFilterProvider("sample"), expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}},
		{directory: "testdata/dir1", filter: architest.FileNameRegexFilterProvider(regexp.MustCompile("sample")), expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}},
		{directory: "testdata/dir1", filter: architest.FileNameFilterWithExclusionsProvider(regexp.MustCompile("sample"), regexp.MustCompile("sample1")), expectedFileNames: []string{"testdata/dir1/sample2.go"}},
		{directory: "testdata/dir2", filter: architest.PackageFilterProvider("dir2"), expectedFileNames: []string{"testdata/dir2/sample1.go"}},
		{directory: "testdata/dir2", filter: architest.PackageFilterProvider("dir2_test"), expectedFileNames: []string{"testdata/dir2/sample1_test.go"}},
		{directory: "testdata/dir2", filter: architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex), expectedFileNames: []string{}},
		{directory: "testdata/dir3", filter: architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex), expectedFileNames: []string{"testdata/dir3/sample1_acceptance_test.go"}},
		{directory: "testdata/dir2", filter: architest.FileNameRegexFilterProvider(architest.TestFileRegex), expectedFileNames: []string{"testdata/dir2/sample1_test.go"}},
		{directory: "testdata/dir3", filter: architest.FileNameRegexFilterProvider(architest.TestFileRegex), expectedFileNames: []string{"testdata/dir3/sample1_acceptance_test.go"}},
	}
	for _, tt := range tests2 {
		t.Run("list only files matching filter in the given directory", func(t *testing.T) {
			dir := architest.NewDirectory(tt.directory)

			filteredFiles := dir.Files(tt.filter)
			assert.Len(t, filteredFiles, len(tt.expectedFileNames))

			fileNames := make([]string, 0, len(filteredFiles))
			for _, f := range filteredFiles {
				fileNames = append(fileNames, f.FileName())
			}
			assert.ElementsMatch(t, fileNames, tt.expectedFileNames)
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

func Test_Assertions(t *testing.T) {
	tests1 := []struct {
		filePath        string
		expectedPackage string
	}{
		{filePath: "testdata/dir1/sample1.go", expectedPackage: "dir1"},
		{filePath: "testdata/dir2/sample1.go", expectedPackage: "dir2"},
		{filePath: "testdata/dir2/sample1_test.go", expectedPackage: "dir2_test"},
	}
	for _, tt := range tests1 {
		t.Run("file package assertions", func(t *testing.T) {
			file := architest.NewFileFromPath(tt.filePath)
			tut1 := &testing.T{}
			tut2 := &testing.T{}

			file.AssertHasPackage(tut1, tt.expectedPackage)
			file.AssertHasPackage(tut2, "some_other_package")

			assert.Equal(t, false, tut1.Failed())
			assert.Equal(t, true, tut2.Failed())
		})
	}

	tests2 := []struct {
		methodName string
		correct    bool
	}{
		{methodName: "TestAcc_abc", correct: true},
		{methodName: "TestAcc_TestAcc_Test", correct: true},
		{methodName: "TestAcc_", correct: false},
		{methodName: "ATestAcc_", correct: false},
		{methodName: "TestAcc", correct: false},
		{methodName: "Test_", correct: false},
		{methodName: "TestAccc_", correct: false},
	}
	for _, tt := range tests2 {
		t.Run(fmt.Sprintf("acceptance test name assertions for method %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut := &testing.T{}

			method.AssertAcceptanceTestNamedCorrectly(tut)

			assert.Equal(t, !tt.correct, tut.Failed())
		})
	}

	tests3 := []struct {
		methodName string
		regexRaw   string
		correct    bool
	}{
		{methodName: "sample1", regexRaw: "sample", correct: true},
		{methodName: "Sample1", regexRaw: "sample", correct: false},
		{methodName: "Sample1", regexRaw: "Sample", correct: true},
	}
	for _, tt := range tests3 {
		t.Run(fmt.Sprintf("matching and not matching method name assertions for %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut1 := &testing.T{}
			tut2 := &testing.T{}

			method.AssertNameMatches(tut1, regexp.MustCompile(tt.regexRaw))
			method.AssertNameDoesNotMatch(tut2, regexp.MustCompile(tt.regexRaw))

			assert.Equal(t, !tt.correct, tut1.Failed())
			assert.Equal(t, tt.correct, tut2.Failed())
		})
	}

	tests4 := []struct {
		methodName string
		correct    bool
	}{
		{methodName: "Test", correct: true},
		{methodName: "aTest", correct: false},
		{methodName: "Test_", correct: true},
		{methodName: "Test_adsfadf", correct: true},
	}
	for _, tt := range tests4 {
		t.Run(fmt.Sprintf("test name assertions for method %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut := &testing.T{}

			method.AssertTestNamedCorrectly(tut)

			assert.Equal(t, !tt.correct, tut.Failed())
		})
	}
}
