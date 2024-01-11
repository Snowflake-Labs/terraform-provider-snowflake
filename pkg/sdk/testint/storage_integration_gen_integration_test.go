package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func TestInt_StorageIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	if !hasExternalEnvironmentVariablesSet {
		t.Skip("Skipping TestInt_StorageIntegrations (External environmental variables are not set)")
	}

	assertStorageIntegrationShowResult := func(t *testing.T, s *sdk.StorageIntegration, name sdk.AccountObjectIdentifier, comment string) {
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, true, s.Enabled)
		assert.Equal(t, "EXTERNAL_STAGE", s.StorageType)
		assert.Equal(t, "STORAGE", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	findProp := func(t *testing.T, props []sdk.StorageIntegrationProperty, name string) *sdk.StorageIntegrationProperty {
		prop, err := collections.FindOne(props, func(property sdk.StorageIntegrationProperty) bool { return property.Name == name })
		require.NoError(t, err)
		return prop
	}

	assertS3StorageIntegrationDescResult := func(
		t *testing.T,
		props []sdk.StorageIntegrationProperty,
		enabled bool,
		allowedLocations []sdk.StorageLocation,
		blockedLocations []sdk.StorageLocation,
		comment string,
	) {
		allowed := make([]string, len(allowedLocations))
		for i, a := range allowedLocations {
			allowed[i] = a.Path
		}
		blocked := make([]string, len(blockedLocations))
		for i, b := range blockedLocations {
			blocked[i] = b.Path
		}
		assert.Equal(t, "Boolean", findProp(t, props, "ENABLED").Type)
		assert.Equal(t, strconv.FormatBool(enabled), findProp(t, props, "ENABLED").Value)
		assert.Equal(t, "false", findProp(t, props, "ENABLED").Default)
		assert.Equal(t, "S3", findProp(t, props, "STORAGE_PROVIDER").Value)
		assert.Equal(t, strings.Join(allowed, ","), findProp(t, props, "STORAGE_ALLOWED_LOCATIONS").Value)
		assert.Equal(t, strings.Join(blocked, ","), findProp(t, props, "STORAGE_BLOCKED_LOCATIONS").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_IAM_USER_ARN").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_ROLE_ARN").Value)
		assert.NotEmpty(t, findProp(t, props, "STORAGE_AWS_EXTERNAL_ID").Value)
		assert.Equal(t, comment, findProp(t, props, "COMMENT").Value)
	}

	allowedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/allowed-location",
			},
			{
				Path: prefix + "/allowed-location2",
			},
		}
	}
	s3AllowedLocations := allowedLocations(awsBucketUrl)
	gcsAllowedLocations := allowedLocations(gcsBucketUrl)

	blockedLocations := func(prefix string) []sdk.StorageLocation {
		return []sdk.StorageLocation{
			{
				Path: prefix + "/blocked-location",
			},
			{
				Path: prefix + "/blocked-location2",
			},
		}
	}
	s3BlockedLocations := blockedLocations(awsBucketUrl)
	gcsBlockedLocations := blockedLocations(gcsBucketUrl)

	createS3StorageIntegration := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := sdk.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, s3AllowedLocations).
			WithIfNotExists(sdk.Bool(true)).
			WithS3StorageProviderParams(sdk.NewS3StorageParamsRequest(awsRoleARN)).
			WithStorageBlockedLocations(s3BlockedLocations).
			WithComment(sdk.String("some comment"))

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
			require.NoError(t, err)
		})

		return id
	}

	createGCSStorageIntegration := func(t *testing.T) sdk.AccountObjectIdentifier {
		t.Helper()

		id := sdk.RandomAccountObjectIdentifier()
		req := sdk.NewCreateStorageIntegrationRequest(id, true, gcsAllowedLocations).
			WithIfNotExists(sdk.Bool(true)).
			WithGCSStorageProviderParams(sdk.NewGCSStorageParamsRequest()).
			WithStorageBlockedLocations(gcsBlockedLocations).
			WithComment(sdk.String("some comment"))

		err := client.StorageIntegrations.Create(ctx, req)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.StorageIntegrations.Drop(ctx, sdk.NewDropStorageIntegrationRequest(id))
			require.NoError(t, err)
		})

		return id
	}

	t.Run("Create - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Create - GCS", func(t *testing.T) {
		id := createGCSStorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Create - Azure", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		storageIntegration, err := client.StorageIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStorageIntegrationShowResult(t, storageIntegration, id, "some comment")
	})

	t.Run("Alter - set - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		changedS3AllowedLocations := append(s3AllowedLocations, sdk.StorageLocation{Path: awsBucketUrl + "/allowed-location3"})
		changedS3BlockedLocations := append(s3BlockedLocations, sdk.StorageLocation{Path: awsBucketUrl + "/blocked-location3"})
		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithSet(
				sdk.NewStorageIntegrationSetRequest().
					WithSetS3Params(sdk.NewSetS3StorageParamsRequest(awsRoleARN)).
					WithEnabled(true).
					WithStorageAllowedLocations(changedS3AllowedLocations).
					WithStorageBlockedLocations(changedS3BlockedLocations).
					WithComment(sdk.String("changed comment")),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, props, true, changedS3AllowedLocations, changedS3BlockedLocations, "changed comment")
	})

	t.Run("Alter - set - Azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter - unset", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		req := sdk.NewAlterStorageIntegrationRequest(id).
			WithUnset(
				sdk.NewStorageIntegrationUnsetRequest().
					WithEnabled(sdk.Bool(true)).
					WithStorageBlockedLocations(sdk.Bool(true)).
					WithComment(sdk.Bool(true)),
			)
		err := client.StorageIntegrations.Alter(ctx, req)
		require.NoError(t, err)

		props, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, props, false, s3AllowedLocations, []sdk.StorageLocation{}, "")
	})

	t.Run("Alter - set and unset tags", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		err := client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).
			WithSetTags([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "tag-value",
				},
			}))
		require.NoError(t, err)

		tagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.NoError(t, err)

		assert.Equal(t, "tag-value", tagValue)

		err = client.StorageIntegrations.Alter(ctx, sdk.NewAlterStorageIntegrationRequest(id).
			WithUnsetTags([]sdk.ObjectIdentifier{
				tag.ID(),
			}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.Error(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe - S3", func(t *testing.T) {
		id := createS3StorageIntegration(t)

		desc, err := client.StorageIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertS3StorageIntegrationDescResult(t, desc, true, s3AllowedLocations, s3BlockedLocations, "some comment")
	})

	t.Run("Describe - GCS", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe - Azure", func(t *testing.T) {
		// TODO: fill me
	})
}
