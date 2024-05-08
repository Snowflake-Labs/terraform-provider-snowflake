package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
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
	ctx := testContext(t)

	createApp := func(t *testing.T) *sdk.Application {
		t.Helper()

		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		testClientHelper().Stage.PutOnStage(t, stage.ID(), "manifest.yml")
		testClientHelper().Stage.PutOnStage(t, stage.ID(), "setup.sql")

		applicationPackage, cleanupApplicationPackage := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
		t.Cleanup(cleanupApplicationPackage)

		testClientHelper().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), "v1")

		application, cleanupApplication := testClientHelper().Application.CreateApplication(t, applicationPackage.ID(), "v1")
		t.Cleanup(cleanupApplication)
		return application
	}

	application := createApp(t)

	assertApplicationRole := func(t *testing.T, appRole *sdk.ApplicationRole, name string, comment string) {
		t.Helper()
		assert.Equal(t, name, appRole.Name)
		assert.Equal(t, application.Name, appRole.Owner)
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

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(application.Name, name)
		ctx := context.Background()

		appRole, err := client.ApplicationRoles.ShowByID(ctx, sdk.NewAccountObjectIdentifier(application.Name), id)
		require.NoError(t, err)

		assertApplicationRole(t, appRole, name, "some comment")
	})

	t.Run("Show", func(t *testing.T) {
		req := sdk.NewShowApplicationRoleRequest().
			WithApplicationName(application.ID()).
			WithLimit(&sdk.LimitFrom{
				Rows: sdk.Int(2),
			})
		appRoles, err := client.ApplicationRoles.Show(ctx, req)
		require.NoError(t, err)

		assertApplicationRoles(t, appRoles, "app_role_1", "some comment")
		assertApplicationRoles(t, appRoles, "app_role_2", "some comment2")
	})

	t.Run("show grants to application - grant to role", func(t *testing.T) {
		role, cleanupRole := testClientHelper().Role.CreateRole(t)
		t.Cleanup(cleanupRole)

		id := sdk.NewDatabaseObjectIdentifier(application.Name, "app_role_1")
		// grant the application role to the role
		kindOfRole := sdk.NewKindOfRoleRequest().WithRoleName(sdk.Pointer(role.ID()))
		gr := sdk.NewGrantApplicationRoleRequest(id).WithTo(*kindOfRole)
		err := client.ApplicationRoles.Grant(ctx, gr)
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Role: role.ID()}})
		require.NoError(t, err)
		require.NotEmpty(t, grants)
		require.NotEmpty(t, grants[0].CreatedOn)
		require.Equal(t, sdk.ObjectPrivilegeUsage.String(), grants[0].Privilege)
		require.Equal(t, sdk.ObjectTypeApplicationRole, grants[0].GrantedOn)
		require.Equal(t, sdk.ObjectTypeRole, grants[0].GrantedTo)

		// revoke the application role from the role
		rr := sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*kindOfRole)
		err = client.ApplicationRoles.Revoke(ctx, rr)
		require.NoError(t, err)

		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Role: role.ID()}})
		require.NoError(t, err)
		require.Equal(t, 0, len(grants))
	})

	t.Run("show grants to application - grant to application", func(t *testing.T) {
		application2 := createApp(t)

		id := sdk.NewDatabaseObjectIdentifier(application.Name, "app_role_1")
		// grant the application role to the application
		kindOfRole := sdk.NewKindOfRoleRequest().WithApplicationName(sdk.Pointer(application2.ID()))
		gr := sdk.NewGrantApplicationRoleRequest(id).WithTo(*kindOfRole)
		err := client.ApplicationRoles.Grant(ctx, gr)
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Application: application2.ID()}})
		require.NoError(t, err)
		require.NotEmpty(t, grants)
		require.NotEmpty(t, grants[0].CreatedOn)
		require.Equal(t, sdk.ObjectPrivilegeUsage.String(), grants[0].Privilege)
		require.Equal(t, sdk.ObjectTypeApplicationRole, grants[0].GrantedOn)
		require.Equal(t, sdk.ObjectTypeApplication, grants[0].GrantedTo)

		// revoke the application role from the role
		rr := sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*kindOfRole)
		err = client.ApplicationRoles.Revoke(ctx, rr)
		require.NoError(t, err)

		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Application: application2.ID()}})
		require.NoError(t, err)
		require.Equal(t, 0, len(grants))
	})

	t.Run("show grants to application role", func(t *testing.T) {
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(application.Name, name)
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
		id := sdk.NewDatabaseObjectIdentifier(application.Name, name)
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
		require.Equal(t, application.ID(), grants[0].GrantedBy)
	})
}
