package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Stages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("CreateInternal", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.DatabaseName(), stage.DatabaseName)
		assert.Equal(t, id.SchemaName(), stage.SchemaName)
		assert.Equal(t, id.Name(), stage.Name)
	})

	t.Run("CreateOnS3", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateOnGCS", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateOnAzure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateOnS3Compatible", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter - rename", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))
		newId := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))
		renamed := false

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			if renamed {
				err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(newId))
				require.NoError(t, err)
			} else {
				err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
				require.NoError(t, err)
			}
		})

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).
			WithIfExists(sdk.Bool(true)).
			WithRenameTo(&newId))
		require.NoError(t, err)
		renamed = true

		stage, err := client.Stages.ShowByID(ctx, newId)
		require.NotNil(t, stage)
		require.NoError(t, err)
	})

	t.Run("Alter - set unset tags", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))
		tag, cleanupTag := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(cleanupTag)

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStage)
		require.Error(t, err)

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: "tag value",
			},
		}))
		require.NoError(t, err)

		value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStage)
		require.NoError(t, err)
		assert.Equal(t, "tag value", value)

		err = client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithUnsetTags([]sdk.ObjectIdentifier{
			tag.ID(),
		}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStage)
		require.Error(t, err)
	})

	t.Run("AlterInternalStage", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithCopyOptions(sdk.NewStageCopyOptionsRequest().WithSizeLimit(sdk.Int(100))).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeJSON)).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "some comment", stage.Comment)

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_COPY_OPTIONS",
			Name:    "SIZE_LIMIT",
			Type:    "Long",
			Value:   "100",
			Default: "",
		})
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_FILE_FORMAT",
			Name:    "TYPE",
			Type:    "String",
			Value:   "JSON",
			Default: "CSV",
		})

		err = client.Stages.AlterInternalStage(ctx, sdk.NewAlterInternalStageStageRequest(id).
			WithIfExists(sdk.Bool(true)).
			WithCopyOptions(sdk.NewStageCopyOptionsRequest().WithSizeLimit(sdk.Int(200))).
			WithFileFormat(sdk.NewStageFileFormatRequest().WithType(&sdk.FileFormatTypeCSV)).
			WithComment(sdk.String("altered comment")))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "altered comment", stage.Comment)

		stageProperties, err = client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_COPY_OPTIONS",
			Name:    "SIZE_LIMIT",
			Type:    "Long",
			Value:   "200",
			Default: "",
		})
		require.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_FILE_FORMAT",
			Name:    "TYPE",
			Type:    "String",
			Value:   "CSV",
			Default: "CSV",
		})
	})

	t.Run("AlterExternalS3Stage", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterExternalGCSStage", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterExternalAzureStage", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterDirectoryTable", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NotNil(t, stage)
		require.NoError(t, err)

		err = client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
		require.NoError(t, err)

		stage, err = client.Stages.ShowByID(ctx, id)
		require.Nil(t, stage)
		require.Error(t, err)
	})

	t.Run("Describe internal", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stageProperties, err := client.Stages.Describe(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, stageProperties)
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "DIRECTORY",
			Name:    "ENABLE",
			Type:    "Boolean",
			Value:   "false",
			Default: "false",
		})
		assert.Contains(t, stageProperties, sdk.StageProperty{
			Parent:  "STAGE_LOCATION",
			Name:    "URL",
			Type:    "String",
			Value:   "",
			Default: "",
		})
	})

	t.Run("Describe external s3", func(t *testing.T) {
	})

	t.Run("Describe external gcs", func(t *testing.T) {
	})

	t.Run("Describe external azure", func(t *testing.T) {
	})

	t.Run("Show internal", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.AlphanumericN(32))

		err := client.Stages.CreateInternal(ctx, sdk.NewCreateInternalStageRequest(id).
			WithDirectoryTableOptions(sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(sdk.Bool(true))).
			WithComment(sdk.String("some comment")))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
			require.NoError(t, err)
		})

		stage, err := client.Stages.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.DatabaseName(), stage.DatabaseName)
		assert.Equal(t, id.SchemaName(), stage.SchemaName)
		assert.Equal(t, id.Name(), stage.Name)
		assert.Empty(t, stage.Url)
		assert.False(t, stage.HasCredentials)
		assert.False(t, stage.HasEncryptionKey)
		assert.Equal(t, "some comment", stage.Comment)
		assert.Nil(t, stage.Region)
		assert.Equal(t, "INTERNAL", stage.Type)
		assert.Nil(t, stage.Cloud)
		assert.Nil(t, stage.StorageIntegration)
		assert.Nil(t, stage.Endpoint)
		assert.True(t, stage.DirectoryEnabled)
	})

	t.Run("Show external s3", func(t *testing.T) {
	})

	t.Run("Show external gcs", func(t *testing.T) {
	})

	t.Run("Show external azure", func(t *testing.T) {
	})
}
