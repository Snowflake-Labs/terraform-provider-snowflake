package tracking

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/stretchr/testify/require"
)

func Test_Context(t *testing.T) {
	metadata := NewMetadata("123", resources.Account, CreateOperation)
	newMetadata := NewMetadata("321", resources.Database, UpdateOperation)
	ctx := context.Background()

	// no metadata in context
	value := ctx.Value(metadataContextKey)
	require.Nil(t, value)

	retrievedMetadata, ok := FromContext(ctx)
	require.False(t, ok)
	require.Empty(t, retrievedMetadata)

	// add metadata by hand
	ctx = context.WithValue(ctx, metadataContextKey, metadata)

	value = ctx.Value(metadataContextKey)
	require.NotNil(t, value)
	require.Equal(t, metadata, value)

	retrievedMetadata, ok = FromContext(ctx)
	require.True(t, ok)
	require.Equal(t, metadata, retrievedMetadata)

	// add metadata with NewContext function (overrides previous value)
	ctx = NewContext(ctx, newMetadata)

	value = ctx.Value(metadataContextKey)
	require.NotNil(t, value)
	require.Equal(t, newMetadata, value)

	retrievedMetadata, ok = FromContext(ctx)
	require.True(t, ok)
	require.Equal(t, newMetadata, retrievedMetadata)
}
