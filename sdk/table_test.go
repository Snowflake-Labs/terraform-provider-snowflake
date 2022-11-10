package sdk

import "context"

func (ts *testSuite) createTable(database string) (*Table, error) {
	options := TableCreateOptions{
		Name:         "TABLE_TEST",
		DatabaseName: database,
		Columns: []string{
			"uncollated_phrase varchar",
			"utf8_phrase varchar collate 'utf8'",
			"english_phrase varchar collate 'en'",
			"spanish_phrase varchar collate 'sp'",
		},
		TableProperties: &TableProperties{},
	}
	return ts.client.Tables.Create(context.Background(), options)
}

func (ts *testSuite) TestListTable() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	table, err := ts.createTable(database.Name)
	ts.NoError(err)

	limit := 1
	tables, err := ts.client.Tables.List(context.Background(), TableListOptions{
		Pattern: "TABLE%",
		Limit:   Int(limit),
	})
	ts.NoError(err)
	ts.Equal(limit, len(tables))

	ts.NoError(ts.client.Tables.Delete(context.Background(), table.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestReadTable() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	table, err := ts.createTable(database.Name)
	ts.NoError(err)

	entity, err := ts.client.Tables.Read(context.Background(), table.Name)
	ts.NoError(err)
	ts.Equal(table.Name, entity.Name)

	ts.NoError(ts.client.Tables.Delete(context.Background(), table.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}

func (ts *testSuite) TestCreateTable() {
	database, err := ts.createDatabase()
	ts.NoError(err)
	table, err := ts.createTable(database.Name)
	ts.NoError(err)
	ts.NoError(ts.client.Tables.Delete(context.Background(), table.Name))
	ts.NoError(ts.client.Databases.Delete(context.Background(), database.Name))
}
