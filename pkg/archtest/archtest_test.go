package archtest_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/archtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Directory(t *testing.T) {
	// TODO: parametrize test
	t.Run("list all files in the given directory", func(t *testing.T) {
		dir := archtest.NewDirectory("testdata/dir1")

		allFiles, err := dir.AllFiles()
		require.NoError(t, err)

		assert.Len(t, allFiles, 2)

		fileNames := make([]string, 0, len(allFiles))
		for _, f := range allFiles {
			fileNames = append(fileNames, f.FileName())
		}
		expectedFileNames := []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}
		assert.ElementsMatch(t, fileNames, expectedFileNames)

		packageNames := make(map[string]bool)
		for _, f := range allFiles {
			packageNames[f.PackageName()] = true
		}
		assert.Len(t, packageNames, 1)
		expectedPackageNames := []string{"dir1"}
		for _, name := range expectedPackageNames {
			assert.Contains(t, packageNames, name)
		}
	})
}
