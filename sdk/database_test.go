package sdk

import "context"

func (ts *testSuite) createDatabase() (*Database, error) {
	options := DatabaseCreateOptions{
		Name: "DATABASE_TEST",
		DatabaseProperties: &DatabaseProperties{
			Comment: String("test database"),
		},
	}
	return ts.client.Databases.Create(context.Background(), options)
}

func (ts *testSuite) TestListDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	databases, err := ts.client.Databases.List(context.Background(), DatabaseListOptions{Pattern: "DATABASE%"})
	ts.NoError(err)
	ts.Equal(1, len(databases))

	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestReadDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	entity, err := ts.client.Databases.Read(context.Background(), database.Name)
	ts.NoError(err)
	ts.Equal(database.Name, entity.Name)

	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestCreateDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestUpdateDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	options := DatabaseUpdateOptions{
		DatabaseProperties: &DatabaseProperties{
			Comment: String("updated database"),
		},
	}
	afterUpdate, err := ts.client.Databases.Update(context.Background(), database.Name, options)
	ts.NoError(err)
	ts.Equal(*options.DatabaseProperties.Comment, afterUpdate.Comment)

	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}
