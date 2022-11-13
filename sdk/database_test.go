package sdk

import (
	"context"
	"strconv"
)

func (ts *testSuite) createDatabase() (*Database, error) {
	options := DatabaseCreateOptions{
		Name: "DATABASE_TEST",
		DatabaseProperties: &DatabaseProperties{
			Comment:                 String("test database"),
			DataRetentionTimeInDays: Int32(5),
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
			Comment:                 String("updated database"),
			DataRetentionTimeInDays: Int32(10),
		},
	}
	afterUpdate, err := ts.client.Databases.Update(context.Background(), database.Name, options)
	ts.NoError(err)
	ts.Equal(*options.DatabaseProperties.Comment, afterUpdate.Comment)
	retentionTime, err := strconv.Atoi(afterUpdate.RetentionTime)
	ts.NoError(err)
	ts.Equal(*options.DatabaseProperties.DataRetentionTimeInDays, int32(retentionTime))

	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestRenameDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	newDB := "NEW_DATABASE_TEST"
	ts.NoError(ts.client.Databases.Rename(context.Background(), database.Name, newDB))
	ts.NoError(ts.client.Databases.Delete(context.Background(), newDB))
}

func (ts *testSuite) TestCloneDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	destDB := "CLONE_DATABASE_TEST"
	ts.NoError(ts.client.Databases.Clone(context.Background(), database.Name, destDB))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), destDB))
}

func (ts *testSuite) TestUseDatabase() {
	database, err := ts.createDatabase()
	ts.NoError(err)

	ts.NoError(ts.client.Databases.Use(context.Background(), database.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}
