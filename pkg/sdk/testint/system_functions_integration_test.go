package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GetTag(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	tagTest, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
	t.Cleanup(tagCleanup)

	t.Run("masking policy tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := createMaskingPolicy(t, client, testDb(t), testSchema(t))
		t.Cleanup(maskingPolicyCleanup)

		tagValue := random.String()
		err := client.MaskingPolicies.Alter(ctx, maskingPolicyTest.ID(), &sdk.AlterMaskingPolicyOptions{
			SetTag: []sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)
		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Equal(t, tagValue, s)
	})

	t.Run("masking policy with no set tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := createMaskingPolicy(t, client, testDb(t), testSchema(t))
		t.Cleanup(maskingPolicyCleanup)

		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.Error(t, err)
		assert.Equal(t, "", s)
	})
}

func TestInt_PipeStatus(t *testing.T) {
	client := testClient(t)

	schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), random.AlphaN(20))
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphaN(20)))
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)
	pipe, pipeCleanup := createPipe(t, client, testDb(t), testSchema(t), random.AlphaN(20), copyStatement)
	t.Cleanup(pipeCleanup)

	pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

	// Pause the pipe
	ctx := context.Background()
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

	// Unpause the pipe
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(false),
		},
	})
	require.NoError(t, err)

	pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)
}

func TestInt_PipeForceResume(t *testing.T) {

}
