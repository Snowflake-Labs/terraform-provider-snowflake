package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * todo: `ALTER SEQUENCE [ IF EXISTS ] <name> UNSET COMMENT` not works, and error: Syntax error: unexpected 'COMMENT'. (line 39)
 */

func TestInt_Sequences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Sequences.Drop(ctx, sdk.NewDropSequenceRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createSequenceHandle := func(t *testing.T) *sdk.Sequence {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		sr := sdk.NewCreateSequenceRequest(id).WithStart(sdk.Int(1)).WithIncrement(sdk.Int(1))
		err := client.Sequences.Create(ctx, sr)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))

		s, err := client.Sequences.ShowByID(ctx, id)
		require.NoError(t, err)
		return s
	}

	assertSequence := func(t *testing.T, id sdk.SchemaObjectIdentifier, interval int, ordered bool, comment string) {
		t.Helper()

		e, err := client.Sequences.ShowByID(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, id.DatabaseName(), e.DatabaseName)
		require.Equal(t, id.SchemaName(), e.SchemaName)
		require.Equal(t, 1, e.NextValue)
		require.Equal(t, interval, e.Interval)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, "ROLE", e.OwnerRoleType)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, ordered, e.Ordered)
	}

	t.Run("create sequence", func(t *testing.T) {
		name := random.StringN(4)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		comment := random.StringN(4)
		request := sdk.NewCreateSequenceRequest(id).
			WithStart(sdk.Int(1)).
			WithIncrement(sdk.Int(1)).
			WithIfNotExists(sdk.Bool(true)).
			WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehaviorOrder)).
			WithComment(&comment)
		err := client.Sequences.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))
		assertSequence(t, id, 1, true, comment)
	})

	t.Run("show event table: without like", func(t *testing.T) {
		e1 := createSequenceHandle(t)
		e2 := createSequenceHandle(t)

		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest())
		require.NoError(t, err)
		require.Equal(t, 2, len(sequences))
		require.Contains(t, sequences, *e1)
		require.Contains(t, sequences, *e2)
	})

	t.Run("show sequence: with like", func(t *testing.T) {
		e1 := createSequenceHandle(t)
		e2 := createSequenceHandle(t)

		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest().WithLike(&sdk.Like{Pattern: &e1.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(sequences))
		require.Contains(t, sequences, *e1)
		require.NotContains(t, sequences, *e2)
	})

	t.Run("show sequence: no matches", func(t *testing.T) {
		sequences, err := client.Sequences.Show(ctx, sdk.NewShowSequenceRequest().WithLike(&sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		require.Equal(t, 0, len(sequences))
	})

	t.Run("describe sequence", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		details, err := client.Sequences.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, e.CreatedOn, details.CreatedOn)
		require.Equal(t, e.Name, details.Name)
		require.Equal(t, e.SchemaName, details.SchemaName)
		require.Equal(t, e.DatabaseName, details.DatabaseName)
		require.Equal(t, e.NextValue, details.NextValue)
		require.Equal(t, e.Interval, details.Interval)
		require.Equal(t, e.Owner, details.Owner)
		require.Equal(t, e.OwnerRoleType, details.OwnerRoleType)
		require.Equal(t, e.Comment, details.Comment)
		require.Equal(t, e.Ordered, details.Ordered)
	})

	t.Run("alter sequence: set options", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		comment := random.StringN(4)
		set := sdk.NewSequenceSetRequest().WithComment(&comment).WithValuesBehavior(sdk.ValuesBehaviorPointer(sdk.ValuesBehaviorNoOrder))
		err := client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithSet(set))
		require.NoError(t, err)

		assertSequence(t, id, 1, false, comment)
	})

	t.Run("alter sequence: set increment", func(t *testing.T) {
		e := createSequenceHandle(t)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		increment := 2
		err := client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithSetIncrement(&increment))
		require.NoError(t, err)
		assertSequence(t, id, 2, false, "")
	})

	t.Run("alter sequence: rename", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.Sequences.Create(ctx, sdk.NewCreateSequenceRequest(id))
		require.NoError(t, err)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		err = client.Sequences.Alter(ctx, sdk.NewAlterSequenceRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupSequenceHandle(t, id))
		} else {
			t.Cleanup(cleanupSequenceHandle(t, nid))
		}
		require.NoError(t, err)

		_, err = client.Sequences.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
		_, err = client.Sequences.ShowByID(ctx, nid)
		require.NoError(t, err)
	})
}

func TestInt_SequencesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Sequences.Drop(ctx, sdk.NewDropSequenceRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createSequenceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		sr := sdk.NewCreateSequenceRequest(id).WithStart(sdk.Int(1)).WithIncrement(sdk.Int(1))
		err := client.Sequences.Create(ctx, sr)
		require.NoError(t, err)
		t.Cleanup(cleanupSequenceHandle(t, id))
	}

	t.Run("show by id", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, client, databaseTest, random.AlphaN(8))
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createSequenceHandle(t, id1)
		createSequenceHandle(t, id2)

		e1, err := client.Sequences.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Sequences.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
