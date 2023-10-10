package sdk_integration_tests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GetTag(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schemaTest, schemaCleanup := sdk.createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	tagTest, tagCleanup := sdk.createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	t.Run("masking policy tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(maskingPolicyCleanup)

		tagValue := sdk.randomString(t)
		err := client.MaskingPolicies.Alter(ctx, maskingPolicyTest.ID(), &sdk.AlterMaskingPolicyOptions{
			Set: &sdk.MaskingPolicySet{
				Tag: []sdk.TagAssociation{
					{
						Name:  tagTest.ID(),
						Value: tagValue,
					},
				},
			},
		})
		require.NoError(t, err)
		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Equal(t, tagValue, s)
	})

	t.Run("masking policy with no set tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := sdk.createMaskingPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(maskingPolicyCleanup)

		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.Error(t, err)
		assert.Equal(t, "", s)
	})
}
