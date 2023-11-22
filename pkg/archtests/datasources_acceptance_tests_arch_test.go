package archtests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/archtest"
	"github.com/stretchr/testify/require"
)

func TestArchCheck_AcceptanceTests_DataSources(t *testing.T) {
	datasourcesDirectory := archtest.NewDirectory("../datasources/")
	datasourcesFiles, err := datasourcesDirectory.AllFiles()
	require.NoError(t, err)

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := archtest.FilterFiles(datasourcesFiles, archtest.FileNameFilterProvider(archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(acceptanceTestFiles, func(file *archtest.File) {
			archtest.AssertPackage(t, file, "datasources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := archtest.FilterFiles(datasourcesFiles, archtest.FileNameFilterProvider(archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(acceptanceTestFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := archtest.FilterFiles(datasourcesFiles, archtest.FileNameFilterWithExclusionsProvider(archtest.TestFileRegex, archtest.AcceptanceTestFileRegex))

		archtest.IterateFiles(otherTestFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertMethodNameDoesNotMatch(t, file, method, archtest.AcceptanceTestNameRegex)
				archtest.AssertMethodNameMatches(t, file, method, archtest.TestNameRegex)
			})
		})
	})

	t.Run("there are only acceptance tests in package datasources_test", func(t *testing.T) {
		packageFiles := archtest.FilterFiles(datasourcesFiles, archtest.PackageFilterProvider("datasources_test"))

		archtest.IterateFiles(packageFiles, func(file *archtest.File) {
			archtest.IterateMethods(file.AllExportedMethods(), func(method *archtest.Method) {
				archtest.AssertAcceptanceTestNamedCorrectly(t, file, method)
			})
		})
	})
}
