package testint

import (
	"context"
	"fmt"
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

	createApp := func(t *testing.T) (*sdk.Application, *sdk.ApplicationPackage) {
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
		return application, applicationPackage
	}

	application, applicationPackage := createApp(t)

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

	assertGrantToRoles := func(t *testing.T, grants []sdk.Grant, id sdk.DatabaseObjectIdentifier, grantee sdk.AccountObjectIdentifier, ot sdk.ObjectType) {
		t.Helper()

		grant, err := collections.FindOne(grants, func(grant sdk.Grant) bool {
			return grant.Name.FullyQualifiedName() == id.FullyQualifiedName()
		})
		require.NoError(t, err)
		require.Equal(t, ot, grant.GrantedOn)
		require.Equal(t, grantee.FullyQualifiedName(), grant.GranteeName.FullyQualifiedName())
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

	t.Run("Grant and Revoke: Role", func(t *testing.T) {
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
		assertGrantToRoles(t, grants, id, role.ID(), sdk.ObjectTypeApplicationRole)

		// revoke the application role from the role
		rr := sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*kindOfRole)
		err = client.ApplicationRoles.Revoke(ctx, rr)
		require.NoError(t, err)

		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Role: role.ID()}})
		require.NoError(t, err)
		require.Equal(t, 0, len(grants))
	})

	t.Run("Grant and Revoke: Application", func(t *testing.T) {
		app, _ := createApp(t)

		id := sdk.NewDatabaseObjectIdentifier(application.Name, "app_role_1")
		// grant the application role to the application
		kindOfRole := sdk.NewKindOfRoleRequest().WithApplicationName(sdk.Pointer(app.ID()))
		gr := sdk.NewGrantApplicationRoleRequest(id).WithTo(*kindOfRole)
		err := client.ApplicationRoles.Grant(ctx, gr)
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Application: app.ID()}})
		require.NoError(t, err)
		assertGrantToRoles(t, grants, id, app.ID(), sdk.ObjectTypeApplicationRole)

		// revoke the application role from the role
		rr := sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*kindOfRole)
		err = client.ApplicationRoles.Revoke(ctx, rr)
		require.NoError(t, err)

		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{To: &sdk.ShowGrantsTo{Application: app.ID()}})
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
	t.Run("show grants to application", func(t *testing.T) {
		// Need second app to be able to grant application role to it. Cannot grant to parent application (098806 (0A000): Cannot grant an APPLICATION ROLE to the parent APPLICATION).
		stage2, cleanupStage2 := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage2)

		testClientHelper().Stage.PutOnStage(t, stage2.ID(), "manifest2.yml")
		testClientHelper().Stage.PutOnStage(t, stage2.ID(), "setup.sql")

		application2, cleanupApplication2 := testClientHelper().Application.CreateApplication(t, applicationPackage.ID(), "v1")
		t.Cleanup(cleanupApplication2)
		name := "app_role_1"
		id := sdk.NewDatabaseObjectIdentifier(application.Name, name)
		ctx := context.Background()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`GRANT APPLICATION ROLE %s TO APPLICATION %s`, id.FullyQualifiedName(), application2.ID().FullyQualifiedName()))
		require.NoError(t, err)
		defer func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`REVOKE APPLICATION ROLE %s FROM APPLICATION %s`, id.FullyQualifiedName(), application2.ID().FullyQualifiedName()))
			require.NoError(t, err)
		}()

		opts := new(sdk.ShowGrantOptions)
		opts.To = &sdk.ShowGrantsTo{
			Application: application2.ID(),
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
