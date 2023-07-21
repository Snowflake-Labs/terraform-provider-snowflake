package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SchemasCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, cleanupSchema := createSchema(t, client, db)
	t.Cleanup(cleanupSchema)

	t.Run("replace", func(t *testing.T) {
		comment := "replaced"
		err := client.Schemas.Create(ctx, schema.ID(), &CreateSchemaOptions{
			OrReplace:                  Bool(true),
			DataRetentionTimeInDays:    Int(10),
			MaxDataExtensionTimeInDays: Int(10),
			DefaultDDLCollation:        String("en_US-trim"),
			WithManagedAccess:          Bool(true),
			Comment:                    String(comment),
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
		err := client.Schemas.Create(ctx, schema.ID(), &CreateSchemaOptions{
			IfNotExists: Bool(true),
			Comment:     String(comment),
		})
		require.NoError(t, err)
		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.NotEqual(t, comment, *s.Comment)
	})

	t.Run("clone", func(t *testing.T) {
		comment := "some_comment"
		schemaID := NewSchemaIdentifier(db.Name, randomAccountObjectIdentifier(t).name)
		err := client.Schemas.Create(ctx, schemaID, &CreateSchemaOptions{
			Comment: String(comment),
		})
		require.NoError(t, err)

		clonedSchemaID := NewSchemaIdentifier(db.Name, randomAccountObjectIdentifier(t).name)
		err = client.Schemas.Create(ctx, clonedSchemaID, &CreateSchemaOptions{
			Comment: String(comment),
			Clone: &Clone{
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
		tagName := randomString(t)
		tagID := NewAccountObjectIdentifier(tagName)
		_, err := client.exec(ctx, fmt.Sprintf(`CREATE TAG "%s"`, tagName))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.exec(ctx, fmt.Sprintf(`DROP TAG "%s"`, tagName))
			require.NoError(t, err)
		})

		schemaID := NewSchemaIdentifier(db.Name, randomAccountObjectIdentifier(t).name)
		tagValue := randomString(t)
		err = client.Schemas.Create(ctx, schemaID, &CreateSchemaOptions{
			Tag: []TagAssociation{
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

		tv, err := client.SystemFunctions.GetTag(ctx, tagID, schemaID, ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, tagValue, tv)
	})
}

func TestInt_SchemasAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	t.Run("rename to", func(t *testing.T) {
		schema, _ := createSchema(t, client, db)
		newID := NewSchemaIdentifier(db.Name, randomString(t))
		err := client.Schemas.Alter(ctx, schema.ID(), &AlterSchemaOptions{
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
		err := client.Schemas.Alter(ctx, schema.ID(), &AlterSchemaOptions{
			SwapWith: swapSchema.ID(),
		})
		require.NoError(t, err)

		schemaDetails, err := client.Schemas.Describe(ctx, swapSchema.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, len(schemaDetails))
		assert.Equal(t, "TABLE", schemaDetails[0].Kind)
		assert.Equal(t, table.Name, schemaDetails[0].Name)

		t.Cleanup(func() {
			_, err := client.exec(ctx, fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", db.Name, swapSchema.Name, table.Name))
			require.NoError(t, err)
		})
	})

	t.Run("set", func(t *testing.T) {
		schema, cleanupSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSchema)

		comment := randomComment(t)
		err := client.Schemas.Alter(ctx, schema.ID(), &AlterSchemaOptions{
			Set: &SchemaSet{
				DataRetentionTimeInDays:    Int(3),
				MaxDataExtensionTimeInDays: Int(3),
				DefaultDDLCollation:        String("en_US-trim"),
				Comment:                    String(comment),
			},
		})
		require.NoError(t, err)

		s, err := client.Schemas.ShowByID(ctx, schema.ID())
		require.NoError(t, err)
		assert.Equal(t, comment, *s.Comment)
	})

	t.Run("unset", func(t *testing.T) {
		schemaID := NewSchemaIdentifier(db.Name, randomString(t))
		comment := randomComment(t)
		err := client.Schemas.Create(ctx, schemaID, &CreateSchemaOptions{
			Comment: String(comment),
		})
		require.NoError(t, err)

		err = client.Schemas.Alter(ctx, schemaID, &AlterSchemaOptions{
			Unset: &SchemaUnset{
				Comment: Bool(true),
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
		schemaID := NewSchemaIdentifier(db.Name, randomString(t))
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
		err = client.Schemas.Alter(ctx, schemaID, &AlterSchemaOptions{
			Set: &SchemaSet{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: tagValue,
					},
				},
			},
		})
		require.NoError(t, err)

		tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), s.ID(), ObjectTypeSchema)
		require.NoError(t, err)
		assert.Equal(t, tagValue, tv)
	})

	t.Run("unset tags", func(t *testing.T) {
		tagName := randomString(t)
		tagID := NewAccountObjectIdentifier(tagName)
		_, err := client.exec(ctx, fmt.Sprintf(`CREATE TAG "%s"`, tagName))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.exec(ctx, fmt.Sprintf(`DROP TAG "%s"`, tagName))
			require.NoError(t, err)
		})

		schemaID := NewSchemaIdentifier(db.Name, randomAccountObjectIdentifier(t).name)
		tagValue := randomString(t)
		err = client.Schemas.Create(ctx, schemaID, &CreateSchemaOptions{
			Tag: []TagAssociation{
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

		err = client.Schemas.Alter(ctx, schemaID, &AlterSchemaOptions{
			Unset: &SchemaUnset{
				Tag: []ObjectIdentifier{
					tagID,
				},
			},
		})
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tagID, schemaID, ObjectTypeSchema)
		require.Error(t, err)
	})

	t.Run("enable manged access", func(t *testing.T) {
		schema, cleanupSchema := createSchema(t, client, db)
		t.Cleanup(cleanupSchema)

		err := client.Schemas.Alter(ctx, schema.ID(), &AlterSchemaOptions{
			EnableManagedAccess: Bool(true),
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
	ctx := context.Background()

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
}

func TestInt_SchemasDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, _ := createSchema(t, client, db)

	s, err := client.Schemas.ShowByID(ctx, schema.ID())
	require.NoError(t, err)
	assert.Equal(t, schema.Name, s.Name)

	err = client.Schemas.Drop(ctx, schema.ID(), nil)
	require.NoError(t, err)

	schemas, err := client.Schemas.Show(ctx, &ShowSchemaOptions{
		Like: &Like{
			Pattern: &schema.Name,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(schemas))
}

func TestInt_SchemasUndrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

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
