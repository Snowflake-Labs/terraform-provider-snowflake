package archtests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/archtest"
	"github.com/stretchr/testify/require"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesDirectory := archtest.NewDirectory("../resources/")
	resourcesFiles, err := resourcesDirectory.AllFiles()
	require.NoError(t, err)

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := archtest.FilterFiles(resourcesFiles, archtest.FileNameFilterProvider(archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(acceptanceTestFiles, func(file *archtest.File) {
			archtest.AssertPackage(t, file, "resources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := archtest.FilterFiles(resourcesFiles, archtest.FileNameFilterProvider(archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(acceptanceTestFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := archtest.FilterFiles(resourcesFiles, archtest.FileNameFilterWithExclusionsProvider(archtest.TestFileRegex, archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(otherTestFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertMethodNameDoesNotMatch(t, file, method, archtest.AcceptanceTestNameRegex)
				archtest.AssertMethodNameMatches(t, file, method, archtest.TestNameRegex)
			})
		})
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := archtest.FilterFiles(resourcesFiles, archtest.PackageFilterProvider("resources_test"))

		archtest.IterateFiles(packageFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})
}
