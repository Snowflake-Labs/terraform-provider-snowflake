package archtests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesDirectory := NewDirectory("../resources/")
	resourcesFiles, err := resourcesDirectory.allFiles()
	require.NoError(t, err)

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := filterFiles(resourcesFiles, fileNameFilterProvider(acceptanceTestFileRegex))

		iterateFiles(acceptanceTestFiles, func(f *File) {
			assertPackage(t, f, "resources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := filterFiles(resourcesFiles, fileNameFilterProvider(acceptanceTestFileRegex))

		iterateFiles(acceptanceTestFiles, func(f *File) {
			for _, method := range f.allExportedMethods() {
				assertAcceptanceTestNamedCorrectly(t, f, method)
			}
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := filterFiles(resourcesFiles, fileNameFilterWithExclusionsProvider(testFileRegex, acceptanceTestFileRegex))

		iterateFiles(otherTestFiles, func(f *File) {
			for _, method := range f.allExportedMethods() {
				assertMethodNameDoesNotMatch(t, f, method, acceptanceTestNameRegex)
				assertMethodNameMatches(t, f, method, testNameRegex)
			}
		})
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := filterFiles(resourcesFiles, packageFilterProvider("resources_test"))

		iterateFiles(packageFiles, func(f *File) {
			for _, method := range f.allExportedMethods() {
				assertAcceptanceTestNamedCorrectly(t, f, method)
			}
		})
	})
}
