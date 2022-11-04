package sdk

import "context"

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

	schemas, err := ts.client.Schemas.List(context.Background(), SchemaListOptions{Pattern: "SCHEMA%"})
	ts.NoError(err)
	ts.Equal(1, len(schemas))

	ts.NoError(ts.client.Schemas.Delete(context.Background(), schema.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}
