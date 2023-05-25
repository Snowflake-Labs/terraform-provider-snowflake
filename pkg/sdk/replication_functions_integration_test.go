package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ShowReplicationFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	accounts, err := client.ReplicationFunctions.ShowReplicationAcccounts(ctx)
	if err != nil {
		t.Skip("replication not enabled in this account")
	}
	assert.NotEmpty(t, accounts)
}

func TestInt_ShowRegions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	t.Run("no options", func(t *testing.T) {
		regions, err := client.ReplicationFunctions.ShowRegions(ctx, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, regions)
	})

	t.Run("with options", func(t *testing.T) {
		regions, err := client.ReplicationFunctions.ShowRegions(ctx, &ShowRegionsOptions{
			Like: &Like{
				Pattern: String("AWS_US_WEST_2"),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(regions))
		region := regions[0]
		assert.Equal(t, "AWS_US_WEST_2", region.SnowflakeRegion)
		assert.Equal(t, CloudTypeAWS, region.CloudType)
		assert.Equal(t, "us-west-2", region.Region)
		assert.Equal(t, "US West (Oregon)", region.DisplayName)
	})
}
