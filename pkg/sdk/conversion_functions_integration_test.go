package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ToTimestampLTZ(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	err := client.Accounts.Alter(ctx, &AlterAccountOptions{
		Set: &AccountSet{
			Parameters: &AccountLevelParameters{
				SessionParameters: &SessionParameters{
					TimestampTypeMapping: String("TIMESTAMP_LTZ"),
					Timezone:             String("UTC"),
				},
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.Accounts.Alter(ctx, &AlterAccountOptions{
			Unset: &AccountUnset{
				Parameters: &AccountLevelParametersUnset{
					SessionParameters: &SessionParametersUnset{
						TimestampTypeMapping: Bool(true),
						Timezone:             Bool(true),
					},
				},
			},
		})
		require.NoError(t, err)
	})
	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	now := time.Now()
	actual, err := client.ConversionFunctions.ToTimestampLTZ(ctx, now)
	require.NoError(t, err)
	expected := now.UTC()
	assert.Equal(t, expected, actual)
}

func TestInt_ToTimestampNTZ(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	err := client.Accounts.Alter(ctx, &AlterAccountOptions{
		Set: &AccountSet{
			Parameters: &AccountLevelParameters{
				SessionParameters: &SessionParameters{
					TimestampTypeMapping: String("TIMESTAMP_NTZ"),
					Timezone:             String("UTC"),
				},
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.Accounts.Alter(ctx, &AlterAccountOptions{
			Unset: &AccountUnset{
				Parameters: &AccountLevelParametersUnset{
					SessionParameters: &SessionParametersUnset{
						TimestampTypeMapping: Bool(true),
						Timezone:             Bool(true),
					},
				},
			},
		})
		require.NoError(t, err)
	})
	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	now := time.Now()
	actual, err := client.ConversionFunctions.ToTimestampNTZ(ctx, now)
	require.NoError(t, err)
	assert.Equal(t, time.UTC, actual.Location())
}
