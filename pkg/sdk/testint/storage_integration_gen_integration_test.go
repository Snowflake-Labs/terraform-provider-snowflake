package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt_StorageIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	if !hasExternalEnvironmentVariablesSet {
		t.Skip("Skipping TestInt_StorageIntegrations (External environmental variables are not set)")
	}

	t.Run("Create - s3", func(t *testing.T) {
		id := sdk.RandomAccountObjectIdentifier()
		err := client.StorageIntegrations.Create(ctx, sdk.NewCreateStorageIntegrationRequest(id, true, []sdk.StorageLocation{
			{
				Path: awsBucketUrl + "/allowed-location",
			},
		}))
		require.NoError(t, err)

	})

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})
}
