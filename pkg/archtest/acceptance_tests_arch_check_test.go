package archtest

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

		iterateFiles(acceptanceTestFiles, func(file *File) {
			assertPackage(t, file, "resources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := filterFiles(resourcesFiles, fileNameFilterProvider(acceptanceTestFileRegex))

		iterateFiles(acceptanceTestFiles, func(file *File) {
			iterateMethods(file.allExportedMethods(), func(method *Method) {
				assertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := filterFiles(resourcesFiles, fileNameFilterWithExclusionsProvider(testFileRegex, acceptanceTestFileRegex))

		iterateFiles(otherTestFiles, func(file *File) {
			iterateMethods(file.allExportedMethods(), func(method *Method) {
				assertMethodNameDoesNotMatch(t, file, method, acceptanceTestNameRegex)
				assertMethodNameMatches(t, file, method, testNameRegex)
			})
		})
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := filterFiles(resourcesFiles, packageFilterProvider("resources_test"))

		iterateFiles(packageFiles, func(file *File) {
			iterateMethods(file.allExportedMethods(), func(method *Method) {
				assertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})
}
