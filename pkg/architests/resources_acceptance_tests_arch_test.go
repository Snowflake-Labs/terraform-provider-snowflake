package architests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/architest"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesDirectory := architest.NewDirectory("../resources/")
	resourcesFiles := resourcesDirectory.AllFiles()

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles := resourcesFiles.Filter(architest.FileNameFilterProvider(architest.AcceptanceTestFileRegex))

		acceptanceTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "resources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles := resourcesFiles.Filter(architest.FileNameFilterProvider(architest.AcceptanceTestFileRegex))

		acceptanceTestFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := resourcesFiles.Filter(architest.FileNameFilterWithExclusionsProvider(architest.TestFileRegex, architest.AcceptanceTestFileRegex))

		otherTestFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertNameDoesNotMatch(t, architest.AcceptanceTestNameRegex)
				method.AssertNameMatches(t, architest.TestNameRegex)
			})
		})
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := resourcesFiles.Filter(architest.PackageFilterProvider("resources_test"))

		packageFiles.All(func(file *architest.File) {
			file.AllExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t)
			})
		})
	})
}
