package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/require"
	"testing"
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

		client.Stages.ShowByID(ctx, id)
	})

	t.Run("CreateOnS3", func(t *testing.T) {

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

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("AlterInternalStage", func(t *testing.T) {
		// TODO: fill me
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

	})

	t.Run("Describe", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})
}
