package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/require"
)

func TestInt_Streamlits(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupStreamlitHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Streamlits.Drop(ctx, sdk.NewDropStreamlitRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createStreamlitHandle := func(t *testing.T, stage *sdk.Stage, mainFile string) *sdk.Streamlit {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		request := sdk.NewCreateStreamlitRequest(id, "@"+stage.ID().FullyQualifiedName(), mainFile)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		e, err := client.Streamlits.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	assertStreamlit := func(t *testing.T, id sdk.SchemaObjectIdentifier, comment string, warehouse string) {
		t.Helper()

		e, err := client.Streamlits.ShowByID(ctx, id)
		require.NoError(t, err)

		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, id.DatabaseName(), e.DatabaseName)
		require.Equal(t, id.SchemaName(), e.SchemaName)
		require.Empty(t, e.Title)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, warehouse, e.QueryWarehouse)
		require.NotEmpty(t, e.UrlId)
		require.Equal(t, "ROLE", e.OwnerRoleType)
	}

	t.Run("create streamlit", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)

		comment := random.StringN(4)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, "@"+stage.ID().FullyQualifiedName(), mainFile).WithComment(&comment)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		assertStreamlit(t, id, comment, "")
	})

	t.Run("alter streamlit: set", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)
		manifest := "manifest.yml"
		e := createStreamlitHandle(t, stage, manifest)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		comment := random.StringN(4)
		set := sdk.NewStreamlitSetRequest(sdk.String("@"+stage.ID().FullyQualifiedName()), &manifest).WithComment(&comment)
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithSet(set))
		require.NoError(t, err)
		assertStreamlit(t, id, comment, "")
	})

	t.Run("alter function: rename", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml")

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(3))
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupStreamlitHandle(id))
		} else {
			t.Cleanup(cleanupStreamlitHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Streamlits.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		o, err := client.Streamlits.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), o.Name)
	})

	t.Run("show streamliit: with like", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml")

		streamlits, err := client.Streamlits.Show(ctx, sdk.NewShowStreamlitRequest().WithLike(&sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(streamlits))
		require.Equal(t, *e, streamlits[0])
	})

	t.Run("describe streamlit", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, databaseTest, schemaTest, "dev_stage_test")
		t.Cleanup(cleanupStage)

		mainFile := "manifest.yml"
		e := createStreamlitHandle(t, stage, mainFile)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		detail, err := client.Streamlits.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, e.Name, detail.Name)
		require.Equal(t, e.UrlId, detail.UrlId)
		require.Equal(t, mainFile, detail.MainFile)
		require.Equal(t, "@"+stage.ID().FullyQualifiedName(), detail.RootLocation)
		require.Empty(t, detail.Title)
		require.Empty(t, detail.QueryWarehouse)
	})
}
