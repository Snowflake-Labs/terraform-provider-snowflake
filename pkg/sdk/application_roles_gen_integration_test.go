package sdk

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func putOnStage(t *testing.T, client *Client, stage *Stage, filename string) {
	t.Helper()
	ctx := context.Background()

	path, err := filepath.Abs("./test-resources/" + filename)
	require.NoError(t, err)
	absPath := "file:/" + path

	_, err = client.exec(ctx, fmt.Sprintf(`PUT '%s' @%s AUTO_COMPRESS = FALSE`, absPath, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func createApplicationPackage(t *testing.T, client *Client, name string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`CREATE APPLICATION PACKAGE %s`, name))
	require.NoError(t, err)
	return func() {
		_, err := client.exec(ctx, fmt.Sprintf(`DROP APPLICATION PACKAGE %s`, name))
		require.NoError(t, err)
	}
}

func addApplicationPackageVersion(t *testing.T, client *Client, stage *Stage, appPackageName string, versionName string) {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE %s ADD VERSION %v USING '@%s'`, appPackageName, versionName, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func createApplication(t *testing.T, client *Client, name string, packageName string, version string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`CREATE APPLICATION %s FROM APPLICATION PACKAGE %s USING VERSION %s`, name, packageName, version))
	require.NoError(t, err)
	return func() {
		_, err := client.exec(ctx, fmt.Sprintf(`DROP APPLICATION %s`, name))
		require.NoError(t, err)
	}
}

func TestInt_ApplicationRoles(t *testing.T) {
	client := testClient(t)

	db, cleanupDB := createDatabase(t, client)
	t.Cleanup(cleanupDB)

	schema, cleanupSchema := createSchema(t, client, db)
	t.Cleanup(cleanupSchema)

	stageName := "stage_name"
	stage, cleanupStage := createStage(t, client, db, schema, stageName)
	t.Cleanup(cleanupStage)

	putOnStage(t, client, stage, "manifest.yml")
	putOnStage(t, client, stage, "setup.sql")

	appPackageName := "snowflake_app_pkg"
	versionName := "v1"
	cleanupAppPackage := createApplicationPackage(t, client, appPackageName)
	t.Cleanup(cleanupAppPackage)
	addApplicationPackageVersion(t, client, stage, appPackageName, versionName)

	appName := "snowflake_app"
	cleanupApp := createApplication(t, client, appName, appPackageName, versionName)
	t.Cleanup(cleanupApp)

	t.Run("Create", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Grant", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Revoke", func(t *testing.T) {
		// TODO: fill me
	})
}
