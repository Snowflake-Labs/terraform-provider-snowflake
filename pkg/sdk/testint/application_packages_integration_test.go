package testint

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
todo: add tests for:
  - Creates a custom release directive for the specified accounts : https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-directive
  - Create application package with insufficient privileges for the following three fields
	-  WithDataRetentionTimeInDays(sdk.Int(1)).
	-  WithMaxDataExtensionTimeInDays(sdk.Int(1)).
	-  WithDefaultDdlCollation(sdk.String("en_US")).
*/

func TestInt_ApplicationPackages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)
	tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	cleanupApplicationPackageHandle := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ApplicationPackages.Drop(ctx, sdk.NewDropApplicationPackageRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createApplicationPackageHandle := func(t *testing.T) *sdk.ApplicationPackage {
		t.Helper()

		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		request := sdk.NewCreateApplicationPackageRequest(id).WithDistribution(sdk.DistributionPointer(sdk.DistributionInternal))
		err := client.ApplicationPackages.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationPackageHandle(id))

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	assertApplicationPackage := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterDataRetentionTimeInDays)
		require.NoError(t, err)

		defaultDataRetentionTimeInDays, err := strconv.Atoi(param.Value)
		require.NoError(t, err)

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.NotEmpty(t, e.CreatedOn)
		assert.Equal(t, id.Name(), e.Name)
		assert.Equal(t, false, e.IsDefault)
		assert.Equal(t, true, e.IsCurrent)
		assert.Equal(t, sdk.DistributionInternal, sdk.Distribution(e.Distribution))
		assert.Equal(t, "ACCOUNTADMIN", e.Owner)
		assert.Empty(t, e.Comment)
		assert.Equal(t, defaultDataRetentionTimeInDays, e.RetentionTime)
		assert.Empty(t, e.Options)
		assert.Empty(t, e.DroppedOn)
		assert.Empty(t, e.ApplicationClass)
	}

	t.Run("create application package", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier(random.StringN(4))
		comment := random.StringN(4)
		request := sdk.NewCreateApplicationPackageRequest(id).
			WithComment(&comment).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			}).
			WithDistribution(sdk.DistributionPointer(sdk.DistributionExternal))
		err := client.ApplicationPackages.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationPackageHandle(id))

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, sdk.DistributionExternal, sdk.Distribution(e.Distribution))
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterDataRetentionTimeInDays)
		require.NoError(t, err)

		defaultDataRetentionTimeInDays, err := strconv.Atoi(param.Value)
		require.NoError(t, err)

		require.Equal(t, defaultDataRetentionTimeInDays, e.RetentionTime)
	})

	t.Run("alter application package: set", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		distribution := sdk.DistributionPointer(sdk.DistributionExternal)
		set := sdk.NewApplicationPackageSetRequest().
			WithDistribution(distribution).
			WithComment(sdk.String("test")).
			WithDataRetentionTimeInDays(sdk.Int(2)).
			WithMaxDataExtensionTimeInDays(sdk.Int(2)).
			WithDefaultDdlCollation(sdk.String("utf8mb4_0900_ai_ci"))
		err := client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithSet(set))
		require.NoError(t, err)

		o, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, *distribution, sdk.Distribution(o.Distribution))
		assert.Equal(t, 2, o.RetentionTime)
		assert.Equal(t, "test", o.Comment)
	})

	t.Run("alter application package: unset", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		// unset comment and distribution
		unset := sdk.NewApplicationPackageUnsetRequest().WithComment(sdk.Bool(true)).WithDistribution(sdk.Bool(true))
		err := client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithUnset(unset))
		require.NoError(t, err)
		o, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Empty(t, o.Comment)
		require.Equal(t, sdk.DistributionInternal, sdk.Distribution(o.Distribution))
	})

	t.Run("alter application package: set and unset tags", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		id := sdk.NewAccountObjectIdentifier(e.Name)

		setTags := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "v1",
			},
		}
		err := client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithSetTags(setTags))
		require.NoError(t, err)
		assertApplicationPackage(t, id)

		unsetTags := []sdk.ObjectIdentifier{
			tagTest.ID(),
		}
		err = client.ApplicationPackages.Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithUnsetTags(unsetTags))
		require.NoError(t, err)
		assertApplicationPackage(t, id)
	})

	t.Run("show application package for SQL: with like", func(t *testing.T) {
		e := createApplicationPackageHandle(t)

		packages, err := client.ApplicationPackages.Show(ctx, sdk.NewShowApplicationPackageRequest().WithLike(&sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(packages))
		require.Equal(t, *e, packages[0])
	})
}

type StagedFile struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type ApplicationPackageVersion struct {
	Version string `json:"version"`
	Patch   int    `json:"patch"`
}

func TestInt_ApplicationPackagesVersionAndReleaseDirective(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupApplicationPackageHandle := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ApplicationPackages.Drop(ctx, sdk.NewDropApplicationPackageRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createApplicationPackageHandle := func(t *testing.T) *sdk.ApplicationPackage {
		t.Helper()

		id := sdk.NewAccountObjectIdentifier("snowflake_package_test")
		request := sdk.NewCreateApplicationPackageRequest(id).WithDistribution(sdk.DistributionPointer(sdk.DistributionInternal))
		err := client.ApplicationPackages.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApplicationPackageHandle(id))

		// grant role "ACCOUNTADMIN" on application package
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`GRANT MANAGE VERSIONS ON APPLICATION PACKAGE "%s" TO ROLE ACCOUNTADMIN;`, id.Name()))
		require.NoError(t, err)

		e, err := client.ApplicationPackages.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	createStageHandle := func(t *testing.T) *sdk.Stage {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, "stage_test")
		co := sdk.NewStageCopyOptionsRequest().WithOnError(sdk.NewStageCopyOnErrorOptionsRequest().WithSkipFile())
		cr := sdk.NewCreateInternalStageRequest(id).WithCopyOptions(co)
		err := client.Stages.CreateInternal(ctx, cr)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		e, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	putOnStageHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier, name string) {
		t.Helper()

		tempFile := fmt.Sprintf("/tmp/%s", name)
		f, err := os.Create(tempFile)
		require.NoError(t, err)
		f.Close()
		defer os.Remove(name)

		_, err = client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @%s AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, f.Name(), id.FullyQualifiedName()))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err = client.ExecForTests(ctx, fmt.Sprintf(`REMOVE @%s/%s`, id.FullyQualifiedName(), name))
			require.NoError(t, err)
		})
	}

	showApplicationPackageVersion := func(t *testing.T, name string) []ApplicationPackageVersion {
		t.Helper()

		var versions []ApplicationPackageVersion
		err := client.QueryForTests(ctx, &versions, fmt.Sprintf(`SHOW VERSIONS IN APPLICATION PACKAGE "%s"`, name))
		require.NoError(t, err)
		return versions
	}

	t.Run("alter application package: add, patch and drop version", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		s := createStageHandle(t)
		putOnStageHandle(t, s.ID(), "manifest.yml")
		putOnStageHandle(t, s.ID(), "setup.sql")

		version := "V001"
		using := "@" + s.ID().FullyQualifiedName()
		// add version to application package
		id := sdk.NewAccountObjectIdentifier(e.Name)
		vr := sdk.NewAddVersionRequest(using).WithVersionIdentifier(&version).WithLabel(sdk.String("add version V001"))
		r1 := sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(vr)
		err := client.ApplicationPackages.Alter(ctx, r1)
		require.NoError(t, err)
		versions := showApplicationPackageVersion(t, e.Name)
		require.Equal(t, 1, len(versions))
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)

		// add patch for application package version
		pr := sdk.NewAddPatchForVersionRequest(&version, using).WithLabel(sdk.String("patch version V001"))
		r2 := sdk.NewAlterApplicationPackageRequest(id).WithAddPatchForVersion(pr)
		err = client.ApplicationPackages.Alter(ctx, r2)
		require.NoError(t, err)
		versions = showApplicationPackageVersion(t, e.Name)
		require.Equal(t, 2, len(versions))
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)
		require.Equal(t, version, versions[1].Version)
		require.Equal(t, 1, versions[1].Patch)

		// drop version from application package
		r3 := sdk.NewAlterApplicationPackageRequest(id).WithDropVersion(sdk.NewDropVersionRequest(version))
		err = client.ApplicationPackages.Alter(ctx, r3)
		require.NoError(t, err)
		versions = showApplicationPackageVersion(t, e.Name)
		require.Equal(t, 0, len(versions))
	})

	t.Run("alter application package: set default release directive", func(t *testing.T) {
		e := createApplicationPackageHandle(t)
		s := createStageHandle(t)
		putOnStageHandle(t, s.ID(), "manifest.yml")
		putOnStageHandle(t, s.ID(), "setup.sql")

		version := "V001"
		using := "@" + s.ID().FullyQualifiedName()
		// add version to application package
		id := sdk.NewAccountObjectIdentifier(e.Name)
		vr := sdk.NewAddVersionRequest(using).WithVersionIdentifier(&version).WithLabel(sdk.String("add version V001"))
		r1 := sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(vr)
		err := client.ApplicationPackages.Alter(ctx, r1)
		require.NoError(t, err)
		versions := showApplicationPackageVersion(t, e.Name)
		require.Equal(t, 1, len(versions))
		require.Equal(t, version, versions[0].Version)
		require.Equal(t, 0, versions[0].Patch)

		// set default release directive
		rr := sdk.NewSetDefaultReleaseDirectiveRequest(version, 0)
		r2 := sdk.NewAlterApplicationPackageRequest(id).WithSetDefaultReleaseDirective(rr)
		err = client.ApplicationPackages.Alter(ctx, r2)
		require.NoError(t, err)
	})
}
