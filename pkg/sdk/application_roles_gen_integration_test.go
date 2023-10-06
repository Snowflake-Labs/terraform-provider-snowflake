package sdk

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func putOnStage(t *testing.T, client *Client, stage *Stage, filename string) {
	t.Helper()
	ctx := context.Background()

	path, err := filepath.Abs("./test-resources/" + filename)
	require.NoError(t, err)
	absPath := "file://" + path

	_, err = client.exec(ctx, fmt.Sprintf(`PUT '%s' @%s AUTO_COMPRESS = FALSE`, absPath, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func createApplicationPackage(t *testing.T, client *Client, name string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`CREATE APPLICATION PACKAGE "%s"`, name))
	require.NoError(t, err)
	return func() {
		_, err := client.exec(ctx, fmt.Sprintf(`DROP APPLICATION PACKAGE "%s"`, name))
		require.NoError(t, err)
	}
}

func addApplicationPackageVersion(t *testing.T, client *Client, stage *Stage, appPackageName string, versionName string) {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE "%s" ADD VERSION %v USING '@%s'`, appPackageName, versionName, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func createApplication(t *testing.T, client *Client, name string, packageName string, version string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf(`CREATE APPLICATION "%s" FROM APPLICATION PACKAGE "%s" USING VERSION %s DEBUG_MODE = TRUE`, name, packageName, version))
	require.NoError(t, err)
	return func() {
		_, err := client.exec(ctx, fmt.Sprintf(`DROP APPLICATION "%s"`, name))
		require.NoError(t, err)
	}
}

func TestInt_ApplicationRoles(t *testing.T) {
	client := testClient(t)

	dbName := randomAlphanumericN(t, 32)
	db, cleanupDB := createDatabaseWithIdentifier(t, client, NewAccountObjectIdentifier(dbName))
	t.Cleanup(cleanupDB)

	schemaName := randomAlphanumericN(t, 32)
	schema, cleanupSchema := createSchemaWithIdentifier(t, client, db, schemaName)
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

	t.Run("Create Drop Show", func(t *testing.T) {
		name := randomAlphanumericN(t, 32)
		id := NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		createReq := NewCreateApplicationRoleRequest(id).
			WithIfNotExists(Bool(true)).
			WithComment(String("some comment"))
		err := client.ApplicationRoles.Create(ctx, createReq)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.ApplicationRoles.Drop(ctx, NewDropApplicationRoleRequest(id).WithIfExists(Bool(true)))
			require.NoError(t, err)
		})

		appRole, err := client.ApplicationRoles.ShowByID(ctx, NewShowByIDApplicationRoleRequest(id, NewAccountObjectIdentifier(appName)))
		require.NoError(t, err)

		assert.Equal(t, name, appRole.Name)
		assert.Equal(t, appName, appRole.Owner)
		assert.Equal(t, "some comment", appRole.Comment)
		assert.Equal(t, "APPLICATION", appRole.OwnerRoleType)
	})

	t.Run("Alter Rename to", func(t *testing.T) {
		name := randomAlphanumericN(t, 32)
		id := NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		err := client.ApplicationRoles.Create(ctx, NewCreateApplicationRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.ApplicationRoles.Drop(ctx, NewDropApplicationRoleRequest(id).WithIfExists(Bool(true)))
			require.NoError(t, err)
		})

		newName := randomAlphanumericN(t, 32)
		newId := NewDatabaseObjectIdentifier(appName, newName)
		err = client.ApplicationRoles.Alter(ctx, NewAlterApplicationRoleRequest(id).WithRenameTo(&newId))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.ApplicationRoles.Drop(ctx, NewDropApplicationRoleRequest(newId))
			require.NoError(t, err)
		})

		_, err = client.ApplicationRoles.ShowByID(ctx, NewShowByIDApplicationRoleRequest(newId, NewAccountObjectIdentifier(appName)))
		require.NoError(t, err)
	})

	t.Run("Alter Set Unset Comment", func(t *testing.T) {
		name := randomAlphanumericN(t, 32)
		id := NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		assertComment := func(comment string) {
			appRole, err := client.ApplicationRoles.ShowByID(ctx, NewShowByIDApplicationRoleRequest(id, NewAccountObjectIdentifier(appName)))
			require.NoError(t, err)
			assert.Equal(t, comment, appRole.Comment)
		}

		err := client.ApplicationRoles.Create(ctx, NewCreateApplicationRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.ApplicationRoles.Drop(ctx, NewDropApplicationRoleRequest(id).WithIfExists(Bool(true)))
			require.NoError(t, err)
		})
		assertComment("")

		err = client.ApplicationRoles.Alter(ctx, NewAlterApplicationRoleRequest(id).WithSetComment(String("some comment")))
		require.NoError(t, err)
		assertComment("some comment")

		err = client.ApplicationRoles.Alter(ctx, NewAlterApplicationRoleRequest(id).WithUnsetComment(Bool(true)))
		require.NoError(t, err)
		assertComment("")
	})

	t.Run("Grant", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Revoke", func(t *testing.T) {
		// TODO: fill me
	})
}
