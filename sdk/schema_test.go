package sdk

import (
	"context"
	"strconv"
)

func (ts *testSuite) createSchema(database string) (*Schema, error) {
	options := SchemaCreateOptions{
		Name:         "SCHEMA_TEST",
		DatabaseName: database,
		SchemaProperties: &SchemaProperties{
			Comment: String("test schema"),
		},
	}
	return ts.client.Schemas.Create(context.Background(), options)
}

func (ts *testSuite) TestListSchema() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	schema, err := ts.createSchema(database.Name)
	ts.NoError(err)

	limit := 1
	schemas, err := ts.client.Schemas.List(context.Background(), SchemaListOptions{
		Pattern: "SCHEMA%",
		Limit:   Int(limit),
	})
	ts.NoError(err)
	ts.Equal(limit, len(schemas))

	ts.NoError(ts.client.Schemas.Delete(context.Background(), schema.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestReadSchema() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	schema, err := ts.createSchema(database.Name)
	ts.NoError(err)

	entity, err := ts.client.Schemas.Read(context.Background(), schema.Name)
	ts.NoError(err)
	ts.Equal(schema.Name, entity.Name)

	ts.NoError(ts.client.Schemas.Delete(context.Background(), schema.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestCreateSchema() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	schema, err := ts.createSchema(database.Name)
	ts.NoError(err)
	ts.NoError(ts.client.Schemas.Delete(context.Background(), schema.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestUpdateSchema() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	schema, err := ts.createSchema(database.Name)
	ts.NoError(err)

	options := SchemaUpdateOptions{
		SchemaProperties: &SchemaProperties{
			Comment:                 String("updated schema"),
			DataRetentionTimeInDays: Int32(10),
		},
	}
	afterUpdate, err := ts.client.Schemas.Update(context.Background(), schema.Name, options)
	ts.NoError(err)
	ts.Equal(*options.SchemaProperties.Comment, afterUpdate.Comment)
	retentionTime, err := strconv.Atoi(afterUpdate.RetentionTime)
	ts.NoError(err)
	ts.Equal(*options.SchemaProperties.DataRetentionTimeInDays, int32(retentionTime))

	ts.NoError(ts.client.Schemas.Delete(context.Background(), schema.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestRenameSchema() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	schema, err := ts.createSchema(database.Name)
	ts.NoError(err)

	newSchema := "NEW_SCHEMA_TEST"
	ts.NoError(ts.client.Schemas.Rename(context.Background(), schema.Name, newSchema))
	ts.NoError(ts.client.Schemas.Delete(context.Background(), newSchema))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}
