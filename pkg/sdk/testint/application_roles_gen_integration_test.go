package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInt_ApplicationRoles setup is a little bit different from usual integration test, because of how native apps work.
// I will try to explain it in a short form, but check out this article for more detailed description (https://docs.snowflake.com/en/developer-guide/native-apps/tutorials/getting-started-tutorial#introduction)
//   - create stage - it is where we will be keeping our application files
//   - put native app specific stuff onto our stage (manifest.yml and setup.sql)
//   - create an application package and a new version of our application
//   - create an application with the application package and the version we just created
//   - while creating the application, the setup.sql script will be run in our application context (and that is where application roles for our tests are created)
//   - we're ready to query application roles we have just created
func TestInt_ApplicationRoles(t *testing.T) {
	client := testClient(t)

	stageName := random.AlphaN(8)
	stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, stageName))
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

	assertApplicationRole := func(t *testing.T, appRole *sdk.ApplicationRole, name string, comment string) {
		t.Helper()
		assert.Equal(t, name, appRole.Name)
		assert.Equal(t, appName, appRole.Owner)
		assert.Equal(t, comment, appRole.Comment)
		assert.Equal(t, "APPLICATION", appRole.OwnerRoleType)
	}

	assertApplicationRoles := func(t *testing.T, appRoles []sdk.ApplicationRole, name string, comment string) {
		t.Helper()
		appRole, err := collections.FindOne(appRoles, func(role sdk.ApplicationRole) bool {
			return role.Name == name
		})
		require.NoError(t, err)
		assertApplicationRole(t, appRole, name, comment)
	}

	t.Run("Show by id", func(t *testing.T) {
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		appRole, err := client.ApplicationRoles.ShowByID(ctx, sdk.NewShowByIDApplicationRoleRequest(id, sdk.NewAccountObjectIdentifier(appName)))
		require.NoError(t, err)

		assertApplicationRole(t, appRole, name, "some comment")
	})

	t.Run("Show", func(t *testing.T) {
		ctx := context.Background()
		req := sdk.NewShowApplicationRoleRequest().
			WithApplicationName(sdk.NewAccountObjectIdentifier(appName)).
			WithLimit(&sdk.LimitFrom{
				Rows: sdk.Int(2),
			})
		appRoles, err := client.ApplicationRoles.Show(ctx, req)
		require.NoError(t, err)

		assertApplicationRoles(t, appRoles, "app_role_1", "some comment")
		assertApplicationRoles(t, appRoles, "app_role_2", "some comment2")
	})

	t.Run("show grants to application role", func(t *testing.T) {
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		opts := new(sdk.ShowGrantOptions)
		opts.To = &sdk.ShowGrantsTo{
			ApplicationRole: id,
		}
		grants, err := client.Grants.Show(ctx, opts)
		require.NoError(t, err)

		require.NotEmpty(t, grants)
		require.NotEmpty(t, grants[0].CreatedOn)
		require.Equal(t, sdk.ObjectPrivilegeUsage.String(), grants[0].Privilege)
		require.Equal(t, sdk.ObjectTypeDatabase, grants[0].GrantedOn)
		require.Equal(t, sdk.ObjectTypeApplicationRole, grants[0].GrantedTo)
	})

	t.Run("show grants of application role", func(t *testing.T) {
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		opts := new(sdk.ShowGrantOptions)
		opts.Of = &sdk.ShowGrantsOf{
			ApplicationRole: id,
		}
		grants, err := client.Grants.Show(ctx, opts)
		require.NoError(t, err)

		require.NotEmpty(t, grants)
		require.NotEmpty(t, grants[0].CreatedOn)
		require.Equal(t, sdk.ObjectTypeRole, grants[0].GrantedTo)
		require.Equal(t, sdk.NewAccountObjectIdentifier(appName), grants[0].GrantedBy)
	})

	t.Run("show grants to application", func(t *testing.T) {
		// Need second app to be able to grant application role to it. Cannot grant to parent application (098806 (0A000): Cannot grant an APPLICATION ROLE to the parent APPLICATION).
		stageName2 := random.AlphaN(8)
		stage2, cleanupStage2 := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, stageName2))
		t.Cleanup(cleanupStage2)

		putOnStage(t, client, stage2, "manifest2.yml")
		putOnStage(t, client, stage2, "setup.sql")

		appName2 := "snowflake_app_2"
		cleanupApp2 := createApplication(t, client, appName2, appPackageName, versionName)
		t.Cleanup(cleanupApp2)

		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(appName, name)
		ctx := context.Background()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`GRANT APPLICATION ROLE %s TO APPLICATION %s`, id.FullyQualifiedName(), sdk.NewAccountObjectIdentifier(appName2).FullyQualifiedName()))
		require.NoError(t, err)
		defer func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`REVOKE APPLICATION ROLE %s FROM APPLICATION %s`, id.FullyQualifiedName(), sdk.NewAccountObjectIdentifier(appName2).FullyQualifiedName()))
			require.NoError(t, err)
		}()

		opts := new(sdk.ShowGrantOptions)
		opts.To = &sdk.ShowGrantsTo{
			Application: sdk.NewAccountObjectIdentifier(appName2),
		}
		grants, err := client.Grants.Show(ctx, opts)
		require.NoError(t, err)

		require.NotEmpty(t, grants)
		require.NotEmpty(t, grants[0].CreatedOn)
		require.Equal(t, sdk.ObjectPrivilegeUsage.String(), grants[0].Privilege)
		require.Equal(t, sdk.ObjectTypeApplicationRole, grants[0].GrantedOn)
		require.Equal(t, sdk.ObjectTypeApplication, grants[0].GrantedTo)
	})
}
