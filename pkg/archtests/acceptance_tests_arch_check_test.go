package archtests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesPath := "../resources/"
	resourcesFiles, err := filesInDirectory(resourcesPath, nil)
	require.NoError(t, err)

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := filterFiles(resourcesFiles, fileNameFilterProvider(acceptanceTestFileRegex))

		for _, file := range acceptanceTestFiles {
			assertPackage(t, &file, "resources_test")
		}
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := filterFiles(resourcesFiles, fileNameFilterProvider(acceptanceTestFileRegex))

		for _, file := range acceptanceTestFiles {
			for _, method := range allExportedMethodsInFile(file.fileSrc) {
				assertAcceptanceTestNamedCorrectly(t, &file, method)
			}
		}
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := filterFiles(resourcesFiles, fileNameFilterWithExclusionsProvider(testFileRegex, acceptanceTestFileRegex))

		for _, file := range otherTestFiles {
			for _, method := range allExportedMethodsInFile(file.fileSrc) {
				assertMethodNameDoesNotMatch(t, &file, method, acceptanceTestNameRegex)
				assertMethodNameMatches(t, &file, method, testNameRegex)
			}
		}
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := filterFiles(resourcesFiles, packageFilterProvider("resources_test"))

		for _, file := range packageFiles {
			for _, method := range allExportedMethodsInFile(file.fileSrc) {
				assertAcceptanceTestNamedCorrectly(t, &file, method)
			}
		}
	})
}
