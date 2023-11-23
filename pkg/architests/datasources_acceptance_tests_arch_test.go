package architests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/architest"
)

func TestArchCheck_AcceptanceTests_DataSources(t *testing.T) {
	datasourcesDirectory := architest.NewDirectory("../datasources/")
	datasourcesFiles := datasourcesDirectory.AllFiles()

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := datasourcesFiles.Filter(architest.FileNameFilterProvider(architest.AcceptanceTestFileRegex))

		acceptanceTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "datasources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := datasourcesFiles.Filter(architest.FileNameFilterProvider(architest.AcceptanceTestFileRegex))

		acceptanceTestFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t, file)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := datasourcesFiles.Filter(architest.FileNameFilterWithExclusionsProvider(architest.TestFileRegex, architest.AcceptanceTestFileRegex))

		otherTestFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertNameDoesNotMatch(t, file, architest.AcceptanceTestNameRegex)
				method.AssertNameMatches(t, file, architest.TestNameRegex)
			})
		})
	})

	t.Run("there are only acceptance tests in package datasources_test", func(t *testing.T) {
		packageFiles := datasourcesFiles.Filter(architest.PackageFilterProvider("datasources_test"))

		packageFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t, file)
			})
		})
	})
}
