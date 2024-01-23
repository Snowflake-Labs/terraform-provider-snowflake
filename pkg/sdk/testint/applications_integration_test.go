package testint

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * todo: (SNOW-1015095) add integration test for `ALTER APPLICATION <name> UPGRADE`
 * todo: (SNOW-1016268) ALTER APPLICATION [ IF EXISTS ] <name> SET [ SHARE_EVENTS_WITH_PROVIDER ]
 *       attention: SHARE_EVENTS_WITH_PROVIDER can only be set/unset if the application is created in a different account from the application package
 */

func TestInt_Applications(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)
	tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	cleanupApplicationHandle := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.Applications.Drop(ctx, sdk.NewDropApplicationRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createApplicationPackageHandle := func(t *testing.T, applicationPackageName, version string, patch int, defaultReleaseDirective bool) *sdk.Stage {
		t.Helper()

		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)
		putOnStageWithContent(t, client, stage.ID(), "manifest.yml", "")
		putOnStageWithContent(t, client, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")
		cleanupApplicationPackage := createApplicationPackage(t, client, applicationPackageName)
		t.Cleanup(cleanupApplicationPackage)
		addApplicationPackageVersion(t, client, stage, applicationPackageName, version)

		// set default release directive for application package
		if defaultReleaseDirective {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE "%s" SET DEFAULT RELEASE DIRECTIVE VERSION = %s PATCH = %d`, applicationPackageName, version, patch))
			require.NoError(t, err)
		}
		return stage
	}

	createApplicationHandle := func(t *testing.T, applicationPackageName, version string, patch int, versionDirectory bool, debug bool, addPatch bool) (*sdk.Stage, *sdk.Application) {
		t.Helper()

		stage := createApplicationPackageHandle(t, applicationPackageName, version, patch, false)

		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		vr := sdk.NewApplicationVersionRequest().WithVersionAndPatch(sdk.NewVersionAndPatchRequest(version, &patch))
		if versionDirectory {
			vr = sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + stage.ID().FullyQualifiedName()))
		}
		request := sdk.NewCreateApplicationRequest(id, sdk.NewAccountObjectIdentifier(applicationPackageName)).WithVersion(vr).WithDebugMode(&debug)
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		if addPatch {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE "%s" ADD PATCH FOR VERSION %s USING '@%s'`, applicationPackageName, version, stage.ID().FullyQualifiedName()))
			require.NoError(t, err)
		}

		application, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		return stage, application
	}

	assertApplication := func(t *testing.T, id sdk.AccountObjectIdentifier, applicationPackageName, version string, patch int, comment string) {
		t.Helper()

		e, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.NotEmpty(t, e.CreatedOn)
		assert.Equal(t, id.Name(), e.Name)
		assert.Equal(t, false, e.IsDefault)
		assert.Equal(t, true, e.IsCurrent)
		assert.Equal(t, "APPLICATION PACKAGE", e.SourceType)
		assert.Equal(t, applicationPackageName, e.Source)
		assert.Equal(t, version, e.Version)
		assert.Equal(t, patch, e.Patch)
		assert.Equal(t, "ACCOUNTADMIN", e.Owner)
		assert.Equal(t, comment, e.Comment)
		assert.Equal(t, 1, e.RetentionTime)
		assert.Empty(t, e.Options)
	}

	t.Run("create application: without version", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		createApplicationPackageHandle(t, applicationPackageName, version, patch, true)

		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		pid := sdk.NewAccountObjectIdentifier(applicationPackageName)
		comment := random.StringN(4)
		request := sdk.NewCreateApplicationRequest(id, pid).
			WithComment(&comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			})
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		e, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, "APPLICATION PACKAGE", e.SourceType)
		require.Equal(t, applicationPackageName, e.Source)
		require.Equal(t, version, e.Version)
		require.Equal(t, patch, e.Patch)
	})

	t.Run("create application: version and patch", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		createApplicationPackageHandle(t, applicationPackageName, version, patch, false)

		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		pid := sdk.NewAccountObjectIdentifier(applicationPackageName)
		vr := sdk.NewApplicationVersionRequest().WithVersionAndPatch(sdk.NewVersionAndPatchRequest(version, &patch))
		comment := random.StringN(4)
		request := sdk.NewCreateApplicationRequest(id, pid).
			WithDebugMode(sdk.Bool(true)).
			WithComment(&comment).
			WithVersion(vr)
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		assertApplication(t, id, applicationPackageName, version, patch, comment)
	})

	t.Run("create application: version directory", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		stage := createApplicationPackageHandle(t, applicationPackageName, version, patch, false)

		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		pid := sdk.NewAccountObjectIdentifier(applicationPackageName)
		vr := sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + stage.ID().FullyQualifiedName()))
		comment := random.StringN(4)
		request := sdk.NewCreateApplicationRequest(id, pid).
			WithDebugMode(sdk.Bool(true)).
			WithComment(&comment).
			WithVersion(vr).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			})
		err := client.Applications.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationHandle(id))

		assertApplication(t, id, applicationPackageName, "UNVERSIONED", patch, comment)
	})

	t.Run("show application: with like", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		packages, err := client.Applications.Show(ctx, sdk.NewShowApplicationRequest().WithLike(&sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(packages))
		require.Equal(t, *e, packages[0])
	})

	t.Run("alter application: set", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		comment, mode := random.StringN(4), true
		set := sdk.NewApplicationSetRequest().
			WithComment(&comment).
			WithDebugMode(&mode)
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithSet(set))
		require.NoError(t, err)

		details, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, e.SourceType, pairs["source_type"])
		require.Equal(t, e.Source, pairs["source"])
		require.Equal(t, e.Version, pairs["version"])
		require.Equal(t, strconv.Itoa(e.Patch), pairs["patch"])
		require.Equal(t, comment, pairs["comment"])
		require.Equal(t, strconv.FormatBool(mode), pairs["debug_mode"])
	})

	t.Run("alter application: unset", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		unset := sdk.NewApplicationUnsetRequest().WithComment(sdk.Bool(true)).WithDebugMode(sdk.Bool(true))
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnset(unset))
		require.NoError(t, err)

		o, err := client.Applications.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Empty(t, o.Comment)

		details, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, strconv.FormatBool(false), pairs["debug_mode"])
	})

	t.Run("alter application: set and unset tags", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		setTags := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "v1",
			},
		}
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithSetTags(setTags))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackageName, version, patch, "")

		// TODO: 391801 (0A000): SQL compilation error: Object tagging not supported for object type APPLICATION.
		// tv, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), id, sdk.ObjectTypeApplication)
		// require.NoError(t, err)
		// assert.Equal(t, "v1", tv)

		unsetTags := []sdk.ObjectIdentifier{
			tagTest.ID(),
		}
		err = client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnsetTags(unsetTags))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackageName, version, patch, "")
	})

	t.Run("alter application: upgrade with version and patch", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, true)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		vr := sdk.NewVersionAndPatchRequest(version, sdk.Int(patch+1))
		av := sdk.NewApplicationVersionRequest().WithVersionAndPatch(vr)
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUpgradeVersion(av))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackageName, version, patch+1, "")
	})

	t.Run("alter application: upgrade with version directory", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		s, e := createApplicationHandle(t, applicationPackageName, version, patch, true, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		av := sdk.NewApplicationVersionRequest().WithVersionDirectory(sdk.String("@" + s.ID().FullyQualifiedName()))
		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUpgradeVersion(av))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackageName, "UNVERSIONED", patch+1, "")
	})

	t.Run("alter application: unset references", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		err := client.Applications.Alter(ctx, sdk.NewAlterApplicationRequest(id).WithUnsetReferences(sdk.NewApplicationReferencesRequest()))
		require.NoError(t, err)
		assertApplication(t, id, applicationPackageName, version, patch, "")
	})

	t.Run("describe application", func(t *testing.T) {
		applicationPackageName, version, patch := "hello_snowflake_package_test", "V001", 0
		_, e := createApplicationHandle(t, applicationPackageName, version, patch, false, true, false)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		properties, err := client.Applications.Describe(ctx, id)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, item := range properties {
			pairs[item.Property] = item.Value
		}
		require.Equal(t, e.SourceType, pairs["source_type"])
		require.Equal(t, e.Source, pairs["source"])
		require.Equal(t, e.Version, pairs["version"])
		require.Equal(t, e.Label, pairs["version_label"])
		require.Equal(t, e.Comment, pairs["comment"])
		require.Equal(t, strconv.Itoa(e.Patch), pairs["patch"])
	})
}
