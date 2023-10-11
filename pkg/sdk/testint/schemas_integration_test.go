package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SchemasCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, cleanupSchema := createSchema(t, client, db)
	t.Cleanup(cleanupSchema)

	t.Run("replace", func(t *testing.T) {
		comment := "replaced"
		err := client.Schemas.Create(ctx, schema.ID(), &sdk.CreateSchemaOptions{
			OrReplace:                  sdk.Bool(true),
			DataRetentionTimeInDays:    sdk.Int(10),
			MaxDataExtensionTimeInDays: sdk.Int(10),
			DefaultDDLCollation:        sdk.String("en_US-trim"),
			WithManagedAccess:          sdk.Bool(true),
			Comment:                    sdk.String(comment),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)
		assert.Equal(t, "MANAGED ACCESS", *s.Options)
		assert.Equal(t, comment, *s.Comment)
	})

	t.Run("if not exists", func(t *testing.T) {
		comment := "some_comment"
		err := client.Schemas.Create(ctx, schema.ID(), &sdk.CreateSchemaOptions{
			IfNotExists: sdk.Bool(true),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.NotEqual(t, comment, *s.Comment)
	})

	t.Run("clone", func(t *testing.T) {
		comment := "some_comment"
		schemaID := sdk.NewDatabaseObjectIdentifier(db.Name, randomAccountObjectIdentifier(t).Name())
		err := client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
		})
		require.NoError(t, err)

		clonedSchemaID := sdk.NewDatabaseObjectIdentifier(db.Name, randomAccountObjectIdentifier(t).Name())
		err = client.Schemas.Create(ctx, clonedSchemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
			Clone: &sdk.Clone{
				SourceObject: schemaID,
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)

		cs, err := client.Schemas.ShowByID(ctx, clonedSchemaID)
		require.NoError(t, err)
		assert.Equal(t, *s.Comment, *cs.Comment)

		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
			err = client.Schemas.Drop(ctx, clonedSchemaID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("with tags", func(t *testing.T) {
		tagName := random.RandomString(t)
		tagID := sdk.NewAccountObjectIdentifier(tagName)
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE TAG "%s"`, tagName))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP TAG "%s"`, tagName))
			require.NoError(t, err)
		})

		schemaID := sdk.NewDatabaseObjectIdentifier(db.Name, randomAccountObjectIdentifier(t).Name())
		tagValue := random.RandomString(t)
		err = client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Tag: []sdk.TagAssociation{
				{
					Name:  tagID,
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})

		tv, err := client.SystemFunctions.GetTag(ctx, tagID, schemaID, sdk.ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, tagValue, tv)
	})
}

func TestInt_SchemasAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	t.Run("rename to", func(t *testing.T) {
		schema, _ := createSchema(t, client, db)
		newID := sdk.NewDatabaseObjectIdentifier(db.Name, random.RandomString(t))
		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			NewName: newID,
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID, s.ID())
		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, newID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("swap with", func(t *testing.T) {
		schema, cleanupSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSchema)

		swapSchema, cleanupSwapSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSwapSchema)

		table, _ := createTable(t, client, db, schema)
		t.Cleanup(func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", db.Name, swapSchema.Name, table.Name))
			require.NoError(t, err)
		})

		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			SwapWith: swapSchema.ID(),
		})
		require.NoError(t, err)

		schemaDetails, err := client.Schemas.Describe(ctx, swapSchema.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(schemaDetails))
		assert.Equal(t, "TABLE", schemaDetails[0].Kind)
		assert.Equal(t, table.Name, schemaDetails[0].Name)
	})

	t.Run("set", func(t *testing.T) {
		schema, cleanupSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSchema)

		comment := random.RandomComment(t)
		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				DataRetentionTimeInDays:    sdk.Int(3),
				MaxDataExtensionTimeInDays: sdk.Int(3),
				DefaultDDLCollation:        sdk.String("en_US-trim"),
				Comment:                    sdk.String(comment),
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, comment, *s.Comment)
	})

	t.Run("unset", func(t *testing.T) {
		schemaID := sdk.NewDatabaseObjectIdentifier(db.Name, random.RandomString(t))
		comment := random.RandomComment(t)
		err := client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Comment: sdk.String(comment),
		})
		require.NoError(t, err)

		err = client.Schemas.Alter(ctx, schemaID, &sdk.AlterSchemaOptions{
			Unset: &sdk.SchemaUnset{
				Comment: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)
		assert.Empty(t, *s.Comment)

		t.Cleanup(func() {
			err := client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})
	})

	t.Run("set tags", func(t *testing.T) {
		schemaID := sdk.NewDatabaseObjectIdentifier(db.Name, random.RandomString(t))
		err := client.Schemas.Create(ctx, schemaID, nil)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})

		s, err := client.Schemas.ShowByID(ctx, schemaID)
		require.NoError(t, err)

		tag, cleanupTag := createTag(t, client, db, s)
		t.Cleanup(cleanupTag)

		tagValue := "tag-value"
		err = client.Schemas.Alter(ctx, schemaID, &sdk.AlterSchemaOptions{
			Set: &sdk.SchemaSet{
				Tag: []sdk.TagAssociation{
					{
						Name:  tag.ID(),
						Value: tagValue,
					},
				},
			},
		})
		require.NoError(t, err)

		tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), s.ID(), sdk.ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, tagValue, tv)
	})

	t.Run("unset tags", func(t *testing.T) {
		tagName := random.RandomString(t)
		tagID := sdk.NewAccountObjectIdentifier(tagName)
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE TAG "%s"`, tagName))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP TAG "%s"`, tagName))
			require.NoError(t, err)
		})

		schemaID := sdk.NewDatabaseObjectIdentifier(db.Name, randomAccountObjectIdentifier(t).Name())
		tagValue := random.RandomString(t)
		err = client.Schemas.Create(ctx, schemaID, &sdk.CreateSchemaOptions{
			Tag: []sdk.TagAssociation{
				{
					Name:  tagID,
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Schemas.Drop(ctx, schemaID, nil)
			require.NoError(t, err)
		})

		err = client.Schemas.Alter(ctx, schemaID, &sdk.AlterSchemaOptions{
			Unset: &sdk.SchemaUnset{
				Tag: []sdk.ObjectIdentifier{
					tagID,
				},
			},
		})
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tagID, schemaID, sdk.ObjectTypeSchema)
		require.Error(t, err)
	})

	t.Run("enable managed access", func(t *testing.T) {
		schema, cleanupSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSchema)

		err := client.Schemas.Alter(ctx, schema.ID(), &sdk.AlterSchemaOptions{
			EnableManagedAccess: sdk.Bool(true),
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, schema.Name, s.Name)
		assert.Equal(t, "MANAGED ACCESS", *s.Options)
	})
}

func TestInt_SchemasShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, cleanupSchema := createSchema(t, client, db)
	t.Cleanup(cleanupSchema)

	t.Run("no options", func(t *testing.T) {
		schemas, err := client.Schemas.Show(ctx, nil)
		require.NoError(t, err)
		schemaNames := make([]string, len(schemas))
		for i, s := range schemas {
			schemaNames[i] = s.Name
		}
		assert.Contains(t, schemaNames, schema.Name)
	})

	t.Run("with options", func(t *testing.T) {
		schemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
			Terse:   sdk.Bool(true),
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(schema.Name),
			},
			In: &sdk.SchemaIn{
				Account: sdk.Bool(true),
			},
			StartsWith: sdk.String(schema.Name),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		})
		require.NoError(t, err)
		schemaNames := make([]string, len(schemas))
		for i, s := range schemas {
			schemaNames[i] = s.Name
		}
		assert.Contains(t, schemaNames, schema.Name)
	})
}

func TestInt_SchemasDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, _ := createSchema(t, client, db)

	s, err := client.Schemas.ShowByID(ctx, schema.ID())
	require.NoError(t, err)
	assert.Equal(t, schema.Name, s.Name)

	err = client.Schemas.Drop(ctx, schema.ID(), nil)
	require.NoError(t, err)

	schemas, err := client.Schemas.Show(ctx, &sdk.ShowSchemaOptions{
		Like: &sdk.Like{
			Pattern: &schema.Name,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(schemas))
}

/*
todo: this test is failing, need to fix
func TestInt_SchemasUndrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, _ := createSchema(t, client, db)

	before, err := client.Schemas.ShowByID(ctx, schema.ID())
	require.NoError(t, err)
	assert.Equal(t, schema.Name, before.Name)

	err = client.Schemas.Drop(ctx, schema.ID(), nil)
	require.NoError(t, err)

	err = client.Schemas.Undrop(ctx, schema.ID())
	require.NoError(t, err)

	after, err := client.Schemas.ShowByID(ctx, schema.ID())
	require.NoError(t, err)
	assert.Equal(t, schema.Name, after.Name)

	t.Cleanup(func() {
		err = client.Schemas.Drop(ctx, schema.ID(), nil)
		require.NoError(t, err)
	})
}
*/
