package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt_TagCreate(t *testing.T) {
	client := testClient(t)

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	_, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	ctx := context.Background()
	t.Run("create with comment", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		comment := randomComment(t)
		err := client.Tags.Create(ctx, NewCreateTagRequest(name).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Tags.Drop(ctx, NewDropTagRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.Tags.Show(ctx, NewShowTagRequest().WithLike(name.Name()))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, comment, entity.Comment)
	})

	t.Run("create with one allowed value", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		values := []string{"value1"}
		err := client.Tags.Create(ctx, NewCreateTagRequest(name).WithOrReplace(true).WithAllowedValues(values))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Tags.Drop(ctx, NewDropTagRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.Tags.Show(ctx, NewShowTagRequest().WithLike(name.Name()))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, values, entity.AllowedValues)
	})

	t.Run("create with two allowed values", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		values := []string{"value1", "value2"}
		err := client.Tags.Create(ctx, NewCreateTagRequest(name).WithOrReplace(true).WithAllowedValues(values))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Tags.Drop(ctx, NewDropTagRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.Tags.Show(ctx, NewShowTagRequest().WithLike(name.Name()))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, values, entity.AllowedValues)
	})

	t.Run("create with comment and allowed values", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		comment := randomComment(t)
		values := []string{"value1"}
		err := client.Tags.Create(ctx, NewCreateTagRequest(name).WithOrReplace(true).WithComment(&comment).WithAllowedValues(values))
		expected := "fields [Comment AllowedValues] are incompatible and cannot be set at once"
		require.Equal(t, expected, err.Error())
	})
}
